package asy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yeqown/log"
	"golang.org/x/sync/errgroup"

	"github.com/yeqown/apollo-synchronizer/pkg/apollo/api"
)

var (
	_ Synchronizer = new(synchronizer)
)

type synchronizer struct {
	apollo api.Client

	// scope injected from Synchronize.
	scope *SynchronizeScope
}

func NewSynchronizer(scope *SynchronizeScope) (Synchronizer, error) {
	// permit scope
	log.WithField("scope", scope).Debug("enter synchronizer.Synchronize")
	if scope == nil {
		return nil, errors.New("scope is nil")
	}
	if err := scope.Valid(); err != nil {
		return nil, errors.Wrap(err, "invalid scope")
	}

	return &synchronizer{
		apollo: api.New(&api.Config{
			Token:         scope.ApolloSecret,
			PortalAddress: scope.ApolloPortalAddr,
			Account:       scope.ApolloAccount,
		}),
		scope: scope,
	}, nil
}

// Synchronize scheduling components to display information and execute CURD action with resources.
// NOTICE: properties will be ignored.
func (s *synchronizer) Synchronize(ctx context.Context) ([]*SynchronizeResult, error) {
	scope := s.scope
	log.Info("prepare to fetching remote namespace resources. please wait")
	// load app/env/cluster/remote info
	namespaceInfos, err := s.apollo.ListNamespaces(ctx, scope.ApolloAppID, scope.ApolloEnv, scope.ApolloClusterName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListNamespaces in synchronizer.Synchronize")
	}
	namespaces := make([]string, 0, len(namespaceInfos))
	for _, v := range namespaceInfos {
		if api.NotAllowedFormat(v.Format) {
			// filter properties
			continue
		}

		namespaces = append(namespaces, v.Name)
	}

	files := make([]string, 0, len(scope.LocalFiles))
	for _, v := range scope.LocalFiles {
		ext := strings.TrimPrefix(filepath.Ext(v), ".")
		if api.NotAllowedFormat(api.Format(ext)) {
			// filter unsupported filetypes by apollo
			continue
		}

		files = append(files, filepath.Base(v))
	}

	// compare and display the synchronization information.
	// 1. direction
	// 2. target resources Mode(C/M/D)
	// 3. local and target resources relationship.
	diffs := s.compare(scope.Mode, scope.Path, scope.Force, scope.Overwrite, files, namespaces)
	log.
		WithFields(log.Fields{
			"diffs":      diffs,
			"files":      files,
			"namespaces": namespaces,
		}).
		Info("compare result")
	// let user Decide what to do next, continue or cancel?
	switch decide, reason := s.renderingDiffs(diffs); decide {
	case Decide_CONFIRMED:
	case Decide_CANCELLED:
		fallthrough
	default:
		// interrupt the synchronization
		log.Info("you cancel the synchronization. quit")
		return nil, fmt.Errorf("synchronization cancelled: %s", reason)
	}

	log.Info("synchronizing ...")
	results := s.doSynchronize(scope, diffs)
	log.Info("synchronization finished, please check the result")
	s.renderingResult(results)

	return results, nil
}

// compare calculates the difference between local and remote.
// `force` indicates to create or delete to keep items consistent.
// `overwrite` indicates cover old version while exists
func (s synchronizer) compare(
	mode SynchronizeMode, parent string, force, overwrite bool, localFiles, remoteNamespaces []string) []Diff1 {

	diff0s := make([]Diff0, 0, len(localFiles)+len(remoteNamespaces))
	switch mode {
	case SynchronizeMode_UPLOAD:
		diff0s = compare(localFiles, remoteNamespaces)
	case SynchronizeMode_DOWNLOAD:
		diff0s = compare(remoteNamespaces, localFiles)
	}

	overwriteFilter := func(d Diff0) bool {
		// if not overwrite, skip modify operations
		if overwrite {
			return false
		}
		return d.Mode == DiffMode_MODIFY
	}

	forceFilter := func(d Diff0) bool {
		// if not force, skip create and delete operations
		if force {
			return false
		}
		return d.Mode == DiffMode_CREATE || d.Mode == DiffMode_DELETE
	}

	diff1s := make([]Diff1, 0, len(diff0s))
	for _, d0 := range diff0s {
		if overwriteFilter(d0) || forceFilter(d0) {
			// skip d0
			continue
		}

		diff1s = append(diff1s, Diff1{
			Diff0:       d0,
			AbsFilepath: filepath.Join(parent, d0.Key),
		})
	}

	return diff1s
}

type SynchronizeResult struct {
	Key       string   `json:"key"`
	Mode      diffMode `json:"mode"`
	Error     string   `json:"error"`     // modified failed reason
	Succeeded bool     `json:"succeeded"` // modified Succeeded
	Published bool     `json:"published"` // changes Published
	Bytes     int      `json:"bytes"`     // file size (byte)
}

// doSynchronize execute synchronization between local and remote.
func (s synchronizer) doSynchronize(scope *SynchronizeScope, diffs []Diff1) []*SynchronizeResult {
	log.
		WithFields(log.Fields{
			"Mode":  scope.Mode,
			"diffs": diffs,
		}).
		Debug("doSynchronize")

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	eg, ctx2 := errgroup.WithContext(ctx)

	var (
		resultCh = make(chan *SynchronizeResult, len(diffs))
		done     = make(chan struct{})
		results  = make([]*SynchronizeResult, 0, len(diffs))
	)

	go func() {
		for result := range resultCh {
			results = append(results, result)
		}
		done <- struct{}{}
	}()

	switch scope.Mode {
	case SynchronizeMode_DOWNLOAD:
		for idx := range diffs {
			d := diffs[idx]
			eg.Go(func() error {
				result := s.download(ctx2, d)
				resultCh <- result
				return nil
			})
		}
	case SynchronizeMode_UPLOAD:
		for idx := range diffs {
			d := diffs[idx]
			eg.Go(func() error {
				result := s.upload(ctx2, d, scope.ApolloAutoPublish)
				resultCh <- result
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	close(resultCh)
	<-done

	return results
}

func (s synchronizer) download(ctx context.Context, d Diff1) (r *SynchronizeResult) {
	r = &SynchronizeResult{
		Key:       d.Key,
		Mode:      d.Mode,
		Error:     "",
		Succeeded: false,
		// download always Published by default since it has no version control mechanism.
		Published: true,
		Bytes:     0,
	}
	var err error

	switch d.Mode {
	case DiffMode_DELETE:
		err = os.Remove(d.AbsFilepath)
	case DiffMode_CREATE:
		fallthrough
	case DiffMode_MODIFY:
		item, err2 := s.apollo.GetNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, d.Key, "content")
		if err2 != nil {
			err = err2
			goto Failed
		}
		r.Bytes = len([]byte(item.Value))
		err = os.WriteFile(d.AbsFilepath, []byte(item.Value), 0644)
	}

Failed:
	if err != nil {
		r.Error = err.Error()
		return
	} else {
		r.Succeeded = true
	}

	return
}

func (s synchronizer) upload(ctx context.Context, d Diff1, autoPublish bool) (r *SynchronizeResult) {
	r = &SynchronizeResult{
		Key:       d.Key,
		Mode:      d.Mode,
		Succeeded: false,
		Published: false,
		Error:     "",
	}
	var (
		err error
		ns  = d.Key
	)

	switch d.Mode {
	case DiffMode_DELETE:
		err = s.apollo.DeleteNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, ns, "content")
	case DiffMode_CREATE:
		ext := filepath.Ext(d.Key)                      // .ext .json .yaml
		namespaceName := strings.TrimSuffix(d.Key, ext) // pure filename without ext
		format := strings.TrimPrefix(ext, ".")          // json, yaml

		bytes, err2 := os.ReadFile(d.AbsFilepath)
		if err2 != nil {
			err = err2
			goto Failed
		}
		r.Bytes = len(bytes)
		_, err = s.apollo.CreateNamespace(ctx,
			namespaceName, s.scope.ApolloAppID, api.Format(format), false, "created by apollo-synchronizer")
		if err != nil {
			goto Failed
		}
		_, err = s.apollo.UpdateNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, ns, "content", string(bytes))
		if err != nil {
			goto Failed
		}
	case DiffMode_MODIFY:
		bytes, err2 := os.ReadFile(d.AbsFilepath)
		if err2 != nil {
			err = err2
			goto Failed
		}

		_, err = s.apollo.UpdateNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, ns, "content", string(bytes))
		if err != nil {
			goto Failed
		}
	}

	if autoPublish {
		if _, err2 := s.apollo.PublishNamespace(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, ns); err2 == nil {
			r.Published = true
		} else {
			log.
				WithFields(log.Fields{
					"namespace": ns,
					"app":       s.scope.ApolloAppID,
				}).
				Warnf("publish namespace failed: %v", err)
		}
	}

Failed:
	if err != nil {
		r.Error = err.Error()
		return
	} else {
		r.Succeeded = true
	}

	return
}

func (s synchronizer) renderingDiffs(diffs []Diff1) (d Decide, reason string) {
	if s.scope.Render == nil {
		return Decide_CONFIRMED, "auto confirmed since there is no render"
	}

	return s.scope.Render.RenderingDiffs(diffs)
}

func (s synchronizer) renderingResult(results []*SynchronizeResult) {
	s.scope.Render.RenderingResult(results)
	return
}

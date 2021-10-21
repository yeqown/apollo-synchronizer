package internal

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

	"github.com/yeqown/apollo-synchronizer/internal/apollo/openapi"
)

var (
	_ Synchronizer = new(synchronizer)
)

type synchronizer struct {
	apollo openapi.Client

	// scope injected from Synchronize.
	scope *SynchronizeScope
}

func NewSynchronizer(token, portalAddress, account string) Synchronizer {
	return &synchronizer{
		apollo: openapi.New(&openapi.Config{
			Token:         token,
			PortalAddress: portalAddress,
			Account:       account,
		}),
		scope: nil,
	}
}

// Synchronize scheduling components to display information and execute CURD action with resources.
// NOTICE: properties will be ignored.
func (s *synchronizer) Synchronize(ctx context.Context, scope *SynchronizeScope) error {
	// permit scope
	log.
		WithFields(log.Fields{
			"scope": scope,
		}).
		Debug("enter synchronizer.Synchronize")
	s.scope = scope

	// load app/env/cluster/remote info
	namespaceInfos, err := s.apollo.ListNamespaces(ctx, scope.ApolloAppID, scope.ApolloEnv, scope.ApolloClusterName)
	if err != nil {
		return errors.Wrap(err, "failed to ListNamespaces in synchronizer.Synchronize")
	}
	namespaces := make([]string, 0, len(namespaceInfos))
	for _, v := range namespaceInfos {
		if openapi.NotAllowedFormat(v.Format) {
			// filter properties
			continue
		}

		namespaces = append(namespaces, v.Name)
	}

	files := make([]string, 0, len(scope.LocalFiles))
	for _, v := range scope.LocalFiles {
		if openapi.NotAllowedFormat(openapi.Format(filepath.Ext(v))) {
			// filter unsupported filetypes by apollo
			continue
		}

		files = append(files, filepath.Base(v))
	}

	// compare and display the synchronization information.
	// 1. direction
	// 2. target resources mode(C/M/D)
	// 3. local and target resources relationship.
	diffs := s.compare(scope.Mode, scope.Path, scope.Force, scope.Overwrite, files, namespaces)
	userDecide := s.renderDiff(diffs)

	switch userDecide {
	case Decide_CONFIRMED:
	case Decide_CANCELLED:
		fallthrough
	default:
		return nil
	}

	syncResults := s.doSynchronize(scope, diffs)
	s.renderSynchronizeResult(syncResults)
	return nil
}

// compare calculates the difference between local and remote.
// `force` indicates to create or delete to keep items consistent.
// `overwrite` indicates cover old version while exists
func (s synchronizer) compare(
	mode SynchronizeMode, parent string, force, overwrite bool, localFiles, remoteNamespaces []string) []diff1 {

	diff0s := make([]diff0, 0, len(localFiles)+len(remoteNamespaces))
	switch mode {
	case SynchronizeMode_UPLOAD:
		diff0s = compare(localFiles, remoteNamespaces)
	case SynchronizeMode_DOWNLOAD:
		diff0s = compare(remoteNamespaces, localFiles)
	}

	overwriteFilter := func(d diff0) bool {
		// if not overwrite, skip modify operations
		if overwrite {
			return false
		}
		return d.mode == diffMode_MODIFY
	}

	forceFilter := func(d diff0) bool {
		// if not force, skip create and delete operations
		if force {
			return false
		}
		return d.mode == diffMode_CREATE || d.mode == diffMode_DELETE
	}

	diff1s := make([]diff1, 0, len(diff0s))
	for _, d0 := range diff0s {
		if overwriteFilter(d0) || forceFilter(d0) {
			// skip d0
			continue
		}

		diff1s = append(diff1s, diff1{
			diff0:       d0,
			absFilepath: filepath.Join(parent, d0.key),
		})
	}

	return diff1s
}

type synchronizeResult struct {
	key       string
	mode      diffMode
	error     string // modified failed reason
	succeeded bool   // modified succeeded
	published bool   // changes published
}

// doSynchronize execute synchronization between local and remote.
func (s synchronizer) doSynchronize(scope *SynchronizeScope, diffs []diff1) []*synchronizeResult {
	log.
		WithFields(log.Fields{
			"mode":  scope.Mode,
			"diffs": diffs,
		}).
		Debug("doSynchronize")

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	eg, ctx2 := errgroup.WithContext(ctx)

	var (
		resultCh = make(chan *synchronizeResult, len(diffs))
		done     = make(chan struct{})
		results  = make([]*synchronizeResult, 0, len(diffs))
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

func (s synchronizer) download(ctx context.Context, d diff1) (r *synchronizeResult) {
	r = &synchronizeResult{
		key:       d.key,
		mode:      d.mode,
		error:     "",
		succeeded: false,
		// download always published by default since it has no version control mechanism.
		published: true,
	}
	var err error

	switch d.mode {
	case diffMode_DELETE:
		err = os.Remove(d.absFilepath)
	case diffMode_CREATE:
		fallthrough
	case diffMode_MODIFY:
		item, err2 := s.apollo.GetNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, d.key, "content")
		if err2 != nil {
			err = err2
			goto Failed
		}
		err = os.WriteFile(d.absFilepath, []byte(item.Value), 0644)
	}

Failed:
	if err != nil {
		r.error = err.Error()
		return
	} else {
		r.succeeded = true
	}

	return
}

func (s synchronizer) upload(ctx context.Context, d diff1, autoPublish bool) (r *synchronizeResult) {
	r = &synchronizeResult{
		key:       d.key,
		mode:      d.mode,
		succeeded: false,
		published: false,
		error:     "",
	}
	var (
		err error
		ns  = d.key
	)

	switch d.mode {
	case diffMode_DELETE:
		err = s.apollo.DeleteNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, ns, "content")
	case diffMode_CREATE:
		ext := filepath.Ext(d.key)                      // .ext .json .yaml
		namespaceName := strings.TrimSuffix(d.key, ext) // pure filename without ext
		format := strings.TrimPrefix(ext, ".")          // json, yaml

		bytes, err2 := os.ReadFile(d.absFilepath)
		if err2 != nil {
			err = err2
			goto Failed
		}
		_, err = s.apollo.CreateNamespace(ctx,
			namespaceName, s.scope.ApolloAppID, openapi.Format(format), false, "created by apollo-synchronizer")
		if err != nil {
			goto Failed
		}
		_, err = s.apollo.UpdateNamespaceItem(
			ctx, s.scope.ApolloAppID, s.scope.ApolloEnv, s.scope.ApolloClusterName, ns, "content", string(bytes))
		if err != nil {
			goto Failed
		}
	case diffMode_MODIFY:
		bytes, err2 := os.ReadFile(d.absFilepath)
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
			r.published = true
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
		r.error = err.Error()
		return
	} else {
		r.succeeded = true
	}

	return
}

// decide confirm synchronize or cancel.
type decide uint8

const (
	Decide_UNKNOWN decide = iota
	Decide_CONFIRMED
	Decide_CANCELLED
)

func (s synchronizer) renderDiff(diffs []diff1) decide {
	return Decide_CONFIRMED
}

func (s synchronizer) renderSynchronizeResult(results []*synchronizeResult) {
	for _, r := range results {
		if r.succeeded {
			fmt.Printf("mode=%s, key=%s, success=%v, published=%v\n", r.mode, r.key, r.succeeded, r.published)
		} else {
			fmt.Printf("mode=%s, key=%s, failed=%s\n, published=%v", r.mode, r.key, r.error, r.published)
		}
	}
}

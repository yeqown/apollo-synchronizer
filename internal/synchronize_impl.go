package internal

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/yeqown/log"

	"github.com/yeqown/apollo-synchronizer/internal/apollo/openapi"
)

var (
	_ Synchronizer = synchronizer{}
)

type synchronizer struct {
	apollo openapi.Client
}

func NewSynchronizer(token, portalAddress, account string) Synchronizer {
	return synchronizer{
		apollo: openapi.New(&openapi.Config{
			Token:         token,
			PortalAddress: portalAddress,
			Account:       account,
		}),
	}
}

// Synchronize scheduling components to display information and execute CURD action with resources.
// NOTICE: properties will be ignored.
func (s synchronizer) Synchronize(ctx context.Context, scope *SynchronizeScope) error {
	// permit scope
	log.
		WithFields(log.Fields{
			"scope": scope,
		}).
		Debug("enter synchronizer.Synchronize")

	// load app/env/cluster/namespaces info
	namespaceInfos, err := s.apollo.ListNamespaces(ctx, scope.ApolloAppID, scope.ApolloEnv, scope.ApolloClusterName)
	if err != nil {
		return errors.Wrap(err, "failed to ListNamespaces in synchronizer.Synchronize")
	}
	namespaces := make([]string, len(namespaceInfos))
	for idx, v := range namespaceInfos {
		if v.Format == openapi.Format_Properties {
			continue
		}

		namespaces[idx] = v.Name
	}

	files := make([]string, len(scope.LocalFiles))
	for idx, v := range scope.LocalFiles {
		if filepath.Ext(v) == string(openapi.Format_Properties) {
			continue
		}
		files[idx] = filepath.Base(v)
	}

	// compare and display the synchronization information.
	// 1. direction
	// 2. target resources mode(C/M/D)
	// 3. local and target resources relationship.
	diffs := s.compare(scope.Mode, scope.Path, files, namespaces)
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

type diff struct {
	key         string
	absFilepath string
	mode        diffMode
}

type diffMode string

const (
	diffMode_CREATE diffMode = "C"
	diffMode_MODIFY diffMode = "M"
	diffMode_DELETE diffMode = "D"
)

// compare calculates the difference between local and remote.
func (s synchronizer) compare(
	mode SynchronizeMode, parent string, localFiles []string, remoteNamespaces []string) []diff {

	localM := make(map[string]struct{}, len(localFiles))
	for _, f := range localFiles {
		localM[f] = struct{}{}
	}

	remoteM := make(map[string]struct{}, len(remoteNamespaces))
	for _, ns := range remoteNamespaces {
		remoteM[ns] = struct{}{}
	}

	diffs := make([]diff, 0, len(localM)+len(remoteM))
	for key := range remoteM {
		_, ok := localM[key]
		d := diff{
			key:         key,
			absFilepath: filepath.Join(parent, key),
			mode:        diffMode_MODIFY,
		}

		if !ok {
			switch mode {
			case SynchronizeMode_DOWNLOAD:
				d.mode = diffMode_CREATE
			case SynchronizeMode_UPLOAD:
				d.mode = diffMode_DELETE
			}
		}

		diffs = append(diffs, d)
	}
	for key := range localM {
		_, ok := remoteM[key]
		if ok {
			continue
		}

		d := diff{
			key:         key,
			absFilepath: filepath.Join(parent, key),
			mode:        diffMode_DELETE,
		}

		if !ok {
			switch mode {
			case SynchronizeMode_DOWNLOAD:
				d.mode = diffMode_DELETE
			case SynchronizeMode_UPLOAD:
				d.mode = diffMode_CREATE
			}
		}

		diffs = append(diffs, d)
	}

	return diffs
}

// doSynchronize execute synchronization between local and remote.
func (s synchronizer) doSynchronize(scope *SynchronizeScope, diffs []diff) []string {
	log.
		WithFields(log.Fields{
			"mode":  scope.Mode,
			"diffs": diffs,
		}).
		Debug("doSynchronize")

	switch scope.Mode {
	case SynchronizeMode_DOWNLOAD:
	case SynchronizeMode_UPLOAD:
	default:
		panic("invalid mode: " + strconv.Itoa(int(scope.Mode)))
	}

	return []string{"TODO"}
}

// decide confirm synchronize or cancel.
type decide uint8

const (
	Decide_UNKNOWN decide = iota
	Decide_CONFIRMED
	Decide_CANCELLED
)

func (s synchronizer) renderDiff(diffs []diff) decide {
	return Decide_CONFIRMED
}

func (s synchronizer) renderSynchronizeResult([]string) {
}

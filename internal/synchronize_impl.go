package internal

import (
	"context"
	"strconv"

	"github.com/yeqown/log"

	"github.com/pkg/errors"

	"github.com/yeqown/apollo-synchronizer/internal/apollo/openapi"
)

var (
	_ Synchronizer = synchronizer{}
)

type synchronizer struct {
	apollo openapi.Client
}

func NewSynchronizer() Synchronizer {
	return synchronizer{}
}

// Synchronize scheduling components to display information and execute CURD action with resources.
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
		namespaces[idx] = v.Name + "." + v.Format
	}

	// compare and display the synchronization information.
	// 1. direction
	// 2. target resources mode(M/N/D)
	// 3. local and target resources relationship.
	diffs := s.compare(scope.LocalFiles, namespaces)
	userDecide := s.renderDiff(diffs)

	switch userDecide {
	case Decide_CONFIRMED:
	case Decide_CANCELLED:
		fallthrough
	default:
		return nil
	}

	syncResults := s.doSynchronize(scope, scope.LocalFiles, namespaces)
	s.renderSynchronizeResult(syncResults)
	return nil
}

type diff struct {
	file      string
	namespace string
	mode      diffMode
}

type diffMode string

const (
	diffMode_CREATE diffMode = "N"
	diffMode_MODIFY diffMode = "M"
	diffMode_DELETE diffMode = "D"
)

func (s synchronizer) compare(localFiles []string, namespaces []string) []diff {
	return []diff{
		{
			file:      "",
			namespace: "",
			mode:      "",
		},
	}
}

// doSynchronize execute synchronization between local and remote.
func (s synchronizer) doSynchronize(scope *SynchronizeScope, localFiles []string, namespaces []string) []string {
	switch scope.Mode {
	case DOWNLOAD:
	case UPLOAD:
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
	return Decide_CANCELLED
}

func (s synchronizer) renderSynchronizeResult([]string) {
}

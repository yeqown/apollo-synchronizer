package internal

import (
	"context"
	"fmt"
)

// Synchronizer 's duty is synchronizing between remote apollo portal and local filesystem.
type Synchronizer interface {
	Synchronize(ctx context.Context, syncCtx *SynchronizeScope) error
}

type SynchronizeMode uint8

const (
	UNKNOWN SynchronizeMode = iota
	UPLOAD
	DOWNLOAD
)

// SynchronizeScope represents the context of Synchronizer.Synchronize.
type SynchronizeScope struct {
	Mode       SynchronizeMode
	LocalFiles []string

	ApolloSecret      string
	ApolloAppID       string
	ApolloEnv         string
	ApolloClusterName string

	// Overwrite indicates whether asy update the target while it exists.
	Overwrite bool
	// Force indicates whether to create the target while it not exists.
	Force bool
}

func (sc *SynchronizeScope) valid() error {
	return fmt.Errorf("TODO this")
}

// FillSynchronizeContext to fill context in order:
// 1. load config file (.asy) from current directory to fill apollo config.
// 2. parse from command line parameters.
func FillSynchronizeContext(scope *SynchronizeScope) error {
	return fmt.Errorf("TODO this")
}

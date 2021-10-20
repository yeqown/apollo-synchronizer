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
	SynchronizeMode_UNKNOWN SynchronizeMode = iota
	SynchronizeMode_UPLOAD
	SynchronizeMode_DOWNLOAD
)

// SynchronizeScope represents the context of Synchronizer.Synchronize.
type SynchronizeScope struct {
	Mode SynchronizeMode
	// Path is the parent directory which holds all downloaded remote.
	// apollo.appid will be used as a sub-directory. [DOWNLOAD REQUIRED]
	Path string
	// LocalFiles represents the absolute file path of local files. [UPLOAD ONLY]
	LocalFiles []string

	ApolloSecret      string
	ApolloAppID       string
	ApolloEnv         string
	ApolloClusterName string
	ApolloPortalAddr  string
	ApolloAccount     string

	// Overwrite indicates whether asy update the target while it exists.
	Overwrite bool
	// Force indicates whether to create the target while it not exists.
	Force bool
}

func (sc SynchronizeScope) Valid() error {
	if sc.ApolloSecret == "" {
		return fmt.Errorf("ApolloSecret could not be empty")
	}
	if sc.ApolloAppID == "" {
		return fmt.Errorf("ApolloAppID could not be empty")
	}
	if sc.ApolloEnv == "" {
		return fmt.Errorf("ApolloEnv could not be empty")
	}
	if sc.ApolloClusterName == "" {
		return fmt.Errorf("ApolloClusterName could not be empty")
	}
	if sc.ApolloPortalAddr == "" {
		return fmt.Errorf("ApolloPortalAddr could not be empty")
	}
	if sc.ApolloAccount == "" {
		return fmt.Errorf("ApolloAccount could not be empty")
	}

	switch sc.Mode {
	case SynchronizeMode_UPLOAD:
	case SynchronizeMode_DOWNLOAD:
		if sc.Path == "" {
			return fmt.Errorf("path can not be empty")
		}
	case SynchronizeMode_UNKNOWN:
		fallthrough
	default:
		return fmt.Errorf("you can only specify upload or download")
	}

	return nil
}

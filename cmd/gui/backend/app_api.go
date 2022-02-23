package backend

import (
	"context"
	"path/filepath"
	"time"

	asy "github.com/yeqown/apollo-synchronizer"
)

// var (
// 	synchronizerCache = make(map[string]asy.Synchronizer) // map[portalHash]synchronizer
// )

// func (b *App) getSynchronizer(portalHash string) asy.Synchronizer {
// 	if s, ok := synchronizerCache[portalHash]; ok {
// 		return s
// 	}

// 	b.infof("getSynchronizer: portalHash=%d, clusterSettings: %+v", portalHash, b.Clusters)
// 	if b.Clusters == nil || portalHash == "" {
// 		return nil
// 	}

// 	scope := b.Clusters[portalHash]
// 	s := asy.NewSynchronizer(scope)
// 	synchronizerCache[portalHash] = s

// 	return s
// }

type SynchronizeScope struct {
}

func (b *App) Synchronize(scope2 *SynchronizeScope) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	scope := asy.SynchronizeScope{
		ApolloPortalAddr:  "",
		ApolloAccount:     "",
		Path:              "",
		LocalFiles:        nil,
		Mode:              0,
		ApolloSecret:      "",
		ApolloAppID:       "",
		ApolloEnv:         "",
		ApolloClusterName: "",
		ApolloAutoPublish: false,
		Overwrite:         false,
		Force:             false,
		Render:            eventsRender{},
	}

	s, err := asy.NewSynchronizer(&scope)
	if err != nil {
		b.errorf("NewSynchronizer: %+v", err)
		return err
	}

	if err := s.Synchronize(ctx); err != nil {
		b.errorf("synchronize failed: %v", err)
		return err
	}

	return nil
}

func (b *App) LoadSetting() []apolloClusterSetting {
	return b.config.Settings
}

func (b *App) SaveSetting(settings []apolloClusterSetting) {
	b.config.Settings = settings
	save(filepath.Join(appConfigRoot(), "asyrc"), b.config.Bytes(), false)
}

func (b *App) Statistics() statistics {
	return *b.statistics
}

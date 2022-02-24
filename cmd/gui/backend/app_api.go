package backend

import (
	"context"
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
	ApolloPortalAddr  string              `json:"portalAddr"`
	ApolloAccount     string              `json:"account"`
	Path              string              `json:"fs"`
	Mode              asy.SynchronizeMode `json:"mode"`
	ApolloSecret      string              `json:"secret"`
	ApolloAppID       string              `json:"appId"`
	ApolloEnv         string              `json:"env"`
	ApolloClusterName string              `json:"cluster"`
	ApolloAutoPublish bool                `json:"isAutoPublish"`
	Overwrite         bool                `json:"isOverwrite"`
	Force             bool                `json:"isForce"`
}

func (b *App) Synchronize(scope2 *SynchronizeScope) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	b.debugf("App.Synchronize called, scope: %+v\n", scope2)
	return nil

	scope := asy.SynchronizeScope{
		ApolloPortalAddr:  scope2.ApolloPortalAddr,
		ApolloAccount:     scope2.ApolloAccount,
		Path:              "",  // TODO(@yeqown)
		LocalFiles:        nil, // TODO(@yeqown)
		Mode:              scope2.Mode,
		ApolloSecret:      scope2.ApolloSecret,
		ApolloAppID:       scope2.ApolloAppID,
		ApolloEnv:         scope2.ApolloEnv,
		ApolloClusterName: scope2.ApolloClusterName,
		ApolloAutoPublish: scope2.ApolloAutoPublish,
		Overwrite:         scope2.Overwrite,
		Force:             scope2.Force,
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
	save(_configFp, b.config, _ext_json)
}

func (b *App) Statistics() statistics {
	b.debugf("App.Statistics called, statistics: %+v\n", b.statistics)
	return *b.statistics
}

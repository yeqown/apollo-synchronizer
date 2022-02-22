package main

import (
	"context"
	"time"

	"github.com/yeqown/apollo-synchronizer/internal"
)

var (
	synchronizerCache = make(map[uint8]internal.Synchronizer)
)

func (b *App) getSynchronizer(clusterIdx uint8) internal.Synchronizer {
	if s, ok := synchronizerCache[clusterIdx]; ok {
		return s
	}

	b.infof("getSynchronizer: clusterIdx=%d, clusterSettings: %+v", clusterIdx, b.Clusters)
	if b.Clusters == nil || int(clusterIdx) > len(b.Clusters) {
		return nil
	}

	setting := b.Clusters[clusterIdx]
	s := internal.NewSynchronizer(setting.Secret, setting.PortalAddress, setting.Account)
	synchronizerCache[clusterIdx] = s

	return s
}

func (b *App) Synchronize() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	scope := internal.SynchronizeScope{
		Mode:              0,
		Path:              "",
		LocalFiles:        nil,
		ApolloSecret:      "",
		ApolloAppID:       "",
		ApolloEnv:         "",
		ApolloClusterName: "",
		ApolloPortalAddr:  "",
		ApolloAccount:     "",
		ApolloAutoPublish: false,
		Overwrite:         false,
		Force:             false,
		EnableTermUI:      false,
	}

	idx := uint8(0)
	s := b.getSynchronizer(idx)
	if s == nil {
		b.debugf("synchronizer is nil: idx=%d", idx)
		return
	}

	if err := s.Synchronize(ctx, &scope); err != nil {
		b.errorf("synchronize failed: %v", err)
	}
}

func (b *App) LoadSetting() []clusterSetting {
	return []clusterSetting{
		{
			Title:         "setting1",
			Secret:        "ebba7e6efa4bb04479eb38464c0e7afc65",
			Clusters:      []string{"default", "preprod"},
			Env:           "DEV",
			PortalAddress: "http://localhost:8080",
			Account:       "apollo",
			LocalDir:      "/Users/jia/.asy/setting1-DEV-$portalHash6",
		},
	}
}

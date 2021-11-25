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
	s := internal.NewSynchronizer(setting.Token, setting.PortalAddress, setting.Account)
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

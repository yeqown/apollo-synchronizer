package backend

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	asy "github.com/yeqown/apollo-synchronizer"
	"github.com/yeqown/apollo-synchronizer/pkg/fs"
)

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

type synchronizeResult struct {
	Succeeded    bool   `json:"succeeded"`
	FailedReason string `json:"failedReason"`
}

func (r *synchronizeResult) markSuccess() {
	r.Succeeded = true
	r.FailedReason = ""
}

func (r *synchronizeResult) markFailure(err error) {
	r.Succeeded = false
	r.FailedReason = "internal error"
	if err != nil {
		r.FailedReason = err.Error()
	}
}

func (b *App) Synchronize(param *SynchronizeScope) (result *synchronizeResult) {
	result = new(synchronizeResult)
	defer func() {
		switch param.Mode {
		case asy.SynchronizeMode_DOWNLOAD:
			b.statistics.DownloadCount++
			if !result.Succeeded {
				b.statistics.DownloadFailedCount++
			}
		case asy.SynchronizeMode_UPLOAD:
			b.statistics.UploadCount++
			if !result.Succeeded {
				b.statistics.UploadFailedCount++
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	b.debugf("App.Synchronize called, param: %+v\n", param)
	// For debug
	// return nil

	apppath := filepath.Join(param.Path,
		fmt.Sprintf("%s-%s-%s", param.ApolloAppID, param.ApolloEnv, param.ApolloClusterName))

	if err := fs.MakeSure(apppath); err != nil {
		result.markFailure(err)
		return result
	}
	localFiles, err := fs.TravelDirectory(apppath, false)
	if err != nil {
		result.markFailure(err)
		return result
	}

	scope := asy.SynchronizeScope{
		ApolloPortalAddr:  param.ApolloPortalAddr,
		ApolloAccount:     param.ApolloAccount,
		Path:              apppath,
		LocalFiles:        localFiles,
		Mode:              param.Mode,
		ApolloSecret:      param.ApolloSecret,
		ApolloAppID:       param.ApolloAppID,
		ApolloEnv:         param.ApolloEnv,
		ApolloClusterName: param.ApolloClusterName,
		ApolloAutoPublish: param.ApolloAutoPublish,
		Overwrite:         param.Overwrite,
		Force:             param.Force,
		Render:            newRender(b),
	}
	b.debugf("App.Synchronize called, scope: %+v\n", scope)

	s, err := asy.NewSynchronizer(&scope)
	if err != nil {
		b.errorf("build synchronizer failed: %+v", err)
		result.markFailure(err)
		return result
	}

	results, err := s.Synchronize(ctx)
	if err != nil {
		b.errorf("synchronize failed: %v", err)
		result.markFailure(err)
		return result
	}

	// file and size statistics
	{
		switch param.Mode {
		case asy.SynchronizeMode_UPLOAD:
			b.statistics.UploadFileCount += int64(len(results))
			for _, v := range results {
				b.statistics.UploadFileSize += int64(v.Bytes)
			}
		case asy.SynchronizeMode_DOWNLOAD:
			b.statistics.DownloadFileCount += int64(len(results))
			for _, v := range results {
				b.statistics.DownloadFileSize += int64(v.Bytes)
			}
		}
	}

	result.markSuccess()
	return result
}

func (b *App) LoadSetting() []apolloClusterSetting {
	return b.config.Settings
}

func (b *App) SaveSetting(settings []apolloClusterSetting) error {
	b.config.Settings = settings
	return save(_configFp, b.config, _ext_json)
}

func (b *App) Statistics() statistics {
	b.debugf("App.Statistics called, statistics: %+v\n", b.statistics)
	return *b.statistics
}

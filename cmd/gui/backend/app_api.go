package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	asy "github.com/yeqown/apollo-synchronizer"
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

	if err := prepareAppPath(apppath); err != nil {
		result.markFailure(err)
		return result
	}
	localFiles := travelDirectory(apppath, false)

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

func prepareAppPath(apppath string) error {
	fi, err := os.Stat(apppath)
	if err == nil && !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", apppath)
	}

	if err != nil {
		if !os.IsNotExist(err) {

			return fmt.Errorf("%s stat failed", apppath)
		}

		if err = os.MkdirAll(apppath, 0755); err != nil {
			return fmt.Errorf("create directory(%s) failed: %v", apppath, err)
		}
	}

	return nil
}

func travelDirectory(root string, recursive bool) []string {
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("failed to travelDirectory: %v\n", err)
	}

	out := make([]string, 0, len(files))
	for _, fp := range files {
		if fp.IsDir() {
			continue
		}

		out = append(out, filepath.Join(root, fp.Name()))
	}

	return out
}

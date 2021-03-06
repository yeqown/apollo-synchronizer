package main

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/yeqown/log"

	asy "github.com/yeqown/apollo-synchronizer"
	"github.com/yeqown/apollo-synchronizer/pkg/fs"
	"github.com/yeqown/apollo-synchronizer/pkg/recommend"
)

// fillSynchronizeScope to fill scope in order:
// 1. load config file (.asy) from current directory to fill apollo config.
// 2. parse from command line parameters.
func fillSynchronizeScope(ctx *cli.Context) (scope *asy.SynchronizeScope) {
	scope = new(asy.SynchronizeScope)
	cwd, _ := os.Getwd()
	tryFromFile(cwd, scope)
	tryFromContext(ctx, scope)

	var err error
	if !path.IsAbs(scope.Path) {
		scope.Path, _ = filepath.Abs(scope.Path)
	}
	scope.Path = recommend.GenerateAppPath(scope.Path, scope.ApolloAppID, scope.ApolloEnv, scope.ApolloClusterName)
	if scope.Mode == asy.SynchronizeMode_DOWNLOAD {
		if err = fs.MakeSure(scope.Path); err != nil {
			log.Fatalf("make sure path %s failed: %v", scope.Path, err)
		}
	}
	scope.LocalFiles, err = fs.TravelDirectory(scope.Path, false)
	if err != nil {
		log.Fatalf("travel directory %s failed: %v", scope.Path, err)
	}

	//switch scope.Mode {
	//case pkg.SynchronizeMode_UPLOAD:
	//	if scope.Path == "" {
	//		scope.LocalFiles = ctx.StringSlice("file")
	//	} else {
	//		scope.Path = path.Join(scope.Path, scope.ApolloAppID)
	//		scope.LocalFiles = travelDirectory(scope.Path, false)
	//	}
	//case pkg.SynchronizeMode_DOWNLOAD:
	//	// use path only
	//	scope.Path = path.Join(scope.Path, scope.ApolloAppID)
	//	scope.LocalFiles = travelDirectory(scope.Path, false)
	//}

	for idx, f := range scope.LocalFiles {
		if scope.LocalFiles[idx], err = filepath.Abs(f); err != nil {
			log.Fatal("stat file failed: %s", f)
		}
	}

	log.
		WithFields(log.Fields{
			"cwd":   cwd,
			"scope": scope,
		}).
		Debug("fillSynchronizeScope")

	return
}

func tryFromFile(cwd string, scope *asy.SynchronizeScope) {
	fp := path.Join(cwd, ".asy.json")
	_, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.
			WithFields(log.Fields{
				"cwd":   cwd,
				"fp":    fp,
				"error": err,
			}).
			Error("tryFromFile failed to stat file")
		return
	}

	fd, err := os.OpenFile(fp, os.O_RDONLY, 0644)
	if err != nil {
		log.Error("tryFromFile failed to open file: %v", err)
		return
	}

	data, err := io.ReadAll(fd)
	if err != nil {
		log.Errorf("tryFromFile failed to read from file: %v", err)
		return
	}

	var apollo = struct {
		ApolloSecret      string `json:"apolloSecret"`
		ApolloAppID       string `json:"apolloAppId"`
		ApolloEnv         string `json:"apolloEnv"`
		ApolloClusterName string `json:"apolloCluster"`
		ApolloPortalAddr  string `json:"apolloPortalAddr"`
		ApolloAccount     string `json:"apolloAccount"`
	}{}

	if err = json.Unmarshal(data, &apollo); err != nil {
		log.Errorf("tryFromFile failed to unmarshal: %v", err)
		return
	}

	// copy
	scope.ApolloSecret = apollo.ApolloSecret
	scope.ApolloAppID = apollo.ApolloAppID
	scope.ApolloEnv = apollo.ApolloEnv
	scope.ApolloClusterName = apollo.ApolloClusterName
	scope.ApolloPortalAddr = apollo.ApolloPortalAddr
	scope.ApolloAccount = apollo.ApolloAccount

	return
}

func tryFromContext(ctx *cli.Context, scope *asy.SynchronizeScope) {
	// mode
	scope.Mode = asy.SynchronizeMode_UNKNOWN
	if ctx.Bool("down") {
		scope.Mode = asy.SynchronizeMode_DOWNLOAD
	}
	if ctx.Bool("up") {
		scope.Mode = asy.SynchronizeMode_UPLOAD
	}

	scope.Force = ctx.Bool("force")
	scope.Overwrite = ctx.Bool("overwrite")

	// apollo api parameter
	scope.ApolloSecret = ctx.String("apollo.secret")
	scope.ApolloAppID = ctx.String("apollo.appid")
	scope.ApolloEnv = ctx.String("apollo.env")
	scope.ApolloClusterName = ctx.String("apollo.cluster")
	scope.ApolloPortalAddr = ctx.String("apollo.portaladdr")
	scope.ApolloAccount = ctx.String("apollo.account")
	scope.ApolloAutoPublish = ctx.Bool("auto-publish")

	scope.Render = newTerminalUI()
	if ctx.Bool("enable-termui") {
		scope.Render = newTermUI(scope)
	}

	// local filesystem
	scope.Path = ctx.String("path")
}

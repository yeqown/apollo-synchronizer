package main

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/yeqown/apollo-synchronizer/internal"

	"github.com/urfave/cli/v2"
	"github.com/yeqown/log"
)

// fillSynchronizeScope to fill scope in order:
// 1. load config file (.asy) from current directory to fill apollo config.
// 2. parse from command line parameters.
func fillSynchronizeScope(ctx *cli.Context) (scope *internal.SynchronizeScope) {
	scope = new(internal.SynchronizeScope)
	cwd, _ := os.Getwd()
	tryFromFile(cwd, scope)
	tryFromContext(ctx, scope)

	log.
		WithFields(log.Fields{
			"cwd":   cwd,
			"scope": scope,
		}).
		Debug("fillSynchronizeScope")

	return
}

func tryFromFile(cwd string, scope *internal.SynchronizeScope) {
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

func tryFromContext(ctx *cli.Context, scope *internal.SynchronizeScope) {
	// mode
	scope.Mode = internal.SynchronizeMode_UNKNOWN
	if ctx.Bool("up") {
		scope.Mode = internal.SynchronizeMode_UPLOAD
	}
	if ctx.Bool("down") {
		scope.Mode = internal.SynchronizeMode_DOWNLOAD
	}

	scope.Force = ctx.Bool("force")
	scope.Overwrite = ctx.Bool("overwrite")

	// apollo openapi parameter
	scope.ApolloSecret = ctx.String("apollo.secret")
	scope.ApolloAppID = ctx.String("apollo.appid")
	scope.ApolloEnv = ctx.String("apollo.env")
	scope.ApolloClusterName = ctx.String("apollo.cluster")
	scope.ApolloPortalAddr = ctx.String("apollo.portaladdr")
	scope.ApolloAccount = ctx.String("apollo.account")

	// local filesystem
	scope.Path = ctx.String("path")
	switch scope.Mode {
	case internal.SynchronizeMode_UPLOAD:
		if scope.Path == "" {
			scope.LocalFiles = ctx.StringSlice("file")
		} else {
			scope.Path = path.Join(scope.Path, scope.ApolloAppID)
			scope.LocalFiles = travelDirectory(scope.Path, false)
		}
	case internal.SynchronizeMode_DOWNLOAD:
		// use path only
		scope.Path = path.Join(scope.Path, scope.ApolloAppID)
	}
}

func travelDirectory(root string, recursive bool) []string {
	files, err := os.ReadDir(root)
	if err != nil {
		log.Fatal("failed to travelDirectory: %v", err)
	}

	out := make([]string, 0, len(files))
	for _, fp := range files {
		if fp.IsDir() {
			continue
		}

		out = append(out, fp.Name())
	}

	return out
}
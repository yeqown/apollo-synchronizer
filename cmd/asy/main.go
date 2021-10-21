package main

import (
	"context"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
	"github.com/yeqown/log"

	"github.com/yeqown/apollo-synchronizer/internal"
)

func main() {
	app := cli.NewApp()
	app.Name = "apollo-synchronizer"
	app.Authors = []*cli.Author{
		{
			Name:  "yeqown",
			Email: "yeqown@gmail.com",
		},
	}
	app.Description = "To help developers synchronize between apollo portal and local filesystem."
	app.Version = "v1.2.0"
	app.Flags = flags
	app.Action = action
	app.Before = before
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func before(c *cli.Context) error {
	log.SetLogLevel(log.LevelInfo)
	if c.Bool("debug") {
		log.
			SetLogLevel(log.LevelDebug)
	}

	return nil
}

func action(c *cli.Context) error {
	scope := fillSynchronizeScope(c)

	if err := scope.Valid(); err != nil {
		return errors.Wrap(err, "scope invalid")
	}

	return internal.
		NewSynchronizer(scope.ApolloSecret, scope.ApolloPortalAddr, scope.ApolloAccount).
		Synchronize(context.Background(), scope)
}

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "up",
		Usage: "upload to apollo portal with local filesystem",
	},
	&cli.BoolFlag{
		Name:        "down",
		DefaultText: "true",
		Value:       true,
		Usage:       "download from apollo portal",
	},
	&cli.BoolFlag{
		Name:    "force",
		Aliases: []string{"f"},
		Usage:   "indicates whether to create the target while it not exists.",
	},
	&cli.BoolFlag{
		Name:        "overwrite",
		Usage:       "indicates whether asy update the target while it exists.",
		DefaultText: "true",
		Value:       true,
	},
	&cli.BoolFlag{
		Name:  "enable-termui",
		Usage: "use terminal ui to display and interact with instead of logs",
	},
	&cli.StringFlag{
		Name:        "path",
		Usage:       "specify the path to synchronize",
		TakesFile:   false,
		Value:       defaultPath(),
		DefaultText: defaultPath(),
	},
	//&cli.StringSliceFlag{
	//	Name:      "file",
	//	Usage:     "specify files to synchronize",
	//	FilePath:  ".",
	//	TakesFile: true,
	//},
	&cli.StringFlag{
		Name:  "apollo.portaladdr",
		Usage: "apollo portal address",
	},
	&cli.StringFlag{
		Name:  "apollo.appid",
		Usage: "the targeted remote app in apollo",
	},
	&cli.StringFlag{
		Name:  "apollo.secret",
		Usage: "openapi app's token",
	},
	&cli.StringFlag{
		Name:        "apollo.account",
		DefaultText: "apollo",
		Value:       "apollo",
		Usage:       "user id in apollo",
	},
	&cli.StringFlag{
		Name:        "apollo.env",
		DefaultText: "DEV",
		Value:       "DEV",
		Usage:       "the environment of target remote app",
	},
	&cli.StringFlag{
		Name:        "apollo.cluster",
		DefaultText: "default",
		Value:       "default",
		Usage:       "the cluster of target remote app",
	},
	&cli.BoolFlag{
		Name:        "auto-publish",
		DefaultText: "false",
		Usage:       "enable auto publish apollo modified namespace.",
	},
	&cli.BoolFlag{
		Name:        "debug",
		Usage:       "print debug logs",
		DefaultText: "false",
		Value:       false,
	},
}

var commands = []*cli.Command{
	genToolCommand(),
}

func defaultPath() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, ".apollo-synchronizer")
}

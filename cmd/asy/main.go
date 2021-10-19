package main

import (
	"github.com/urfave/cli/v2"
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
	app.Version = "v1.0.0"
	app.Flags = flags
	app.Action = action
}

func action(c *cli.Context) error {
	panic("implement me")
}

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "up",
		Usage: "upload to apollo portal with local filesystem",
	},
	&cli.BoolFlag{
		Name:  "down",
		Usage: "download from apollo portal",
	},
}

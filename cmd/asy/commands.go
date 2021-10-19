package main

import "github.com/urfave/cli/v2"

// genToolCommand to help developers create, delete and read resources from apollo portal.
func genToolCommand() *cli.Command {
	return &cli.Command{
		Name:                   "tool",
		Aliases:                nil,
		Usage:                  "To help developers create, delete and read resources from apollo portal.",
		UsageText:              "",
		Description:            "",
		ArgsUsage:              "",
		Category:               "",
		BashComplete:           nil,
		Before:                 nil,
		After:                  nil,
		Action:                 nil,
		OnUsageError:           nil,
		Subcommands:            nil,
		Flags:                  nil,
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
}

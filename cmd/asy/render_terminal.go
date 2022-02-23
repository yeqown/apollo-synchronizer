package main

import (
	"fmt"

	asy "github.com/yeqown/apollo-synchronizer"
)

var _ asy.Renderer = (*terminalRenderer)(nil)

type terminalRenderer struct{}

func newTerminalUI() asy.Renderer {
	return &terminalRenderer{}
}

func (t terminalRenderer) RenderingDiffs(diffs []asy.Diff1) asy.Decide {
	fmt.Printf("=================== synchronize differences ===================\n")

	display := func(d asy.Diff1) {
		modeText := ""
		switch d.Mode {
		case asy.DiffMode_DELETE:
			modeText = red(string(d.Mode))
		case asy.DiffMode_MODIFY:
			modeText = yellow(string(d.Mode))
		case asy.DiffMode_CREATE:
			modeText = green(string(d.Mode))
		}

		fmt.Printf("%s %15s\t%s\n", modeText, d.Key, d.AbsFilepath)
	}

	for _, d := range diffs {
		display(d)
	}

	return asy.Decide_CONFIRMED
}

func (t terminalRenderer) RenderingResult(results []*asy.SynchronizeResult) {
	fmt.Printf("=================== synchronization results ===================\n")

	display := func(r asy.SynchronizeResult) {
		modeText := ""
		switch r.Mode {
		case asy.DiffMode_DELETE:
			modeText = red(string(r.Mode))
		case asy.DiffMode_MODIFY:
			modeText = yellow(string(r.Mode))
		case asy.DiffMode_CREATE:
			modeText = green(string(r.Mode))
		}
		uploadResultText := red("FAILED ")
		publishResultText := gray("UNDONE")
		ext := yellow("REASON: " + r.Error)
		if r.Succeeded {
			uploadResultText = green("SUCCESS")
			publishResultText = green("DONE   ")
			ext = ""
		}

		fmt.Printf("%s %15s\tSyncStatus: %s\tPubStatus: %s\t%s\n",
			modeText, r.Key, uploadResultText, publishResultText, ext)
	}

	for _, r := range results {
		display(*r)
	}
}

func red(text string) string {
	return colored(color_RED, text)
}

func green(text string) string {
	return colored(color_GREEN, text)
}

func yellow(text string) string {
	return colored(color_YELLOW, text)
}

func gray(text string) string {
	return colored(color_GRAY, text)
}

type color string

const (
	color_RED    = "0;31"
	color_GREEN  = "0;32"
	color_YELLOW = "0;33"
	color_GRAY   = "0;37"
)

func colored(c color, text string) string {
	return "\033[" + string(c) + "m" + text + "\033[0m"
}

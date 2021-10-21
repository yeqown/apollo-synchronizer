package internal

import (
	"fmt"
)

type terminalRenderer struct{}

func (t terminalRenderer) renderingDiffs(diffs []diff1) decide {
	fmt.Printf("=================== synchronize differences ===================\n")

	display := func(d diff1) {
		modeText := ""
		switch d.mode {
		case diffMode_DELETE:
			modeText = red(string(d.mode))
		case diffMode_MODIFY:
			modeText = yellow(string(d.mode))
		case diffMode_CREATE:
			modeText = green(string(d.mode))
		}

		fmt.Printf("%s %15s\t%s\n", modeText, d.key, d.absFilepath)
	}

	for _, d := range diffs {
		display(d)
	}

	return Decide_CONFIRMED
}

func (t terminalRenderer) renderingResult(results []*synchronizeResult) {
	fmt.Printf("=================== synchronization results ===================\n")

	display := func(r synchronizeResult) {
		modeText := ""
		switch r.mode {
		case diffMode_DELETE:
			modeText = red(string(r.mode))
		case diffMode_MODIFY:
			modeText = yellow(string(r.mode))
		case diffMode_CREATE:
			modeText = green(string(r.mode))
		}
		uploadResultText := red("FAILED ")
		publishResultText := gray("UNDONE")
		ext := yellow("REASON: " + r.error)
		if r.succeeded {
			uploadResultText = green("SUCCESS")
			publishResultText = green("DONE   ")
			ext = ""
		}

		fmt.Printf("%s %15s\tSyncStatus: %s\tPubStatus: %s\t%s\n",
			modeText, r.key, uploadResultText, publishResultText, ext)
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

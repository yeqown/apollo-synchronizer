package internal

import (
	"errors"
	"fmt"
	"image"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/yeqown/log"
)

type termuiRenderer struct {
	scope *SynchronizeScope
}

func newTermUI(scope *SynchronizeScope) termuiRenderer {
	return termuiRenderer{
		scope: scope,
	}
}

func (t termuiRenderer) initUI() {
	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
}

var (
	errEnter = errors.New("enter key")
)

// wait any key break view loop
func (t termuiRenderer) wait() (err error) {
	defer ui.Close()
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<Enter>" {
				return errEnter
			}
			break
		}
	}

	return nil
}

func (t termuiRenderer) renderingDiffs(diffs []diff1) decide {
	t.initUI()
	x, y := ui.TerminalDimensions()
	p := widgets.NewParagraph()

	layout := `        AppID: %s
          Env: %s
       luster: %s
PortalAddress: %s
  AutoPublish: %v
        Force: %v
    Overwrite: %v`

	p.Text = fmt.Sprintf(layout, t.scope.ApolloAppID, t.scope.ApolloEnv, t.scope.ApolloClusterName,
		t.scope.ApolloPortalAddr, t.scope.ApolloAutoPublish, t.scope.Force, t.scope.Overwrite,
	)
	p.Border = true
	p.Title = "synchronize settings"
	p.TitleStyle = ui.NewStyle(ui.ColorCyan)
	p.SetRect(0, 0, x, y/4+1)

	l1 := widgets.NewList()
	l1.Title = "Local"
	l1.TitleStyle = ui.NewStyle(ui.ColorMagenta)
	l2 := widgets.NewList()
	l2.Title = "Remote"
	l2.TitleStyle = ui.NewStyle(ui.ColorBlue)
	l3 := widgets.NewList()
	l3.TitleStyle = ui.NewStyle(ui.ColorCyan)
	l3.Title = "Absolute Path"

	dire := widgets.NewImage(directionImage(x/6+1, 2*y/4+1, 2*x/6+1, 3*y/4+1, t.scope.Mode))
	dire.Border = false

	for _, d := range diffs {
		l3.Rows = append(l3.Rows, d.absFilepath)
		switch d.mode {
		case diffMode_CREATE:
			l1.Rows = append(l1.Rows, d.key)
			l2.Rows = append(l2.Rows, "-")
		case diffMode_DELETE:
			l1.Rows = append(l1.Rows, "-")
			l2.Rows = append(l2.Rows, d.key)
		case diffMode_MODIFY:
			l1.Rows = append(l1.Rows, d.key)
			l2.Rows = append(l2.Rows, d.key)
		}
	}

	l1.SetRect(0, y/4+1, x/6+1, y)
	dire.SetRect(x/6+1, 2*y/4+1, 2*x/6+1, 3*y/4+1)
	l2.SetRect(2*x/6+1, y/4+1, 3*x/6+1, y)
	l3.SetRect(3*x/6+1, y/4+1, x, y)

	ui.Render(p, l1, dire, l2, l3)
	if err := t.wait(); err != nil {
		if errors.Is(err, errEnter) {
			return Decide_CONFIRMED
		}
	}

	return Decide_CANCELLED
}

func (t termuiRenderer) renderingResult(results []*synchronizeResult) {
	t.initUI()
	x, y := ui.TerminalDimensions()

	tb := widgets.NewTable()
	tb.Title = "synchronization Result"
	tb.TitleStyle = ui.NewStyle(ui.ColorCyan)
	tb.Border = true
	tb.RowSeparator = true
	tb.TextAlignment = ui.AlignCenter
	tb.FillRow = true
	tb.Rows = [][]string{
		{"Mode", "Key", "Synchronize Status", "Publish Status", "Reason"},
	}

	for idx, r := range results {
		row := []string{
			string(r.mode),
			r.key,
			strconv.FormatBool(r.succeeded),
			strconv.FormatBool(r.published),
			r.error,
		}
		tb.Rows = append(tb.Rows, row)

		tb.RowStyles[idx+1] = ui.NewStyle(ui.ColorGreen)
		if !r.succeeded {
			tb.RowStyles[idx+1] = ui.NewStyle(ui.ColorWhite, ui.ColorRed)
		}
	}
	tb.SetRect(0, 0, x, y)

	ui.Render(tb)
	_ = t.wait()
}

// draw arrow direction (left arrow, right arrow) image
func directionImage(x1, y1, x2, y2 int, mode SynchronizeMode) image.Image {
	img := image.NewGray(image.Rect(x1, y1, x2, y2))
	// TODO(@yeqown) draw the arrow
	return img
}

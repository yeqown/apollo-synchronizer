package internal

import (
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
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
		log.WithField("event", e).Debug("an event received")
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<Enter>" {
				return errEnter
			}
			if e.ID == "<C-c>" {
				ui.Close()
				os.Exit(0)
				return
			}
			break
		}
	}

	return nil
}

func (t termuiRenderer) renderingDiffs(diffs []diff1) decide {
	t.initUI()
	x, y := ui.TerminalDimensions()

	p0 := widgets.NewParagraph()
	layout := `        AppID: %s
          Env: %s
       luster: %s
PortalAddress: %s
  AutoPublish: %v
        Force: %v
    Overwrite: %v`
	p0.Text = fmt.Sprintf(layout, t.scope.ApolloAppID, t.scope.ApolloEnv, t.scope.ApolloClusterName,
		t.scope.ApolloPortalAddr, t.scope.ApolloAutoPublish, t.scope.Force, t.scope.Overwrite,
	)
	p0.Border = true
	p0.Title = "synchronize settings"
	p0.TitleStyle = ui.NewStyle(ui.ColorCyan)
	p0.SetRect(0, 0, x/2+1, y/4+1)

	p1 := widgets.NewParagraph()
	p1.Text = "Press [Enter] to continue, [^C] to quit, other keys to cancel"
	p1.Title = "Tips"
	p1.TitleStyle = ui.NewStyle(ui.ColorCyan)
	//p1.TextStyle =
	p1.SetRect(x/2+1, 0, x, y/4+1)

	l1 := widgets.NewList()
	l1.Title = "Local"
	l1.TitleStyle = ui.NewStyle(ui.ColorCyan)
	l1.TextStyle = ui.NewStyle(ui.ColorCyan)
	l1.SelectedRow = 9

	l2 := widgets.NewList()
	l2.Title = "Remote"
	l2.TitleStyle = ui.NewStyle(ui.ColorCyan)
	l2.TextStyle = ui.NewStyle(ui.ColorCyan)
	l2.SelectedRow = 9

	l3 := widgets.NewList()
	l3.Title = "Absolute Path"
	l3.TitleStyle = ui.NewStyle(ui.ColorCyan)
	l3.TextStyle = ui.NewStyle(ui.ColorWhite)
	l3.SelectedRow = 9

	//img := widgets.NewImage(directionImage(0, 0, x/6-2, y/4-1, t.scope.Mode))
	//img.Border = false

	p2 := widgets.NewParagraph()
	p2.Border = false
	switch t.scope.Mode {
	case SynchronizeMode_DOWNLOAD:
		p2.Text = `      ⬅️⬅️`

	case SynchronizeMode_UPLOAD:
		p2.Text = "      ➡️➡️"
	}

	for idx, d := range diffs {
		row1, row2 := d.key, d.key
		color1, color2 := "yellow", "yellow"

		switch d.mode {
		case diffMode_CREATE:
			row2 = "✖✖✖"
			color1 = "green"
			color2 = "green"
		case diffMode_DELETE:
			row1 = "✖✖✖"
			color1 = "red"
			color2 = "red"
		case diffMode_MODIFY:
		}

		row1 = fmt.Sprintf("[%d] [%s %s](fg:%s)", idx+1, d.mode, row1, color1)
		row2 = fmt.Sprintf("[%d] [%s %s](fg:%s)", idx+1, d.mode, row2, color2)
		l1.Rows = append(l1.Rows, row1)
		l2.Rows = append(l2.Rows, row2)
		l3.Rows = append(l3.Rows, d.absFilepath)
	}

	l1.SetRect(0, y/4+1, x/6+1, y)
	p2.SetRect(x/6+1, 2*y/4+1, 2*x/6+1, 3*y/4+1)
	l2.SetRect(2*x/6+1, y/4+1, 3*x/6+1, y)
	l3.SetRect(3*x/6+1, y/4+1, x, y)

	ui.Render(p0, p1, l1, p2, l2, l3)
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
	tb.ColumnWidths = []int{6, 20, 20, 20, x - 66}
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
			tb.RowStyles[idx+1] = ui.NewStyle(ui.ColorRed, ui.ColorBlack, ui.ModifierBold)
		}
	}
	tb.SetRect(0, 0, x, y)

	ui.Render(tb)
	_ = t.wait()
}

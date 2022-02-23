package backend

import (
	asy "github.com/yeqown/apollo-synchronizer"
)

var _ asy.Renderer = (*eventsRender)(nil)

const (
	eventNameRenderDiff   = "event.render.diff"
	eventNameRenderResult = "event.render.result"
)

func newRender(app *App) asy.Renderer {
	return eventsRender{
		app: app,
	}
}

// eventsRender implements the asy.Renderer interface based on the wails EventsOn.
//
// https://wails.io/zh-Hans/docs/reference/runtime/events
type eventsRender struct {
	app *App
}

func (r eventsRender) RenderingDiffs(diffs []asy.Diff1) asy.Decide {
	r.app.infof("RenderingDiffs: %+v", diffs)
	return asy.Decide_CONFIRMED
}

func (r eventsRender) RenderingResult(results []*asy.SynchronizeResult) {
	r.app.infof("RenderingResult: %+v", results)
	return
}

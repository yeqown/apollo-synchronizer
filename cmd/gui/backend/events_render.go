package backend

import (
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	asy "github.com/yeqown/apollo-synchronizer"
)

var _ asy.Renderer = (*eventsRender)(nil)

const (
	eventNameRenderDiff   = "event.render.diff"
	eventNameRenderResult = "event.render.result"

	// eventNameInputDecide means the user input the decide whether to continue
	// or not. Frontend will emit this event to backend.
	eventNameInputDecide = "event.input.decide"
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

func (r eventsRender) RenderingDiffs(diffs []asy.Diff1) (d asy.Decide, reason string) {
	r.app.infof("RenderingDiffs: %+v", diffs)

	decideChan := make(chan asy.Decide, 1)
	timeout := time.NewTimer(30 * time.Second)

	runtime.EventsOnce(r.app.ctx, eventNameInputDecide, func(optionalData ...interface{}) {
		r.app.debugf("eventNameInputDecide called: %+v", optionalData)
		// DONE(@yeqown): get actual data from optionalData
		if len(optionalData) <= 0 {
			r.app.debugf("eventNameInputDecide called with no data")
			decideChan <- asy.Decide_CANCELLED
			return
		}

		v := optionalData[0].([]interface{})
		// fmt.Printf("%t\n", v[0])
		decide, ok2 := v[0].(float64)
		if !ok2 || asy.Decide(int(decide)) == asy.Decide_UNKNOWN {
			r.app.debugf("eventNameInputDecide called with invalid data 2")
			decideChan <- asy.Decide_CANCELLED
			return
		}

		decideChan <- asy.Decide(int(decide))
	})
	// trigger frontend to render diffs
	runtime.EventsEmit(r.app.ctx, eventNameRenderDiff, diffs)

	r.app.infof("Waiting for decide...")
	select {
	case d = <-decideChan:
		r.app.infof("Input decide: %+v", d)
		reason = "user decided"
	case <-timeout.C:
		r.app.errorf("Timeout waiting for decide")
		d = asy.Decide_CANCELLED
		reason = "timeout"
	}
	r.app.infof("final decide: %+v, reason: %s", d, reason)

	return d, reason
}

func (r eventsRender) RenderingResult(results []*asy.SynchronizeResult) {
	r.app.infof("RenderingResult: %+v", results)
	runtime.EventsEmit(r.app.ctx, eventNameRenderResult, results)

	return
}

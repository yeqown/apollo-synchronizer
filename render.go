package asy

// Decide confirm synchronize or cancel.
type Decide uint8

const (
	Decide_UNKNOWN Decide = iota
	Decide_CONFIRMED
	Decide_CANCELLED
)

// Renderer declares the display component which display some information.
type Renderer interface {
	// RenderingDiffs display the diffs and make a decision to control the synchronization
	// should go on or abort now.
	RenderingDiffs(diffs []Diff1) Decide

	// RenderingResult display the result of synchronization.
	RenderingResult(results []*SynchronizeResult)
}

package internal

// decide confirm synchronize or cancel.
type decide uint8

const (
	Decide_UNKNOWN decide = iota
	Decide_CONFIRMED
	Decide_CANCELLED
)

// renderer declares the display component which display some information.
type renderer interface {
	renderingDiffs(diffs []diff1) decide
	renderingResult(results []*synchronizeResult)
}

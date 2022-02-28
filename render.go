package asy

// Decide confirm synchronize or cancel.
type Decide uint8

const (
	Decide_UNKNOWN Decide = iota
	Decide_CONFIRMED
	Decide_CANCELLED
)

func (d Decide) String() string {
	switch d {
	case Decide_CONFIRMED:
		return "CONFIRMED"
	case Decide_CANCELLED:
		return "CANCELLED"
	}
	return "UNKNOWN"
}

// Renderer declares the display component which display some information.
type Renderer interface {
	// RenderingDiffs display the diffs and make a decision to control the synchronization
	// should go on or abort now.
	RenderingDiffs(diffs []Diff1) (d Decide, reason string)

	// RenderingResult display the result of synchronization.
	RenderingResult(results []*SynchronizeResult)
}

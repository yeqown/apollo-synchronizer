package internal

type diffMode string

const (
	diffMode_CREATE diffMode = "C"
	diffMode_MODIFY diffMode = "M"
	diffMode_DELETE diffMode = "D"
)

type diff0 struct {
	key  string
	mode diffMode
}

type diff1 struct {
	diff0

	absFilepath string
}

type sorterDiff1 []diff1

func (d sorterDiff1) Len() int {
	return len(d)
}

func (d sorterDiff1) Less(i, j int) bool {
	return d[i].key < d[j].key
}

func (d sorterDiff1) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// compare from and to, calculate the diffs between from and to obey the rule:
// 1. key exists in `from` and `to` get 'M' mode
// 2. key exists in `from` only get 'C' mode
// 3. key exists in `to` only get 'D' mode.
func compare(from []string, to []string) (diffs []diff0) {
	fromMapping := make(map[string]struct{}, len(from))
	for _, f := range from {
		fromMapping[f] = struct{}{}
	}

	toMapping := make(map[string]struct{}, len(to))
	for _, ns := range to {
		toMapping[ns] = struct{}{}
	}

	diffs = make([]diff0, 0, len(fromMapping)+len(toMapping))
	for k1 := range fromMapping {
		d := diff0{
			key:  k1,
			mode: diffMode_MODIFY,
		}
		if _, ok := toMapping[k1]; !ok {
			// exists in `from` only
			d.mode = diffMode_CREATE
		}
		diffs = append(diffs, d)
	}
	for k2 := range toMapping {
		if _, ok := fromMapping[k2]; ok {
			// has been set already.
			continue
		}
		// exists in `to` only
		d := diff0{
			key:  k2,
			mode: diffMode_DELETE,
		}
		diffs = append(diffs, d)
	}

	return diffs
}

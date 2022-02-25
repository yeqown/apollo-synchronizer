package asy

type diffMode string

const (
	DiffMode_CREATE diffMode = "C+"
	DiffMode_MODIFY diffMode = "M~"
	DiffMode_DELETE diffMode = "D-"
)

type Diff0 struct {
	Key  string   `json:"key"`
	Mode diffMode `json:"mode"`
}

type Diff1 struct {
	Diff0

	AbsFilepath string `json:"absFilepath"`
}

type sorterDiff1 []Diff1

func (d sorterDiff1) Len() int {
	return len(d)
}

func (d sorterDiff1) Less(i, j int) bool {
	return d[i].Key < d[j].Key
}

func (d sorterDiff1) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// compare from and to, calculate the diffs between from and to obey the rule:
// 1. Key exists in `from` and `to` get 'M' Mode
// 2. Key exists in `from` only get 'C' Mode
// 3. Key exists in `to` only get 'D' Mode.
func compare(from []string, to []string) (diffs []Diff0) {
	fromMapping := make(map[string]struct{}, len(from))
	for _, f := range from {
		fromMapping[f] = struct{}{}
	}

	toMapping := make(map[string]struct{}, len(to))
	for _, ns := range to {
		toMapping[ns] = struct{}{}
	}

	diffs = make([]Diff0, 0, len(fromMapping)+len(toMapping))
	for k1 := range fromMapping {
		d := Diff0{
			Key:  k1,
			Mode: DiffMode_MODIFY,
		}
		if _, ok := toMapping[k1]; !ok {
			// exists in `from` only
			d.Mode = DiffMode_CREATE
		}
		diffs = append(diffs, d)
	}
	for k2 := range toMapping {
		if _, ok := fromMapping[k2]; ok {
			// has been set already.
			continue
		}
		// exists in `to` only
		d := Diff0{
			Key:  k2,
			Mode: DiffMode_DELETE,
		}
		diffs = append(diffs, d)
	}

	return diffs
}

package gitdiff

type patienceDiffer struct{}

func (p patienceDiffer) ComputeDiff(previous Document, current Document, contextLinesCount int) ([]Diff, error) {
	panic("TODO")
}

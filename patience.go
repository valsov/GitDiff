package gitdiff

type patienceDiffer struct{}

func (p patienceDiffer) ComputeDiff(previous Document, current Document) ([]Diff, error) {
	panic("TODO")
}

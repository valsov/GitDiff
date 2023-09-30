package gitdiff

import "strings"

const (
	UNCHANGED DiffType = "unchanged"
	ADDED     DiffType = "added"
	DELETED   DiffType = "deleted"
)

type DiffType string

type Differ interface {
	ComputeDiff(previous, current Document) ([]Diff, error)
}

type Diff struct {
	Old, New Line
	Type     DiffType
}

func NewDiff(old, new Line, dType DiffType) Diff {
	return Diff{
		Old:  old,
		New:  new,
		Type: dType,
	}
}

func (d Diff) GetText() string {
	if d.New.lineNumber != -1 {
		return d.New.text
	}
	return d.Old.text
}

type Document []Line

type Line struct {
	text       string
	lineNumber int
}

func NewDocument(text string) Document {
	d := Document{}
	split := strings.Split(text, "\n")
	for i, t := range split {
		d = append(d, Line{text: t, lineNumber: i})
	}
	return d
}

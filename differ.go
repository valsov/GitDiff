package gitdiff

import (
	"errors"
	"strings"
)

const (
	MYERS_DIFF    DiffAlgorithm = "myers"
	PATIENCE_DIFF DiffAlgorithm = "patience"

	UNCHANGED DiffType = "unchanged"
	ADDED     DiffType = "added"
	DELETED   DiffType = "deleted"
)

type DiffAlgorithm string

type DiffType string

type Differ interface {
	ComputeDiff(previous, current Document) ([]Diff, error)
}

func NewDiffProducer(algorithm DiffAlgorithm) (Differ, error) {
	switch algorithm {
	case MYERS_DIFF:
		return myersDiffer{}, nil
	case PATIENCE_DIFF:
		return patienceDiffer{}, nil
	default:
		return nil, errors.New("unknown diff algorithm")
	}
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

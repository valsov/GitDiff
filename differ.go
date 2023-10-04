package gitdiff

import (
	"errors"
	"strings"
)

const (
	MYERS_DIFF DiffAlgorithm = "myers"
)

const (
	UNCHANGED DiffType = "unchanged"
	ADDED     DiffType = "added"
	DELETED   DiffType = "deleted"
)

type DiffAlgorithm string

type DiffType string

type Differ interface {
	// Produce a slice of Diff by comparing input documents.
	// 'contextLinesCount' indicates the expected number of lines surrounding edits. If negative, all lines will be included.
	ComputeDiff(previous, current Document, contextLinesCount int) ([]Diff, error)
}

func NewDiffProducer(algorithm DiffAlgorithm) (Differ, error) {
	switch algorithm {
	case MYERS_DIFF:
		return myersDiffer{}, nil
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
	if text == "" {
		return d
	}
	split := strings.Split(text, "\n")
	for i, t := range split {
		d = append(d, Line{text: t, lineNumber: i + 1})
	}
	return d
}

func handleSimpleCases(previous, current Document) ([]Diff, bool) {
	if len(previous) == 0 {
		// All added
		diffs := make([]Diff, len(current))
		for i, line := range current {
			diffs[i] = NewDiff(Line{lineNumber: -1}, line, ADDED)
		}
		return diffs, true
	}
	if len(current) == 0 {
		// All deleted
		diffs := make([]Diff, len(previous))
		for i, line := range previous {
			diffs[i] = NewDiff(line, Line{lineNumber: -1}, DELETED)
		}
		return diffs, true
	}
	return nil, false
}

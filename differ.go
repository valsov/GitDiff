package gitdiff

import "strings"

type Differ interface {
	ComputeDiff(previous, current Document) ([]Diff, error)
}

type Diff struct {
	OldLineNumber, NewLineNumber int
	Text                         string
}

type Line struct {
	text       string
	lineNumber int
}

type Document []Line

func NewDocument(text string) Document {
	d := Document{}
	split := strings.Split(text, "\n")
	for i, t := range split {
		d = append(d, Line{text: t, lineNumber: i})
	}
	return d
}

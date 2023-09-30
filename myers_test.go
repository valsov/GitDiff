package gitdiff

import "testing"

func TestComputeDiff(t *testing.T) {
	previousText := `A
B
C
A
B
B
A`
	previous := NewDocument(previousText)

	currentText := `C
B
A
B
A
C`
	current := NewDocument(currentText)

	alg := MyersDiffer{}
	diffs, err := alg.ComputeDiff(previous, current)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if len(diffs) != 9 {
		t.Fatalf("wrong diffs element count. expected=%d, got=%d", 9, len(diffs))
	}

	for i, diffType := range []DiffType{
		DELETED,
		DELETED,
		UNCHANGED,
		ADDED,
		UNCHANGED,
		UNCHANGED,
		DELETED,
		UNCHANGED,
		ADDED,
	} {
		if diffs[i].Type != diffType {
			t.Fatalf("wrong diffs type at index %d. expected=%s, got=%s", i, diffType, diffs[i].Type)
		}
	}
}

func TestShortestEdition(t *testing.T) {
	previousText := `A
B
C
A
B
B
A`
	previous := NewDocument(previousText)

	currentText := `C
B
A
B
A
C`
	current := NewDocument(currentText)

	alg := MyersDiffer{}
	traces, err := alg.shortestEdition(previous, current)

	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if len(traces) != 6 {
		t.Fatalf("wrong diff count. expected=%d, got=%d", 6, len(traces))
	}
}

func TestBacktrack(t *testing.T) {
	previousText := `A
B
C
A
B
B
A`
	previous := NewDocument(previousText)

	currentText := `C
B
A
B
A
C`
	current := NewDocument(currentText)

	alg := MyersDiffer{}
	traces, err := alg.shortestEdition(previous, current)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	path := alg.backtrack(previous, current, traces)
	if len(path) != 9 {
		t.Fatalf("wrong path length. expected=%d, got=%d", 9, len(path))
	}
	if path[0].x2 != 7 || path[0].y2 != 6 {
		t.Fatalf("wrong path end. expected=[%d, %d], got=[%d, %d]", 7, 6, path[0].x2, path[0].y2)
	}
	previousX, previousY := path[0].x1, path[0].y1
	for i := 1; i < len(path); i++ {
		if path[i].x2 != previousX || path[i].y2 != previousY {
			t.Fatalf("path is not contiguous between index %d and %d. expected=[%d, %d], got=[%d, %d]",
				i-1, i, previousX, previousY, path[i].x2, path[i].y2)
		}
		previousX, previousY = path[i].x1, path[i].y1
	}
}

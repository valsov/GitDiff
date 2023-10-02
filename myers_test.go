package gitdiff

import "testing"

func TestComputeDiff(t *testing.T) {
	for _, tc := range []struct {
		previous, current string
		contextLinesCount int
		expectedDiffs     []DiffType
	}{
		{
			previous:          "",
			current:           "",
			contextLinesCount: 0,
			expectedDiffs:     []DiffType{},
		},
		{
			previous:          "A\nB\nC",
			current:           "",
			contextLinesCount: 0,
			expectedDiffs:     []DiffType{DELETED, DELETED, DELETED},
		},
		{
			previous:          "",
			current:           "A\nB\nC",
			contextLinesCount: 0,
			expectedDiffs:     []DiffType{ADDED, ADDED, ADDED},
		},
		{
			previous:          "A\nB\nC",
			current:           "A\nB\nC",
			contextLinesCount: 0,
			expectedDiffs:     []DiffType{},
		},
		{
			previous:          "A\nA\nB\nC",
			current:           "A\nA\nD\nC",
			contextLinesCount: 0,
			expectedDiffs:     []DiffType{DELETED, ADDED},
		},
		{
			previous:          "A\nA\nB\nC",
			current:           "A\nA\nD\nC",
			contextLinesCount: 1, // With 1 line of context
			expectedDiffs:     []DiffType{UNCHANGED, DELETED, ADDED, UNCHANGED},
		},
		{
			previous:          "A\nA\nB\nC",
			current:           "A\nA\nD\nC",
			contextLinesCount: -1, // All lines
			expectedDiffs:     []DiffType{UNCHANGED, UNCHANGED, DELETED, ADDED, UNCHANGED},
		},
	} {
		previous, current := NewDocument(tc.previous), NewDocument(tc.current)
		alg := myersDiffer{}
		diffs, err := alg.ComputeDiff(previous, current, tc.contextLinesCount)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if len(diffs) != len(tc.expectedDiffs) {
			t.Fatalf("wrong diffs element count. expected=%d, got=%d", len(tc.expectedDiffs), len(diffs))
		}

		for i, diffType := range tc.expectedDiffs {
			if diffs[i].Type != diffType {
				t.Fatalf("wrong diffs type at index %d. expected=%s, got=%s", i, diffType, diffs[i].Type)
			}
		}
	}
}

func TestShortestEdition(t *testing.T) {
	previousText := `DELETED
UNCHANGED
UNCHANGED`
	previous := NewDocument(previousText)

	currentText := `ADDED
UNCHANGED
UNCHANGED
ADDED`
	current := NewDocument(currentText)

	alg := myersDiffer{}
	traces, err := alg.shortestEdition(previous, current)

	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if len(traces) != 4 {
		t.Fatalf("wrong diff count. expected=%d, got=%d", 4, len(traces))
	}
}

func TestBacktrack(t *testing.T) {
	previousText := `DELETED
UNCHANGED
UNCHANGED`
	previous := NewDocument(previousText)

	currentText := `ADDED
UNCHANGED
UNCHANGED
ADDED`
	current := NewDocument(currentText)

	alg := myersDiffer{}
	traces, err := alg.shortestEdition(previous, current)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	path := alg.backtrack(previous, current, traces)
	if len(path) != 5 {
		t.Fatalf("wrong path length. expected=%d, got=%d", 5, len(path))
	}
	if path[0].x2 != 3 || path[0].y2 != 4 {
		t.Fatalf("wrong path end. expected=[%d, %d], got=[%d, %d]", 3, 4, path[0].x2, path[0].y2)
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

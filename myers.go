package gitdiff

import (
	"errors"
	"math"
	"slices"
)

type myersDiffer struct{}

type positionsPair struct {
	x1, y1 int
	x2, y2 int
}

func (m myersDiffer) ComputeDiff(previous, current Document, contextLinesCount int) ([]Diff, error) {
	if len(previous) == 0 {
		// All added
		diffs := make([]Diff, len(current))
		for i, line := range current {
			diffs[i] = NewDiff(Line{lineNumber: -1}, line, ADDED)
		}
		return diffs, nil
	}
	if len(current) == 0 {
		// All deleted
		diffs := make([]Diff, len(current))
		for i, line := range current {
			diffs[i] = NewDiff(line, Line{lineNumber: -1}, DELETED)
		}
		return diffs, nil
	}

	traces, err := m.shortestEdition(previous, current)
	if err != nil {
		return nil, err
	}

	path := m.backtrack(previous, current, traces)
	slices.Reverse(path)

	// Build diffs slice
	var (
		diffs          = []Diff{}
		tempDiffBuffer = []Diff{}
		ahead          = false
		insertCount    = 0
	)

	if contextLinesCount < 0 {
		contextLinesCount = math.MaxInt
	}

	for _, pos := range path {
		var diff Diff
		switch {
		case pos.x1 == pos.x2:
			diff = NewDiff(Line{lineNumber: -1}, current[pos.y1], ADDED)
		case pos.y1 == pos.y2:
			diff = NewDiff(previous[pos.x1], Line{lineNumber: -1}, DELETED)
		default:
			diff = NewDiff(previous[pos.x1], current[pos.y1], UNCHANGED)
		}

		if diff.Type != UNCHANGED {
			for j := len(tempDiffBuffer) - contextLinesCount; j < len(tempDiffBuffer) && j >= 0; j++ {
				diffs = append(diffs, tempDiffBuffer[j])
			}
			tempDiffBuffer = nil // Reset slice and set len to 0
			insertCount = 0
			ahead = true

			diffs = append(diffs, diff)
			continue
		}

		if ahead {
			if insertCount == contextLinesCount {
				ahead = false
				insertCount = 0
				tempDiffBuffer = append(tempDiffBuffer, diff)
			} else {
				diffs = append(diffs, diff)
				insertCount++
			}
		} else {
			tempDiffBuffer = append(tempDiffBuffer, diff)
		}
	}

	return diffs, nil
}

// Find shortest edition, produce a list of traces that led to the result
func (m myersDiffer) shortestEdition(previous, current Document) ([][]int, error) {
	maxDepth := max(len(previous), len(current))
	if maxDepth == 0 {
		return [][]int{}, nil
	}

	arr := make([]int, maxDepth*2+1)
	trace := [][]int{}
	for depth := 0; depth <= maxDepth; depth++ {
		for k := -depth; k <= depth; k += 2 {
			// Move by one
			var x int
			if k == -depth || (k != depth && getAtIndex(arr, k-1) < getAtIndex(arr, k+1)) {
				x = getAtIndex(arr, k+1)
			} else {
				x = getAtIndex(arr, k-1) + 1
			}
			y := x - k

			// Diagonal moves
			for x < len(previous) && y < len(current) && previous[x].text == current[y].text {
				x, y = x+1, y+1
			}

			// Record move
			setAtIndex(arr, k, x)

			if x >= len(previous) && y >= len(current) {
				trace = append(trace, arr)
				return trace, nil
			}
		}

		newTrace := make([]int, len(arr))
		copy(newTrace, arr)
		trace = append(trace, newTrace)
	}

	return nil, errors.New("failed to find a path")
}

// Find shortest path by backtracking traces
func (m myersDiffer) backtrack(previous, current Document, shortestEdit [][]int) []positionsPair {
	path := []positionsPair{}
	x, y := len(previous), len(current)
	for depth := 0; depth < len(shortestEdit); depth++ {
		reverseD := len(shortestEdit) - 1 - depth // Start from the end
		arr := shortestEdit[reverseD]
		k := x - y

		var previousK int
		if k == -depth || (k != depth && getAtIndex(arr, k-1) < getAtIndex(arr, k+1)) {
			previousK = k + 1
		} else {
			previousK = k - 1
		}

		previousX := getAtIndex(arr, previousK)
		previousY := previousX - previousK
		// Handle diagonal backtracing
		for x > previousX && y > previousY {
			path = append(path, positionsPair{x1: x - 1, y1: y - 1, x2: x, y2: y})
			x, y = x-1, y-1
		}

		if reverseD > 0 {
			path = append(path, positionsPair{x1: previousX, y1: previousY, x2: x, y2: y})
		} else if path[len(path)-1].x1 != 0 && path[len(path)-1].y1 != 0 {
			// Add {0, 0, x, y} at the end if it wasn't already added
			path = append(path, positionsPair{x1: 0, y1: 0, x2: x, y2: y})
		}
		x, y = previousX, previousY
	}

	return path
}

// Utility function wrapping slice index access, allow negative numbers (starting from end)
func getAtIndex(s []int, index int) int {
	return s[getIndex(index, len(s))]
}

// Utility function wrapping slice index set, allow negative numbers (starting from end)
func setAtIndex(s []int, index, value int) {
	s[getIndex(index, len(s))] = value
}

// Get index, normal access if positive, access from the end if negative
func getIndex(index, length int) int {
	if index >= 0 {
		return index
	}
	return length + index
}

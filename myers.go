package gitdiff

import "errors"

type MyersDiffer struct{}

type positionsPair struct {
	x1, y1 int
	x2, y2 int
}

func (m MyersDiffer) ComputeDiff(previous, current Document) ([]Diff, error) {
	traces, err := m.shortestEdition(previous, current)
	if err != nil {
		return nil, err
	}

	path := m.backtrack(previous, current, traces)
	diffs := make([]Diff, len(path))
	for i, pos := range path {
		writeIndex := len(path) - i - 1
		if pos.x1 == pos.x2 {
			diffs[writeIndex] = NewDiff(Line{lineNumber: -1}, current[pos.y1], ADDED)
		} else if pos.y1 == pos.y2 {
			diffs[writeIndex] = NewDiff(previous[pos.x1], Line{lineNumber: -1}, DELETED)
		} else {
			diffs[writeIndex] = NewDiff(previous[pos.x1], current[pos.y1], UNCHANGED)
		}
	}

	return diffs, nil
}

func (m MyersDiffer) shortestEdition(previous, current Document) ([][]int, error) {
	maxSize := getMaxSize(previous, current)
	if maxSize == 0 {
		return [][]int{}, nil
	}

	arr := make([]int, maxSize*2+1)
	trace := [][]int{}
	for d := 0; d <= maxSize; d++ {
		for k := -d; k <= d; k += 2 {
			// Move by one
			var x int
			if k == -d || (k != d && getAtIndex(arr, k-1) < getAtIndex(arr, k+1)) {
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

	return nil, errors.New("failed to find the end")
}

func (m MyersDiffer) backtrack(previous, current Document, shortestEdit [][]int) []positionsPair {
	path := []positionsPair{}
	x, y := len(previous), len(current)
	for d := 0; d < len(shortestEdit); d++ {
		reverseD := len(shortestEdit) - 1 - d // Start from the end
		arr := shortestEdit[reverseD]
		k := x - y

		var previousK int
		if k == -d || (k != d && getAtIndex(arr, k-1) < getAtIndex(arr, k+1)) {
			previousK = k + 1
		} else {
			previousK = k - 1
		}

		previousX := getAtIndex(arr, previousK)
		previousY := previousX - previousK
		for x > previousX && y > previousY {
			path = append(path, positionsPair{x1: x - 1, y1: y - 1, x2: x, y2: y})
			x, y = x-1, y-1
		}

		if reverseD > 0 {
			path = append(path, positionsPair{x1: previousX, y1: previousY, x2: x, y2: y})
		}
		x, y = previousX, previousY
	}

	return path
}

func getMaxSize(a, b Document) int {
	return max(len(a), len(b))
}

func getAtIndex(s []int, index int) int {
	return s[getIndex(index, len(s))]
}

func setAtIndex(s []int, index, value int) {
	s[getIndex(index, len(s))] = value
}

func getIndex(index, length int) int {
	if index >= 0 {
		return index
	}
	// Negative index indicates to get from the end
	return length + index
}

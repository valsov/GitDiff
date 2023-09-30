package gitdiff

import (
	"fmt"
	"strconv"
)

const (
	UNCHANGED_COLOR string = "\033[39m"
	ADDED_COLOR     string = "\033[32m"
	DELETED_COLOR   string = "\033[31m"
)

func PrintDiffs(diffs []Diff) {
	for _, diff := range diffs {
		var (
			color                     string
			symbol                    rune
			linePrevious, lineCurrent string
		)

		if diff.Type == ADDED {
			color = ADDED_COLOR
			linePrevious = " "
			lineCurrent = strconv.Itoa(diff.New.lineNumber)
			symbol = '+'
		} else if diff.Type == DELETED {
			color = DELETED_COLOR
			linePrevious = strconv.Itoa(diff.Old.lineNumber)
			lineCurrent = " "
			symbol = '-'
		} else {
			color = UNCHANGED_COLOR
			linePrevious = strconv.Itoa(diff.Old.lineNumber)
			lineCurrent = strconv.Itoa(diff.New.lineNumber)
			symbol = ' '
		}

		fmt.Printf("%v%s %s\t%s\t  %s%s\n",
			color, string(symbol), linePrevious, lineCurrent, diff.GetText(), UNCHANGED_COLOR)
	}
}

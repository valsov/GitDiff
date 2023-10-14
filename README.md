# GitDiff
Git diff algorithm implementation. The program structure allows for easy addition of diff algorithms.

Currently, only [Myers diff](http://www.xmailserver.org/diff2.pdf) is implemented.

## Example use

```go
input1 := `DELETED
UNCHANGED
UNCHANGED`
previous := gitdiff.NewDocument(input1)

input2 := `ADDED
UNCHANGED
UNCHANGED
ADDED`
current := gitdiff.NewDocument(input2)

differ, _ := gitdiff.NewDiffProducer(gitdiff.MYERS_DIFF)
diffs, err := differ.ComputeDiff(previous, current, -1) // -1 to include all lines, not only changed ones
```

The diffs can then be handled as you like, or printed using the `printer` utility. The above example would produce this:
```diff
- 1               DELETED
+       1         ADDED
  2     2         UNCHANGED
  3     3         UNCHANGED
+       4         ADDED
```

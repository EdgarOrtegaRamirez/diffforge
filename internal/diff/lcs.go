// Package diff provides core diffing algorithms and data structures.
package diff

// OpType represents the type of diff operation.
type OpType int

const (
	OpEqual OpType = iota
	OpInsert
	OpDelete
)

// DiffOp represents a single diff operation.
type DiffOp struct {
	Type    OpType
	OldIdx  int // Index in old sequence (-1 for inserts)
	NewIdx  int // Index in new sequence (-1 for deletes)
	OldLine string
	NewLine string
}

// DiffResult holds the complete diff result.
type DiffResult struct {
	Ops      []DiffOp
	OldLines int
	NewLines int
	Added    int
	Removed  int
	Modified int
}

// Stats returns a summary of changes.
func (d *DiffResult) Stats() (added, removed, modified int) {
	for _, op := range d.Ops {
		switch op.Type {
		case OpInsert:
			added++
		case OpDelete:
			removed++
		case OpEqual:
			// no-op
		}
	}
	return
}

// LCS computes the longest common subsequence table.
// Returns a 2D table where lcs[i][j] is the length of LCS of a[0..i-1] and b[0..j-1].
func LCS(a, b []string) [][]int {
	n, m := len(a), len(b)
	table := make([][]int, n+1)
	for i := range table {
		table[i] = make([]int, m+1)
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				table[i][j] = table[i-1][j-1] + 1
			} else if table[i-1][j] > table[i][j-1] {
				table[i][j] = table[i-1][j]
			} else {
				table[i][j] = table[i][j-1]
			}
		}
	}
	return table
}

// LCSBacktrack traces back through the LCS table to produce diff operations.
func LCSBacktrack(table [][]int, a, b []string) []DiffOp {
	ops := make([]DiffOp, 0, len(a)+len(b))
	i, j := len(a), len(b)

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			ops = append(ops, DiffOp{
				Type:    OpEqual,
				OldIdx:  i - 1,
				NewIdx:  j - 1,
				OldLine: a[i-1],
				NewLine: b[j-1],
			})
			i--
			j--
		} else if j > 0 && (i == 0 || table[i][j-1] >= table[i-1][j]) {
			ops = append(ops, DiffOp{
				Type:    OpInsert,
				OldIdx:  -1,
				NewIdx:  j - 1,
				NewLine: b[j-1],
			})
			j--
		} else {
			ops = append(ops, DiffOp{
				Type:    OpDelete,
				OldIdx:  i - 1,
				NewIdx:  -1,
				OldLine: a[i-1],
			})
			i--
		}
	}

	// Reverse to get chronological order
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}

	return ops
}

// TextDiff computes a line-by-line diff between two texts.
func TextDiff(oldText, newText string) *DiffResult {
	oldLines := SplitLines(oldText)
	newLines := SplitLines(newText)

	table := LCS(oldLines, newLines)
	ops := LCSBacktrack(table, oldLines, newLines)

	result := &DiffResult{
		Ops:      ops,
		OldLines: len(oldLines),
		NewLines: len(newLines),
	}

	for _, op := range ops {
		switch op.Type {
		case OpInsert:
			result.Added++
		case OpDelete:
			result.Removed++
		}
	}

	return result
}

// SplitLines splits text into lines, preserving line endings.
func SplitLines(text string) []string {
	if text == "" {
		return nil
	}
	lines := make([]string, 0)
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			lines = append(lines, text[start:i+1])
			start = i + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}

// WordDiff computes a word-level diff between two lines.
func WordDiff(oldLine, newLine string) []DiffOp {
	oldWords := SplitWords(oldLine)
	newWords := SplitWords(newLine)

	table := LCS(oldWords, newWords)
	return LCSBacktrack(table, oldWords, newWords)
}

// SplitWords splits a line into words for word-level diffing.
func SplitWords(line string) []string {
	words := make([]string, 0)
	start := -1
	for i := 0; i < len(line); i++ {
		if isWordChar(line[i]) {
			if start == -1 {
				start = i
			}
		} else {
			if start != -1 {
				words = append(words, line[start:i])
				start = -1
			}
			words = append(words, string(line[i]))
		}
	}
	if start != -1 {
		words = append(words, line[start:])
	}
	return words
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

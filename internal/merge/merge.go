// Package merge provides three-way merge capabilities.
package merge

import (
	"fmt"
	"strings"
)

// MergeOp represents a merge operation.
type MergeOp struct {
	Type string // "keep_base", "keep_old", "keep_new", "conflict"
	Base string
	Old  string
	New  string
	Line int
}

// MergeResult holds the result of a three-way merge.
type MergeResult struct {
	Ops       []MergeOp
	Output    string
	Conflicts int
	Resolved  int
}

// ThreeWayMerge performs a three-way merge between base, old, and new.
func ThreeWayMerge(base, old, new string) *MergeResult {
	baseLines := splitLines(base)
	oldLines := splitLines(old)
	newLines := splitLines(new)

	result := &MergeResult{}

	// Simple line-by-line three-way merge
	maxLen := len(baseLines)
	if len(oldLines) > maxLen {
		maxLen = len(oldLines)
	}
	if len(newLines) > maxLen {
		maxLen = len(newLines)
	}

	var output []string

	for i := 0; i < maxLen; i++ {
		var baseLine, oldLine, newLine string
		if i < len(baseLines) {
			baseLine = baseLines[i]
		}
		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}

		// Three cases
		if baseLine == oldLine && baseLine == newLine {
			// All same - keep it
			output = append(output, baseLine)
			result.Ops = append(result.Ops, MergeOp{
				Type: "keep_base",
				Base: baseLine,
				Line: i + 1,
			})
			result.Resolved++
		} else if baseLine == oldLine {
			// Only new changed - take new
			output = append(output, newLine)
			result.Ops = append(result.Ops, MergeOp{
				Type: "keep_new",
				Base: baseLine,
				New:  newLine,
				Line: i + 1,
			})
			result.Resolved++
		} else if baseLine == newLine {
			// Only old changed - take old
			output = append(output, oldLine)
			result.Ops = append(result.Ops, MergeOp{
				Type: "keep_old",
				Base: baseLine,
				Old:  oldLine,
				Line: i + 1,
			})
			result.Resolved++
		} else {
			// Both changed - conflict
			result.Conflicts++
			result.Ops = append(result.Ops, MergeOp{
				Type: "conflict",
				Base: baseLine,
				Old:  oldLine,
				New:  newLine,
				Line: i + 1,
			})
			// Use conflict markers
			output = append(output, "<<<<<<< OLD")
			output = append(output, oldLine)
			output = append(output, "=======")
			output = append(output, newLine)
			output = append(output, ">>>>>>> NEW")
		}
	}

	result.Output = strings.Join(output, "\n")
	if len(output) > 0 && !strings.HasSuffix(result.Output, "\n") {
		result.Output += "\n"
	}

	return result
}

// FormatMergeResult formats a merge result.
func FormatMergeResult(result *MergeResult) string {
	if result == nil {
		return "No merge result."
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Merge: %d resolved, %d conflicts\n\n",
		result.Resolved, result.Conflicts)

	for _, op := range result.Ops {
		switch op.Type {
		case "conflict":
			fmt.Fprintf(&sb, "<<<<<<< OLD (line %d)\n", op.Line)
			sb.WriteString(op.Old + "\n")
			sb.WriteString("=======\n")
			sb.WriteString(op.New + "\n")
			sb.WriteString(">>>>>>> NEW\n\n")
		}
	}

	if result.Conflicts == 0 {
		sb.WriteString("Merge completed successfully.\n")
	}

	return sb.String()
}

func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	// Remove trailing empty string from split
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// Package format provides output formatters for diff results.
package format

import (
	"fmt"
	"strings"

	"github.com/EdgarOrtegaRamirez/diffforge/internal/diff"
)

// FormatType represents the output format type.
type FormatType string

const (
	FormatUnified     FormatType = "unified"
	FormatSideBySide  FormatType = "side-by-side"
	FormatJSON        FormatType = "json"
	FormatHTML        FormatType = "html"
	FormatContext     FormatType = "context"
	FormatMinimal     FormatType = "minimal"
)

// Unified formats a text diff in unified diff format.
func Unified(result *diff.DiffResult, oldName, newName string) string {
	if result == nil || len(result.Ops) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- %s\n", oldName))
	sb.WriteString(fmt.Sprintf("+++ %s\n", newName))

	// Group operations into hunks
	hunks := groupHunks(result.Ops, 3) // 3 lines of context
	for _, hunk := range hunks {
		sb.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
			hunk.oldStart, hunk.oldCount, hunk.newStart, hunk.newCount))
		for _, op := range hunk.ops {
			switch op.Type {
			case diff.OpEqual:
				sb.WriteString(" " + op.OldLine)
			case diff.OpDelete:
				sb.WriteString("-" + op.OldLine)
			case diff.OpInsert:
				sb.WriteString("+" + op.NewLine)
			}
		}
	}

	return sb.String()
}

// SideBySide formats a text diff in side-by-side format.
func SideBySide(result *diff.DiffResult, oldName, newName string, width int) string {
	if result == nil || len(result.Ops) == 0 {
		return ""
	}

	if width <= 0 {
		width = 80
	}
	halfWidth := (width - 3) / 2 // -3 for the separator column

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s vs %s ===\n\n", oldName, newName))

	lineNum := 1
	for _, op := range result.Ops {
		switch op.Type {
		case diff.OpEqual:
			oldLine := truncate(op.OldLine, halfWidth)
			newLine := truncate(op.NewLine, halfWidth)
			sb.WriteString(fmt.Sprintf("%4d | %-*s | %s\n", lineNum, halfWidth, oldLine, newLine))
			lineNum++
		case diff.OpDelete:
			oldLine := truncate(op.OldLine, halfWidth)
			sb.WriteString(fmt.Sprintf("%4d | %-*s | %-*s\n", lineNum, halfWidth, oldLine, halfWidth, ""))
			lineNum++
		case diff.OpInsert:
			newLine := truncate(op.NewLine, halfWidth)
			sb.WriteString(fmt.Sprintf("%4d | %-*s | %s\n", lineNum, halfWidth, "", newLine))
			lineNum++
		}
	}

	return sb.String()
}

// Minimal formats a text diff with minimal output.
func Minimal(result *diff.DiffResult) string {
	if result == nil || len(result.Ops) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	added, removed, _ := result.Stats()
	sb.WriteString(fmt.Sprintf("%d files changed, %d insertions(+), %d deletions(-)\n",
		1, added, removed))

	for _, op := range result.Ops {
		switch op.Type {
		case diff.OpInsert:
			sb.WriteString("+" + op.NewLine)
		case diff.OpDelete:
			sb.WriteString("-" + op.OldLine)
		}
	}

	return sb.String()
}

// Context formats a text diff in context diff format.
func Context(result *diff.DiffResult, oldName, newName string, contextLines int) string {
	if result == nil || len(result.Ops) == 0 {
		return ""
	}

	if contextLines <= 0 {
		contextLines = 3
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*** %s\n", oldName))
	sb.WriteString(fmt.Sprintf("--- %s\n", newName))

	hunks := groupHunks(result.Ops, contextLines)
	for _, hunk := range hunks {
		sb.WriteString("***************\n")
		for _, op := range hunk.ops {
			switch op.Type {
			case diff.OpEqual:
				sb.WriteString("  " + op.OldLine)
			case diff.OpDelete:
				sb.WriteString("- " + op.OldLine)
			case diff.OpInsert:
				sb.WriteString("+ " + op.NewLine)
			}
		}
	}

	return sb.String()
}

// hunk represents a group of related diff operations.
type hunk struct {
	oldStart int
	oldCount int
	newStart int
	newCount int
	ops      []diff.DiffOp
}

// groupHunks groups diff operations into hunks with context.
func groupHunks(ops []diff.DiffOp, context int) []hunk {
	if len(ops) == 0 {
		return nil
	}

	var hunks []hunk
	var current hunk
	inChange := false

	oldLine := 1
	newLine := 1

	for _, op := range ops {
		switch op.Type {
		case diff.OpEqual:
			if inChange {
				// Add context line to current hunk
				current.ops = append(current.ops, op)
				current.oldCount++
				current.newCount++
				context--
				if context <= 0 {
					hunks = append(hunks, current)
					current = hunk{}
					inChange = false
					context = 3
				}
			}
			oldLine++
			newLine++
		case diff.OpDelete:
			if !inChange {
				// Start new hunk
				inChange = true
				context = 3
				current.oldStart = oldLine
				current.newStart = newLine
				// Add preceding context lines
				// (simplified: just start the hunk)
			}
			current.ops = append(current.ops, op)
			current.oldCount++
			oldLine++
		case diff.OpInsert:
			if !inChange {
				inChange = true
				context = 3
				current.oldStart = oldLine
				current.newStart = newLine
			}
			current.ops = append(current.ops, op)
			current.newCount++
			newLine++
		}
	}

	// Add the last hunk
	if inChange {
		hunks = append(hunks, current)
	}

	return hunks
}

func truncate(s string, maxLen int) string {
	s = strings.TrimRight(s, "\n")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

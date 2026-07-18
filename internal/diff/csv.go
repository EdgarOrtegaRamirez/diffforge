package diff

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CSVDiffOp represents a diff operation on CSV data.
type CSVDiffOp struct {
	Type   OpType
	Row    int
	Old    []string
	New    []string
	Column string
	OldVal string
	NewVal string
}

// CSVDiffResult holds the result of a CSV diff.
type CSVDiffResult struct {
	Ops      []CSVDiffOp
	Headers  []string
	Added    int
	Removed  int
	Modified int
}

// CSVDiff computes a row-level diff between two CSV strings.
func CSVDiff(oldCSV, newCSV string) (*CSVDiffResult, error) {
	oldRows, err := parseCSV(oldCSV)
	if err != nil {
		return nil, fmt.Errorf("failed to parse old CSV: %w", err)
	}
	newRows, err := parseCSV(newCSV)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new CSV: %w", err)
	}

	result := &CSVDiffResult{}

	// Use first row as headers if present
	if len(oldRows) > 0 && len(newRows) > 0 {
		result.Headers = oldRows[0]
		oldRows = oldRows[1:]
		newRows = newRows[1:]
	}

	maxRows := len(oldRows)
	if len(newRows) > maxRows {
		maxRows = len(newRows)
	}

	for i := 0; i < maxRows; i++ {
		switch {
		case i >= len(oldRows):
			result.Ops = append(result.Ops, CSVDiffOp{
				Type: OpInsert,
				Row:  i + 1,
				New:  newRows[i],
			})
			result.Added++
		case i >= len(newRows):
			result.Ops = append(result.Ops, CSVDiffOp{
				Type: OpDelete,
				Row:  i + 1,
				Old:  oldRows[i],
			})
			result.Removed++
		default:
			// Compare cell by cell
			diffRow(oldRows[i], newRows[i], i+1, result)
		}
	}

	return result, nil
}

func diffRow(oldRow, newRow []string, rowIdx int, result *CSVDiffResult) {
	maxCols := len(oldRow)
	if len(newRow) > maxCols {
		maxCols = len(newRow)
	}

	modified := false
	for c := 0; c < maxCols; c++ {
		var oldVal, newVal string
		var colName string

		if c < len(oldRow) {
			oldVal = oldRow[c]
		}
		if c < len(newRow) {
			newVal = newRow[c]
		}
		if c < len(result.Headers) {
			colName = result.Headers[c]
		} else {
			colName = fmt.Sprintf("col_%d", c)
		}

		if oldVal != newVal {
			result.Ops = append(result.Ops, CSVDiffOp{
				Type:   OpModify,
				Row:    rowIdx,
				Old:    oldRow,
				New:    newRow,
				Column: colName,
				OldVal: oldVal,
				NewVal: newVal,
			})
			modified = true
		}
	}

	if modified {
		result.Modified++
	}
}

// OpModify is a special operation type for cell modifications.
const OpModify OpType = 3

func parseCSV(text string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(text))
	records, err := reader.ReadAll()
	if err != nil && err != io.EOF {
		return nil, err
	}
	return records, nil
}

// FormatCSVDiff formats a CSV diff result.
func FormatCSVDiff(result *CSVDiffResult) string {
	if result == nil || len(result.Ops) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "CSV Diff: %d added, %d removed, %d modified\n\n",
		result.Added, result.Removed, result.Modified)

	if len(result.Headers) > 0 {
		sb.WriteString("Headers: ")
		sb.WriteString(strings.Join(result.Headers, " | "))
		sb.WriteString("\n\n")
	}

	for _, op := range result.Ops {
		switch op.Type {
		case OpInsert:
			fmt.Fprintf(&sb, "+ Row %d: %s\n", op.Row, strings.Join(op.New, " | "))
		case OpDelete:
			fmt.Fprintf(&sb, "- Row %d: %s\n", op.Row, strings.Join(op.Old, " | "))
		case OpModify:
			fmt.Fprintf(&sb, "~ Row %d, %s: %q → %q\n", op.Row, op.Column, op.OldVal, op.NewVal)
		}
	}

	return sb.String()
}

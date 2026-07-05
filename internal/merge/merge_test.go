package merge

import (
	"strings"
	"testing"
)

func TestThreeWayMergeNoChanges(t *testing.T) {
	base := "line1\nline2\nline3"
	old := "line1\nline2\nline3"
	new := "line1\nline2\nline3"

	result := ThreeWayMerge(base, old, new)
	if result.Conflicts != 0 {
		t.Errorf("conflicts = %d, want 0", result.Conflicts)
	}
	if result.Resolved != 3 {
		t.Errorf("resolved = %d, want 3", result.Resolved)
	}
}

func TestThreeWayMergeOnlyOldChanges(t *testing.T) {
	base := "line1\nline2\nline3"
	old := "modified\nline2\nline3"
	new := "line1\nline2\nline3"

	result := ThreeWayMerge(base, old, new)
	if result.Conflicts != 0 {
		t.Errorf("conflicts = %d, want 0", result.Conflicts)
	}
	if !strings.Contains(result.Output, "modified") {
		t.Error("expected 'modified' in output")
	}
}

func TestThreeWayMergeOnlyNewChanges(t *testing.T) {
	base := "line1\nline2\nline3"
	old := "line1\nline2\nline3"
	new := "line1\nmodified\nline3"

	result := ThreeWayMerge(base, old, new)
	if result.Conflicts != 0 {
		t.Errorf("conflicts = %d, want 0", result.Conflicts)
	}
	if !strings.Contains(result.Output, "modified") {
		t.Error("expected 'modified' in output")
	}
}

func TestThreeWayMergeBothChangeSameLine(t *testing.T) {
	base := "line1\nline2\nline3"
	old := "old_change\nline2\nline3"
	new := "new_change\nline2\nline3"

	result := ThreeWayMerge(base, old, new)
	if result.Conflicts != 1 {
		t.Errorf("conflicts = %d, want 1", result.Conflicts)
	}
	if !strings.Contains(result.Output, "<<<<<<< OLD") {
		t.Error("expected conflict markers in output")
	}
}

func TestThreeWayMergeEmptyFiles(t *testing.T) {
	result := ThreeWayMerge("", "", "")
	if result.Conflicts != 0 {
		t.Errorf("conflicts = %d, want 0", result.Conflicts)
	}
}

func TestThreeWayMergeDifferentLengths(t *testing.T) {
	base := "line1\nline2"
	old := "line1\nline2\nline3"
	new := "line1\nline2\nline4"

	result := ThreeWayMerge(base, old, new)
	// line3 and line4 are both additions - they conflict
	if result.Conflicts < 1 {
		t.Errorf("conflicts = %d, want >= 1", result.Conflicts)
	}
}

func TestFormatMergeResult(t *testing.T) {
	result := &MergeResult{
		Resolved:  2,
		Conflicts: 1,
		Ops: []MergeOp{
			{Type: "conflict", Line: 1, Old: "old", New: "new"},
		},
	}

	output := FormatMergeResult(result)
	if !strings.Contains(output, "<<<<<<< OLD") {
		t.Error("expected conflict markers in output")
	}
	if !strings.Contains(output, "1 conflicts") {
		t.Error("expected conflict count in output")
	}
}

func TestFormatMergeResultNoConflicts(t *testing.T) {
	result := &MergeResult{
		Resolved:  3,
		Conflicts: 0,
	}

	output := FormatMergeResult(result)
	if !strings.Contains(output, "successfully") {
		t.Error("expected success message")
	}
}

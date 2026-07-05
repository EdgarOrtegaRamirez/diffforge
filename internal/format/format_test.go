package format

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/diffforge/internal/diff"
)

func TestUnified(t *testing.T) {
	result := &diff.DiffResult{
		Ops: []diff.DiffOp{
			{Type: diff.OpEqual, OldLine: "line1\n", NewLine: "line1\n"},
			{Type: diff.OpDelete, OldLine: "old\n"},
			{Type: diff.OpInsert, NewLine: "new\n"},
		},
	}

	output := Unified(result, "old.txt", "new.txt")
	if !strings.Contains(output, "--- old.txt") {
		t.Error("expected old file header")
	}
	if !strings.Contains(output, "+++ new.txt") {
		t.Error("expected new file header")
	}
	if !strings.Contains(output, "-old") {
		t.Error("expected delete marker")
	}
	if !strings.Contains(output, "+new") {
		t.Error("expected insert marker")
	}
}

func TestUnifiedEmpty(t *testing.T) {
	result := &diff.DiffResult{Ops: []diff.DiffOp{}}
	output := Unified(result, "a", "b")
	if output != "" {
		t.Error("expected empty output for no changes")
	}
}

func TestSideBySide(t *testing.T) {
	result := &diff.DiffResult{
		Ops: []diff.DiffOp{
			{Type: diff.OpEqual, OldLine: "same\n", NewLine: "same\n"},
			{Type: diff.OpDelete, OldLine: "old\n"},
			{Type: diff.OpInsert, NewLine: "new\n"},
		},
	}

	output := SideBySide(result, "a", "b", 80)
	if !strings.Contains(output, "same") {
		t.Error("expected equal line in output")
	}
}

func TestMinimal(t *testing.T) {
	result := &diff.DiffResult{
		Ops: []diff.DiffOp{
			{Type: diff.OpInsert, NewLine: "added\n"},
			{Type: diff.OpDelete, OldLine: "removed\n"},
		},
	}

	output := Minimal(result)
	if !strings.Contains(output, "+added") {
		t.Error("expected insert in output")
	}
	if !strings.Contains(output, "-removed") {
		t.Error("expected delete in output")
	}
}

func TestMinimalEmpty(t *testing.T) {
	result := &diff.DiffResult{Ops: []diff.DiffOp{}}
	output := Minimal(result)
	if !strings.Contains(output, "No differences") {
		t.Error("expected no differences message")
	}
}

func TestContext(t *testing.T) {
	result := &diff.DiffResult{
		Ops: []diff.DiffOp{
			{Type: diff.OpEqual, OldLine: "line1\n", NewLine: "line1\n"},
			{Type: diff.OpDelete, OldLine: "old\n"},
			{Type: diff.OpInsert, NewLine: "new\n"},
			{Type: diff.OpEqual, OldLine: "line3\n", NewLine: "line3\n"},
		},
	}

	output := Context(result, "old", "new", 2)
	if !strings.Contains(output, "*** old") {
		t.Error("expected old file header")
	}
	if !strings.Contains(output, "--- new") {
		t.Error("expected new file header")
	}
}

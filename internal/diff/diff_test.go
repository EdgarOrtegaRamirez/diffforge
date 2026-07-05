package diff

import (
	"os"
	"strings"
	"testing"
)

func TestTextDiff(t *testing.T) {
	tests := []struct {
		name    string
		old     string
		new     string
		added   int
		removed int
	}{
		{
			name:    "identical",
			old:     "line1\nline2\n",
			new:     "line1\nline2\n",
			added:   0,
			removed: 0,
		},
		{
			name:    "add line",
			old:     "line1\n",
			new:     "line1\nline2\n",
			added:   1,
			removed: 0,
		},
		{
			name:    "remove line",
			old:     "line1\nline2\n",
			new:     "line1\n",
			added:   0,
			removed: 1,
		},
		{
			name:    "modify line",
			old:     "hello\n",
			new:     "world\n",
			added:   1, // LCS sees this as delete+insert
			removed: 1,
		},
		{
			name:    "empty files",
			old:     "",
			new:     "",
			added:   0,
			removed: 0,
		},
		{
			name:    "completely different",
			old:     "aaa\nbbb\n",
			new:     "ccc\nddd\n",
			added:   2, // No common lines
			removed: 2,
		},
		{
			name:    "partial match",
			old:     "aaa\nbbb\nccc\n",
			new:     "aaa\nddd\nccc\n",
			added:   1, // bbb removed, ddd added
			removed: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TextDiff(tt.old, tt.new)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			added, removed, _ := result.Stats()
			if added != tt.added {
				t.Errorf("added = %d, want %d", added, tt.added)
			}
			if removed != tt.removed {
				t.Errorf("removed = %d, want %d", removed, tt.removed)
			}
		})
	}
}

func TestJSONDiff(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
	}{
		{
			name: "identical",
			old:  `{"a": 1}`,
			new:  `{"a": 1}`,
		},
		{
			name: "add key",
			old:  `{"a": 1}`,
			new:  `{"a": 1, "b": 2}`,
		},
		{
			name: "remove key",
			old:  `{"a": 1, "b": 2}`,
			new:  `{"a": 1}`,
		},
		{
			name: "modify value",
			old:  `{"a": 1}`,
			new:  `{"a": 2}`,
		},
		{
			name: "nested change",
			old:  `{"a": {"b": 1}}`,
			new:  `{"a": {"b": 2}}`,
		},
		{
			name: "array change",
			old:  `{"items": [1, 2, 3]}`,
			new:  `{"items": [1, 4, 3]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := JSONDiff(tt.old, tt.new)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}
		})
	}
}

func TestJSONDiffInvalid(t *testing.T) {
	_, err := JSONDiff("not json", `{"a": 1}`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCSVDiff(t *testing.T) {
	oldCSV := "name,age,city\nAlice,30,NYC\nBob,25,LA\n"
	newCSV := "name,age,city\nAlice,31,NYC\nBob,25,LA\nCharlie,35,Chicago\n"

	result, err := CSVDiff(oldCSV, newCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Added != 1 {
		t.Errorf("added = %d, want 1", result.Added)
	}
	if result.Modified != 1 {
		t.Errorf("modified = %d, want 1", result.Modified)
	}
}

func TestCSVDiffIdentical(t *testing.T) {
	csv := "a,b\n1,2\n"
	result, err := CSVDiff(csv, csv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Added != 0 || result.Removed != 0 || result.Modified != 0 {
		t.Error("expected no changes for identical CSV")
	}
}

func TestBinaryDiff(t *testing.T) {
	old := []byte("Hello, World!")
	new := []byte("Hello, Universe!")

	result := BinaryDiff(old, new)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Added == 0 && result.Removed == 0 {
		t.Error("expected some changes")
	}
}

func TestBinaryDiffIdentical(t *testing.T) {
	data := []byte("Hello, World!")
	result := BinaryDiff(data, data)

	if result.Added != 0 || result.Removed != 0 {
		t.Error("expected no changes for identical data")
	}
}

func TestWordDiff(t *testing.T) {
	old := "the cat sat on the mat"
	new := "the dog sat on the hat"

	ops := WordDiff(old, new)
	if len(ops) == 0 {
		t.Error("expected some operations")
	}

	changes := 0
	for _, op := range ops {
		if op.Type != OpEqual {
			changes++
		}
	}

	if changes == 0 {
		t.Error("expected some changes")
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input string
		count int
	}{
		{"", 0},
		{"line1", 1},
		{"line1\n", 1},
		{"line1\nline2", 2},
		{"line1\nline2\n", 2},
		{"line1\nline2\nline3", 3},
	}

	for _, tt := range tests {
		lines := SplitLines(tt.input)
		if len(lines) != tt.count {
			t.Errorf("SplitLines(%q) returned %d lines, want %d", tt.input, len(lines), tt.count)
		}
	}
}

func TestSplitWords(t *testing.T) {
	words := SplitWords("hello world 123")
	if len(words) == 0 {
		t.Error("expected some words")
	}

	found := false
	for _, w := range words {
		if w == "hello" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find 'hello' in words")
	}
}

func TestFormatJSONDiff(t *testing.T) {
	result := &JSONDiffResult{
		Ops: []JSONDiffOp{
			{Type: OpInsert, Path: "$.b", New: 2},
			{Type: OpDelete, Path: "$.c", Old: 3},
		},
		Added:   1,
		Removed: 1,
	}

	output := FormatJSONDiff(result, "")
	if !strings.Contains(output, "+ $.b") {
		t.Error("expected insert operation in output")
	}
	if !strings.Contains(output, "- $.c") {
		t.Error("expected delete operation in output")
	}
}

func TestFormatJSONDiffEmpty(t *testing.T) {
	result := &JSONDiffResult{Ops: []JSONDiffOp{}}
	output := FormatJSONDiff(result, "")
	if !strings.Contains(output, "No differences") {
		t.Error("expected 'No differences' message")
	}
}

func TestFormatCSVDiff(t *testing.T) {
	result := &CSVDiffResult{
		Ops: []CSVDiffOp{
			{Type: OpInsert, Row: 1, New: []string{"Alice", "30"}},
		},
		Headers: []string{"name", "age"},
		Added:   1,
	}

	output := FormatCSVDiff(result)
	if !strings.Contains(output, "+ Row 1") {
		t.Error("expected insert operation in output")
	}
}

func TestFormatBinaryDiff(t *testing.T) {
	result := &BinaryDiffResult{
		Chunks: []BinaryChunk{
			{Offset: 0, Length: 5, Type: "insert", New: []byte("Hello")},
		},
		TotalSize: 5,
		Added:     5,
	}

	output := FormatBinaryDiff(result)
	if !strings.Contains(output, "+ Offset") {
		t.Error("expected insert operation in output")
	}
}

func TestFormatDirDiff(t *testing.T) {
	result := &DirDiffResult{
		Ops: []DirDiffOp{
			{Type: OpInsert, Path: "new.txt", NewSize: 100},
			{Type: OpDelete, Path: "old.txt", OldSize: 200},
		},
		Added:   1,
		Removed: 1,
	}

	output := FormatDirDiff(result, "dir1", "dir2")
	if !strings.Contains(output, "+ new.txt") {
		t.Error("expected insert in output")
	}
	if !strings.Contains(output, "- old.txt") {
		t.Error("expected delete in output")
	}
}

func TestLCS(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "c", "e", "d"}

	table := LCS(a, b)

	// LCS should be "b", "c", "d" = length 3
	if table[len(a)][len(b)] != 3 {
		t.Errorf("LCS length = %d, want 3", table[len(a)][len(b)])
	}
}

func TestFilesIdentical(t *testing.T) {
	a := []byte("hello")
	b := []byte("hello")
	c := []byte("world")

	if !FilesIdentical(a, b) {
		t.Error("expected identical files")
	}
	if FilesIdentical(a, c) {
		t.Error("expected different files")
	}
}

func TestYAMLDiff(t *testing.T) {
	old := "name: Alice\nage: 30\n"
	new := "name: Alice\nage: 31\ncity: NYC\n"

	result, err := YAMLDiff(old, new)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Should have modifications (age changed) and additions (city added)
	if result.Added == 0 && result.Modified == 0 {
		t.Error("expected some changes")
	}
}

func TestFormatYAMLDiff(t *testing.T) {
	result := &JSONDiffResult{
		Ops: []JSONDiffOp{
			{Type: OpInsert, Path: "$.city", New: "NYC"},
		},
		Added: 1,
	}

	output := FormatYAMLDiff(result)
	if !strings.Contains(output, "+ $.city") {
		t.Error("expected insert operation in output")
	}
}

func TestDirDiff(t *testing.T) {
	// Test with actual temp directories
	oldDir := t.TempDir()
	newDir := t.TempDir()

	// Create files in old
	os.WriteFile(oldDir+"/file1.txt", []byte("old content"), 0644)
	os.WriteFile(oldDir+"/file2.txt", []byte("same content"), 0644)

	// Create files in new
	os.WriteFile(newDir+"/file1.txt", []byte("new content"), 0644)
	os.WriteFile(newDir+"/file2.txt", []byte("same content"), 0644)
	os.WriteFile(newDir+"/file3.txt", []byte("added content"), 0644)

	result, err := DirDiff(oldDir, newDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Added != 1 {
		t.Errorf("added = %d, want 1", result.Added)
	}
	if result.Modified != 1 {
		t.Errorf("modified = %d, want 1", result.Modified)
	}
}

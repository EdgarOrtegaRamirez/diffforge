package diff

import (
	"bytes"
	"fmt"
	"strings"
)

// BinaryChunk represents a chunk of binary data in a diff.
type BinaryChunk struct {
	Offset int
	Length int
	Type   string // "equal", "insert", "delete", "modify"
	Old    []byte
	New    []byte
}

// BinaryDiffResult holds the result of a binary diff.
type BinaryDiffResult struct {
	Chunks    []BinaryChunk
	TotalSize int
	Added     int
	Removed   int
	Modified  int
}

// BinaryDiff computes a byte-level diff between two binary slices.
func BinaryDiff(oldData, newData []byte) *BinaryDiffResult {
	result := &BinaryDiffResult{
		TotalSize: len(oldData),
	}

	i, j := 0, 0
	for i < len(oldData) || j < len(newData) {
		if i < len(oldData) && j < len(newData) && oldData[i] == newData[j] {
			// Find the extent of equal bytes
			start := i
			for i < len(oldData) && j < len(newData) && oldData[i] == newData[j] {
				i++
				j++
			}
			result.Chunks = append(result.Chunks, BinaryChunk{
				Offset: start,
				Length: i - start,
				Type:   "equal",
				Old:    oldData[start:i],
			})
		} else if j < len(newData) && (i >= len(oldData) || countMatches(oldData, i, newData, j) < 3) {
			// Insert: new bytes that don't match old
			start := j
			for j < len(newData) && (i >= len(oldData) || newData[j] != oldData[i]) {
				j++
			}
			chunk := BinaryChunk{
				Offset: start,
				Length: j - start,
				Type:   "insert",
				New:    newData[start:j],
			}
			result.Chunks = append(result.Chunks, chunk)
			result.Added += j - start
		} else if i < len(oldData) {
			// Delete: old bytes that don't match new
			start := i
			for i < len(oldData) && (j >= len(newData) || oldData[i] != newData[j]) {
				i++
			}
			result.Chunks = append(result.Chunks, BinaryChunk{
				Offset: start,
				Length: i - start,
				Type:   "delete",
				Old:    oldData[start:i],
			})
			result.Removed += i - start
		}
	}

	return result
}

func countMatches(a []byte, ai int, b []byte, bi int) int {
	count := 0
	for ai < len(a) && bi < len(b) && a[ai] == b[bi] {
		count++
		ai++
		bi++
	}
	return count
}

// FormatBinaryDiff formats a binary diff result as hex dump.
func FormatBinaryDiff(result *BinaryDiffResult) string {
	if result == nil || len(result.Chunks) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Binary Diff: %d bytes total, %d added, %d removed\n\n",
		result.TotalSize, result.Added, result.Removed))

	for _, chunk := range result.Chunks {
		switch chunk.Type {
		case "insert":
			sb.WriteString(fmt.Sprintf("+ Offset 0x%04X (%d bytes)\n", chunk.Offset, chunk.Length))
			sb.WriteString(hexDump(chunk.New))
		case "delete":
			sb.WriteString(fmt.Sprintf("- Offset 0x%04X (%d bytes)\n", chunk.Offset, chunk.Length))
			sb.WriteString(hexDump(chunk.Old))
		}
	}

	return sb.String()
}

func hexDump(data []byte) string {
	var sb strings.Builder
	for i := 0; i < len(data); i += 16 {
		sb.WriteString(fmt.Sprintf("  %04X: ", i))
		// Hex part
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		for j := i; j < end; j++ {
			sb.WriteString(fmt.Sprintf("%02X ", data[j]))
		}
		// Pad if needed
		for j := end; j < i+16; j++ {
			sb.WriteString("   ")
		}
		// ASCII part
		sb.WriteString(" |")
		for j := i; j < end; j++ {
			if data[j] >= 32 && data[j] < 127 {
				sb.WriteByte(data[j])
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteString("|\n")
	}
	return sb.String()
}

// FilesIdentical checks if two byte slices are identical.
func FilesIdentical(a, b []byte) bool {
	return bytes.Equal(a, b)
}

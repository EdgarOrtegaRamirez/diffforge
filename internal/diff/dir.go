package diff

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DirDiffOp represents a file-level diff operation.
type DirDiffOp struct {
	Type     OpType
	Path     string
	OldSize  int64
	NewSize  int64
	OldHash  string
	NewHash  string
	IsDir    bool
}

// DirDiffResult holds the result of a directory diff.
type DirDiffResult struct {
	Ops       []DirDiffOp
	Added     int
	Removed   int
	Modified  int
	Identical int
}

// DirDiff computes a diff between two directories.
func DirDiff(oldDir, newDir string) (*DirDiffResult, error) {
	result := &DirDiffResult{}

	oldFiles, err := walkDir(oldDir)
	if err != nil {
		return nil, fmt.Errorf("failed to walk old directory: %w", err)
	}

	newFiles, err := walkDir(newDir)
	if err != nil {
		return nil, fmt.Errorf("failed to walk new directory: %w", err)
	}

	// Compare files
	allPaths := make(map[string]bool)
	for p := range oldFiles {
		allPaths[p] = true
	}
	for p := range newFiles {
		allPaths[p] = true
	}

	paths := make([]string, 0, len(allPaths))
	for p := range allPaths {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, path := range paths {
		oldInfo, oldExists := oldFiles[path]
		newInfo, newExists := newFiles[path]

		switch {
		case oldExists && !newExists:
			result.Ops = append(result.Ops, DirDiffOp{
				Type:    OpDelete,
				Path:    path,
				OldSize: oldInfo.Size,
				OldHash: oldInfo.Hash,
				IsDir:   oldInfo.IsDir,
			})
			result.Removed++
		case !oldExists && newExists:
			result.Ops = append(result.Ops, DirDiffOp{
				Type:    OpInsert,
				Path:    path,
				NewSize: newInfo.Size,
				NewHash: newInfo.Hash,
				IsDir:   newInfo.IsDir,
			})
			result.Added++
		default:
			if oldInfo.Hash != newInfo.Hash {
				result.Ops = append(result.Ops, DirDiffOp{
					Type:    OpDelete, // Show as modification
					Path:    path,
					OldSize: oldInfo.Size,
					NewSize: newInfo.Size,
					OldHash: oldInfo.Hash,
					NewHash: newInfo.Hash,
				})
				result.Modified++
			} else {
				result.Identical++
			}
		}
	}

	return result, nil
}

type fileInfo struct {
	Size  int64
	Hash  string
	IsDir bool
}

func walkDir(root string) (map[string]fileInfo, error) {
	files := make(map[string]fileInfo)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		fi := fileInfo{
			Size:  info.Size(),
			IsDir: info.IsDir(),
		}

		if !info.IsDir() {
			hash, err := fileHash(path)
			if err != nil {
				return err
			}
			fi.Hash = hash
		}

		files[relPath] = fi
		return nil
	})

	return files, err
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// FormatDirDiff formats a directory diff result.
func FormatDirDiff(result *DirDiffResult, oldDir, newDir string) string {
	if result == nil || len(result.Ops) == 0 {
		return "Directories are identical."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Directory Diff: %s vs %s\n", oldDir, newDir))
	sb.WriteString(fmt.Sprintf("%d added, %d removed, %d modified, %d identical\n\n",
		result.Added, result.Removed, result.Modified, result.Identical))

	for _, op := range result.Ops {
		switch op.Type {
		case OpInsert:
			sb.WriteString(fmt.Sprintf("+ %s (%d bytes)\n", op.Path, op.NewSize))
		case OpDelete:
			if op.NewSize > 0 {
				sb.WriteString(fmt.Sprintf("~ %s (%d → %d bytes)\n", op.Path, op.OldSize, op.NewSize))
			} else {
				sb.WriteString(fmt.Sprintf("- %s (%d bytes)\n", op.Path, op.OldSize))
			}
		}
	}

	return sb.String()
}

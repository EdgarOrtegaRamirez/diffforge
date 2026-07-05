package diff

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// JSONDiffOp represents a diff operation on a JSON structure.
type JSONDiffOp struct {
	Type    OpType
	Path    string
	Old     interface{}
	New     interface{}
	Key     string
}

// JSONDiffResult holds the result of a JSON structural diff.
type JSONDiffResult struct {
	Ops      []JSONDiffOp
	Added    int
	Removed  int
	Modified int
}

// JSONDiff computes a structural diff between two JSON values.
func JSONDiff(oldJSON, newJSON string) (*JSONDiffResult, error) {
	var oldVal, newVal interface{}

	if err := json.Unmarshal([]byte(oldJSON), &oldVal); err != nil {
		return nil, fmt.Errorf("failed to parse old JSON: %w", err)
	}
	if err := json.Unmarshal([]byte(newJSON), &newVal); err != nil {
		return nil, fmt.Errorf("failed to parse new JSON: %w", err)
	}

	result := &JSONDiffResult{}
	diffJSONValue(oldVal, newVal, "$", result)
	return result, nil
}

func diffJSONValue(oldVal, newVal interface{}, path string, result *JSONDiffResult) {
	oldMap, oldIsMap := oldVal.(map[string]interface{})
	newMap, newIsMap := newVal.(map[string]interface{})
	oldArr, oldIsArr := oldVal.([]interface{})
	newArr, newIsArr := newVal.([]interface{})

	switch {
	case oldIsMap && newIsMap:
		diffJSONObject(oldMap, newMap, path, result)
	case oldIsArr && newIsArr:
		diffJSONArray(oldArr, newArr, path, result)
	case oldVal != newVal:
		result.Ops = append(result.Ops, JSONDiffOp{
			Type: OpDelete,
			Path: path,
			Old:  oldVal,
			New:  newVal,
		})
		result.Modified++
	}
}

func diffJSONObject(oldMap, newMap map[string]interface{}, path string, result *JSONDiffResult) {
	// Collect all keys
	allKeys := make(map[string]bool)
	for k := range oldMap {
		allKeys[k] = true
	}
	for k := range newMap {
		allKeys[k] = true
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		childPath := path + "." + key
		oldVal, oldExists := oldMap[key]
		newVal, newExists := newMap[key]

		switch {
		case oldExists && !newExists:
			result.Ops = append(result.Ops, JSONDiffOp{
				Type: OpDelete,
				Path: childPath,
				Key:  key,
				Old:  oldVal,
			})
			result.Removed++
		case !oldExists && newExists:
			result.Ops = append(result.Ops, JSONDiffOp{
				Type: OpInsert,
				Path: childPath,
				Key:  key,
				New:  newVal,
			})
			result.Added++
		default:
			diffJSONValue(oldVal, newVal, childPath, result)
		}
	}
}

func diffJSONArray(oldArr, newArr []interface{}, path string, result *JSONDiffResult) {
	maxLen := len(oldArr)
	if len(newArr) > maxLen {
		maxLen = len(newArr)
	}

	for i := 0; i < maxLen; i++ {
		childPath := fmt.Sprintf("%s[%d]", path, i)
		if i >= len(oldArr) {
			result.Ops = append(result.Ops, JSONDiffOp{
				Type: OpInsert,
				Path: childPath,
				New:  newArr[i],
			})
			result.Added++
		} else if i >= len(newArr) {
			result.Ops = append(result.Ops, JSONDiffOp{
				Type: OpDelete,
				Path: childPath,
				Old:  oldArr[i],
			})
			result.Removed++
		} else {
			diffJSONValue(oldArr[i], newArr[i], childPath, result)
		}
	}
}

// FormatJSONDiff formats a JSON diff result as a human-readable string.
func FormatJSONDiff(result *JSONDiffResult, format string) string {
	if result == nil || len(result.Ops) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("JSON Diff: %d added, %d removed, %d modified\n\n",
		result.Added, result.Removed, result.Modified))

	for _, op := range result.Ops {
		switch op.Type {
		case OpInsert:
			newVal, _ := json.MarshalIndent(op.New, "  ", "  ")
			sb.WriteString(fmt.Sprintf("+ %s = %s\n", op.Path, string(newVal)))
		case OpDelete:
			oldVal, _ := json.MarshalIndent(op.Old, "  ", "  ")
			sb.WriteString(fmt.Sprintf("- %s = %s\n", op.Path, string(oldVal)))
		case OpEqual:
			// Skip equal operations in output
		}
	}

	return sb.String()
}

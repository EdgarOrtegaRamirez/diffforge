package diff

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAMLDiff computes a structural diff between two YAML documents.
func YAMLDiff(oldYAML, newYAML string) (*JSONDiffResult, error) {
	var oldVal, newVal interface{}

	if err := yaml.Unmarshal([]byte(oldYAML), &oldVal); err != nil {
		return nil, fmt.Errorf("failed to parse old YAML: %w", err)
	}
	if err := yaml.Unmarshal([]byte(newYAML), &newVal); err != nil {
		return nil, fmt.Errorf("failed to parse new YAML: %w", err)
	}

	result := &JSONDiffResult{}
	diffJSONValue(oldVal, newVal, "$", result)
	return result, nil
}

// FormatYAMLDiff formats a YAML diff result.
func FormatYAMLDiff(result *JSONDiffResult) string {
	if result == nil || len(result.Ops) == 0 {
		return "No differences found."
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "YAML Diff: %d added, %d removed, %d modified\n\n",
		result.Added, result.Removed, result.Modified)

	for _, op := range result.Ops {
		switch op.Type {
		case OpInsert:
			fmt.Fprintf(&sb, "+ %s = %v\n", op.Path, op.New)
		case OpDelete:
			fmt.Fprintf(&sb, "- %s = %v\n", op.Path, op.Old)
		}
	}

	return sb.String()
}

package rules

import (
	"testing"
)

func TestIgnoreWhitespace(t *testing.T) {
	rule := IgnoreWhitespace()
	result := rule.Transform("  hello  ")
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestIgnoreBlankLines(t *testing.T) {
	rule := IgnoreBlankLines()
	result := rule.Transform("   ")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
	result = rule.Transform("hello")
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestIgnoreComments(t *testing.T) {
	rule := IgnoreComments()

	tests := []struct {
		input string
		want  string
	}{
		{"// this is a comment", ""},
		{"# python comment", ""},
		{"-- sql comment", ""},
		{"code here", "code here"},
	}

	for _, tt := range tests {
		result := rule.Transform(tt.input)
		if result != tt.want {
			t.Errorf("Transform(%q) = %q, want %q", tt.input, result, tt.want)
		}
	}
}

func TestIgnoreTrailingWhitespace(t *testing.T) {
	rule := IgnoreTrailingWhitespace()
	result := rule.Transform("hello   ")
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestRuleSetApply(t *testing.T) {
	rs := NewRuleSet()
	rs.AddRule("whitespace", IgnoreWhitespace().Transform)
	rs.AddRule("comments", IgnoreComments().Transform)

	// Should ignore comments
	result := rs.Apply("// comment")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}

	// Should trim whitespace
	result = rs.Apply("  hello  ")
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestRuleSetApplyToLines(t *testing.T) {
	rs := NewRuleSet()
	rs.AddRule("whitespace", IgnoreWhitespace().Transform)

	lines := []string{"  hello  ", "  world  ", "  test  "}
	result := rs.ApplyToLines(lines)

	expected := []string{"hello", "world", "test"}
	for i, r := range result {
		if r != expected[i] {
			t.Errorf("line %d: got %q, want %q", i, r, expected[i])
		}
	}
}

func TestPredefinedRuleSets(t *testing.T) {
	// Test code rule set
	codeRules := PredefinedRuleSets["code"]
	if codeRules == nil {
		t.Error("expected code rule set")
	}

	// Should handle comments
	result := codeRules.Apply("// comment")
	if result != "" {
		t.Errorf("expected empty for comment, got %q", result)
	}

	// Test whitespace rule set
	whitespaceRules := PredefinedRuleSets["whitespace"]
	if whitespaceRules == nil {
		t.Error("expected whitespace rule set")
	}

	// Test strict rule set
	strictRules := PredefinedRuleSets["strict"]
	if strictRules == nil {
		t.Error("expected strict rule set")
	}
}

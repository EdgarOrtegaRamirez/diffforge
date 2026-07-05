// Package rules provides ignore rules for diff operations.
package rules

import (
	"regexp"
	"strings"
)

// Rule defines an ignore rule for diff operations.
type Rule struct {
	Name    string
	Pattern *regexp.Regexp
	Transform func(string) string
}

// RuleSet is a collection of rules.
type RuleSet struct {
	Rules []Rule
}

// NewRuleSet creates a new rule set with default rules.
func NewRuleSet() *RuleSet {
	return &RuleSet{
		Rules: []Rule{},
	}
}

// AddRule adds a rule to the set.
func (rs *RuleSet) AddRule(name string, transform func(string) string) {
	rs.Rules = append(rs.Rules, Rule{
		Name:      name,
		Transform: transform,
	})
}

// AddPatternRule adds a regex-based rule.
func (rs *RuleSet) AddPatternRule(name, pattern string) {
	re := regexp.MustCompile(pattern)
	rs.Rules = append(rs.Rules, Rule{
		Name:    name,
		Pattern: re,
		Transform: func(s string) string {
			return re.ReplaceAllString(s, "")
		},
	})
}

// IgnoreWhitespace creates a rule that ignores whitespace differences.
func IgnoreWhitespace() Rule {
	return Rule{
		Name: "whitespace",
		Transform: func(s string) string {
			return strings.TrimSpace(s)
		},
	}
}

// IgnoreBlankLines creates a rule that ignores blank lines.
func IgnoreBlankLines() Rule {
	return Rule{
		Name: "blank-lines",
		Transform: func(s string) string {
			if strings.TrimSpace(s) == "" {
				return ""
			}
			return s
		},
	}
}

// IgnoreComments creates a rule that ignores comments in various languages.
func IgnoreComments() Rule {
	return Rule{
		Name: "comments",
		Transform: func(s string) string {
			trimmed := strings.TrimSpace(s)
			// Single-line comments
			for _, prefix := range []string{"//", "#", "--", ";"} {
				if strings.HasPrefix(trimmed, prefix) {
					return ""
				}
			}
			return s
		},
	}
}

// IgnoreTrailingWhitespace creates a rule that ignores trailing whitespace.
func IgnoreTrailingWhitespace() Rule {
	return Rule{
		Name: "trailing-whitespace",
		Transform: func(s string) string {
			return strings.TrimRight(s, " \t")
		},
	}
}

// IgnoreLeadingWhitespace creates a rule that ignores leading whitespace.
func IgnoreLeadingWhitespace() Rule {
	return Rule {
		Name: "leading-whitespace",
		Transform: func(s string) string {
			return strings.TrimLeft(s, " \t")
		},
	}
}

// IgnoreAllWhitespace creates a rule that ignores all whitespace.
func IgnoreAllWhitespace() Rule {
	return Rule{
		Name: "all-whitespace",
		Transform: func(s string) string {
			return strings.Map(func(r rune) rune {
				if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
					return -1
				}
				return r
			}, s)
		},
	}
}

// IgnoreLineOrder creates a rule that sorts lines (for set-based comparison).
func IgnoreLineOrder() Rule {
	return Rule{
		Name: "line-order",
		Transform: func(s string) string {
			return s // This rule is handled specially in the diff engine
		},
	}
}

// Apply applies all rules to a line and returns the transformed result.
func (rs *RuleSet) Apply(line string) string {
	result := line
	for _, rule := range rs.Rules {
		if rule.Transform != nil {
			result = rule.Transform(result)
		}
	}
	return result
}

// ApplyToLines applies all rules to a slice of lines.
func (rs *RuleSet) ApplyToLines(lines []string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = rs.Apply(line)
	}
	return result
}

// PredefinedRuleSets provides common rule sets.
var PredefinedRuleSets = map[string]*RuleSet{
	"code": func() *RuleSet {
		rs := NewRuleSet()
		rs.AddRule("whitespace", IgnoreWhitespace().Transform)
		rs.AddRule("blank-lines", IgnoreBlankLines().Transform)
		rs.AddRule("comments", IgnoreComments().Transform)
		return rs
	}(),
	"whitespace": func() *RuleSet {
		rs := NewRuleSet()
		rs.AddRule("whitespace", IgnoreWhitespace().Transform)
		return rs
	}(),
	"strict": NewRuleSet(),
}

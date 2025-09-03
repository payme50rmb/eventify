package eventify

import "strings"

// Matcher is a struct that represents a string matcher.
type Matcher struct {
	pattern   string // The pattern without wildcards
	matchType int    // 0: exact, 1: prefix, 2: suffix, 3: wildcard
}

const (
	matchExact    = iota // exact match
	matchPrefix          // prefix*
	matchSuffix          // *suffix
	matchContains        // *middle*
	matchWildcard        // *
)

// NewMatcher creates a new matcher with the specified pattern.
// The pattern can be:
// - "*" matches any string
// - "prefix*" matches strings starting with "prefix"
// - "*suffix" matches strings ending with "suffix"
// - "exact" matches exactly "exact"
func NewMatcher(s string) *Matcher {
	m := &Matcher{}

	switch {
	case s == "*":
		m.matchType = matchWildcard
	case len(s) > 1 && s[0] == '*' && s[len(s)-1] == '*':
		// Handle patterns like *middle* if needed in the future
		m.matchType = matchContains // Not currently supported, fallback to exact
		m.pattern = s[1 : len(s)-1]
	case s[0] == '*':
		m.matchType = matchSuffix
		m.pattern = s[1:]
	case len(s) > 0 && s[len(s)-1] == '*':
		m.matchType = matchPrefix
		m.pattern = s[:len(s)-1]
	default:
		m.matchType = matchExact
		m.pattern = s
	}

	return m
}

// Match returns true if the target string matches the matcher.
func (m *Matcher) Match(target string) bool {
	switch m.matchType {
	case matchWildcard:
		return true
	case matchPrefix:
		return len(target) >= len(m.pattern) && target[:len(m.pattern)] == m.pattern
	case matchSuffix:
		return len(target) >= len(m.pattern) && target[len(target)-len(m.pattern):] == m.pattern
	case matchContains:
		return strings.Contains(target, m.pattern)
	default: // matchExact
		return target == m.pattern
	}
}

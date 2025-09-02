package eventify

import "strings"

// Matcher is a struct that represents a string matcher.
type Matcher struct {
	s string
}

// NewMatcher creates a new matcher with the specified string.
func NewMatcher(s string) *Matcher {
	return &Matcher{
		s: s,
	}
}

// Match returns true if the target string matches the matcher
// if s == "*", it matches any string
// if s == "foo*", it matches any string that starts with "foo"
// if s == "*bar", it matches any string that ends with "bar"
func (m *Matcher) Match(target string) bool {
	if m.s == "*" {
		return true
	}
	if strings.HasPrefix(m.s, "*") {
		return strings.HasSuffix(target, m.s[1:])
	}
	if strings.HasSuffix(m.s, "*") {
		return strings.HasPrefix(target, m.s[:len(m.s)-1])
	}
	return m.s == target
}

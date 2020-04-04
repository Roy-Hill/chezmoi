package chezmoi

import "github.com/bmatcuk/doublestar"

// An PatternSet is a set of patterns.
type PatternSet struct {
	includes StringSet
	excludes StringSet
}

// A PatternSetOption sets an option on a pattern set.
type PatternSetOption func(*PatternSet)

// NewPatternSet returns a new PatternSet.
func NewPatternSet(options ...PatternSetOption) *PatternSet {
	ps := &PatternSet{
		includes: NewStringSet(),
		excludes: NewStringSet(),
	}
	for _, option := range options {
		option(ps)
	}
	return ps
}

// Add adds a pattern to ps.
func (ps *PatternSet) Add(pattern string, include bool) error {
	if _, err := doublestar.Match(pattern, ""); err != nil {
		return nil
	}
	if include {
		ps.includes.Add(pattern)
	} else {
		ps.excludes.Add(pattern)
	}
	return nil
}

// Match returns if name matches any pattern in ps.
func (ps *PatternSet) Match(name string) bool {
	for pattern := range ps.excludes {
		if ok, _ := doublestar.Match(pattern, name); ok {
			return false
		}
	}
	for pattern := range ps.includes {
		if ok, _ := doublestar.Match(pattern, name); ok {
			return true
		}
	}
	return false
}

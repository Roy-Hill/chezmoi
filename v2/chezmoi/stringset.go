package chezmoi

// A StringSet is a set of strings.
type StringSet map[string]struct{}

// NewStringSet returns a new StringSet containing elements.
func NewStringSet(elements ...string) StringSet {
	s := make(StringSet)
	s.Add(elements...)
	return s
}

// Add adds elements to s.
func (s StringSet) Add(elements ...string) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Contains returns true if element is in s.
func (s StringSet) Contains(element string) bool {
	_, ok := s[element]
	return ok
}

// Elements returns all the elements of s.
func (s StringSet) Elements() []string {
	elements := make([]string, 0, len(s))
	for element := range s {
		elements = append(elements, element)
	}
	return elements
}

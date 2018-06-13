package testsql

import (
	"sync"
)

// Set represents a set container
type Set struct {
	m map[string]struct{}
	l *sync.RWMutex
}

// Add adds value to set.
func (s *Set) Add(value string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.m[value] = struct{}{}
}

// Values returns a slice of values.
func (s *Set) Values() []string {
	s.l.RLock()
	defer s.l.RUnlock()
	values := make([]string, 0)
	for v := range s.m {
		values = append(values, v)
	}
	return values
}

// NewSet returns a set.
func NewSet() *Set {
	return &Set{l: new(sync.RWMutex), m: make(map[string]struct{})}
}

// package safe provides a thread-safe holder for a value that may be updated
// by multiple goroutines.
package safe

import (
	"sync"
)

// Safe contains a thread-safe value, also allowing getters to wait until the value is non-nil.
type Safe struct {
	value interface{}
	cond  *sync.Cond
}

// New create a new Safe instance given a value
func New(value interface{}) *Safe {
	cond := sync.NewCond(&sync.Mutex{})
	return &Safe{value: value, cond: cond}
}

// Get returns the value, even if nil.
func (s *Safe) Get() interface{} {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	return s.value
}

// GetWhenDefined returns the value, waiting until the value is not nil.
func (s *Safe) GetWhenDefined() interface{} {
	s.cond.L.Lock()
	for s.value == nil {
		s.cond.Wait()
	}
	v := s.value
	s.cond.L.Unlock()
	return v
}

// Put sets a new value, which may be nil. If non-nil, any getters waiting via GetWhenDefined will
// wake up with the new value.
func (s *Safe) Put(value interface{}) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	broadcast := s.value == nil && value != nil
	s.value = value

	if broadcast {
		s.cond.Broadcast()
	}
}

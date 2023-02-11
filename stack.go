package chipeight

import (
	"errors"
	"sync"
)

// TODO: Update to use generics
type Stack struct {
	lock     sync.Mutex
	elements []interface{}
}

func (s *Stack) Push(d interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.elements = append(s.elements, d)
}

func (s *Stack) Pop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.elements)
	if l < 1 {
		return
	}

	s.elements = s.elements[:l-1]
}

func (s *Stack) Top() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.elements)
	if l < 1 {
		return nil, errors.New("stack is currently empty")
	}

	return s.elements[l-1], nil
}

func (s *Stack) Empty() bool {
	return len(s.elements) == 0
}

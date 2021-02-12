package stack

import (
	"errors"
	"sync"
)

type stack struct {
	lock sync.Mutex
	elements []interface{}
}

func (s *stack) Push(d interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.elements = append(s.elements, d)
}

func (s *stack) Pop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.elements)
	if l < 1 {
		return
	}

	s.elements = s.elements[:l-1]
}

func (s *stack) Top() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.elements)
	if l < 1 {
		return nil, errors.New("stack is currently empty")
	}

	return s.elements[l - 1], nil
}

func (s *stack) Empty() bool {
	return len(s.elements) == 0
}
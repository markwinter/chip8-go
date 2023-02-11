package chipeight

import (
	"reflect"
	"testing"
)

func TestStack_Push(t *testing.T) {
	s := Stack{}
	s.Push(0)
	s.Push(1)
	s.Push(2)

	want := []interface{}{0, 1, 2}

	if !reflect.DeepEqual(s.elements, want) {
		t.Errorf("Expected %v but got %v", want, s.elements)
	}
}

func TestStack_Pop(t *testing.T) {
	s := Stack{}
	s.Push(0)
	s.Push(1)
	s.Push(2)

	s.Pop()

	want := []interface{}{0, 1}

	if !reflect.DeepEqual(s.elements, want) {
		t.Errorf("Expected %v but got %v", want, s.elements)
	}
}

func TestStack_PopEmpty(t *testing.T) {
	s := Stack{}
	s.Pop()
}

func TestStack_Top(t *testing.T) {
	s := Stack{}
	s.Push(0)
	s.Push(1)

	if value, _ := s.Top(); value != 1 {
		t.Errorf("Expected 1 but got %v", value)
	}
}

func TestStack_TopEmpty(t *testing.T) {
	s := Stack{}

	if _, err := s.Top(); err == nil {
		t.Errorf("expected error but it was nil")
	}
}

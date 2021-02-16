package chipeight

import (
	"testing"
)

func Test_opA000(t *testing.T) {
	data := []byte{0xA1, 0x11}

	c8 := NewChipeight()
	c8.LoadBytes(data)

	c8.Step()

	want := uint16(0x111)

	if c8.GetIndexRegister() != want {
		t.Errorf("Expected 0x%X but got 0x%X", want, c8.GetIndexRegister())
	}
}

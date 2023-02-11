package chipeight

import (
	"testing"
)

func Test_opA000(t *testing.T) {
	data := []byte{0xA1, 0x11}

	c8 := NewChipeight()
	c8.LoadBytes(data)

	c8.Step()

	if c8.indexRegister != 0x111 {
		t.Errorf("Expected 0x111 but got 0x%X", c8.indexRegister)
	}
}

func Test_CallRet(t *testing.T) {
	/*
		0x2202  -> current PC is pushed to stack (0x200), PC set to 0x202 (the RET op)
		0x00EE  -> PC set to top of stack (0x200)
	*/
	data := []byte{0x22, 0x02, 0x00, 0xEE}

	c8 := NewChipeight()
	c8.LoadBytes(data)

	c8.Step()

	if c8.programCounter != 0x202 {
		t.Errorf("PC should be 0x202 but got 0x%X", c8.programCounter)
	}
	if value, _ := c8.stack.Top(); value.(uint16) != 0x200 {
		t.Errorf("Top of stack should be 0x200 but got 0x%X", value)
	}

	c8.Step()

	if c8.programCounter != 0x202 {
		t.Errorf("PC should be 0x202 but got 0x%X", c8.programCounter)
	}
}

package chipeight

import (
	"log"
	"math/rand"
)

// Returns VX
func getRegisterX(opcode uint16) uint16 {
	return (opcode & 0x0F00) >> 8
}

// Returns VY
func getRegisterY(opcode uint16) uint16 {
	return (opcode & 0x00F0) >> 4
}

// 00E0: Clear the screen
// 00EE: Return from subroutine
func op0000(c *Chipeight) {
	switch c.currentOpcode & 0x000F {
	case 0x0:
		c.screen = [64 * 32]uint8{}
		c.programCounter += 2
	case 0xE:
		value, err := c.stack.Top()
		if err != nil {
			log.Panicf("Tried to return from subroutine but stack was empty")
		}

		c.programCounter = value.(uint16)

		c.stack.Pop()
	}
}

// 1NNN: Jump to NNN
func op1000(c *Chipeight) {
	c.programCounter = c.currentOpcode & 0x0FFF
}

// 2NNN: Call subroutine at NNN
func op2000(c *Chipeight) {
	c.stack.Push(c.programCounter)
	c.programCounter = c.currentOpcode & 0x0FFF
}

// 3XNN: Skip next instruction if VX == NN
func op3000(c *Chipeight) {
	value := uint8(c.currentOpcode & 0x00FF)
	register := getRegisterX(c.currentOpcode)

	if value == c.registers[register] {
		c.programCounter += 2
	}

	c.programCounter += 2
}

// 4XNN: Skips next instruction if VX != NN
func op4000(c *Chipeight) {
	value := uint8(c.currentOpcode & 0x00FF)
	register := getRegisterX(c.currentOpcode)

	if value != c.registers[register] {
		c.programCounter += 2
	}

	c.programCounter += 2
}

// 5XY0: Skips next instruction if VX equals VY
func op5000(c *Chipeight) {
	registerX := getRegisterX(c.currentOpcode)
	registerY := getRegisterY(c.currentOpcode)

	if registerX == registerY {
		c.programCounter += 2
	}

	c.programCounter += 2
}

// 6XNN: Sets VX to NN
func op6000(c *Chipeight) {
	register := getRegisterX(c.currentOpcode)
	c.registers[register] = uint8(c.currentOpcode & 0x00FF)
	c.programCounter += 2
}

// 7XNN: Adds NN to VX
func op7000(c *Chipeight) {
	register := getRegisterX(c.currentOpcode)
	c.registers[register] += uint8(c.currentOpcode & 0x00FF)
	c.programCounter += 2
}

// 8XY0: Set VX = VY
// 8XY1: Set VX = VX|VY
// 8XY2: Set VX = VX&VY
// 8XY3: Set VX = VX^VY
// 8XY4: Set VX += VY
// 8XY5: Set VX -= VY
// 8XY6: Store least significant bit of VX in VF and shift VX right 1
// 8XY7: Set VX = VY - VX. Set VF=0 when there's a borrow, else 1
// 8XYE: Store most significant bit of VX in VF then shift VX left 1
func op8000(c *Chipeight) {
	switch c.currentOpcode & 0x000F {
	case 0x0:

	}

	c.programCounter += 2
}

// 9XY0: Skips next instruction if VX doesn't equal VY
func op9000(c *Chipeight) {
	registerX := getRegisterX(c.currentOpcode)
	registerY := getRegisterY(c.currentOpcode)

	if registerX != registerY {
		c.programCounter += 2
	}

	c.programCounter += 2
}

// ANNN: Sets I to NNN
func opA000(c *Chipeight) {
	c.indexRegister = c.currentOpcode & 0x0FFF
	c.programCounter += 2
}

// BNNN: Jump to the address NNN plus V0
func opB000(c *Chipeight) {
	value := uint16(c.registers[0])
	address := c.currentOpcode & 0x0FFF
	c.programCounter = value + address
}

// CXNN: Sets VX to the result of a NN & randomNumber
func opC000(c *Chipeight) {
	randomNumber := uint8(rand.Intn(255))
	nn := uint8(c.currentOpcode & 0x00FF)

	register := getRegisterX(c.currentOpcode)

	c.registers[register] = nn & randomNumber

	c.programCounter += 2
}

// DXYN: Draw at (VX, VY) with width=8, height=N+1
func opD000(c *Chipeight) {
	c.programCounter += 2
}

// EX9E: Skip next instruction if key stored in VX is pressed
// EXA1: Skip next instruction if key stored in VX isn't pressed
func opE000(c *Chipeight) {
	register := getRegisterX(c.currentOpcode)
	key := c.registers[register]

	switch c.currentOpcode & 0x000F {
	case 0xE:
		if c.keys[key] == 1 {
			c.programCounter += 2
		}
	case 0x1:
		if c.keys[key] == 0 {
			c.programCounter += 2
		}
	}
	c.programCounter += 2
}

// FX33: Stores the binary-coded decimal representation of VX,
// with the most significant of three digits at the address in I,
// the middle digit at I plus 1, and the least significant digit at I plus 2
func opF000(c *Chipeight) {
	opcode := c.currentOpcode & 0x00FF
	switch opcode {
	case 0x33:
		register := getRegisterX(c.currentOpcode)
		c.memory[c.indexRegister] = c.registers[register] / 100
		c.memory[c.indexRegister+1] = (c.registers[register] / 10) % 10
		c.memory[c.indexRegister+2] = (c.registers[register] % 100) % 10
	}

	c.programCounter += 2
}

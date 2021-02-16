package chipeight

import (
	"log"
	"math/rand"
)

// Returns VX
func getRegisterX(opcode uint16) uint8 {
	return uint8((opcode & 0x0F00) >> 8)
}

// Returns VY
func getRegisterY(opcode uint16) uint8 {
	return uint8((opcode & 0x00F0) >> 4)
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

		c.programCounter = value.(uint16) + 2 // +2 so we skip over the CALL instruction

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

	if c.registers[registerX] == c.registers[registerY] {
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
	registerX := getRegisterX(c.currentOpcode)
	registerY := getRegisterY(c.currentOpcode)

	switch c.currentOpcode & 0x000F {
	case 0x0:
		c.registers[registerX] = c.registers[registerY]
	case 0x1:
		c.registers[registerX] = c.registers[registerX] | c.registers[registerY]
	case 0x2:
		c.registers[registerX] = c.registers[registerX] & c.registers[registerY]
	case 0x3:
		c.registers[registerX] = c.registers[registerX] ^ c.registers[registerY]
	case 0x4:
		sum := uint16(c.registers[registerX]) + uint16(c.registers[registerY])
		if sum > 255 {
			c.registers[registerVF] = 1
		} else {
			c.registers[registerVF] = 0
		}
		c.registers[registerX] += c.registers[registerY]
	case 0x5:
		if c.registers[registerY] > c.registers[registerX] {
			c.registers[registerVF] = 0
		} else {
			c.registers[registerVF] = 1
		}
		c.registers[registerX] -= c.registers[registerY]
	case 0x6:
		c.registers[registerVF] = c.registers[registerX] & 0x01
		c.registers[registerX] >>= 1
	case 0x7:
		if c.registers[registerY] > c.registers[registerX] {
			c.registers[registerVF] = 1
		} else {
			c.registers[registerVF] = 0
		}
		c.registers[registerX] = c.registers[registerY] - c.registers[registerX]
	case 0xE:
		c.registers[registerVF] = (c.registers[registerX] & 0x80) >> 7
		c.registers[registerX] <<= 1
	}

	c.programCounter += 2
}

// 9XY0: Skips next instruction if VX doesn't equal VY
func op9000(c *Chipeight) {
	registerX := getRegisterX(c.currentOpcode)
	registerY := getRegisterY(c.currentOpcode)

	if c.registers[registerX] != c.registers[registerY] {
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
// Each row of 8 pixels is read as bit-coded starting from memory location I
// VF is set to 1 if any screen pixels are flipped from set to unset
func opD000(c *Chipeight) {
	registerX := getRegisterX(c.currentOpcode)
	registerY := getRegisterY(c.currentOpcode)

	x := uint16(c.registers[registerX] % screenWidth)
	y := uint16(c.registers[registerY] % screenHeight)

	width := uint16(spriteWidth)
	height := c.currentOpcode & 0x000F

	c.registers[registerVF] = 0

	for row := uint16(0); row < height; row++ {
		spriteByte := c.memory[c.indexRegister+row]

		for col := uint16(0); col < width; col++ {
			spritePixel := spriteByte & (0x80 >> col)

			if spritePixel == 0 {
				continue
			}

			screenLoc := (y+row)*screenWidth + (x + col)

			if c.screen[screenLoc] == 1 {
				c.registers[registerVF] = 1
			}

			c.screen[screenLoc] ^= 1
		}
	}

	c.shouldDraw = true
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

// FX07: Set VX = delay timer
// FX0A: A key press is awaited, then stored in VX (blocking)
// FX15: Set delay timer = VX
// FX18: Set sound timer = VX
// FX1E: Set I += VX
// FX29: Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
// FX33: Stores the binary-coded decimal representation of VX,
// with the most significant of three digits at the address in I,
// the middle digit at I plus 1, and the least significant digit at I plus 2
// FX55: Store V0-VX (inclusive) starting at memory address I
// FX65: Fill V0-VX (inclusive) with values from memory address I
func opF000(c *Chipeight) {
	opcode := c.currentOpcode & 0x00FF
	switch opcode {
	case 0x07:
		register := getRegisterX(c.currentOpcode)
		c.registers[register] = c.delayTimer
	case 0x0A:
		log.Printf("opcode unimplemented: 0x%X", c.currentOpcode)
	case 0x15:
		register := getRegisterX(c.currentOpcode)
		c.delayTimer = c.registers[register]
	case 0x18:
		register := getRegisterX(c.currentOpcode)
		c.soundTimer = c.registers[register]
	case 0x1E:
		register := getRegisterX(c.currentOpcode)
		c.indexRegister += uint16(c.registers[register])
	case 0x29:
		register := getRegisterX(c.currentOpcode)
		character := c.registers[register]
		c.indexRegister = uint16(fontStartLoc + (5 * character))
	case 0x33:
		register := getRegisterX(c.currentOpcode)
		value := c.registers[register]

		c.memory[c.indexRegister+2] = value % 10
		value /= 10

		c.memory[c.indexRegister+1] = value % 10
		value /= 10

		c.memory[c.indexRegister] = value % 10
	case 0x55:
		register := getRegisterX(c.currentOpcode)
		for i := uint8(0); i <= register; i++ {
			c.memory[c.indexRegister+uint16(i)] = c.registers[i]
		}
	case 0x65:
		register := getRegisterX(c.currentOpcode)
		for i := uint8(0); i <= register; i++ {
			c.registers[i] = c.memory[c.indexRegister+uint16(i)]
		}
	}

	c.programCounter += 2
}

package chipeight

import (
	"io/ioutil"
	"log"
)

const (
	fontStartLoc    = 0x50
	programStartLoc = 0x200
	registerVF      = 15
	screenWidth     = 64
	screenHeight    = 32
	spriteWidth     = 8
)

var (
	fontSet = [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
)

type Chipeight struct {
	memory [4096]uint8

	registers     [16]uint8
	indexRegister uint16

	screen [screenWidth * screenHeight]uint8
	keys   [16]uint8

	currentOpcode  uint16
	programCounter uint16

	delayTimer uint8
	soundTimer uint8

	shouldDraw bool

	stack Stack
}

func NewChipeight() *Chipeight {
	c := &Chipeight{
		programCounter: programStartLoc,
	}

	for i := 0; i < len(fontSet); i++ {
		c.memory[fontStartLoc+i] = fontSet[i]
	}

	return c
}

func (c *Chipeight) LoadROM(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i++ {
		c.memory[programStartLoc+i] = data[i]
	}

	return nil
}

func (c *Chipeight) LoadBytes(data []byte) error {
	for i := 0; i < len(data); i++ {
		c.memory[programStartLoc+i] = data[i]
	}

	return nil
}

// For testing without GUI
func (c *Chipeight) Run() {
	for {
		c.Step()
	}
}

func (c *Chipeight) Step() {
	c.currentOpcode = uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])

	switch c.currentOpcode & 0xF000 {
	case 0x0000:
		op0000(c)
	case 0x1000:
		op1000(c)
	case 0x2000:
		op2000(c)
	case 0x3000:
		op3000(c)
	case 0x4000:
		op4000(c)
	case 0x5000:
		op5000(c)
	case 0x6000:
		op6000(c)
	case 0x7000:
		op7000(c)
	case 0x8000:
		op8000(c)
	case 0x9000:
		op9000(c)
	case 0xA000:
		opA000(c)
	case 0xB000:
		opB000(c)
	case 0xC000:
		opC000(c)
	case 0xD000:
		opD000(c)
	case 0xE000:
		opE000(c)
	case 0xF000:
		opF000(c)
	default:
		log.Printf("unknown opcode: 0x%X", c.currentOpcode)
	}

	// Increase PC except for special opcodes: CALL, JMP
	if (c.currentOpcode&0xF000) != 0x1000 && (c.currentOpcode&0xF000) != 0x2000 {
		c.programCounter += 2
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}

	if c.soundTimer > 0 {
		c.soundTimer--
	}
}

func (c *Chipeight) ShouldDraw() bool {
	sd := c.shouldDraw
	if sd {
		c.shouldDraw = false
	}
	return sd
}

func (c *Chipeight) GetScreen() [screenWidth * screenHeight]uint8 {
	return c.screen
}

package chipeight

type chipeight struct {
	memory [4096]uint8
	registers [16]uint8
	screen [64 * 32]uint8

	currentOpcode  uint16
	programCounter uint16
	indexRegister  uint16

	delayTimer uint8
	soundTimer uint8
}
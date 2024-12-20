package machine

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MAX_ADDRESS int = 0xb800
	REG_NUM     int = 10
	DEV_NUM     int = 256
)

const (
	MASK_WORD = 0xFFFFFF
	MAX_WORD  = (1 << 23) - 1
	MIN_WORD  = -(1 << 23)
)

var RegisterNames [REG_NUM]string = [REG_NUM]string{"A", "X", "L", "B", "S", "T", "F", "", "PC", "SW"}

const (
	A = 0
	X = 1
	L = 2
	B = 3
	S = 4
	T = 5
	F = 6
	// preskočimo 7
	PC = 8
	SW = 9
)

type Machine struct {
	registers []int
	memory    []uint8
	devices   []Device
}

func NewMachine() Machine {
	m := Machine{
		registers: make([]int, REG_NUM),
		memory:    make([]uint8, MAX_ADDRESS),
		devices:   make([]Device, DEV_NUM),
	}

	//default devices
	m.devices[0] = newInputDevice()
	m.devices[1] = newOutputDevice()
	m.devices[2] = newErrorDevice()

	return m
}

func (m Machine) initDevice(ind int) {
	fileName := fmt.Sprintf("%X.dev", ind)

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0o600) // add the filemode argument to this function

	if err != nil {
		panic("Failed to create file: " + err.Error())
	}

	//defer file.Close()

	input := bufio.NewReader(file)
	output := bufio.NewWriter(file)

	m.devices[ind].input = input
	m.devices[ind].output = output

}

func (m Machine) readDevice(ind int) byte {
	if m.devices[ind].input == nil {
		m.initDevice(ind)
	}

	res := m.devices[ind].Read()

	//fmt.Printf("Reading dev %d, val %d\n", ind, res)

	return res
}

func (m Machine) writeDevice(ind int, b byte) {
	if m.devices[ind].output == nil {
		m.initDevice(ind)
	}

	//fmt.Printf("Writing dev %d, val %d\n", ind, b)

	m.devices[ind].Write(b)
}

func validAddress(addr int) bool {
	return (addr >= 0 && addr < MAX_ADDRESS)
}

func (m Machine) GetByte(addr int) int {
	if !validAddress(addr) {
		fmt.Printf("Invalid memory address: %d\n", addr)

		return 0
	}
	return int(m.memory[addr]) & 0xFF
}

func (m Machine) setByte(addr, val int) {
	m.memory[addr] = byte(val & 0xFF)
}

func (m Machine) GetWord(addr int) int {
	// val := int(m.memory[addr])
	// val = val << 8
	// val += int(m.memory[addr+1])
	// val = val << 8
	// val += int(m.memory[addr+2])
	// return val

	return (m.GetByte(addr) << 16) | (m.GetByte(addr+1) << 8) | (m.GetByte(addr + 2))
}

func (m Machine) GetWordInt(addr int) int {
	return signedWordToInt((m.GetByte(addr) << 16) | (m.GetByte(addr+1) << 8) | (m.GetByte(addr + 2)))
}

func (m Machine) setWord(addr, val int) {
	// m.memory[addr+2] = uint8(val % (1 << 8))
	// val = val >> 8
	// m.memory[addr+1] = uint8(val % (1 << 8))
	// val = val >> 8
	// m.memory[addr+2] = uint8(val % (1 << 8))

	m.memory[addr+2] = uint8(val & 0xFF)
	val = val >> 8
	m.memory[addr+1] = uint8(val & 0xFF)
	val = val >> 8
	m.memory[addr] = uint8(val & 0xFF)

}

func intToSignedWord(val int) int {
	if val >= 0 {
		return val & MASK_WORD
	}
	return ^(-val - 1) & MASK_WORD
}

func signedWordToInt(word int) int {
	if word <= MAX_WORD {
		return word
	}
	return -(^word & MASK_WORD) - 1
}

func (m Machine) getReg(reg int) int {
	return (m.registers[reg])
}

func (m Machine) setReg(reg int, val int) {
	m.registers[reg] = val
}

func (m Machine) setAByte(val int) {
	val &= 0xFF
	m.registers[A] &= 0xFFFF00
	m.registers[A] |= val
}

func (m Machine) getAByte() int {
	return m.registers[A] & 0xFF
}

func (m Machine) GetRegInt(reg int) int {
	return signedWordToInt(m.registers[reg])
}

func (m Machine) GetRegistersRaw() []int {
	return m.registers
}

func (m Machine) equalSW() bool {
	return m.registers[SW] == 0
}

func (m Machine) greaterSW() bool {
	return m.registers[SW] == 1
}

func (m Machine) lessSW() bool {
	return m.registers[SW] == -1
}

package machine

const (
	MAX_ADDRESS int = 0xb800
	REG_NUM     int = 10
	DEV_NUM     int = 256
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

func (m Machine) getByte(addr int) int {
	return int(m.memory[addr])
}

func (m Machine) setByte(addr, val int) {
	m.memory[addr] = byte(val)
}

func (m Machine) getWord(addr int) int {
	val := int(m.memory[addr])
	val = val << 8
	val += int(m.memory[addr+1])
	val = val << 8
	val += int(m.memory[addr+2])
	return val
}

func (m Machine) setWord(addr, val int) {
	m.memory[addr+2] = uint8(val % (1 << 8))
	val = val >> 8
	m.memory[addr+1] = uint8(val % (1 << 8))
	val = val >> 8
	m.memory[addr+2] = uint8(val % (1 << 8))

}

func (m Machine) GetRegisters() []int {
	return m.registers
}

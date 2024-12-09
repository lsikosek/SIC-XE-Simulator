package machine

import (
	"fmt"
	"time"
)

func notImplemented(mnemonic string) {}

func invalidOpcode(opcode int) {

}

func invalidAddressing() {}

func (m Machine) fetch() int {

	ind := m.registers[PC]
	m.registers[PC]++

	return m.GetByte(ind)

}

func (m Machine) step() {
	m.execute()
}

func (m Machine) stepper(quitChan chan int, speedChan chan int, spd int) {
	var dur int = 1000000000 / (spd)
	var timeDuration time.Duration = time.Duration(dur)
	check := time.After(timeDuration)

	cnt := 0

	for {
		select {
		case <-quitChan:
			return
		case spd = <-speedChan:
			dur = 1000000000 / (spd)
			timeDuration = time.Duration(dur)
		case <-check:
			fmt.Println(cnt)
			cnt++
			m.step()
			check = time.After(timeDuration)
		}
	}

}

func (m Machine) execute() {
	opcode := m.fetch()

	if m.execF1(opcode) {
		return
	}

	op := m.fetch()

	if m.execF2(opcode, op) {
		return
	}

	operand := m.fetch()

	// Pripravimo operand
	temp := (op & ((1 << 4) - 1)) << 8
	operand += temp

	// Pripravimo nibble in opcode
	nixbpe := op >> 4
	temp = (opcode & 3) << 4
	nixbpe += temp

	opcode >>= 2
	opcode <<= 2 //LATEST CHANGEEEEEE

	if nixbpe%2 == 1 {
		temp = m.fetch()
		operand <<= 8
		operand += temp
	}

	//fmt.Printf("%s(%X) %d\n", InstructionMap[opcode], opcode, operand)

	m.execSICF3F4(opcode, nixbpe, operand)
	return
}

func (m Machine) execF1(opcode int) bool {
	switch opcode {
	case FIX:
		notImplemented("FIX")
	case FLOAT:
		notImplemented("FLOAT")
	case HIO:
		// TODO: halt IO channel number (A)
		notImplemented("HIO")
	case NORM:
		notImplemented("NORM za floate")
	case SIO:
		// TODO: start IO channel number (A); address of channel given by (S)
		notImplemented("SIO")
	case TIO:
		// Test IO channel number (A)
		notImplemented("TIO")
	default:
		return false
	}
	return true
}

func (m Machine) execF2(opcode, op int) bool {
	var r1, r2 int
	r1 = op & 0x00F0
	r2 = op & 0x000F

	r1 >>= 4
	if r1 >= len(m.registers) || r2 >= len(m.registers) {
		return false
	}
	switch opcode {
	case ADDR:
		m.registers[r2] = m.registers[r2] + m.registers[r1]
	case CLEAR:
		m.registers[r1] = 0
	case COMPR:
		v1, v2 := m.registers[r1], m.registers[r2]
		if v1 < v2 {
			m.registers[SW] = -1
		} else if v1 > v2 {
			m.registers[SW] = 1
		} else {
			m.registers[SW] = 0
		}
	case MULR:
		m.registers[r2] = m.registers[r2] * m.registers[r1]
	case RMO:
		m.registers[r2] = m.registers[r1]
	case SHIFTL:
		v := m.registers[r1]
		n := r2

		temp := (1 << n) - 1
		temp <<= 24 - n
		temp = temp & v
		temp >>= 24 - n

		v = (v << n) | temp

		m.registers[r1] = v
	case SHIFTR:
		v := m.registers[r1]
		n := r2
		temp := 0
		if check := v & (1 << 23); check != 0 {
			temp = (1 << n) - 1
			temp <<= 24 - n
		}

		v = (v >> n) | temp

		m.registers[r1] = v
	case SUBR:
		m.registers[r2] = m.registers[r2] - m.registers[r1]
	case SVC:
		//TODO: generiraj SVC interrupt
	case TIXR:
		m.registers[X]++
		v1, v2 := m.registers[X], m.registers[r2]
		if v1 < v2 {
			m.registers[SW] = -1
		} else if v1 > v2 {
			m.registers[SW] = 1
		} else {
			m.registers[SW] = 0
		}
	default:
		return false
	}

	return true

}

func (m Machine) execSICF3F4(opcode, nixbpe, operand int) bool {

	un := m.getAddressFromInstruction(nixbpe, operand)

	ni := nixbpe >> 4

	var valW, valB int

	rA := &m.registers[A]

	fmt.Printf("%s(%X) %d\n", InstructionMap[opcode], opcode, un)

	switch ni {
	case 0b_01:
		valW = un
		valB = un
	case 0b_10:
		valW = m.GetWord(m.GetWord(un))
		valB = m.GetByte(m.GetWord(un))
	default:
		valW = m.GetWord(un)
		valB = m.GetByte(un)
	}

	switch opcode {
	case ADD:
		m.registers[A] += valW
	case AND:
		m.registers[A] = m.registers[A] & valW

	case COMP:

		v1, v2 := m.registers[A], valW
		if v1 < v2 {
			m.registers[SW] = -1
		} else if v1 > v2 {
			m.registers[SW] = 1
		} else {
			m.registers[SW] = 0
		}
	case COMPF:
		notImplemented("COMPF")
	case DIV:
		*rA = *rA / valW
	case DIVF:
		notImplemented("DIVF")
	case J:
		fmt.Printf("UN: %d, valW: %d\n", un, valW)
		m.registers[PC] = valW
	case JEQ:
		if m.registers[SW] == 0 {
			m.registers[PC] = valW
		}
	case JGT:
		if m.registers[SW] == 1 {
			m.registers[PC] = valW
		}
	case JLT:
		if m.registers[SW] == -1 {
			m.registers[PC] = valW
		}
	case JSUB:
		m.registers[L] = m.registers[PC]
		m.registers[PC] = valW
	case LDA:
		*rA = valW
	case LDB:
		m.registers[B] = valW
	case LDCH:
		*rA >>= 8
		*rA <<= 8
		*rA += valB
	case LDF:
		notImplemented("LDF")
	case LDL:
		m.registers[L] = valW
	case LDS:
		m.registers[S] = valW
	case LDT:
		m.registers[T] = valW
	case LDX:
		m.registers[X] = valW
	case LPS:
		// TODO load processor status
	case MUL:
		*rA *= valW
	case MULF:
		notImplemented("MULF")
	case OR:
		*rA |= valW
	case RD:
		*rA >>= 8
		*rA <<= 8
		*rA += int(m.devices[valW].Read())
	case RSUB:
		m.registers[PC] = m.registers[L]
	//case SSK:
	case STA:
		m.setWord(un, *rA)
	case STB:
		m.setWord(un, m.registers[B])
	case STCH:
		temp := *rA
		temp <<= 16
		temp >>= 16
		m.setByte(un, temp)
	case STF:
		m.setWord(un, m.registers[F])
	case STL:
		m.setWord(un, m.registers[L])
	case STS:
		m.setWord(un, m.registers[S])
	case STSW:
		m.setWord(un, m.registers[SW])
	case STT:
		m.setWord(un, m.registers[T])
	case STX:
		m.setWord(un, m.registers[X])
	case STI:
		// interval timer value <- valW
	case SUB:
		*rA -= valW
	case SUBF:
		notImplemented("SUBF")
	case TD:
		if m.devices[valW].Test() {
			m.registers[SW] = -1
		} else {
			m.registers[SW] = 0
		}
	case TIX:
		m.registers[X]++
		m.registers[X]++
		v1, v2 := m.registers[X], valW
		if v1 < v2 {
			m.registers[SW] = -1
		} else if v1 > v2 {
			m.registers[SW] = 1
		} else {
			m.registers[SW] = 0
		}
	case WD:
		temp := *rA
		temp <<= 16
		temp >>= 16

		m.devices[valW].Write(byte(temp))
	default:
		return false
	}
	return true
}

//func (m Machine) computeAddress(nixbpe, operand int) int {}

func (m Machine) getAddressFromInstruction(nixbpe, operand int) int { // PROBABLY WRONG
	var un int = 0
	var opLen int = 12 // 6 mest za opcode in 3 za nix

	ni := nixbpe >> 4

	xBit := (nixbpe & 0b1000) != 0
	eBit := (nixbpe & 0b0100) != 0
	pBit := (nixbpe & 0b0010) != 0
	bBit := (nixbpe & 0b0001) != 0
	if !eBit {
		opLen += 8 // 8 dodatnih mest, ki jih NIMAMO pri F3
	}

	un = operand

	if ni == 0b_00 { // SIC
		un |= (nixbpe << opLen)
	}

	if xBit {
		if ni == 0b_01 || ni == 0b_10 {
			invalidAddressing()
		} else {
			un += m.registers[X]
		}
	}

	if bBit {
		un += m.registers[B]
	}

	if pBit {
		un += m.registers[PC]
	}

	return un
}

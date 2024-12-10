package machine

import (
	"fmt"
	"time"
)

func notImplemented(mnemonic string) {}

func invalidOpcode(opcode int) {

}

func invalidAddressing() string {
	return "Invalid address."
}

func (m Machine) fetch() int {

	ind := m.registers[PC]
	m.registers[PC]++

	return m.GetByte(ind)

}

func signedDispToInt(word int) int {
	word &= 0x0FFF
	if word > 0x07FF {
		return -((^word) & 0x07FF)
	}
	return word
}

// func signedDispToInt(disp int) int {
// 	disp &= 0xFFF
// 	return disp - 2048
// }

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
			//fmt.Println(cnt)
			cnt++
			m.step()
			check = time.After(timeDuration)
		}
	}

}

func (m Machine) execute() string {
	opcode := m.fetch()

	if m.execF1(opcode) {
		return ""
	}

	op := m.fetch()

	if m.execF2(opcode, op) {
		return ""
	}

	flags := NewFlags(opcode, op)

	var operand int

	if flags.isSIC() { // SIC
		operand = op & 0b_01111111 // dodamo vse, razen x
		if flags.isIndexed() {
			operand += m.registers[X]
		}
	} else {
		operand = op & 0b_00001111 // dodamo vse razen xbpe

		if !flags.isExtended() { // F3
			operand = (operand << 8) | m.fetch()

			if flags.isPCRelative() && flags.isBaseRelative() {
				return invalidAddressing()
			}

			if flags.isPCRelative() {
				//fmt.Printf("operand: %d PC: %d\n", operand, m.registers[PC])
				if operand > MAX_PC_REL_ADDR {
					operand = -(^operand & MASK_PC_REL_ADDR) - 1
				}
				//fmt.Printf("operand: %d PC: %d\n", operand, m.registers[PC])

				operand += m.registers[PC]
			}

			if flags.isBaseRelative() {
				operand += m.registers[B]
			}

			if flags.isIndexed() && !flags.isSimple() {
				return invalidAddressing()
			}

			if flags.isIndexed() {
				operand += m.registers[X]
			}

		} else { // F4
			temp1 := m.fetch()
			temp2 := m.fetch()
			operand = (operand << 16) | temp1<<8 | temp2
			if flags.isRelative() {
				return invalidAddressing()
			}
			if !flags.isSimple() && flags.isIndexed() {
				return invalidAddressing()
			} else if flags.isIndexed() {
				operand += m.registers[X]
			}
		}
	}

	// // Pripravimo operand
	// temp := (op & ((1 << 4) - 1)) << 8
	// operand += temp

	// // Pripravimo nibble in opcode
	// nixbpe := op >> 4
	// temp = (opcode & 3) << 4
	// nixbpe += temp

	// opcode >>= 2
	// opcode <<= 2 //LATEST CHANGEEEEEE

	// if nixbpe%2 == 1 {
	// 	temp = m.fetch()
	// 	operand <<= 8
	// 	operand += temp
	// }

	//fmt.Printf("%s(%X) %d\n", InstructionMap[opcode], opcode, operand)

	opcode &= ^3

	m.execSICF3F4(opcode, flags, operand)
	return ""
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

func (m Machine) setSWForCompare(i1, i2 int) {
	temp := i1 - i2

	if temp < 0 {
		temp = -1
	} else if temp > 0 {
		temp = 1
	}

	m.setReg(SW, temp)
}

func (m Machine) execSICF3F4(opcode int, flags Flags, operand int) bool {

	var value, valueByte int

	if flags.isImmediate() {
		value = operand
		valueByte = operand
	} else if flags.isSIC() {
		value = m.GetWord(operand)
		valueByte = m.GetByte(operand)
	} else if flags.isSimple() {
		//fmt.Println("yes")
		value = m.GetWord(operand)
		valueByte = m.GetByte(operand)
	} else if flags.isIndirect() {
		operand = m.GetWord(operand)
		value = m.GetWord(operand)
		valueByte = m.GetByte(operand)

	}
	// DEBUG
	if opcode == MUL {
		fmt.Printf("ni: %b, xbpe: %b, value: %d, valueByte: %d, operand: %d\nisSimple: %b\n", flags.ni, flags.xbpe, value, valueByte, operand, flags.isSimple())
	}
	// DEBUG
	//rA := &m.registers[A]

	// Operand is now the address, value is the value

	switch opcode {
	// SETS -------------------------------------------
	case STA:
		m.setWord(operand, m.getReg(A))
	case STX:
		m.setWord(operand, m.getReg(X))
	case STL:
		m.setWord(operand, m.getReg(L))
	case STCH:
		m.setByte(operand, m.getAByte())
	case STB:
		m.setWord(operand, m.getReg(B))
	case STS:
		m.setWord(operand, m.getReg(S))
	case STF:
		m.setWord(operand, m.getReg(F))
	case STT:
		m.setWord(operand, m.getReg(T))
	case STSW:
		m.setWord(operand, m.getReg(SW))
	// JUMPS ---------------------------------------------
	case JEQ:
		if m.equalSW() {
			m.setReg(PC, operand)
		}
	case JGT:
		if m.greaterSW() {
			m.setReg(PC, operand)
		}
	case JLT:
		if m.lessSW() {
			m.setReg(PC, operand)
		}
	case J:
		m.setReg(PC, operand)
	case RSUB:
		m.setReg(PC, m.getReg(L))
	case JSUB:
		m.setReg(L, m.getReg(PC))
		m.setReg(PC, operand)
	// LOADS ----------------------------------
	case LDA:
		m.setReg(A, value)
	case LDX:
		m.setReg(X, value)
	case LDL:
		m.setReg(L, value)
	case LDCH:
		m.setAByte(valueByte)
	case LDB:
		m.setReg(B, value)
	case LDS:
		m.setReg(S, value)
	case LDF:
		notImplemented("LDF")
	case LDT:
		m.setReg(T, value)

	// ARITHMETIC ------------------------------------------
	case ADD:
		m.setReg(A, m.getReg(A)+value)
	case SUB:
		m.setReg(A, m.getReg(A)-value)
	case MUL:
		m.setReg(A, m.getReg(A)*value)
	case DIV:
		svalue := signedWordToInt(value)
		if svalue == 0 {
			fmt.Println("NaN encountered. (Division by zero)")
		} else {
			m.setReg(A, m.GetRegInt(A)/svalue)
		}
	case AND:
		m.setReg(A, m.getReg(A)&value)
	case OR:
		m.setReg(A, m.getReg(A)|value)
	case COMP:
		m.setSWForCompare(m.GetRegInt(A), signedWordToInt(value))
	case TIX:
		m.setReg(X, m.getReg(X)+1)
		m.setSWForCompare(m.GetRegInt(X), signedWordToInt(value))

	// IO ---------------------------------------------------
	case RD:
		m.setAByte(int(m.devices[valueByte].Read())) // device index is specified by byte at address "operand"
	case WD:
		m.devices[valueByte].Write(uint8(m.getAByte()))
	case TD:
		if m.devices[valueByte].Test() {
			m.setSWForCompare(-1, 0)
		} else {
			m.setSWForCompare(0, 0)
		}
	// FLOATS ----------------------------------------
	case ADDF:
		notImplemented("ADDF")
	case SUBF:
		notImplemented("SUBF")
	case MULF:
		notImplemented("MULF")
	case DIVF:
		notImplemented("DIVF")
	case COMPF:
		notImplemented("COMPF")

	// MISC --------------------------------------
	case LPS:
		notImplemented("LPS")
	case STI:
		notImplemented("STI")
	default:
		return false

	}

	// switch opcode {
	// case ADD:
	// 	m.registers[A] += signedWordToInt(value)
	// case AND:
	// 	m.registers[A] = signedWordToInt(intToSignedWord(m.registers[A]) & value)

	// case COMP:

	// 	v1, v2 := m.registers[A], signedWordToInt(value)
	// 	if v1 < v2 {
	// 		m.registers[SW] = -1
	// 	} else if v1 > v2 {
	// 		m.registers[SW] = 1
	// 	} else {
	// 		m.registers[SW] = 0
	// 	}
	// case COMPF:
	// 	notImplemented("COMPF")
	// case DIV:
	// 	*rA = *rA / signedWordToInt(value)
	// case DIVF:
	// 	notImplemented("DIVF")
	// case J:
	// 	//fmt.Printf("%s(%X) %d\n", InstructionMap[opcode], opcode, value)

	// 	//fmt.Printf("UN: %d, value: %d, PC: %d\n", operand, value, m.registers[PC])
	// 	m.registers[PC] = operand //value --- IF IT DOESNT WORK, REVERT ||| might have to change to setReg(...)
	// case JEQ:
	// 	if m.registers[SW] == 0 {
	// 		m.registers[PC] = operand //value --- IF IT DOESNT WORK, REVERT
	// 	}
	// case JGT:
	// 	if m.registers[SW] == 1 {
	// 		m.registers[PC] = operand //value --- IF IT DOESNT WORK, REVERT
	// 	}
	// case JLT:
	// 	if m.registers[SW] == -1 {
	// 		m.registers[PC] = operand //value --- IF IT DOESNT WORK, REVERT
	// 	}
	// case JSUB:
	// 	m.registers[L] = m.registers[PC]
	// 	m.registers[PC] = operand //value --- IF IT DOESNT WORK, REVERT
	// case LDA:
	// 	//*rA = value
	// 	m.setReg(A, value)
	// case LDB:
	// 	//m.registers[B] = value
	// 	m.setReg(B, value)
	// case LDCH:
	// 	*rA = signedWordToInt((intToSignedWord(*rA) & ^(0xFF)) | (value & 0xFF))
	// case LDF:
	// 	notImplemented("LDF")
	// case LDL:
	// 	//m.registers[L] = value
	// 	m.setReg(L, value)
	// case LDS:
	// 	//m.registers[S] = value
	// 	m.setReg(S, value)
	// case LDT:
	// 	//m.registers[T] = value
	// 	m.setReg(T, value)
	// case LDX:
	// 	m.registers[X] = value
	// case LPS:
	// 	// TODO load processor status
	// case MUL:
	// 	*rA *= signedWordToInt(value)
	// case MULF:
	// 	notImplemented("MULF")
	// case OR:
	// 	*rA = signedWordToInt(intToSignedWord(*rA) | value)
	// case RD:
	// 	*rA >>= 8
	// 	*rA <<= 8
	// 	*rA += int(m.devices[value].Read())
	// case RSUB:
	// 	m.registers[PC] = m.registers[L]
	// //case SSK:
	// case STA:
	// 	m.setWord(operand, intToSignedWord(*rA)) //value --- IF IT DOESNT WORK, REVERT
	// 	fmt.Printf("A: %d\n", m.GetReg(A))
	// case STB:
	// 	m.setWord(operand, m.registers[B]) //value --- IF IT DOESNT WORK, REVERT
	// case STCH:
	// 	temp := *rA
	// 	temp <<= 16
	// 	temp >>= 16
	// 	m.setByte(operand, temp) //value --- IF IT DOESNT WORK, REVERT
	// case STF:
	// 	m.setWord(operand, m.registers[F]) //value --- IF IT DOESNT WORK, REVERT
	// case STL:
	// 	m.setWord(operand, m.registers[L]) //value --- IF IT DOESNT WORK, REVERT
	// case STS:
	// 	m.setWord(operand, m.registers[S]) //value --- IF IT DOESNT WORK, REVERT
	// case STSW:
	// 	m.setWord(operand, m.registers[SW]) //value --- IF IT DOESNT WORK, REVERT
	// case STT:
	// 	m.setWord(operand, m.registers[T]) //value --- IF IT DOESNT WORK, REVERT
	// case STX:
	// 	m.setWord(operand, m.registers[X]) //value --- IF IT DOESNT WORK, REVERT
	// case STI:
	// 	// interval timer value <- value
	// 	notImplemented("STI")
	// case SUB:
	// 	*rA -= value
	// case SUBF:
	// 	notImplemented("SUBF")
	// case TD:
	// 	if m.devices[value].Test() {
	// 		m.registers[SW] = -1
	// 	} else {
	// 		m.registers[SW] = 0
	// 	}
	// case TIX:
	// 	m.registers[X]++
	// 	m.registers[X]++
	// 	v1, v2 := m.registers[X], value
	// 	if v1 < v2 {
	// 		m.registers[SW] = -1
	// 	} else if v1 > v2 {
	// 		m.registers[SW] = 1
	// 	} else {
	// 		m.registers[SW] = 0
	// 	}
	// case WD:
	// 	temp := *rA
	// 	temp <<= 16
	// 	temp >>= 16

	// 	m.devices[value].Write(byte(temp))
	// default:
	// 	return false
	// }
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

package machine

import (
	"bufio"
	"fmt"
	"os"

	"simulator/simUtils"
)

func (m Machine) LoadFile(file *os.File) bool {
	var r *bufio.Reader = bufio.NewReader(file)
	return m.loadSection(r)
}

func (m Machine) loadSection(r *bufio.Reader) bool {
	t, _, _ := r.ReadRune()

	if t != 'H' {
		fmt.Println("Doesnt start with H record.")
		return false
	}
	simUtils.ReadString(r, 6)
	codeAddress := simUtils.ReadWord(r)
	codeLength := simUtils.ReadWord(r)

	//fmt.Printf("codeAddress: %d, codeLength: %d\n", codeAddress, codeLength)

	// NE VEM ČE JE POTREBNO BRATI NEW LINE -- TESTIRAJ ... TESTIRAL - SE IZKAZE DA JE POTREBNO
	if temp, _, _ := r.ReadRune(); temp == '\r' {
		r.ReadRune()
	}

	t, _, _ = r.ReadRune()
	//fmt.Printf("t = %c\n", t)

	for t == 'T' {
		location := simUtils.ReadWord(r)
		len := simUtils.ReadByte(r)
		for i := 0; i < len; i++ {
			if location < codeAddress || location > codeLength {
				fmt.Println("Code address out of bounds.")
				return false
			}

			b := simUtils.ReadByte(r)
			m.setByte(location+i, b)

		}

		if temp, _, _ := r.ReadRune(); temp == '\r' {
			r.ReadRune()
		}

		t, _, _ = r.ReadRune()
		//fmt.Printf("t = %c\n", t)

	}

	for t == 'M' {
		simUtils.ReadString(r, 6)
		simUtils.ReadByte(r)
		//len := simUtils.ReadByte(r)
		//simUtils.ReadString(r, len)

		if temp, _, _ := r.ReadRune(); temp == '\r' {
			r.ReadRune()
		}

		t, _, _ = r.ReadRune()
		//fmt.Printf("t = %c\n", t)

	}

	if t != 'E' {
		fmt.Printf("Missing E record. (t = %c)\n", t)
		return false
	}

	start := simUtils.ReadWord(r)
	m.registers[PC] = start

	fmt.Println(m.registers)

	return true
}

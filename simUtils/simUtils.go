package simUtils

import (
	"bufio"
	"strconv"
)

func ReadString(r *bufio.Reader, len int) string {
	var res string = ""
	for ; len > 0; len-- {
		s, _, _ := r.ReadRune()
		res += string(s)

	}

	return res
}

func ReadByte(r *bufio.Reader) int {
	s := ReadString(r, 2)
	res, _ := strconv.ParseInt(s, 16, 32)
	return int(res)
}

func ReadWord(r *bufio.Reader) int {
	s := ReadString(r, 6)
	res, _ := strconv.ParseInt(s, 16, 32)
	return int(res)

}

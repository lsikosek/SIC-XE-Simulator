package main

import (
	"bufio"
	"os"
)

func ReadString(r bufio.Reader, len int) string {
	var res string = ""
	for ; len > 0; len-- {
		s, _, _ := r.ReadRune()
		res += string(s)

	}

	return res
}

func main() {
	w := bufio.NewWriter(os.Stdout)

	var c byte = 'a'
	w.WriteByte(c)

	w.Flush()

}

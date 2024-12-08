package main

import (
	"bufio"
	"fmt"
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
	r := bufio.NewReader(os.Stdin)
	t, _, _ := r.ReadRune()

	fmt.Print(t)

	str := ReadString(*r, 6)

	fmt.Print(str)

	if temp, _, _ := r.ReadRune(); temp == '\r' {
		r.ReadRune()
		fmt.Println("Found")
	}

}

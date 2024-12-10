package main

import (
	"fmt"
	"os"
	"simulator/machine"
)

func main() {
	var m machine.Machine = machine.NewMachine()

	if len(os.Args) < 2 {
		fmt.Println("No object file name provided.")
		return
	}

	//	fmt.Printf("DEBUG: %s\n", os.Args[1])

	f, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Printf("Could not open %s\n", f.Name())
		return
	} else if !m.LoadFile(f) {
		fmt.Printf("Could not load %s\n", f.Name())
	}

	for i := 0; i <= 100; i += 3 {
		fmt.Printf("%d ", m.GetWordInt(i))
	}

	m.Start()

	for m.IsRunning() {
	}
	//time.Sleep(5 * time.Second)

	fmt.Println(machine.RegisterNames)
	fmt.Println(m.GetRegistersRaw())

	for i := 0; i <= 100; i += 3 {
		fmt.Printf("%d ", m.GetWordInt(i))
	}

}

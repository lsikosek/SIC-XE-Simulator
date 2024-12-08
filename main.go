package main

import (
	"fmt"
	"os"
	"simulator/machine"
	"time"
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

	m.Start()

	time.Sleep(2 * time.Second)

	m.Stop()
	fmt.Println(machine.RegisterNames)
	fmt.Println(m.GetRegisters())

}

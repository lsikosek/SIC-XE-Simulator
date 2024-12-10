package main

import (
	"flag"
	"fmt"
	"os"
	"simulator/machine"
)

func main() {
	var m machine.Machine = machine.NewMachine()

	var sFlag = flag.Int("s", 100, "Sets the clock speed in Hz")
	var mFlag = flag.Int("m", 0, "Prints memory up to a certain address")
	flag.Parse()

	m.SetSpeed(*sFlag)

	if len(flag.Args()) < 1 {
		fmt.Println("No object file name provided.")
		return
	}

	//	fmt.Printf("DEBUG: %s\n", os.Args[1])

	f, err := os.Open(flag.Arg(0))

	if err != nil {
		fmt.Printf("Could not open %s\n", f.Name())
		return
	} else if !m.LoadFile(f) {
		fmt.Printf("Could not load %s\n", f.Name())
		return
	} else {
		fmt.Printf("File %s loaded. Starting...\n", f.Name())
	}

	// for i := 0; i < *mFlag; i += 3 {
	// 	//fmt.Printf("%d ", m.GetWordInt(i))
	// }

	m.Start()

	for m.IsRunning() {
	}
	//time.Sleep(5 * time.Second)

	fmt.Println()
	fmt.Println(machine.RegisterNames)
	fmt.Println(m.GetRegistersRaw())

	for i := 0; i < *mFlag; i += 3 {
		if ((i / 3) % 6) == 0 {

			fmt.Println()
			fmt.Printf("[%8X - %8X] | ", i, i+5)
		}
		fmt.Printf("%d ", m.GetWordInt(i))

	}

}

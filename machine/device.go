package machine

import (
	"bufio"
	"os"
)

type DeviceInterface interface {
	Test() bool
	Read() uint8
	Write()
}

type Device struct {
	output *bufio.Writer
	input  *bufio.Reader
}

func (d Device) Write(b uint8) {
	//fmt.Printf("WROTE: %c\n", b)
	//fmt.Fprintf(d.output, "%c\n", rune(b))

	d.output.WriteByte(byte(b))
	d.output.Flush()
}

func (d Device) Read() uint8 {

	b, err := d.input.ReadByte()
	if err == nil {
		//fmt.Printf("READ: %c\n", b)
		return byte(b)
	}
	return 0
}

func (d Device) Test() bool {
	return (d.output != nil || d.input != nil)
}

func newInputDevice() Device {
	return Device{
		output: nil,
		input:  bufio.NewReader(os.Stdin),
	}
}

func newOutputDevice() Device {
	return Device{
		output: bufio.NewWriter(os.Stdout),
		input:  nil,
	}
}

func newErrorDevice() Device {
	return Device{
		output: bufio.NewWriter(os.Stderr),
		input:  nil,
	}
}

func newDevice(file *os.File) Device {
	return Device{
		output: bufio.NewWriter(file),
		input:  bufio.NewReader(file),
	}
}

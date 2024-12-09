package machine

const (
	NONE         int = 0x00
	MASK_NI      int = 0x03
	SIC          int = 0x00
	IMMEDIATE    int = 0x01
	INDIRECT     int = 0x02
	SIMPLE       int = 0x03
	MASK_XBPE    int = 0xF0
	MASK_XBP     int = 0xE0
	MASK_BP      int = 0x60
	INDEXED      int = 0x80
	BASERELATIVE int = 0x40
	PCRELATIVE   int = 0x20
	EXTENDED     int = 0x10

	MAX_PC_REL_ADDR  int = 2047
	MASK_PC_REL_ADDR int = 0x7FF // 0111 1111 1111 ...
)

type Flags struct {
	ni   int
	xbpe int
}

func NewFlags(opcode, op int) Flags {
	return Flags{
		opcode & 0x03,
		op & 0xF0,
	}
}

func (f Flags) isExtended() bool {
	return (f.xbpe&EXTENDED != 0)
}

func (f Flags) isPCRelative() bool {
	return (f.xbpe&PCRELATIVE != 0)
}

func (f Flags) isBaseRelative() bool {
	return (f.xbpe&BASERELATIVE != 0)
}

func (f Flags) isRelative() bool {
	return f.isBaseRelative() || f.isPCRelative()
}

func (f Flags) isIndexed() bool {
	return (f.xbpe&INDEXED != 0)
}

func (f Flags) isSIC() bool {
	return (f.ni&SIC != 0)
}

func (f Flags) isSimple() bool {
	return (f.ni&SIMPLE != 0)
}

func (f Flags) isIndirect() bool {
	return (f.ni&INDIRECT != 0)
}

func (f Flags) isImmediate() bool {
	return (f.ni&IMMEDIATE != 0)
}

func (f Flags) getOperandF3(op int, disp int) {
	return
}

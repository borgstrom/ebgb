package cpu

// CPU implements the 8-bit Sharp LR35902
// See: https://gbdev.io/pandocs
//
// There is no explicit reset function, to reset the CPU just create a new instance.
type CPU struct {
	// 8-bit Registers
	a, b, c, d, e, h, l, f uint8

	// 16-bit Registers
	pc, sp uint16

	ram Memory
}

type Memory interface {
	Read(a uint16) byte
	Write(a uint16, v byte)
}

func New(ram Memory) *CPU {
	return &CPU{
		pc:  0x100,
		sp:  0xfffe,
		ram: ram,
	}
}

func (c *CPU) Tick() {
}

// instruction is a function that takes a CPU pointer and returns the number of cycles taken to execute
type instruction func(*CPU) int

var instructions = map[byte]instruction{
	// NOP
	0x00: func(c *CPU) int {
		return 1
	},

	// STOP
	0x10: func(c *CPU) int {
		// TODO: read another byte
		panic("STOP not implemented yet")
	},

	// LD B, X
	0x40: func(c *CPU) int {
		// c.b = c.b
		return 1
	},
	0x41: func(c *CPU) int {
		c.b = c.c
		return 1
	},
	0x42: func(c *CPU) int {
		c.b = c.d
		return 1
	},
	0x43: func(c *CPU) int {
		c.b = c.e
		return 1
	},
	0x44: func(c *CPU) int {
		c.b = c.h
		return 1
	},
	0x45: func(c *CPU) int {
		c.b = c.l
		return 1
	},
	0x47: func(c *CPU) int {
		c.b = c.a
		return 1
	},

	// LD C, X
	0x48: func(c *CPU) int {
		c.c = c.b
		return 1
	},
	0x49: func(c *CPU) int {
		// c.c = c.c
		return 1
	},
	0x4a: func(c *CPU) int {
		c.c = c.d
		return 1
	},
	0x4b: func(c *CPU) int {
		c.c = c.e
		return 1
	},
	0x4c: func(c *CPU) int {
		c.c = c.h
		return 1
	},
	0x4d: func(c *CPU) int {
		c.c = c.l
		return 1
	},
	0x4f: func(c *CPU) int {
		c.c = c.a
		return 1
	},

	// LD D, X
	0x50: func(c *CPU) int {
		c.d = c.b
		return 1
	},
	0x51: func(c *CPU) int {
		c.d = c.c
		return 1
	},
	0x52: func(c *CPU) int {
		// c.d = c.d
		return 1
	},
	0x53: func(c *CPU) int {
		c.d = c.e
		return 1
	},
	0x54: func(c *CPU) int {
		c.d = c.h
		return 1
	},
	0x55: func(c *CPU) int {
		c.d = c.l
		return 1
	},
	0x57: func(c *CPU) int {
		c.d = c.a
		return 1
	},

	// LD E, X
	0x58: func(c *CPU) int {
		c.e = c.b
		return 1
	},
	0x59: func(c *CPU) int {
		c.e = c.c
		return 1
	},
	0x5a: func(c *CPU) int {
		c.e = c.d
		return 1
	},
	0x5b: func(c *CPU) int {
		// c.e = c.e
		return 1
	},
	0x5c: func(c *CPU) int {
		c.e = c.h
		return 1
	},
	0x5d: func(c *CPU) int {
		c.e = c.l
		return 1
	},
	0x5f: func(c *CPU) int {
		c.e = c.a
		return 1
	},

	// LD H, X
	0x60: func(c *CPU) int {
		c.h = c.b
		return 1
	},
	0x61: func(c *CPU) int {
		c.h = c.c
		return 1
	},
	0x62: func(c *CPU) int {
		c.h = c.d
		return 1
	},
	0x63: func(c *CPU) int {
		c.h = c.e
		return 1
	},
	0x64: func(c *CPU) int {
		// c.h = c.h
		return 1
	},
	0x65: func(c *CPU) int {
		c.h = c.l
		return 1
	},
	0x67: func(c *CPU) int {
		c.h = c.a
		return 1
	},

	// LD L, X
	0x68: func(c *CPU) int {
		c.l = c.b
		return 1
	},
	0x69: func(c *CPU) int {
		c.l = c.c
		return 1
	},
	0x6a: func(c *CPU) int {
		c.l = c.d
		return 1
	},
	0x6b: func(c *CPU) int {
		c.l = c.e
		return 1
	},
	0x6c: func(c *CPU) int {
		c.l = c.h
		return 1
	},
	0x6d: func(c *CPU) int {
		// c.l = c.l
		return 1
	},
	0x6f: func(c *CPU) int {
		c.l = c.a
		return 1
	},

	// LD A, X
	0x78: func(c *CPU) int {
		c.a = c.b
		return 1
	},
	0x79: func(c *CPU) int {
		c.a = c.c
		return 1
	},
	0x7a: func(c *CPU) int {
		c.a = c.d
		return 1
	},
	0x7b: func(c *CPU) int {
		c.a = c.e
		return 1
	},
	0x7c: func(c *CPU) int {
		c.a = c.h
		return 1
	},
	0x7d: func(c *CPU) int {
		c.a = c.l
		return 1
	},
	0x7f: func(c *CPU) int {
		// c.a = c.a
		return 1
	},

	// LD X, (HL)
	0x46: func(c *CPU) int {
		c.b = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},
	0x4e: func(c *CPU) int {
		c.c = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},
	0x56: func(c *CPU) int {
		c.d = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},
	0x5e: func(c *CPU) int {
		c.e = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},
	0x66: func(c *CPU) int {
		c.h = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},
	0x6e: func(c *CPU) int {
		c.l = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},
	0x76: func(c *CPU) int {
		c.a = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},

	// LD (HL), X
	0x70: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.b)
		return 2
	},
	0x71: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.c)
		return 2
	},
	0x72: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.d)
		return 2
	},
	0x73: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.e)
		return 2
	},
	0x74: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.h)
		return 2
	},
	0x75: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.l)
		return 2
	},
	0x77: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.a)
		return 2
	},

	// LD B, d8
	0x06: func(c *CPU) int {
		c.b = c.ram.Read(c.pc)
		c.pc++
		return 2
	},
	0x0e: func(c *CPU) int {
		c.c = c.ram.Read(c.pc)
		c.pc++
		return 2
	},
	0x16: func(c *CPU) int {
		c.d = c.ram.Read(c.pc)
		c.pc++
		return 2
	},
	0x1e: func(c *CPU) int {
		c.e = c.ram.Read(c.pc)
		c.pc++
		return 2
	},
	0x26: func(c *CPU) int {
		c.h = c.ram.Read(c.pc)
		c.pc++
		return 2
	},
	0x2e: func(c *CPU) int {
		c.l = c.ram.Read(c.pc)
		c.pc++
		return 2
	},
	0x3e: func(c *CPU) int {
		c.a = c.ram.Read(c.pc)
		c.pc++
		return 2
	},

	// LD (HL), d8
	0x36: func(c *CPU) int {
		c.ram.Write(bb2i(c.h, c.l), c.ram.Read(c.pc))
		c.pc++
		return 3
	},

	// LD (BC), A
	0x02: func(c *CPU) int {
		c.ram.Write(bb2i(c.b, c.c), c.a)
		return 2
	},

	// LD (DE), A
	0x12: func(c *CPU) int {
		c.ram.Write(bb2i(c.d, c.e), c.a)
		return 2
	},

	// LD (HL+), A
	0x22: func(c *CPU) int {
		hl := bb2i(c.h, c.l)
		c.ram.Write(hl, c.a)
		hl++
		c.h, c.l = i2bb(hl)
		return 2
	},

	// LD (HL-), A
	0x32: func(c *CPU) int {
		hl := bb2i(c.h, c.l)
		c.ram.Write(hl, c.a)
		hl--
		c.h, c.l = i2bb(hl)
		return 2
	},

	// LD A, (BC)
	0x0a: func(c *CPU) int {
		c.a = c.ram.Read(bb2i(c.b, c.c))
		return 2
	},

	// LD A, (DE)
	0x1a: func(c *CPU) int {
		c.a = c.ram.Read(bb2i(c.d, c.e))
		return 2
	},

	// LD A, (HL)
	0x7e: func(c *CPU) int {
		c.a = c.ram.Read(bb2i(c.h, c.l))
		return 2
	},

	// LD A, (HL+)
	0x2a: func(c *CPU) int {
		hl := bb2i(c.h, c.l)
		c.a = c.ram.Read(hl)
		hl++
		c.h, c.l = i2bb(hl)
		return 2
	},

	// LD A, (HL-)
	0x3a: func(c *CPU) int {
		hl := bb2i(c.h, c.l)
		c.a = c.ram.Read(hl)
		hl--
		c.h, c.l = i2bb(hl)
		return 2
	},
}

// bb2i converts two separate byte into an unsigned 16-bit integer
func bb2i(a, b byte) uint16 {
	return uint16(a)<<8 | uint16(b)
}

// i2bb converts an unsigned 16-bit integer into two bytes
func i2bb(v uint16) (byte, byte) {
	return byte(v >> 8), byte(v & 0xff)
}

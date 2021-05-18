package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// CPU implements the 8-bit Sharp LR35902
// See: https://gbdev.io/pandocs
//
// There is no explicit reset function, to reset the CPU just create a new instance.
type CPU struct {
	// 8-bit Registers
	a, b, c, d, e, h, l, f uint8

	// 16-bit Registers
	pc, sp uint16

	// Flags
	fz, fn, fh, fc bool

	cycles uint64

	ram Memory
}

type Memory interface {
	Read(a uint16) uint8
	Write(a uint16, v uint8)
}

func New(ram Memory) *CPU {
	return &CPU{
		a: 0x01,
		b: 0x00,
		c: 0x13,
		d: 0x00,
		e: 0xdb,
		h: 0x01,
		l: 0x4d,

		sp: 0xfffe,
		pc: 0x0000,

		fz: true,
		fn: false,
		fh: true,
		fc: true,

		cycles: 0,

		ram: ram,
	}
}

// Next runs a single iteration of the CPU and returns the number of cycles taken
func (c *CPU) Next() uint8 {
	opCode := c.PC()
	cycles := instructions[opCode].execute(c)
	c.cycles = c.cycles + uint64(cycles)
	return cycles
}

func (c *CPU) PC() uint8 {
	v := c.ram.Read(c.pc)
	c.pc++
	return v
}

// TODO move this somewhere
type testRAM []uint8

func (r testRAM) Read(a uint16) uint8 {
	return r[a]
}

func (r testRAM) Write(a uint16, v uint8) {
	r[a] = v
}

// instructionFunc is a function that takes a CPU pointer and returns the number of cycles taken to execute
type instructionFunc func(c *CPU) uint8

type instructionTestFunc func(t *testing.T, c *CPU)

type instruction struct {
	mnemonic   string
	execute    instructionFunc
	test       instructionTestFunc
	testMemory testRAM
}

var instructions map[uint8]*instruction

func init() {
	instructions = map[uint8]*instruction{
		0x00: {
			mnemonic: "NOP",
			execute: func(c *CPU) uint8 {
				return 1
			},
			testMemory: testRAM{0x00},
			test: func(t *testing.T, c *CPU) {
				c.Next()
				require.Equal(t, uint64(1), c.cycles)
			},
		},
		0x01: {
			mnemonic: "LD BC, d16A",
			execute: func(c *CPU) uint8 {
				c.c = c.PC()
				c.b = c.PC()
				return 3
			},
			testMemory: testRAM{0x01, 0x7b, 0x00},
			test: func(t *testing.T, c *CPU) {
				c.Next()
				require.Equal(t, uint16(123), bb2i(c.b, c.c))
				require.Equal(t, uint64(3), c.cycles)
			},
		},

		0x02: {
			mnemonic: "LD (BC), A",
			execute: func(c *CPU) uint8 {
				c.ram.Write(bb2i(c.b, c.c), c.a)
				return 2
			},
			testMemory: testRAM{0x02, 0x00},
			test: func(t *testing.T, c *CPU) {
				c.b = 0x0
				c.c = 0x1
				c.a = 0xff
				c.Next()
				require.Equal(t, uint8(0xff), c.ram.Read(0x0001))
				require.Equal(t, uint64(2), c.cycles)
			},
		},
		//
		//// INC BC
		//0x03: func(c *CPU) uint8 {
		//	bc := bb2i(c.b, c.c)
		//	c.b, c.c = i2bb(bc + 1)
		//	return 2
		//},
		//
		//// INC B
		//0x04: func(c *CPU) uint8 {
		//	c.b++
		//	c.fz = c.b == 0x00
		//	c.fn = false
		//	c.fh = (c.b & 0x0f) == 0x00
		//	return 1
		//},
		//
		//// DEC B
		//0x05: func(c *CPU) uint8 {
		//	c.b--
		//	c.fz = c.b == 0x00
		//	c.fn = false
		//	c.fh = (c.b & 0x0f) == 0x00
		//	return 1
		//},
		//
		//// LD B, d8
		//0x06: func(c *CPU) uint8 {
		//	c.b = c.PC()
		//	return 1
		//},
		//
		//// RLCA
		//0x07: func(c *CPU) uint8 {
		//	// R_A = (R_A << 1) | (R_A >> 7);
		//	// F_Z = 0;
		//	// F_N = 0;
		//	// F_H = 0;
		//	// F_C = (R_A & 0x01);
		//	c.a = c.a<<1 | c.a>>7
		//	c.fz = false
		//	c.fn = false
		//	c.fh = false
		//	c.fc = c.a&0x01 != 0x00
		//	return 1
		//},
		//
		//// LD (a16), SP
		//0x08: func(c *CPU) uint8 {
		//	a := bb2i(c.PC(), c.PC())
		//	b1, b2 := i2bb(c.sp)
		//	c.ram.Write(a, b1)
		//	c.ram.Write(a+1, b2)
		//	return 5
		//},
		//
		//// ADD HL, BC
		//0x09: func(c *CPU) uint8 {
		//	bc := bb2i(c.b, c.c)
		//	hl := bb2i(c.h, c.l)
		//	v := bc + hl
		//	c.h, c.l = i2bb(v)
		//
		//	c.fz = v == 0
		//	c.fc = v > 0xffff
		//	c.fn = false
		//	c.fh = v&0x10 == 0x10
		//	return 2
		//},
		//
		//// LD A, (BC)
		//0x0a: func(c *CPU) uint8 {
		//	c.a = c.ram.Read(bb2i(c.b, c.c))
		//	return 2
		//},
		//
		//// DEC BC
		//0x0b: func(c *CPU) uint8 {
		//	v := bb2i(c.b, c.c)
		//	c.b, c.c = i2bb(v - 1)
		//	return 2
		//},
		//
		//// STOP
		//0x10: func(c *CPU) uint8 {
		//	// TODO: read another uint8, which should be 0x00
		//	panic("STOP not implemented yet")
		//},
		//
		//// LD B, X
		//0x40: func(c *CPU) uint8 {
		//	// c.b = c.b
		//	return 1
		//},
		//0x41: func(c *CPU) uint8 {
		//	c.b = c.c
		//	return 1
		//},
		//0x42: func(c *CPU) uint8 {
		//	c.b = c.d
		//	return 1
		//},
		//0x43: func(c *CPU) uint8 {
		//	c.b = c.e
		//	return 1
		//},
		//0x44: func(c *CPU) uint8 {
		//	c.b = c.h
		//	return 1
		//},
		//0x45: func(c *CPU) uint8 {
		//	c.b = c.l
		//	return 1
		//},
		//0x47: func(c *CPU) uint8 {
		//	c.b = c.a
		//	return 1
		//},
		//
		//// LD C, X
		//0x48: func(c *CPU) uint8 {
		//	c.c = c.b
		//	return 1
		//},
		//0x49: func(c *CPU) uint8 {
		//	// c.c = c.c
		//	return 1
		//},
		//0x4a: func(c *CPU) uint8 {
		//	c.c = c.d
		//	return 1
		//},
		//0x4b: func(c *CPU) uint8 {
		//	c.c = c.e
		//	return 1
		//},
		//0x4c: func(c *CPU) uint8 {
		//	c.c = c.h
		//	return 1
		//},
		//0x4d: func(c *CPU) uint8 {
		//	c.c = c.l
		//	return 1
		//},
		//0x4f: func(c *CPU) uint8 {
		//	c.c = c.a
		//	return 1
		//},
		//
		//// LD D, X
		//0x50: func(c *CPU) uint8 {
		//	c.d = c.b
		//	return 1
		//},
		//0x51: func(c *CPU) uint8 {
		//	c.d = c.c
		//	return 1
		//},
		//0x52: func(c *CPU) uint8 {
		//	// c.d = c.d
		//	return 1
		//},
		//0x53: func(c *CPU) uint8 {
		//	c.d = c.e
		//	return 1
		//},
		//0x54: func(c *CPU) uint8 {
		//	c.d = c.h
		//	return 1
		//},
		//0x55: func(c *CPU) uint8 {
		//	c.d = c.l
		//	return 1
		//},
		//0x57: func(c *CPU) uint8 {
		//	c.d = c.a
		//	return 1
		//},
		//
		//// LD E, X
		//0x58: func(c *CPU) uint8 {
		//	c.e = c.b
		//	return 1
		//},
		//0x59: func(c *CPU) uint8 {
		//	c.e = c.c
		//	return 1
		//},
		//0x5a: func(c *CPU) uint8 {
		//	c.e = c.d
		//	return 1
		//},
		//0x5b: func(c *CPU) uint8 {
		//	// c.e = c.e
		//	return 1
		//},
		//0x5c: func(c *CPU) uint8 {
		//	c.e = c.h
		//	return 1
		//},
		//0x5d: func(c *CPU) uint8 {
		//	c.e = c.l
		//	return 1
		//},
		//0x5f: func(c *CPU) uint8 {
		//	c.e = c.a
		//	return 1
		//},
		//
		//// LD H, X
		//0x60: func(c *CPU) uint8 {
		//	c.h = c.b
		//	return 1
		//},
		//0x61: func(c *CPU) uint8 {
		//	c.h = c.c
		//	return 1
		//},
		//0x62: func(c *CPU) uint8 {
		//	c.h = c.d
		//	return 1
		//},
		//0x63: func(c *CPU) uint8 {
		//	c.h = c.e
		//	return 1
		//},
		//0x64: func(c *CPU) uint8 {
		//	// c.h = c.h
		//	return 1
		//},
		//0x65: func(c *CPU) uint8 {
		//	c.h = c.l
		//	return 1
		//},
		//0x67: func(c *CPU) uint8 {
		//	c.h = c.a
		//	return 1
		//},
		//
		//// LD L, X
		//0x68: func(c *CPU) uint8 {
		//	c.l = c.b
		//	return 1
		//},
		//0x69: func(c *CPU) uint8 {
		//	c.l = c.c
		//	return 1
		//},
		//0x6a: func(c *CPU) uint8 {
		//	c.l = c.d
		//	return 1
		//},
		//0x6b: func(c *CPU) uint8 {
		//	c.l = c.e
		//	return 1
		//},
		//0x6c: func(c *CPU) uint8 {
		//	c.l = c.h
		//	return 1
		//},
		//0x6d: func(c *CPU) uint8 {
		//	// c.l = c.l
		//	return 1
		//},
		//0x6f: func(c *CPU) uint8 {
		//	c.l = c.a
		//	return 1
		//},
		//
		//// LD A, X
		//0x78: func(c *CPU) uint8 {
		//	c.a = c.b
		//	return 1
		//},
		//0x79: func(c *CPU) uint8 {
		//	c.a = c.c
		//	return 1
		//},
		//0x7a: func(c *CPU) uint8 {
		//	c.a = c.d
		//	return 1
		//},
		//0x7b: func(c *CPU) uint8 {
		//	c.a = c.e
		//	return 1
		//},
		//0x7c: func(c *CPU) uint8 {
		//	c.a = c.h
		//	return 1
		//},
		//0x7d: func(c *CPU) uint8 {
		//	c.a = c.l
		//	return 1
		//},
		//0x7f: func(c *CPU) uint8 {
		//	// c.a = c.a
		//	return 1
		//},
		//
		//// LD X, (HL)
		//0x46: func(c *CPU) uint8 {
		//	c.b = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//0x4e: func(c *CPU) uint8 {
		//	c.c = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//0x56: func(c *CPU) uint8 {
		//	c.d = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//0x5e: func(c *CPU) uint8 {
		//	c.e = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//0x66: func(c *CPU) uint8 {
		//	c.h = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//0x6e: func(c *CPU) uint8 {
		//	c.l = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//0x76: func(c *CPU) uint8 {
		//	c.a = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//
		//// LD (HL), X
		//0x70: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.b)
		//	return 2
		//},
		//0x71: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.c)
		//	return 2
		//},
		//0x72: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.d)
		//	return 2
		//},
		//0x73: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.e)
		//	return 2
		//},
		//0x74: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.h)
		//	return 2
		//},
		//0x75: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.l)
		//	return 2
		//},
		//0x77: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.a)
		//	return 2
		//},
		//
		//// LD X, d8
		//0x0e: func(c *CPU) uint8 {
		//	c.c = c.PC()
		//	return 2
		//},
		//0x16: func(c *CPU) uint8 {
		//	c.d = c.PC()
		//	return 2
		//},
		//0x1e: func(c *CPU) uint8 {
		//	c.e = c.PC()
		//	return 2
		//},
		//0x26: func(c *CPU) uint8 {
		//	c.h = c.PC()
		//	return 2
		//},
		//0x2e: func(c *CPU) uint8 {
		//	c.l = c.PC()
		//	return 2
		//},
		//0x3e: func(c *CPU) uint8 {
		//	c.a = c.PC()
		//	return 2
		//},
		//
		//// LD (HL), d8
		//0x36: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.h, c.l), c.PC())
		//	return 3
		//},
		//
		//// LD (DE), A
		//0x12: func(c *CPU) uint8 {
		//	c.ram.Write(bb2i(c.d, c.e), c.a)
		//	return 2
		//},
		//
		//// LD (HL+), A
		//0x22: func(c *CPU) uint8 {
		//	hl := bb2i(c.h, c.l)
		//	c.ram.Write(hl, c.a)
		//	hl++
		//	c.h, c.l = i2bb(hl)
		//	return 2
		//},
		//
		//// LD (HL-), A
		//0x32: func(c *CPU) uint8 {
		//	hl := bb2i(c.h, c.l)
		//	c.ram.Write(hl, c.a)
		//	hl--
		//	c.h, c.l = i2bb(hl)
		//	return 2
		//},
		//
		//// LD A, (DE)
		//0x1a: func(c *CPU) uint8 {
		//	c.a = c.ram.Read(bb2i(c.d, c.e))
		//	return 2
		//},
		//
		//// LD A, (HL)
		//0x7e: func(c *CPU) uint8 {
		//	c.a = c.ram.Read(bb2i(c.h, c.l))
		//	return 2
		//},
		//
		//// LD A, (HL+)
		//0x2a: func(c *CPU) uint8 {
		//	hl := bb2i(c.h, c.l)
		//	c.a = c.ram.Read(hl)
		//	hl++
		//	c.h, c.l = i2bb(hl)
		//	return 2
		//},
		//
		//// LD A, (HL-)
		//0x3a: func(c *CPU) uint8 {
		//	hl := bb2i(c.h, c.l)
		//	c.a = c.ram.Read(hl)
		//	hl--
		//	c.h, c.l = i2bb(hl)
		//	return 2
		//},
	}
}

// bb2i converts two separate uint8 into an unsigned 16-bit integer
func bb2i(a, b uint8) uint16 {
	return uint16(a)<<8 | uint16(b)
}

// i2bb converts an unsigned 16-bit integer into two bytes
func i2bb(v uint16) (uint8, uint8) {
	return uint8(v >> 8), uint8(v & 0xff)
}

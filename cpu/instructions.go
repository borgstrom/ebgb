package cpu

// instructionFunc is a function that takes a CPU pointer and returns the number of cycles taken to execute
type instructionFunc func(c *CPU) uint8

// instructionsByOpcode allows us to layout our op codes in a table format
// See: https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html
var instructionsByOpcode = [][]instructionFunc{
	// 0x00 - 0x0f
	{nop, ldBcD16, ldBcA, incBc, incB, decB, ldBD8}, // rlca, ldA16Sp, addHlBc, ldABc, decBc, incC, decC, ldCD8, rrca},
	// 0x01 - 0x1f
	//{STOP, LD_DE_d16, LD_DE_A, INC_DE, INC_D, DEC_D, LD_D_d8, RLA, JR_s8, ADD_HL_DE, LD_A_DE, DEC_DE, INC_E, DEC_E, LD_E_d8, RRA},
}

// ----- common helpers -----

func inc(c *CPU, get func() uint8, set func(v uint8)) {
	v := get() + 1
	set(v)

	if c.isFlagSet(flagCarry) {
		c.initFlags(flagCarry)
	} else {
		c.initFlags(flagNone)
	}

	if v == 0 {
		c.enableFlag(flagZero)
	}

	if v&0x0f == 0x00 {
		c.enableFlag(flagHalfCarry)
	}
}

func dec(c *CPU, get func() uint8, set func(v uint8)) {
	v := get() - 1
	set(v)

	if c.isFlagSet(flagCarry) {
		c.initFlags(flagCarry)
	} else {
		c.initFlags(flagNone)
	}

	if v == 0 {
		c.enableFlag(flagZero)
	}

	c.enableFlag(flagSubtraction)

	if v&0x0f == 0x0f {
		c.enableFlag(flagHalfCarry)
	}
}

// ----- instructions -----

func nop(c *CPU) uint8 {
	return 1
}

func ldBcD16(c *CPU) uint8 {
	c.bc.SetHigh(c.PC())
	c.bc.SetLow(c.PC())
	return 3
}

func ldBcA(c *CPU) uint8 {
	c.ram.Write(uint16(c.bc), c.af.GetHigh())
	return 2
}

func incBc(c *CPU) uint8 {
	c.bc++
	return 2
}

func incB(c *CPU) uint8 {
	inc(c, c.bc.GetHigh, c.bc.SetHigh)
	return 1
}

func decB(c *CPU) uint8 {
	dec(c, c.bc.GetHigh, c.bc.SetHigh)
	return 1
}

func ldBD8(c *CPU) uint8 {
	c.bc.SetHigh(c.PC())
	return 1
}

//		//
//		//// LD B, d8
//		//0x06: func(c *CPU) uint8 {
//		//	c.b = c.PC()
//		//	return 1
//		//},
//		//
//		//// RLCA
//		//0x07: func(c *CPU) uint8 {
//		//	// R_A = (R_A << 1) | (R_A >> 7);
//		//	// F_Z = 0;
//		//	// F_N = 0;
//		//	// F_H = 0;
//		//	// F_C = (R_A & 0x01);
//		//	c.a = c.a<<1 | c.a>>7
//		//	c.fz = false
//		//	c.fn = false
//		//	c.fh = false
//		//	c.fc = c.a&0x01 != 0x00
//		//	return 1
//		//},
//		//
//		//// LD (a16), SP
//		//0x08: func(c *CPU) uint8 {
//		//	a := bb2i(c.PC(), c.PC())
//		//	b1, b2 := i2bb(c.sp)
//		//	c.ram.Write(a, b1)
//		//	c.ram.Write(a+1, b2)
//		//	return 5
//		//},
//		//
//		//// ADD HL, BC
//		//0x09: func(c *CPU) uint8 {
//		//	bc := bb2i(c.b, c.c)
//		//	hl := bb2i(c.h, c.l)
//		//	v := bc + hl
//		//	c.h, c.l = i2bb(v)
//		//
//		//	c.fz = v == 0
//		//	c.fc = v > 0xffff
//		//	c.fn = false
//		//	c.fh = v&0x10 == 0x10
//		//	return 2
//		//},
//		//
//		//// LD A, (BC)
//		//0x0a: func(c *CPU) uint8 {
//		//	c.a = c.ram.Read(bb2i(c.b, c.c))
//		//	return 2
//		//},
//		//
//		//// DEC BC
//		//0x0b: func(c *CPU) uint8 {
//		//	v := bb2i(c.b, c.c)
//		//	c.b, c.c = i2bb(v - 1)
//		//	return 2
//		//},
//		//
//		//// STOP
//		//0x10: func(c *CPU) uint8 {
//		//	// TODO: read another uint8, which should be 0x00
//		//	panic("STOP not implemented yet")
//		//},
//		//
//		//// LD B, X
//		//0x40: func(c *CPU) uint8 {
//		//	// c.b = c.b
//		//	return 1
//		//},
//		//0x41: func(c *CPU) uint8 {
//		//	c.b = c.c
//		//	return 1
//		//},
//		//0x42: func(c *CPU) uint8 {
//		//	c.b = c.d
//		//	return 1
//		//},
//		//0x43: func(c *CPU) uint8 {
//		//	c.b = c.e
//		//	return 1
//		//},
//		//0x44: func(c *CPU) uint8 {
//		//	c.b = c.h
//		//	return 1
//		//},
//		//0x45: func(c *CPU) uint8 {
//		//	c.b = c.l
//		//	return 1
//		//},
//		//0x47: func(c *CPU) uint8 {
//		//	c.b = c.a
//		//	return 1
//		//},
//		//
//		//// LD C, X
//		//0x48: func(c *CPU) uint8 {
//		//	c.c = c.b
//		//	return 1
//		//},
//		//0x49: func(c *CPU) uint8 {
//		//	// c.c = c.c
//		//	return 1
//		//},
//		//0x4a: func(c *CPU) uint8 {
//		//	c.c = c.d
//		//	return 1
//		//},
//		//0x4b: func(c *CPU) uint8 {
//		//	c.c = c.e
//		//	return 1
//		//},
//		//0x4c: func(c *CPU) uint8 {
//		//	c.c = c.h
//		//	return 1
//		//},
//		//0x4d: func(c *CPU) uint8 {
//		//	c.c = c.l
//		//	return 1
//		//},
//		//0x4f: func(c *CPU) uint8 {
//		//	c.c = c.a
//		//	return 1
//		//},
//		//
//		//// LD D, X
//		//0x50: func(c *CPU) uint8 {
//		//	c.d = c.b
//		//	return 1
//		//},
//		//0x51: func(c *CPU) uint8 {
//		//	c.d = c.c
//		//	return 1
//		//},
//		//0x52: func(c *CPU) uint8 {
//		//	// c.d = c.d
//		//	return 1
//		//},
//		//0x53: func(c *CPU) uint8 {
//		//	c.d = c.e
//		//	return 1
//		//},
//		//0x54: func(c *CPU) uint8 {
//		//	c.d = c.h
//		//	return 1
//		//},
//		//0x55: func(c *CPU) uint8 {
//		//	c.d = c.l
//		//	return 1
//		//},
//		//0x57: func(c *CPU) uint8 {
//		//	c.d = c.a
//		//	return 1
//		//},
//		//
//		//// LD E, X
//		//0x58: func(c *CPU) uint8 {
//		//	c.e = c.b
//		//	return 1
//		//},
//		//0x59: func(c *CPU) uint8 {
//		//	c.e = c.c
//		//	return 1
//		//},
//		//0x5a: func(c *CPU) uint8 {
//		//	c.e = c.d
//		//	return 1
//		//},
//		//0x5b: func(c *CPU) uint8 {
//		//	// c.e = c.e
//		//	return 1
//		//},
//		//0x5c: func(c *CPU) uint8 {
//		//	c.e = c.h
//		//	return 1
//		//},
//		//0x5d: func(c *CPU) uint8 {
//		//	c.e = c.l
//		//	return 1
//		//},
//		//0x5f: func(c *CPU) uint8 {
//		//	c.e = c.a
//		//	return 1
//		//},
//		//
//		//// LD H, X
//		//0x60: func(c *CPU) uint8 {
//		//	c.h = c.b
//		//	return 1
//		//},
//		//0x61: func(c *CPU) uint8 {
//		//	c.h = c.c
//		//	return 1
//		//},
//		//0x62: func(c *CPU) uint8 {
//		//	c.h = c.d
//		//	return 1
//		//},
//		//0x63: func(c *CPU) uint8 {
//		//	c.h = c.e
//		//	return 1
//		//},
//		//0x64: func(c *CPU) uint8 {
//		//	// c.h = c.h
//		//	return 1
//		//},
//		//0x65: func(c *CPU) uint8 {
//		//	c.h = c.l
//		//	return 1
//		//},
//		//0x67: func(c *CPU) uint8 {
//		//	c.h = c.a
//		//	return 1
//		//},
//		//
//		//// LD L, X
//		//0x68: func(c *CPU) uint8 {
//		//	c.l = c.b
//		//	return 1
//		//},
//		//0x69: func(c *CPU) uint8 {
//		//	c.l = c.c
//		//	return 1
//		//},
//		//0x6a: func(c *CPU) uint8 {
//		//	c.l = c.d
//		//	return 1
//		//},
//		//0x6b: func(c *CPU) uint8 {
//		//	c.l = c.e
//		//	return 1
//		//},
//		//0x6c: func(c *CPU) uint8 {
//		//	c.l = c.h
//		//	return 1
//		//},
//		//0x6d: func(c *CPU) uint8 {
//		//	// c.l = c.l
//		//	return 1
//		//},
//		//0x6f: func(c *CPU) uint8 {
//		//	c.l = c.a
//		//	return 1
//		//},
//		//
//		//// LD A, X
//		//0x78: func(c *CPU) uint8 {
//		//	c.a = c.b
//		//	return 1
//		//},
//		//0x79: func(c *CPU) uint8 {
//		//	c.a = c.c
//		//	return 1
//		//},
//		//0x7a: func(c *CPU) uint8 {
//		//	c.a = c.d
//		//	return 1
//		//},
//		//0x7b: func(c *CPU) uint8 {
//		//	c.a = c.e
//		//	return 1
//		//},
//		//0x7c: func(c *CPU) uint8 {
//		//	c.a = c.h
//		//	return 1
//		//},
//		//0x7d: func(c *CPU) uint8 {
//		//	c.a = c.l
//		//	return 1
//		//},
//		//0x7f: func(c *CPU) uint8 {
//		//	// c.a = c.a
//		//	return 1
//		//},
//		//
//		//// LD X, (HL)
//		//0x46: func(c *CPU) uint8 {
//		//	c.b = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//0x4e: func(c *CPU) uint8 {
//		//	c.c = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//0x56: func(c *CPU) uint8 {
//		//	c.d = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//0x5e: func(c *CPU) uint8 {
//		//	c.e = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//0x66: func(c *CPU) uint8 {
//		//	c.h = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//0x6e: func(c *CPU) uint8 {
//		//	c.l = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//0x76: func(c *CPU) uint8 {
//		//	c.a = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//
//		//// LD (HL), X
//		//0x70: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.b)
//		//	return 2
//		//},
//		//0x71: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.c)
//		//	return 2
//		//},
//		//0x72: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.d)
//		//	return 2
//		//},
//		//0x73: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.e)
//		//	return 2
//		//},
//		//0x74: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.h)
//		//	return 2
//		//},
//		//0x75: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.l)
//		//	return 2
//		//},
//		//0x77: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.a)
//		//	return 2
//		//},
//		//
//		//// LD X, d8
//		//0x0e: func(c *CPU) uint8 {
//		//	c.c = c.PC()
//		//	return 2
//		//},
//		//0x16: func(c *CPU) uint8 {
//		//	c.d = c.PC()
//		//	return 2
//		//},
//		//0x1e: func(c *CPU) uint8 {
//		//	c.e = c.PC()
//		//	return 2
//		//},
//		//0x26: func(c *CPU) uint8 {
//		//	c.h = c.PC()
//		//	return 2
//		//},
//		//0x2e: func(c *CPU) uint8 {
//		//	c.l = c.PC()
//		//	return 2
//		//},
//		//0x3e: func(c *CPU) uint8 {
//		//	c.a = c.PC()
//		//	return 2
//		//},
//		//
//		//// LD (HL), d8
//		//0x36: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.h, c.l), c.PC())
//		//	return 3
//		//},
//		//
//		//// LD (DE), A
//		//0x12: func(c *CPU) uint8 {
//		//	c.ram.Write(bb2i(c.d, c.e), c.a)
//		//	return 2
//		//},
//		//
//		//// LD (HL+), A
//		//0x22: func(c *CPU) uint8 {
//		//	hl := bb2i(c.h, c.l)
//		//	c.ram.Write(hl, c.a)
//		//	hl++
//		//	c.h, c.l = i2bb(hl)
//		//	return 2
//		//},
//		//
//		//// LD (HL-), A
//		//0x32: func(c *CPU) uint8 {
//		//	hl := bb2i(c.h, c.l)
//		//	c.ram.Write(hl, c.a)
//		//	hl--
//		//	c.h, c.l = i2bb(hl)
//		//	return 2
//		//},
//		//
//		//// LD A, (DE)
//		//0x1a: func(c *CPU) uint8 {
//		//	c.a = c.ram.Read(bb2i(c.d, c.e))
//		//	return 2
//		//},
//		//
//		//// LD A, (HL)
//		//0x7e: func(c *CPU) uint8 {
//		//	c.a = c.ram.Read(bb2i(c.h, c.l))
//		//	return 2
//		//},
//		//
//		//// LD A, (HL+)
//		//0x2a: func(c *CPU) uint8 {
//		//	hl := bb2i(c.h, c.l)
//		//	c.a = c.ram.Read(hl)
//		//	hl++
//		//	c.h, c.l = i2bb(hl)
//		//	return 2
//		//},
//		//
//		//// LD A, (HL-)
//		//0x3a: func(c *CPU) uint8 {
//		//	hl := bb2i(c.h, c.l)
//		//	c.a = c.ram.Read(hl)
//		//	hl--
//		//	c.h, c.l = i2bb(hl)
//		//	return 2
//		//},
//	}
//}

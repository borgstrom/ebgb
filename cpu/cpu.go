package cpu

import (
	"encoding/binary"

	"github.com/borgstrom/ebgb/mmu"
)

// Register represents a 16-bit register that allows for reading and writing its high and low bits
type Register uint16

// GetLow returns the low byte for the Register
func (r Register) GetLow() uint8 {
	return uint8(r)
}

// GetHigh returns the high byte for the Register
func (r Register) GetHigh() uint8 {
	return uint8(r >> 8)
}

// SetLow sets the low byte on the Register
func (r *Register) SetLow(v uint8) {
	*r = Register(binary.LittleEndian.Uint16([]byte{v, uint8(*r >> 8)}))
}

// SetHigh sets the high byte on the Register
func (r *Register) SetHigh(v uint8) {
	*r = Register(binary.LittleEndian.Uint16([]byte{uint8(*r), v}))
}

// CPU implements the 8-bit Sharp LR35902
// See: https://gbdev.io/pandocs
//
// There is no explicit reset function, to reset the CPU just create a new instance.
type CPU struct {
	af Register
	bc Register
	de Register
	hl Register
	sp Register
	pc Register

	halt   bool
	cycles uint64
	ram    mmu.ReadWriter
}

// initFlags set
func (c *CPU) initFlags(flag flag) {
	c.af.SetLow(uint8(flag))
}

func (c *CPU) flipFlag(flag flag) {
	c.af.SetLow(c.af.GetLow() ^ uint8(flag))
}

func (c *CPU) enableFlag(flag flag) {
	c.af.SetLow(c.af.GetLow() | uint8(flag))
}

func (c *CPU) disableFlag(flag flag) {
	c.af.SetLow(c.af.GetLow() & (^uint8(flag)))
}

func (c *CPU) isFlagSet(flag flag) bool {
	return c.af.GetLow()&uint8(flag) != 0
}

// setIncFlags flags takes a value that is the result of an increment/addition and sets the appropriate flags
func (c *CPU) setIncFlags(v uint8) {
}

func (c *CPU) setDecFlags(v uint8) {
	c.setIncFlags(v)
	c.enableFlag(flagSubtraction)
}

func New(ram mmu.ReadWriter) *CPU {
	return &CPU{
		af: 0x01b0,
		bc: 0x0013,
		de: 0x00d8,
		hl: 0x014d,

		sp: 0xfffe,
		pc: 0x0100,

		cycles: 0,

		ram: ram,
	}
}

type flag uint8
type interrupt uint8

const (
	flagNone        flag = 0x00
	flagCarry       flag = 0x10
	flagHalfCarry   flag = 0x20
	flagSubtraction flag = 0x40
	flagZero        flag = 0x80

	interruptNone    interrupt = 0x00
	interruptVBlank  interrupt = 0x01
	interruptLCDSTAT interrupt = 0x02
	interruptTimer   interrupt = 0x04
	interruptSerial  interrupt = 0x08
	interruptJoypad  interrupt = 0x10
)

// Next runs a single iteration of the CPU and returns the number of cycles taken
func (c *CPU) Next() uint8 {
	opCode := c.PC()
	// Break the op code into two parts to look up from our map
	instruction := instructionsByOpcode[opCode&0xf0][opCode&0x0f]
	return c.exec(instruction)
}

func (c *CPU) exec(instruction instructionFunc) uint8 {
	cycles := instruction(c)
	c.cycles = c.cycles + uint64(cycles)
	return cycles
}

func (c *CPU) PC() uint8 {
	v := c.ram.Read(uint16(c.pc))
	c.pc++
	return v
}

func (c *CPU) StackPush(v Register) {
	c.sp--
	c.ram.Write(uint16(c.sp), v.GetHigh())
	c.sp--
	c.ram.Write(uint16(c.sp), v.GetLow())
}

func (c *CPU) StackPop(v Register) {
	v.SetLow(c.ram.Read(uint16(c.sp)))
	c.sp--
	v.SetHigh(c.ram.Read(uint16(c.sp)))
	c.sp--
}

// bb2i converts two separate uint8 into an unsigned 16-bit integer
func bb2i(a, b uint8) uint16 {
	return uint16(a)<<8 | uint16(b)
}

// i2bb converts an unsigned 16-bit integer into two bytes
func i2bb(v uint16) (uint8, uint8) {
	return uint8(v >> 8), uint8(v & 0xff)
}

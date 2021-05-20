package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/borgstrom/ebgb/mmu"
)

func TestInstructions(t *testing.T) {
	var tests = []struct {
		name   string
		memory mmu.RAM
		test   func(t *testing.T, c *CPU)
	}{
		{
			name: "nop",
			test: func(t *testing.T, c *CPU) {
				c.exec(nop)
				require.EqualValues(t, 1, c.cycles)
			},
		},
		{
			name:   "ldBcD16",
			memory: mmu.RAM{0x00, 0x7b},
			test: func(t *testing.T, c *CPU) {
				c.exec(ldBcD16)
				require.EqualValues(t, 123, c.bc)
				require.EqualValues(t, 3, c.cycles)
			},
		},
		{
			name:   "ldBcA",
			memory: mmu.RAM{0x00},
			test: func(t *testing.T, c *CPU) {
				c.bc = 0x0000
				c.af.SetHigh(0xff)
				require.EqualValues(t, 0x00, c.ram.Read(0x0000))
				c.exec(ldBcA)
				require.EqualValues(t, 0xff, c.ram.Read(0x0000))
				require.EqualValues(t, 2, c.cycles)
			},
		},
		{
			name: "incBc",
			test: func(t *testing.T, c *CPU) {
				c.bc.SetLow(0x7b)
				c.bc.SetHigh(0x00)
				c.exec(incBc)
				require.EqualValues(t, 124, c.bc)
				require.EqualValues(t, 2, c.cycles)
			},
		},
		{
			name: "incB",
			test: func(t *testing.T, c *CPU) {
				c.bc.SetHigh(254)

				c.exec(incB)
				require.EqualValues(t, 255, c.bc.GetHigh())
				require.False(t, c.isFlagSet(flagSubtraction))
				require.False(t, c.isFlagSet(flagZero))
				require.False(t, c.isFlagSet(flagHalfCarry))
				require.True(t, c.isFlagSet(flagCarry))

				c.exec(incB)
				require.EqualValues(t, 0, c.bc.GetHigh())
				require.False(t, c.isFlagSet(flagSubtraction))
				require.True(t, c.isFlagSet(flagZero))
				require.True(t, c.isFlagSet(flagHalfCarry))
				require.True(t, c.isFlagSet(flagCarry))

				require.EqualValues(t, 2, c.cycles)
			},
		},
		{
			name: "decB",
			test: func(t *testing.T, c *CPU) {
				c.bc.SetHigh(1)

				c.exec(decB)
				require.EqualValues(t, 0, c.bc.GetHigh())
				require.True(t, c.isFlagSet(flagSubtraction))
				require.True(t, c.isFlagSet(flagZero))
				require.False(t, c.isFlagSet(flagHalfCarry))
				require.True(t, c.isFlagSet(flagCarry))

				c.exec(decB)
				require.EqualValues(t, 255, c.bc.GetHigh())
				require.True(t, c.isFlagSet(flagSubtraction))
				require.False(t, c.isFlagSet(flagZero))
				require.True(t, c.isFlagSet(flagHalfCarry))
				require.True(t, c.isFlagSet(flagCarry))

				require.EqualValues(t, 2, c.cycles)
			},
		},
		{
			name:   "ldBD8",
			memory: mmu.RAM{0x12},
			test: func(t *testing.T, c *CPU) {
				c.exec(ldBD8)
				require.EqualValues(t, 0x12, c.bc.GetHigh())
				require.EqualValues(t, 1, c.cycles)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the CPU with the provided memory
			c := New(test.memory)

			// Set our program counter to 0x0
			c.pc = 0x0000

			// Run the test
			test.test(t, c)
		})
	}
}

package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	var tests = []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Read",
			test: func(t *testing.T) {
				// 0x8000
				r := Register(32768)
				require.EqualValues(t, 0x80, r.GetHigh())
				require.EqualValues(t, 0x00, r.GetLow())
			},
		},
		{
			name: "Write",
			test: func(t *testing.T) {
				r := Register(32768)
				require.EqualValues(t, 0x8000, r)
				r.SetLow(0x80)
				require.EqualValues(t, 0x8080, r)
				r.SetHigh(0x79)
				require.EqualValues(t, 0x7980, r)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.test)
	}
}

//func TestInstructions(t *testing.T) {
//	for opcode, i := range instructions {
//		if i.test == nil {
//			log.Printf("Bad programmer! Write a test for %#x %s", opcode, i.mnemonic)
//			continue
//		}
//		if i.testMemory == nil {
//			log.Printf("Bad programmer! You didn't specify test memory for %#x %s, and now your test will crash", opcode, i.mnemonic)
//		}
//
//		t.Run(fmt.Sprintf("%#x %s", opcode, i.mnemonic), func(t *testing.T) {
//			// Create the CPU with the provided memory
//			c := New(i.testMemory)
//
//			// Run the test function
//			i.test(t, c)
//		})
//	}
//}

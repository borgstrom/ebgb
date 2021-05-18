package cpu

import (
	"fmt"
	"log"
	"testing"
)

func TestInstructions(t *testing.T) {
	for opcode, i := range instructions {
		if i.test == nil {
			log.Printf("Bad programmer! Write a test for %#x %s", opcode, i.mnemonic)
			continue
		}
		if i.testMemory == nil {
			log.Printf("Bad programmer! You didn't specify test memory for %#x %s, and now your test will crash", opcode, i.mnemonic)
		}

		t.Run(fmt.Sprintf("%#x %s", opcode, i.mnemonic), func(t *testing.T) {
			// Create the CPU with the provided memory
			c := New(i.testMemory)

			// Run the test function
			i.test(t, c)
		})
	}
}

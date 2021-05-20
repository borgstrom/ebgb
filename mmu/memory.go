package mmu

// ReadWriter is the interface used by the MMU to Read and Write different Registers
type ReadWriter interface {
	Read(a uint16) uint8
	Write(a uint16, v uint8)
}

// RAM is a slice backed container that can be used as a ReadWriter
type RAM []uint8

func (r RAM) Read(a uint16) uint8 {
	return r[a]
}

func (r RAM) Write(a uint16, v uint8) {
	r[a] = v
}

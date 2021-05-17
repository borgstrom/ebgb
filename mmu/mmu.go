package mmu

import (
	"io"
)

type MMU struct {
	rom  io.ReadSeeker
	eRAM [8192]byte
	wRAM [32768]byte
	zRAM [127]byte

	biosLoaded bool
}

func New(rom io.ReadSeeker) *MMU {
	return &MMU{
		rom: rom,
	}
}

func (m *MMU) Reset() {
	for i := 0; i < len(m.eRAM); i++ {
		m.eRAM[i] = 0x00
	}
	for i := 0; i < len(m.wRAM); i++ {
		m.wRAM[i] = 0x00
	}
	for i := 0; i < len(m.zRAM); i++ {
		m.zRAM[i] = 0x00
	}
}

func (m *MMU) Read(a uint16) byte {
	b := make([]byte, 1)

	switch a & 0xF000 {
	case 0x0000, 0x1000, 0x2000, 0x3000:
		// ROM banks 0 through 3
		m.rom.Seek(int64(a), io.SeekStart)
		m.rom.Read(b)
		return b[0]
	}
}

func (m *MMU) Write(a uint16, v byte) {

}

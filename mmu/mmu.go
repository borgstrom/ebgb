package mmu

type MMU struct {
	rom  []uint8
	eRAM [8192]uint8
	wRAM [32768]uint8
	zRAM [127]uint8

	biosEnabled bool
}

func New(rom []uint8) *MMU {
	return &MMU{
		rom:         rom,
		biosEnabled: true,
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
	m.biosEnabled = true
}

func (m *MMU) Read(a uint16) uint8 {
	switch a & 0xf000 {
	case 0x0000:
		// ROM bank 0, except it returns the bios for the first 256 bytes during boot-up
		if m.biosEnabled && a < 0x0100 {
			return bios[a]
		}
		fallthrough

	case 0x1000, 0x2000, 0x3000:
		// ROM banks 1 through 3
		return m.rom[a]
	}

	panic("Invalid read")
}

func (m *MMU) Write(a uint16, v uint8) {
	switch a & 0xf000 {
	case 0xf000:
		// 0xff00 ... 0xffff
		switch a & 0xff {
		case 0x50:
			m.biosEnabled = false
		}
	}

	panic("Invalid write")
}

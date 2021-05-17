package emulator

import (
	"context"
	"encoding/binary"
	"io"
	"log"
	"time"

	"github.com/borgstrom/ebgb/cpu"
	"github.com/borgstrom/ebgb/gpu"
	"github.com/borgstrom/ebgb/memory"
	"github.com/borgstrom/ebgb/mmu"
)

type Emulator struct {
	rom    io.ReadSeeker
	header CartridgeHeader
	mmu    *mmu.MMU
	cpu    cpu.CPU
	gpu    gpu.GPU
	ram    memory.RAM

	fps           int
	currentSecond int
}

func New(rom io.ReadSeeker) Emulator {
	return Emulator{
		rom: rom,
		mmu: mmu.New(rom),
	}
}

type CartridgeHeader struct {
	EntryPoint      [4]byte
	Logo            [48]byte
	Title           [16]byte
	CGB             byte
	NewLicenseeCode [2]byte
	SGB             byte
	Type            byte
	ROMSize         byte
	RAMSize         byte
	DestinationCode byte
	OldLicenseeCode byte
	MaskROMVersion  byte
	HeaderChecksum  byte
	GlobalChecksum  [2]byte
}

func (e Emulator) Reset() {
	e.ram = memory.RAM{}
	e.cpu = cpu.New(e.ram)
	e.gpu = gpu.New(1)

	e.header = CartridgeHeader{}
	e.rom.Seek(0x0100, io.SeekStart)
	err := binary.Read(e.rom, binary.LittleEndian, e.header)
	if err != nil {
		log.Fatalf("Failed to read rom: %s", err)
	}

	//if hdr.Type == 0x01 {
	//	log.Printf("Yep")
	//}
}
func (e Emulator) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			e.cycle()
		}
	}
}

const (
	lines         = 154
	clocks        = 114
	ticksPerCycle = lines * clocks
)

func (e Emulator) cycle() {
	for i := 0; i < ticksPerCycle; i++ {
		e.cpu.Tick()
		e.gpu.Tick()
	}

	now := time.Now().Second()
	if e.currentSecond == now {
		e.fps++
	} else {
		e.fps = 0
		e.currentSecond = now
	}
}

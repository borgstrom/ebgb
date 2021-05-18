package emulator

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/borgstrom/ebgb/cpu"
	"github.com/borgstrom/ebgb/gpu"
	"github.com/borgstrom/ebgb/mmu"
)

type Emulator struct {
	cartridge *Cartridge

	mmu *mmu.MMU
	cpu *cpu.CPU
	gpu *gpu.GPU

	fps           int
	currentSecond int
}

func New(file io.ReadSeeker) *Emulator {
	cartridge, err := Load(file)
	if err != nil {
		log.Fatalf("Failed to load rom: %s", err)
	}

	e := &Emulator{
		cartridge: cartridge,
	}
	e.Reset()
	return e
}

func (e *Emulator) Reset() {
	e.mmu = mmu.New(e.cartridge.ROM)
	e.cpu = cpu.New(e.mmu)
	e.gpu = gpu.New(e.mmu)
}

// Run starts the *Emulator, must be run in the main thread to satisfy SDL
func (e *Emulator) Run(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"ebgb",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(width),
		int32(height),
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, g0)

	//rect := sdl.Rect{0, 0, 200, 200}
	//surface.FillRect(&rect, g0)
	window.UpdateSurface()

	//lastTick := sdl.GetTicks()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// no-op
		}

		// handle input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return

			case *sdl.KeyboardEvent:
				switch e.Type {
				case sdl.KEYDOWN:
					// Key press
				case sdl.KEYUP:
					// Key release
				}
			}
		}

		e.frame()
	}
}

const (
	// The CPU runs as 4.194304 MHz with a vsync of 59.73 Hz, which is ~70221.06144316 cycles per frame
	// To avoid working with a float we make 59.73 a whole number and shift the cpu speed out two places
	// resulting an int with a value of 70221
	cyclesPerFrame = 419430400 / 5973
)

func (e *Emulator) frame() {
	var cycles uint32
	for cycles < cyclesPerFrame {
		cycles += uint32(e.cpu.Next())
	}

	now := time.Now().Second()
	if e.currentSecond == now {
		e.fps++
	} else {
		e.fps = 0
		e.currentSecond = now
	}
}

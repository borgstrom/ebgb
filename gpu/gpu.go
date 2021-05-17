package gpu

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Tiles
// Palettes
// Layers -> Background, Window, Objects

const (
	// 160x144 pixel display
	width  = 160
	height = 144

	// Gameboy colors, as ints
	// https://1997kherreegoldie.files.wordpress.com/2019/01/gameboy-screen-color-palette-by-ooloopa.studio.jpg
	g0 = (0xca << 16) | (0xdc << 8) | 0x9f
	g1 = (0x0f << 16) | (0x38 << 8) | 0x0f
	g2 = (0x30 << 16) | (0x62 << 8) | 0x30
	g3 = (0x8b << 16) | (0xac << 8) | 0x0f
	g4 = (0x9b << 16) | (0xbc << 8) | 0x0f
)

type GPU struct {
	displayMultiplier int
}

func New(displayMultiplier int) GPU {
	g := GPU{
		displayMultiplier: displayMultiplier,
	}
	g.init()
	return g
}

func (g GPU) Tick() {

}

func (g GPU) init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"ebgb",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(width*g.displayMultiplier),
		int32(height*g.displayMultiplier),
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

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

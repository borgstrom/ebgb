package emulator

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

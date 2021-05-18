package gpu

// Tiles
// Palettes
// Layers -> Background, Window, Objects

type GPU struct {
	ram Memory
}

type Memory interface {
	Read(a uint16) uint8
	Write(a uint16, v uint8)
}

func New(ram Memory) *GPU {
	g := &GPU{
		ram: ram,
	}
	return g
}

func (g *GPU) Next() {

}

package emulator

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

type Cartridge struct {
	Size   int
	ROM    []uint8
	Header CartridgeHeader
}

func Load(file io.ReadSeeker) (*Cartridge, error) {
	rom, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	header := CartridgeHeader{}
	file.Seek(0x0100, io.SeekStart)
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	return &Cartridge{
		ROM:    rom,
		Size:   len(rom),
		Header: header,
	}, nil
}

type CartridgeHeader struct {
	EntryPoint      [4]uint8
	Logo            [48]uint8
	Title           [16]uint8
	CGB             uint8
	NewLicenseeCode [2]uint8
	SGB             uint8
	Type            uint8
	ROMSize         uint8
	RAMSize         uint8
	DestinationCode uint8
	OldLicenseeCode uint8
	MaskROMVersion  uint8
	HeaderChecksum  uint8
	GlobalChecksum  [2]uint8
}

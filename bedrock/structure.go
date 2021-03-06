package bedrock

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"io"
	"os"
)

// Structure holds the data of an .mcstructure file. Structure implements the world.Structure interface. It
// may be built in a world by using (world.World).BuildStructure.
type Structure struct {
	*structure
}

// WriteStructure writes a Structure to the io.Writer passed. If successful, the error returned is nil.
func WriteStructure(w io.Writer, s Structure) error {
	s.Structure.Palettes[s.paletteName] = *s.palette

	if err := nbt.NewEncoderWithEncoding(w, nbt.LittleEndian).Encode(s.structure); err != nil {
		return fmt.Errorf("encode structure: %w", err)
	}
	return nil
}

// WriteStructureToFile writes a Structure to the file passed. If successful, the error returned is nil. WriteFile
// creates a file if it doesn't yet exist and truncates it if one does exist.
func WriteStructureToFile(file string, s Structure) error {
	f, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	return WriteStructure(f, s)
}

// NewStructure creates a new Structure and initialises it with air blocks. The Structure returned may be written to
// using Structure.Set and Structure.SetAdditionalLiquid and the palette may be changed by using UsePalette.
func NewStructure(dimensions [3]int) Structure {
	front := make([]int32, dimensions[0]*dimensions[1]*dimensions[2])
	liquids := make([]int32, dimensions[0]*dimensions[1]*dimensions[2])
	for i := range liquids {
		liquids[i] = -1
	}

	s := Structure{structure: &structure{
		FormatVersion: version,
		Size:          []int32{int32(dimensions[0]), int32(dimensions[1]), int32(dimensions[2])},
		Origin:        []int32{0, 0, 0},
		Structure: structureData{
			BlockIndices: [][]int32{front, liquids},
			Palettes:     map[string]palette{},
		},
	}}
	s.UsePalette("default")
	s.palette.BlockPalette = append(s.palette.BlockPalette, block{
		Name:    "minecraft:air",
		States:  map[string]interface{}{},
		Version: CurrentBlockVersion,
	})
	return s
}

// UsePalette changes the palette name to use for the Structure. When reading a Structure, this will change
// the palette used to read blocks from. When writing a Structure, the palette will be written with this name,
// so that subsequent readers of the Structure must first call UsePalette with this name to get the right
// palette.
func (s Structure) UsePalette(name string) {
	if current := s.palette; current != nil {
		s.Structure.Palettes[s.paletteName] = *s.palette
	}

	p, _ := s.Structure.Palettes[name]
	if p.BlockPositionData == nil {
		p.BlockPositionData = map[string]blockPositionData{}
	}
	s.palette = &p
	s.paletteName = name
}

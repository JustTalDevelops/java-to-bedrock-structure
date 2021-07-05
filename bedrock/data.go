package bedrock

import (
	"fmt"
	"github.com/justtaldevelops/java-to-bedrock-structure/states"
)

const (
	// CurrentBlockVersion is the current version of blocks (states) of the game. This version is composed
	// of 4 bytes indicating a version, interpreted as a big endian int. The current version represents
	// 1.16.0.14 {1, 16, 0, 14}.
	CurrentBlockVersion int32 = 17825806
)

// structure is the outer wrapper of the structure. It holds the version of the structure, its dimensions,
// origin in the world it was created in, and the actual structure data.
type structure struct {
	FormatVersion int32         `nbt:"format_version"`
	Size          []int32       `nbt:"size"`
	Origin        []int32       `nbt:"structure_world_origin"`
	Structure     structureData `nbt:"structure"`

	palette     *palette
	paletteName string
}

const version = 1

// check checks if the structure is valid. It returns an error if anything in the structure was found to be
// incorrect.
func (s *structure) check() error {
	if s.FormatVersion != version {
		return fmt.Errorf("unsupported format version %v: expected version %v", s.FormatVersion, version)
	}
	if l := len(s.Size); l != 3 {
		return fmt.Errorf("structure size must have 3 values, but got %v (%v)", l, s.Size)
	}
	if l := len(s.Origin); l != 3 {
		return fmt.Errorf("structure origin must have 3 values, but got %v (%v)", l, s.Origin)
	}
	if s.Structure.Palettes == nil {
		s.Structure.Palettes = map[string]palette{}
	}
	if len(s.Structure.BlockIndices) == 0 {
		return fmt.Errorf("structure has no blocks in it")
	}
	if len(s.Structure.Palettes) == 0 {
		return fmt.Errorf("structure has no palettes in it")
	}
	for _, indices := range s.Structure.BlockIndices {
		size := int(s.Size[0] * s.Size[1] * s.Size[2])
		if len(indices) != size {
			return fmt.Errorf("structure is %vx%vx%v and should have %v blocks, but got only %v", s.Size[0], s.Size[1], s.Size[2], size, len(indices))
		}
	}
	paletteLen := -1
	for _, p := range s.Structure.Palettes {
		if paletteLen == -1 {
			paletteLen = len(p.BlockPalette)
			continue
		}
		if len(p.BlockPalette) != paletteLen {
			return fmt.Errorf("all palettes must have the same length, but got one with length %v and one with length %v", paletteLen, len(p.BlockPalette))
		}
	}
	return nil
}

// Dimensions returns the dimensions of the structure as set in the Origin field.
func (s *structure) Dimensions() [3]int {
	return [3]int{int(s.Size[0]), int(s.Size[1]), int(s.Size[2])}
}

// Set sets the block at a specific position within the structure to the world.Block passed. Set will panic
// if the x, y or z exceed the bounds of the structure. The world.Liquid passed may be nil to avoid waterlogging the
// block.
func (s *structure) Set(x, y, z int, state states.BedrockState) {
	l, h := int(s.Size[2]), int(s.Size[1])
	offset := (x * l * h) + (y * l) + z

	name, properties := state.Name, state.Properties
	ptr := s.lookup(name, properties)
	if ptr == -1 {
		ptr = int32(len(s.palette.BlockPalette))
		s.palette.BlockPalette = append(s.palette.BlockPalette, block{
			Name:    name,
			States:  properties,
			Version: CurrentBlockVersion,
		})
	}

	s.Structure.BlockIndices[0][offset] = ptr
	s.Structure.BlockIndices[1][offset] = -1
}

// lookup looks up the world.Block passed in the palette of the structure. If not found, the value returned is
// -1.
func (s *structure) lookup(name string, properties map[string]interface{}) int32 {
	for index, block := range s.palette.BlockPalette {
		if block.Name == name {
			allEqual := true
			for k, v := range block.States {
				if bVal, _ := properties[k]; bVal != v {
					allEqual = false
					break
				}
			}
			if allEqual {
				return int32(index)
			}
		}
	}
	return -1
}

// structureData holds the actual data of the structure. This includes both blocks and entities.
type structureData struct {
	// BlockIndices holds the actual block data. This is a two-dimensional slice, where the first indicates
	// the layer that these blocks are in. The int32s held are pointers to blocks in the palette used. Note
	// that an index may be -1 to indicate that neither air nor any other block should be placed on this
	// position.
	BlockIndices [][]int32                `nbt:"block_indices"`
	Entities     []map[string]interface{} `nbt:"entities"`
	Palettes     map[string]palette       `nbt:"palette"`
}

// palette represents the palette of a single structure.
type palette struct {
	BlockPalette      []block                      `nbt:"block_palette"`
	BlockPositionData map[string]blockPositionData `nbt:"block_position_data"`
}

// block represents a single block entry, holding a block name and its states. These entries also hold a
// version.
type block struct {
	Name    string                 `nbt:"name"`
	States  map[string]interface{} `nbt:"states"`
	Version int32                  `nbt:"version"`
}

// blockPositionData holds additional data associated with specific block positions in the structure. At the
// moment, these appear to be limited to block entity data.
type blockPositionData struct {
	BlockEntityData map[string]interface{} `nbt:"block_entity_data"`
}

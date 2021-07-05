package java

import (
	"bytes"
	"compress/gzip"
	"github.com/justtaldevelops/java-to-bedrock-structure/bedrock"
	"github.com/justtaldevelops/java-to-bedrock-structure/states"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"io/ioutil"
)

// BlockEntry is a block entry in a structure.
type BlockEntry struct {
	// Pos is the position of the block entry.
	Pos []int32 `nbt:"pos"`
	// State is the assigned state.
	State int32 `nbt:"state"`

	// NBT is unused.
	NBT interface{} `nbt:"nbt,omitempty"`
}

// Structure is a Java Edition structure.
type Structure struct {
	// Blocks are positions with an assigned state.
	Blocks []BlockEntry `nbt:"blocks"`
	// Palette is the block palette. It maps indexes to states.
	Palette []states.JavaState `nbt:"palette"`
	// Size are the bounds of the structure.
	Size []int32 `nbt:"size"`

	// Entities is unused.
	Entities interface{} `nbt:"entities"`
	// DataVersion is unused.
	DataVersion interface{}
}

// convert converts a Java edition structure to a Bedrock edition structure.
func (structure Structure) Convert() bedrock.Structure {
	size := structure.Size

	bedrockStructure := bedrock.NewStructure([3]int{int(size[0]), int(size[1]), int(size[2])})
	for _, b := range structure.Blocks {
		converted := convert(structure.Palette[b.State])
		if converted.Name == "minecraft:structure_block" {
			continue
		}

		bedrockStructure.Set(int(b.Pos[0]), int(b.Pos[1]), int(b.Pos[2]), converted)
	}

	return bedrockStructure
}

// DecodeStructure decodes a structure from bytes into a Java Edition structure.
func DecodeStructure(data []byte) (javaStructure Structure) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	parsedData, err := ioutil.ReadAll(gr)
	if err != nil {
		panic(err)
	}
	err = gr.Close()
	if err != nil {
		panic(err)
	}

	err = nbt.UnmarshalEncoding(parsedData, &javaStructure, nbt.BigEndian)
	if err != nil {
		panic(err)
	}

	return
}

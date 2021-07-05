package java

import (
	_ "embed"
	"encoding/json"
	"github.com/justtaldevelops/java-to-bedrock-structure/states"
	"github.com/tidwall/gjson"
)

var (
	//go:embed mappings.json
	mappingsData []byte
	// conversionTable converts between a Java encoded state and a Bedrock state.
	conversionTable = make(map[string]states.BedrockState)
)

func init() {
	parsedData := gjson.ParseBytes(mappingsData)
	parsedData.ForEach(func(key, value gjson.Result) bool {
		var bState states.BedrockState

		err := json.Unmarshal([]byte(value.String()), &bState)
		if err != nil {
			panic(err)
		}

		conversionTable[key.String()] = bState
		return true
	})
}

// convert converts a states.JavaState to a states.BedrockState.
func convert(state states.JavaState) states.BedrockState {
	return conversionTable[state.Encode()]
}
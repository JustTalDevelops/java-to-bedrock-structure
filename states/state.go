package states

import (
	"fmt"
	"strings"
)

// BedrockState is the Bedrock edition state.
type BedrockState struct {
	// Name is the Bedrock edition name.
	Name string `json:"bedrock_identifier"`
	// Properties is the Bedrock edition properties.
	Properties map[string]interface{} `json:"bedrock_states"`
}

// JavaState is the Java edition state.
type JavaState struct {
	// Name is the name of the Java state.
	Name string
	// Properties is the Java state properties.
	Properties map[string]interface{}
}

// Encode encodes the Java state into a string.
func (state JavaState) Encode() string {
	sb := &strings.Builder{}

	sb.WriteString(state.Name)

	if len(state.Properties) > 0 {
		sb.WriteString("[")
	}

	var index int
	for k, v := range state.Properties {
		sb.WriteString(fmt.Sprintf("%v=%v", k, v))

		if index < len(state.Properties) - 1 {
			sb.WriteString(",")
		}
		index++
	}

	if len(state.Properties) > 0 {
		sb.WriteString("]")
	}

	return sb.String()
}

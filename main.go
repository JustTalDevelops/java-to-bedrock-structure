package main

import (
	"github.com/justtaldevelops/java-to-bedrock-structure/bedrock"
	"github.com/justtaldevelops/java-to-bedrock-structure/java"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("java_test.nbt")
	if err != nil {
		panic(err)
	}

	err = bedrock.WriteStructureToFile("bedrock_test.mcstructure", java.DecodeStructure(data).Convert())
	if err != nil {
		panic(err)
	}
}
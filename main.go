package main

import (
	"encoding/json"
	"log"

	ccjson "github.com/tuananhlai/cc-json-parser/json"
)

func main() {
	output, err := ccjson.Parse(`{"foo": "bar", "baz": 3.14}`)
	if err != nil {
		log.Fatalf("error parsing json: %v", err)
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("error marshaling json: %v", err)
	}

	log.Println(string(jsonOutput))
}

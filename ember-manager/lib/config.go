package lib

import (
	"encoding/json"
	"log"
	"os"
)

type Processor struct {
	Name     string
	Dir      string
	FileType string
}

type config struct {
	Processors []Processor
	Vendors    []string
}

var Config config

func init() {
	file, err := os.Open("ember-manager-config.json")
	if err != nil {
		log.Fatal("ember-manager-config.json file needed")
	}
	defer file.Close()

	Config = config{}
	if err = json.NewDecoder(file).Decode(&Config); err != nil {
		log.Fatal("error reading ember-manager-config.json ", err)
	}
}

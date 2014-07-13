package lib

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Processors []*Processor
	Vendors    []string
}

func (c *config) GetChanFor(name string) (ch chan *File) {
	for _, p2 := range c.Processors {
		if name == p2.Name {
			ch = p2.InC
			break
		}
	}
	return
}

var Config config

func init() {
	MainChan = make(chan *File)

	file, err := os.Open("ember-manager-config.json")
	if err != nil {
		log.Fatal("ember-manager-config.json file needed")
	}
	defer file.Close()

	Config = config{}
	if err = json.NewDecoder(file).Decode(&Config); err != nil {
		log.Fatal("error reading ember-manager-config.json ", err)
	}

	for _, p := range Config.Processors {
		p.Start()
	}
}

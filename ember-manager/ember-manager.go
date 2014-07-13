package main

import (
	"log"

	"github.com/monocle/ember-tools/ember-manager/lib"
)

func main() {
	log.SetFlags(log.Ltime)

	for _, p := range lib.Config.Processors {
		lib.NewAppWatcher(p)
	}

	lib.StartServer("3000", lib.MainChan)
}

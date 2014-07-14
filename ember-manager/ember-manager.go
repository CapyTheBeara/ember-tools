package main

import (
	"flag"
	"log"

	"github.com/monocle/ember-tools/ember-manager/lib"
)

var (
	proxy = flag.String("proxy", "", "Proxy for requests")
	port  = flag.String("port", "4200", "Server port")
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	flag.Parse()
	lib.Config.Parse()

	for _, p := range lib.Config.Processors {
		lib.NewAppWatcher(p)
	}

	lib.StartServer(*port, *proxy, lib.MainChan)
}

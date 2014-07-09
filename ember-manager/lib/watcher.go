package lib

import (
	"log"
	"strings"
	"time"

	"code.google.com/p/go.exp/fsnotify"
)

var lastTime = time.Now()

func NewAppWatcher(paths []string, fileExt string, evtC chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case evt := <-watcher.Event:
				now := time.Now()

				// TODO - watch newly created folders
				if strings.HasSuffix(evt.Name, fileExt) && now.Sub(lastTime).Seconds() > 0.01 {
					lastTime = now
					evtC <- evt.Name
				}

			case err := <-watcher.Error:
				log.Println(Color("[watch error]", "red"), err)
			}
		}
	}()

	for _, path := range paths {
		err = watcher.Watch(path)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(Color("[watching]", "cyan"), path)
	}
}

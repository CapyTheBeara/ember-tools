package lib

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"code.google.com/p/go.exp/fsnotify"
)

var lastTime = time.Now()

func NewAppWatcher(paths []string, fileExt string, evtC chan File) {
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

					file, err := ioutil.ReadFile(evt.Name)
					if err != nil {
						log.Fatal("watcher error", err)
					}

					evtC <- File{evt.Name, file, evt}
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

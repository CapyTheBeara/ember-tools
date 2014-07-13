package lib

import (
	"log"
	"time"

	"code.google.com/p/go.exp/fsnotify"
)

var lastTime = time.Now()

func NewAppWatcher(p *Processor) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range GetAllDirs(p.Dirs) {
		err = watcher.Watch(path)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(Color("[watching]", "cyan"), path)
	}

	go func() {
		for {
			select {
			case evt := <-watcher.Event:
				now := time.Now()
				// TODO - watch newly created folders
				f := File{Path: evt.Name, Event: evt}
				if p.FileHit(&f) && now.Sub(lastTime).Seconds() > 0.01 {
					lastTime = now
					f.SetContent()
					p.InC <- &f
				}

			case err := <-watcher.Error:
				log.Println(Color("[watch error]", "red"), err)
			}
		}
	}()
}

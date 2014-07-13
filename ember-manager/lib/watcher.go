package lib

import (
	"log"
	"time"

	"code.google.com/p/go.exp/fsnotify"
)

var lastTimes map[string]time.Time

var zeroTime = time.Time{}

func init() {
	lastTimes = make(map[string]time.Time)
}

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
				f := File{Path: evt.Name, Event: evt}

				// TODO - watch newly created folders
				if p.FileHit(&f) && (now.Sub(lastTimes[p.Name]).Seconds() > 0.01 || lastTimes[p.Name] == zeroTime) {
					lastTimes[p.Name] = now
					f.SetContent()
					p.InC <- &f
				}

			case err := <-watcher.Error:
				log.Println(Color("[watch error]", "red"), err)
			}
		}
	}()
}

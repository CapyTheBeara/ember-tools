package lib

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var MainChan chan *File

type Processor struct {
	Name    string
	Dirs    []string
	FileExt string

	PipeTo  []string
	Type    string
	Command *Command
	InC     chan *File
	OutCs   []chan *File
}

func (p *Processor) FileHit(f *File) bool {
	return strings.HasSuffix(f.Path, p.FileExt)
}

func (p *Processor) IsReload() bool {
	return !p.Command.NoPipe
}

func (p *Processor) Start() {
	if len(p.Dirs) == 0 {
		p.Dirs = []string{"app"}
	}

	if p.FileExt == "" {
		p.FileExt = "js"
	}

	if p.Type == "" {
		p.Type = "file"
	}

	p.makeChans()
	p.makeCommand()
	p.listen()
	p.sendAllFiles()
}

func (p *Processor) makeChans() {
	p.InC = make(chan *File)

	for _, to := range p.PipeTo {
		p.OutCs = append(p.OutCs, Config.GetChanFor(to))
	}

	p.OutCs = append(p.OutCs, MainChan)
}

func (p *Processor) makeCommand() {
	if p.Type == "reload" {
		p.Command = &Command{}
	} else if p.Command == nil {
		dir := "processors"

		f, err := ioutil.ReadFile(filepath.Join(dir, p.Name+".js"))
		if err != nil {
			log.Fatalf("Error opening processor file %s/%s", dir, p.Name)
		}

		p.Command = &Command{
			Name:   "node",
			Args:   []string{"-e"},
			Source: f,
		}
	}
}

func (p *Processor) listen() {
	go func() {
		for {
			select {
			case file := <-p.InC:
				f := p.Command.Run(file)
				if p.IsReload() {
					p.OutCs[0] <- f
				}
			}
		}
	}()
}

func (p *Processor) sendAllFiles() {
	var res []*File
	files := make(chan *File)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for file := range files {
				res = append(res, p.Command.Run(file))
			}
			wg.Done()
		}()
	}

	for _, dir := range p.Dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal("[error] reading file", err)
			}

			f := File{Path: path}
			f.SetContent()

			if !info.IsDir() && p.FileHit(&f) {
				files <- &f
			}

			return err
		})
	}

	close(files)
	wg.Wait()

	go func() {
		for _, f := range res {
			if !f.IsEmpty() {
				p.OutCs[0] <- f
			}
		}
	}()
}

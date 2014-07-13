package lib

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		p.Command = &Command{OutCs: p.OutCs}
	} else if p.Command == nil {
		dir := "js"

		f, err := ioutil.ReadFile(filepath.Join(dir, p.Name+".js"))
		if err != nil {
			log.Fatalf("Error opening processor file %s/%s", dir, p.Name)
		}

		p.Command = &Command{
			Name:   "node",
			Args:   []string{"-e"},
			Source: f,
			OutCs:  p.OutCs,
		}
	}
}

func (p *Processor) listen() {
	go func() {
		for {
			select {
			case file := <-p.InC:
				go p.Command.Run(file)
			}
		}
	}()
}

func (p *Processor) sendAllFiles() {
	for _, dir := range p.Dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal("[error] reading file", err)
			}

			f := File{Path: path}
			f.SetContent()

			if !info.IsDir() && p.FileHit(&f) {
				go p.Command.Run(&f)
			}
			return err
		})
	}
}

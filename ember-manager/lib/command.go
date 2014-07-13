package lib

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	Name   string
	Args   []string
	Source []byte
	OutCs  []chan *File
}

func (c *Command) Run(f *File) {
	start := time.Now()

	if c.Name == "" {
		c.OutCs[0] <- &File{}
		return
	}

	args := append(c.Args, string(c.Source), f.Path, string(f.Content))
	cmd := exec.Command(c.Name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}

	if len(content) < 1 {
		msg, err := ioutil.ReadAll(stderr)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(Color("[processor error]", "red"), f.Path, string(msg))

	} else {
		log.Printf("%s %.2fs for %s\n", Color("[processing]", "magenta"), time.Since(start).Seconds(), f.Path)

		// TODO - do mo' better
		split := strings.Split(string(content), "OUTPUT_PATH=")
		if len(split) > 1 {
			f.Path = split[1]
		}

		f.Content = []byte(split[0])
		ch := c.OutCs[0]
		ch <- f
	}
}

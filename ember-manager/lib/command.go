package lib

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	Name     string
	Args     []string
	Source   []byte
	PathOnly bool
	NoPipe   bool
}

func (c *Command) Run(f *File) (res *File) {
	start := time.Now()

	if c.Name == "" {
		return &File{Path: "reload"}
	}

	args := append(c.Args, string(c.Source), f.Path, string(f.Content))
	if c.PathOnly {
		args = append(c.Args, f.Path)
	}

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

		if len(msg) > 0 {
			log.Println(Color("[processor error]", "red"), c.Name, f.Path, string(msg))
		}
	} else {
		log.Printf("%s %.2fs for %s\n", Color("[processing]", "magenta"), time.Since(start).Seconds(), f.Path)

		if c.NoPipe {
			log.Println(Color(string(content), "magenta"))
		}

		// TODO - do mo' better
		split := strings.Split(string(content), "OUTPUT_PATH=")
		if len(split) > 1 {
			f.Path = split[1]
		}

		f.Content = []byte(split[0])
		return f
	}
	return &File{}
}

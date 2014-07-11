package lib

import (
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

func CommandFn(commandName string, args ...string) func(string, []byte, chan File) {
	return func(path string, input []byte, res chan File) {
		log.Println(Color("[processing]", "magenta"), path)

		start := time.Now()
		cmd := exec.Command(commandName, append(args, path, string(input))...)

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

			log.Println(Color("[processor error]", "red"), path, string(msg))

		} else {
			log.Printf("%s %.2fs for %s\n", Color("[processing]", "magenta"), time.Since(start).Seconds(), path)
			res <- File{path, content, nil}
		}
	}

}

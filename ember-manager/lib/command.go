package lib

import (
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

type CommandRes struct {
	FileName string
	Output   string
}

func CommandFn(commandName string, args ...string) func(string, chan CommandRes) {
	return func(path string, res chan CommandRes) {
		log.Println(Color("[processing]", "magenta"), path)

		start := time.Now()
		cmd := exec.Command(commandName, append(args, path)...)

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

		out, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}

		if len(out) < 1 {
			msg, err := ioutil.ReadAll(stderr)
			if err != nil {
				log.Fatal(err)
			}

			if string(msg) == "FILE_NOT_FOUND" {
				log.Println(path, " was removed")
			} else {
				log.Println("[js error]", path, string(msg))
			}
		} else {
			log.Printf("%s %.2fs for %s\n", Color("[processing]", "magenta"), time.Since(start).Seconds(), path)
			res <- CommandRes{path, string(out)}
		}
	}

}

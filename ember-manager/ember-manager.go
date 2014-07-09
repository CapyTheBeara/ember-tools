package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/monocle/ember-tools/ember-manager/lib"
)

var jsCompilerSource map[string]string
var jsSource map[string]string

func compileAll(dirs []string, fileType string, callback func(string, chan lib.CommandRes), c chan lib.CommandRes) {
	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), "."+fileType) {
				go callback(filepath.Join(dir, info.Name()), c)
			}
			return err
		})
	}
}

func main() {
	// log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	jsCompilerSource = lib.GetJS("./js")
	jsSource = make(map[string]string)

	es6Transpiler := lib.CommandFn("node", "-e", jsCompilerSource["es6-transpile"])

	// TODO - recursively watch folders
	jsDirs := []string{"app/controllers", "app/models"}

	jsC := make(chan lib.CommandRes)
	compileAll(jsDirs, "js", es6Transpiler, jsC)

	jsEventC := make(chan string)
	lib.NewAppWatcher(jsDirs, "js", jsEventC)

	serverC := make(chan map[string]string)
	go lib.StartServer("3000", serverC)

	for {
		select {
		case fileName := <-jsEventC:
			go es6Transpiler(fileName, jsC)
		case res := <-jsC:
			jsSource[res.FileName] = res.Output
			serverC <- jsSource
		}
	}

}

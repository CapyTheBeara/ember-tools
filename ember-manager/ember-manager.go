package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/monocle/ember-tools/ember-manager/lib"
)

var jsCompilerSource map[string]string

// TODO - handle concurrency issues
// http://blog.golang.org/go-maps-in-action
var jsSource map[string]string

func compileAll(dirs []string, fileType string, callback func(string, []byte, chan lib.File), c chan lib.File) {
	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), "."+fileType) {
				file, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}

				go callback(path, file, c)
			}
			return err
		})
	}
}

func main() {
	jsCompilerSource = lib.GetJS("./js")
	jsSource = make(map[string]string)

	es6C := make(chan lib.File)
	processor := lib.Config.Processors[0]

	es6Transpiler := lib.CommandFn("node", "-e", jsCompilerSource[processor.Name])
	// // TODO - recursively watch folders
	// jsDirs := []string{"app", "app/controllers", "app/models", "app/routes"}
	lib.NewAppWatcher([]string{processor.Dir}, "js", es6C)

	emberTemplateCompiler := lib.CommandFn("node", "-e", jsCompilerSource["ember-template-compiler"])
	hbsC := make(chan lib.File)
	lib.NewAppWatcher([]string{"app/templates"}, "hbs", hbsC)

	reloadC := make(chan lib.File)
	lib.NewAppWatcher([]string{"app/styles"}, "css", reloadC)

	jsC := make(chan lib.File)
	compileAll([]string{"app"}, "js", es6Transpiler, jsC)
	compileAll([]string{"app/templates"}, "hbs", emberTemplateCompiler, es6C)

	serverC := make(chan map[string]string)
	go lib.StartServer("3000", serverC, reloadC)

	for {
		select {
		case file := <-hbsC:
			if file.IsEmpty() {
				go func() {
					es6C <- file
				}()
				continue
			}
			go emberTemplateCompiler(file.Path, file.Content, es6C)

		case file := <-es6C:
			path := file.Path

			if strings.HasSuffix(path, ".hbs") {
				path = strings.Replace(path, ".hbs", ".js", 1)
			}

			if file.IsEmpty() {
				go func() {
					jsC <- lib.File{Path: path}
				}()
				continue
			}

			go es6Transpiler(path, file.Content, jsC)

		case file := <-jsC:
			if file.IsEmpty() {
				delete(jsSource, file.Path)
			} else {
				jsSource[file.Path] = string(file.Content)
			}
			serverC <- jsSource
		}
	}

}

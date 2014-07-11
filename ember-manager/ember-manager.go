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
	// log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	jsCompilerSource = lib.GetJS("./js")
	jsSource = make(map[string]string)

	es6Transpiler := lib.CommandFn("node", "-e", jsCompilerSource["es6-transpiler"])
	emberTemplateCompiler := lib.CommandFn("node", "-e", jsCompilerSource["ember-template-compiler"])

	templateC := make(chan lib.File)
	compileAll([]string{"app/templates"}, "hbs", emberTemplateCompiler, templateC)

	jsC := make(chan lib.File)
	compileAll([]string{"app"}, "js", es6Transpiler, jsC)

	// TODO - recursively watch folders
	jsDirs := []string{"app", "app/controllers", "app/models", "app/routes"}
	es6C := make(chan lib.File)
	lib.NewAppWatcher(jsDirs, "js", es6C)

	hbsCC := make(chan lib.File)
	lib.NewAppWatcher([]string{"app/templates"}, "hbs", hbsCC)

	serverC := make(chan map[string]string)
	go lib.StartServer("3000", serverC)

	for {
		select {
		case file := <-hbsCC:
			go emberTemplateCompiler(file.Path, file.Content, es6C)

		case file := <-es6C:
			// TODO - handle file DELETE events here
			path := file.Path

			if strings.HasSuffix(path, ".hbs") {
				path = strings.Replace(path, ".hbs", ".js", 1)
			}

			go es6Transpiler(path, file.Content, jsC)

		case res := <-jsC:
			jsSource[res.Path] = string(res.Content)
			serverC <- jsSource
		}
	}

}

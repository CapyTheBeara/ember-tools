package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetJS(dir string) map[string]string {
	files := make(map[string]string)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".js") {
			file, err := ioutil.ReadFile(filepath.Join(dir, info.Name()))
			if err != nil {
				panic(err)
			}

			files[strings.Split(info.Name(), ".")[0]] = string(file)
		}
		return err
	})

	return files
}

package lib

import (
	"io/ioutil"
	"log"
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

func GetVendorJS() (vendorScripts []string, err error) {
	dir := "vendor/"

	for vendor, path := range Config.Vendors {
		file, err := ioutil.ReadFile(filepath.Join(dir, path))
		if err != nil {
			log.Fatal("error reading vendor file:", vendor, err)
		}

		vendorScripts = append(vendorScripts, string(file))
	}
	return
}

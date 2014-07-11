package lib

import (
	"encoding/json"
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
	file, err := os.Open("ember-manager-config.json")
	if err != nil {
		log.Fatal("vendor_config.json file needed")
	}
	defer file.Close()

	config := make(map[string]interface{})
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal("error reading vendor_config.json", err)
	}

	dir := "vendor/"

	for vendor, path := range config["vendor"].([]interface{}) {
		file, err := ioutil.ReadFile(filepath.Join(dir, path.(string)))
		if err != nil {
			log.Fatal("error reading vendor file:", vendor, err)
		}

		vendorScripts = append(vendorScripts, string(file))
	}

	return
}

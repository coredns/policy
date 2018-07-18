// ETC (Execute Template Catalog) executes all templates recursively in a given directory.
package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	parseFlags()
	execute(conf.template, load(conf.data, conf.selector), makePrefix(conf.template))
}

func execute(name string, data interface{}, prefix []string) {
	absName, err := filepath.Abs(name)
	if err != nil {
		log.Fatal(err)
	}

	relName := getRelName(absName, prefix)
	log.Printf("executing %q", relName)

	fi, err := os.Stat(absName)
	if err != nil {
		log.Fatal(err)
	}

	newName := executeString(absName, data)
	if newName == absName {
		log.Fatalf("instance name is the same as template name %q", newName)
	}

	if fi.IsDir() {
		log.Printf("creating directory %q -> %q", relName, getRelName(newName, prefix))
		executeDir(newName)

		lst, err := ioutil.ReadDir(absName)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range lst {
			execute(path.Join(absName, item.Name()), data, prefix)
		}
	} else {
		log.Printf("creating file %q -> %q", relName, getRelName(newName, prefix))
		executeFile(newName, absName, data)
	}
}

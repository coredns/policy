package main

import (
	"log"
	"os"
	"text/template"
)

func executeDir(name string) {
	if err := os.MkdirAll(name, 0755); err != nil {
		log.Fatal(err)
	}
}

func executeFile(out, in string, data interface{}) {
	t, err := template.ParseFiles(in)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		log.Fatal(err)
	}
}

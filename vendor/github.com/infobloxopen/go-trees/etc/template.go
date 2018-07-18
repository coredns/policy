package main

import (
	"bytes"
	"log"
	"text/template"
)

func executeString(s string, data interface{}) string {
	t := template.New("string")
	t, err := t.Parse(s)
	if err != nil {
		log.Fatal(err)
	}

	b := new(bytes.Buffer)
	if err := t.Execute(b, data); err != nil {
		log.Fatal(err)
	}

	return b.String()
}

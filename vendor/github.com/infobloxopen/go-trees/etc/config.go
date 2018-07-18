package main

import (
	"flag"
	"log"
)

type config struct {
	template string
	data     string
	selector string
}

var conf config

func parseFlags() {
	flag.StringVar(&conf.template, "t", "", "path to template (required)")
	flag.StringVar(&conf.data, "d", "", "path to data")
	flag.StringVar(&conf.selector, "s", "", "path in YAML to get data")

	flag.Parse()

	if len(conf.template) <= 0 {
		log.Fatal("No path to template - nothing to execute")
	}
}

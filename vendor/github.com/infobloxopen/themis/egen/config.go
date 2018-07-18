package main

import "flag"

const defaultInput = "errors.yaml"

type config struct {
	input string
	check bool
}

var conf config

func init() {
	flag.StringVar(&conf.input, "i", defaultInput, "path to input file")
	flag.BoolVar(&conf.check, "c", false,
		"don't generate output just check usage of any defined errors in current directory *.go files")
	flag.Parse()
}

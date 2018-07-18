package main

import (
	"fmt"
	"os"
)

func main() {
	e, err := unmarshal(conf.input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if conf.check {
		e.check()
	} else {
		e.generate()
	}
}

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

const grepBin = "/usr/bin/grep"

func (e *errors) check() {
	if len(e.Errors) <= 0 {
		fmt.Println("Nothing to check")
		return
	}

	count := 0
	for _, item := range e.Errors {
		r := fmt.Sprintf("\\<new%s%s\\>", strings.ToUpper(string(item.ID[0])), item.ID[1:])
		if err := exec.Command(grepBin, "-qr", "--include=*.go", r, ".").Run(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					if status != 0 {
						fmt.Printf("%q hasn't been found\n", item.ID)
						count++
					}
				} else {
					panic(fmt.Errorf("'%s -qr --include=*.go %s .': %s", grepBin, r, err))
				}
			} else {
				panic(fmt.Errorf("'%s -qr --include=*.go %s .': %s", grepBin, r, err))
			}
		}
	}

	if count <= 0 {
		fmt.Printf("All %d entries of new<Error-id> has been found\n", len(e.Errors))
	}
}

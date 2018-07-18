package main

import (
	"log"
	"path"
	"path/filepath"
)

func makePrefix(p string) []string {
	a, err := filepath.Abs(p)
	if err != nil {
		log.Fatal(err)
	}

	d, f := filepath.Split(a)
	if len(f) <= 0 {
		log.Fatalf("can't execute template at %q (%q) as it has no base name", p, a)
	}

	return fullSplit(d)
}

func getRelName(p string, prefix []string) string {
	seq := fullSplit(p)
	if hasPrefix(seq, prefix) {
		return path.Join(seq[len(prefix):]...)
	}

	return filepath.Base(p)
}

const preAllocStrings = 10

func fullSplit(p string) []string {
	d := filepath.Clean(p)

	out := make([]string, preAllocStrings)
	i := preAllocStrings
	for {
		i--
		if i < 0 {
			out = append(make([]string, preAllocStrings), out...)
			i = preAllocStrings
			continue
		}

		parent, f := filepath.Split(d)
		if len(f) <= 0 {
			out[i] = filepath.Clean(parent)
			break
		}

		if len(parent) <= 0 {
			out[i] = f
			break
		}

		out[i] = f

		d = filepath.Clean(parent)
	}

	return out[i:]
}

func hasPrefix(p, prefix []string) bool {
	if len(p) < len(prefix) {
		return false
	}

	for i, item := range prefix {
		if p[i] != item {
			return false
		}
	}

	return true
}

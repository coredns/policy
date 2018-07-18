package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func load(name string, path string) interface{} {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)

	m := make(map[interface{}]interface{})
	if err := d.Decode(m); err != nil {
		log.Fatalf("can't decode data from %q: %s", name, err)
	}

	out, err := sel(path, m)
	if err != nil {
		log.Fatalf("no data for %q in %q: %s", path, name, err)
	}

	return out
}

func sel(path string, m map[interface{}]interface{}) (interface{}, error) {
	if len(path) > 0 {
		parts := strings.SplitN(path, ".", 2)
		v, ok := m[parts[0]]
		if !ok {
			return nil, fmt.Errorf("missing %q", parts[0])
		}

		if len(parts) > 1 {
			m, ok := v.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("can't go deeper than %q", parts[0])
			}

			return sel(parts[1], m)
		}

		return v, nil
	}

	return m, nil
}

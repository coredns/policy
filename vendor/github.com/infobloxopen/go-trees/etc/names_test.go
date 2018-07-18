package main

import (
	"path"
	"path/filepath"
	"testing"
)

func TestNameGetRelName(t *testing.T) {
	e := filepath.FromSlash("d/e")
	s := getRelName(
		filepath.FromSlash("/a/b/c/d/e"),
		[]string{filepath.FromSlash("/"), "a", "b", "c"},
	)

	if s != e {
		t.Errorf("expected %q as relative path but got %q", e, s)
	}

	e = filepath.FromSlash("e")
	s = getRelName(
		filepath.FromSlash("/a/b/c/d/e"),
		[]string{filepath.FromSlash("/"), "a", "B", "c"},
	)

	if s != e {
		t.Errorf("expected %q as relative path but got %q", e, s)
	}
}

func TestNameFullSplit(t *testing.T) {
	assertPathSlice(t, "3-level abs path",
		fullSplit(filepath.FromSlash("/a/b/c")),
		filepath.FromSlash("/"), "a", "b", "c",
	)

	assertPathSlice(t, "11-level abs path",
		fullSplit(filepath.FromSlash("/a/b/c/d/e/f/g/h/i/j/k")),
		filepath.FromSlash("/"), "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
	)

	assertPathSlice(t, "rel path",
		fullSplit(filepath.FromSlash("a/b/c")),
		"a", "b", "c",
	)

	assertPathSlice(t, "dirty path",
		fullSplit(filepath.FromSlash("a//b/./c/")),
		"a", "b", "c",
	)
}

func TestNameHasPrefix(t *testing.T) {
	p := []string{"/", "a", "b", "c", "d"}
	prefix := []string{"/", "a", "b", "c"}

	if !hasPrefix(p, prefix) {
		t.Errorf("expected %q to have prefix %q", path.Join(p...), path.Join(prefix...))
	}

	p = []string{"/", "a", "b"}
	if hasPrefix(p, prefix) {
		t.Errorf("expected %q not to have prefix %q", path.Join(p...), path.Join(prefix...))
	}

	p = []string{"/", "a", "B", "c", "d"}
	if hasPrefix(p, prefix) {
		t.Errorf("expected %q not to have prefix %q", path.Join(p...), path.Join(prefix...))
	}
}

func assertPathSlice(t *testing.T, desc string, v []string, e ...string) {
	if len(v) != len(e) {
		t.Errorf("expected %q for %s but got %q", path.Join(e...), desc, path.Join(v...))
	}

	for i, item := range e {
		if v[i] != item {
			t.Errorf("expected %q for %s but got %q", path.Join(e...), desc, path.Join(v...))
		}
	}
}

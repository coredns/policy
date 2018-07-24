package server

import "testing"

func TestBytePool(t *testing.T) {
	p := makeBytePool(10, false)

	b := p.Get()
	if len(b) != 10 {
		t.Errorf("expected buffer of %d bytes but got %d %#v", 10, len(b), b)
	}

	p.Put(b)
}

func TestDummyBytePool(t *testing.T) {
	p := makeBytePool(10, true)

	b := p.Get()
	if len(b) != 10 {
		t.Errorf("expected buffer of %d bytes but got %d %#v", 10, len(b), b)
	}

	p.Put(b)
}

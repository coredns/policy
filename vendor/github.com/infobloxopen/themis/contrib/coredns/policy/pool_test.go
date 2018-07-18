package policy

import "testing"

func TestAttrPool(t *testing.T) {
	p := makeAttrPool(10, false)

	a := p.Get()
	if len(a) != 10 {
		t.Errorf("expected buffer of %d attributes but got %d %#v", 10, len(a), a)
	}

	p.Put(a)
}

func TestDummyAttrPool(t *testing.T) {
	p := makeAttrPool(10, true)

	a := p.Get()
	if len(a) != 10 {
		t.Errorf("expected buffer of %d attributes but got %d %#v", 10, len(a), a)
	}

	p.Put(a)
}

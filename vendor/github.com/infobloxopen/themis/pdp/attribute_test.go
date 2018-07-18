package pdp

import "testing"

func TestAttribute(t *testing.T) {
	a := MakeAttribute("test", TypeString)
	if a.id != "test" {
		t.Errorf("Expected \"test\" as attribute id but got %q", a.id)
	}

	at := a.GetType()
	if at != TypeString {
		t.Errorf("Expected %q as attribute type but got %q", TypeString, at)
	}
}

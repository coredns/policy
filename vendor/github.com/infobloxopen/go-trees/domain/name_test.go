package domain

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNameMakeNameFromString(t *testing.T) {
	s := "wiki.example.com"
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	if n.h != s {
		t.Errorf("expected %q as human-readable name but got %q", s, n.h)
	}

	e := "\x03COM\x07EXAMPLE\x04WIKI"
	if n.c != e {
		t.Errorf("expected %q as name for comparison but got %q", e, n.c)
	}
}

func TestNameMakeNameFromStringEmpty(t *testing.T) {
	s := ""
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	if n.h != s {
		t.Errorf("expected %q as human-readable name but got %q", s, n.h)
	}

	e := ""
	if n.c != e {
		t.Errorf("expected %q as name for comparison but got %q", e, n.c)
	}
}

func TestNameMakeNameFromStringDot(t *testing.T) {
	s := "."
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	if n.h != s {
		t.Errorf("expected %q as human-readable name but got %q", s, n.h)
	}

	e := ""
	if n.c != e {
		t.Errorf("expected %q as name for comparison but got %q", e, n.c)
	}
}

func TestNameMakeNameFromStringWithNameTooLong(t *testing.T) {
	s := "toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo." +
		"loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong." +
		"doooooooooooooooooooooooooooooooooooooooooooooooooooooooooomain." +
		"naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaame"

	n, err := MakeNameFromString(s)
	if err == nil {
		t.Fatalf("expected error but got name %q", n.c)
	}

	if err != ErrNameTooLong {
		t.Fatalf("expected ErrNameTooLong but got %q (%T)", err, err)
	}
}

func TestNameMakeNameFromStringWithTooManyLabels(t *testing.T) {
	s := "0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9." +
		"0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9." +
		"0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9.0.1.2.3.4.5.6.7.8.9." +
		"0.1.2.3.4.5.6.7"

	n, err := MakeNameFromString(s)
	if err == nil {
		t.Fatalf("expected error but got name %q", n.c)
	}

	if err != ErrTooManyLabels {
		t.Fatalf("expected ErrTooManyLabels but got %q (%T)", err, err)
	}
}

func TestNameMakeNameFromStringWithTooLongLabel(t *testing.T) {
	s := "www.looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong.com"

	n, err := MakeNameFromString(s)
	if err == nil {
		t.Fatalf("expected error but got name %q", n.c)
	}

	if err != ErrLabelTooLong {
		t.Fatalf("expected ErrLabelTooLong but got %q (%T)", err, err)
	}
}

func TestNameMakeNameFromReflection(t *testing.T) {
	n, err := MakeNameFromString("www.example.com")
	if err != nil {
		t.Fatal(err)
	}

	nr := MakeNameFromReflection(reflect.ValueOf(n))
	if nr.h != n.h || nr.c != n.c {
		t.Errorf("expected %q (%q) but got %q (%q)", n, n.c, nr, nr.c)
	}

	nr = MakeNameFromReflection(reflect.ValueOf("www.example.com"))
	if nr.String() != "www.example.com" {
		t.Errorf("expected %q but got %q", "www.example.com", nr)
	}
}

func TestNameMakeNameFromReflectionPtr(t *testing.T) {
	n, err := MakeNameFromString("www.example.com")
	if err != nil {
		t.Fatal(err)
	}

	nr := MakeNameFromReflection(reflect.ValueOf(&n))
	if nr.h != n.h || nr.c != n.c {
		t.Errorf("expected %q (%q) but got %q (%q)", n, n.c, nr, nr.c)
	}
}

func TestNameMakeNameFromReflectionPanicWrongType(t *testing.T) {
	var (
		i int
		n Name
	)

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				name := reflect.TypeOf(i).Name()
				if !strings.Contains(err.Error(), name) {
					t.Errorf("expected %q in error message but got %s", name, err)
				}
			} else {
				t.Errorf("expected panic on error but got %T (%#v)", r, r)
			}
		} else {
			t.Fatalf("expected panic but got name %q", n)
		}
	}()

	n = MakeNameFromReflection(reflect.ValueOf(&i))
}

func TestNameMakeNameFromReflectionPanicWrongString(t *testing.T) {
	s := "empty..label"
	var n Name

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if !strings.Contains(err.Error(), s) {
					t.Errorf("expected %q in error message but got %s", s, err)
				}
			} else {
				t.Errorf("expected panic on error but got %T (%#v)", r, r)
			}
		} else {
			t.Fatalf("expected panic but got name %q", n)
		}
	}()

	n = MakeNameFromReflection(reflect.ValueOf(s))
}

func TestNameString(t *testing.T) {
	s := "wiki.example.com"
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	h := n.String()
	if h != s {
		t.Errorf("expected %q as human-readable name but got %q", s, n.h)
	}
}

func TestGetLabel(t *testing.T) {
	s := "wiki.example.com"
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	lbls := []string{}
	off := 0
	for {
		lbl, next := n.GetLabel(off)
		if next < 0 {
			t.Fatalf("expected nonnegative offset but got %d after %d (%#v)", next, off, lbls)
		}

		lbls = append(lbls, lbl)
		off = next
		if off == 0 {
			break
		}
	}

	assertLabels(t, lbls, []string{"COM", "EXAMPLE", "WIKI"})
}

func TestGetLabelWithRoot(t *testing.T) {
	s := ""
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	lbls := []string{}
	off := 0
	for {
		lbl, next := n.GetLabel(off)
		if next < 0 {
			t.Fatalf("expected nonnegative offset but got %d after %d (%#v)", next, off, lbls)
		}

		lbls = append(lbls, lbl)
		off = next
		if off == 0 {
			break
		}
	}

	assertLabels(t, lbls, []string{""})
}

func TestGetLabelWithInvalidOffset(t *testing.T) {
	s := "wiki.example.com"
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	lbl, off := n.GetLabel(-1)
	if off >= 0 {
		t.Errorf("expected negative offset but got label %q", lbl)
	}

	lbl, off = n.GetLabel(len(n.c))
	if off >= 0 {
		t.Errorf("expected negative offset but got label %q", lbl)
	}

	lbl, off = n.GetLabel(2)
	if off >= 0 {
		t.Errorf("expected negative offset but got label %q", lbl)
	}
}

func TestGetLabels(t *testing.T) {
	s := "wiki.example.com"
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	lbls := []string{}
	if err := n.GetLabels(func(lbl string) error {
		lbls = append(lbls, lbl)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	assertLabels(t, lbls, []string{"COM", "EXAMPLE", "WIKI"})
}

func TestGetLabelsWithError(t *testing.T) {
	s := "wiki.example.com"
	n, err := MakeNameFromString(s)
	if err != nil {
		t.Fatal(err)
	}

	stop := fmt.Errorf("stop iteration")

	lbls := []string{}
	err = n.GetLabels(func(lbl string) error {
		lbls = append(lbls, lbl)
		if len(lbls) >= 2 {
			return stop
		}

		return nil
	})
	if err == nil {
		t.Fatalf("expected error but got %d labels:\n%#v", len(lbls), lbls)
	}

	if err != stop {
		t.Errorf("expected \"stop\" error but got %q (%T)", err, err)
	}

	assertLabels(t, lbls, []string{"COM", "EXAMPLE"})
}

func assertLabels(t *testing.T, v, e []string) {
	if len(v) != len(e) {
		t.Errorf("expected %d labels\n\t%#v\nbut got %d\n\t%#v", len(e), e, len(v), v)
		return
	}

	for i, b := range e {
		if v[i] != b {
			t.Errorf("expected labels\n\t%#v\nbut got\n\t%#v", e, v)
			return
		}
	}
}

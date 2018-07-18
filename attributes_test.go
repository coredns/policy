package policy

import (
	"strings"
	"testing"

	"github.com/infobloxopen/themis/pdp"
)

func TestSerializeOrPanic(t *testing.T) {
	s := serializeOrPanic(pdp.MakeStringAssignment("s", "test"))
	if s != "test" {
		t.Errorf("expected %q but got %q", "test", s)
	}

	assertPanicWithErrorContains(t, "serializeOrPanic(expression)", func() {
		serializeOrPanic(pdp.MakeExpressionAssignment("s", pdp.MakeStringDesignator("s")))
	}, "pdp.AttributeDesignator")

	assertPanicWithErrorContains(t, "serializeOrPanic(undefined)", func() {
		serializeOrPanic(pdp.MakeExpressionAssignment("s", pdp.UndefinedValue))
	}, "Undefined")
}

func TestCustAttr(t *testing.T) {
	if !custAttr(custAttrEdns).isEdns() {
		t.Errorf("expected %d is EDNS", custAttrEdns)
	}

	if !custAttr(custAttrTransfer).isTransfer() {
		t.Errorf("expected %d is transfer", custAttrTransfer)
	}

	if !custAttr(custAttrDnstap).isDnstap() {
		t.Errorf("expected %d is DNStap", custAttrDnstap)
	}
}

func assertPanicWithErrorContains(t *testing.T, desc string, f func(), e string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				t.Errorf("excpected error containing %q on panic for %q but got %T (%#v)", e, desc, r, r)
			} else if !strings.Contains(err.Error(), e) {
				t.Errorf("excpected error containing %q on panic for %q but got %q", e, desc, r)
			}
		} else {
			t.Errorf("expected error containing %q on panic for %q", e, desc)
		}
	}()

	f()
}

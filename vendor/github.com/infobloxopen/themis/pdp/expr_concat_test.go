package pdp

import (
	"fmt"
	"testing"
)

func TestFunctionConcat(t *testing.T) {
	s := MakeStringValue("teststring")

	ss := MakeSetOfStringsValue(newStrTree(
		"set-first",
		"set-second",
		"set-third",
	))

	ls := MakeListOfStringsValue([]string{
		"list-first",
		"list-second",
		"list-third",
	})

	ft8, err := NewFlagsType("8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Fatal(err)
	}
	f8 := MakeFlagsValue8(0x5, ft8)

	ft16, err := NewFlagsType("16flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
	)
	if err != nil {
		t.Fatal(err)
	}
	f16 := MakeFlagsValue16(0x500, ft16)

	ft32, err := NewFlagsType("32flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
	)
	if err != nil {
		t.Fatal(err)
	}
	f32 := MakeFlagsValue32(0x50000, ft32)

	ft64, err := NewFlagsType("64flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77",
	)
	if err != nil {
		t.Fatal(err)
	}
	f64 := MakeFlagsValue64(0x500000000, ft64)

	f := makeFunctionConcat([]Expression{s, ss, ls, f8, f16, f32, f64})
	if f, ok := f.(functionConcat); ok {
		e := "concat"
		desc := f.describe()
		if desc != e {
			t.Errorf("expected %q description but got %q", e, desc)
		}
	} else {
		t.Errorf("expected functionConcat but got %T (%#v)", f, f)
	}

	rt := f.GetResultType()
	if rt != TypeListOfStrings {
		t.Errorf("expected %q type but got %q", TypeListOfStrings, rt)
	} else {
		v, err := f.Calculate(nil)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				e := "\"teststring\",\"set-first\",\"set-second\",\"set-third\"," +
					"\"list-first\",\"list-second\",\"list-third\"," +
					"\"f00\",\"f02\",\"f10\",\"f12\",\"f20\",\"f22\",\"f40\",\"f42\""
				if e != s {
					t.Errorf("Expected list of strings %q but got %q", e, s)
				}
			}
		}
	}

	m := findValidator("concat", s, ss, ls, f8, f16, f32, f64)
	if m == nil {
		t.Errorf("expected makeFunctionConcat but got %#v", m)
	} else {
		f := m([]Expression{s, ss, ls, f8, f16, f32, f64})
		if _, ok := f.(functionConcat); !ok {
			t.Errorf("expected functionConcat but got %T (%#v)", f, f)
		}
	}

	m = findValidator("concat", s, ss, ls, f8, f16, f32, f64, MakeIntegerValue(5))
	if m != nil {
		t.Errorf("expected nothing but got %#v", m)
	}

	m = findValidator("concat")
	if m != nil {
		t.Errorf("expected nothing but got %#v", m)
	}
}

func TestFunctionConcatWithMissingValues(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	ls1 := MakeListOfStringsValue([]string{"1", "2", "3"})
	ls2 := MakeListOfStringsValue([]string{"A", "B", "C"})
	mlose := failExpr{t: TypeListOfStrings, err: newMissingValueError()}
	mse := failExpr{t: TypeString, err: newMissingValueError()}
	fsose := failExpr{t: TypeSetOfStrings, err: fmt.Errorf("test error")}

	v, err := makeFunctionConcat([]Expression{ls1, ls2}).Calculate(ctx)
	if err != nil {
		t.Error(err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"1\",\"2\",\"3\",\"A\",\"B\",\"C\""
			if e != s {
				t.Errorf("Expected list of strings %q but got %q", e, s)
			}
		}
	}

	v, err = makeFunctionConcat([]Expression{mlose, ls2}).Calculate(ctx)
	if err != nil {
		t.Error(err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"A\",\"B\",\"C\""
			if e != s {
				t.Errorf("Expected list of strings %q but got %q", e, s)
			}
		}
	}

	v, err = makeFunctionConcat([]Expression{ls1, mse}).Calculate(ctx)
	if err != nil {
		t.Error(err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"1\",\"2\",\"3\""
			if e != s {
				t.Errorf("Expected list of strings %q but got %q", e, s)
			}
		}
	}

	v, err = makeFunctionConcat([]Expression{mse, mlose}).Calculate(ctx)
	if err == nil {
		t.Errorf("Expected *MissingValueError but got value %s", v.describe())
	} else if _, ok := err.(*MissingValueError); !ok {
		t.Errorf("Expected *MissingValueError but got %T: %s", err, err)
	}

	v, err = makeFunctionConcat([]Expression{ls1, fsose}).Calculate(ctx)
	if err == nil {
		t.Errorf("Expected *externalError but got value %s", v.describe())
	} else if _, ok := err.(*externalError); !ok {
		t.Errorf("Expected *externalError but got %T: %s", err, err)
	}
}

type failExpr struct {
	t   Type
	err error
}

func (f failExpr) GetResultType() Type {
	return f.t
}

func (f failExpr) Calculate(ctx *Context) (AttributeValue, error) {
	return UndefinedValue, f.err
}

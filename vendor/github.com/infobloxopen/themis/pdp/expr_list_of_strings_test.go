package pdp

import (
	"testing"

	"github.com/infobloxopen/go-trees/strtree"
)

func TestFunctionListOfStrings(t *testing.T) {
	l := MakeListOfStringsValue([]string{
		"first",
		"second",
		"third",
	})

	f := makeFunctionListOfStrings(l)

	if f, ok := f.(functionListOfStrings); ok {
		expDesc := "list of strings"
		desc := f.describe()
		if desc != expDesc {
			t.Errorf("Expected %q description but got %q", expDesc, desc)
		}
	} else {
		t.Errorf("Expected functionListOfStrings but got %T (%#v)", f, f)
	}

	rt := f.GetResultType()
	if rt != TypeListOfStrings {
		t.Errorf("Expected %q type but got %q", TypeListOfStrings, rt)
	} else {
		v, err := f.Calculate(nil)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				e := "\"first\",\"second\",\"third\""
				if e != s {
					t.Errorf("Expected list of strings %q but got %q", e, s)
				}
			}
		}
	}

	m := findValidator("list of strings", l)
	if m == nil {
		t.Errorf("Expected makeFunctionListOfStringsAlt but got %v", m)
	} else {
		f := m([]Expression{l})
		if _, ok := f.(functionListOfStrings); !ok {
			t.Errorf("Expected functionListOfStrings but got %T (%v)", f, f)
		}
	}

	strTree := strtree.NewTree()
	strTree.InplaceInsert("fourth", 4)
	strTree.InplaceInsert("third", 3)
	strTree.InplaceInsert("second", 2)
	strTree.InplaceInsert("first", 1)

	s := MakeSetOfStringsValue(strTree)

	v, err := makeFunctionListOfStrings(s).Calculate(nil)
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"first\",\"second\",\"third\",\"fourth\""
			if e != s {
				t.Errorf("Expected list of strings %q but got %q", e, s)
			}
		}
	}

	m = findValidator("list of strings", s)
	if m == nil {
		t.Errorf("Expected makeFunctionListOfStringsAlt but got %v", m)
	} else {
		f := m([]Expression{s})
		if _, ok := f.(functionListOfStrings); !ok {
			t.Errorf("Expected functionListOfStrings but got %T (%v)", f, f)
		}
	}

	ft8, err := NewFlagsType("8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v, err := makeFunctionListOfStrings(MakeFlagsValue8(0x55, ft8)).Calculate(nil)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				e := "\"f00\",\"f02\",\"f04\",\"f06\""
				if e != s {
					t.Errorf("Expected list of strings %q but got %q", e, s)
				}
			}
		}

		m := findValidator("list of strings", MakeFlagsValue8(0x55, ft8))
		if m == nil {
			t.Errorf("Expected makeFunctionListOfStringsAlt but got %v", m)
		} else {
			f := m([]Expression{MakeFlagsValue8(0x55, ft8)})
			if _, ok := f.(functionListOfStrings); !ok {
				t.Errorf("Expected functionListOfStrings but got %T (%v)", f, f)
			}
		}
	}

	ft16, err := NewFlagsType("16flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v, err := makeFunctionListOfStrings(MakeFlagsValue16(0x5555, ft16)).Calculate(nil)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				e := "\"f00\",\"f02\",\"f04\",\"f06\",\"f10\",\"f12\",\"f14\",\"f16\""
				if e != s {
					t.Errorf("Expected list of strings %q but got %q", e, s)
				}
			}
		}
	}

	ft32, err := NewFlagsType("32flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v, err := makeFunctionListOfStrings(MakeFlagsValue32(0x55555555, ft32)).Calculate(nil)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				e := "\"f00\",\"f02\",\"f04\",\"f06\",\"f10\",\"f12\",\"f14\",\"f16\"," +
					"\"f20\",\"f22\",\"f24\",\"f26\",\"f30\",\"f32\",\"f34\",\"f36\""
				if e != s {
					t.Errorf("Expected list of strings %q but got %q", e, s)
				}
			}
		}
	}

	ft64, err := NewFlagsType("64flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73", "f74", "f75", "f76",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v, err := makeFunctionListOfStrings(MakeFlagsValue64(0x5555555555555555, ft64)).Calculate(nil)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				e := "\"f00\",\"f02\",\"f04\",\"f06\",\"f10\",\"f12\",\"f14\",\"f16\"," +
					"\"f20\",\"f22\",\"f24\",\"f26\",\"f30\",\"f32\",\"f34\",\"f36\"," +
					"\"f40\",\"f42\",\"f44\",\"f46\",\"f50\",\"f52\",\"f54\",\"f56\"," +
					"\"f60\",\"f62\",\"f64\",\"f66\",\"f70\",\"f72\",\"f74\",\"f76\""
				if e != s {
					t.Errorf("Expected list of strings %q but got %q", e, s)
				}
			}
		}
	}

	m = findValidator("list of strings", s, s)
	if m != nil {
		t.Errorf("Expected nil but got %v", m)
	}

	m = findValidator("list of strings", MakeStringValue("test"))
	if m != nil {
		t.Errorf("Expected nil but got %v", m)
	}
}

func findValidator(n string, args ...Expression) functionMaker {
	if v, ok := FunctionArgumentValidators[n]; ok {
		for _, v := range v {
			if m := v(args); m != nil {
				return m
			}
		}
	}

	return nil
}

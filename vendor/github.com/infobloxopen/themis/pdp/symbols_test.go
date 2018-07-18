package pdp

import "testing"

func TestSymbolsType(t *testing.T) {
	s := MakeSymbols()

	bt := s.GetType(TypeUndefined.GetKey())
	if bt == nil {
		t.Errorf("Expected %q but got nothing", TypeUndefined)
	} else if bt != TypeUndefined {
		t.Errorf("Expected %q but got %q", TypeUndefined, bt)
	}

	bt = s.GetType(TypeDomain.GetKey())
	if bt == nil {
		t.Errorf("Expected %q but got nothing", TypeDomain)
	} else if bt != TypeDomain {
		t.Errorf("Expected %q but got %q", TypeDomain, bt)
	}

	st := s.GetType("flags")
	if st != nil {
		t.Errorf("Expected nothing but got %q", st)
	}

	ft, err := NewFlagsType("Flags", "first", "second", "third")
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		err = s.PutType(ft)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			st := s.GetType("flags")
			if st == nil {
				t.Errorf("Expected %q but got nothin", ft)
			} else if st != ft {
				t.Errorf("Expected %q but got %q", ft, st)
			}
		}

		err = s.PutType(ft)
		if err == nil {
			t.Error("Expected *duplicateCustomTypeError but got nothing")
		} else if _, ok := err.(*duplicateCustomTypeError); !ok {
			t.Errorf("Expected *duplicateCustomTypeError but got %T (%s)", err, err)
		}
	}

	err = s.PutType(nil)
	if err == nil {
		t.Error("Expected *nilTypeError but got nothing")
	} else if _, ok := err.(*nilTypeError); !ok {
		t.Errorf("Expected *nilTypeError but got %T (%s)", err, err)
	}

	err = s.PutType(TypeString)
	if err == nil {
		t.Error("Expected *builtinCustomTypeError but got nothing")
	} else if _, ok := err.(*builtinCustomTypeError); !ok {
		t.Errorf("Expected *builtinCustomTypeError but got %T (%s)", err, err)
	}

	ros := s.makeROCopy()
	ft1, err := NewFlagsType("OtherFlags", "first", "second", "third", "fourth")
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		err = ros.PutType(ft1)
		if err == nil {
			t.Error("Expected *ReadOnlySymbolsChangeError but got nothing")
		} else if _, ok := err.(*ReadOnlySymbolsChangeError); !ok {
			t.Errorf("Expected *builtinCustomTypeError but got %T (%s)", err, err)
		}
	}
}

func TestSymbolsAttribute(t *testing.T) {
	s := MakeSymbols()

	a := MakeAttribute("a", TypeString)
	err := s.PutAttribute(a)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		A := MakeAttribute("A", TypeDomain)
		err := s.PutAttribute(A)
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		if sa, ok := s.GetAttribute("a"); ok {
			if sa.id != a.id || sa.t != a.t {
				t.Errorf("Expected %s but got %s", a.describe(), sa.describe())
			}
		} else {
			t.Errorf("Expected attribute %q but got nothing", a.id)
		}

		err = s.PutAttribute(a)
		if err == nil {
			t.Error("Expected *duplicateAttributeError but got nothing")
		} else if _, ok := err.(*duplicateAttributeError); !ok {
			t.Errorf("Expected *duplicateAttributeError but got %T (%s)", err, err)
		}
	}

	if sa, ok := s.GetAttribute("b"); ok {
		t.Errorf("Expected no attribute but got %s", sa.describe())
	}

	n := MakeAttribute("n", nil)
	err = s.PutAttribute(n)
	if err == nil {
		t.Error("Expected *noTypedAttributeError but got nothing")
	} else if _, ok := err.(*noTypedAttributeError); !ok {
		t.Errorf("Expected *noTypedAttributeError but got %T (%s)", err, err)
	}

	u := MakeAttribute("u", TypeUndefined)
	err = s.PutAttribute(u)
	if err == nil {
		t.Error("Expected *undefinedAttributeTypeError but got nothing")
	} else if _, ok := err.(*undefinedAttributeTypeError); !ok {
		t.Errorf("Expected *undefinedAttributeTypeError but got %T (%s)", err, err)
	}

	ft, err := NewFlagsType("Flags", "first", "second", "third")
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		fa := MakeAttribute("fa", ft)
		err = s.PutAttribute(fa)
		if err == nil {
			t.Error("Expected *unknownAttributeTypeError but got nothing")
		} else if _, ok := err.(*unknownAttributeTypeError); !ok {
			t.Errorf("Expected *unknownAttributeTypeError but got %T (%s)", err, err)
		}
	}

	ros := s.makeROCopy()
	err = ros.PutAttribute(MakeAttribute("x", TypeBoolean))
	if err == nil {
		t.Error("Expected *ReadOnlySymbolsChangeError but got nothing")
	} else if _, ok := err.(*ReadOnlySymbolsChangeError); !ok {
		t.Errorf("Expected *builtinCustomTypeError but got %T (%s)", err, err)
	}
}

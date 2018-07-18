package pdp

import "testing"

func TestFunctionTry(t *testing.T) {
	sAttr := MakeStringDesignator("x")
	sVal := MakeStringValue("test")

	f := makeFunctionTry([]Expression{sAttr, sVal})
	if f, ok := f.(functionTry); ok {
		e := "try"
		desc := f.describe()
		if desc != e {
			t.Errorf("expected %q description but got %q", e, desc)
		}
	} else {
		t.Errorf("expected functionTry but got %T (%#v)", f, f)
	}

	rt := f.GetResultType()
	if rt != sAttr.GetResultType() {
		t.Errorf("expected type of arguments (%s) but got %q", sAttr.GetResultType(), rt)
	} else {
		ctx, err := NewContext(nil, 0, nil)
		if err != nil {
			t.Error(err)
		} else {
			v, err := f.Calculate(ctx)
			if err != nil {
				t.Error(err)
			} else {
				s, err := v.Serialize()
				if err != nil {
					t.Error(err)
				} else {
					e := "test"
					if s != e {
						t.Errorf("expected string %q but got %q", e, s)
					}
				}
			}
		}

		ctx, err = NewContext(nil, 1, func(i int) (string, AttributeValue, error) {
			return "x", MakeStringValue("value"), nil
		})
		if err != nil {
			t.Error(err)
		} else {
			v, err := f.Calculate(ctx)
			if err != nil {
				t.Error(err)
			} else {
				s, err := v.Serialize()
				if err != nil {
					t.Error(err)
				} else {
					e := "value"
					if s != e {
						t.Errorf("expected string %q but got %q", e, s)
					}
				}
			}
		}
	}

	f = makeFunctionTry([]Expression{sAttr})
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Error(err)
	} else {
		v, err := f.Calculate(ctx)
		if err == nil {
			t.Errorf("expected *missingAttributeError but got value %s", v.describe())
		} else if _, ok := err.(*missingAttributeError); !ok {
			t.Errorf("expected *missingAttributeError but got %T: %s", err, err)
		}
	}

	m := findValidator("try", sAttr, sVal)
	if m == nil {
		t.Errorf("expected makeFunctionTry but got %#v", m)
	} else {
		f := m([]Expression{sAttr, sVal})
		if _, ok := f.(functionTry); !ok {
			t.Errorf("expected functionTry but got %T (%#v)", f, f)
		}
	}

	m = findValidator("try")
	if m != nil {
		t.Errorf("expected nothing but got %#v", m)
	}

	m = findValidator("try", MakeStringValue("test"), MakeIntegerValue(5))
	if m != nil {
		t.Errorf("expected nothing but got %#v", m)
	}
}

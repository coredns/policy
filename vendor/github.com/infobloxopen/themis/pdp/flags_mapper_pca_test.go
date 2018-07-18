package pdp

import "testing"

func TestFlagsMapperPCAOrdering(t *testing.T) {
	ft, err := NewFlagsType("flags", "third", "first", "second")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	c := &Context{
		a: map[string]interface{}{
			"f": MakeFlagsValue8(7, ft),
		},
	}

	p := NewPolicySet("test", false, Target{},
		[]Evaluable{
			makeSimplePermitPolicyWithObligations(
				"first",
				makeSingleStringObligation("order", "first"),
			),
			makeSimplePermitPolicyWithObligations(
				"second",
				makeSingleStringObligation("order", "second"),
			),
			makeSimplePermitPolicyWithObligations(
				"third",
				makeSingleStringObligation("order", "third"),
			),
		},
		makeMapperPCA, MapperPCAParams{
			Argument:  MakeDesignator("f", ft),
			Order:     MapperPCAInternalOrder,
			Algorithm: firstApplicableEffectPCA{},
		},
		nil,
	)

	r := p.Calculate(c)
	if len(r.Obligations) != 1 {
		t.Errorf("Expected the only obligation got %#v", r.Obligations)
	} else {
		ID, ot, s, err := r.Obligations[0].Serialize(c)
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		eID := "order"
		eot := TypeString
		es := "first"
		if ID != eID || ot != eot.GetKey() || s != es {
			t.Errorf("Expected %q = %q.(%s) obligation but got %q = %q.(%s)", eID, es, eot, ID, s, ot)
		}
	}

	p = NewPolicySet("test", false, Target{},
		[]Evaluable{
			makeSimplePermitPolicyWithObligations(
				"first",
				makeSingleStringObligation("order", "first"),
			),
			makeSimplePermitPolicyWithObligations(
				"second",
				makeSingleStringObligation("order", "second"),
			),
			makeSimplePermitPolicyWithObligations(
				"third",
				makeSingleStringObligation("order", "third"),
			),
		},
		makeMapperPCA, MapperPCAParams{
			Argument:  MakeDesignator("f", ft),
			Order:     MapperPCAExternalOrder,
			Algorithm: firstApplicableEffectPCA{},
		},
		nil,
	)

	r = p.Calculate(c)
	if len(r.Obligations) != 1 {
		t.Errorf("Expected the only obligation got %#v", r.Obligations)
	} else {
		ID, ot, s, err := r.Obligations[0].Serialize(c)
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		eID := "order"
		eot := TypeString
		es := "third"
		if ID != eID || ot != eot.GetKey() || s != es {
			t.Errorf("Expected %q = %q.(%s) obligation but got %q = %q.(%s)", eID, es, eot, ID, s, ot)
		}
	}
}

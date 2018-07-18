package pdp

import "testing"

func TestMapperPCAOrders(t *testing.T) {
	if totalMapperPCAOrders != len(MapperPCAOrderNames) {
		t.Errorf("Expected total number of order values to be equal to number of their names "+
			"but got totalMapperPCAOrders = %d and len(MapperPCAOrderNames) = %d",
			totalMapperPCAOrders, len(MapperPCAOrderNames))
	}
}

func TestMapperPCAOrdering(t *testing.T) {
	c := &Context{
		a: map[string]interface{}{
			"k": MakeListOfStringsValue([]string{"third", "first", "second"}),
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
			Argument:  MakeListOfStringsDesignator("k"),
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
			Argument:  MakeListOfStringsDesignator("k"),
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

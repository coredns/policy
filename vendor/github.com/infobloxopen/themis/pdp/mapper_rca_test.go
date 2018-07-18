package pdp

import "testing"

func TestMapperRCAOrders(t *testing.T) {
	if totalMapperRCAOrders != len(MapperRCAOrderNames) {
		t.Errorf("Expected total number of order values to be equal to number of their names "+
			"but got totalMapperRCAOrders = %d and len(MapperRCAOrderNames) = %d",
			totalMapperRCAOrders, len(MapperRCAOrderNames))
	}
}

func TestMapperRCAOrdering(t *testing.T) {
	c := &Context{
		a: map[string]interface{}{
			"k": MakeListOfStringsValue([]string{"third", "first", "second"}),
		},
	}

	p := NewPolicy("test", false, Target{},
		[]*Rule{
			makeSimplePermitRuleWithObligations(
				"first",
				makeSingleStringObligation("order", "first"),
			),
			makeSimplePermitRuleWithObligations(
				"second",
				makeSingleStringObligation("order", "second"),
			),
			makeSimplePermitRuleWithObligations(
				"third",
				makeSingleStringObligation("order", "third"),
			),
		},
		makeMapperRCA, MapperRCAParams{
			Argument:  MakeListOfStringsDesignator("k"),
			Order:     MapperRCAInternalOrder,
			Algorithm: firstApplicableEffectRCA{},
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

	p = NewPolicy("test", false, Target{},
		[]*Rule{
			makeSimplePermitRuleWithObligations(
				"first",
				makeSingleStringObligation("order", "first"),
			),
			makeSimplePermitRuleWithObligations(
				"second",
				makeSingleStringObligation("order", "second"),
			),
			makeSimplePermitRuleWithObligations(
				"third",
				makeSingleStringObligation("order", "third"),
			),
		},
		makeMapperRCA, MapperRCAParams{
			Argument:  MakeListOfStringsDesignator("k"),
			Order:     MapperRCAExternalOrder,
			Algorithm: firstApplicableEffectRCA{},
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

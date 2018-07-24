package pdp

import (
	"encoding/json"
	"testing"
)

func TestFlagsMapperRCAOrdering(t *testing.T) {
	ft, err := NewFlagsType("flags", "third", "first", "second")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	c := &Context{
		a: map[string]interface{}{
			"f": MakeFlagsValue8(7, ft),
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
			Argument:  MakeDesignator("f", ft),
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
			Argument:  MakeDesignator("f", ft),
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

func TestFlagsMapperRCAMarshal(t *testing.T) {
	const (
		expectEmptyFlagMapperRCAJSON = `{"type":"flagsMapperRCA","def":"\"\"","err":"\"\"","alg":{"type":"firstApplicableEffectRCA"}}`
		expectFlagMapperRCAJSON      = `{"type":"flagsMapperRCA","def":"\"first\"","err":"\"second\"","alg":{"type":"firstApplicableEffectRCA"}}`
	)

	ft, err := NewFlagsType("flags", "third", "first", "second")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}
	rules := []*Rule{
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
	}
	algParam := MapperRCAParams{
		Argument:  AttributeDesignator{a: Attribute{id: "f", t: ft}},
		Order:     MapperRCAExternalOrder,
		Algorithm: firstApplicableEffectRCA{},
	}
	alg := makeMapperRCA(rules, algParam)
	b, _ := json.Marshal(alg)
	if expectEmptyFlagMapperRCAJSON != string(b) {
		t.Errorf("Expected marshalled %s\nGot %s", expectEmptyFlagMapperRCAJSON, string(b))
	}

	algParam2 := MapperRCAParams{
		Argument:  AttributeDesignator{a: Attribute{id: "f", t: ft}},
		Order:     MapperRCAExternalOrder,
		Algorithm: firstApplicableEffectRCA{},
		Def:       "first",
		DefOk:     true,
		Err:       "second",
		ErrOk:     true,
	}
	alg = makeMapperRCA(rules, algParam2)
	b, _ = json.Marshal(alg)
	if expectFlagMapperRCAJSON != string(b) {
		t.Errorf("Expected marshalled %s\nGot %s", expectFlagMapperRCAJSON, string(b))
	}
}

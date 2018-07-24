package pdp

import (
	"bytes"
	"sort"
	"strings"
	"testing"
)

func TestPolicySet(t *testing.T) {
	c := &Context{
		a: map[string]interface{}{
			"missing-type":   MakeBooleanValue(false),
			"test-string":    MakeStringValue("test"),
			"example-string": MakeStringValue("example")}}

	testID := "Test Policy"

	p := makeSimplePolicySet(testID)
	ID, ok := p.GetID()
	if !ok {
		t.Errorf("Expected policy set ID %q but got hidden policy set", testID)
	} else if ID != testID {
		t.Errorf("Expected policy set ID %q but got %q", testID, ID)
	}

	r := p.Calculate(c)
	if r.Effect != EffectNotApplicable {
		t.Errorf("Expected %q for empty policy set but got %q",
			effectNames[EffectNotApplicable], effectNames[r.Effect])
	}

	p = &PolicySet{
		id:        testID,
		target:    makeSimpleStringTarget("missing", "test"),
		algorithm: firstApplicableEffectPCA{}}

	r = p.Calculate(c)
	if r.Effect != EffectNotApplicable {
		t.Errorf("Expected %q for policy set with FirstApplicableEffectPCA and not found attribute but got %q",
			effectNames[EffectNotApplicable], effectNames[r.Effect])
	}

	_, ok = r.Status.(*missingAttributeError)
	if !ok {
		t.Errorf("Expected missing attribute status for policy set with FirstApplicableEffectPCA and "+
			"not found attribute but got %T (%s)", r.Status, r.Status)
	}

	p = &PolicySet{
		id:        testID,
		target:    makeSimpleStringTarget("missing-type", "test"),
		algorithm: firstApplicableEffectPCA{}}

	r = p.Calculate(c)
	if r.Effect != EffectNotApplicable {
		t.Errorf("Expected %q for policy set with FirstApplicableEffectPCA and attribute with wrong type but got %q",
			effectNames[EffectNotApplicable], effectNames[r.Effect])
	}

	_, ok = r.Status.(*missingAttributeError)
	if !ok {
		t.Errorf("Expected missing attribute status for policy with FirstApplicableEffectPCA and "+
			"attribute with wrong type but got %T (%s)", r.Status, r.Status)
	}

	p = &PolicySet{
		id:        testID,
		target:    makeSimpleStringTarget("example-string", "test"),
		algorithm: firstApplicableEffectPCA{}}

	r = p.Calculate(c)
	if r.Effect != EffectNotApplicable {
		t.Errorf("Expected %q for policy set with FirstApplicableEffectPCA and "+
			"attribute with not maching value but got %q",
			effectNames[EffectNotApplicable], effectNames[r.Effect])
	}

	if r.Status != nil {
		t.Errorf("Expected no error status for policy set with FirstApplicableEffectPCA and "+
			"attribute with not maching value but got %T (%s)", r.Status, r.Status)
	}

	p = &PolicySet{
		id:     testID,
		target: makeSimpleStringTarget("test-string", "test"),
		policies: []Evaluable{
			makeSimpleHiddenPolicy(makeSimpleHiddenRule(EffectPermit)),
		},
		obligations: makeSingleStringObligation("obligation", "test"),
		algorithm:   firstApplicableEffectPCA{}}

	r = p.Calculate(c)
	if r.Effect != EffectPermit {
		t.Errorf("Expected %q for policy with rule and obligations but got %q",
			effectNames[EffectPermit], effectNames[r.Effect])
	}

	if r.Status != nil {
		t.Errorf("Expected no error status for policy rule and obligations but got %T (%s)",
			r.Status, r.Status)
	}

	defaultPolicy := makeSimplePolicy("Default", makeSimpleHiddenRule(EffectDeny))
	errorPolicy := makeSimplePolicy("Error", makeSimpleHiddenRule(EffectDeny))
	permitPolicy := makeSimplePolicy("Permit", makeSimpleHiddenRule(EffectPermit))
	p = &PolicySet{
		id:       testID,
		policies: []Evaluable{defaultPolicy, errorPolicy, permitPolicy},
		algorithm: makeMapperPCA(
			[]Evaluable{defaultPolicy, errorPolicy, permitPolicy},
			MapperPCAParams{
				Argument: MakeSetOfStringsDesignator("x"),
				DefOk:    true,
				Def:      "Default",
				ErrOk:    true,
				Err:      "Error",
				Algorithm: makeMapperPCA(
					nil,
					MapperPCAParams{
						Argument: MakeStringDesignator("y")})})}

	c = &Context{
		a: map[string]interface{}{
			"x": MakeSetOfStringsValue(newStrTree("Permit", "Default")),
			"y": MakeStringValue("Permit")}}

	r = p.Calculate(c)
	if r.Effect != EffectPermit {
		t.Errorf("Expected %q for policy with rule and obligations but got %q",
			effectNames[EffectPermit], effectNames[r.Effect])
	}

	if r.Status != nil {
		t.Errorf("Expected no error status for policy rule and obligations but got %T (%s)",
			r.Status, r.Status)
	}

	c = &Context{
		a: map[string]interface{}{
			"x": MakeSetOfStringsValue(newStrTree("Permit", "Default")),
			"y": MakeSetOfStringsValue(newStrTree("Permit", "Default"))}}

	r = p.Calculate(c)
	if r.Effect != EffectIndeterminate {
		t.Errorf("Expected %q for policy with rule and obligations but got %q",
			effectNames[EffectIndeterminate], effectNames[r.Effect])
	}

	_, ok = r.Status.(*missingAttributeError)
	if !ok {
		t.Errorf("Expected missing attribute status for policy with rule and obligations but got %T (%s)",
			r.Status, r.Status)
	}
}

func TestPolicySetAppend(t *testing.T) {
	testPermitPol := makeSimplePolicySet("test",
		makeSimplePolicy("permit", makeSimpleRule("permit", EffectPermit)),
	)

	p := makeSimplePolicySet("test")
	p.ord = 5

	newE, err := p.Append([]string{}, testPermitPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if p == newP {
			t.Errorf("Expected different new policy set but got the same")
		}

		if newP.ord != p.ord {
			t.Errorf("Expected unchanged order %d but got %d", p.ord, newP.ord)
		}

		if len(newP.policies) == 1 {
			p := newP.policies[0]
			ord := p.getOrder()
			if ord != 0 {
				t.Errorf("Expected index of the only index to be 0 but got %d", ord)
			}
		} else {
			t.Errorf("Expected only appended item but got %d items", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = p.Append([]string{"test"}, testPermitPol)
	if err == nil {
		t.Errorf("Expected error but got policy %#v", newE)
	} else if _, ok := err.(*missingPolicySetChildError); !ok {
		t.Errorf("Expected *missingPolicySetChildError but got %T (%s)", err, err)
	}

	newE, err = p.Append([]string{}, &Rule{id: "permit", effect: EffectPermit})
	if err == nil {
		t.Errorf("Expected error but got policy %#v", newE)
	} else if _, ok := err.(*invalidPolicySetItemTypeError); !ok {
		t.Errorf("Expected *invalidPolicySetItemTypeError but got %T (%s)", err, err)
	}

	newE, err = p.Append([]string{}, &PolicySet{hidden: true, algorithm: firstApplicableEffectPCA{}})
	if err == nil {
		t.Errorf("Expected error but got policy %#v", newE)
	} else if _, ok := err.(*hiddenPolicyAppendError); !ok {
		t.Errorf("Expected *hiddenPolicyAppendError but got %T (%s)", err, err)
	}

	p = makeSimpleHiddenPolicySet()
	newE, err = p.Append([]string{}, testPermitPol)
	if err == nil {
		t.Errorf("Expected error but got policy %#v", newE)
	} else if _, ok := err.(*hiddenPolicySetModificationError); !ok {
		t.Errorf("Expected *hiddenPolicySetModificationError but got %T (%s)", err, err)
	}

	p = makeSimplePolicySet("test", makeSimplePolicy("test"))
	newE, err = p.Append([]string{"test"},
		makeSimpleHiddenRule(EffectPermit),
	)
	if err == nil {
		t.Errorf("Expected error but got policy %#v", newE)
	} else if _, ok := err.(*hiddenRuleAppendError); !ok {
		t.Errorf("Expected *hiddenRuleAppendError but got %T (%s)", err, err)
	}

	_, err = p.Append([]string{"test"},
		makeSimpleRule("test", EffectPermit),
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	testFirstPol := makeSimplePolicy("first", makeSimpleRule("deny", EffectDeny))
	testSecondPol := makeSimplePolicy("second", makeSimpleRule("deny", EffectDeny))
	testThirdPermitPol := makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit))
	testThirdDenyPol := makeSimplePolicy("third", makeSimpleRule("deny", EffectDeny))

	p = makeSimplePolicySet("test",
		makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
		makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
	)
	if len(p.policies) == 2 {
		for i, e := range p.policies {
			ord := e.getOrder()
			if ord != i {
				id, ok := e.GetID()
				if !ok {
					id = "hidden"
				}

				t.Errorf("Expected %q policy to get %d order but got %d", id, i, ord)
			}
		}
	} else {
		t.Errorf("Expected 2 policies in the policy set but got %d", len(p.policies))
	}

	newE, err = p.Append([]string{}, testThirdPermitPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 3 {
			e := newP.policies[2]
			if p, ok := e.(*Policy); ok {
				if p.id != "third" {
					t.Errorf("Expected \"third\" policy added to the end but got %q", p.id)
				}

				if p.ord != 2 {
					t.Errorf("Expected the last rule to get order 2 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as third item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected three policies after append but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = newE.Append([]string{}, testFirstPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 3 {
			e := newP.policies[0]
			if p, ok := e.(*Policy); ok {
				if p.id != "first" {
					t.Errorf("Expected \"first\" policy replaced at the begining but got %q", p.id)
				} else if len(p.rules) == 1 {
					r := p.rules[0]
					if r.effect != EffectDeny {
						t.Errorf("Expected \"first\" policy became deny but it's still %s", effectNames[r.effect])
					}
				} else {
					t.Errorf("Expected \"first\" policy to have the only rule but got %d", len(p.rules))
				}

				if p.ord != 0 {
					t.Errorf("Expected the first policy to keep order 0 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as first item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected three policies after append but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = newE.Append([]string{}, testSecondPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 3 {
			e := newP.policies[1]
			if p, ok := e.(*Policy); ok {
				if p.id != "second" {
					t.Errorf("Expected \"second\" policy replaced at the middle but got %q", p.id)
				} else if len(p.rules) == 1 {
					r := p.rules[0]
					if r.effect != EffectDeny {
						t.Errorf("Expected \"second\" policy became deny but it's still %s", effectNames[r.effect])
					}
				} else {
					t.Errorf("Expected \"second\" policy to have the only rule but got %d", len(p.rules))
				}

				if p.ord != 1 {
					t.Errorf("Expected second policy to keep order 1 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as second item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected three policies after append but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = newE.Append([]string{}, testThirdDenyPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 3 {
			e := newP.policies[2]
			if p, ok := e.(*Policy); ok {
				if p.id != "third" {
					t.Errorf("Expected \"third\" policy replaced at the end but got %q", p.id)
				} else if len(p.rules) == 1 {
					r := p.rules[0]
					if r.effect != EffectDeny {
						t.Errorf("Expected \"third\" policy became deny but it's still %s", effectNames[r.effect])
					}
				} else {
					t.Errorf("Expected \"third\" policy to have the only rule but got %d", len(p.rules))
				}

				if p.ord != 2 {
					t.Errorf("Expected third policy to keep order 2 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as third item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected three policies after append but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy but got %T (%#v)", newE, newE)
	}

	testFourthPol := makeSimplePolicy("fourth", makeSimpleRule("permit", EffectPermit))

	p = NewPolicySet("test", false, Target{},
		[]Evaluable{
			makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit)),
		},
		makeMapperPCA, MapperPCAParams{
			Argument: MakeStringDesignator("k"),
			DefOk:    true,
			Def:      "first",
			ErrOk:    true,
			Err:      "second"},
		nil)
	newE, err = p.Append([]string{}, testFourthPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 4 {
			e := newP.policies[3]
			if p, ok := e.(*Policy); ok {
				if p.id != "fourth" {
					t.Errorf("Expected \"fourth\" policy added to the end but got %q", p.id)
				}

				if p.ord != 3 {
					t.Errorf("Expected fourth policy to get order 3 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as fourth item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected four policies after append but got %d", len(newP.policies))
		}

		assertMapperPCAMapKeys(newP.algorithm, "after insert \"fourth\"", t, "first", "fourth", "second", "third")
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = newE.Append([]string{}, testFirstPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 4 {
			e := newP.policies[0]
			if p, ok := e.(*Policy); ok {
				if p.id != "first" {
					t.Errorf("Expected \"first\" policy replaced at the begining but got %q", p.id)
				} else if len(p.rules) == 1 {
					r := p.rules[0]
					if r.effect != EffectDeny {
						t.Errorf("Expected \"first\" policy became deny but it's still %s", effectNames[r.effect])
					}
				} else {
					t.Errorf("Expected \"first\" policy to have the only rule but got %d", len(p.rules))
				}

				if p.ord != 0 {
					t.Errorf("Expected the first policy to keep order 0 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as first item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected four policies after append but got %d", len(newP.policies))
		}

		assertMapperPCAMapKeys(newP.algorithm, "after insert another \"first\"", t,
			"first", "fourth", "second", "third")

		if m, ok := newP.algorithm.(mapperPCA); ok {
			if m.def != testFirstPol {
				t.Errorf("Expected default policy to be new \"first\" policy %p but got %p", testFirstPol, m.def)
			}
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	ft, err := NewFlagsType("flags", "first", "second", "third", "fourth")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	p = NewPolicySet("test", false, Target{},
		[]Evaluable{
			makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit)),
		},
		makeMapperPCA, MapperPCAParams{
			Argument: MakeDesignator("f", ft),
			DefOk:    true,
			Def:      "first",
			ErrOk:    true,
			Err:      "second"},
		nil)
	newE, err = p.Append([]string{}, testFourthPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 4 {
			e := newP.policies[3]
			if p, ok := e.(*Policy); ok {
				if p.id != "fourth" {
					t.Errorf("Expected \"fourth\" policy added to the end but got %q", p.id)
				}

				if p.ord != 3 {
					t.Errorf("Expected fourth policy to get order 3 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as fourth item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected four policies after append but got %d", len(newP.policies))
		}

		assertFlagsMapperPCAMapKeys(newP.algorithm, "after insert \"fourth\"", t, "first", "second", "third", "fourth")
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	testFifthPol := makeSimplePolicy("fifth", makeSimpleRule("permit", EffectPermit))
	newE, err = newE.Append([]string{}, testFifthPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 5 {
			e := newP.policies[4]
			if p, ok := e.(*Policy); ok {
				if p.id != "fifth" {
					t.Errorf("Expected \"fifth\" policy added to the end but got %q", p.id)
				}

				if p.ord != 4 {
					t.Errorf("Expected fourth policy to get order 4 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as fifth item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected five policies after append but got %d", len(newP.policies))
		}

		assertFlagsMapperPCAMapKeys(newP.algorithm, "after insert \"fifth\"", t, "first", "second", "third", "fourth")
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = newE.Append([]string{}, testFirstPol)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 5 {
			e := newP.policies[0]
			if p, ok := e.(*Policy); ok {
				if p.id != "first" {
					t.Errorf("Expected \"first\" policy replaced at the begining but got %q", p.id)
				} else if len(p.rules) == 1 {
					r := p.rules[0]
					if r.effect != EffectDeny {
						t.Errorf("Expected \"first\" policy became deny but it's still %s", effectNames[r.effect])
					}
				} else {
					t.Errorf("Expected \"first\" policy to have the only rule but got %d", len(p.rules))
				}

				if p.ord != 0 {
					t.Errorf("Expected the first policy to keep order 0 but got %d", p.ord)
				}
			} else {
				t.Errorf("Expected policy as first item of policy set but got %T (%#v)", e, e)
			}
		} else {
			t.Errorf("Expected four policies after append but got %d", len(newP.policies))
		}

		assertFlagsMapperPCAMapKeys(newP.algorithm, "after insert other \"first\"", t,
			"first", "second", "third", "fourth")

		if m, ok := newP.algorithm.(mapperPCA); ok {
			if m.def != testFirstPol {
				t.Errorf("Expected default policy to be new \"first\" policy %p but got %p", testFirstPol, m.def)
			}
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}
}

func TestPolicySetDelete(t *testing.T) {
	p := makeSimplePolicySet("test",
		makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
		makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
		makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit)),
	)
	if len(p.policies) == 3 {
		for i, e := range p.policies {
			ord := e.getOrder()
			if ord != i {
				id, ok := e.GetID()
				if !ok {
					id = "hidden"
				}

				t.Errorf("Expected %q policy to get %d order but got %d", id, i, ord)
			}
		}
	} else {
		t.Errorf("Expected 3 policies in the policy set but got %d", len(p.policies))
	}

	newE, err := p.Delete([]string{"second"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 2 {
			e1 := newP.policies[0]
			e3 := newP.policies[1]

			p1, ok1 := e1.(*Policy)
			p3, ok3 := e3.(*Policy)
			if ok1 && ok3 {
				if p1.id != "first" || p3.id != "third" {
					t.Errorf("Expected \"first\" and \"third\" policies remaining but got %q and %q", p1.id, p3.id)
				}

				if p1.ord != 0 || p3.ord != 2 {
					t.Errorf("Expected remaining policies to keep their orders but got %d and %d", p1.ord, p3.ord)
				}
			} else {
				t.Errorf("Expected two policies after delete but got %T and %T", e1, e3)
			}
		} else {
			t.Errorf("Expected two policies after delete but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = p.Delete([]string{"first"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 2 {
			e2 := newP.policies[0]
			e3 := newP.policies[1]

			p2, ok2 := e2.(*Policy)
			p3, ok3 := e3.(*Policy)
			if ok2 && ok3 {
				if p2.id != "second" || p3.id != "third" {
					t.Errorf("Expected \"second\" and \"third\" policies remaining but got %q and %q", p2.id, p3.id)
				}

				if p2.ord != 1 || p3.ord != 2 {
					t.Errorf("Expected remaining policies to keep their orders but got %d and %d", p2.ord, p3.ord)
				}
			} else {
				t.Errorf("Expected two policies after delete but got %T and %T", e2, e3)
			}
		} else {
			t.Errorf("Expected two policies after delete but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = p.Delete([]string{"third"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 2 {
			e1 := newP.policies[0]
			e2 := newP.policies[1]

			p1, ok1 := e1.(*Policy)
			p2, ok2 := e2.(*Policy)
			if ok1 && ok2 {
				if p1.id != "first" || p2.id != "second" {
					t.Errorf("Expected \"first\" and \"second\" policies remaining but got %q and %q", p1.id, p2.id)
				}

				if p1.ord != 0 || p2.ord != 1 {
					t.Errorf("Expected remaining policies to keep their orders but got %d and %d", p1.ord, p2.ord)
				}
			} else {
				t.Errorf("Expected two policies after delete but got %T and %T", e1, e2)
			}
		} else {
			t.Errorf("Expected two policies after delete but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = p.Delete([]string{"first", "permit"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 3 {
			p := newP.policies[0]
			if p, ok := p.(*Policy); ok {
				if p.id == "first" {
					if len(p.rules) > 0 {
						t.Errorf("Expected no rules after nested delete but got %d", len(p.rules))
					}
				} else {
					t.Errorf("Expected \"first\" policy at the beginning but got %q", p.id)
				}
			} else {
				t.Errorf("Expected policy as first item of policy set but got %T (%#v)", newP, newP)
			}
		} else {
			t.Errorf("Expected three policies after nested delete but got %d", len(newP.policies))
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = p.Delete([]string{"fourth"})
	if err == nil {
		t.Errorf("Expected error but got new policy %T (%#v)", newE, newE)
	} else if _, ok := err.(*missingPolicySetChildError); !ok {
		t.Errorf("Expected *missingPolicySetChildError but got %T (%s)", err, err)
	}

	newE, err = p.Delete([]string{"fourth", "permit"})
	if err == nil {
		t.Errorf("Expected error but got new policy %T (%#v)", newE, newE)
	} else if _, ok := err.(*missingPolicySetChildError); !ok {
		t.Errorf("Expected *missingPolicySetChildError but got %T (%s)", err, err)
	}

	newE, err = p.Delete([]string{"first", "deny"})
	if err == nil {
		t.Errorf("Expected error but got new policy %T (%#v)", newE, newE)
	} else if _, ok := err.(*missingPolicyChildError); !ok {
		t.Errorf("Expected *missingPolicyChildError but got %T (%s)", err, err)
	}

	p = makeSimpleHiddenPolicySet(
		makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
	)
	newE, err = p.Delete([]string{"first"})
	if err == nil {
		t.Errorf("Expected error but got new policy %T (%#v)", newE, newE)
	} else if _, ok := err.(*hiddenPolicySetModificationError); !ok {
		t.Errorf("Expected *hiddenPolicySetModificationError but got %T (%s)", err, err)
	}

	p = makeSimplePolicySet("test",
		makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
	)
	newE, err = p.Delete([]string{})
	if err == nil {
		t.Errorf("Expected error but got new policy %T (%#v)", newE, newE)
	} else if _, ok := err.(*tooShortPathPolicySetModificationError); !ok {
		t.Errorf("Expected *tooShortPathPolicySetModificationError but got %T (%s)", err, err)
	}

	p = NewPolicySet("test", false, Target{},
		[]Evaluable{
			makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit)),
		},
		makeMapperPCA, MapperPCAParams{
			Argument: MakeStringDesignator("k"),
			DefOk:    true,
			Def:      "first",
			ErrOk:    true,
			Err:      "second"},
		nil)
	newE, err = p.Delete([]string{"second"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 2 {
			e1 := newP.policies[0]
			e3 := newP.policies[1]

			p1, ok1 := e1.(*Policy)
			p3, ok3 := e3.(*Policy)
			if ok1 && ok3 {
				if p1.id != "first" || p3.id != "third" {
					t.Errorf("Expected \"first\" and \"third\" policies remaining but got %q and %q", p1.id, p3.id)
				}

				if p1.ord != 0 || p3.ord != 2 {
					t.Errorf("Expected remaining policies to keep their orders but got %d and %d", p1.ord, p3.ord)
				}
			} else {
				t.Errorf("Expected two policies after delete but got %T and %T", e1, e3)
			}
		} else {
			t.Errorf("Expected two policies after delete but got %d", len(newP.policies))
		}

		assertMapperPCAMapKeys(newP.algorithm, "after deletion", t, "first", "third")

		if m, ok := newP.algorithm.(mapperPCA); ok {
			if m.err != nil {
				t.Errorf("Expected error policy to be nil but got %p", m.err)
			}
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	ft, err := NewFlagsType("flags", "first", "second", "third", "fourth")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	p = NewPolicySet("test", false, Target{},
		[]Evaluable{
			makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("fifth", makeSimpleRule("permit", EffectPermit)),
		},
		makeMapperPCA, MapperPCAParams{
			Argument: MakeDesignator("f", ft),
			DefOk:    true,
			Def:      "first",
			ErrOk:    true,
			Err:      "second"},
		nil)
	newE, err = p.Delete([]string{"second"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 3 {
			e1 := newP.policies[0]
			e3 := newP.policies[1]
			e5 := newP.policies[2]

			p1, ok1 := e1.(*Policy)
			p3, ok3 := e3.(*Policy)
			p5, ok5 := e5.(*Policy)
			if ok1 && ok3 && ok5 {
				if p1.id != "first" || p3.id != "third" || p5.id != "fifth" {
					t.Errorf("Expected \"first\", \"third\" and \"fifth\" policies remaining but got %q, %q and %q",
						p1.id, p3.id, p5.id)
				}

				if p1.ord != 0 || p3.ord != 2 || p5.ord != 3 {
					t.Errorf("Expected remaining policies to keep their orders but got %d, %d and %d",
						p1.ord, p3.ord, p5.ord)
				}
			} else {
				t.Errorf("Expected three policies after delete but got %T, %T and %T", e1, e3, e5)
			}
		} else {
			t.Errorf("Expected three policies after delete but got %d", len(newP.policies))
		}

		assertFlagsMapperPCAMapKeys(newP.algorithm, "after \"second\" deletion from flags policy set", t,
			"first", "third")

		if m, ok := newP.algorithm.(flagsMapperPCA); ok {
			if m.err != nil {
				t.Errorf("Expected error policy to be nil but got %p", m.err)
			}
		}
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}

	newE, err = newE.Delete([]string{"fifth"})
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if newP, ok := newE.(*PolicySet); ok {
		if len(newP.policies) == 2 {
			e1 := newP.policies[0]
			e3 := newP.policies[1]

			p1, ok1 := e1.(*Policy)
			p3, ok3 := e3.(*Policy)
			if ok1 && ok3 {
				if p1.id != "first" || p3.id != "third" {
					t.Errorf("Expected \"first\" and \"third\" policies remaining but got %q and %q", p1.id, p3.id)
				}

				if p1.ord != 0 || p3.ord != 2 {
					t.Errorf("Expected remaining policies to keep their orders but got %d and %d", p1.ord, p3.ord)
				}
			} else {
				t.Errorf("Expected two policies after delete but got %T and %T", e1, e3)
			}
		} else {
			t.Errorf("Expected two policies after delete but got %d", len(newP.policies))
		}

		assertFlagsMapperPCAMapKeys(newP.algorithm, "after \"fifth\" deletion from flags policy set", t,
			"first", "third")
	} else {
		t.Errorf("Expected new policy set but got %T (%#v)", newE, newE)
	}
}

func TestSortPoliciesByOrder(t *testing.T) {
	policies := []Evaluable{
		&PolicySet{
			ord: 1,
			id:  "second",
		},
		&PolicySet{
			ord: 3,
			id:  "fourth",
		},
		&Policy{
			ord: 0,
			id:  "first",
		},
		&Policy{
			ord: 2,
			id:  "third",
		},
	}

	sort.Sort(byPolicyOrder(policies))

	ids := make([]string, len(policies))
	for i, p := range policies {
		id, ok := p.GetID()
		if !ok {
			id = "hidden"
		}

		ids[i] = id
	}
	s := strings.Join(ids, ", ")
	e := "first, second, third, fourth"
	if s != e {
		t.Errorf("Expected policies in order \"%s\" but got \"%s\"", e, s)
	}
}

func TestPolicySetMarshalWithDepth(t *testing.T) {
	var (
		buf  bytes.Buffer
		buf2 bytes.Buffer
		buf3 bytes.Buffer
		p    = makeSimplePolicySet("test",
			makeSimplePolicy("first", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("second", makeSimpleRule("permit", EffectPermit)),
			makeSimplePolicy("third", makeSimpleRule("permit", EffectPermit)),
		)
	)

	// bad depth
	err := p.MarshalWithDepth(&buf, -1)
	expectErr := newMarshalInvalidDepthError(-1)
	if err == nil {
		t.Errorf("Expecting error %v, got nil error", expectErr)
	} else if err.Error() != expectErr.Error() {
		t.Errorf("Expecting error %v, got %v", expectErr, err)
	}

	// depth = 0, visible policy
	policySetExtra := `"target":{},"obligations":null,"algorithm":{"type":"firstApplicableEffectPCA"},`
	expectMarshal := `{"ord":0,"id":"test",` + policySetExtra + `"policies":[]}`
	err = p.MarshalWithDepth(&buf, 0)
	if err != nil {
		t.Errorf("Expecting no error, got %v", err)
	} else {
		gotMarshal := buf.String()
		if 0 != strings.Compare(gotMarshal, expectMarshal) {
			t.Errorf("Expecting marshal output %s, got %s", expectMarshal, gotMarshal)
		}
	}

	// show children, visible policy
	policyExtra := `"target":{},"obligations":null,"algorithm":{"type":"firstApplicableEffectRCA"}`
	cRule := `,"rules":[{"ord":0,"id":"permit","target":{},"obligations":null,"effect":"Permit"}]}`
	expectChildren := `{"ord":0,"id":"first",` + policyExtra + cRule +
		`,{"ord":1,"id":"second",` + policyExtra + cRule +
		`,{"ord":2,"id":"third",` + policyExtra + cRule
	expectWithC := `{"ord":0,"id":"test",` + policySetExtra + `"policies":[` + expectChildren + `]}`
	err = p.MarshalWithDepth(&buf2, 2)
	if err != nil {
		t.Errorf("Expecting no error, got %v", err)
	} else {
		gotMarshal := buf2.String()
		if 0 != strings.Compare(gotMarshal, expectWithC) {
			t.Errorf("Expecting marshal output %s, got %s",
				expectWithC, gotMarshal)
		}
	}

	// depth beyond maximum, visible policy
	err = p.MarshalWithDepth(&buf3, 100)
	if err != nil {
		t.Errorf("Expecting no error, got %v", err)
	} else {
		gotMarshal := buf3.String()
		if 0 != strings.Compare(gotMarshal, expectWithC) {
			t.Errorf("Expecting marshal output %s, got %s",
				expectWithC, gotMarshal)
		}
	}
}

func makeSimplePolicySet(ID string, policies ...Evaluable) *PolicySet {
	return NewPolicySet(
		ID, false,
		Target{},
		policies,
		makeFirstApplicableEffectPCA,
		nil,
		nil,
	)
}

func makeSimpleHiddenPolicySet(policies ...Evaluable) *PolicySet {
	return NewPolicySet(
		"", true,
		Target{},
		policies,
		makeFirstApplicableEffectPCA,
		nil,
		nil,
	)
}

func assertMapperPCAMapKeys(a PolicyCombiningAlg, desc string, t *testing.T, expected ...string) {
	if m, ok := a.(mapperPCA); ok {
		keys := []string{}
		for p := range m.policies.Enumerate() {
			keys = append(keys, p.Key)
		}

		assertStrings(keys, expected, desc, t)
	} else {
		t.Errorf("Expected mapper as policy combining algorithm but got %T for %s", a, desc)
	}
}

func assertFlagsMapperPCAMapKeys(a PolicyCombiningAlg, desc string, t *testing.T, expected ...string) {
	if m, ok := a.(flagsMapperPCA); ok {
		keys := []string{}
		for _, p := range m.policies {
			if p != nil {
				if id, ok := p.GetID(); ok {
					keys = append(keys, id)
				}
			}
		}

		assertStrings(keys, expected, desc, t)
	} else {
		t.Errorf("Expected flags mapper as policy combining algorithm but got %T for %s", a, desc)
	}
}

package pdp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestStorage(t *testing.T) {
	root := &PolicySet{
		id: "test",
		policies: []Evaluable{
			&Policy{
				id:        "first",
				rules:     []*Rule{{id: "permit", effect: EffectPermit}},
				algorithm: firstApplicableEffectRCA{}}},
		algorithm: firstApplicableEffectPCA{}}

	s := NewPolicyStorage(root, Symbols{}, nil)
	sr := s.Root()
	if sr != root {
		t.Errorf("Expected stored root policy to be exactly root policy but got different")
	}
}

func TestStorageNewTransaction(t *testing.T) {
	initialTag := uuid.New()

	s := &PolicyStorage{tag: &initialTag}
	tr, err := s.NewTransaction(&initialTag)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if tr == nil {
		t.Errorf("Expected transaction but got nothing")
	}

	s = &PolicyStorage{}
	tr, err = s.NewTransaction(&initialTag)
	if err == nil {
		t.Errorf("Expected error but got transaction %#v", tr)
	} else if _, ok := err.(*UntaggedPolicyModificationError); !ok {
		t.Errorf("Expected *untaggedPolicyModificationError but got %T (%s)", err, err)
	}

	s = &PolicyStorage{tag: &initialTag}
	tr, err = s.NewTransaction(nil)
	if err == nil {
		t.Errorf("Expected error but got transaction %#v", tr)
	} else if _, ok := err.(*MissingPolicyTagError); !ok {
		t.Errorf("Expected *missingPolicyTagError but got %T (%s)", err, err)
	}

	otherTag := uuid.New()
	s = &PolicyStorage{tag: &initialTag}
	tr, err = s.NewTransaction(&otherTag)
	if err == nil {
		t.Errorf("Expected error but got transaction %#v", tr)
	} else if _, ok := err.(*PolicyTagsNotMatchError); !ok {
		t.Errorf("Expected *policyTagsNotMatchError but got %T (%s)", err, err)
	}
}

func TestStorageCommitTransaction(t *testing.T) {
	initialTag := uuid.New()
	newTag := uuid.New()

	s := &PolicyStorage{tag: &initialTag}
	tr, err := s.NewTransaction(&initialTag)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		u := NewPolicyUpdate(initialTag, newTag)
		err := tr.Apply(u)
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else {
			newS, err := tr.Commit()
			if err != nil {
				t.Errorf("Expected no error but got %s", err)
			} else {
				if &newS == &s {
					t.Errorf("Expected other storage instance but got the same")
				}

				if newS.tag.String() != newTag.String() {
					t.Errorf("Expected tag %s but got %s", newTag.String(), newS.tag.String())
				}
			}
		}
	}

	tr, err = s.NewTransaction(&initialTag)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		tr.err = newUnknownPolicyUpdateOperationError(-1)
		s, err := tr.Commit()
		if err == nil {
			t.Errorf("Expected error but got storage %#v", s)
		} else if _, ok := err.(*failedPolicyTransactionError); !ok {
			t.Errorf("Expected *failedPolicyTransactionError but got %T (%s)", err, err)
		}
	}
}

func TestStorageModifications(t *testing.T) {
	tag := uuid.New()

	s := &PolicyStorage{
		tag: &tag,
		policies: &Policy{
			id:        "test",
			algorithm: firstApplicableEffectRCA{}}}
	tr, err := s.NewTransaction(&tag)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		err := tr.appendItem([]string{"test"}, &Rule{id: "permit", effect: EffectPermit})
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		err = tr.appendItem(nil, &Rule{id: "permit", effect: EffectPermit})
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*invalidRootPolicyItemTypeError); !ok {
			t.Errorf("Expected *invalidRootPolicyItemTypeError but got %T (%s)", err, err)
		}

		err = tr.appendItem(nil, &Policy{hidden: true})
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*hiddenRootPolicyAppendError); !ok {
			t.Errorf("Expected *hiddenRootPolicyAppendError but got %T (%s)", err, err)
		}

		err = tr.appendItem(nil, &Policy{
			id:        "test",
			rules:     []*Rule{{id: "permit", effect: EffectPermit}},
			algorithm: firstApplicableEffectRCA{}})
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		err = tr.appendItem([]string{"example"}, &Rule{id: "permit", effect: EffectPermit})
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*invalidRootPolicyError); !ok {
			t.Errorf("Expected *invalidRootPolicyError but got %T (%s)", err, err)
		}

		err = tr.appendItem([]string{"test"}, &Rule{hidden: true})
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*hiddenRuleAppendError); !ok {
			t.Errorf("Expected *hiddenRuleAppendError but got %T (%s)", err, err)
		}

		err = tr.del(nil)
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*emptyPathModificationError); !ok {
			t.Errorf("Expected *emptyPathModificationError but got %T (%s)", err, err)
		}

		err = tr.del([]string{"example"})
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*invalidRootPolicyError); !ok {
			t.Errorf("Expected *invalidRootPolicyError but got %T (%s)", err, err)
		}

		err = tr.del([]string{"test", "example"})
		if err == nil {
			t.Errorf("Expected error but got nothing")
		} else if _, ok := err.(*missingPolicyChildError); !ok {
			t.Errorf("Expected *missingPolicyChildError but got %T (%s)", err, err)
		}

		err = tr.del([]string{"test", "permit"})
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		err = tr.del([]string{"test"})
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		}

		if tr.policies != nil {
			t.Errorf("Expected no root policy but got %#v", tr.policies)
		}
	}
}

func TestStorageTransactionalUpdate(t *testing.T) {
	tag := uuid.New()

	root := &PolicySet{
		id: "test",
		policies: []Evaluable{
			&Policy{
				id:     "first",
				target: makeSimpleStringTarget("s", "test"),
				rules: []*Rule{
					{
						id:          "permit",
						effect:      EffectPermit,
						obligations: makeSingleStringObligation("s", "permit")}},
				algorithm: denyOverridesRCA{}},
			&Policy{
				id: "del",
				rules: []*Rule{
					{
						id:          "permit",
						effect:      EffectPermit,
						obligations: makeSingleStringObligation("s", "del-permit")}},
				algorithm: firstApplicableEffectRCA{}}},
		algorithm: firstApplicableEffectPCA{}}

	ft, err := NewFlagsType("f", "first")
	if err != nil {
		t.Fatalf("failed to create custom type: %s", err)
	}
	s := NewPolicyStorage(
		root,
		makeSymbols(
			map[string]Type{
				"f": ft,
			},
			map[string]Attribute{
				"s": MakeAttribute("s", TypeString),
			},
		),
		&tag,
	)

	newTag := uuid.New()

	u := NewPolicyUpdate(tag, newTag)
	u.Append(UOAdd, []string{"test", "first"}, &Rule{
		id:          "deny",
		effect:      EffectDeny,
		obligations: makeSingleStringObligation("s", "deny")})
	u.Append(UODelete, []string{"test", "del"}, nil)

	eUpd := fmt.Sprintf("policy update: %s - %s\n"+
		"commands:\n"+
		"- Add path (\"test\"/\"first\")\n"+
		"- Delete path (\"test\"/\"del\")", tag.String(), newTag.String())
	sUpd := u.String()
	if sUpd != eUpd {
		t.Errorf("Expected:\n%s\n\nupdate but got:\n%s\n\n", eUpd, sUpd)
	}

	tr, err := s.NewTransaction(&tag)
	if err != nil {
		t.Fatalf("Expected no error but got %T (%s)", err, err)
	}

	symbols := tr.Symbols()
	if len(symbols.types) != 1 {
		t.Errorf("Expected one custom type but got %#v", symbols.types)
	} else if _, ok := symbols.types["f"]; !ok {
		t.Errorf("Expected %q custom type but got %#v", "f", symbols.types)
	}

	if len(symbols.attrs) != 1 {
		t.Errorf("Expected one attribute but got %#v", symbols.attrs)
	} else if _, ok := symbols.attrs["s"]; !ok {
		t.Errorf("Expected %q attribute but got %#v", "s", symbols.attrs)
	}

	err = tr.Apply(u)
	if err != nil {
		t.Fatalf("Expected no error but got %T (%s)", err, err)
	}

	s, err = tr.Commit()
	if err != nil {
		t.Fatalf("Expected no error but got %T (%s)", err, err)
	}

	ctx, err := NewContext(nil, 1, func(i int) (string, AttributeValue, error) {
		return "s", MakeStringValue("test"), nil
	})
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		r := s.Root().Calculate(ctx)
		if r.Status != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		}

		if r.Effect != EffectDeny {
			t.Errorf("Expected deny effect but got %d", r.Effect)
		}

		if len(r.Obligations) < 1 {
			t.Error("Expected at least one obligation")
		} else {
			_, _, v, err := r.Obligations[0].Serialize(ctx)
			if err != nil {
				t.Errorf("Expected no error but got %T (%s)", err, err)
			} else {
				e := "deny"
				if v != e {
					t.Errorf("Expected %q obligation but got %q", e, v)
				}
			}
		}
	}

	ctx, err = NewContext(nil, 1, func(i int) (string, AttributeValue, error) {
		return "s", MakeStringValue("no test"), nil
	})
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		r := s.Root().Calculate(ctx)
		if r.Status != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		}

		if r.Effect != EffectNotApplicable {
			t.Errorf("Expected \"not applicable\" effect but got %d", r.Effect)
		}
	}
}

func makeSymbols(t map[string]Type, a map[string]Attribute) Symbols {
	return Symbols{
		types: t,
		attrs: a,
	}
}

func TestStorageGetAtPath(t *testing.T) {
	tag := uuid.New()

	r1 := &Rule{
		id:          "permit",
		effect:      EffectPermit,
		obligations: makeSingleStringObligation("s", "permit")}
	r2 := &Rule{
		id:          "permit2",
		effect:      EffectPermit,
		obligations: makeSingleStringObligation("s", "del-permit")}
	p1 := &Policy{
		id:        "first",
		target:    makeSimpleStringTarget("s", "test"),
		rules:     []*Rule{r1},
		algorithm: denyOverridesRCA{}}
	p2 := &Policy{
		id:        "del",
		rules:     []*Rule{r2},
		algorithm: firstApplicableEffectRCA{}}
	root := &PolicySet{
		id:        "test",
		policies:  []Evaluable{p1, p2},
		algorithm: firstApplicableEffectPCA{}}

	ft, err := NewFlagsType("f", "first")
	if err != nil {
		t.Fatalf("failed to create custom type: %s", err)
	}
	s := NewPolicyStorage(
		root,
		makeSymbols(
			map[string]Type{
				"f": ft,
			},
			map[string]Attribute{
				"s": MakeAttribute("s", TypeString),
			},
		),
		&tag,
	)

	expectRoot, expectNil := s.GetAtPath([]string{"test"})
	if expectNil != nil {
		t.Errorf("expected nil error, got %v", expectNil)
	} else if expectRoot != root {
		id, _ := expectRoot.GetID()
		t.Errorf("expected to find root \"test\", got %s", id)
	}

	expectFirst, expectNil := s.GetAtPath([]string{"test", "first"})
	if expectNil != nil {
		t.Errorf("expected nil error, got %v", expectNil)
	} else if expectFirst != p1 {
		id, _ := expectFirst.GetID()
		t.Errorf("expected to find policy \"first\", got %s", id)
	}

	expectDel, expectNil := s.GetAtPath([]string{"test", "del"})
	if expectNil != nil {
		t.Errorf("expected nil error, got %v", expectNil)
	} else if expectDel != p2 {
		id, _ := expectDel.GetID()
		t.Errorf("expected to find policy \"del\", got %s", id)
	}

	expectPermit, expectNil := s.GetAtPath([]string{"test", "first", "permit"})
	if expectNil != nil {
		t.Errorf("expected nil error, got %v", expectNil)
	} else if expectPermit != r1 {
		id, _ := expectPermit.GetID()
		t.Errorf("expected to find rule \"permit\", got %s", id)
	}

	expectPermit2, expectNil := s.GetAtPath([]string{"test", "del", "permit2"})
	if expectNil != nil {
		t.Errorf("expected nil error, got %v", expectNil)
	} else if expectPermit2 != r2 {
		id, _ := expectPermit2.GetID()
		t.Errorf("expected to find rule \"permit2\", got %s", id)
	}

	// error condition 1: path longer than max depth
	badPath := []string{"test", "del", "permit2", "nonexist"}
	expectErr := newPathNotFoundError(badPath)
	expectNilM, expectPathNFE := s.GetAtPath(badPath)
	if !reflect.DeepEqual(expectPathNFE, expectErr) {
		t.Errorf("expected %v, got %v", expectErr, expectPathNFE)
	}
	if expectNilM != nil {
		id, _ := expectNilM.GetID()
		t.Errorf("expected to find nil, got rule/policy with id %s", id)
	}

	// error condition 2: path includes invalid id
	badPath2 := []string{"test", "del2", "permit2"}
	expectErr2 := newPathNotFoundError(badPath2)
	expectNilM, expectPathNFE2 := s.GetAtPath(badPath2)
	if !reflect.DeepEqual(expectPathNFE2, expectErr2) {
		t.Errorf("expected %v, got %v", expectErr2, expectPathNFE2)
	}
	if expectNilM != nil {
		id, _ := expectNilM.GetID()
		t.Errorf("expected to find nil, got rule/policy with id %s", id)
	}
}

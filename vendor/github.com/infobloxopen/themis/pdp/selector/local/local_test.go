package local

import (
	"net/url"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/uintX/domaintree8"
	"github.com/infobloxopen/themis/pdp"
)

func TestLocalSelectorRegistry(t *testing.T) {
	pdp.InitializeSelectors()

	s := pdp.GetSelector(localSelectorScheme)
	if s == nil {
		t.Errorf("Expected selector for %q scheme but got nothing", localSelectorScheme)
	}

	if _, ok := s.(*selector); !ok {
		t.Errorf("Expected local implementation *selector to be registered for %q but got %T (%#v)",
			localSelectorScheme, s, s)
	}
}

func TestMakeLocalSelector(t *testing.T) {
	path := []pdp.Expression{
		pdp.MakeAttributeDesignator(pdp.MakeAttribute("domain", pdp.TypeDomain)),
	}

	uri, err := url.Parse("local:content/item")
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		e, err := pdp.MakeSelector(uri, path, pdp.TypeString)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else if _, ok := e.(LocalSelector); !ok {
			t.Errorf("Expected LocalSelector expression but got %T (%#v)", e, e)
		} else {
			st := e.GetResultType()
			if st != pdp.TypeString {
				t.Errorf("Expected %q as selector result type but got %q", pdp.TypeString, st)
			}
		}
	}

	uri, err = url.Parse("local:content")
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		e, err := pdp.MakeSelector(uri, path, pdp.TypeString)
		if err == nil {
			t.Errorf("Expected error but got selector expression %T (%#v)", e, e)
		}
	}
}

func TestSelectorCalculate(t *testing.T) {
	uri, err := url.Parse("local:test-content/test-item")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	dn := pdp.MakeAttribute("domain", pdp.TypeDomain)
	path := []pdp.Expression{
		pdp.MakeAttributeDesignator(dn),
	}

	sft, err := pdp.NewFlagsType("flags", "first", "second", "third")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	e, err := pdp.MakeSelector(uri, path, sft)
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	cft, err := pdp.NewFlagsType("flags", "first", "second", "third")
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	st := pdp.MakeSymbols()
	if err := st.PutType(cft); err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	dTree8 := &domaintree8.Node{}
	dTree8.InplaceInsert(makeTestDN(t, "example.com"), 1)
	dTree8.InplaceInsert(makeTestDN(t, "example.net"), 3)
	dTree8.InplaceInsert(makeTestDN(t, "example.org"), 5)

	cs := pdp.NewLocalContentStorage([]*pdp.LocalContent{
		pdp.NewLocalContent("test-content", nil, st,
			[]*pdp.ContentItem{
				pdp.MakeContentMappingItem(
					"test-item",
					cft,
					pdp.MakeSignature(pdp.TypeDomain),
					pdp.MakeContentDomainFlags8Map(dTree8),
				),
			},
		),
	})
	ctx, err := pdp.NewContext(cs, 1, func(i int) (string, pdp.AttributeValue, error) {
		v, err := pdp.MakeValueFromString(pdp.TypeDomain, "example.net")
		if err != nil {
			return "", pdp.UndefinedValue, err
		}

		return "domain", v, nil
	})
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	v, err := e.Calculate(ctx)
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"first\",\"second\""
			if s != e {
				t.Errorf("Expected %q value from selector but got %q", e, s)
			}
		}
	}

	e, err = pdp.MakeSelector(uri, path, pdp.TypeString)
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		v, err = e.Calculate(ctx)
		if err == nil {
			s, err := v.Serialize()
			if err != nil {
				t.Errorf("Expected error but got value which can't be serialized: %s", err)
			} else {
				t.Errorf("Expected error but got value %q", s)
			}
		}
	}

	ctx, err = pdp.NewContext(cs, 1, func(i int) (string, pdp.AttributeValue, error) {
		v, err := pdp.MakeValueFromString(pdp.TypeDomain, "example.gov")
		if err != nil {
			return "", pdp.UndefinedValue, err
		}

		return "domain", v, nil
	})
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}

	v, err = e.Calculate(ctx)
	if err == nil {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected error but got value which can't be serialized: %s", err)
		} else {
			t.Errorf("Expected error but got value %q", s)
		}
	} else if _, ok := err.(*pdp.MissingValueError); !ok {
		t.Errorf("Expected *MissingValueError but got %T (%s)", err, err)
	}

	uri, err = url.Parse("local:test-content/missing-item")
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		e, err := pdp.MakeSelector(uri, path, sft)
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			v, err = e.Calculate(ctx)
			if err == nil {
				s, err := v.Serialize()
				if err != nil {
					t.Errorf("Expected error but got value which can't be serialized: %s", err)
				} else {
					t.Errorf("Expected error but got value %q", s)
				}
			} else if _, ok := err.(*pdp.MissingContentItemError); !ok {
				t.Errorf("Expected *MissingContentItemError but got %T (%s)", err, err)
			}

		}
	}
}

func makeTestDN(t *testing.T, s string) domain.Name {
	d, err := domain.MakeNameFromString(s)
	if err != nil {
		t.Fatalf("can't create domain name from string %q: %s", s, err)
	}

	return d
}

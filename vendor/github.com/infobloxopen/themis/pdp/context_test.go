package pdp

import (
	"fmt"
	"net"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

func TestContext(t *testing.T) {
	st := strtree.NewTree()
	st.InplaceInsert("1 - one", 1)
	st.InplaceInsert("2 - two", 2)
	st.InplaceInsert("3 - three", 3)

	nt := iptree.NewTree()
	nt.InplaceInsertNet(makeTestNetwork("192.0.2.0/28"), 1)
	nt.InplaceInsertNet(makeTestNetwork("192.0.2.16/28"), 2)
	nt.InplaceInsertNet(makeTestNetwork("192.0.2.32/28"), 3)

	dt := &domaintree.Node{}
	dt.InplaceInsert(makeTestDomain("example.com"), 1)
	dt.InplaceInsert(makeTestDomain("example.org"), 2)
	dt.InplaceInsert(makeTestDomain("example.net"), 3)

	lt := []string{"1 - one", "2 - two", "3 - three"}

	ct := strtree.NewTree()
	ct.InplaceInsert("test-notag-content", &LocalContent{id: "test-notag-content"})
	var nullContent *LocalContent
	ct.InplaceInsert("test-null-content", nullContent)
	tag := uuid.New()
	ct.InplaceInsert("test-content", &LocalContent{id: "test-content", tag: &tag})
	lcs := &LocalContentStorage{r: ct}

	ctx, err := NewContext(lcs, 9, func(i int) (string, AttributeValue, error) {
		switch i {
		default:
			return "", UndefinedValue, fmt.Errorf("no attribute for index: %d", i)

		case 0:
			return "b", MakeBooleanValue(true), nil

		case 1:
			return "s", MakeStringValue("test"), nil

		case 2:
			return "a", MakeAddressValue(net.ParseIP("192.0.2.1")), nil

		case 3:
			return "n", MakeNetworkValue(makeTestNetwork("192.0.2.0/24")), nil

		case 4:
			return "d", MakeDomainValue(makeTestDomain("example.com")), nil

		case 5:
			return "ss", MakeSetOfStringsValue(st), nil

		case 6:
			return "sn", MakeSetOfNetworksValue(nt), nil

		case 7:
			return "sd", MakeSetOfDomainsValue(dt), nil

		case 8:
			return "ls", MakeListOfStringsValue(lt), nil
		}
	})
	if err != nil {
		t.Error(err)
	} else {
		assertContextAttributes(t, "NewContext", ctx,
			MakeAddressAssignment("a", net.ParseIP("192.0.2.1")),
			MakeBooleanAssignment("b", true),
			MakeDomainAssignment("d", makeTestDomain("example.com")),
			MakeListOfStringsAssignment("ls", lt),
			MakeNetworkAssignment("n", makeTestNetwork("192.0.2.0/24")),
			MakeStringAssignment("s", "test"),
			MakeSetOfDomainsAssignment("sd", dt),
			MakeSetOfNetworksAssignment("sn", nt),
			MakeSetOfStringsAssignment("ss", st),
		)
	}

	ctx, err = NewContextFromBytes(lcs, testWireRequest)
	if err != nil {
		t.Error(err)
	} else {
		assertContextAttributes(t, "NewContextFromBytes", ctx, testRequestAssignments...)
	}
}

func TestResponse(t *testing.T) {
	r := Response{
		Effect:      EffectPermit,
		Obligations: testRequestAssignments,
	}

	var b [258]byte

	n, err := r.MarshalToBuffer(b[:], nil)
	assertRequestBytesBuffer(t, "r.MarshalToBuffer", err, b[:47], n,
		append([]byte{1, 0, 1, 0, 0}, testWireAttributes...)...,
	)

	r = Response{
		Effect: EffectIndeterminate,
		Status: bindError(newMissingAttributeError(), MakeAttribute("t", TypeInteger).describe()),
		Obligations: append(
			testRequestAssignments,
			MakeExpressionAssignment(
				"x",
				makeFunctionIntegerAdd(
					MakeAttributeDesignator(MakeAttribute("y", TypeInteger)),
					MakeAttributeDesignator(MakeAttribute("z", TypeInteger)),
				),
			),
		),
	}

	ctx, _ := NewContext(nil, 0, nil)
	n, err = r.MarshalToBuffer(b[:], ctx)
	assertRequestBytesBuffer(t, "r(multiple errors).MarshalToBuffer", err, b[:258], n,
		append([]byte{
			1, 0, 3,
			211, 0,
			'm', 'u', 'l', 't', 'i', 'p', 'l', 'e', ' ', 'e', 'r', 'r', 'o',
			'r', 's', ':', ' ', '"', '#', 'a', '3', ':', ' ', 'F', 'a', 'i',
			'l', 'e', 'd', ' ', 't', 'o', ' ', 'p', 'r', 'o', 'c', 'e', 's',
			's', ' ', 'r', 'e', 'q', 'u', 'e', 's', 't', ':', ' ', '#', '0',
			'2', ' ', '(', 'a', 't', 't', 'r', '(', 't', '.', 'I', 'n', 't',
			'e', 'g', 'e', 'r', ')', ')', ':', ' ', 'M', 'i', 's', 's', 'i',
			'n', 'g', ' ', 'a', 't', 't', 'r', 'i', 'b', 'u', 't', 'e', '"',
			',', ' ', '"', '#', 'a', '4', ':', ' ', 'F', 'a', 'i', 'l', 'e',
			'd', ' ', 't', 'o', ' ', 'c', 'a', 'l', 'c', 'u', 'l', 'a', 't',
			'e', ' ', 'o', 'b', 'l', 'i', 'g', 'a', 't', 'i', 'o', 'n', ' ',
			'f', 'o', 'r', ' ', 'a', 't', 't', 'r', '(', 'x', '.', 'I', 'n',
			't', 'e', 'g', 'e', 'r', ')', ':', ' ', '#', '0', '2', ' ', '(',
			'a', 'd', 'd', '>', 'f', 'i', 'r', 's', 't', ' ', 'a', 'r', 'g',
			'u', 'm', 'e', 'n', 't', '>', 'a', 't', 't', 'r', '(', 'y', '.',
			'I', 'n', 't', 'e', 'g', 'e', 'r', ')', ')', ':', ' ', 'M', 'i',
			's', 's', 'i', 'n', 'g', ' ', 'a', 't', 't', 'r', 'i', 'b', 'u',
			't', 'e', '"',
		}, testWireAttributes...)...,
	)
}

func assertContextAttributes(t *testing.T, desc string, ctx *Context, attrs ...AttributeAssignment) {
	if ctx == nil {
		t.Errorf("expected some context for %q but got nothing", desc)
		return
	}

	sa := []string{}
	for k, v := range ctx.a {
		switch v := v.(type) {
		default:
			sa = append(sa, fmt.Sprintf("unknown entry %q: %T (%#v)", k, v, v))

		case AttributeValue:
			sa = append(sa, fmt.Sprintf("- %s.(%s): %s", k, v.t, v.describe()))

		case map[Type]AttributeValue:
			for _, v := range v {
				sa = append(sa, fmt.Sprintf("- %s.(%s): %s", k, v.t, v.describe()))
			}
		}
	}
	sort.Strings(sa)

	se := make([]string, len(attrs))
	for i, a := range attrs {
		v, _ := a.e.Calculate(nil)
		se[i] = fmt.Sprintf("- %s.(%s): %s", a.a.id, v.t, v.describe())
	}
	sort.Strings(se)

	assertStrings(sa, se, desc, t)
}

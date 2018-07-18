package pep

import (
	"bytes"
	"math"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

type DummyStruct struct {
}

type TestStruct struct {
	Bool     bool
	Int      int
	Float    float64
	String   string
	If       interface{}
	Address  net.IP
	hidden   string
	Network  *net.IPNet
	Slice    []int
	Struct   DummyStruct
	Domain   domain.Name
	Strings  *strtree.Tree
	Networks *iptree.Tree
	Domains  *domaintree.Node
	StrList  []string
}

type TestTaggedStruct struct {
	bool1    bool
	Bool2    bool             `pdp:""`
	bool3    bool             `pdp:"flag"`
	str      string           `pdp:"s,string"`
	integer  int              `pdp:"i,integer"`
	float    float64          `pdp:"f,float"`
	domain   domain.Name      `pdp:"d,domain"`
	address  net.IP           `pdp:"Address"`
	network  *net.IPNet       `pdp:"net,network"`
	strings  *strtree.Tree    `pdp:"ss,set of strings"`
	networks *iptree.Tree     `pdp:"sn,set of networks"`
	domains  *domaintree.Node `pdp:"sd,set of domains"`
	strlist  []string         `pdp:"ls,list of strings"`
}

type TestInvalidStruct1 struct {
	String string `pdp:",address"`
}

type TestInvalidStruct2 struct {
	If interface{} `pdp:""`
}

var (
	testStruct = TestStruct{
		Bool:    true,
		Int:     5,
		Float:   555.5,
		String:  "test",
		If:      "interface",
		Address: net.ParseIP("1.2.3.4"),
		hidden:  "hide",
		Network: makeTestNetwork("1.2.3.4/32"),
		Slice:   []int{1, 2, 3, 4},
		Struct:  DummyStruct{},
		Domain:  makeTestDomain("example.com"),
		Strings: newStrTree("one", "two", "three"),
		Networks: newIPTree(
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		),
		Domains: newDomainTree(
			makeTestDomain("example.com"),
			makeTestDomain("example.gov"),
			makeTestDomain("www.example.com"),
		),
		StrList: []string{"one", "two", "three"},
	}

	testRequestBuffer = []byte{
		1, 0,
		11, 0,
		4, 'B', 'o', 'o', 'l', 1,
		3, 'I', 'n', 't', 3, 5, 0, 0, 0, 0, 0, 0, 0,
		5, 'F', 'l', 'o', 'a', 't', 4, 0, 0, 0, 0, 0, 92, 129, 64,
		6, 'S', 't', 'r', 'i', 'n', 'g', 2, 4, 0, 't', 'e', 's', 't',
		7, 'A', 'd', 'd', 'r', 'e', 's', 's', 5, 1, 2, 3, 4,
		7, 'N', 'e', 't', 'w', 'o', 'r', 'k', 7, 32, 1, 2, 3, 4,
		6, 'D', 'o', 'm', 'a', 'i', 'n', 9, 11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		7, 'S', 't', 'r', 'i', 'n', 'g', 's', 10, 3, 0,
		3, 0, 'o', 'n', 'e', 3, 0, 't', 'w', 'o', 5, 0, 't', 'h', 'r', 'e', 'e',
		8, 'N', 'e', 't', 'w', 'o', 'r', 'k', 's', 11, 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
		7, 'D', 'o', 'm', 'a', 'i', 'n', 's', 12, 3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		7, 'S', 't', 'r', 'L', 'i', 's', 't', 13, 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	}
)

func TestMarshalUntaggedStruct(t *testing.T) {
	var b [249]byte

	n, err := marshalValue(reflect.ValueOf(testStruct), b[:])
	assertBytesBuffer(t, "marshalValue(TestStruct)", err, b[:], n, testRequestBuffer...)
}

func TestMarshalTaggedStruct(t *testing.T) {
	var b [215]byte

	v := TestTaggedStruct{
		bool1:   true,
		Bool2:   false,
		bool3:   true,
		str:     "test",
		integer: math.MaxInt32,
		float:   12345.6789,
		domain:  makeTestDomain("example.com"),
		address: net.ParseIP("1.2.3.4"),
		network: makeTestNetwork("1.2.3.4/32"),
		strings: newStrTree("one", "two", "three"),
		networks: newIPTree(
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		),
		domains: newDomainTree(
			makeTestDomain("example.com"),
			makeTestDomain("example.gov"),
			makeTestDomain("www.example.com"),
		),
		strlist: []string{"one", "two", "three"},
	}

	n, err := marshalValue(reflect.ValueOf(v), b[:])
	assertBytesBuffer(t, "marshalValue(TestTaggedStruct)", err, b[:], n,
		1, 0,
		12, 0,
		5, 'B', 'o', 'o', 'l', '2', 0,
		4, 'f', 'l', 'a', 'g', 1,
		1, 's', 2, 4, 0, 't', 'e', 's', 't',
		1, 'i', 3, 255, 255, 255, 127, 00, 00, 00, 00,
		1, 'f', 4, 161, 248, 49, 230, 214, 28, 200, 64,
		1, 'd', 9, 11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		7, 'A', 'd', 'd', 'r', 'e', 's', 's', 5, 1, 2, 3, 4,
		3, 'n', 'e', 't', 7, 32, 1, 2, 3, 4,
		2, 's', 's', 10, 3, 0, 3, 0, 'o', 'n', 'e', 3, 0, 't', 'w', 'o', 5, 0, 't', 'h', 'r', 'e', 'e',
		2, 's', 'n', 11, 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
		2, 's', 'd', 12, 3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		2, 'l', 's', 13, 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)
}

func TestMarshalInvalidStructs(t *testing.T) {
	var b [100]byte

	n, err := marshalValue(reflect.ValueOf(TestInvalidStruct1{}), b[:])
	if err == nil {
		t.Errorf("exepcted \"can't marshal\" error but got %d bytes in buffer: [% x]", n, b[:n])
	} else if !strings.Contains(err.Error(), "can't marshal") {
		t.Errorf("Exepcted \"can't marshal\" error but got:\n%s", err)
	}

	n, err = marshalValue(reflect.ValueOf(TestInvalidStruct2{}), b[:])
	if err == nil {
		t.Errorf("exepcted \"can't marshal\" error but got %d bytes in buffer: [% x]", n, b[:n])
	} else if !strings.Contains(err.Error(), "can't marshal") {
		t.Errorf("Exepcted \"can't marshal\" error but got:\n%s", err)
	}
}

func TestMakeRequest(t *testing.T) {
	var b [249]byte

	m, err := makeRequest(pb.Msg{Body: testRequestBuffer}, b[:])
	assertBytesBuffer(t, "makeRequest(pb.Msg)", err, m.Body, len(m.Body), testRequestBuffer...)

	m, err = makeRequest(&pb.Msg{Body: testRequestBuffer}, b[:])
	assertBytesBuffer(t, "makeRequest(&pb.Msg)", err, m.Body, len(m.Body), testRequestBuffer...)

	m, err = makeRequest(testRequestBuffer, b[:])
	assertBytesBuffer(t, "makeRequest(testRequestBuffer)", err, m.Body, len(m.Body), testRequestBuffer...)

	m, err = makeRequest(testStruct, b[:])
	assertBytesBuffer(t, "makeRequest(testStruct)", err, m.Body, len(m.Body), testRequestBuffer...)

	m, err = makeRequest([]pdp.AttributeAssignment{
		pdp.MakeBooleanAssignment("Bool", true),
		pdp.MakeIntegerAssignment("Int", 5),
		pdp.MakeFloatAssignment("Float", 555.5),
		pdp.MakeStringAssignment("String", "test"),
		pdp.MakeAddressAssignment("Address", net.ParseIP("1.2.3.4")),
		pdp.MakeNetworkAssignment("Network", makeTestNetwork("1.2.3.4/32")),
		pdp.MakeDomainAssignment("Domain", makeTestDomain("example.com")),
		pdp.MakeSetOfStringsAssignment("Strings", newStrTree("one", "two", "three")),
		pdp.MakeSetOfNetworksAssignment("Networks", newIPTree(
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		)),
		pdp.MakeSetOfDomainsAssignment("Domains", newDomainTree(
			makeTestDomain("example.com"),
			makeTestDomain("example.gov"),
			makeTestDomain("www.example.com"),
		)),
		pdp.MakeListOfStringsAssignment("StrList", []string{"one", "two", "three"}),
	}, b[:])
	assertBytesBuffer(t, "makeRequest(assignments)", err, m.Body, len(m.Body), testRequestBuffer...)
}

func makeTestNetwork(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}

	return n
}

func makeTestDomain(s string) domain.Name {
	d, err := domain.MakeNameFromString(s)
	if err != nil {
		panic(err)
	}

	return d
}

func newStrTree(args ...string) *strtree.Tree {
	t := strtree.NewTree()
	for i, s := range args {
		t.InplaceInsert(s, i)
	}

	return t
}

func newIPTree(args ...*net.IPNet) *iptree.Tree {
	t := iptree.NewTree()
	for i, n := range args {
		t.InplaceInsertNet(n, i)
	}

	return t
}

func newDomainTree(args ...domain.Name) *domaintree.Node {
	t := new(domaintree.Node)
	for i, d := range args {
		t.InplaceInsert(d, i)
	}

	return t
}

func assertBytesBuffer(t *testing.T, desc string, err error, b []byte, n int, e ...byte) {
	if err != nil {
		t.Errorf("expected no error for %s but got: %s", desc, err)
	} else if n != len(b) {
		t.Errorf("expected exactly all buffer used (%d bytes) for %s but got %d bytes", len(b), desc, n)
	} else {
		if bytes.Compare(b[:], e) != 0 {
			t.Errorf("expected [% x] for %s but got [% x]", e, desc, b)
		}
	}
}

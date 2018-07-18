package pep

import (
	"fmt"
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

type TestResponseStruct struct {
	Effect  bool
	Int     int
	Float   float64
	Bool    bool
	String  string
	Address net.IP
	Network *net.IPNet
	Strings *strtree.Tree
}

type TestTaggedResponseStruct struct {
	Result  string `pdp:"Effect"`
	Error   string `pdp:"Reason"`
	Bool1   bool
	Bool2   bool       `pdp:""`
	Bool3   bool       `pdp:"flag"`
	Int     int        `pdp:"i,integer"`
	Float   float64    `pdp:"f,float"`
	Domain  string     `pdp:"d,domain"`
	Address net.IP     `pdp:""`
	Network *net.IPNet `pdp:"net,network"`
}

type TestTaggedAllTypesResponseStruct struct {
	Effect     string           `pdp:"Effect"`
	Reason     string           `pdp:"Reason"`
	BoolFalse  bool             `pdp:"baf"`
	BoolTrue   bool             `pdp:"bat"`
	String     string           `pdp:"sa"`
	Int        int              `pdp:"ia"`
	Int8       int8             `pdp:"i8a"`
	Int16      int16            `pdp:"i16a"`
	Int32      int32            `pdp:"i32a"`
	Int64      int64            `pdp:"i64a"`
	Uint       uint             `pdp:"uia"`
	Uint8      uint8            `pdp:"ui8a"`
	Uint16     uint16           `pdp:"ui16a"`
	Uint32     uint32           `pdp:"ui32a"`
	Uint64     uint64           `pdp:"ui64a"`
	Float32    float32          `pdp:"f32a"`
	Float64    float64          `pdp:"f64a"`
	Address4   net.IP           `pdp:"aa4"`
	Address6   net.IP           `pdp:"aa6"`
	Network4   net.IPNet        `pdp:"na4"`
	Network6   net.IPNet        `pdp:"na6"`
	NetworkPtr *net.IPNet       `pdp:"pna"`
	DomainS    string           `pdp:"das,domain"`
	DomainD    domain.Name      `pdp:"dad,domain"`
	Strings    *strtree.Tree    `pdp:"ssa,set of strings"`
	Networks   *iptree.Tree     `pdp:"sna,set of networks"`
	Domains    *domaintree.Node `pdp:"sda,set of domains"`
	StrList    []string         `pdp:"lsa,list of strings"`
}

type TestInvalidResponseStruct1 struct {
	Effect bool `pdp:"Effect,string"`
}

type TestInvalidResponseStruct2 struct {
	Reason string `pdp:"Reason,string"`
}

type TestInvalidResponseStruct3 struct {
	Attribute string `pdp:",unknown"`
}

type TestInvalidResponseStruct4 struct {
	Attribute bool `pdp:",address"`
}

type TestInvalidResponseStruct5 struct {
	flag bool `pdp:"flag"`
}

var (
	TestResponse = []byte{
		1, 0, 1,
		0, 0,
		6, 0,
		4, 'B', 'o', 'o', 'l', 1,
		6, 'S', 't', 'r', 'i', 'n', 'g', 2, 4, 0, 't', 'e', 's', 't',
		3, 'I', 'n', 't', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		5, 'F', 'l', 'o', 'a', 't', 4, 233, 72, 46, 63, 164, 84, 33, 65,
		7, 'A', 'd', 'd', 'r', 'e', 's', 's', 5, 1, 2, 3, 4,
		7, 'N', 'e', 't', 'w', 'o', 'r', 'k', 7, 32, 1, 2, 3, 4,
	}

	TestTaggedResponse = []byte{
		1, 0, 3,
		11, 0, 'T', 'e', 's', 't', ' ', 'E', 'r', 'r', 'o', 'r', '!',
		5, 0,
		5, 'B', 'o', 'o', 'l', '2', 0,
		4, 'f', 'l', 'a', 'g', 1,
		1, 'd', 9, 11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		7, 'A', 'd', 'd', 'r', 'e', 's', 's', 5, 1, 2, 3, 4,
		3, 'n', 'e', 't', 7, 32, 1, 2, 3, 4,
	}

	TestTaggedAllTypesResponse = []byte{
		1, 0, 3,
		11, 0, 'T', 'e', 's', 't', ' ', 'E', 'r', 'r', 'o', 'r', '!',
		26, 0,
		3, 'b', 'a', 'f', 0,
		3, 'b', 'a', 't', 1,
		2, 's', 'a', 2, 4, 0, 't', 'e', 's', 't',
		2, 'i', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		3, 'i', '8', 'a', 3, 64, 0, 0, 0, 0, 0, 0, 0,
		4, 'i', '1', '6', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		4, 'i', '3', '2', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		4, 'i', '6', '4', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		3, 'u', 'i', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		4, 'u', 'i', '8', 'a', 3, 64, 0, 0, 0, 0, 0, 0, 0,
		5, 'u', 'i', '1', '6', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		5, 'u', 'i', '3', '2', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		5, 'u', 'i', '6', '4', 'a', 3, 210, 4, 0, 0, 0, 0, 0, 0,
		4, 'f', '3', '2', 'a', 4, 190, 193, 23, 38, 3, 133, 186, 64,
		4, 'f', '6', '4', 'a', 4, 190, 193, 23, 38, 3, 133, 186, 64,
		3, 'a', 'a', '4', 5, 192, 0, 2, 1,
		3, 'a', 'a', '6', 6, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		3, 'n', 'a', '4', 7, 24, 192, 0, 2, 0,
		3, 'n', 'a', '6', 8, 32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 'p', 'n', 'a', 8, 32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		3, 'd', 'a', 's', 9, 11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		3, 'd', 'a', 'd', 9, 11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		3, 's', 's', 'a', 10, 3, 0, 3, 0, 'o', 'n', 'e', 3, 0, 't', 'w', 'o', 5, 0, 't', 'h', 'r', 'e', 'e',
		3, 's', 'n', 'a', 11, 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
		3, 's', 'd', 'a', 12, 3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		3, 'l', 's', 'a', 13, 3, 0, 3, 0, 'o', 'n', 'e', 3, 0, 't', 'w', 'o', 5, 0, 't', 'h', 'r', 'e', 'e',
	}
)

func TestUnmarshalUntaggedStruct(t *testing.T) {
	v := TestResponseStruct{}

	err := unmarshalToValue(TestResponse, reflect.ValueOf(&v))
	if err != nil {
		t.Error(err)
	} else {
		assertTestResponseStruct(t, v,
			TestResponseStruct{
				Effect:  true,
				Int:     1234,
				Float:   567890.1234,
				Bool:    true,
				String:  "test",
				Address: net.ParseIP("1.2.3.4"),
				Network: makeTestNetwork("1.2.3.4/32"),
				Strings: newStrTree("one", "two", "three"),
			},
		)
	}
}

func TestUnmarshalTaggedStruct(t *testing.T) {
	v := TestTaggedResponseStruct{}

	err := unmarshalToValue(TestTaggedResponse, reflect.ValueOf(&v))
	if err != nil {
		t.Error(err)
	} else {
		assertTestTaggedStruct(t, v,
			TestTaggedResponseStruct{
				Result:  pdp.EffectNameFromEnum(pdp.EffectIndeterminate),
				Error:   "Test Error!",
				Bool1:   false,
				Bool2:   false,
				Bool3:   true,
				Int:     0,
				Float:   0.,
				Domain:  "example.com",
				Address: net.ParseIP("1.2.3.4"),
				Network: makeTestNetwork("1.2.3.4/32"),
			},
		)
	}

	vAllTypes := TestTaggedAllTypesResponseStruct{}
	err = unmarshalToValue(TestTaggedAllTypesResponse, reflect.ValueOf(&vAllTypes))
	if err != nil {
		t.Error(err)
	} else {
		assertTestTaggedAllTypesStruct(t, vAllTypes,
			TestTaggedAllTypesResponseStruct{
				Effect:     pdp.EffectNameFromEnum(pdp.EffectIndeterminate),
				Reason:     "Test Error!",
				BoolFalse:  false,
				BoolTrue:   true,
				String:     "test",
				Int:        1234,
				Int8:       64,
				Int16:      1234,
				Int32:      1234,
				Int64:      1234,
				Uint:       1234,
				Uint8:      64,
				Uint16:     1234,
				Uint32:     1234,
				Uint64:     1234,
				Float32:    6789.012,
				Float64:    6789.0123,
				Address4:   net.ParseIP("192.0.2.1"),
				Address6:   net.ParseIP("2001:db8::1"),
				Network4:   *makeTestNetwork("192.0.2.1/24"),
				Network6:   *makeTestNetwork("2001:db8::/32"),
				NetworkPtr: makeTestNetwork("2001:db8::/32"),
				DomainS:    "example.com",
				DomainD:    makeTestDomain("example.com"),
				Strings:    newStrTree("one", "two", "three"),
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
			},
		)
	}
}

type TestIntResponse struct {
	Effect int8
}

type TestUintResponse struct {
	Effect uint8
}

func TestUnmarshalEffectTypes(t *testing.T) {
	r := []byte{1, 0, 5, 0, 0, 0, 0}

	v1 := TestIntResponse{}
	err := unmarshalToValue(r, reflect.ValueOf(&v1))
	if err != nil {
		t.Error(err)
	} else if v1.Effect != int8(pdp.EffectIndeterminateP) {
		t.Errorf("expected %d %q effect but got %d %q",
			pdp.EffectIndeterminateP, pdp.EffectNameFromEnum(pdp.EffectIndeterminateP),
			v1.Effect, pdp.EffectNameFromEnum(int(v1.Effect)),
		)
	}

	v2 := TestUintResponse{}
	err = unmarshalToValue(r, reflect.ValueOf(&v2))
	if err != nil {
		t.Error(err)
	} else if v2.Effect != uint8(pdp.EffectIndeterminateP) {
		t.Errorf("expected %d %q effect but got %d %q",
			pdp.EffectIndeterminateP, pdp.EffectNameFromEnum(pdp.EffectIndeterminateP),
			v2.Effect, pdp.EffectNameFromEnum(int(v1.Effect)),
		)
	}
}

type TestErrorResponse struct {
	Reason error
}

func TestUnmarshalReasonErrorType(t *testing.T) {
	r := []byte{
		1, 0, 3,
		11, 0, 'T', 'e', 's', 't', ' ', 'E', 'r', 'r', 'o', 'r', '!',
		0, 0,
	}

	v := TestErrorResponse{}
	err := unmarshalToValue(r, reflect.ValueOf(&v))
	if err != nil {
		t.Error(err)
	} else {
		if v.Reason == nil {
			t.Error("expected \"Test Error!\" but got nothing")
		} else if !strings.Contains(v.Reason.Error(), "Test Error!") {
			t.Errorf("expected \"Test Error!\" but got %q", v.Reason.Error())
		}
	}
}

func TestUnmarshalInvalidStructures(t *testing.T) {
	r := []byte{1, 0, 5, 0, 0, 0, 0}
	v1 := TestInvalidResponseStruct1{}
	err := unmarshalToValue(r, reflect.ValueOf(&v1))
	if err != nil {
		if !strings.Contains(err.Error(), "don't support type definition") {
			t.Errorf("expected \"don't support type definition\" error but got: %s", err)
		}
	} else {
		t.Errorf("expected \"don't support type definition\" error")
	}

	v2 := TestInvalidResponseStruct2{}
	err = unmarshalToValue(r, reflect.ValueOf(&v2))
	if err != nil {
		if !strings.Contains(err.Error(), "don't support type definition") {
			t.Errorf("expected \"don't support type definition\" error but got: %s", err)
		}
	} else {
		t.Errorf("expected \"don't support type definition\" error")
	}

	v3 := TestInvalidResponseStruct3{}
	err = unmarshalToValue(r, reflect.ValueOf(&v3))
	if err != nil {
		if !strings.Contains(err.Error(), "unknown type") {
			t.Errorf("expected \"unknown type\" error but got: %s", err)
		}
	} else {
		t.Errorf("expected \"unknown type\" error")
	}

	v4 := TestInvalidResponseStruct4{}
	err = unmarshalToValue(r, reflect.ValueOf(&v4))
	if err != nil {
		if !strings.Contains(err.Error(), "tagged type") {
			t.Errorf("expected \"tagged type\" error but got: %s", err)
		}
	} else {
		t.Errorf("expected \"tagged type\" error")
	}

	v5 := TestInvalidResponseStruct5{}
	err = unmarshalToValue(TestTaggedResponse, reflect.ValueOf(&v5))
	if err != nil {
		if !strings.Contains(err.Error(), "can't be set") {
			t.Errorf("expected \"can't be set\" error but got: %s", err)
		}
	} else {
		t.Errorf("expected \"can't be set\" error")
	}
}

func TestFillResponse(t *testing.T) {
	r := pb.Msg{
		Body: TestResponse,
	}

	m := new(pb.Msg)
	err := fillResponse(r, m)
	if err != nil {
		t.Error(err)
	} else if string(m.Body) != string(r.Body) {
		t.Errorf("expected same body %p as original response but got %p", r.Body, m.Body)
	}

	o := make([]pdp.AttributeAssignment, 10)
	pr := &pdp.Response{
		Obligations: o,
	}
	err = fillResponse(r, pr)
	if err != nil {
		t.Error(err)
	} else if len(pr.Obligations) != 6 {
		t.Errorf("expected %d obligations but got %d", 6, len(pr.Obligations))
	}

	v := TestResponseStruct{}
	err = fillResponse(r, &v)
	if err != nil {
		t.Error(err)
	} else {
		assertTestResponseStruct(t, v,
			TestResponseStruct{
				Effect:  true,
				Int:     1234,
				Float:   567890.1234,
				Bool:    true,
				String:  "test",
				Address: net.ParseIP("1.2.3.4"),
				Network: makeTestNetwork("1.2.3.4/32"),
			},
		)
	}
}

func assertTestResponseStruct(t *testing.T, v, e TestResponseStruct) {
	if v.Effect != e.Effect ||
		v.Int != e.Int ||
		v.Float != e.Float ||
		v.Bool != e.Bool ||
		v.String != e.String ||
		v.Address.String() != e.Address.String() ||
		v.Network.String() != e.Network.String() {
		t.Errorf("expected:\n%v\nbut got:\n%v\n", SprintfTestResponseStruct(e), SprintfTestResponseStruct(v))
	}
}

func assertTestTaggedStruct(t *testing.T, v, e TestTaggedResponseStruct) {
	if v.Result != e.Result ||
		v.Error != e.Error ||
		v.Bool1 != e.Bool1 ||
		v.Bool2 != e.Bool2 ||
		v.Bool3 != e.Bool3 ||
		v.Int != e.Int ||
		v.Float != e.Float ||
		v.Domain != e.Domain ||
		v.Address.String() != e.Address.String() ||
		v.Network.String() != e.Network.String() {
		t.Errorf("expected:\n%v\nbut got:\n%v\n", SprintfTestTaggedStruct(e), SprintfTestTaggedStruct(v))
	}
}

func assertTestTaggedAllTypesStruct(t *testing.T, v, e TestTaggedAllTypesResponseStruct) {
	vN := "<nil>"
	if v.NetworkPtr != nil {
		vN = v.NetworkPtr.String()
	}

	eN := "<nil>"
	if e.NetworkPtr != nil {
		eN = e.NetworkPtr.String()
	}

	if v.Effect != e.Effect ||
		v.Reason != e.Reason ||
		v.BoolFalse != v.BoolFalse ||
		v.BoolTrue != e.BoolTrue ||
		v.String != e.String ||
		v.Int != e.Int ||
		v.Int8 != e.Int8 ||
		v.Int16 != e.Int16 ||
		v.Int32 != e.Int32 ||
		v.Int64 != e.Int64 ||
		v.Uint != e.Uint ||
		v.Uint8 != e.Uint8 ||
		v.Uint16 != e.Uint16 ||
		v.Uint32 != e.Uint32 ||
		v.Uint64 != e.Uint64 ||
		v.Float32 != e.Float32 ||
		v.Float64 != e.Float64 ||
		v.Address4.String() != e.Address4.String() ||
		v.Address6.String() != e.Address6.String() ||
		v.Network4.String() != e.Network4.String() ||
		v.Network6.String() != e.Network6.String() ||
		vN != eN ||
		v.DomainS != e.DomainS ||
		v.DomainD.String() != e.DomainD.String() ||
		!setOfStringsEqual(v.Strings, e.Strings) ||
		!setOfNetworksEqual(v.Networks, e.Networks) ||
		!setOfDomainsEqual(v.Domains, e.Domains) ||
		!listOfStringsEqual(v.StrList, e.StrList) {
		t.Errorf("expected:\n%v\nbut got:\n%v\n",
			SprintfTestTaggedAllTypesStruct(e),
			SprintfTestTaggedAllTypesStruct(v),
		)
	}
}

func SprintfTestResponseStruct(v TestResponseStruct) string {
	return fmt.Sprintf(
		"\tEffect.: %v\n"+
			"\tInt....: %v\n"+
			"\tFloat..: %v\n"+
			"\tBool...: %v\n"+
			"\tString.: %v\n"+
			"\tAddress: %s\n"+
			"\tNetwork: %s\n"+
			"\tStrings: %v\n",
		v.Effect,
		v.Int, v.Float, v.Bool, v.String, v.Address.String(), v.Network.String(), pdp.SortSetOfStrings(v.Strings))
}

func SprintfTestTaggedStruct(v TestTaggedResponseStruct) string {
	return fmt.Sprintf(
		"\tResult.: %v\n"+
			"\tError..: %v\n"+
			"\tBool1..: %v\n"+
			"\tBool2..: %v\n"+
			"\tBool3..: %v\n"+
			"\tDomain.: %v\n"+
			"\tAddress: %v\n"+
			"\tNetwork: %v\n",
		v.Result, v.Error,
		v.Bool1, v.Bool2, v.Bool3, v.Domain, v.Address.String(), v.Network.String())
}

func SprintfTestTaggedAllTypesStruct(v TestTaggedAllTypesResponseStruct) string {
	n := "<nil>"
	if v.NetworkPtr != nil {
		n = v.NetworkPtr.String()
	}

	return fmt.Sprintf(
		"\tEffect....: %v\n"+
			"\tReason....: %v\n"+
			"\tBoolFalse.: %v\n"+
			"\tBoolTrue..: %v\n"+
			"\tString....: %v\n"+
			"\tInt.......: %v\n"+
			"\tInt8......: %v\n"+
			"\tInt16.....: %v\n"+
			"\tInt32.....: %v\n"+
			"\tInt64.....: %v\n"+
			"\tUint......: %v\n"+
			"\tUint8.....: %v\n"+
			"\tUint16....: %v\n"+
			"\tUint32....: %v\n"+
			"\tUint64....: %v\n"+
			"\tFloat32...: %v\n"+
			"\tFloat64...: %v\n"+
			"\tAddress4..: %v\n"+
			"\tAddress6..: %v\n"+
			"\tNetwork4..: %v\n"+
			"\tNetwork6..: %v\n"+
			"\tNetworkPtr: %v\n"+
			"\tDomainS...: %q\n"+
			"\tDomainD...: %q\n"+
			"\tStrings...: %v\n"+
			"\tNetworks..: %v\n"+
			"\tDomains...: %v\n"+
			"\tStrList...: %v\n",
		v.Effect, v.Reason, v.BoolFalse, v.BoolTrue, v.String,
		v.Int, v.Int8, v.Int16, v.Int32, v.Int64,
		v.Uint, v.Uint8, v.Uint16, v.Uint32, v.Uint64,
		v.Float32, v.Float64,
		v.Address4.String(), v.Address6.String(), v.Network4.String(), v.Network6.String(), n,
		v.DomainS, v.DomainD,
		pdp.SortSetOfStrings(v.Strings),
		makeStringFromNetworks(v.Networks),
		pdp.SortSetOfDomains(v.Domains),
		v.StrList,
	)
}

func setOfStringsEqual(v, e *strtree.Tree) bool {
	ss := pdp.SortSetOfStrings(v)
	se := pdp.SortSetOfStrings(e)

	if len(ss) != len(se) {
		return false
	}

	for i, s := range ss {
		if s != se[i] {
			return false
		}
	}

	return true
}

func makeStringFromNetworks(t *iptree.Tree) []string {
	sn := pdp.SortSetOfNetworks(t)

	out := make([]string, len(sn))
	for i, n := range sn {
		out[i] = n.String()
	}

	return out
}

func setOfNetworksEqual(v, e *iptree.Tree) bool {
	sn := makeStringFromNetworks(v)
	se := makeStringFromNetworks(e)

	if len(sn) != len(se) {
		return false
	}

	for i, n := range sn {
		if n != se[i] {
			return false
		}
	}

	return true
}

func setOfDomainsEqual(v, e *domaintree.Node) bool {
	sd := pdp.SortSetOfDomains(v)
	se := pdp.SortSetOfDomains(e)

	if len(sd) != len(se) {
		return false
	}

	for i, d := range sd {
		if d != se[i] {
			return false
		}
	}

	return true
}

func listOfStringsEqual(v, e []string) bool {
	if len(v) != len(e) {
		return false
	}

	for i, s := range v {
		if s != e[i] {
			return false
		}
	}

	return true
}

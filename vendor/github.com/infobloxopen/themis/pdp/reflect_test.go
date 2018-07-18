package pdp

import (
	"math"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

func TestSetEffect(t *testing.T) {
	var b bool
	if err := setEffect(reflect.Indirect(reflect.ValueOf(&b)), EffectPermit); err != nil {
		t.Error(err)
	} else if !b {
		t.Errorf("expected %v for %q but got %v", true, EffectNameFromEnum(EffectPermit), b)
	}

	var i int
	if err := setEffect(reflect.Indirect(reflect.ValueOf(&i)), EffectPermit); err != nil {
		t.Error(err)
	} else if i != EffectPermit {
		t.Errorf("expected %d for %q but got %d", EffectPermit, EffectNameFromEnum(EffectPermit), i)
	}

	var u uint
	if err := setEffect(reflect.Indirect(reflect.ValueOf(&u)), EffectPermit); err != nil {
		t.Error(err)
	} else if u != EffectPermit {
		t.Errorf("expected %d for %q but got %d", EffectPermit, EffectNameFromEnum(EffectPermit), u)
	}

	var s string
	if err := setEffect(reflect.Indirect(reflect.ValueOf(&s)), EffectPermit); err != nil {
		t.Error(err)
	} else if s != EffectNameFromEnum(EffectPermit) {
		t.Errorf("expected %q but got %q", EffectNameFromEnum(EffectPermit), s)
	}

	if err := setEffect(reflect.ValueOf(nil), EffectPermit); err != nil {
		t.Error(err)
	}

	err := setEffect(reflect.ValueOf(i), EffectNotApplicable)
	if err == nil {
		t.Errorf("expected *requestUnmarshalEffectConstError but got %d (%q)", i, EffectNameFromEnum(i))
	} else if _, ok := err.(*requestUnmarshalEffectConstError); !ok {
		t.Errorf("expected *requestUnmarshalEffectConstError but got %T (%s)", err, err)
	}

	var f float64
	err = setEffect(reflect.Indirect(reflect.ValueOf(&f)), EffectPermit)
	if err == nil {
		t.Errorf("expected *requestUnmarshalEffectTypeError but got %g", f)
	} else if _, ok := err.(*requestUnmarshalEffectTypeError); !ok {
		t.Errorf("expected *requestUnmarshalEffectTypeError but got %T (%s)", err, err)
	}
}

func TestSetStatus(t *testing.T) {
	var s string
	if err := setStatus(reflect.Indirect(reflect.ValueOf(&s)), "testError"); err != nil {
		t.Error(err)
	} else if s != "testError" {
		t.Errorf("expected %q but got %q", "testError", s)
	}

	var sErr error
	if err := setStatus(reflect.Indirect(reflect.ValueOf(&sErr)), "testError"); err != nil {
		t.Error(err)
	} else if !strings.Contains(sErr.Error(), "testError") {
		t.Errorf("expected %q in %q", "testError", sErr)
	}

	if err := setStatus(reflect.ValueOf(nil), "testError"); err != nil {
		t.Error(err)
	}

	err := setStatus(reflect.ValueOf(s), "testError")
	if err == nil {
		t.Errorf("expected *requestUnmarshalStatusConstError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalStatusConstError); !ok {
		t.Errorf("expected *requestUnmarshalStatusConstError but got %T (%s)", err, err)
	}

	var f float64
	err = setStatus(reflect.Indirect(reflect.ValueOf(&f)), "testError")
	if err == nil {
		t.Errorf("expected *requestUnmarshalStatusTypeError but got %g", f)
	} else if _, ok := err.(*requestUnmarshalStatusTypeError); !ok {
		t.Errorf("expected *requestUnmarshalStatusTypeError but got %T (%s)", err, err)
	}
}

func TestSetBool(t *testing.T) {
	var b bool
	if err := setBool(reflect.Indirect(reflect.ValueOf(&b)), true); err != nil {
		t.Error(err)
	} else if !b {
		t.Errorf("expected %v but got %v", true, b)
	}

	if err := setBool(reflect.ValueOf(nil), true); err != nil {
		t.Error(err)
	}

	err := setBool(reflect.ValueOf(b), true)
	if err == nil {
		t.Errorf("expected *requestUnmarshalBooleanConstError but got %v", b)
	} else if _, ok := err.(*requestUnmarshalBooleanConstError); !ok {
		t.Errorf("expected *requestUnmarshalBooleanConstError but got %T (%s)", err, err)
	}

	var s string
	err = setBool(reflect.Indirect(reflect.ValueOf(&s)), true)
	if err == nil {
		t.Errorf("expected *requestUnmarshalBooleanTypeError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalBooleanTypeError); !ok {
		t.Errorf("expected *requestUnmarshalBooleanTypeError but got %T (%s)", err, err)
	}
}

func TestSetString(t *testing.T) {
	var s string
	if err := setString(reflect.Indirect(reflect.ValueOf(&s)), "test"); err != nil {
		t.Error(err)
	} else if s != "test" {
		t.Errorf("expected %q but got %q", "test", s)
	}

	if err := setString(reflect.ValueOf(nil), "test"); err != nil {
		t.Error(err)
	}

	err := setString(reflect.ValueOf(s), "test")
	if err == nil {
		t.Errorf("expected *requestUnmarshalStringConstError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalStringConstError); !ok {
		t.Errorf("expected *requestUnmarshalStringConstError but got %T (%s)", err, err)
	}

	var b bool
	err = setString(reflect.Indirect(reflect.ValueOf(&b)), "test")
	if err == nil {
		t.Errorf("expected *requestUnmarshalStringTypeError but got %v", b)
	} else if _, ok := err.(*requestUnmarshalStringTypeError); !ok {
		t.Errorf("expected *requestUnmarshalStringTypeError but got %T (%s)", err, err)
	}
}

func TestSetInt(t *testing.T) {
	var i int
	if err := setInt(reflect.Indirect(reflect.ValueOf(&i)), math.MinInt32); err != nil {
		t.Error(err)
	} else if i != math.MinInt32 {
		t.Errorf("expected %d value but got %d", math.MinInt32, i)
	}

	var i8 int8
	if err := setInt(reflect.Indirect(reflect.ValueOf(&i8)), math.MinInt8); err != nil {
		t.Error(err)
	} else if i8 != math.MinInt8 {
		t.Errorf("expected %d value but got %d", math.MinInt8, i8)
	}

	var i16 int16
	if err := setInt(reflect.Indirect(reflect.ValueOf(&i16)), math.MinInt16); err != nil {
		t.Error(err)
	} else if i16 != math.MinInt16 {
		t.Errorf("expected %d value but got %d", math.MinInt16, i16)
	}

	var i32 int32
	if err := setInt(reflect.Indirect(reflect.ValueOf(&i32)), math.MinInt32); err != nil {
		t.Error(err)
	} else if i32 != math.MinInt32 {
		t.Errorf("expected %d value but got %d", math.MinInt32, i32)
	}

	var i64 int64
	if err := setInt(reflect.Indirect(reflect.ValueOf(&i64)), math.MinInt64); err != nil {
		t.Error(err)
	} else if i64 != math.MinInt64 {
		t.Errorf("expected %d value but got %d", math.MinInt64, i64)
	}

	var u uint
	if err := setInt(reflect.Indirect(reflect.ValueOf(&u)), math.MaxInt32); err != nil {
		t.Error(err)
	} else if u != math.MaxInt32 {
		t.Errorf("expected %d value but got %d", math.MaxInt32, u)
	}

	var u8 uint8
	if err := setInt(reflect.Indirect(reflect.ValueOf(&u8)), math.MaxInt8); err != nil {
		t.Error(err)
	} else if u8 != math.MaxInt8 {
		t.Errorf("expected %d value but got %d", math.MaxInt8, u8)
	}

	var u16 uint16
	if err := setInt(reflect.Indirect(reflect.ValueOf(&u16)), math.MaxInt16); err != nil {
		t.Error(err)
	} else if u16 != math.MaxInt16 {
		t.Errorf("expected %d value but got %d", math.MaxInt16, u16)
	}

	var u32 uint32
	if err := setInt(reflect.Indirect(reflect.ValueOf(&u32)), math.MaxInt32); err != nil {
		t.Error(err)
	} else if u32 != math.MaxInt32 {
		t.Errorf("expected %d value but got %d", math.MaxInt32, u32)
	}

	var u64 uint64
	if err := setInt(reflect.Indirect(reflect.ValueOf(&u64)), math.MaxInt64); err != nil {
		t.Error(err)
	} else if u64 != math.MaxInt64 {
		t.Errorf("expected %d value but got %d", math.MaxInt64, u64)
	}

	if err := setInt(reflect.ValueOf(nil), 0); err != nil {
		t.Error(err)
	}

	err := setInt(reflect.ValueOf(i), 0)
	if err == nil {
		t.Errorf("expected *requestUnmarshalIntegerConstError but got %d", i)
	} else if _, ok := err.(*requestUnmarshalIntegerConstError); !ok {
		t.Errorf("expected *requestUnmarshalIntegerConstError but got %T (%s)", err, err)
	}

	var s string
	err = setInt(reflect.Indirect(reflect.ValueOf(&s)), 0)
	if err == nil {
		t.Errorf("expected *requestUnmarshalIntegerTypeError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalIntegerTypeError); !ok {
		t.Errorf("expected *requestUnmarshalIntegerTypeError but got %T (%s)", err, err)
	}
}

func TestSetIntOverflow(t *testing.T) {
	var i int
	v := reflect.Indirect(reflect.ValueOf(&i))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var i8 int8
	v = reflect.Indirect(reflect.ValueOf(&i8))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var i16 int16
	v = reflect.Indirect(reflect.ValueOf(&i16))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var i32 int32
	v = reflect.Indirect(reflect.ValueOf(&i32))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var u uint
	v = reflect.Indirect(reflect.ValueOf(&u))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var u8 uint8
	v = reflect.Indirect(reflect.ValueOf(&u8))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var u16 uint16
	v = reflect.Indirect(reflect.ValueOf(&u16))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)

	var u32 uint32
	v = reflect.Indirect(reflect.ValueOf(&u32))

	assertUnmarshalIntegerOverflowError(t, setInt(v, math.MaxInt64), v)
}

func TestSetIntUnderflow(t *testing.T) {
	var u uint
	v := reflect.Indirect(reflect.ValueOf(&u))

	assertUnmarshalIntegerUnderflowError(t, setInt(v, -1), v)

	var u8 uint8
	v = reflect.Indirect(reflect.ValueOf(&u8))

	assertUnmarshalIntegerUnderflowError(t, setInt(v, -1), v)

	var u16 uint16
	v = reflect.Indirect(reflect.ValueOf(&u16))

	assertUnmarshalIntegerUnderflowError(t, setInt(v, -1), v)

	var u32 uint32
	v = reflect.Indirect(reflect.ValueOf(&u32))

	assertUnmarshalIntegerUnderflowError(t, setInt(v, -1), v)

	var u64 uint64
	v = reflect.Indirect(reflect.ValueOf(&u64))

	assertUnmarshalIntegerUnderflowError(t, setInt(v, -1), v)
}

func TestSetFloat(t *testing.T) {
	var f32 float32
	v := reflect.Indirect(reflect.ValueOf(&f32))
	if err := setFloat(v, math.MaxFloat32); err != nil {
		t.Error(err)
	} else if f32 != math.MaxFloat32 {
		t.Errorf("expected %g but got %g", math.MaxFloat32, f32)
	}

	var f64 float64
	if err := setFloat(reflect.Indirect(reflect.ValueOf(&f64)), math.MaxFloat64); err != nil {
		t.Error(err)
	} else if f64 != math.MaxFloat64 {
		t.Errorf("expected %g but got %g", math.MaxFloat64, f64)
	}

	if err := setFloat(v, math.MaxFloat64); err != nil {
		t.Error(err)
	} else if !math.IsInf(float64(f32), 1) {
		t.Errorf("expected positive infinity but got %g", f32)
	}

	if err := setFloat(v, math.SmallestNonzeroFloat64); err != nil {
		t.Error(err)
	} else if f32 != 0 {
		t.Errorf("expected zero but got %g", f32)
	}

	if err := setFloat(reflect.ValueOf(nil), 0); err != nil {
		t.Error(err)
	}

	err := setFloat(reflect.ValueOf(f64), 0)
	if err == nil {
		t.Errorf("expected *requestUnmarshalFloatConstError but got %g", f64)
	} else if _, ok := err.(*requestUnmarshalFloatConstError); !ok {
		t.Errorf("expected *requestUnmarshalFloatConstError but got %T (%s)", err, err)
	}

	var s string
	err = setFloat(reflect.Indirect(reflect.ValueOf(&s)), 0)
	if err == nil {
		t.Errorf("expected *requestUnmarshalFloatTypeError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalFloatTypeError); !ok {
		t.Errorf("expected *requestUnmarshalFloatTypeError but got %T (%s)", err, err)
	}
}

func TestSetAddress(t *testing.T) {
	ea := net.ParseIP("192.0.2.1")
	var a net.IP
	if err := setAddress(reflect.Indirect(reflect.ValueOf(&a)), ea); err != nil {
		t.Error(err)
	} else if !a.Equal(ea) {
		t.Errorf("expected %q but got %q", ea, a)
	}

	if err := setAddress(reflect.ValueOf(nil), ea); err != nil {
		t.Error(err)
	}

	err := setAddress(reflect.ValueOf(a), ea)
	if err == nil {
		t.Errorf("expected *requestUnmarshalAddressConstError but got %q", a)
	} else if _, ok := err.(*requestUnmarshalAddressConstError); !ok {
		t.Errorf("expected *requestUnmarshalAddressConstError but got %T (%s)", err, err)
	}

	var s string
	err = setAddress(reflect.Indirect(reflect.ValueOf(&s)), ea)
	if err == nil {
		t.Errorf("expected *requestUnmarshalAddressTypeError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalAddressTypeError); !ok {
		t.Errorf("expected *requestUnmarshalAddressTypeError but got %T (%s)", err, err)
	}
}

func TestGetNetwork(t *testing.T) {
	en := makeTestNetwork("192.0.2.0/24")

	n := getNetwork(reflect.ValueOf(en))
	if n.String() != en.String() {
		t.Errorf("expected %q but got %q", en, n)
	}

	n = getNetwork(reflect.ValueOf(*en))
	if n.String() != en.String() {
		t.Errorf("expected %q but got %q", en, n)
	}

	n = getNetwork(reflect.ValueOf(nil))
	if n != nil {
		t.Errorf("expected nil but got %q", n)
	}

	assertPanicWithError(t, func() {
		n = getNetwork(reflect.ValueOf("string"))
	}, "can't marshal %s as network value", reflect.TypeOf("string"))
}

func TestSetNetwork(t *testing.T) {
	en := makeTestNetwork("192.0.2.0/24")
	var nPtr *net.IPNet
	if err := setNetwork(reflect.Indirect(reflect.ValueOf(&nPtr)), en); err != nil {
		t.Error(err)
	} else if nPtr.String() != en.String() {
		t.Errorf("expected %q but got %q", en, nPtr)
	}

	var n net.IPNet
	if err := setNetwork(reflect.Indirect(reflect.ValueOf(&n)), en); err != nil {
		t.Error(err)
	} else if n.String() != en.String() {
		t.Errorf("expected %q but got %q", en, n)
	}

	if err := setNetwork(reflect.ValueOf(nil), en); err != nil {
		t.Error(err)
	}

	err := setNetwork(reflect.ValueOf(n), en)
	if err == nil {
		t.Errorf("expected *requestUnmarshalNetworkConstError but got %q", n)
	} else if _, ok := err.(*requestUnmarshalNetworkConstError); !ok {
		t.Errorf("expected *requestUnmarshalNetworkConstError but got %T (%s)", err, err)
	}

	var s string
	err = setNetwork(reflect.Indirect(reflect.ValueOf(&s)), en)
	if err == nil {
		t.Errorf("expected *requestUnmarshalNetworkTypeError but got %q", s)
	} else if _, ok := err.(*requestUnmarshalNetworkTypeError); !ok {
		t.Errorf("expected *requestUnmarshalNetworkTypeError but got %T (%s)", err, err)
	}
}

func TestSetDomain(t *testing.T) {
	eDn := makeTestDomain("example.com")
	var dn domain.Name
	if err := setDomain(reflect.Indirect(reflect.ValueOf(&dn)), eDn); err != nil {
		t.Error(err)
	} else if dn.String() != eDn.String() {
		t.Errorf("expected %q but got %q", eDn, dn)
	}

	var s string
	if err := setDomain(reflect.Indirect(reflect.ValueOf(&s)), eDn); err != nil {
		t.Error(err)
	} else if s != eDn.String() {
		t.Errorf("expected %q but got %q", eDn, s)
	}

	if err := setDomain(reflect.ValueOf(nil), eDn); err != nil {
		t.Error(err)
	}

	err := setDomain(reflect.ValueOf(dn), eDn)
	if err == nil {
		t.Errorf("expected *requestUnmarshalDomainConstError but got %q", dn)
	} else if _, ok := err.(*requestUnmarshalDomainConstError); !ok {
		t.Errorf("expected *requestUnmarshalDomainConstError but got %T (%s)", err, err)
	}

	var n int64
	err = setDomain(reflect.Indirect(reflect.ValueOf(&n)), eDn)
	if err == nil {
		t.Errorf("expected *requestUnmarshalDomainTypeError but got %d", n)
	} else if _, ok := err.(*requestUnmarshalDomainTypeError); !ok {
		t.Errorf("expected *requestUnmarshalDomainTypeError but got %T (%s)", err, err)
	}
}

func TestGetSetOfStrings(t *testing.T) {
	eSs := newStrTree("one", "two", "three")

	ss := getSetOfStrings(reflect.ValueOf(eSs))
	assertStrings(SortSetOfStrings(ss), SortSetOfStrings(eSs), "getSetOfStrings", t)

	ss = getSetOfStrings(reflect.ValueOf(nil))
	if ss != nil {
		t.Errorf("expected nil but got %#v", SortSetOfStrings(ss))
	}

	assertPanicWithError(t, func() {
		ss = getSetOfStrings(reflect.ValueOf(true))
	}, "can't marshal %s as set of strings value", reflect.TypeOf(true))
}

func TestSetSetOfStrings(t *testing.T) {
	eSs := newStrTree("one", "two", "three")
	var ss *strtree.Tree
	if err := setSetOfStrings(reflect.Indirect(reflect.ValueOf(&ss)), eSs); err != nil {
		t.Error(err)
	} else {
		assertStrings(SortSetOfStrings(ss), SortSetOfStrings(eSs), "setSetOfStrings", t)
	}

	if err := setSetOfStrings(reflect.ValueOf(nil), eSs); err != nil {
		t.Error(err)
	}

	err := setSetOfStrings(reflect.ValueOf(ss), eSs)
	if err == nil {
		t.Errorf("expected *requestUnmarshalSetOfStringsConstError but got %#v", SortSetOfStrings(ss))
	} else if _, ok := err.(*requestUnmarshalSetOfStringsConstError); !ok {
		t.Errorf("expected *requestUnmarshalSetOfStringsConstError but got %T (%s)", err, err)
	}

	var n int64
	err = setSetOfStrings(reflect.Indirect(reflect.ValueOf(&n)), eSs)
	if err == nil {
		t.Errorf("expected *requestUnmarshalSetOfStringsTypeError but got %d", n)
	} else if _, ok := err.(*requestUnmarshalSetOfStringsTypeError); !ok {
		t.Errorf("expected *requestUnmarshalSetOfStringsTypeError but got %T (%s)", err, err)
	}
}

func TestGetSetOfNetworks(t *testing.T) {
	eSn := newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	)

	sn := getSetOfNetworks(reflect.ValueOf(eSn))
	assertNetworks(SortSetOfNetworks(sn), SortSetOfNetworks(eSn), "getSetOfNetworks", t)

	sn = getSetOfNetworks(reflect.ValueOf(nil))
	if sn != nil {
		t.Errorf("expected nil but got %#v", SortSetOfNetworks(sn))
	}

	assertPanicWithError(t, func() {
		sn = getSetOfNetworks(reflect.ValueOf(true))
	}, "can't marshal %s as set of networks value", reflect.TypeOf(true))
}

func TestSetSetOfNetworks(t *testing.T) {
	eSn := newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	)

	var sn *iptree.Tree
	if err := setSetOfNetworks(reflect.Indirect(reflect.ValueOf(&sn)), eSn); err != nil {
		t.Error(err)
	} else {
		assertNetworks(SortSetOfNetworks(sn), SortSetOfNetworks(eSn), "setSetOfNetworks", t)
	}

	if err := setSetOfNetworks(reflect.ValueOf(nil), eSn); err != nil {
		t.Error(err)
	}

	err := setSetOfNetworks(reflect.ValueOf(sn), eSn)
	if err == nil {
		t.Errorf("expected *requestUnmarshalSetOfNetworksConstError but got %#v", SortSetOfNetworks(sn))
	} else if _, ok := err.(*requestUnmarshalSetOfNetworksConstError); !ok {
		t.Errorf("expected *requestUnmarshalSetOfNetworksConstError but got %T (%s)", err, err)
	}

	var n int64
	err = setSetOfNetworks(reflect.Indirect(reflect.ValueOf(&n)), eSn)
	if err == nil {
		t.Errorf("expected *requestUnmarshalSetOfNetworksTypeError but got %d", n)
	} else if _, ok := err.(*requestUnmarshalSetOfNetworksTypeError); !ok {
		t.Errorf("expected *requestUnmarshalSetOfNetworksTypeError but got %T (%s)", err, err)
	}
}

func TestGetSetOfDomains(t *testing.T) {
	eSd := newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	)

	sd := getSetOfDomains(reflect.ValueOf(eSd))
	assertStrings(SortSetOfDomains(sd), SortSetOfDomains(eSd), "getSetOfDomains", t)

	sd = getSetOfDomains(reflect.ValueOf(nil))
	if sd != nil {
		t.Errorf("expected nil but got %#v", SortSetOfDomains(sd))
	}

	assertPanicWithError(t, func() {
		sd = getSetOfDomains(reflect.ValueOf(true))
	}, "can't marshal %s as set of domains value", reflect.TypeOf(true))
}

func TestSetSetOfDomains(t *testing.T) {
	eSd := newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	)

	var sd *domaintree.Node
	if err := setSetOfDomains(reflect.Indirect(reflect.ValueOf(&sd)), eSd); err != nil {
		t.Error(err)
	} else {
		assertStrings(SortSetOfDomains(sd), SortSetOfDomains(eSd), "setSetOfDomains", t)
	}

	if err := setSetOfDomains(reflect.ValueOf(nil), eSd); err != nil {
		t.Error(err)
	}

	err := setSetOfDomains(reflect.ValueOf(sd), eSd)
	if err == nil {
		t.Errorf("expected *requestUnmarshalSetOfDomainsConstError but got %#v", SortSetOfDomains(sd))
	} else if _, ok := err.(*requestUnmarshalSetOfDomainsConstError); !ok {
		t.Errorf("expected *requestUnmarshalSetOfDomainsConstError but got %T (%s)", err, err)
	}

	var n int64
	err = setSetOfDomains(reflect.Indirect(reflect.ValueOf(&n)), eSd)
	if err == nil {
		t.Errorf("expected *requestUnmarshalSetOfDomainsTypeError but got %d", n)
	} else if _, ok := err.(*requestUnmarshalSetOfDomainsTypeError); !ok {
		t.Errorf("expected *requestUnmarshalSetOfDomainsTypeError but got %T (%s)", err, err)
	}
}

func TestGetListOfStrings(t *testing.T) {
	eLs := []string{"one", "two", "three"}

	ls := getListOfStrings(reflect.ValueOf(eLs))
	assertStrings(ls, eLs, "getListOfStrings", t)

	ls = getListOfStrings(reflect.ValueOf(nil))
	if ls != nil {
		t.Errorf("expected nil but got %#v", ls)
	}

	assertPanicWithError(t, func() {
		ls = getListOfStrings(reflect.ValueOf(true))
	}, "can't marshal %s as list of strings value", reflect.TypeOf(true))
}

func TestSetListOfStrings(t *testing.T) {
	eLs := []string{"one", "two", "three"}
	var ls []string
	if err := setListOfStrings(reflect.Indirect(reflect.ValueOf(&ls)), eLs); err != nil {
		t.Error(err)
	} else {
		assertStrings(ls, eLs, "setListOfStrings", t)
	}

	if err := setListOfStrings(reflect.ValueOf(nil), eLs); err != nil {
		t.Error(err)
	}

	err := setListOfStrings(reflect.ValueOf(ls), eLs)
	if err == nil {
		t.Errorf("expected *requestUnmarshalListOfStringsConstError but got %#v", ls)
	} else if _, ok := err.(*requestUnmarshalListOfStringsConstError); !ok {
		t.Errorf("expected *requestUnmarshalListOfStringsConstError but got %T (%s)", err, err)
	}

	var n int64
	err = setListOfStrings(reflect.Indirect(reflect.ValueOf(&n)), eLs)
	if err == nil {
		t.Errorf("expected *requestUnmarshalListOfStringsTypeError but got %d", n)
	} else if _, ok := err.(*requestUnmarshalListOfStringsTypeError); !ok {
		t.Errorf("expected *requestUnmarshalListOfStringsTypeError but got %T (%s)", err, err)
	}
}

func assertUnmarshalIntegerOverflowError(t *testing.T, err error, v reflect.Value) {
	if err == nil {
		t.Errorf("expected *requestUnmarshalIntegerOverflowError but got value %#v", v)
	} else if _, ok := err.(*requestUnmarshalIntegerOverflowError); !ok {
		t.Errorf("expected *requestUnmarshalIntegerOverflowError but got %T (%s)", err, err)
	}
}

func assertUnmarshalIntegerUnderflowError(t *testing.T, err error, v reflect.Value) {
	if err == nil {
		t.Errorf("expected *requestUnmarshalIntegerUnderflowError but got value %#v", v)
	} else if _, ok := err.(*requestUnmarshalIntegerUnderflowError); !ok {
		t.Errorf("expected *requestUnmarshalIntegerUnderflowError but got %T (%s)", err, err)
	}
}

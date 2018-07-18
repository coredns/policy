package pdp

import (
	"fmt"
	"net"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

func TestAttributeValue(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	tt := Type(&builtinType{
		n: "Test-Type",
		k: "test-type",
	})
	v := AttributeValue{t: tt, v: nil}
	expDesc := "val(unknown type)"
	d := v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	r, err := v.Calculate(ctx)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if r.t != v.t || r.v != v.v {
		t.Errorf("Expected the same attribute with type %d and value %T (%#v) but got %d and %T (%#v)",
			v.t, v.v, v.v, r.t, r.v, r.v)
	}

	v = UndefinedValue
	vt := v.GetResultType()
	if vt != TypeUndefined {
		t.Errorf("Expected %q as value type but got %q", TypeUndefined, vt)
	}

	expDesc = "val(undefined)"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeBooleanValue(true)
	vt = v.GetResultType()
	if vt != TypeBoolean {
		t.Errorf("Expected %q as value type but got %q", TypeBoolean, vt)
	}

	expDesc = "true"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeStringValue("test")
	vt = v.GetResultType()
	if vt != TypeString {
		t.Errorf("Expected %q as value type but got %q", TypeString, vt)
	}

	expDesc = "\"test\""
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeIntegerValue(123)
	vt = v.GetResultType()
	if vt != TypeInteger {
		t.Errorf("Expected %q as value type but got %q", TypeInteger, vt)
	}

	expDesc = "123"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeFloatValue(123.456)
	vt = v.GetResultType()
	if vt != TypeFloat {
		t.Errorf("Expected %q as value type but got %q", TypeFloat, vt)
	}

	expDesc = "123.456"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeAddressValue(net.ParseIP("192.0.2.1"))
	vt = v.GetResultType()
	if vt != TypeAddress {
		t.Errorf("Expected %q as value type but got %q", TypeAddress, vt)
	}

	expDesc = "192.0.2.1"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeNetworkValue(makeTestNetwork("192.0.2.0/24"))
	vt = v.GetResultType()
	if vt != TypeNetwork {
		t.Errorf("Expected %q as value type but got %q", TypeNetwork, vt)
	}

	expDesc = "192.0.2.0/24"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeDomainValue(makeTestDomain("example.com"))
	vt = v.GetResultType()
	if vt != TypeDomain {
		t.Errorf("Expected %q as value type but got %q", TypeDomain, vt)
	}

	expDesc = "domain(example.com)"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	sTree := strtree.NewTree()
	sTree.InplaceInsert("1 - one", 1)
	sTree.InplaceInsert("2 - two", 2)
	sTree.InplaceInsert("3 - three", 3)
	sTree.InplaceInsert("4 - four", 4)
	v = MakeSetOfStringsValue(sTree)
	vt = v.GetResultType()
	if vt != TypeSetOfStrings {
		t.Errorf("Expected %q as value type but got %q", TypeSetOfStrings, vt)
	}

	expDesc = "set(\"1 - one\", \"2 - two\", ...)"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	nTree := iptree.NewTree()
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.16/28"), 1)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.32/28"), 2)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.48/28"), 3)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.64/28"), 4)
	v = MakeSetOfNetworksValue(nTree)
	vt = v.GetResultType()
	if vt != TypeSetOfNetworks {
		t.Errorf("Expected %q as value type but got %q", TypeSetOfNetworks, vt)
	}

	expDesc = "set(192.0.2.16/28, 192.0.2.32/28, ...)"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	dTree := &domaintree.Node{}
	dTree.InplaceInsert(makeTestDomain("example.com"), 1)
	dTree.InplaceInsert(makeTestDomain("example.gov"), 2)
	dTree.InplaceInsert(makeTestDomain("example.net"), 3)
	dTree.InplaceInsert(makeTestDomain("example.org"), 4)
	v = MakeSetOfDomainsValue(dTree)
	vt = v.GetResultType()
	if vt != TypeSetOfDomains {
		t.Errorf("Expected %q as value type but got %q", TypeSetOfDomains, vt)
	}

	expDesc = "domains(\"example.com\", \"example.gov\", ...)"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	v = MakeListOfStringsValue([]string{"one", "two", "three", "four"})
	vt = v.GetResultType()
	if vt != TypeListOfStrings {
		t.Errorf("Expected %q as value type but got %q", TypeListOfStrings, vt)
	}

	expDesc = "[\"one\", \"two\", ...]"
	d = v.describe()
	if d != expDesc {
		t.Errorf("Expected %q as value description but got %q", expDesc, d)
	}

	ft8, err := NewFlagsType("8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue8(3, ft8)
		expDesc = "flags<\"8flags\">(\"f00\", \"f01\")"
		d = v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}

		assertPanicWithError(t, func() {
			MakeFlagsValue8(3, TypeBoolean)
		}, "can't make flags value for type %q", TypeBoolean)

		assertPanicWithError(t, func() {
			MakeFlagsValue16(3, ft8)
		}, "expected 8 bits value for \"8flags\" but got 16")
	}

	ft16, err := NewFlagsType("16flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue16(7, ft16)
		expDesc = "flags<\"16flags\">(\"f00\", \"f01\", ...)"
		d = v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}

		assertPanicWithError(t, func() {
			MakeFlagsValue16(3, TypeBoolean)
		}, "can't make flags value for type %q", TypeBoolean)

		assertPanicWithError(t, func() {
			MakeFlagsValue32(3, ft16)
		}, "expected 16 bits value for \"16flags\" but got 32")
	}

	ft32, err := NewFlagsType("32flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue32(3, ft32)
		expDesc = "flags<\"32flags\">(\"f00\", \"f01\")"
		d = v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}

		assertPanicWithError(t, func() {
			MakeFlagsValue32(3, TypeBoolean)
		}, "can't make flags value for type %q", TypeBoolean)

		assertPanicWithError(t, func() {
			MakeFlagsValue64(3, ft32)
		}, "expected 32 bits value for \"32flags\" but got 64")
	}

	ft64, err := NewFlagsType("64flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue64(7, ft64)
		expDesc = "flags<\"64flags\">(\"f00\", \"f01\", ...)"
		d = v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}

		assertPanicWithError(t, func() {
			MakeFlagsValue64(3, TypeBoolean)
		}, "can't make flags value for type %q", TypeBoolean)

		assertPanicWithError(t, func() {
			MakeFlagsValue8(3, ft64)
		}, "expected 64 bits value for \"64flags\" but got 8")
	}
}

func TestMakeValueFromSting(t *testing.T) {
	tt := Type(&builtinType{
		n: "Test-Type",
		k: "test-type",
	})
	v, err := MakeValueFromString(tt, "test")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*unknownTypeStringCastError); !ok {
		t.Errorf("Expected *unknownTypeStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeUndefined, "test")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidTypeStringCastError); !ok {
		t.Errorf("Expected *invalidTypeStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeSetOfStrings, "test")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*notImplementedStringCastError); !ok {
		t.Errorf("Expected *notImplementedStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeSetOfNetworks, "test")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*notImplementedStringCastError); !ok {
		t.Errorf("Expected *notImplementedStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeSetOfDomains, "test")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*notImplementedStringCastError); !ok {
		t.Errorf("Expected *notImplementedStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeListOfStrings, "test")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*notImplementedStringCastError); !ok {
		t.Errorf("Expected *notImplementedStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeBoolean, "true")
	if err != nil {
		t.Errorf("Expected boolean attribute value but got error: %s", err)
	} else {
		expDesc := "true"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeBoolean, "not boolean value")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidBooleanStringCastError); !ok {
		t.Errorf("Expected *invalidBooleanStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeString, "test")
	if err != nil {
		t.Errorf("Expected string attribute value but got error: %s", err)
	} else {
		expDesc := "\"test\""
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeInteger, "654321")
	if err != nil {
		t.Errorf("Expected integer attribute value but got error: %s", err)
	} else {
		expDesc := "654321"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %s as value description but got %s", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeInteger, "not integer value")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidIntegerStringCastError); !ok {
		t.Errorf("Expected *invalidIntegerStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeFloat, "654.321")
	if err != nil {
		t.Errorf("Expected integer attribute value but got error: %s", err)
	} else {
		expDesc := "654.321"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %s as value description but got %s", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeFloat, "0.00000000000654321")
	if err != nil {
		t.Errorf("Expected integer attribute value but got error: %s", err)
	} else {
		expDesc := "6.54321E-12"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %s as value description but got %s", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeFloat, "not float value")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidFloatStringCastError); !ok {
		t.Errorf("Expected *invalidFloatStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeAddress, "192.0.2.1")
	if err != nil {
		t.Errorf("Expected address attribute value but got error: %s", err)
	} else {
		expDesc := "192.0.2.1"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeAddress, "999.999.999.999")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidAddressStringCastError); !ok {
		t.Errorf("Expected *invalidAddressStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeNetwork, "192.0.2.0/24")
	if err != nil {
		t.Errorf("Expected network attribute value but got error: %s", err)
	} else {
		expDesc := "192.0.2.0/24"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeNetwork, "999.999.999.999/999")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidNetworkStringCastError); !ok {
		t.Errorf("Expected *invalidNetworkStringCastError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeDomain, "example.com")
	if err != nil {
		t.Errorf("Expected domain attribute value but got error: %s", err)
	} else {
		expDesc := "domain(example.com)"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	v, err = MakeValueFromString(TypeDomain, "..")
	if err == nil {
		t.Errorf("Expected error but got value: %s", v.describe())
	} else if _, ok := err.(*invalidDomainNameStringCastError); !ok {
		t.Errorf("Expected *invalidDomainNameStringCastError but got %T (%s)", err, err)
	}

	ft, err := NewFlagsType("flags", "first", "second", "third")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v, err = MakeValueFromString(ft, "")
		if err == nil {
			t.Errorf("Expected error but got value: %s", v.describe())
		} else if _, ok := err.(*notImplementedStringCastError); !ok {
			t.Errorf("Expected *notImplementedStringCastError but got %T (%s)", err, err)
		}
	}
}

func TestAttributeValueTypeCast(t *testing.T) {
	v, err := MakeValueFromString(TypeBoolean, "true")
	if err != nil {
		t.Errorf("Expected boolean attribute value but got error: %s", err)
	} else {
		b, err := v.boolean()
		if err != nil {
			t.Errorf("Expected boolean value but got error: %s", err)
		} else if !b {
			t.Errorf("Expected true as attribute value but got %#v", b)
		}

		s, err := v.str()
		if err == nil {
			t.Errorf("Expected error but got string %q", s)
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	v, err = MakeValueFromString(TypeString, "test")
	if err != nil {
		t.Errorf("Expected string attribute value but got error: %s", err)
	} else {
		s, err := v.str()
		if err != nil {
			t.Errorf("Expected string value but got error: %s", err)
		} else if s != "test" {
			t.Errorf("Expected \"test\" as attribute value but got %q", s)
		}

		i, err := v.integer()
		if err == nil {
			t.Errorf("Expected error but got integer %d", i)
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	v, err = MakeValueFromString(TypeInteger, "34567")
	if err != nil {
		t.Errorf("Expected integer attribute value but got error: %s", err)
	} else {
		i, err := v.integer()
		if err != nil {
			t.Errorf("Expected integer value but got error: %s", err)
		} else if i != 34567 {
			t.Errorf("Expected %d as attribute value but got %#v", 34567, i)
		}

		f, err := v.float()
		if err == nil {
			t.Errorf("Expected error but got float %g", f)
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	v, err = MakeValueFromString(TypeFloat, "12345.678")
	if err != nil {
		t.Errorf("Expected float attribute value but got error: %s", err)
	} else {
		f, err := v.float()
		if err != nil {
			t.Errorf("Expected float value but got error: %s", err)
		} else if f != 12345.678 {
			t.Errorf("Expected %g as attribute value but got %#v", 12345.678, f)
		}

		a, err := v.address()
		if err == nil {
			t.Errorf("Expected error but got address %s", a)
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	v, err = MakeValueFromString(TypeAddress, "192.0.2.1")
	if err != nil {
		t.Errorf("Expected address attribute value but got error: %s", err)
	} else {
		_, err := v.address()
		if err != nil {
			t.Errorf("Expected address value but got error: %s", err)
		} else {
			expDesc := "192.0.2.1"
			d := v.describe()
			if d != expDesc {
				t.Errorf("Expected %q as value description but got %q", expDesc, d)
			}
		}

		n, err := v.network()
		if err == nil {
			t.Errorf("Expected error but got network %s", n)
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	v, err = MakeValueFromString(TypeNetwork, "192.0.2.0/24")
	if err != nil {
		t.Errorf("Expected network attribute value but got error: %s", err)
	} else {
		_, err := v.network()
		if err != nil {
			t.Errorf("Expected network value but got error: %s", err)
		} else {
			expDesc := "192.0.2.0/24"
			d := v.describe()
			if d != expDesc {
				t.Errorf("Expected %q as value description but got %q", expDesc, d)
			}
		}

		d, err := v.domain()
		if err == nil {
			t.Errorf("Expected error but got domain %s", d)
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	v, err = MakeValueFromString(TypeDomain, "example.com")
	if err != nil {
		t.Errorf("Expected domain attribute value but got error: %s", err)
	} else {
		_, err := v.domain()
		if err != nil {
			t.Errorf("Expected domain value but got error: %s", err)
		} else {
			expDesc := "domain(example.com)"
			d := v.describe()
			if d != expDesc {
				t.Errorf("Expected %q as attribute value but got %s", expDesc, d)
			}
		}

		_, err = v.setOfStrings()
		if err == nil {
			t.Errorf("Expected error but got set of strings %s", v.describe())
		} else if _, ok := err.(*attributeValueTypeError); !ok {
			t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
		}
	}

	sTree := strtree.NewTree()
	sTree.InplaceInsert("1 - one", 1)
	sTree.InplaceInsert("2 - two", 2)
	sTree.InplaceInsert("3 - three", 3)
	sTree.InplaceInsert("4 - four", 4)
	v = MakeSetOfStringsValue(sTree)

	_, err = v.setOfStrings()
	if err != nil {
		t.Errorf("Expected set of strings value but got error: %s", err)
	} else {
		expDesc := "set(\"1 - one\", \"2 - two\", ...)"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	_, err = v.setOfNetworks()
	if err == nil {
		t.Errorf("Expected error but got set of networks %s", v.describe())
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	nTree := iptree.NewTree()
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.16/28"), 1)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.32/28"), 2)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.48/28"), 3)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.48/28"), 4)
	v = MakeSetOfNetworksValue(nTree)

	_, err = v.setOfNetworks()
	if err != nil {
		t.Errorf("Expected set of networks value but got error: %s", err)
	} else {
		expDesc := "set(192.0.2.16/28, 192.0.2.32/28, ...)"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	_, err = v.setOfDomains()
	if err == nil {
		t.Errorf("Expected error but got set of domains %s", v.describe())
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	dTree := &domaintree.Node{}
	dTree.InplaceInsert(makeTestDomain("example.com"), 1)
	dTree.InplaceInsert(makeTestDomain("example.gov"), 2)
	dTree.InplaceInsert(makeTestDomain("example.net"), 3)
	dTree.InplaceInsert(makeTestDomain("example.org"), 4)
	v = MakeSetOfDomainsValue(dTree)
	_, err = v.setOfDomains()
	if err != nil {
		t.Errorf("Expected set of domains value but got error: %s", err)
	} else {
		expDesc := "domains(\"example.com\", \"example.gov\", ...)"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	_, err = v.listOfStrings()
	if err == nil {
		t.Errorf("Expected error but got list of strings %s", v.describe())
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	v = MakeListOfStringsValue([]string{"one", "two", "three", "four"})
	_, err = v.listOfStrings()
	if err != nil {
		t.Errorf("Expected list of strings value but got error: %s", err)
	} else {
		expDesc := "[\"one\", \"two\", ...]"
		d := v.describe()
		if d != expDesc {
			t.Errorf("Expected %q as value description but got %q", expDesc, d)
		}
	}

	b, err := v.boolean()
	if err == nil {
		t.Errorf("Expected error but got boolean %#v", b)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	ft8, err := NewFlagsType("8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue8(3, ft8)
		n8, err := v.flags8()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n8 != 3 {
			t.Errorf("Expected 3 as attribute value but got %d", n8)
		}

		n864, err := v.flags()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n864 != uint64(n8) {
			t.Errorf("Expected %d as attribute value but got %d", n8, n864)
		}

		n16, err := v.flags16()
		if err == nil {
			t.Errorf("Expected error but got 16 bits flags %d from 8 bits flags %s", n16, v.describe())
		} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
			t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
		}

		v = MakeBooleanValue(true)
		n8, err = v.flags8()
		if err == nil {
			t.Errorf("Expected error but got flags %d from boolean %s", n8, v.describe())
		} else if _, ok := err.(*attributeValueFlagsTypeError); !ok {
			t.Errorf("Expected *attributeValueFlagsTypeError but got %T (%s)", err, err)
		}

		n64, err := v.flags()
		if err == nil {
			t.Errorf("Expected error but got flags %d from boolean %s", n64, v.describe())
		} else if _, ok := err.(*attributeValueFlagsTypeError); !ok {
			t.Errorf("Expected *attributeValueFlagsTypeError but got %T (%s)", err, err)
		}
	}

	ft16, err := NewFlagsType("16flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue16(3, ft16)
		n16, err := v.flags16()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n16 != 3 {
			t.Errorf("Expected 3 as attribute value but got %d", n16)
		}

		n1664, err := v.flags()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n1664 != uint64(n16) {
			t.Errorf("Expected %d as attribute value but got %d", n16, n1664)
		}

		n32, err := v.flags32()
		if err == nil {
			t.Errorf("Expected error but got 32 bits flags %d from 16 bits flags %s", n32, v.describe())
		} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
			t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
		}

		v = MakeBooleanValue(true)
		n16, err = v.flags16()
		if err == nil {
			t.Errorf("Expected error but got flags %d from boolean %s", n16, v.describe())
		} else if _, ok := err.(*attributeValueFlagsTypeError); !ok {
			t.Errorf("Expected *attributeValueFlagsTypeError but got %T (%s)", err, err)
		}
	}

	ft32, err := NewFlagsType("32flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue32(3, ft32)
		n32, err := v.flags32()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n32 != 3 {
			t.Errorf("Expected 3 as attribute value but got %d", n32)
		}

		n3264, err := v.flags()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n3264 != uint64(n32) {
			t.Errorf("Expected %d as attribute value but got %d", n32, n3264)
		}

		n64, err := v.flags64()
		if err == nil {
			t.Errorf("Expected error but got 64 bits flags %d from 32 bits flags %s", n64, v.describe())
		} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
			t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
		}

		v = MakeBooleanValue(true)
		n32, err = v.flags32()
		if err == nil {
			t.Errorf("Expected error but got flags %d from boolean %s", n32, v.describe())
		} else if _, ok := err.(*attributeValueFlagsTypeError); !ok {
			t.Errorf("Expected *attributeValueFlagsTypeError but got %T (%s)", err, err)
		}
	}

	ft64, err := NewFlagsType("64flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue64(3, ft64)
		n64, err := v.flags64()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n64 != 3 {
			t.Errorf("Expected 3 as attribute value but got %d", n64)
		}

		n6464, err := v.flags()
		if err != nil {
			t.Errorf("Expected flags value but got error: %s", err)
		} else if n6464 != n64 {
			t.Errorf("Expected %d as attribute value but got %d", n64, n6464)
		}

		n8, err := v.flags8()
		if err == nil {
			t.Errorf("Expected error but got 8 bits flags %d from 64 bits flags %s", n8, v.describe())
		} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
			t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
		}

		v = MakeBooleanValue(true)
		n64, err = v.flags64()
		if err == nil {
			t.Errorf("Expected error but got flags %d from boolean %s", n64, v.describe())
		} else if _, ok := err.(*attributeValueFlagsTypeError); !ok {
			t.Errorf("Expected *attributeValueFlagsTypeError but got %T (%s)", err, err)
		}
	}
}

func TestAttributeValueSerialize(t *testing.T) {
	tt := Type(&builtinType{
		n: "Test-Type",
		k: "test-type",
	})
	v := AttributeValue{t: tt, v: nil}
	s, err := v.Serialize()
	if err == nil {
		t.Errorf("Expected error but got string %q", s)
	} else if _, ok := err.(*unknownTypeSerializationError); !ok {
		t.Errorf("Expected *unknownTypeSerializationError but got %T (%s)", err, err)
	}

	v = UndefinedValue
	s, err = v.Serialize()
	if err == nil {
		t.Errorf("Expected error but got string %q", s)
	} else if _, ok := err.(*invalidTypeSerializationError); !ok {
		t.Errorf("Expected *invalidTypeSerializationError but got %T (%s)", err, err)
	}

	v, err = MakeValueFromString(TypeBoolean, "true")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "true" {
			t.Errorf("Expected \"true\" but got %q", s)
		}
	}

	v, err = MakeValueFromString(TypeInteger, "47238")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "47238" {
			t.Errorf("Expected \"true\" but got %q", s)
		}
	}

	v, err = MakeValueFromString(TypeFloat, "3.1415927")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "3.1415927" {
			t.Errorf("Expected \"true\" but got %q", s)
		}
	}

	v, err = MakeValueFromString(TypeString, "test")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "test" {
			t.Errorf("Expected \"test\" but got %q", s)
		}
	}

	v, err = MakeValueFromString(TypeAddress, "192.0.2.1")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "192.0.2.1" {
			t.Errorf("Expected \"192.0.2.1\" but got %q", s)
		}
	}

	v, err = MakeValueFromString(TypeNetwork, "192.0.2.0/24")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "192.0.2.0/24" {
			t.Errorf("Expected \"192.0.2.0/24\" but got %q", s)
		}
	}

	v, err = MakeValueFromString(TypeDomain, "example.com")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else if s != "example.com" {
			t.Errorf("Expected \"example.com\" but got %q", s)
		}
	}

	sTree := strtree.NewTree()
	sTree.InplaceInsert("1 - one", 1)
	sTree.InplaceInsert("2 - two", 2)
	sTree.InplaceInsert("3 - three", 3)
	sTree.InplaceInsert("4 - four", 4)
	v = MakeSetOfStringsValue(sTree)
	s, err = v.Serialize()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		e := "\"1 - one\",\"2 - two\",\"3 - three\",\"4 - four\""
		if s != e {
			t.Errorf("Expected %q but got %q", e, s)
		}
	}

	nTree := iptree.NewTree()
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.16/28"), 1)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.32/28"), 2)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.48/28"), 3)
	nTree.InplaceInsertNet(makeTestNetwork("192.0.2.64/28"), 4)
	v = MakeSetOfNetworksValue(nTree)
	s, err = v.Serialize()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		e := "\"192.0.2.16/28\",\"192.0.2.32/28\",\"192.0.2.48/28\",\"192.0.2.64/28\""
		if s != e {
			t.Errorf("Expected %q but got %q", e, s)
		}
	}

	dTree := &domaintree.Node{}
	dTree.InplaceInsert(makeTestDomain("example.com"), 1)
	dTree.InplaceInsert(makeTestDomain("example.gov"), 2)
	dTree.InplaceInsert(makeTestDomain("example.net"), 3)
	dTree.InplaceInsert(makeTestDomain("example.org"), 4)
	v = MakeSetOfDomainsValue(dTree)
	s, err = v.Serialize()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		e := "\"example.com\",\"example.gov\",\"example.net\",\"example.org\""
		if s != e {
			t.Errorf("Expected %q but got %q", e, s)
		}
	}

	v = MakeListOfStringsValue([]string{"one", "two", "three", "four"})
	s, err = v.Serialize()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		e := "\"one\",\"two\",\"three\",\"four\""
		if s != e {
			t.Errorf("Expected %q but got %q", e, s)
		}
	}

	ft8, err := NewFlagsType("flags",
		"f00", "f01", "f02", "f03",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue8(7, ft8)
		s, err = v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"f00\",\"f01\",\"f02\""
			if s != e {
				t.Errorf("Expected %q but got %q", e, s)
			}
		}
	}

	ft16, err := NewFlagsType("flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue16(7, ft16)
		s, err = v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"f00\",\"f01\",\"f02\""
			if s != e {
				t.Errorf("Expected %q but got %q", e, s)
			}
		}
	}

	ft32, err := NewFlagsType("flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue32(7, ft32)
		s, err = v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"f00\",\"f01\",\"f02\""
			if s != e {
				t.Errorf("Expected %q but got %q", e, s)
			}
		}
	}

	ft64, err := NewFlagsType("flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		v = MakeFlagsValue64(7, ft64)
		s, err = v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got: %s", err)
		} else {
			e := "\"f00\",\"f01\",\"f02\""
			if s != e {
				t.Errorf("Expected %q but got %q", e, s)
			}
		}
	}
}

type testMetaType struct{}

func (t *testMetaType) String() string     { return "Test" }
func (t *testMetaType) GetKey() string     { return "test" }
func (t *testMetaType) Match(ot Type) bool { return true }

func TestAttributeValueRebind(t *testing.T) {
	s := MakeStringValue("test")
	v, err := s.Rebind(TypeString)
	if err != nil {
		t.Errorf("Expected no error but got: %s", err)
	} else {
		vt := v.GetResultType()
		if vt != TypeString {
			t.Errorf("Expected %q as value type but got %q", TypeString, vt)
		}

		if s, ok := v.v.(string); ok {
			e := "test"
			if s != e {
				t.Errorf("Expected string value %q but got %q", e, s)
			}
		} else {
			t.Errorf("Expected string value but got %T (%#v)", v.v, v.v)
		}
	}

	v, err = s.Rebind(TypeBoolean)
	if err == nil {
		t.Errorf("Expected error but got value %s", v.describe())
	} else if _, ok := err.(*notMatchingTypeRebindError); !ok {
		t.Errorf("Expected *notMatchingTypeRebindError but got %T (%s)", err, err)
	}

	ft8, err := NewFlagsType("F8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		gt8, err := NewFlagsType("G8flags",
			"g00", "g01", "g02", "g03", "g04", "g05", "g06", "g07",
		)
		if err != nil {
			t.Errorf("Expected no error but got %s", err)
		} else {
			f := MakeFlagsValue8(7, ft8)
			v, err := f.Rebind(gt8)
			if err != nil {
				t.Errorf("Expected no error but got: %s", err)
			} else {
				vt := v.GetResultType()
				if vt != gt8 {
					t.Errorf("Expected %q as value type but got %q", gt8, vt)
				}

				if n, ok := v.v.(uint8); ok {
					if n != 7 {
						t.Errorf("Expected 8 bits flags value %d but got %d", 7, n)
					}
				} else {
					t.Errorf("Expected 8 bits flags value but got %T (%#v)", v.v, v.v)
				}
			}
		}
	}

	ft8, err = NewFlagsType("F8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		gt8, err := NewFlagsType("G8flags",
			"g00", "g01", "g02", "g03",
		)
		f := MakeFlagsValue8(7, ft8)
		v, err := f.Rebind(gt8)
		if err == nil {
			t.Errorf("Expected error but got value %s", v.describe())
		} else if _, ok := err.(*notMatchingTypeRebindError); !ok {
			t.Errorf("Expected *notMatchingTypeRebindError but got %T (%s)", err, err)
		}
	}

	testType := new(testMetaType)
	mv := AttributeValue{t: testType}
	v, err = mv.Rebind(TypeBoolean)
	if err == nil {
		t.Errorf("Expected error but got value %s", v.describe())
	} else if _, ok := err.(*unknownMetaType); !ok {
		t.Errorf("Expected *unknownMetaType but got %T (%s)", err, err)
	}
}

func assertPanicWithError(t *testing.T, f func(), format string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			e := fmt.Sprintf(format, args...)
			err, ok := r.(error)
			if !ok {
				t.Errorf("Excpected error %q on panic but got %T (%#v)", e, r, r)
			} else if err.Error() != e {
				t.Errorf("Excpected error %q on panic but got %q", e, r)
			}
		} else {
			t.Errorf("Expected panic %q", fmt.Sprintf(format, args...))
		}
	}()

	f()
}

func makeTestNetwork(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}

	return n
}

func makeTestDomain(s string) domain.Name {
	n, err := domain.MakeNameFromString(s)
	if err != nil {
		panic(err)
	}

	return n
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

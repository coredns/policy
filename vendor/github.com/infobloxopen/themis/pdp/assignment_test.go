package pdp

import (
	"math"
	"net"
	"testing"
)

func TestAttributeAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	expect := "test-value"
	v := MakeStringValue(expect)
	a := Attribute{
		id: "test-id",
		t:  TypeString}

	aa := MakeAttributeAssignment(a, v)
	id, tKey, s, err := aa.Serialize(ctx)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if id != a.id || tKey != a.t.GetKey() || s != expect {
		t.Errorf("Expected %q, %q, %q but got %q, %q, %q", a.id, a.t.GetKey(), expect, id, tKey, s)
	}

	dv := MakeDomainValue(makeTestDomain("example.com"))
	v = MakeStringValue(expect)
	e := makeFunctionStringEqual(v, dv)
	a = Attribute{
		id: "test-id",
		t:  TypeBoolean}

	aa = MakeAttributeAssignment(a, e)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError error but got %T (%s)", err, err)
	}

	expect = "test-value"
	v = MakeStringValue(expect)
	a = Attribute{
		id: "test-id",
		t:  TypeBoolean}
	aa = MakeAttributeAssignment(a, v)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*assignmentTypeMismatch); !ok {
		t.Errorf("Expected *ssignmentTypeMismatch error but got %T (%s)", err, err)
	}

	fv := MakeFloatValue(2.718282)
	v = MakeStringValue(expect)
	e = makeFunctionStringEqual(v, fv)
	a = Attribute{
		id: "test-id",
		t:  TypeInteger}

	aa = MakeAttributeAssignment(a, e)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError error but got %T (%s)", err, err)
	}

	v = MakeFloatValue(1234.567)
	a = Attribute{
		id: "test-id",
		t:  TypeInteger}
	aa = MakeAttributeAssignment(a, v)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*assignmentTypeMismatch); !ok {
		t.Errorf("Expected *ssignmentTypeMismatch error but got %T (%s)", err, err)
	}

	iv := MakeIntegerValue(45678)
	v = MakeStringValue(expect)
	e = makeFunctionStringEqual(v, iv)
	a = Attribute{
		id: "test-id",
		t:  TypeFloat}

	aa = MakeAttributeAssignment(a, e)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError error but got %T (%s)", err, err)
	}

	expect = "45679.23"
	iv = MakeIntegerValue(45678)
	fv = MakeFloatValue(1.23)
	e = makeFunctionFloatAdd(fv, iv)
	a = Attribute{
		id: "test-id",
		t:  TypeFloat}

	aa = MakeAttributeAssignment(a, e)
	id, tKey, s, err = aa.Serialize(ctx)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if id != a.id || tKey != a.t.GetKey() || s != expect {
		t.Errorf("Expected %q, %q, %q but got %q, %q, %q", a.id, a.t.GetKey(), expect, id, tKey, s)
	}

	v = MakeIntegerValue(12345)
	a = Attribute{
		id: "test-id",
		t:  TypeFloat}
	aa = MakeAttributeAssignment(a, v)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*assignmentTypeMismatch); !ok {
		t.Errorf("Expected *ssignmentTypeMismatch error but got %T (%s)", err, err)
	}

	v = UndefinedValue
	a = Attribute{
		id: "test-id",
		t:  TypeUndefined}
	aa = MakeAttributeAssignment(a, v)
	id, tKey, s, err = aa.Serialize(ctx)
	if err == nil {
		t.Errorf("Expected error but got %q, %q, %q", id, tKey, s)
	} else if _, ok := err.(*invalidTypeSerializationError); !ok {
		t.Errorf("Expected *invalidTypeSerializationError error but got %T (%s)", err, err)
	}
}

func TestMakeExpressionAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeExpressionAssignment(
		"test-id",
		makeFunctionStringEqual(
			MakeStringValue("test"),
			MakeStringValue("example"),
		),
	)

	v, err := aa.calculate(ctx)
	if err != nil {
		t.Error(err)
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "false"
			if e != s {
				t.Errorf("Expected boolean value of %q but got %q", e, s)
			}
		}
	}

	aa = MakeExpressionAssignment(
		"test-id",
		makeFunctionStringEqual(MakeStringValue("test"), MakeStringDesignator("s")),
	)

	v, err = aa.calculate(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", v.describe())
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestBooleanAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeBooleanAssignment("test-id", true)

	b, err := aa.GetBoolean(ctx)
	if err != nil {
		t.Error(err)
	} else if !b {
		t.Errorf("Expected boolean value of %#v but got %#v", true, b)
	}

	s, err := aa.GetString(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %q", s)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeBooleanDesignator("b"),
	)

	b, err = aa.GetBoolean(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %#v", b)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestStringAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeStringAssignment("test-id", "test")

	s, err := aa.GetString(ctx)
	if err != nil {
		t.Error(err)
	} else if s != "test" {
		t.Errorf("Expected string value of %q but got %q", "test", s)
	}

	i, err := aa.GetInteger(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %d", i)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeStringDesignator("s"),
	)

	s, err = aa.GetString(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %q", s)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestIntegerAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeIntegerAssignment("test-id", math.MaxInt64)

	i, err := aa.GetInteger(ctx)
	if err != nil {
		t.Error(err)
	} else if i != math.MaxInt64 {
		t.Errorf("Expected integer value of %d but got %d", math.MaxInt64, i)
	}

	f, err := aa.GetFloat(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %g", f)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeIntegerDesignator("i"),
	)

	i, err = aa.GetInteger(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %d", i)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestFloatAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeFloatAssignment("test-id", math.SmallestNonzeroFloat64)

	f, err := aa.GetFloat(ctx)
	if err != nil {
		t.Error(err)
	} else if f != math.SmallestNonzeroFloat64 {
		t.Errorf("Expected float value of %g but got %g", math.SmallestNonzeroFloat64, f)
	}

	a, err := aa.GetAddress(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", a)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeFloatDesignator("f"),
	)

	f, err = aa.GetFloat(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %g", f)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestAddressAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeAddressAssignment("test-id", net.ParseIP("192.0.2.1"))

	a, err := aa.GetAddress(ctx)
	if err != nil {
		t.Error(err)
	} else if !a.Equal(net.ParseIP("192.0.2.1")) {
		t.Errorf("Expected address value of %s but got %s", "192.0.2.1", a)
	}

	n, err := aa.GetNetwork(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", n)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeAddressDesignator("a"),
	)

	a, err = aa.GetAddress(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", a)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestNetworkAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeNetworkAssignment("test-id", makeTestNetwork("192.0.2.0/24"))

	n, err := aa.GetNetwork(ctx)
	if err != nil {
		t.Error(err)
	} else if n.String() != "192.0.2.0/24" {
		t.Errorf("Expected network value of %s but got %s", "192.0.2.0/24", n)
	}

	d, err := aa.GetDomain(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", d)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeNetworkDesignator("n"),
	)

	n, err = aa.GetNetwork(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", n)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestDomainAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeDomainAssignment("test-id", makeTestDomain("example.com"))

	d, err := aa.GetDomain(ctx)
	if err != nil {
		t.Error(err)
	} else if d.String() != "example.com" {
		t.Errorf("Expected domain value of %s but got %s", "example.com", d)
	}

	ss, err := aa.GetSetOfStrings(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", serializeSetOfStrings(ss))
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeDomainDesignator("d"),
	)

	d, err = aa.GetDomain(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", d)
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestSetOfStringsAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeSetOfStringsAssignment("test-id", newStrTree("one", "two", "three"))

	ss, err := aa.GetSetOfStrings(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeSetOfStrings(ss) != "\"one\",\"two\",\"three\"" {
		t.Errorf("Expected set of strings value of %s but got %s",
			"\"one\",\"two\",\"three\"", serializeSetOfStrings(ss))
	}

	sn, err := aa.GetSetOfNetworks(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", serializeSetOfNetworks(sn))
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeSetOfStringsDesignator("ss"),
	)

	ss, err = aa.GetSetOfStrings(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeSetOfStrings(ss))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestSetOfNetworksAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeSetOfNetworksAssignment("test-id",
		newIPTree(
			makeTestNetwork("192.0.2.16/28"),
			makeTestNetwork("192.0.2.32/28"),
			makeTestNetwork("192.0.2.48/28"),
		),
	)

	sn, err := aa.GetSetOfNetworks(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeSetOfNetworks(sn) != "\"192.0.2.16/28\",\"192.0.2.32/28\",\"192.0.2.48/28\"" {
		t.Errorf("Expected set of networks value of %s but got %s",
			"\"192.0.2.16/28\",\"192.0.2.32/28\",\"192.0.2.48/28\"", serializeSetOfNetworks(sn))
	}

	sd, err := aa.GetSetOfDomains(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", serializeSetOfDomains(sd))
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeSetOfNetworksDesignator("sn"),
	)

	sn, err = aa.GetSetOfNetworks(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeSetOfNetworks(sn))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestSetOfDomainsAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeSetOfDomainsAssignment("test-id",
		newDomainTree(
			makeTestDomain("example.com"),
			makeTestDomain("example.net"),
			makeTestDomain("example.org"),
		),
	)

	sd, err := aa.GetSetOfDomains(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeSetOfDomains(sd) != "\"example.com\",\"example.net\",\"example.org\"" {
		t.Errorf("Expected set of domains value of %s but got %s",
			"\"example.com\",\"example.net\",\"example.org\"", serializeSetOfDomains(sd))
	}

	ls, err := aa.GetListOfStrings(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %s", serializeListOfStrings(ls))
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeSetOfDomainsDesignator("sd"),
	)

	sd, err = aa.GetSetOfDomains(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeSetOfDomains(sd))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestListOfStringsAssignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	aa := MakeListOfStringsAssignment("test-id", []string{"one", "two", "three"})

	ls, err := aa.GetListOfStrings(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeListOfStrings(ls) != "\"one\",\"two\",\"three\"" {
		t.Errorf("Expected list of strings value of %s but got %s",
			"\"one\",\"two\",\"three\"", serializeListOfStrings(ls))
	}

	b, err := aa.GetBoolean(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueTypeError but got value %#v", b)
	} else if _, ok := err.(*attributeValueTypeError); !ok {
		t.Errorf("Expected *attributeValueTypeError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeListOfStringsDesignator("ls"),
	)

	ls, err = aa.GetListOfStrings(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeListOfStrings(ls))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestFlags8Assignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	nt, err := NewFlagsType("F8flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Fatal(err)
	}

	ft8, ok := nt.(*FlagsType)
	if !ok {
		t.Fatalf("Expected *FlagsType but got %T", nt)
	}

	aa := MakeFlags8Assignment("test-id", nt, 21)

	f8, err := aa.GetFlags8(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeFlags(uint64(f8), ft8) != "\"f00\",\"f02\",\"f04\"" {
		t.Errorf("Expected 8 bits flags value of %s but got %s",
			"\"f00\",\"f02\",\"f04\"", serializeFlags(uint64(f8), ft8))
	}

	f16, err := aa.GetFlags16(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueFlagsBitsError but got value %s", serializeFlags(uint64(f16), ft8))
	} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
		t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeDesignator("f8", ft8),
	)

	f8, err = aa.GetFlags8(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeFlags(uint64(f8), ft8))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestFlags16Assignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	nt, err := NewFlagsType("F16flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
	)
	if err != nil {
		t.Fatal(err)
	}

	ft16, ok := nt.(*FlagsType)
	if !ok {
		t.Fatalf("Expected *FlagsType but got %T", nt)
	}

	aa := MakeFlags16Assignment("test-id", nt, 21)

	f16, err := aa.GetFlags16(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeFlags(uint64(f16), ft16) != "\"f00\",\"f02\",\"f04\"" {
		t.Errorf("Expected 16 bits flags value of %s but got %s",
			"\"f00\",\"f02\",\"f04\"", serializeFlags(uint64(f16), ft16))
	}

	f32, err := aa.GetFlags32(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueFlagsBitsError but got value %s", serializeFlags(uint64(f32), ft16))
	} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
		t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeDesignator("f16", ft16),
	)

	f16, err = aa.GetFlags16(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeFlags(uint64(f16), ft16))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestFlags32Assignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	nt, err := NewFlagsType("F32flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
	)
	if err != nil {
		t.Fatal(err)
	}

	ft32, ok := nt.(*FlagsType)
	if !ok {
		t.Fatalf("Expected *FlagsType but got %T", nt)
	}

	aa := MakeFlags32Assignment("test-id", nt, 21)

	f32, err := aa.GetFlags32(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeFlags(uint64(f32), ft32) != "\"f00\",\"f02\",\"f04\"" {
		t.Errorf("Expected 32 bits flags value of %s but got %s",
			"\"f00\",\"f02\",\"f04\"", serializeFlags(uint64(f32), ft32))
	}

	f64, err := aa.GetFlags64(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueFlagsBitsError but got value %s", serializeFlags(f64, ft32))
	} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
		t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeDesignator("f32", ft32),
	)

	f32, err = aa.GetFlags32(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeFlags(uint64(f32), ft32))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

func TestFlags64Assignment(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	nt, err := NewFlagsType("F64flags",
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
		t.Fatal(err)
	}

	ft64, ok := nt.(*FlagsType)
	if !ok {
		t.Fatalf("Expected *FlagsType but got %T", nt)
	}

	aa := MakeFlags64Assignment("test-id", nt, 21)

	f64, err := aa.GetFlags64(ctx)
	if err != nil {
		t.Error(err)
	} else if serializeFlags(f64, ft64) != "\"f00\",\"f02\",\"f04\"" {
		t.Errorf("Expected 64 bits flags value of %s but got %s",
			"\"f00\",\"f02\",\"f04\"", serializeFlags(f64, ft64))
	}

	f8, err := aa.GetFlags8(ctx)
	if err == nil {
		t.Errorf("Expected *attributeValueFlagsBitsError but got value %s", serializeFlags(uint64(f8), ft64))
	} else if _, ok := err.(*attributeValueFlagsBitsError); !ok {
		t.Errorf("Expected *attributeValueFlagsBitsError but got %T (%s)", err, err)
	}

	aa = MakeExpressionAssignment(
		"test-id",
		MakeDesignator("f64", ft64),
	)

	f64, err = aa.GetFlags64(ctx)
	if err == nil {
		t.Errorf("Expected *missingAttributeError but got value %s", serializeFlags(uint64(f64), ft64))
	} else if _, ok := err.(*missingAttributeError); !ok {
		t.Errorf("Expected *missingAttributeError but got %T (%s)", err, err)
	}
}

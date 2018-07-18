package pdp

import "testing"

func TestAttributeDesignator(t *testing.T) {
	testAttributes := []struct {
		id  string
		val AttributeValue
	}{
		{
			id:  "test-id",
			val: MakeStringValue("test-value"),
		},
		{
			id:  "test-id-i",
			val: MakeIntegerValue(12345),
		},
		{
			id:  "test-id-f",
			val: MakeFloatValue(67.89),
		},
	}

	ctx, err := NewContext(nil, len(testAttributes), func(i int) (string, AttributeValue, error) {
		return testAttributes[i].id, testAttributes[i].val, nil
	})
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	a := Attribute{
		id: "test-id",
		t:  TypeString}
	d := MakeAttributeDesignator(a)
	dai := d.GetID()
	if dai != "test-id" {
		t.Errorf("Expected %q id but got %q", "test-id", dai)
	}

	dat := d.GetResultType()
	if dat != TypeString {
		t.Errorf("Expected %q type but got %q", TypeString, dat)
	}

	_, err = d.Calculate(ctx)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
}

func TestMakeDesignator(t *testing.T) {
	d := MakeDesignator("test-id", TypeString)
	dai := d.GetID()
	if dai != "test-id" {
		t.Errorf("Expected %q id but got %q", "test-id", dai)
	}

	dat := d.GetResultType()
	if dat != TypeString {
		t.Errorf("Expected %q type but got %q", TypeString, dat)
	}

	testCases := []struct {
		f func(string) AttributeDesignator
		i string
		t Type
	}{
		{
			f: MakeBooleanDesignator,
			i: "b",
			t: TypeBoolean,
		},
		{
			f: MakeStringDesignator,
			i: "s",
			t: TypeString,
		},
		{
			f: MakeIntegerDesignator,
			i: "i",
			t: TypeInteger,
		},
		{
			f: MakeFloatDesignator,
			i: "f",
			t: TypeFloat,
		},
		{
			f: MakeAddressDesignator,
			i: "a",
			t: TypeAddress,
		},
		{
			f: MakeNetworkDesignator,
			i: "n",
			t: TypeNetwork,
		},
		{
			f: MakeDomainDesignator,
			i: "d",
			t: TypeDomain,
		},
		{
			f: MakeSetOfStringsDesignator,
			i: "ss",
			t: TypeSetOfStrings,
		},
		{
			f: MakeSetOfNetworksDesignator,
			i: "sn",
			t: TypeSetOfNetworks,
		},
		{
			f: MakeSetOfDomainsDesignator,
			i: "sd",
			t: TypeSetOfDomains,
		},
		{
			f: MakeListOfStringsDesignator,
			i: "ls",
			t: TypeListOfStrings,
		},
	}

	for i, c := range testCases {
		d := c.f(c.i)

		dai := d.GetID()
		if dai != c.i {
			t.Errorf("Expected %q id for %d designator but got %q", c.i, i+1, dai)
		}

		dat := d.GetResultType()
		if dat != c.t {
			t.Errorf("Expected %q type for %d %q designator but got %q", c.t, i+1, dai, dat)
		}
	}
}

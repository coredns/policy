package pdp

import (
	"fmt"
	"net"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
	"github.com/pmezard/go-difflib/difflib"
)

// AttributeAssignment represents assignment of arbitrary result to
// an attribute.
type AttributeAssignment struct {
	a Attribute
	e Expression
}

// MakeAttributeAssignment creates assignment of given expression to given
// attribute.
func MakeAttributeAssignment(a Attribute, e Expression) AttributeAssignment {
	return AttributeAssignment{
		a: a,
		e: e,
	}
}

// MakeExpressionAssignment creates attribute assignment for attribute with
// given id and type derived from given expression.
func MakeExpressionAssignment(id string, e Expression) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, e.GetResultType()),
		e,
	)
}

// MakeBooleanAssignment creates attribute assignment for boolean value.
func MakeBooleanAssignment(id string, v bool) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeBoolean),
		MakeBooleanValue(v),
	)
}

// MakeStringAssignment creates attribute assignment for string value.
func MakeStringAssignment(id string, v string) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeString),
		MakeStringValue(v),
	)
}

// MakeIntegerAssignment creates attribute assignment for integer value.
func MakeIntegerAssignment(id string, v int64) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeInteger),
		MakeIntegerValue(v),
	)
}

// MakeFloatAssignment creates attribute assignment for float value.
func MakeFloatAssignment(id string, v float64) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeFloat),
		MakeFloatValue(v),
	)
}

// MakeAddressAssignment creates attribute assignment for address value.
func MakeAddressAssignment(id string, v net.IP) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeAddress),
		MakeAddressValue(v),
	)
}

// MakeNetworkAssignment creates attribute assignment for network value.
func MakeNetworkAssignment(id string, v *net.IPNet) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeNetwork),
		MakeNetworkValue(v),
	)
}

// MakeDomainAssignment creates attribute assignment for domain value.
func MakeDomainAssignment(id string, v domain.Name) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeDomain),
		MakeDomainValue(v),
	)
}

// MakeSetOfStringsAssignment creates attribute assignment for set of strings
// value.
func MakeSetOfStringsAssignment(id string, v *strtree.Tree) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeSetOfStrings),
		MakeSetOfStringsValue(v),
	)
}

// MakeSetOfNetworksAssignment creates attribute assignment for set of networks
// value.
func MakeSetOfNetworksAssignment(id string, v *iptree.Tree) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeSetOfNetworks),
		MakeSetOfNetworksValue(v),
	)
}

// MakeSetOfDomainsAssignment creates attribute assignment for set of domains
// value.
func MakeSetOfDomainsAssignment(id string, v *domaintree.Node) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeSetOfDomains),
		MakeSetOfDomainsValue(v),
	)
}

// MakeListOfStringsAssignment creates attribute assignment for list of strings
// value.
func MakeListOfStringsAssignment(id string, v []string) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, TypeListOfStrings),
		MakeListOfStringsValue(v),
	)
}

// MakeFlags8Assignment creates attribute assignment for flags value which fits
// 8 bits integer.
func MakeFlags8Assignment(id string, t Type, v uint8) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, t),
		MakeFlagsValue8(v, t),
	)
}

// MakeFlags16Assignment creates attribute assignment for flags value which fits
// 16 bits integer.
func MakeFlags16Assignment(id string, t Type, v uint16) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, t),
		MakeFlagsValue16(v, t),
	)
}

// MakeFlags32Assignment creates attribute assignment for flags value which fits
// 32 bits integer.
func MakeFlags32Assignment(id string, t Type, v uint32) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, t),
		MakeFlagsValue32(v, t),
	)
}

// MakeFlags64Assignment creates attribute assignment for flags value which fits
// 64 bits integer.
func MakeFlags64Assignment(id string, t Type, v uint64) AttributeAssignment {
	return MakeAttributeAssignment(
		MakeAttribute(id, t),
		MakeFlagsValue64(v, t),
	)
}

// GetID returns id of assignment's attribute.
func (a AttributeAssignment) GetID() string {
	return a.a.id
}

func (a AttributeAssignment) calculate(ctx *Context) (AttributeValue, error) {
	v, err := a.e.Calculate(ctx)
	if err != nil {
		return UndefinedValue, a.bindError(err)
	}

	return v, nil
}

func (a AttributeAssignment) GetValue() (AttributeValue, error) {
	v, ok := a.e.(AttributeValue)
	if !ok {
		return UndefinedValue, a.bindError(newRequestInvalidExpressionError(a))
	}

	return v, nil
}

// GetBoolean returns boolean value of assignment. It returns error if type of
// assignment is not boolean.
func (a AttributeAssignment) GetBoolean(ctx *Context) (bool, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return false, err
	}

	b, err := v.boolean()
	if err != nil {
		return false, a.bindError(err)
	}

	return b, nil
}

// GetString retruns string value of assignment. It returns error if type of
// assignment is not string.
func (a AttributeAssignment) GetString(ctx *Context) (string, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return "", err
	}

	s, err := v.str()
	if err != nil {
		return "", a.bindError(err)
	}

	return s, nil
}

// GetInteger retruns integer value of assignment. It returns error if type of
// assignment is not integer.
func (a AttributeAssignment) GetInteger(ctx *Context) (int64, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return 0, err
	}

	i, err := v.integer()
	if err != nil {
		return 0, a.bindError(err)
	}

	return i, nil
}

// GetFloat retruns float value of assignment. It returns error if type of
// assignment is not float.
func (a AttributeAssignment) GetFloat(ctx *Context) (float64, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return 0, err
	}

	f, err := v.float()
	if err != nil {
		return 0, a.bindError(err)
	}

	return f, nil
}

// GetAddress retruns address value of assignment. It returns error if type of
// assignment is not address.
func (a AttributeAssignment) GetAddress(ctx *Context) (net.IP, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return nil, err
	}

	ip, err := v.address()
	if err != nil {
		return nil, a.bindError(err)
	}

	return ip, nil
}

// GetNetwork retruns network value of assignment. It returns error if type of
// assignment is not network.
func (a AttributeAssignment) GetNetwork(ctx *Context) (*net.IPNet, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return nil, err
	}

	n, err := v.network()
	if err != nil {
		return nil, a.bindError(err)
	}

	return n, nil
}

// GetDomain retruns domain value of assignment. It returns error if type of
// assignment is not domain.
func (a AttributeAssignment) GetDomain(ctx *Context) (domain.Name, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return domain.Name{}, err
	}

	d, err := v.domain()
	if err != nil {
		return domain.Name{}, a.bindError(err)
	}

	return d, nil
}

// GetSetOfStrings retruns set of strings value of assignment. It returns error
// if type of assignment is not set of strings.
func (a AttributeAssignment) GetSetOfStrings(ctx *Context) (*strtree.Tree, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return nil, err
	}

	ss, err := v.setOfStrings()
	if err != nil {
		return nil, a.bindError(err)
	}

	return ss, nil
}

// GetSetOfNetworks retruns set of networks value of assignment. It returns
// error if type of assignment is not set of networks.
func (a AttributeAssignment) GetSetOfNetworks(ctx *Context) (*iptree.Tree, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return nil, err
	}

	sn, err := v.setOfNetworks()
	if err != nil {
		return nil, a.bindError(err)
	}

	return sn, nil
}

// GetSetOfDomains retruns set of networks value of assignment. It returns error
// if type of assignment is not set of domains.
func (a AttributeAssignment) GetSetOfDomains(ctx *Context) (*domaintree.Node, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return nil, err
	}

	sd, err := v.setOfDomains()
	if err != nil {
		return nil, a.bindError(err)
	}

	return sd, nil
}

// GetListOfStrings retruns list of strings value of assignment. It returns
// error if type of assignment is not list of strings.
func (a AttributeAssignment) GetListOfStrings(ctx *Context) ([]string, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return nil, err
	}

	ls, err := v.listOfStrings()
	if err != nil {
		return nil, a.bindError(err)
	}

	return ls, nil
}

// GetFlags8 retruns flags value of assignment which fits 8 bits integer.
// It returns error if type of assignment is not appropriate flags.
func (a AttributeAssignment) GetFlags8(ctx *Context) (uint8, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return 0, err
	}

	f8, err := v.flags8()
	if err != nil {
		return 0, a.bindError(err)
	}

	return f8, nil
}

// GetFlags16 retruns flags value of assignment which fits 16 bits integer.
// It returns error if type of assignment is not appropriate flags.
func (a AttributeAssignment) GetFlags16(ctx *Context) (uint16, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return 0, err
	}

	f16, err := v.flags16()
	if err != nil {
		return 0, a.bindError(err)
	}

	return f16, nil
}

// GetFlags32 retruns flags value of assignment which fits 32 bits integer.
// It returns error if type of assignment is not appropriate flags.
func (a AttributeAssignment) GetFlags32(ctx *Context) (uint32, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return 0, err
	}

	f32, err := v.flags32()
	if err != nil {
		return 0, a.bindError(err)
	}

	return f32, nil
}

// GetFlags64 retruns flags value of assignment which fits 64 bits integer.
// It returns error if type of assignment is not appropriate flags.
func (a AttributeAssignment) GetFlags64(ctx *Context) (uint64, error) {
	v, err := a.calculate(ctx)
	if err != nil {
		return 0, err
	}

	f64, err := v.flags64()
	if err != nil {
		return 0, a.bindError(err)
	}

	return f64, nil
}

// Serialize evaluates assignment and returns string representation of
// resulting attribute name, type and value or error if the evaluaction
// can't be done.
func (a AttributeAssignment) Serialize(ctx *Context) (string, string, string, error) {
	ID := a.GetID()
	k := a.a.GetType().GetKey()

	v, err := a.calculate(ctx)
	if err != nil {
		return ID, k, "", err
	}

	t := v.GetResultType()
	if a.a.GetType() != t {
		return ID, k, "", a.bindError(newAssignmentTypeMismatch(a.a, t))
	}

	s, err := v.Serialize()
	if err != nil {
		return ID, k, "", a.bindError(err)
	}

	return ID, k, s, nil
}

func (a AttributeAssignment) String() string {
	name, valueType, value, err := a.Serialize(nil)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("(%s)%s:%s", valueType, name, value)
}

func (a AttributeAssignment) bindError(err error) error {
	return bindErrorf(err, "assignment to %q", a.GetID())
}

func serializeAssignmentsForAssert(desc string, expected bool, a []AttributeAssignment) ([]string, error) {
	ctx, _ := NewContext(nil, 0, nil)

	out := make([]string, len(a))
	for i, a := range a {
		id, tName, s, err := a.Serialize(ctx)
		if err != nil {
			attr := "attribute"
			if expected {
				attr = "expected " + attr
			}

			return out, fmt.Errorf("can't serialize %s %d %q for %s: %s", attr, i+1, id, desc, err)
		}

		out[i] = fmt.Sprintf("%q.(%q) = %q\n", id, tName, s)
	}

	return out, nil
}

func AssertAttributeAssignments(t *testing.T, desc string, a []AttributeAssignment, e ...AttributeAssignment) {
	sa, err := serializeAssignmentsForAssert(desc, false, a)
	if err != nil {
		t.Error(err)
		return
	}

	se, err := serializeAssignmentsForAssert(desc, true, e)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := difflib.ContextDiff{
		A:        se,
		B:        sa,
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("can't compare \"%s\": %s", desc, err))
	}

	if len(diff) > 0 {
		t.Errorf("\"%s\" doesn't match:\n%s", desc, diff)
	}
}

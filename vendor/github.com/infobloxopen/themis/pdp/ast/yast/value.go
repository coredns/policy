package yast

import (
	"net"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx context) unmarshalStringValue(v interface{}) (pdp.AttributeValue, boundError) {
	s, err := ctx.validateString(v, "value of string type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	return pdp.MakeStringValue(s), nil
}

func (ctx context) unmarshalIntegerValue(v interface{}) (pdp.AttributeValue, boundError) {
	n, err := ctx.validateInteger(v, "value of integer type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	return pdp.MakeIntegerValue(n), nil
}

func (ctx context) unmarshalFloatValue(v interface{}) (pdp.AttributeValue, boundError) {
	n, err := ctx.validateFloat(v, "value of float type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	return pdp.MakeFloatValue(n), nil
}

func (ctx context) unmarshalAddressValue(v interface{}) (pdp.AttributeValue, boundError) {
	s, err := ctx.validateString(v, "value of address type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	a := net.ParseIP(s)
	if a == nil {
		return pdp.UndefinedValue, newInvalidAddressError(s)
	}

	return pdp.MakeAddressValue(a), nil
}

func (ctx context) unmarshalNetworkValue(v interface{}) (pdp.AttributeValue, boundError) {
	s, err := ctx.validateString(v, "value of network type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	_, n, ierr := net.ParseCIDR(s)
	if ierr != nil {
		return pdp.UndefinedValue, newInvalidNetworkError(s, ierr)
	}

	return pdp.MakeNetworkValue(n), nil
}

func (ctx context) unmarshalDomainValue(v interface{}) (pdp.AttributeValue, boundError) {
	s, err := ctx.validateString(v, "value of domain type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	d, ierr := domain.MakeNameFromString(s)
	if ierr != nil {
		return pdp.UndefinedValue, newInvalidDomainError(s, ierr)
	}

	return pdp.MakeDomainValue(d), nil
}

func (ctx context) unmarshalSetOfStringsValueItem(v interface{}, i int, set *strtree.Tree) boundError {
	s, err := ctx.validateString(v, "element")
	if err != nil {
		return err
	}

	set.InplaceInsert(s, i)
	return nil
}

func (ctx context) unmarshalSetOfStringsValue(v interface{}) (pdp.AttributeValue, boundError) {
	items, err := ctx.validateList(v, "set of strings")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	set := strtree.NewTree()
	for i, item := range items {
		err = ctx.unmarshalSetOfStringsValueItem(item, i, set)
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "set of strings")
		}
	}

	return pdp.MakeSetOfStringsValue(set), nil
}

func (ctx context) unmarshalSetOfNetworksValueItem(v interface{}, i int, set *iptree.Tree) boundError {
	s, err := ctx.validateString(v, "element")
	if err != nil {
		return err
	}

	_, n, ierr := net.ParseCIDR(s)
	if ierr != nil {
		return newInvalidNetworkError(s, ierr)
	}

	set.InplaceInsertNet(n, i)

	return nil
}

func (ctx context) unmarshalSetOfNetworksValue(v interface{}) (pdp.AttributeValue, boundError) {
	items, err := ctx.validateList(v, "set of networks")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	set := iptree.NewTree()
	for i, item := range items {
		err = ctx.unmarshalSetOfNetworksValueItem(item, i, set)
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "set of networks")
		}
	}

	return pdp.MakeSetOfNetworksValue(set), nil
}

func (ctx context) unmarshalSetOfDomainsValueItem(v interface{}, i int, set *domaintree.Node) boundError {
	s, err := ctx.validateString(v, "element")
	if err != nil {
		return err
	}

	d, ierr := domain.MakeNameFromString(s)
	if ierr != nil {
		return newInvalidDomainError(s, ierr)
	}

	set.InplaceInsert(d, i)

	return nil
}

func (ctx context) unmarshalSetOfDomainsValue(v interface{}) (pdp.AttributeValue, boundError) {
	items, err := ctx.validateList(v, "")
	if err != nil {
		return pdp.UndefinedValue, nil
	}

	set := &domaintree.Node{}
	for i, item := range items {
		err = ctx.unmarshalSetOfDomainsValueItem(item, i, set)
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "set of domains")
		}
	}

	return pdp.MakeSetOfDomainsValue(set), nil
}

func (ctx context) unmarshalListOfStringsValueItem(v interface{}, list []string) ([]string, boundError) {
	s, err := ctx.validateString(v, "element")
	if err != nil {
		return list, err
	}

	return append(list, s), nil
}

func (ctx context) unmarshalListOfStringsValue(v interface{}) (pdp.AttributeValue, boundError) {
	items, err := ctx.validateList(v, "list of strings")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	list := []string{}
	for i, item := range items {
		list, err = ctx.unmarshalListOfStringsValueItem(item, list)
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "list of strings")
		}
	}

	return pdp.MakeListOfStringsValue(list), nil
}

func (ctx context) unmarshalFlagsValue(v interface{}, t *pdp.FlagsType) (pdp.AttributeValue, boundError) {
	f, err := ctx.validateList(v, "flag names")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	switch t.Capacity() {
	case 8:
		return ctx.unmarshalFlags8Value(f, t)

	case 16:
		return ctx.unmarshalFlags16Value(f, t)

	case 32:
		return ctx.unmarshalFlags32Value(f, t)
	}

	return ctx.unmarshalFlags64Value(f, t)
}

func (ctx context) unmarshalFlags8Value(v []interface{}, t *pdp.FlagsType) (pdp.AttributeValue, boundError) {
	var n uint8

	for i, v := range v {
		f, err := ctx.validateString(v, "flag name")
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "flag names")
		}

		b := t.GetFlagBit(f)
		if b < 0 {
			return pdp.UndefinedValue, bindError(bindErrorf(newUnknownFlagNameError(f, t), "%d", i), "flag names")
		}

		n |= 1 << uint(b)
	}

	return pdp.MakeFlagsValue8(n, t), nil
}

func (ctx context) unmarshalFlags16Value(v []interface{}, t *pdp.FlagsType) (pdp.AttributeValue, boundError) {
	var n uint16

	for i, v := range v {
		f, err := ctx.validateString(v, "flag name")
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "flag names")
		}

		b := t.GetFlagBit(f)
		if b < 0 {
			return pdp.UndefinedValue, bindError(bindErrorf(newUnknownFlagNameError(f, t), "%d", i), "flag names")
		}

		n |= 1 << uint(b)
	}

	return pdp.MakeFlagsValue16(n, t), nil
}

func (ctx context) unmarshalFlags32Value(v []interface{}, t *pdp.FlagsType) (pdp.AttributeValue, boundError) {
	var n uint32

	for i, v := range v {
		f, err := ctx.validateString(v, "flag name")
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "flag names")
		}

		b := t.GetFlagBit(f)
		if b < 0 {
			return pdp.UndefinedValue, bindError(bindErrorf(newUnknownFlagNameError(f, t), "%d", i), "flag names")
		}

		n |= 1 << uint(b)
	}

	return pdp.MakeFlagsValue32(n, t), nil
}

func (ctx context) unmarshalFlags64Value(v []interface{}, t *pdp.FlagsType) (pdp.AttributeValue, boundError) {
	var n uint64

	for i, v := range v {
		f, err := ctx.validateString(v, "flag name")
		if err != nil {
			return pdp.UndefinedValue, bindError(bindErrorf(err, "%d", i), "flag names")
		}

		b := t.GetFlagBit(f)
		if b < 0 {
			return pdp.UndefinedValue, bindError(bindErrorf(newUnknownFlagNameError(f, t), "%d", i), "flag names")
		}

		n |= 1 << uint(b)
	}

	return pdp.MakeFlagsValue64(n, t), nil
}

func (ctx context) unmarshalValueByType(t pdp.Type, v interface{}) (pdp.AttributeValue, boundError) {
	if t, ok := t.(*pdp.FlagsType); ok {
		return ctx.unmarshalFlagsValue(v, t)
	}

	switch t {
	case pdp.TypeString:
		return ctx.unmarshalStringValue(v)

	case pdp.TypeInteger:
		return ctx.unmarshalIntegerValue(v)

	case pdp.TypeFloat:
		return ctx.unmarshalFloatValue(v)

	case pdp.TypeAddress:
		return ctx.unmarshalAddressValue(v)

	case pdp.TypeNetwork:
		return ctx.unmarshalNetworkValue(v)

	case pdp.TypeDomain:
		return ctx.unmarshalDomainValue(v)

	case pdp.TypeSetOfStrings:
		return ctx.unmarshalSetOfStringsValue(v)

	case pdp.TypeSetOfNetworks:
		return ctx.unmarshalSetOfNetworksValue(v)

	case pdp.TypeSetOfDomains:
		return ctx.unmarshalSetOfDomainsValue(v)

	case pdp.TypeListOfStrings:
		return ctx.unmarshalListOfStringsValue(v)
	}

	return pdp.UndefinedValue, newNotImplementedValueTypeError(t)
}

func (ctx context) unmarshalValue(v interface{}) (pdp.AttributeValue, boundError) {
	m, err := ctx.validateMap(v, "value attributes")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	strT, err := ctx.extractString(m, yastTagType, "type")
	if err != nil {
		return pdp.UndefinedValue, err
	}

	t := ctx.symbols.GetType(strT)
	if t == nil {
		return pdp.UndefinedValue, newUnknownTypeError(strT)
	}

	if t == pdp.TypeUndefined {
		return pdp.UndefinedValue, newInvalidTypeError(t)
	}

	c, ok := m[yastTagContent]
	if !ok {
		return pdp.UndefinedValue, newMissingContentError()
	}

	return ctx.unmarshalValueByType(t, c)
}

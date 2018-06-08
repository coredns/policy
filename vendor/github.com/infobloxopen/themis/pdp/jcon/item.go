package jcon

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

type contentItem struct {
	id string
	s  pdp.Symbols

	k      pdp.Signature
	keysOk bool

	t   pdp.Type
	tOk bool

	v      interface{}
	vOk    bool
	vReady bool
}

func (c *contentItem) unmarshalTypeField(d *json.Decoder) error {
	if c.tOk {
		return newDuplicateContentItemFieldError("type")
	}

	token, err := d.Token()
	if err != nil {
		return err
	}

	switch v := token.(type) {
	default:
		return newInvalidTypeFormatError(token)

	case string:
		c.t = c.s.GetType(v)
		if c.t == nil {
			return newUnknownTypeError(v)
		}

		if c.t == pdp.TypeUndefined {
			return newInvalidContentItemTypeError(c.t)
		}

		c.tOk = true

	case json.Delim:
		if v.String() != jparser.DelimObjectStart {
			return newInvalidTypeFormatError(token)
		}

		if err := c.unmarshalTypeDeclaration(d); err != nil {
			return err
		}
	}

	return nil
}

func (c *contentItem) unmarshalTypeDeclaration(d *json.Decoder) error {
	var (
		metaOk bool
		meta   string

		nameOk bool
		name   string

		flags []string
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		switch k {
		default:
			return newUnknownTypeFieldError(k)

		case "meta":
			s, err := jparser.GetString(d, "meta type name")
			if err != nil {
				return err
			}

			meta = s
			metaOk = true

		case "name":
			s, err := jparser.GetString(d, "type name")
			if err != nil {
				return err
			}

			name = s
			nameOk = true

		case "flags":
			flags = []string{}
			if err := jparser.GetStringSequence(d, func(i int, s string) error {
				flags = append(flags, s)
				return nil
			}, "list of flag names"); err != nil {
				return err
			}
		}

		return nil
	}, "type declarations"); err != nil {
		return err
	}

	if !metaOk {
		return newMissingMetaTypeNameError()
	}

	switch strings.ToLower(meta) {
	default:
		return newUnknownMetaTypeError(meta)

	case "flags":
		if flags == nil {
			return newMissingFlagNameListError()
		}

		t, err := pdp.NewFlagsType(name, flags...)
		if err != nil {
			return err
		}

		if nameOk {
			if err := c.s.PutType(t); err != nil {
				if _, ok := err.(*pdp.ReadOnlySymbolsChangeError); ok {
					return newNewTypeOnUpdateError()
				}

				return err
			}
		}

		c.t = t
		c.tOk = true
	}

	return nil
}

func (c *contentItem) unmarshalKeysField(d *json.Decoder) error {
	if c.keysOk {
		return newDuplicateContentItemFieldError("keys")
	}

	err := jparser.CheckArrayStart(d, "content item keys")
	if err != nil {
		return err
	}

	k := pdp.MakeSignature()
	i := 1
	for {
		src := fmt.Sprintf("key %d", i)
		t, err := d.Token()
		if err != nil {
			return bindError(err, src)
		}

		switch s := t.(type) {
		default:
			return newStringCastError(t, src)

		case string:
			t, ok := pdp.BuiltinTypes[strings.ToLower(s)]
			if !ok {
				return bindError(newUnknownTypeError(s), src)
			}

			if !pdp.ContentKeyTypes.Contains(t) {
				return bindError(newInvalidContentKeyTypeError(t, pdp.ContentKeyTypes), src)
			}

			k = append(k, t)
			i++

		case json.Delim:
			if s.String() != jparser.DelimArrayEnd {
				return newArrayEndDelimiterError(s, jparser.DelimArrayEnd, src)
			}

			c.k = k
			c.keysOk = true

			return nil
		}
	}
}

func (c *contentItem) unmarshalMap(d *json.Decoder, keyIdx int) (interface{}, error) {
	src := fmt.Sprintf("level %d map", keyIdx+1)
	err := jparser.CheckObjectStart(d, src)
	if err != nil {
		return nil, err
	}

	m, err := newTypedMap(c, keyIdx)
	if err != nil {
		return nil, err
	}

	err = jparser.UnmarshalObject(d, m.unmarshal, src)
	if err != nil {
		return nil, err
	}

	return m.get(), nil
}

func (c *contentItem) unmarshalValue(d *json.Decoder) (interface{}, error) {
	if t, ok := c.t.(*pdp.FlagsType); ok {
		var n uint64
		err := jparser.GetStringSequence(d, func(idx int, s string) error {
			i := t.GetFlagBit(s)
			if i < 0 {
				return newUnknownFlagNameError(s)
			}

			n |= 1 << uint(i)

			return nil
		}, "flag names")
		if err != nil {
			return 0, err
		}

		switch t.Capacity() {
		case 8:
			return uint8(n), nil

		case 16:
			return uint16(n), nil

		case 32:
			return uint32(n), nil
		}

		return n, nil
	}

	switch c.t {
	case pdp.TypeBoolean:
		return jparser.GetBoolean(d, "value")

	case pdp.TypeString:
		return jparser.GetString(d, "value")

	case pdp.TypeInteger:
		x, err := jparser.GetNumber(d, "value")
		if err != nil {
			return nil, err
		}

		if x < -9007199254740992 || x > 9007199254740992 {
			return nil, newIntegerOverflowError(x)
		}

		return int64(x), nil

	case pdp.TypeFloat:
		x, err := jparser.GetNumber(d, "value")
		if err != nil {
			return nil, err
		}

		return float64(x), nil

	case pdp.TypeAddress:
		s, err := jparser.GetString(d, "address value")
		if err != nil {
			return nil, err
		}

		a := net.ParseIP(s)
		if a == nil {
			return nil, newAddressCastError(s)
		}

		return a, nil

	case pdp.TypeNetwork:
		s, err := jparser.GetString(d, "network value")
		if err != nil {
			return nil, err
		}

		_, n, err := net.ParseCIDR(s)
		if err != nil {
			return nil, newNetworkCastError(s, err)
		}

		return n, nil

	case pdp.TypeDomain:
		s, err := jparser.GetString(d, "domain value")
		if err != nil {
			return nil, err
		}

		d, err := domain.MakeNameFromString(s)
		if err != nil {
			return nil, newDomainCastError(s, err)
		}

		return d, nil

	case pdp.TypeSetOfStrings:
		m := strtree.NewTree()
		i := 0
		err := jparser.GetStringSequence(d, func(idx int, s string) error {
			if _, ok := m.Get(s); !ok {
				m.InplaceInsert(s, i)
				i++
			}

			return nil
		}, "set of strings value")
		if err != nil {
			return nil, err
		}

		return m, nil

	case pdp.TypeSetOfNetworks:
		m := iptree.NewTree()
		err := jparser.GetStringSequence(d, func(idx int, s string) error {
			a := net.ParseIP(s)
			if a != nil {
				m.InplaceInsertIP(a, nil)
			} else {
				_, n, err := net.ParseCIDR(s)
				if err != nil {
					return bindErrorf(newAddressNetworkCastError(s, err), "%d", idx)
				}

				m.InplaceInsertNet(n, nil)
			}

			return nil
		}, "set of networks value")
		if err != nil {
			return nil, err
		}

		return m, nil

	case pdp.TypeSetOfDomains:
		m := &domaintree.Node{}
		err := jparser.GetStringSequence(d, func(idx int, s string) error {
			dn, err := domain.MakeNameFromString(s)
			if err != nil {
				return bindErrorf(newDomainCastError(s, err), "%d", idx)
			}

			m.InplaceInsert(dn, nil)
			return nil
		}, "set of domains value")
		if err != nil {
			return nil, err
		}

		return m, nil

	case pdp.TypeListOfStrings:
		lst := []string{}
		err := jparser.GetStringSequence(d, func(idx int, s string) error {
			lst = append(lst, s)
			return nil
		}, "list of strings value")
		if err != nil {
			return nil, err
		}

		return lst, nil
	}

	return nil, newInvalidContentItemTypeError(c.t)
}

func (c *contentItem) unmarshalFlags8Value(d *json.Decoder) (uint8, error) {
	t, ok := c.t.(*pdp.FlagsType)
	if !ok {
		return 0, newInvalidContentItemTypeError(c.t)
	}

	var n uint8
	err := jparser.GetStringSequence(d, func(idx int, s string) error {
		i := t.GetFlagBit(s)
		if i < 0 {
			return newUnknownFlagNameError(s)
		}

		n |= 1 << uint(i)

		return nil
	}, "flag names")
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (c *contentItem) unmarshalFlags16Value(d *json.Decoder) (uint16, error) {
	t, ok := c.t.(*pdp.FlagsType)
	if !ok {
		return 0, newInvalidContentItemTypeError(c.t)
	}

	var n uint16
	err := jparser.GetStringSequence(d, func(idx int, s string) error {
		i := t.GetFlagBit(s)
		if i < 0 {
			return newUnknownFlagNameError(s)
		}

		n |= 1 << uint(i)

		return nil
	}, "flag names")
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (c *contentItem) unmarshalFlags32Value(d *json.Decoder) (uint32, error) {
	t, ok := c.t.(*pdp.FlagsType)
	if !ok {
		return 0, newInvalidContentItemTypeError(c.t)
	}

	var n uint32
	err := jparser.GetStringSequence(d, func(idx int, s string) error {
		i := t.GetFlagBit(s)
		if i < 0 {
			return newUnknownFlagNameError(s)
		}

		n |= 1 << uint(i)

		return nil
	}, "flag names")
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (c *contentItem) unmarshalFlags64Value(d *json.Decoder) (uint64, error) {
	t, ok := c.t.(*pdp.FlagsType)
	if !ok {
		return 0, newInvalidContentItemTypeError(c.t)
	}

	var n uint64
	err := jparser.GetStringSequence(d, func(idx int, s string) error {
		i := t.GetFlagBit(s)
		if i < 0 {
			return newUnknownFlagNameError(s)
		}

		n |= 1 << uint(i)

		return nil
	}, "flag names")
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (c *contentItem) unmarshalTypedData(d *json.Decoder, keyIdx int) (interface{}, error) {
	if len(c.k) > keyIdx {
		return c.unmarshalMap(d, keyIdx)
	}

	return c.unmarshalValue(d)
}

func (c *contentItem) unmarshalDataField(d *json.Decoder) error {
	if c.vOk {
		return newDuplicateContentItemFieldError("type")
	}

	c.vReady = c.keysOk && c.tOk
	if c.vReady {
		v, err := c.unmarshalTypedData(d, 0)
		if err != nil {
			return err
		}

		if len(c.k) <= 0 {
			c.v = pdp.MakeContentValue(v)
		} else {
			c.v = v
		}
	} else {
		v, err := jparser.GetUndefined(d, "content")
		if err != nil {
			return err
		}

		c.v = v
	}

	c.vOk = true
	return nil
}

func (c *contentItem) unmarshal(k string, d *json.Decoder) error {
	switch strings.ToLower(k) {
	case "type":
		return c.unmarshalTypeField(d)

	case "keys":
		return c.unmarshalKeysField(d)

	case "data":
		return c.unmarshalDataField(d)
	}

	return newUnknownContentItemFieldError(k)
}

func (c *contentItem) adjustValue(v interface{}) pdp.ContentSubItem {
	cv, ok := v.(pdp.ContentSubItem)
	if !ok {
		panic(fmt.Errorf("expected value of type ContentSubItem when item is ready but got %T", v))
	}

	return cv
}

func (c *contentItem) get() (*pdp.ContentItem, error) {
	if !c.vOk {
		return nil, newMissingContentDataError()
	}

	if !c.tOk {
		return nil, newMissingContentTypeError()
	}

	if c.vReady {
		return pdp.MakeContentMappingItem(c.id, c.t, c.k, c.adjustValue(c.v)), nil
	}

	v, err := c.postProcess(c.v, 0)
	if err != nil {
		return nil, err
	}

	if len(c.k) <= 0 {
		return pdp.MakeContentValueItem(c.id, c.t, v), nil
	}

	return pdp.MakeContentMappingItem(c.id, c.t, c.k, c.adjustValue(v)), nil
}

func unmarshalContentItem(id string, s pdp.Symbols, d *json.Decoder) (*pdp.ContentItem, error) {
	err := jparser.CheckObjectStart(d, "content item")
	if err != nil {
		return nil, err
	}

	item := &contentItem{
		id: id,
		s:  s,
	}
	err = jparser.UnmarshalObject(d, item.unmarshal, "content item")
	if err != nil {
		return nil, err
	}

	return item.get()
}

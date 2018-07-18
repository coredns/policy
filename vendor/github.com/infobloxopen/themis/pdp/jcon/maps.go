package jcon

import (
	"encoding/json"
	"net"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
	"github.com/infobloxopen/go-trees/uintX/domaintree16"
	"github.com/infobloxopen/go-trees/uintX/domaintree32"
	"github.com/infobloxopen/go-trees/uintX/domaintree64"
	"github.com/infobloxopen/go-trees/uintX/domaintree8"
	"github.com/infobloxopen/go-trees/uintX/iptree16"
	"github.com/infobloxopen/go-trees/uintX/iptree32"
	"github.com/infobloxopen/go-trees/uintX/iptree64"
	"github.com/infobloxopen/go-trees/uintX/iptree8"
	"github.com/infobloxopen/go-trees/uintX/strtree16"
	"github.com/infobloxopen/go-trees/uintX/strtree32"
	"github.com/infobloxopen/go-trees/uintX/strtree64"
	"github.com/infobloxopen/go-trees/uintX/strtree8"
	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

type mapUnmarshaller interface {
	get() interface{}
	unmarshal(k string, d *json.Decoder) error
	postProcess(p jparser.Pair) error
}

func newTypedMap(c *contentItem, keyIdx int) (mapUnmarshaller, error) {
	t := c.k[keyIdx]

	switch t {
	case pdp.TypeString:
		if t, ok := c.t.(*pdp.FlagsType); ok {
			switch t.Capacity() {
			case 8:
				return &string8Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               strtree8.NewTree()}, nil

			case 16:
				return &string16Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               strtree16.NewTree()}, nil

			case 32:
				return &string32Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               strtree32.NewTree()}, nil
			}

			return &string64Map{
				contentItemLink: contentItemLink{c: c, i: keyIdx},
				m:               strtree64.NewTree()}, nil
		}

		return &stringMap{
			contentItemLink: contentItemLink{c: c, i: keyIdx},
			m:               strtree.NewTree()}, nil

	case pdp.TypeAddress, pdp.TypeNetwork:
		if t, ok := c.t.(*pdp.FlagsType); ok {
			switch t.Capacity() {
			case 8:
				return &network8Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               iptree8.NewTree()}, nil

			case 16:
				return &network16Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               iptree16.NewTree()}, nil

			case 32:
				return &network32Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               iptree32.NewTree()}, nil
			}

			return &network64Map{
				contentItemLink: contentItemLink{c: c, i: keyIdx},
				m:               iptree64.NewTree()}, nil
		}

		return &networkMap{
			contentItemLink: contentItemLink{c: c, i: keyIdx},
			m:               iptree.NewTree()}, nil

	case pdp.TypeDomain:
		if t, ok := c.t.(*pdp.FlagsType); ok {
			switch t.Capacity() {
			case 8:
				return &domain8Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               &domaintree8.Node{}}, nil

			case 16:
				return &domain16Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               &domaintree16.Node{}}, nil

			case 32:
				return &domain32Map{
					contentItemLink: contentItemLink{c: c, i: keyIdx},
					m:               &domaintree32.Node{}}, nil
			}

			return &domain64Map{
				contentItemLink: contentItemLink{c: c, i: keyIdx},
				m:               &domaintree64.Node{}}, nil
		}

		return &domainMap{
			contentItemLink: contentItemLink{c: c, i: keyIdx},
			m:               &domaintree.Node{}}, nil
	}

	return nil, newInvalidContentKeyTypeError(t, pdp.ContentKeyTypes)
}

type contentItemLink struct {
	c *contentItem
	i int
}

type stringMap struct {
	contentItemLink
	m *strtree.Tree
}

func (m *stringMap) get() interface{} {
	return pdp.MakeContentStringMap(m.m)
}

func (m *stringMap) unmarshal(k string, d *json.Decoder) error {
	v, err := m.c.unmarshalTypedData(d, m.i+1)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(k, v)

	return nil
}

func (m *stringMap) postProcess(p jparser.Pair) error {
	v, err := m.c.postProcess(p.V, m.i+1)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(p.K, v)

	return nil
}

type string8Map struct {
	contentItemLink
	m *strtree8.Tree
}

func (m *string8Map) get() interface{} {
	return pdp.MakeContentStringFlags8Map(m.m)
}

func (m *string8Map) unmarshal(k string, d *json.Decoder) error {
	v, err := m.c.unmarshalFlags8Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(k, v)

	return nil
}

func (m *string8Map) postProcess(p jparser.Pair) error {
	v, err := m.c.postProcessFlags8Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(p.K, v)

	return nil
}

type string16Map struct {
	contentItemLink
	m *strtree16.Tree
}

func (m *string16Map) get() interface{} {
	return pdp.MakeContentStringFlags16Map(m.m)
}

func (m *string16Map) unmarshal(k string, d *json.Decoder) error {
	v, err := m.c.unmarshalFlags16Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(k, v)

	return nil
}

func (m *string16Map) postProcess(p jparser.Pair) error {
	v, err := m.c.postProcessFlags16Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(p.K, v)

	return nil
}

type string32Map struct {
	contentItemLink
	m *strtree32.Tree
}

func (m *string32Map) get() interface{} {
	return pdp.MakeContentStringFlags32Map(m.m)
}

func (m *string32Map) unmarshal(k string, d *json.Decoder) error {
	v, err := m.c.unmarshalFlags32Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(k, v)

	return nil
}

func (m *string32Map) postProcess(p jparser.Pair) error {
	v, err := m.c.postProcessFlags32Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(p.K, v)

	return nil
}

type string64Map struct {
	contentItemLink
	m *strtree64.Tree
}

func (m *string64Map) get() interface{} {
	return pdp.MakeContentStringFlags64Map(m.m)
}

func (m *string64Map) unmarshal(k string, d *json.Decoder) error {
	v, err := m.c.unmarshalFlags64Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(k, v)

	return nil
}

func (m *string64Map) postProcess(p jparser.Pair) error {
	v, err := m.c.postProcessFlags64Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(p.K, v)

	return nil
}

type networkMap struct {
	contentItemLink
	m *iptree.Tree
}

func (m *networkMap) unmarshal(k string, d *json.Decoder) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(k)
	if a == nil {
		_, n, err = net.ParseCIDR(k)
		if err != nil {
			return newAddressNetworkCastError(k, err)
		}
	}

	v, err := m.c.unmarshalTypedData(d, m.i+1)
	if err != nil {
		return bindError(err, k)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *networkMap) postProcess(p jparser.Pair) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(p.K)
	if a == nil {
		_, n, err = net.ParseCIDR(p.K)
		if err != nil {
			return newAddressNetworkCastError(p.K, err)
		}
	}

	v, err := m.c.postProcess(p.V, m.i+1)
	if err != nil {
		return bindError(err, p.K)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *networkMap) get() interface{} {
	return pdp.MakeContentNetworkMap(m.m)
}

type network8Map struct {
	contentItemLink
	m *iptree8.Tree
}

func (m *network8Map) unmarshal(k string, d *json.Decoder) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(k)
	if a == nil {
		_, n, err = net.ParseCIDR(k)
		if err != nil {
			return newAddressNetworkCastError(k, err)
		}
	}

	v, err := m.c.unmarshalFlags8Value(d)
	if err != nil {
		return bindError(err, k)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network8Map) postProcess(p jparser.Pair) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(p.K)
	if a == nil {
		_, n, err = net.ParseCIDR(p.K)
		if err != nil {
			return newAddressNetworkCastError(p.K, err)
		}
	}

	v, err := m.c.postProcessFlags8Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network8Map) get() interface{} {
	return pdp.MakeContentNetworkFlags8Map(m.m)
}

type network16Map struct {
	contentItemLink
	m *iptree16.Tree
}

func (m *network16Map) unmarshal(k string, d *json.Decoder) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(k)
	if a == nil {
		_, n, err = net.ParseCIDR(k)
		if err != nil {
			return newAddressNetworkCastError(k, err)
		}
	}

	v, err := m.c.unmarshalFlags16Value(d)
	if err != nil {
		return bindError(err, k)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network16Map) postProcess(p jparser.Pair) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(p.K)
	if a == nil {
		_, n, err = net.ParseCIDR(p.K)
		if err != nil {
			return newAddressNetworkCastError(p.K, err)
		}
	}

	v, err := m.c.postProcessFlags16Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network16Map) get() interface{} {
	return pdp.MakeContentNetworkFlags16Map(m.m)
}

type network32Map struct {
	contentItemLink
	m *iptree32.Tree
}

func (m *network32Map) unmarshal(k string, d *json.Decoder) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(k)
	if a == nil {
		_, n, err = net.ParseCIDR(k)
		if err != nil {
			return newAddressNetworkCastError(k, err)
		}
	}

	v, err := m.c.unmarshalFlags32Value(d)
	if err != nil {
		return bindError(err, k)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network32Map) postProcess(p jparser.Pair) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(p.K)
	if a == nil {
		_, n, err = net.ParseCIDR(p.K)
		if err != nil {
			return newAddressNetworkCastError(p.K, err)
		}
	}

	v, err := m.c.postProcessFlags32Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network32Map) get() interface{} {
	return pdp.MakeContentNetworkFlags32Map(m.m)
}

type network64Map struct {
	contentItemLink
	m *iptree64.Tree
}

func (m *network64Map) unmarshal(k string, d *json.Decoder) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(k)
	if a == nil {
		_, n, err = net.ParseCIDR(k)
		if err != nil {
			return newAddressNetworkCastError(k, err)
		}
	}

	v, err := m.c.unmarshalFlags64Value(d)
	if err != nil {
		return bindError(err, k)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network64Map) postProcess(p jparser.Pair) error {
	var (
		n   *net.IPNet
		err error
	)

	a := net.ParseIP(p.K)
	if a == nil {
		_, n, err = net.ParseCIDR(p.K)
		if err != nil {
			return newAddressNetworkCastError(p.K, err)
		}
	}

	v, err := m.c.postProcessFlags64Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	if a != nil {
		m.m.InplaceInsertIP(a, v)
	} else {
		m.m.InplaceInsertNet(n, v)
	}

	return nil
}

func (m *network64Map) get() interface{} {
	return pdp.MakeContentNetworkFlags64Map(m.m)
}

type domainMap struct {
	contentItemLink
	m *domaintree.Node
}

func (m *domainMap) unmarshal(k string, d *json.Decoder) error {
	dn, err := domain.MakeNameFromString(k)
	if err != nil {
		return newDomainCastError(k, err)
	}

	v, err := m.c.unmarshalTypedData(d, m.i+1)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domainMap) postProcess(p jparser.Pair) error {
	dn, err := domain.MakeNameFromString(p.K)
	if err != nil {
		return newDomainCastError(p.K, err)
	}

	v, err := m.c.postProcess(p.V, m.i+1)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domainMap) get() interface{} {
	return pdp.MakeContentDomainMap(m.m)
}

type domain8Map struct {
	contentItemLink
	m *domaintree8.Node
}

func (m *domain8Map) unmarshal(k string, d *json.Decoder) error {
	dn, err := domain.MakeNameFromString(k)
	if err != nil {
		return newDomainCastError(k, err)
	}

	v, err := m.c.unmarshalFlags8Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain8Map) postProcess(p jparser.Pair) error {
	dn, err := domain.MakeNameFromString(p.K)
	if err != nil {
		return newDomainCastError(p.K, err)
	}

	v, err := m.c.postProcessFlags8Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain8Map) get() interface{} {
	return pdp.MakeContentDomainFlags8Map(m.m)
}

type domain16Map struct {
	contentItemLink
	m *domaintree16.Node
}

func (m *domain16Map) unmarshal(k string, d *json.Decoder) error {
	dn, err := domain.MakeNameFromString(k)
	if err != nil {
		return newDomainCastError(k, err)
	}

	v, err := m.c.unmarshalFlags16Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain16Map) postProcess(p jparser.Pair) error {
	dn, err := domain.MakeNameFromString(p.K)
	if err != nil {
		return newDomainCastError(p.K, err)
	}

	v, err := m.c.postProcessFlags16Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain16Map) get() interface{} {
	return pdp.MakeContentDomainFlags16Map(m.m)
}

type domain32Map struct {
	contentItemLink
	m *domaintree32.Node
}

func (m *domain32Map) unmarshal(k string, d *json.Decoder) error {
	dn, err := domain.MakeNameFromString(k)
	if err != nil {
		return newDomainCastError(k, err)
	}

	v, err := m.c.unmarshalFlags32Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain32Map) postProcess(p jparser.Pair) error {
	dn, err := domain.MakeNameFromString(p.K)
	if err != nil {
		return newDomainCastError(p.K, err)
	}

	v, err := m.c.postProcessFlags32Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain32Map) get() interface{} {
	return pdp.MakeContentDomainFlags32Map(m.m)
}

type domain64Map struct {
	contentItemLink
	m *domaintree64.Node
}

func (m *domain64Map) unmarshal(k string, d *json.Decoder) error {
	dn, err := domain.MakeNameFromString(k)
	if err != nil {
		return newDomainCastError(k, err)
	}

	v, err := m.c.unmarshalFlags64Value(d)
	if err != nil {
		return bindError(err, k)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain64Map) postProcess(p jparser.Pair) error {
	dn, err := domain.MakeNameFromString(p.K)
	if err != nil {
		return newDomainCastError(p.K, err)
	}

	v, err := m.c.postProcessFlags64Value(p.V)
	if err != nil {
		return bindError(err, p.K)
	}

	m.m.InplaceInsert(dn, v)

	return nil
}

func (m *domain64Map) get() interface{} {
	return pdp.MakeContentDomainFlags64Map(m.m)
}

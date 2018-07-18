package pdp

import (
	"fmt"
	"math"
	"net"
	"reflect"
	"unsafe"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

var (
	reflectValueNil = reflect.ValueOf(nil)

	reflectTypeString     = reflect.TypeOf("")
	reflectTypeIP         = reflect.TypeOf(net.IP{})
	reflectTypeIPNet      = reflect.TypeOf(net.IPNet{})
	reflectTypePtrIPNet   = reflect.TypeOf((*net.IPNet)(nil))
	reflectTypeDomain     = reflect.TypeOf(domain.Name{})
	reflectTypeStrtree    = reflect.TypeOf((*strtree.Tree)(nil))
	reflectTypeIPTree     = reflect.TypeOf((*iptree.Tree)(nil))
	reflectTypeDomaintree = reflect.TypeOf((*domaintree.Node)(nil))
	reflectTypeStrings    = reflect.TypeOf([]string(nil))
)

func setEffect(v reflect.Value, effect int) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalEffectConstError(v)
	}

	switch v.Kind() {
	default:
		return newRequestUnmarshalEffectTypeError(v)

	case reflect.Bool:
		v.SetBool(effect == EffectPermit)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(effect))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(effect))

	case reflect.String:
		v.SetString(EffectNameFromEnum(effect))
	}

	return nil
}

func setStatus(v reflect.Value, s string) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalStatusConstError(v)
	}

	if v.Kind() == reflect.String {
		v.SetString(s)
		return nil
	}

	t := v.Type()
	if t.PkgPath() == "" && t.Name() == "error" {
		if len(s) > 0 {
			v.Set(reflect.ValueOf(newResponseServerError(s)))
		}

		return nil
	}

	return newRequestUnmarshalStatusTypeError(v)
}

func setBool(v reflect.Value, b bool) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalBooleanConstError(v)
	}

	if v.Kind() != reflect.Bool {
		return newRequestUnmarshalBooleanTypeError(v)
	}

	v.SetBool(b)
	return nil
}

func setString(v reflect.Value, s string) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalStringConstError(v)
	}

	if v.Kind() != reflect.String {
		return newRequestUnmarshalStringTypeError(v)
	}

	v.SetString(s)
	return nil
}

func setInt(v reflect.Value, i int64) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalIntegerConstError(v)
	}

	switch v.Kind() {
	default:
		return newRequestUnmarshalIntegerTypeError(v)

	case reflect.Int:
		if i < math.MinInt32 || i > math.MaxInt32 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetInt(i)

	case reflect.Int8:
		if i < math.MinInt8 || i > math.MaxInt8 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetInt(i)

	case reflect.Int16:
		if i < math.MinInt16 || i > math.MaxInt16 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetInt(i)

	case reflect.Int32:
		if i < math.MinInt32 || i > math.MaxInt32 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetInt(i)

	case reflect.Int64:
		v.SetInt(i)

	case reflect.Uint:
		if i < 0 {
			return newRequestUnmarshalIntegerUnderflowError(i, v)
		}

		if i > math.MaxUint32 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetUint(uint64(i))

	case reflect.Uint8:
		if i < 0 {
			return newRequestUnmarshalIntegerUnderflowError(i, v)
		}

		if i > math.MaxUint8 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetUint(uint64(i))

	case reflect.Uint16:
		if i < 0 {
			return newRequestUnmarshalIntegerUnderflowError(i, v)
		}

		if i > math.MaxUint16 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetUint(uint64(i))

	case reflect.Uint32:
		if i < 0 {
			return newRequestUnmarshalIntegerUnderflowError(i, v)
		}

		if i > math.MaxUint32 {
			return newRequestUnmarshalIntegerOverflowError(i, v)
		}

		v.SetUint(uint64(i))

	case reflect.Uint64:
		if i < 0 {
			return newRequestUnmarshalIntegerUnderflowError(i, v)
		}

		v.SetUint(uint64(i))
	}

	return nil
}

func setFloat(v reflect.Value, f float64) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalFloatConstError(v)
	}

	switch v.Kind() {
	default:
		return newRequestUnmarshalFloatTypeError(v)

	case reflect.Float32, reflect.Float64:
		v.SetFloat(f)
	}

	return nil
}

func setAddress(v reflect.Value, a net.IP) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalAddressConstError(v)
	}

	if v.Type() != reflectTypeIP {
		return newRequestUnmarshalAddressTypeError(v)
	}

	v.Set(reflect.ValueOf(a))
	return nil
}

func getNetwork(v reflect.Value) *net.IPNet {
	if v == reflectValueNil {
		return nil
	}

	t := v.Type()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Type() != reflectTypeIPNet {
		panic(fmt.Errorf("can't marshal %s as network value", t))
	}

	return &net.IPNet{
		IP:   net.IP(v.FieldByName("IP").Bytes()),
		Mask: net.IPMask(v.FieldByName("Mask").Bytes()),
	}
}

func setNetwork(v reflect.Value, n *net.IPNet) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalNetworkConstError(v)
	}

	switch v.Type() {
	default:
		return newRequestUnmarshalNetworkTypeError(v)

	case reflectTypeIPNet:
		v.Set(reflect.ValueOf(*n))

	case reflectTypePtrIPNet:
		v.Set(reflect.ValueOf(n))
	}

	return nil
}

func setDomain(v reflect.Value, d domain.Name) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalDomainConstError(v)
	}

	switch v.Type() {
	default:
		return newRequestUnmarshalDomainTypeError(v)

	case reflectTypeString:
		v.SetString(d.String())

	case reflectTypeDomain:
		v.Set(reflect.ValueOf(d))
	}

	return nil
}

func getSetOfStrings(v reflect.Value) *strtree.Tree {
	if v == reflectValueNil {
		return nil
	}

	t := v.Type()
	if t == reflectTypeStrtree {
		return (*strtree.Tree)(unsafe.Pointer(v.Pointer()))
	}

	panic(fmt.Errorf("can't marshal %s as set of strings value", t))
}

func setSetOfStrings(v reflect.Value, ss *strtree.Tree) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalSetOfStringsConstError(v)
	}

	if v.Type() != reflectTypeStrtree {
		return newRequestUnmarshalSetOfStringsTypeError(v)
	}

	v.Set(reflect.ValueOf(ss))
	return nil
}

func getSetOfNetworks(v reflect.Value) *iptree.Tree {
	if v == reflectValueNil {
		return nil
	}

	t := v.Type()
	if t == reflectTypeIPTree {
		return (*iptree.Tree)(unsafe.Pointer(v.Pointer()))
	}

	panic(fmt.Errorf("can't marshal %s as set of networks value", t))
}

func setSetOfNetworks(v reflect.Value, sn *iptree.Tree) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalSetOfNetworksConstError(v)
	}

	if v.Type() != reflectTypeIPTree {
		return newRequestUnmarshalSetOfNetworksTypeError(v)
	}

	v.Set(reflect.ValueOf(sn))
	return nil
}

func getSetOfDomains(v reflect.Value) *domaintree.Node {
	if v == reflectValueNil {
		return nil
	}

	t := v.Type()
	if t == reflectTypeDomaintree {
		return (*domaintree.Node)(unsafe.Pointer(v.Pointer()))
	}

	panic(fmt.Errorf("can't marshal %s as set of domains value", t))
}

func setSetOfDomains(v reflect.Value, sd *domaintree.Node) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalSetOfDomainsConstError(v)
	}

	if v.Type() != reflectTypeDomaintree {
		return newRequestUnmarshalSetOfDomainsTypeError(v)
	}

	v.Set(reflect.ValueOf(sd))
	return nil
}

func getListOfStrings(v reflect.Value) []string {
	if v == reflectValueNil {
		return nil
	}

	t := v.Type()
	if t == reflectTypeStrings {
		out := make([]string, v.Len())
		for i := range out {
			out[i] = v.Index(i).String()
		}

		return out
	}

	panic(fmt.Errorf("can't marshal %s as list of strings value", t))
}

func setListOfStrings(v reflect.Value, ls []string) error {
	if v == reflectValueNil {
		return nil
	}

	if !v.CanSet() {
		return newRequestUnmarshalListOfStringsConstError(v)
	}

	if v.Type() != reflectTypeStrings {
		return newRequestUnmarshalListOfStringsTypeError(v)
	}

	v.Set(reflect.ValueOf(ls))
	return nil
}

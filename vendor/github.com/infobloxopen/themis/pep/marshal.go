package pep

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

var (
	// ErrorInvalidSource indicates that input value of validate method is not
	// a structure.
	ErrorInvalidSource = errors.New("given value is not a structure")
	// ErrorInvalidSlice indicates that input structure has slice field
	// (client can't marshal slices).
	ErrorInvalidSlice = errors.New("marshalling for the slice hasn't been implemented")
	// ErrorInvalidStruct indicates that input structure has struct field
	// (client can't marshal nested structures).
	ErrorInvalidStruct = errors.New("marshalling for the struct hasn't been implemented")
	// ErrorIntegerOverflow indicates that input structure contains integer
	// which doesn't fit to int64.
	ErrorIntegerOverflow = errors.New("integer overflow")
)

var (
	boolType        = reflect.TypeOf(false)
	stringType      = reflect.TypeOf("")
	intType         = reflect.TypeOf(0)
	int8Type        = reflect.TypeOf(int8(0))
	int16Type       = reflect.TypeOf(int16(0))
	int32Type       = reflect.TypeOf(int32(0))
	int64Type       = reflect.TypeOf(int64(0))
	uintType        = reflect.TypeOf(uint(0))
	uint8Type       = reflect.TypeOf(uint8(0))
	uint16Type      = reflect.TypeOf(uint16(0))
	uint32Type      = reflect.TypeOf(uint32(0))
	uint64Type      = reflect.TypeOf(uint64(0))
	float32Type     = reflect.TypeOf(float32(0))
	float64Type     = reflect.TypeOf(float64(0))
	netIPType       = reflect.TypeOf(net.IP{})
	netIPNetType    = reflect.TypeOf(net.IPNet{})
	ptrNetIPNetType = reflect.TypeOf((*net.IPNet)(nil))
	domainType      = reflect.TypeOf(domain.Name{})
	strtreeType     = reflect.TypeOf((*strtree.Tree)(nil))
	iptreeType      = reflect.TypeOf((*iptree.Tree)(nil))
	domaintreeType  = reflect.TypeOf((*domaintree.Node)(nil))
	stringsType     = reflect.TypeOf([]string(nil))

	attrTypeByType = map[reflect.Type]pdp.Type{
		boolType:        pdp.TypeBoolean,
		stringType:      pdp.TypeString,
		intType:         pdp.TypeInteger,
		int8Type:        pdp.TypeInteger,
		int16Type:       pdp.TypeInteger,
		int32Type:       pdp.TypeInteger,
		int64Type:       pdp.TypeInteger,
		uintType:        pdp.TypeInteger,
		uint8Type:       pdp.TypeInteger,
		uint16Type:      pdp.TypeInteger,
		uint32Type:      pdp.TypeInteger,
		uint64Type:      pdp.TypeInteger,
		float32Type:     pdp.TypeFloat,
		float64Type:     pdp.TypeFloat,
		netIPType:       pdp.TypeAddress,
		netIPNetType:    pdp.TypeNetwork,
		ptrNetIPNetType: pdp.TypeNetwork,
		domainType:      pdp.TypeDomain,
		strtreeType:     pdp.TypeSetOfStrings,
		iptreeType:      pdp.TypeSetOfNetworks,
		domaintreeType:  pdp.TypeSetOfDomains,
		stringsType:     pdp.TypeListOfStrings,
	}

	attrTypeByTag = map[string]pdp.Type{
		pdp.TypeBoolean.GetKey():       pdp.TypeBoolean,
		pdp.TypeString.GetKey():        pdp.TypeString,
		pdp.TypeInteger.GetKey():       pdp.TypeInteger,
		pdp.TypeFloat.GetKey():         pdp.TypeFloat,
		pdp.TypeAddress.GetKey():       pdp.TypeAddress,
		pdp.TypeNetwork.GetKey():       pdp.TypeNetwork,
		pdp.TypeDomain.GetKey():        pdp.TypeDomain,
		pdp.TypeSetOfStrings.GetKey():  pdp.TypeSetOfStrings,
		pdp.TypeSetOfNetworks.GetKey(): pdp.TypeSetOfNetworks,
		pdp.TypeSetOfDomains.GetKey():  pdp.TypeSetOfDomains,
		pdp.TypeListOfStrings.GetKey(): pdp.TypeListOfStrings,
	}

	typeByAttrType = map[pdp.Type]map[reflect.Type]struct{}{
		pdp.TypeBoolean: {
			boolType: {},
		},
		pdp.TypeString: {
			stringType: {},
		},
		pdp.TypeInteger: {
			intType:    {},
			int8Type:   {},
			int16Type:  {},
			int32Type:  {},
			int64Type:  {},
			uintType:   {},
			uint8Type:  {},
			uint16Type: {},
			uint32Type: {},
			uint64Type: {},
		},
		pdp.TypeFloat: {
			float32Type: {},
			float64Type: {},
		},
		pdp.TypeAddress: {
			netIPType: {},
		},
		pdp.TypeNetwork: {
			netIPNetType:    {},
			ptrNetIPNetType: {},
		},
		pdp.TypeDomain: {
			stringType: {},
			domainType: {},
		},
		pdp.TypeSetOfStrings: {
			strtreeType: {},
		},
		pdp.TypeSetOfNetworks: {
			iptreeType: {},
		},
		pdp.TypeSetOfDomains: {
			domaintreeType: {},
		},
		pdp.TypeListOfStrings: {
			stringsType: {},
		},
	}

	typeByTag = map[string]map[reflect.Type]struct{}{}
)

type reqFieldInfo struct {
	idx int
	tag string
	at  pdp.Type
}

type reqFieldsInfo struct {
	fields []reqFieldInfo
	err    error
}

func init() {
	for k, v := range typeByAttrType {
		typeByTag[k.GetKey()] = v
	}
}

func makeTaggedFieldsInfo(fields []reflect.StructField, typeName string) reqFieldsInfo {
	var out []reqFieldInfo
	for i, f := range fields {
		tag, ok := getTag(f)
		if !ok {
			continue
		}

		var at pdp.Type
		items := strings.Split(tag, ",")
		if len(items) > 1 {
			tag = items[0]
			t := items[1]

			at, ok = attrTypeByTag[strings.ToLower(t)]
			if !ok {
				return makeReqsFieldsInfoErr("unknown type %q (%s.%s)", t, typeName, f.Name)
			}

			if _, ok := typeByAttrType[at][f.Type]; !ok {
				return makeReqsFieldsInfoErr("can't marshal %q as %q (%s.%s)", f.Type, t, typeName, f.Name)
			}

		} else {
			at, ok = attrTypeByType[f.Type]
			if !ok {
				return makeReqsFieldsInfoErr("can't marshal %q (%s.%s)", f.Type, typeName, f.Name)
			}
		}

		if len(tag) <= 0 {
			tag, ok = getName(f)
			if !ok {
				continue
			}
		}

		out = append(out, reqFieldInfo{i, tag, at})
	}

	return reqFieldsInfo{fields: out}
}

func makeUntaggedFieldsInfo(fields []reflect.StructField) reqFieldsInfo {
	var out []reqFieldInfo
	for i, f := range fields {
		name, ok := getName(f)
		if !ok {
			continue
		}

		t, ok := attrTypeByType[f.Type]
		if !ok {
			continue
		}

		out = append(out, reqFieldInfo{i, name, t})
	}

	return reqFieldsInfo{fields: out}
}

var (
	typeCache     = map[string]reqFieldsInfo{}
	typeCacheLock = sync.RWMutex{}
)

func makeRequest(v interface{}, b []byte) (pb.Msg, error) {
	switch v := v.(type) {
	case []byte:
		return pb.Msg{Body: v}, nil

	case pb.Msg:
		return v, nil

	case *pb.Msg:
		return *v, nil
	}

	var (
		n   int
		err error
	)

	if a, ok := v.([]pdp.AttributeAssignment); ok {
		n, err = pdp.MarshalRequestAssignments(b, a)
	} else {
		n, err = marshalValue(reflect.ValueOf(v), b)
	}
	if err != nil {
		return pb.Msg{}, err
	}

	return pb.Msg{Body: b[:n]}, nil
}

func marshalValue(v reflect.Value, b []byte) (int, error) {
	if v.Kind() != reflect.Struct {
		return 0, ErrorInvalidSource
	}

	return marshalStruct(v, getFields(v.Type()), b)
}

func getFields(t reflect.Type) reqFieldsInfo {
	key := t.PkgPath() + "." + t.Name()
	typeCacheLock.RLock()
	if info, ok := typeCache[key]; ok {
		typeCacheLock.RUnlock()
		return info
	}
	typeCacheLock.RUnlock()

	fields := make([]reflect.StructField, 0)
	tagged := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		_, ok := getTag(f)
		tagged = tagged || ok

		fields = append(fields, f)
	}

	typeCacheLock.Lock()
	var info reqFieldsInfo
	if tagged {
		info = makeTaggedFieldsInfo(fields, t.Name())
	} else {
		info = makeUntaggedFieldsInfo(fields)
	}
	typeCache[key] = info
	typeCacheLock.Unlock()

	return info
}

func getName(f reflect.StructField) (string, bool) {
	name := f.Name
	if len(name) <= 0 {
		return "", false
	}

	c := name[:1]
	if c != strings.ToUpper(c) {
		return "", false
	}

	return name, true
}

func getTag(f reflect.StructField) (string, bool) {
	if f.Tag == "pdp" {
		return "", true
	}

	return f.Tag.Lookup("pdp")
}

func marshalStruct(v reflect.Value, info reqFieldsInfo, b []byte) (int, error) {
	if info.err != nil {
		return 0, info.err
	}

	return pdp.MarshalRequestReflection(b, len(info.fields), func(i int) (string, pdp.Type, reflect.Value, error) {
		f := info.fields[i]
		return f.tag, f.at, v.Field(f.idx), nil
	})
}

func makeReqsFieldsInfoErr(s string, args ...interface{}) reqFieldsInfo {
	return reqFieldsInfo{
		err: fmt.Errorf(s, args...),
	}
}

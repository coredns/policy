// Package requests provides loader for YAML formatted authorization requests file.
package requests

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/infobloxopen/themis/pdp-service && protoc -I $GOPATH/src/github.com/infobloxopen/themis/proto/ $GOPATH/src/github.com/infobloxopen/themis/proto/service.proto --go_out=plugins=grpc:$GOPATH/src/github.com/infobloxopen/themis/pdp-service && ls $GOPATH/src/github.com/infobloxopen/themis/pdp-service"

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

type requests struct {
	Attributes map[string]string
	Requests   []map[string]interface{}
}

// Load reads given YAML file and porduces list of requests to run.
func Load(name string, size uint32) ([]pb.Msg, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	in := &requests{}
	err = yaml.Unmarshal(b, in)
	if err != nil {
		return nil, err
	}

	symbols := make(map[string]pdp.Type)
	for k, v := range in.Attributes {
		t, ok := pdp.BuiltinTypes[strings.ToLower(v)]
		if !ok {
			return nil, fmt.Errorf("unknown type %q of %q attribute", v, k)
		}

		symbols[k] = t
	}

	out := make([]pb.Msg, len(in.Requests))
	for i, r := range in.Requests {
		attrs := make([]pdp.AttributeAssignment, len(r))
		j := 0
		for k, v := range r {
			a, err := makeAttribute(k, v, symbols)
			if err != nil {
				return nil, fmt.Errorf("invalid attribute in request %d: %s", i+1, err)
			}

			attrs[j] = a
			j++
		}

		b := make([]byte, 10240)
		n, err := pdp.MarshalRequestAssignments(b, attrs)
		if err != nil {
			return nil, fmt.Errorf("can't create request: %s", err)
		}

		out[i] = pb.Msg{Body: b[:n]}
	}

	return out, nil
}

type attributeMarshaller func(value interface{}) (pdp.AttributeValue, error)

var marshallers = map[pdp.Type]attributeMarshaller{
	pdp.TypeBoolean: booleanMarshaller,
	pdp.TypeString:  stringMarshaller,
	pdp.TypeInteger: integerMarshaller,
	pdp.TypeFloat:   floatMarshaller,
	pdp.TypeAddress: addressMarshaller,
	pdp.TypeNetwork: networkMarshaller,
	pdp.TypeDomain:  domainMarshaller}

func makeAttribute(name string, value interface{}, symbols map[string]pdp.Type) (pdp.AttributeAssignment, error) {
	t, ok := symbols[name]
	var err error
	if !ok {
		t, err = guessType(value)
		if err != nil {
			return pdp.AttributeAssignment{},
				fmt.Errorf("type of %q attribute isn't defined and can't be derived: %s", name, err)
		}
	}

	marshaller, ok := marshallers[t]
	if !ok {
		return pdp.AttributeAssignment{},
			fmt.Errorf("marshaling hasn't been implemented for type %q of %q attribute", t, name)
	}

	v, err := marshaller(value)
	if err != nil {
		return pdp.AttributeAssignment{},
			fmt.Errorf("can't marshal %q attribute as %q: %s", name, t, err)
	}

	return pdp.MakeExpressionAssignment(name, v), nil
}

func guessType(value interface{}) (pdp.Type, error) {
	switch value.(type) {
	case bool:
		return pdp.TypeBoolean, nil
	case string:
		return pdp.TypeString, nil
	case net.IP:
		return pdp.TypeAddress, nil
	case net.IPNet:
		return pdp.TypeNetwork, nil
	case *net.IPNet:
		return pdp.TypeNetwork, nil
	}

	return pdp.TypeUndefined, fmt.Errorf("marshaling hasn't been implemented for %T", value)
}

func booleanMarshaller(value interface{}) (pdp.AttributeValue, error) {
	switch value := value.(type) {
	case bool:
		return pdp.MakeBooleanValue(value), nil
	case string:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return pdp.UndefinedValue, fmt.Errorf("can't marshal \"%s\" as boolean", value)
		}

		return pdp.MakeBooleanValue(b), nil
	}

	return pdp.UndefinedValue, fmt.Errorf("can't marshal %T as boolean", value)
}

func stringMarshaller(value interface{}) (pdp.AttributeValue, error) {
	s, ok := value.(string)
	if !ok {
		return pdp.UndefinedValue, fmt.Errorf("can't marshal %T as string", value)
	}

	return pdp.MakeStringValue(s), nil
}

func integerMarshaller(value interface{}) (pdp.AttributeValue, error) {
	var (
		i   int64
		err error
	)

	switch value := value.(type) {
	case int:
		return pdp.MakeIntegerValue(int64(value)), nil
	case int64:
		return pdp.MakeIntegerValue(value), nil
	case uint:
		if value <= math.MaxInt64 {
			return pdp.MakeIntegerValue(int64(value)), nil
		}
		err = fmt.Errorf("can't marshal %T (%d) as int64", value, value)

	case uint64:
		if value <= math.MaxInt64 {
			return pdp.MakeIntegerValue(int64(value)), nil
		}
		err = fmt.Errorf("can't marshal %T (%d) as int64", value, value)

	case float64:
		if value > -9007199254740992 && value < 9007199254740992 {
			return pdp.MakeIntegerValue(int64(value)), nil
		}
		err = fmt.Errorf("can't marshal %T (%g) as int64", value, value)

	case string:
		i, err = strconv.ParseInt(value, 10, 64)
		if err == nil {
			return pdp.MakeIntegerValue(i), nil
		}
		err = fmt.Errorf("can't marshal \"%s\" as int64", value)
	}

	return pdp.UndefinedValue, err
}

func floatMarshaller(value interface{}) (pdp.AttributeValue, error) {
	var (
		f   float64
		err error
	)

	switch value := value.(type) {
	case int:
		return pdp.MakeFloatValue(float64(value)), nil
	case int64:
		return pdp.MakeFloatValue(float64(value)), nil
	case uint:
		return pdp.MakeFloatValue(float64(value)), nil
	case uint64:
		return pdp.MakeFloatValue(float64(value)), nil
	case float64:
		return pdp.MakeFloatValue(float64(value)), nil
	case string:
		f, err = strconv.ParseFloat(value, 64)
		if err == nil {
			return pdp.MakeFloatValue(f), nil
		}
		err = fmt.Errorf("can't marshal \"%s\" as float64", value)
	}

	return pdp.UndefinedValue, err
}

func addressMarshaller(value interface{}) (pdp.AttributeValue, error) {
	switch value := value.(type) {
	case net.IP:
		return pdp.MakeAddressValue(value), nil
	case string:
		addr := net.ParseIP(value)
		if addr == nil {
			return pdp.UndefinedValue, fmt.Errorf("can't marshal \"%s\" as IP address", value)
		}

		return pdp.MakeAddressValue(addr), nil
	}

	return pdp.UndefinedValue, fmt.Errorf("can't marshal %T as IP address", value)
}

func networkMarshaller(value interface{}) (pdp.AttributeValue, error) {
	switch value := value.(type) {
	case net.IPNet:
		return pdp.MakeNetworkValue(&value), nil
	case *net.IPNet:
		return pdp.MakeNetworkValue(value), nil
	case string:
		_, n, err := net.ParseCIDR(value)
		if err != nil {
			return pdp.UndefinedValue, fmt.Errorf("can't marshal \"%s\" as network", value)
		}

		return pdp.MakeNetworkValue(n), nil
	}

	return pdp.UndefinedValue, fmt.Errorf("can't marshal %T as network", value)
}

func domainMarshaller(value interface{}) (pdp.AttributeValue, error) {
	s, ok := value.(string)
	if !ok {
		return pdp.UndefinedValue, fmt.Errorf("can't marshal %T as domain", value)
	}

	d, err := domain.MakeNameFromString(s)
	if err != nil {
		return pdp.UndefinedValue, fmt.Errorf("can't marshal %q as domain: %s", s, err)
	}

	return pdp.MakeDomainValue(d), nil
}

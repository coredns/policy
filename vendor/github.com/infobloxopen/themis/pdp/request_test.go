package pdp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"strconv"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

var (
	testWireRequest = append([]byte{1, 0}, testWireAttributes...)

	testRequestAssignments = []AttributeAssignment{
		MakeStringAssignment("string", "test"),
		MakeBooleanAssignment("boolean", true),
		MakeIntegerAssignment("integer", 9223372036854775807),
	}
)

func TestRequestWireTypesTotal(t *testing.T) {
	if requestWireTypesTotal != len(requestWireTypeNames) {
		t.Errorf("Expected number of wire type names %d to be equal to total number of wire types %d",
			len(requestWireTypeNames), requestWireTypesTotal)
	}
}

func TestMarshalRequestAssignments(t *testing.T) {
	b, err := MarshalRequestAssignments(testRequestAssignments)
	assertRequestBytesBuffer(t, "MarshalRequestAssignments", err, b, len(b), testWireRequest...)

	b, err = MarshalRequestAssignments([]AttributeAssignment{
		MakeExpressionAssignment("test", UndefinedValue),
	})
	if err == nil {
		t.Errorf("expected requestAttributeMarshallingNotImplementedError but got %d bytes in request", len(b))
	} else if _, ok := err.(*requestAttributeMarshallingNotImplementedError); !ok {
		t.Errorf("expected *requestAttributeMarshallingNotImplementedError but got %T (%s)", err, err)
	}
}

func TestMarshalRequestAssignmentsWithAllocator(t *testing.T) {
	b, err := MarshalRequestAssignmentsWithAllocator(testRequestAssignments, func(n int) ([]byte, error) {
		return make([]byte, n), nil
	})
	assertRequestBytesBuffer(t, "MarshalRequestAssignments", err, b, len(b), testWireRequest...)

	testFuncErr := errors.New("test function error")
	b, err = MarshalRequestAssignmentsWithAllocator(testRequestAssignments, func(n int) ([]byte, error) {
		return nil, testFuncErr
	})
	if err == nil {
		t.Errorf("expected testFuncErr but got %d bytes in request", len(b))
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}

	b, err = MarshalRequestAssignmentsWithAllocator([]AttributeAssignment{
		MakeExpressionAssignment("test", UndefinedValue),
	}, func(n int) ([]byte, error) {
		return make([]byte, n), nil
	})
	if err == nil {
		t.Errorf("expected requestAttributeMarshallingNotImplementedError but got %d bytes in request", len(b))
	} else if _, ok := err.(*requestAttributeMarshallingNotImplementedError); !ok {
		t.Errorf("expected *requestAttributeMarshallingNotImplementedError but got %T (%s)", err, err)
	}

	b, err = MarshalRequestAssignmentsWithAllocator(testRequestAssignments, func(n int) ([]byte, error) {
		return make([]byte, 2), nil
	})
	assertRequestBufferOverflow(t, "MarshalRequestAssignmentsWithAllocator", err, len(b))
}

func TestMarshalRequestAssignmentsToBuffer(t *testing.T) {
	var b [44]byte
	n, err := MarshalRequestAssignmentsToBuffer(b[:], testRequestAssignments)
	assertRequestBytesBuffer(t, "MarshalRequestAssignmentsToBuffer", err, b[:], n, testWireRequest...)

	n, err = MarshalRequestAssignmentsToBuffer([]byte{}, testRequestAssignments)
	assertRequestBufferOverflow(t, "MarshalRequestAssignmentsToBuffer(count)", err, n)

	n, err = MarshalRequestAssignmentsToBuffer(b[:2], testRequestAssignments)
	assertRequestBufferOverflow(t, "MarshalRequestAssignmentsToBuffer(first value)", err, n)
}

func TestMarshalRequestReflection(t *testing.T) {
	b, err := MarshalRequestReflection(1, func(i int) (string, Type, reflect.Value, error) {
		return "boolean", TypeBoolean, reflect.ValueOf(true), nil
	})
	assertRequestBytesBuffer(t, "MarshalRequestReflectionToBuffer", err, b, len(b),
		1, 0, 1, 0,
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
	)

	testFuncErr := errors.New("test function error")
	b, err = MarshalRequestReflection(1, func(i int) (string, Type, reflect.Value, error) {
		return "", TypeUndefined, reflectValueNil, testFuncErr
	})
	if err == nil {
		t.Errorf("expected testFuncErr but got %d bytes in request", len(b))
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}
}

func TestMarshalRequestReflectionWithAllocator(t *testing.T) {
	b, err := MarshalRequestReflectionWithAllocator(1, func(i int) (string, Type, reflect.Value, error) {
		return "boolean", TypeBoolean, reflect.ValueOf(true), nil
	}, func(n int) ([]byte, error) {
		return make([]byte, n), nil
	})
	assertRequestBytesBuffer(t, "MarshalRequestReflectionToBuffer", err, b, len(b),
		1, 0, 1, 0,
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
	)

	testFuncErr := errors.New("test function error")
	b, err = MarshalRequestReflectionWithAllocator(1, func(i int) (string, Type, reflect.Value, error) {
		return "", TypeUndefined, reflectValueNil, testFuncErr
	}, func(n int) ([]byte, error) {
		return make([]byte, n), nil
	})
	if err == nil {
		t.Errorf("expected testFuncErr but got %d bytes in request", len(b))
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}

	b, err = MarshalRequestReflectionWithAllocator(1, func(i int) (string, Type, reflect.Value, error) {
		return "boolean", TypeBoolean, reflect.ValueOf(true), nil
	}, func(n int) ([]byte, error) {
		return nil, testFuncErr
	})
	if err == nil {
		t.Errorf("expected testFuncErr but got %d bytes in request", len(b))
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}

	b, err = MarshalRequestReflectionWithAllocator(1, func(i int) (string, Type, reflect.Value, error) {
		return "boolean", TypeBoolean, reflect.ValueOf(true), nil
	}, func(n int) ([]byte, error) {
		return make([]byte, 2), nil
	})
	assertRequestBufferOverflow(t, "MarshalRequestReflectionWithAllocator", err, len(b))
}

func TestMarshalRequestReflectionToBuffer(t *testing.T) {
	var b [13]byte

	f := func(i int) (string, Type, reflect.Value, error) {
		return "boolean", TypeBoolean, reflect.ValueOf(true), nil
	}
	n, err := MarshalRequestReflectionToBuffer(b[:], 1, f)
	assertRequestBytesBuffer(t, "MarshalRequestReflectionToBuffer", err, b[:], n,
		1, 0, 1, 0,
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
	)

	n, err = MarshalRequestReflectionToBuffer([]byte{}, 1, f)
	assertRequestBufferOverflow(t, "MarshalRequestReflectionToBuffer(version)", err, n)

	n, err = MarshalRequestReflectionToBuffer(b[:2], 1, f)
	assertRequestBufferOverflow(t, "MarshalRequestReflectionToBuffer(collection)", err, n)
}

func TestUnmarshalRequestAssignments(t *testing.T) {
	a, err := UnmarshalRequestAssignments(testWireRequest)
	assertRequestAssignmentExpressions(t, "UnmarshalRequestAssignments", err, a, len(a), testRequestAssignments...)

	a, err = UnmarshalRequestAssignments([]byte{0, 0, 0, 0})
	if err == nil {
		t.Errorf("expected *requestVersionError but got %d attributes", len(a))
	} else if _, ok := err.(*requestVersionError); !ok {
		t.Errorf("expected *requestVersionError but got %T (%s)", err, err)
	}
}

func TestUnmarshalRequestAssignmentsWithAllocator(t *testing.T) {
	a, err := UnmarshalRequestAssignmentsWithAllocator(testWireRequest, func(n int) ([]AttributeAssignment, error) {
		return make([]AttributeAssignment, n), nil
	})
	assertRequestAssignmentExpressions(t, "UnmarshalRequestAssignmentsWithAllocator", err, a, len(a),
		testRequestAssignments...)

	a, err = UnmarshalRequestAssignmentsWithAllocator([]byte{0, 0, 0, 0}, func(n int) ([]AttributeAssignment, error) {
		return make([]AttributeAssignment, n), nil
	})
	if err == nil {
		t.Errorf("expected *requestVersionError but got %d attributes", len(a))
	} else if _, ok := err.(*requestVersionError); !ok {
		t.Errorf("expected *requestVersionError but got %T (%s)", err, err)
	}
}

func TestUnmarshalRequestToAssignmentsArray(t *testing.T) {
	var a [3]AttributeAssignment

	n, err := UnmarshalRequestToAssignmentsArray(testWireRequest, a[:])
	assertRequestAssignmentExpressions(t, "UnmarshalRequestToAssignmentsArray", err, a[:], n, testRequestAssignments...)

	n, err = UnmarshalRequestToAssignmentsArray([]byte{}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %d attributes", n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	n, err = UnmarshalRequestToAssignmentsArray([]byte{1, 0}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %d attributes", n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestUnmarshalRequestReflection(t *testing.T) {
	var (
		names [3]string
	)

	i := 0

	var (
		str     string
		boolean bool
		num     int64
	)

	types := []Type{
		TypeString,
		TypeBoolean,
		TypeInteger,
	}

	values := []reflect.Value{
		reflect.Indirect(reflect.ValueOf(&str)),
		reflect.Indirect(reflect.ValueOf(&boolean)),
		reflect.Indirect(reflect.ValueOf(&num)),
	}

	err := UnmarshalRequestReflection(testWireRequest, func(id string, t Type) (reflect.Value, error) {
		if i >= len(names) || i >= len(values) || i >= len(types) {
			return reflect.ValueOf(nil), fmt.Errorf("requested invalid value number: %d", i)
		}

		if et := types[i]; t != et {
			return reflect.ValueOf(nil), fmt.Errorf("expected %q for %d but got %q", et, i, t)
		}

		names[i] = id
		v := values[i]
		i++

		return v, nil
	})

	a := []AttributeAssignment{
		MakeStringAssignment(names[0], str),
		MakeBooleanAssignment(names[1], boolean),
		MakeIntegerAssignment(names[2], num),
	}

	assertRequestAssignmentExpressions(t, "UnmarshalRequestReflection", err, a, i, testRequestAssignments...)

	err = UnmarshalRequestReflection([]byte{}, func(id string, t Type) (reflect.Value, error) {
		return reflect.ValueOf(nil), fmt.Errorf("in unreacheable place with id %q and type %q", id, t)
	})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError but got nothing")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = UnmarshalRequestReflection([]byte{
		1, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return reflect.ValueOf(nil), fmt.Errorf("in unreacheable place with id %q and type %q", id, t)
	})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError but got nothing")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestCheckRequestVersion(t *testing.T) {
	n, err := checkRequestVersion(testWireRequest)
	if err != nil {
		t.Error(err)
	} else if n <= 0 {
		t.Errorf("expected some bytes consumed but got %d", n)
	} else if n > len(testWireRequest) {
		t.Errorf("not expected more bytes consumed (%d) than buffer has (%d)", n, len(testWireRequest))
	}

	_, err = checkRequestVersion([]byte{})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError but got nothing")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	_, err = checkRequestVersion([]byte{2, 0})
	if err == nil {
		t.Error("expected *requestVersionError but got nothing")
	} else if _, ok := err.(*requestVersionError); !ok {
		t.Errorf("expected *requestVersionError but got %T (%s)", err, err)
	}
}

func TestGetRequestAttributeCount(t *testing.T) {
	off, err := checkRequestVersion(testWireRequest)
	if err != nil {
		t.Fatal(err)
	}

	c, n, err := getRequestAttributeCount(testWireRequest[off:])
	if err != nil {
		t.Error(err)
	} else if n <= 0 {
		t.Errorf("expected some bytes consumed but got %d", n)
	} else if n > len(testWireRequest) {
		t.Errorf("not expected more bytes consumed (%d) than buffer has (%d)", n, len(testWireRequest))
	} else if c != 3 {
		t.Errorf("expected %d as attribute count but got %d", 1, c)
	}

	c, _, err = getRequestAttributeCount([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got count %d", c)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestAttributeName(t *testing.T) {
	off, err := checkRequestVersion(testWireRequest)
	if err != nil {
		t.Fatal(err)
	}

	c, n, err := getRequestAttributeCount(testWireRequest[off:])
	if err != nil {
		t.Fatal(err)
	} else if c != 3 {
		t.Fatalf("expected %d as attribute count but got %d", 1, c)
	}

	off += n

	name, n, err := getRequestAttributeName(testWireRequest[off:])
	if err != nil {
		t.Error(err)
	} else if n <= 0 {
		t.Errorf("expected some bytes consumed but got %d", n)
	} else if n > len(testWireRequest) {
		t.Errorf("not expected more bytes consumed (%d) than buffer has (%d)", n, len(testWireRequest))
	} else if name != "string" {
		t.Errorf("expected %q as attribute name but got %q", "test", name)
	}

	name, _, err = getRequestAttributeName([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got name %q", name)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, _, err = getRequestAttributeName([]byte{4, 't', 'e'})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got name %q", name)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestAttributeType(t *testing.T) {
	off, err := checkRequestVersion(testWireRequest)
	if err != nil {
		t.Fatal(err)
	}

	c, n, err := getRequestAttributeCount(testWireRequest[off:])
	if err != nil {
		t.Fatal(err)
	} else if c != 3 {
		t.Fatalf("expected %d as attribute count but got %d", 1, c)
	}

	off += n

	name, n, err := getRequestAttributeName(testWireRequest[off:])
	if err != nil {
		t.Fatal(err)
	} else if name != "string" {
		t.Fatalf("expected %q as attribute name but got %q", "test", name)
	}

	off += n

	at, n, err := getRequestAttributeType(testWireRequest[off:])
	if err != nil {
		t.Error(err)
	} else if n <= 0 {
		t.Errorf("expected some bytes consumed but got %d", n)
	} else if n > len(testWireRequest) {
		t.Errorf("not expected more bytes consumed (%d) than buffer has (%d)", n, len(testWireRequest))
	} else if at != requestWireTypeString {
		tn := "unknown"
		if at >= 0 || at < len(requestWireTypeNames) {
			tn = requestWireTypeNames[at]
		}

		t.Errorf("expected %q (%d) as attribute type but got %q (%d)",
			requestWireTypeNames[requestWireTypeString], requestWireTypeString, tn, at)
	}

	at, _, err = getRequestAttributeType([]byte{})
	if err == nil {
		tn := "unknown"
		if at >= 0 || at < len(requestWireTypeNames) {
			tn = requestWireTypeNames[at]
		}

		t.Errorf("expected *requestBufferUnderflowError but got type %q (%d)", tn, at)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestStringValue(t *testing.T) {
	off, err := checkRequestVersion(testWireRequest)
	if err != nil {
		t.Fatal(err)
	}

	c, n, err := getRequestAttributeCount(testWireRequest[off:])
	if err != nil {
		t.Fatal(err)
	} else if c != 3 {
		t.Fatalf("expected %d as attribute count but got %d", 1, c)
	}

	off += n

	name, n, err := getRequestAttributeName(testWireRequest[off:])
	if err != nil {
		t.Fatal(err)
	} else if name != "string" {
		t.Fatalf("expected %q as attribute name but got %q", "test", name)
	}

	off += n

	at, n, err := getRequestAttributeType(testWireRequest[off:])
	if err != nil {
		t.Fatal(err)
	} else if at != requestWireTypeString {
		tn := "unknown"
		if at >= 0 || at < len(requestWireTypeNames) {
			tn = requestWireTypeNames[at]
		}

		t.Fatalf("expected %q (%d) as attribute type but got %q (%d)",
			requestWireTypeNames[requestWireTypeString], requestWireTypeString, tn, at)
	}

	off += n

	v, n, err := getRequestStringValue(testWireRequest[off:])
	if err != nil {
		t.Error(err)
	} else if n <= 0 {
		t.Errorf("expected some bytes consumed but got %d", n)
	} else if off+n > len(testWireRequest) {
		t.Errorf("not expected more bytes consumed (%d) than buffer has (%d)", n, len(testWireRequest))
	} else if v != "test" {
		t.Errorf("expected string %q as attribute value but got %q", "test value", v)
	}

	v, _, err = getRequestStringValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got string %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestStringValue([]byte{10, 0, 't', 'e', 's', 't'})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got string %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestIntegerValue(t *testing.T) {
	testWireIntegerValue := []byte{
		0, 0, 0, 0, 0, 0, 0, 128,
	}
	v, n, err := getRequestIntegerValue(testWireIntegerValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIntegerValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIntegerValue), n)
	} else if v != -9223372036854775808 {
		t.Errorf("expected integer %d as attribute value but got %d", -9223372036854775808, v)
	}

	v, _, err = getRequestIntegerValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got integer %d", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestFloatValue(t *testing.T) {
	testWireFloatValue := []byte{
		24, 45, 68, 84, 251, 33, 9, 64,
	}
	v, n, err := getRequestFloatValue(testWireFloatValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireFloatValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireFloatValue), n)
	} else if v != float64(math.Pi) {
		t.Errorf("expected float %g as attribute value but got %g", float64(math.Pi), v)
	}

	v, _, err = getRequestFloatValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got float %g", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestIPv4AddressValue(t *testing.T) {
	testWireIPv4AddressValue := []byte{
		192, 0, 2, 1,
	}
	v, n, err := getRequestIPv4AddressValue(testWireIPv4AddressValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv4AddressValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv4AddressValue), n)
	} else if !v.Equal(net.ParseIP("192.0.2.1")) {
		t.Errorf("expected IPv4 address %q as attribute value but got %q", net.ParseIP("192.0.2.1"), v)
	}

	v, _, err = getRequestIPv4AddressValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got IPv4 address %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestIPv6AddressValue(t *testing.T) {
	testWireIPv6AddressValue := []byte{
		32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}
	v, n, err := getRequestIPv6AddressValue(testWireIPv6AddressValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv6AddressValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv6AddressValue), n)
	} else if !v.Equal(net.ParseIP("2001:db8::1")) {
		t.Errorf("expected IPv6 address %q as attribute value but got %q", net.ParseIP("2001:db8::1"), v)
	}

	v, _, err = getRequestIPv6AddressValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got IPv6 address %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestIPv4NetworkValue(t *testing.T) {
	testWireIPv4NetworkValue := []byte{
		24, 192, 0, 2, 1,
	}
	v, n, err := getRequestIPv4NetworkValue(testWireIPv4NetworkValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv4NetworkValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv4NetworkValue), n)
	} else if v.String() != "192.0.2.0/24" {
		t.Errorf("expected IPv4 network %q as attribute value but got %q", "192.0.2.0/24", v)
	}

	v, _, err = getRequestIPv4NetworkValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got IPv4 network %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestIPv4NetworkValue([]byte{
		255, 192, 0, 2, 1,
	})
	if err == nil {
		t.Errorf("expected *requestIPv4InvalidMaskError but got IPv4 network %q", v)
	} else if _, ok := err.(*requestIPv4InvalidMaskError); !ok {
		t.Errorf("expected *requestIPv4InvalidMaskError but got %T (%s)", err, err)
	}
}

func TestGetRequestIPv6NetworkValue(t *testing.T) {
	testWireIPv6NetworkValue := []byte{
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}
	v, n, err := getRequestIPv6NetworkValue(testWireIPv6NetworkValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv6NetworkValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv6NetworkValue), n)
	} else if v.String() != "2001:db8::/32" {
		t.Errorf("expected IPv6 network %q as attribute value but got %q", "2001:db8::/32", v)
	}

	v, _, err = getRequestIPv6NetworkValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got IPv6 network %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestIPv6NetworkValue([]byte{
		255, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	})
	if err == nil {
		t.Errorf("expected *requestIPv6InvalidMaskError but got IPv6 network %q", v)
	} else if _, ok := err.(*requestIPv6InvalidMaskError); !ok {
		t.Errorf("expected *requestIPv6InvalidMaskError but got %T (%s)", err, err)
	}
}

func TestGetRequestDomainValue(t *testing.T) {
	testWireDomainValue := []byte{
		8, 0, 't', 'e', 's', 't', '.', 'c', 'o', 'm',
	}
	v, n, err := getRequestDomainValue(testWireDomainValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireDomainValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireDomainValue), n)
	} else if v.String() != "test.com" {
		t.Errorf("expected domain %q as attribute value but got %q", "test.com", v)
	}

	v, _, err = getRequestDomainValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got domain %q", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestDomainValue([]byte{
		3, 0, '.', '.', '.',
	})
	if err == nil {
		t.Errorf("expected domain.ErrEmptyLabel error but got domain %q", v)
	} else if err != domain.ErrEmptyLabel {
		t.Errorf("expected domain.ErrEmptyLabel error but got %T (%s)", err, err)
	}
}

func TestGetRequestSetOfStringsValue(t *testing.T) {
	testWireSetOfStringsValue := []byte{
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	}

	v, n, err := getRequestSetOfStringsValue(testWireSetOfStringsValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireSetOfStringsValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireSetOfStringsValue), n)
	} else {
		assertStrings(SortSetOfStrings(v), []string{"one", "two", "three"}, "getRequestSetOfStringsValue", t)
	}

	v, _, err = getRequestSetOfStringsValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got set of strings %#v", SortSetOfStrings(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestSetOfStringsValue([]byte{
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got set of strings %#v", SortSetOfStrings(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestSetOfNetworksValue(t *testing.T) {
	testWireSetOfNetworksValue := []byte{
		3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
	}

	v, n, err := getRequestSetOfNetworksValue(testWireSetOfNetworksValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireSetOfNetworksValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireSetOfNetworksValue), n)
	} else {
		assertNetworks(SortSetOfNetworks(v), []*net.IPNet{
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		}, "getRequestSetOfNetworksValue", t)
	}

	v, _, err = getRequestSetOfNetworksValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got set of networks %#v", SortSetOfNetworks(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestSetOfNetworksValue([]byte{
		3, 0,
		216, 192, 0, 2, 0,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got set of networks %#v", SortSetOfNetworks(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestSetOfNetworksValue([]byte{
		3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got set of networks %#v", SortSetOfNetworks(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestSetOfNetworksValue([]byte{
		1, 0,
		225, 192, 0, 2, 0,
	})
	if err == nil {
		t.Errorf("expected *requestIPv4InvalidMaskError but got set of networks %#v", SortSetOfNetworks(v))
	} else if _, ok := err.(*requestIPv4InvalidMaskError); !ok {
		t.Errorf("expected *requestIPv4InvalidMaskError but got %T (%s)", err, err)
	}

	v, _, err = getRequestSetOfNetworksValue([]byte{
		1, 0,
		129, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})
	if err == nil {
		t.Errorf("expected *requestIPv6InvalidMaskError but got set of networks %#v", SortSetOfNetworks(v))
	} else if _, ok := err.(*requestIPv6InvalidMaskError); !ok {
		t.Errorf("expected *requestIPv6InvalidMaskError but got %T (%s)", err, err)
	}
}

func TestGetRequestSetOfDomainsValue(t *testing.T) {
	testWireSetOfDomainsValue := []byte{
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	}

	v, n, err := getRequestSetOfDomainsValue(testWireSetOfDomainsValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireSetOfDomainsValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireSetOfDomainsValue), n)
	} else {
		assertStrings(SortSetOfDomains(v), []string{
			"example.com",
			"example.gov",
			"www.example.com",
		}, "getRequestSetOfDomainsValue", t)
	}

	v, _, err = getRequestSetOfDomainsValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got domains %#v", SortSetOfDomains(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestSetOfDomainsValue([]byte{
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got domains %#v", SortSetOfDomains(v))
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestListOfStringsValue(t *testing.T) {
	testWireListOfStringsValue := []byte{
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	}

	v, n, err := getRequestListOfStringsValue(testWireListOfStringsValue)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireListOfStringsValue) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireListOfStringsValue), n)
	} else {
		assertStrings(v, []string{"one", "two", "three"}, "getRequestListOfStringsValue", t)
	}

	v, _, err = getRequestListOfStringsValue([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got list of strings %#v", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	v, _, err = getRequestListOfStringsValue([]byte{
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got list of strings %#v", v)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetRequestAttribute(t *testing.T) {
	testWireBooleanFalseAttribute := []byte{
		2, 'n', 'o', byte(requestWireTypeBooleanFalse),
	}

	name, v, n, err := getRequestAttribute(testWireBooleanFalseAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireBooleanFalseAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireBooleanFalseAttribute), n)
	} else if name != "no" {
		t.Errorf("expected %q as attribute name but got %q", "no", name)
	} else if vt := v.GetResultType(); vt != TypeBoolean {
		t.Errorf("expected value of %q type but got %q %s", TypeBoolean, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "false"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireBooleanTrueAttribute := []byte{
		3, 'y', 'e', 's', byte(requestWireTypeBooleanTrue),
	}

	name, v, n, err = getRequestAttribute(testWireBooleanTrueAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireBooleanTrueAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireBooleanTrueAttribute), n)
	} else if name != "yes" {
		t.Errorf("expected %q as attribute name but got %q", "yes", name)
	} else if vt := v.GetResultType(); vt != TypeBoolean {
		t.Errorf("expected value of %q type but got %q %s", TypeBoolean, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "true"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireStringAttribute := []byte{
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
	}

	name, v, n, err = getRequestAttribute(testWireStringAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireStringAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireStringAttribute), n)
	} else if name != "string" {
		t.Errorf("expected %q as attribute name but got %q", "string", name)
	} else if vt := v.GetResultType(); vt != TypeString {
		t.Errorf("expected value of %q type but got %q %s", TypeString, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "test"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireIntegerAttribute := []byte{
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 0, 0, 0, 0, 0, 0, 0, 128,
	}

	name, v, n, err = getRequestAttribute(testWireIntegerAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIntegerAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIntegerAttribute), n)
	} else if name != "integer" {
		t.Errorf("expected %q as attribute name but got %q", "integer", name)
	} else if vt := v.GetResultType(); vt != TypeInteger {
		t.Errorf("expected value of %q type but got %q %s", TypeInteger, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "-9223372036854775808"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireFloatAttribute := []byte{
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84, 251, 33, 9, 64,
	}

	name, v, n, err = getRequestAttribute(testWireFloatAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireFloatAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireFloatAttribute), n)
	} else if name != "float" {
		t.Errorf("expected %q as attribute name but got %q", "float", name)
	} else if vt := v.GetResultType(); vt != TypeFloat {
		t.Errorf("expected value of %q type but got %q %s", TypeFloat, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "3.141592653589793"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireIPv4AddressAttribute := []byte{
		4, 'I', 'P', 'v', '4', byte(requestWireTypeIPv4Address), 192, 0, 2, 1,
	}

	name, v, n, err = getRequestAttribute(testWireIPv4AddressAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv4AddressAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv4AddressAttribute), n)
	} else if name != "IPv4" {
		t.Errorf("expected %q as attribute name but got %q", "IPv4", name)
	} else if vt := v.GetResultType(); vt != TypeAddress {
		t.Errorf("expected value of %q type but got %q %s", TypeAddress, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "192.0.2.1"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireIPv6AddressAttribute := []byte{
		4, 'I', 'P', 'v', '6', byte(requestWireTypeIPv6Address), 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}

	name, v, n, err = getRequestAttribute(testWireIPv6AddressAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv6AddressAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv6AddressAttribute), n)
	} else if name != "IPv6" {
		t.Errorf("expected %q as attribute name but got %q", "IPv6", name)
	} else if vt := v.GetResultType(); vt != TypeAddress {
		t.Errorf("expected value of %q type but got %q %s", TypeAddress, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "2001:db8::1"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireIPv4NetworkAttribute := []byte{
		11, 'I', 'P', 'v', '4', 'N', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv4Network), 24, 192, 0, 2, 1,
	}

	name, v, n, err = getRequestAttribute(testWireIPv4NetworkAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv4NetworkAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv4NetworkAttribute), n)
	} else if name != "IPv4Network" {
		t.Errorf("expected %q as attribute name but got %q", "IPv4Network", name)
	} else if vt := v.GetResultType(); vt != TypeNetwork {
		t.Errorf("expected value of %q type but got %q %s", TypeNetwork, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "192.0.2.0/24"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireIPv6NetworkAttribute := []byte{
		11, 'I', 'P', 'v', '6', 'N', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv6Network),
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}

	name, v, n, err = getRequestAttribute(testWireIPv6NetworkAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireIPv6NetworkAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireIPv6NetworkAttribute), n)
	} else if name != "IPv6Network" {
		t.Errorf("expected %q as attribute name but got %q", "IPv6Network", name)
	} else if vt := v.GetResultType(); vt != TypeNetwork {
		t.Errorf("expected value of %q type but got %q %s", TypeNetwork, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "2001:db8::/32"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireDomainAttribute := []byte{
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain), 8, 0, 't', 'e', 's', 't', '.', 'c', 'o', 'm',
	}

	name, v, n, err = getRequestAttribute(testWireDomainAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireDomainAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireDomainAttribute), n)
	} else if name != "domain" {
		t.Errorf("expected %q as attribute name but got %q", "domain", name)
	} else if vt := v.GetResultType(); vt != TypeDomain {
		t.Errorf("expected value of %q type but got %q %s", TypeDomain, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "test.com"
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireSetOfStringsAttribute := []byte{
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings),
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	}
	name, v, n, err = getRequestAttribute(testWireSetOfStringsAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireSetOfStringsAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireSetOfStringsAttribute), n)
	} else if name != "set of strings" {
		t.Errorf("expected %q as attribute name but got %q", "set of strings", name)
	} else if vt := v.GetResultType(); vt != TypeSetOfStrings {
		t.Errorf("expected value of %q type but got %q %s", TypeSetOfStrings, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "\"one\",\"two\",\"three\""
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireSetOfNetworksAttribute := []byte{
		15, 's', 'e', 't', ' ', 'o', 'f', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', 's',
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
	}
	name, v, n, err = getRequestAttribute(testWireSetOfNetworksAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireSetOfNetworksAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireSetOfNetworksAttribute), n)
	} else if name != "set of networks" {
		t.Errorf("expected %q as attribute name but got %q", "set of networks", name)
	} else if vt := v.GetResultType(); vt != TypeSetOfNetworks {
		t.Errorf("expected value of %q type but got %q %s", TypeSetOfNetworks, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "\"192.0.2.0/24\",\"2001:db8::/32\",\"192.0.2.16/28\""
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireSetOfDomainsAttribute := []byte{
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains),
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	}
	name, v, n, err = getRequestAttribute(testWireSetOfDomainsAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireSetOfDomainsAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireSetOfDomainsAttribute), n)
	} else if name != "set of domains" {
		t.Errorf("expected %q as attribute name but got %q", "set of strings", name)
	} else if vt := v.GetResultType(); vt != TypeSetOfDomains {
		t.Errorf("expected value of %q type but got %q %s", TypeSetOfDomains, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "\"example.com\",\"example.gov\",\"www.example.com\""
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	testWireListOfStringsAttribute := []byte{
		15, 'l', 'i', 's', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's',
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	}
	name, v, n, err = getRequestAttribute(testWireListOfStringsAttribute)
	if err != nil {
		t.Error(err)
	} else if n != len(testWireListOfStringsAttribute) {
		t.Errorf("expected whole buffer consumed (%d) but got (%d)", len(testWireListOfStringsAttribute), n)
	} else if name != "list of strings" {
		t.Errorf("expected %q as attribute name but got %q", "list of strings", name)
	} else if vt := v.GetResultType(); vt != TypeListOfStrings {
		t.Errorf("expected value of %q type but got %q %s", TypeListOfStrings, vt, v.describe())
	} else {
		s, err := v.Serialize()
		if err != nil {
			t.Error(err)
		} else {
			e := "\"one\",\"two\",\"three\""
			if s != e {
				t.Errorf("expected %q but got %q", e, s)
			}
		}
	}

	name, v, _, err = getRequestAttribute([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		6, 'n', 'o', 't', 'y', 'p', 'e',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		7, 'u', 'n', 'k', 'n', 'o', 'w', 'n', 255,
	})
	if err == nil {
		t.Errorf("expected *requestAttributeUnmarshallingTypeError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestAttributeUnmarshallingTypeError); !ok {
		t.Errorf("expected *requestAttributeUnmarshallingTypeError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 0, 0, 0, 0,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		4, 'I', 'P', 'v', '4', byte(requestWireTypeIPv4Address), 192, 0,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		4, 'I', 'P', 'v', '6', byte(requestWireTypeIPv6Address), 32, 1, 13, 184, 0, 0, 0, 0,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		11, 'I', 'P', 'v', '4', 'N', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv4Network), 192, 0, 2, 1,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		11, 'I', 'P', 'v', '6', 'N', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv6Network),
		32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain), 8, 0, 't', 'e', 's', 't',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings),
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		15, 's', 'e', 't', ' ', 'o', 'f', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', 's',
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains),
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	name, v, _, err = getRequestAttribute([]byte{
		15, 'l', 'i', 's', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's',
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got attribute %q = %s", name, v.describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestPutRequestVersion(t *testing.T) {
	var b [2]byte

	n, err := putRequestVersion(b[:])
	assertRequestBytesBuffer(t, "putRequestVersion", err, b[:], n, 1, 0)

	n, err = putRequestVersion(nil)
	assertRequestBufferOverflow(t, "putRequestVersion", err, n)
}

func TestPutRequestAttributeCount(t *testing.T) {
	var b [2]byte

	n, err := putRequestAttributeCount(b[:], 0)
	assertRequestBytesBuffer(t, "putRequestAttributeCount", err, b[:], n, 0, 0)

	n, err = putRequestAttributeCount(nil, 0)
	assertRequestBufferOverflow(t, "putRequestAttributeCount", err, n)

	n, err = putRequestAttributeCount(b[:], -1)
	if err == nil {
		t.Errorf("expected *requestInvalidAttributeCountError for negative count but got %d bytes in buffer", n)
	} else if _, ok := err.(*requestInvalidAttributeCountError); !ok {
		t.Errorf("expected *requestInvalidAttributeCountError but got %T (%s)", err, err)
	}

	n, err = putRequestAttributeCount(b[:], math.MaxUint16+1)
	if err == nil {
		t.Errorf("expected *requestTooManyAttributesError for large count but got %d bytes in buffer", n)
	} else if _, ok := err.(*requestTooManyAttributesError); !ok {
		t.Errorf("expected *requestTooManyAttributesError but got %T (%s)", err, err)
	}
}

func TestPutRequestAttributeName(t *testing.T) {
	var b [5]byte

	n, err := putRequestAttributeName(b[:], "test")
	assertRequestBytesBuffer(t, "putRequestAttributeName", err, b[:], n, 4, 't', 'e', 's', 't')

	n, err = putRequestAttributeName(nil, "test")
	assertRequestBufferOverflow(t, "putRequestAttributeName", err, n)

	n, err = putRequestAttributeName(b[:],
		"01234567890123456789012345678901234567890123456789012345678901234567890123456789"+
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789"+
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789"+
			"0123456789012345",
	)
	if err == nil {
		t.Errorf("expected *requestTooLongAttributeNameError for long name but got %d bytes in buffer", n)
	} else if _, ok := err.(*requestTooLongAttributeNameError); !ok {
		t.Errorf("expected *requestTooLongAttributeNameError but got %T (%s)", err, err)
	}
}

func TestPutRequestAttributeType(t *testing.T) {
	var b [1]byte

	n, err := putRequestAttributeType(b[:], requestWireTypeString)
	assertRequestBytesBuffer(t, "putRequestAttributeType", err, b[:], n, byte(requestWireTypeString))

	n, err = putRequestAttributeType(nil, requestWireTypeString)
	assertRequestBufferOverflow(t, "putRequestAttributeType", err, n)

	n, err = putRequestAttributeType(b[:], 2147483647)
	if err == nil {
		t.Errorf("expected *requestAttributeMarshallingTypeError for long name but got %d bytes in buffer", n)
	} else if _, ok := err.(*requestAttributeMarshallingTypeError); !ok {
		t.Errorf("expected *requestAttributeMarshallingTypeError but got %T (%s)", err, err)
	}
}

func TestPutRequestAttribute(t *testing.T) {
	var b [61]byte

	n, err := putRequestAttribute(b[:9], "boolean", MakeBooleanValue(true))
	assertRequestBytesBuffer(t, "putRequestAttribute(boolean)", err, b[:9], n,
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
	)

	n, err = putRequestAttribute(b[:14], "string", MakeStringValue("test"))
	assertRequestBytesBuffer(t, "putRequestAttribute(string)", err, b[:14], n,
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
	)

	n, err = putRequestAttribute(b[:17], "integer", MakeIntegerValue(-9223372036854775808))
	assertRequestBytesBuffer(t, "putRequestAttribute(integer)", err, b[:17], n,
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 0, 0, 0, 0, 0, 0, 0, 0x80,
	)

	n, err = putRequestAttribute(b[:15], "float", MakeFloatValue(math.Pi))
	assertRequestBytesBuffer(t, "putRequestAttribute(float)", err, b[:15], n,
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84, 251, 33, 9, 64,
	)

	n, err = putRequestAttribute(b[:13], "address", MakeAddressValue(net.ParseIP("192.0.2.1")))
	assertRequestBytesBuffer(t, "putRequestAttribute(address)", err, b[:13], n,
		7, 'a', 'd', 'd', 'r', 'e', 's', 's', byte(requestWireTypeIPv4Address), 192, 0, 2, 1,
	)

	n, err = putRequestAttribute(b[:14], "network", MakeNetworkValue(makeTestNetwork("192.0.2.0/24")))
	assertRequestBytesBuffer(t, "putRequestAttribute(network)", err, b[:14], n,
		7, 'n', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv4Network), 24, 192, 0, 2, 0,
	)

	n, err = putRequestAttribute(b[:25], "domain", MakeDomainValue(makeTestDomain("www.example.com")))
	assertRequestBytesBuffer(t, "putRequestAttribute(domain)", err, b[:25], n,
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain),
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	)

	n, err = putRequestAttribute(b[:35], "set of strings", MakeSetOfStringsValue(newStrTree("one", "two", "three")))
	assertRequestBytesBuffer(t, "putRequestAttribute(set of strings)", err, b[:35], n,
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings),
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putRequestAttribute(b[:46], "set of networks", MakeSetOfNetworksValue(newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	)))
	assertRequestBytesBuffer(t, "putRequestAttribute(set of networks)", err, b[:46], n,
		15, 's', 'e', 't', ' ', 'o', 'f', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', 's',
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
	)

	n, err = putRequestAttribute(b[:61], "set of domains", MakeSetOfDomainsValue(newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	)))
	assertRequestBytesBuffer(t, "putRequestAttribute(set of domains)", err, b[:61], n,
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains),
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	)

	n, err = putRequestAttribute(b[:36], "list of strings", MakeListOfStringsValue([]string{"one", "two", "three"}))
	assertRequestBytesBuffer(t, "putRequestAttribute(list of strings)", err, b[:36], n,
		15, 'l', 'i', 's', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's',
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putRequestAttribute(b[:], "undefined", UndefinedValue)
	if err == nil {
		t.Errorf("expected no data put to buffer for undefined value but got %d", n)
	} else if _, ok := err.(*requestAttributeMarshallingNotImplementedError); !ok {
		t.Errorf("expected *requestAttributeMarshallingNotImplementedError but got %T (%s)", err, err)
	}
}

func TestPutRequestAttributeBoolean(t *testing.T) {
	var b [9]byte

	n, err := putRequestAttributeBoolean(b[:], "boolean", true)
	assertRequestBytesBuffer(t, "putRequestAttributeBoolean", err, b[:], n,
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
	)

	n, err = putRequestAttributeBoolean(b[:5], "boolean", true)
	assertRequestBufferOverflow(t, "putRequestAttributeBoolean(name)", err, n)

	n, err = putRequestAttributeBoolean(b[:8], "boolean", true)
	assertRequestBufferOverflow(t, "putRequestAttributeBoolean(value)", err, n)
}

func TestPutRequestAttributeString(t *testing.T) {
	var b [14]byte

	n, err := putRequestAttributeString(b[:], "string", "test")
	assertRequestBytesBuffer(t, "putRequestAttributeString", err, b[:], n,
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
	)

	n, err = putRequestAttributeString(b[:5], "string", "test")
	assertRequestBufferOverflow(t, "putRequestAttributeString(name)", err, n)

	n, err = putRequestAttributeString(b[:9], "string", "test")
	assertRequestBufferOverflow(t, "putRequestAttributeString(value)", err, n)
}

func TestPutRequestAttributeInteger(t *testing.T) {
	var b [17]byte

	n, err := putRequestAttributeInteger(b[:], "integer", -9223372036854775808)
	assertRequestBytesBuffer(t, "putRequestAttributeInteger", err, b[:], n,
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 0, 0, 0, 0, 0, 0, 0, 0x80,
	)

	n, err = putRequestAttributeInteger(b[:5], "integer", 0)
	assertRequestBufferOverflow(t, "putRequestAttributeInteger(name)", err, n)

	n, err = putRequestAttributeInteger(b[:9], "integer", 0)
	assertRequestBufferOverflow(t, "putRequestAttributeInteger(value)", err, n)
}

func TestPutRequestAttributeFloat(t *testing.T) {
	var b [15]byte

	n, err := putRequestAttributeFloat(b[:], "float", math.Pi)
	assertRequestBytesBuffer(t, "putRequestAttributeFloat", err, b[:], n,
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84, 251, 33, 9, 64,
	)

	n, err = putRequestAttributeFloat(b[:4], "float", 0)
	assertRequestBufferOverflow(t, "putRequestAttributeFloat(name)", err, n)

	n, err = putRequestAttributeFloat(b[:7], "float", 0)
	assertRequestBufferOverflow(t, "putRequestAttributeFloat(value)", err, n)
}

func TestPutRequestAttributeAddress(t *testing.T) {
	var b [13]byte

	n, err := putRequestAttributeAddress(b[:], "address", net.ParseIP("192.0.2.1"))
	assertRequestBytesBuffer(t, "putRequestAttributeFloat", err, b[:], n,
		7, 'a', 'd', 'd', 'r', 'e', 's', 's', byte(requestWireTypeIPv4Address), 192, 0, 2, 1,
	)

	n, err = putRequestAttributeAddress(b[:4], "address", net.ParseIP("192.0.2.1"))
	assertRequestBufferOverflow(t, "putRequestAttributeAddress(name)", err, n)

	n, err = putRequestAttributeAddress(b[:9], "address", net.ParseIP("192.0.2.1"))
	assertRequestBufferOverflow(t, "putRequestAttributeAddress(value)", err, n)
}

func TestPutRequestAttributeNetwork(t *testing.T) {
	var b [14]byte

	n, err := putRequestAttributeNetwork(b[:], "network", makeTestNetwork("192.0.2.0/24"))
	assertRequestBytesBuffer(t, "putRequestAttributeFloat", err, b[:], n,
		7, 'n', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv4Network), 24, 192, 0, 2, 0,
	)

	n, err = putRequestAttributeNetwork(b[:4], "network", makeTestNetwork("192.0.2.0/24"))
	assertRequestBufferOverflow(t, "putRequestAttributeNetwork(name)", err, n)

	n, err = putRequestAttributeNetwork(b[:10], "network", makeTestNetwork("192.0.2.0/24"))
	assertRequestBufferOverflow(t, "putRequestAttributeNetwork(value)", err, n)
}

func TestPutRequestAttributeDomain(t *testing.T) {
	var b [25]byte

	n, err := putRequestAttributeDomain(b[:], "domain", makeTestDomain("www.example.com"))
	assertRequestBytesBuffer(t, "putRequestAttributeFloat", err, b[:], n,
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain),
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	)

	n, err = putRequestAttributeDomain(b[:4], "domain", makeTestDomain("www.example.com"))
	assertRequestBufferOverflow(t, "putRequestAttributeDomain(name)", err, n)

	n, err = putRequestAttributeDomain(b[:10], "domain", makeTestDomain("www.example.com"))
	assertRequestBufferOverflow(t, "putRequestAttributeDomain(value)", err, n)
}

func TestPutRequestAttributeSetOfStrings(t *testing.T) {
	var b [28]byte

	n, err := putRequestAttributeSetOfStrings(b[:], "strings", newStrTree("one", "two", "three"))
	assertRequestBytesBuffer(t, "putRequestAttributeSetOfStrings", err, b[:], n,
		7, 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putRequestAttributeSetOfStrings(b[:5], "strings", newStrTree("one", "two", "three"))
	assertRequestBufferOverflow(t, "putRequestAttributeSetOfStrings(name)", err, n)

	n, err = putRequestAttributeSetOfStrings(b[:13], "strings", newStrTree("one", "two", "three"))
	assertRequestBufferOverflow(t, "putRequestAttributeSetOfStrings(value)", err, n)
}

func TestPutRequestAttributeSetOfNetworks(t *testing.T) {
	var b [39]byte

	n, err := putRequestAttributeSetOfNetworks(b[:], "networks", newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	assertRequestBytesBuffer(t, "putRequestAttributeSetOfNetworks", err, b[:], n,
		8, 'n', 'e', 't', 'w', 'o', 'r', 'k', 's', byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
	)

	n, err = putRequestAttributeSetOfNetworks(b[:5], "networks", newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	assertRequestBufferOverflow(t, "putRequestAttributeSetOfNetworks(name)", err, n)

	n, err = putRequestAttributeSetOfNetworks(b[:15], "networks", newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	assertRequestBufferOverflow(t, "putRequestAttributeSetOfNetworks(value)", err, n)
}

func TestPutRequestAttributeSetOfDomains(t *testing.T) {
	var b [54]byte

	n, err := putRequestAttributeSetOfDomains(b[:], "domains", newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	assertRequestBytesBuffer(t, "putRequestAttributeSetOfDomains", err, b[:], n,
		7, 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains), 3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	)

	n, err = putRequestAttributeSetOfDomains(b[:5], "domains", newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	assertRequestBufferOverflow(t, "putRequestAttributeSetOfDomains(name)", err, n)

	n, err = putRequestAttributeSetOfDomains(b[:13], "domains", newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	assertRequestBufferOverflow(t, "putRequestAttributeSetOfDomains(value)", err, n)
}

func TestPutRequestAttributeListOfStrings(t *testing.T) {
	var b [28]byte

	n, err := putRequestAttributeListOfStrings(b[:], "strings", []string{"one", "two", "three"})
	assertRequestBytesBuffer(t, "putRequestAttributeListOfStrings", err, b[:], n,
		7, 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putRequestAttributeListOfStrings(b[:5], "strings", []string{"one", "two", "three"})
	assertRequestBufferOverflow(t, "putRequestAttributeListOfStrings(name)", err, n)

	n, err = putRequestAttributeListOfStrings(b[:13], "strings", []string{"one", "two", "three"})
	assertRequestBufferOverflow(t, "putRequestAttributeListOfStrings(value)", err, n)
}

func TestPutRequestBooleanValue(t *testing.T) {
	var b [1]byte

	n, err := putRequestBooleanValue(b[:], true)
	assertRequestBytesBuffer(t, "putRequestBooleanValue(true)", err, b[:], n, byte(requestWireTypeBooleanTrue))

	n, err = putRequestBooleanValue(b[:], false)
	assertRequestBytesBuffer(t, "putRequestBooleanValue(false)", err, b[:], n, byte(requestWireTypeBooleanFalse))
}

func TestPutRequestStringValue(t *testing.T) {
	var b [7]byte

	n, err := putRequestStringValue(b[:], "test")
	assertRequestBytesBuffer(t, "putRequestStringValue", err, b[:], n,
		byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
	)

	n, err = putRequestStringValue(nil, "test")
	assertRequestBufferOverflow(t, "putRequestStringValue", err, n)

	n, err = putRequestStringValue(b[:], string(make([]byte, math.MaxUint16+1)))
	if err == nil {
		t.Errorf("expected no data put to buffer for large string value but got %d", n)
	} else if _, ok := err.(*requestTooLongStringValueError); !ok {
		t.Errorf("expected *requestTooLongStringValueError but got %T (%s)", err, err)
	}

	n, err = putRequestStringValue(b[:3], "test")
	assertRequestBufferOverflow(t, "putRequestStringValue(buffer)", err, n)
}

func TestPutRequestIntegerValue(t *testing.T) {
	var b [9]byte

	n, err := putRequestIntegerValue(b[:], -9223372036854775808)
	assertRequestBytesBuffer(t, "putRequestIntegerValue", err, b[:], n,
		byte(requestWireTypeInteger), 0, 0, 0, 0, 0, 0, 0, 0x80,
	)

	n, err = putRequestIntegerValue(nil, 0)
	assertRequestBufferOverflow(t, "putRequestIntegerValue", err, n)

	n, err = putRequestIntegerValue(b[:5], 0)
	assertRequestBufferOverflow(t, "putRequestIntegerValue(buffer)", err, n)
}

func TestPutRequestFloatValue(t *testing.T) {
	var b [9]byte

	n, err := putRequestFloatValue(b[:], math.Pi)
	assertRequestBytesBuffer(t, "putRequestFloatValue", err, b[:], n,
		byte(requestWireTypeFloat), 24, 45, 68, 84, 251, 33, 9, 64,
	)

	n, err = putRequestFloatValue(nil, 0)
	assertRequestBufferOverflow(t, "putRequestFloatValue", err, n)

	n, err = putRequestFloatValue(b[:5], 0)
	assertRequestBufferOverflow(t, "putRequestFloatValue(buffer)", err, n)
}

func TestPutRequestAddressValue(t *testing.T) {
	var b [17]byte

	n, err := putRequestAddressValue(b[:5], net.ParseIP("192.0.2.1"))
	assertRequestBytesBuffer(t, "putRequestAddressValue(IPv4)", err, b[:5], n,
		byte(requestWireTypeIPv4Address), 192, 0, 2, 1,
	)

	n, err = putRequestAddressValue(b[:17], net.ParseIP("2001:db8::1"))
	assertRequestBytesBuffer(t, "putRequestAddressValue(IPv6)", err, b[:17], n,
		byte(requestWireTypeIPv6Address), 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	)

	n, err = putRequestAddressValue(b[:], []byte{192, 0, 2, 1, 0})
	if err == nil {
		t.Errorf("expected no data put to buffer for invalid IP address but got %d", n)
	} else if _, ok := err.(*requestAddressValueError); !ok {
		t.Errorf("expected *requestAddressValueError but got %T (%s)", err, err)
	}

	n, err = putRequestAddressValue(nil, net.ParseIP("192.0.2.1"))
	assertRequestBufferOverflow(t, "putRequestAddressValue", err, n)

	n, err = putRequestAddressValue(b[:3], net.ParseIP("192.0.2.1"))
	assertRequestBufferOverflow(t, "putRequestAddressValue(buffer)", err, n)
}

func TestPutRequestNetworkValue(t *testing.T) {
	var b [18]byte

	n, err := putRequestNetworkValue(b[:6], makeTestNetwork("192.0.2.0/24"))
	assertRequestBytesBuffer(t, "putRequestNetworkValue(IPv4)", err, b[:6], n,
		byte(requestWireTypeIPv4Network), 24, 192, 0, 2, 0,
	)

	n, err = putRequestNetworkValue(b[:18], makeTestNetwork("2001:db8::/32"))
	assertRequestBytesBuffer(t, "putRequestNetworkValue(IPv6)", err, b[:18], n,
		byte(requestWireTypeIPv6Network), 32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	)

	n, err = putRequestNetworkValue(b[:], nil)
	if err == nil {
		t.Errorf("expected no data put to buffer for invalid IP network but got %d", n)
	} else if _, ok := err.(*requestInvalidNetworkValueError); !ok {
		t.Errorf("expected *requestInvalidNetworkValueError but got %T (%s)", err, err)
	}

	n, err = putRequestNetworkValue(b[:], &net.IPNet{IP: nil, Mask: net.CIDRMask(24, 32)})
	if err == nil {
		t.Errorf("expected no data put to buffer for invalid IP network but got %d", n)
	} else if _, ok := err.(*requestInvalidNetworkValueError); !ok {
		t.Errorf("expected *requestInvalidNetworkValueError but got %T (%s)", err, err)
	}

	n, err = putRequestNetworkValue(b[:], &net.IPNet{IP: net.ParseIP("192.0.2.1").To4(), Mask: nil})
	if err == nil {
		t.Errorf("expected no data put to buffer for invalid IP network but got %d", n)
	} else if _, ok := err.(*requestInvalidNetworkValueError); !ok {
		t.Errorf("expected *requestInvalidNetworkValueError but got %T (%s)", err, err)
	}

	n, err = putRequestNetworkValue(nil, makeTestNetwork("192.0.2.0/24"))
	assertRequestBufferOverflow(t, "putRequestNetworkValue", err, n)

	n, err = putRequestNetworkValue(b[:4], makeTestNetwork("192.0.2.0/24"))
	assertRequestBufferOverflow(t, "putRequestNetworkValue(buffer)", err, n)
}

func TestPutRequestDomainValue(t *testing.T) {
	var b [18]byte

	n, err := putRequestDomainValue(b[:], makeTestDomain("www.example.com"))
	assertRequestBytesBuffer(t, "putRequestDomainValue", err, b[:], n,
		byte(requestWireTypeDomain), 15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	)

	n, err = putRequestDomainValue(nil, makeTestDomain("www.example.com"))
	assertRequestBufferOverflow(t, "putRequestDomainValue", err, n)

	n, err = putRequestDomainValue(b[:7], makeTestDomain("www.example.com"))
	assertRequestBufferOverflow(t, "putRequestDomainValue(buffer)", err, n)
}

func TestPutRequestSetOfStringsValue(t *testing.T) {
	var b [20]byte

	n, err := putRequestSetOfStringsValue(b[:], newStrTree("one", "two", "three"))
	assertRequestBytesBuffer(t, "putRequestSetOfStringsValue", err, b[:], n,
		byte(requestWireTypeSetOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putRequestSetOfStringsValue(nil, newStrTree("one", "two", "three"))
	assertRequestBufferOverflow(t, "putRequestSetOfStringsValue", err, n)

	n, err = putRequestSetOfStringsValue(b[:10], newStrTree("one", "two", "three"))
	assertRequestBufferOverflow(t, "putRequestSetOfStringsValue(buffer)", err, n)

	ss := strtree.NewTree()
	for i := 0; i < math.MaxUint16+1; i++ {
		ss.InplaceInsert(strconv.Itoa(i), i)
	}

	n, err = putRequestSetOfStringsValue(b[:], ss)
	if err == nil {
		t.Errorf("expected no data put from too big set but got %d", n)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}

	n, err = putRequestSetOfStringsValue(b[:], newStrTree("one", "two", string(make([]byte, math.MaxUint16+1))))
	if err == nil {
		t.Errorf("expected no data put with too big set of strings element but got %d", n)
	} else if _, ok := err.(*requestTooLongStringValueError); !ok {
		t.Errorf("expected *requestTooLongStringValueError but got %T (%s)", err, err)
	}
}

func TestPutRequestSetOfNetworksValue(t *testing.T) {
	var b [30]byte

	n, err := putRequestSetOfNetworksValue(b[:], newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	assertRequestBytesBuffer(t, "putRequestSetOfNetworksValue", err, b[:], n,
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
	)

	n, err = putRequestSetOfNetworksValue(nil, newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	assertRequestBufferOverflow(t, "putRequestSetOfNetworksValue", err, n)

	n, err = putRequestSetOfNetworksValue(b[:8], newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	assertRequestBufferOverflow(t, "putRequestSetOfNetworksValue", err, n)

	sn := iptree.NewTree()
	var ip [4]byte
	for i := 0; i < math.MaxUint16+1; i++ {
		binary.BigEndian.PutUint32(ip[:], uint32(i+1))
		ip[0] = 127
		sn.InplaceInsertIP(net.IP(ip[:]), i)
	}

	n, err = putRequestSetOfNetworksValue(b[:], sn)
	if err == nil {
		t.Errorf("expected no data put from too big set but got %d", n)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}
}

func TestPutRequestSetOfDomainsValue(t *testing.T) {
	var b [46]byte

	n, err := putRequestSetOfDomainsValue(b[:], newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	assertRequestBytesBuffer(t, "putRequestSetOfDomainsValue", err, b[:], n,
		byte(requestWireTypeSetOfDomains), 3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
	)

	n, err = putRequestSetOfDomainsValue(nil, newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	assertRequestBufferOverflow(t, "putRequestSetOfDomainsValue", err, n)

	n, err = putRequestSetOfDomainsValue(b[:29], newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	assertRequestBufferOverflow(t, "putRequestSetOfDomainsValue(buffer)", err, n)

	sd := new(domaintree.Node)
	for i := 0; i < math.MaxUint16+1; i++ {
		sd.InplaceInsert(makeTestDomain(strconv.Itoa(i)+".com"), i)
	}

	n, err = putRequestSetOfDomainsValue(b[:], sd)
	if err == nil {
		t.Errorf("expected no data put from too big set but got %d", n)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}
}

func TestPutRequestListOfStringsValue(t *testing.T) {
	var b [20]byte

	n, err := putRequestListOfStringsValue(b[:], []string{"one", "two", "three"})
	assertRequestBytesBuffer(t, "putRequestListOfStringsValue", err, b[:], n,
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putRequestListOfStringsValue(nil, []string{"one", "two", "three"})
	assertRequestBufferOverflow(t, "putRequestListOfStringsValue", err, n)

	n, err = putRequestListOfStringsValue(b[:10], []string{"one", "two", "three"})
	assertRequestBufferOverflow(t, "putRequestListOfStringsValue(buffer)", err, n)

	ls := make([]string, math.MaxUint16+1)
	for i := range ls {
		ls[i] = strconv.Itoa(i)
	}

	n, err = putRequestListOfStringsValue(b[:], ls)
	if err == nil {
		t.Errorf("expected no data put from too big list but got %d", n)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}

	n, err = putRequestListOfStringsValue(b[:], []string{"one", "two", string(make([]byte, math.MaxUint16+1))})
	if err == nil {
		t.Errorf("expected no data put with too big list of strings element but got %d", n)
	} else if _, ok := err.(*requestTooLongStringValueError); !ok {
		t.Errorf("expected *requestTooLongStringValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestSize(t *testing.T) {
	s, err := calcRequestSize(testRequestAssignments)
	if err != nil {
		t.Error(err)
	} else if s != len(testWireRequest) {
		t.Errorf("expected %d bytes in request but got %d", testWireRequest, s)
	}

	s, err = calcRequestSize([]AttributeAssignment{
		MakeExpressionAssignment("test", UndefinedValue),
	})
	if err == nil {
		t.Errorf("expected requestAttributeMarshallingNotImplementedError but got %d bytes in request", s)
	} else if _, ok := err.(*requestAttributeMarshallingNotImplementedError); !ok {
		t.Errorf("expected *requestAttributeMarshallingNotImplementedError but got %T (%s)", err, err)
	}
}

func TestCalcRequestSizeFromReflection(t *testing.T) {
	s, err := calcRequestSizeFromReflection(11, testReflectAttributes)
	if err != nil {
		t.Error(err)
	} else if s != len(testWireReflectAttributes)+reqVersionSize {
		t.Errorf("expected %d bytes in request but got %d", len(testWireReflectAttributes)+reqVersionSize, s)
	}

	testFuncErr := errors.New("test function error")
	s, err = calcRequestSizeFromReflection(1, func(i int) (string, Type, reflect.Value, error) {
		return "", TypeUndefined, reflectValueNil, testFuncErr
	})
	if err == nil {
		t.Errorf("expected testFuncErr but got %d bytes in request", s)
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeSize(t *testing.T) {
	s, err := calcRequestAttributeSize(MakeBooleanValue(true))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize, s)
	}

	s, err = calcRequestAttributeSize(MakeStringValue("test"))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+reqBigCounterSize+4 {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+reqBigCounterSize+4, s)
	}

	s, err = calcRequestAttributeSize(MakeIntegerValue(0))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+reqIntegerValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+reqIntegerValueSize, s)
	}

	s, err = calcRequestAttributeSize(MakeFloatValue(0))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+reqFloatValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+reqFloatValueSize, s)
	}

	s, err = calcRequestAttributeSize(MakeAddressValue(net.ParseIP("192.0.2.1")))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+reqIPv4AddressValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+reqIPv4AddressValueSize, s)
	}

	s, err = calcRequestAttributeSize(MakeNetworkValue(makeTestNetwork("192.0.2.0/24")))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+reqIPv4AddressValueSize+reqNetworkCIDRSize {
		t.Errorf("expected %d bytes for value but got %d",
			reqTypeSize+reqIPv4AddressValueSize+reqNetworkCIDRSize, s)
	}

	s, err = calcRequestAttributeSize(MakeDomainValue(makeTestDomain("example.com")))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+reqBigCounterSize+11 {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+reqBigCounterSize+11, s)
	}

	s, err = calcRequestAttributeSize(MakeSetOfStringsValue(newStrTree("one", "two", "three")))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+4*reqBigCounterSize+3+3+5 {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+4*reqBigCounterSize+3+3+5, s)
	}

	s, err = calcRequestAttributeSize(MakeSetOfNetworksValue(newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	)))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+
		reqBigCounterSize+
		2*(reqIPv4AddressValueSize+reqNetworkCIDRSize)+
		reqIPv6AddressValueSize+reqNetworkCIDRSize {
		t.Errorf("expected %d bytes for value but got %d",
			reqTypeSize+
				reqBigCounterSize+
				2*(reqIPv4AddressValueSize+reqNetworkCIDRSize)+
				reqIPv6AddressValueSize+reqNetworkCIDRSize,
			s)
	}

	s, err = calcRequestAttributeSize(MakeSetOfDomainsValue(newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	)))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+4*reqBigCounterSize+2*11+15 {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+4*reqBigCounterSize+2*11+15, s)
	}

	s, err = calcRequestAttributeSize(MakeListOfStringsValue([]string{"one", "two", "three"}))
	if err != nil {
		t.Error(err)
	} else if s != reqTypeSize+4*reqBigCounterSize+3+3+5 {
		t.Errorf("expected %d bytes for value but got %d", reqTypeSize+4*reqBigCounterSize+3+3+5, s)
	}

	s, err = calcRequestAttributeSize(UndefinedValue)
	if err == nil {
		t.Errorf("expected requestAttributeMarshallingNotImplementedError but got %d bytes in request", s)
	} else if _, ok := err.(*requestAttributeMarshallingNotImplementedError); !ok {
		t.Errorf("expected *requestAttributeMarshallingNotImplementedError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeNameSize(t *testing.T) {
	s, err := calcRequestAttributeNameSize("test")
	if err != nil {
		t.Error(err)
	} else if s != reqSmallCounterSize+4 {
		t.Errorf("expected %d bytes for value but got %d", reqSmallCounterSize+4, s)
	}

	s, err = calcRequestAttributeNameSize(
		"01234567890123456789012345678901234567890123456789012345678901234567890123456789" +
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789" +
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789" +
			"0123456789012345",
	)
	if err == nil {
		t.Errorf("expected *requestTooLongAttributeNameError for long name but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongAttributeNameError); !ok {
		t.Errorf("expected *requestTooLongAttributeNameError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeStringSize(t *testing.T) {
	s, err := calcRequestAttributeStringSize("test")
	if err != nil {
		t.Error(err)
	} else if s != reqBigCounterSize+4 {
		t.Errorf("expected %d bytes for value but got %d", reqBigCounterSize+4, s)
	}

	s, err = calcRequestAttributeStringSize(string(make([]byte, math.MaxUint16+1)))
	if err == nil {
		t.Errorf("expected *requestTooLongStringValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongStringValueError); !ok {
		t.Errorf("expected *requestTooLongStringValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeIntegerSize(t *testing.T) {
	s, err := calcRequestAttributeIntegerSize(0)
	if err != nil {
		t.Error(err)
	} else if s != reqIntegerValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqIntegerValueSize, s)
	}
}

func TestCalcRequestAttributeFloatSize(t *testing.T) {
	s, err := calcRequestAttributeFloatSize(0)
	if err != nil {
		t.Error(err)
	} else if s != reqFloatValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqFloatValueSize, s)
	}
}

func TestCalcRequestAttributeAddressSize(t *testing.T) {
	s, err := calcRequestAttributeAddressSize(net.ParseIP("192.0.2.1"))
	if err != nil {
		t.Error(err)
	} else if s != reqIPv4AddressValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqIPv4AddressValueSize, s)
	}

	s, err = calcRequestAttributeAddressSize(net.ParseIP("2001:db8::1"))
	if err != nil {
		t.Error(err)
	} else if s != reqIPv6AddressValueSize {
		t.Errorf("expected %d bytes for value but got %d", reqIPv6AddressValueSize, s)
	}

	s, err = calcRequestAttributeAddressSize(net.IP([]byte{0, 1, 2, 3, 4, 5, 6, 7}))
	if err == nil {
		t.Errorf("expected *requestAddressValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestAddressValueError); !ok {
		t.Errorf("expected *requestAddressValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeNetworkSize(t *testing.T) {
	s, err := calcRequestAttributeNetworkSize(makeTestNetwork("192.0.2.0/24"))
	if err != nil {
		t.Error(err)
	} else if s != reqIPv4AddressValueSize+reqNetworkCIDRSize {
		t.Errorf("expected %d bytes for value but got %d", reqIPv4AddressValueSize+reqNetworkCIDRSize, s)
	}

	s, err = calcRequestAttributeNetworkSize(makeTestNetwork("2001:db8::/32"))
	if err != nil {
		t.Error(err)
	} else if s != reqIPv6AddressValueSize+reqNetworkCIDRSize {
		t.Errorf("expected %d bytes for value but got %d", reqIPv6AddressValueSize+reqNetworkCIDRSize, s)
	}

	s, err = calcRequestAttributeNetworkSize(nil)
	if err == nil {
		t.Errorf("expected *requestInvalidNetworkValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestInvalidNetworkValueError); !ok {
		t.Errorf("expected *requestInvalidNetworkValueError but got %T (%s)", err, err)
	}

	s, err = calcRequestAttributeNetworkSize(&net.IPNet{
		IP:   net.IP([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	})
	if err == nil {
		t.Errorf("expected *requestInvalidNetworkValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestInvalidNetworkValueError); !ok {
		t.Errorf("expected *requestInvalidNetworkValueError but got %T (%s)", err, err)
	}

	s, err = calcRequestAttributeNetworkSize(&net.IPNet{
		IP:   net.ParseIP("192.0.2.0").To4(),
		Mask: net.IPMask([]byte{0, 1, 2, 3, 4, 5, 6, 7}),
	})
	if err == nil {
		t.Errorf("expected *requestInvalidNetworkValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestInvalidNetworkValueError); !ok {
		t.Errorf("expected *requestInvalidNetworkValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeDomainSize(t *testing.T) {
	s, err := calcRequestAttributeDomainSize(makeTestDomain("www.example.com"))
	if err != nil {
		t.Error(err)
	} else if s != reqBigCounterSize+15 {
		t.Errorf("expected %d bytes for value but got %d", reqBigCounterSize+15, s)
	}
}

func TestCalcRequestAttributeSetOfStringsSize(t *testing.T) {
	s, err := calcRequestAttributeSetOfStringsSize(newStrTree("one", "two", "three"))
	if err != nil {
		t.Error(err)
	} else if s != 4*reqBigCounterSize+3+3+5 {
		t.Errorf("expected %d bytes for value but got %d", 4*reqBigCounterSize+3+3+5, s)
	}

	ss := strtree.NewTree()
	for i := 0; i < math.MaxUint16+1; i++ {
		ss.InplaceInsert(strconv.Itoa(i), i)
	}

	s, err = calcRequestAttributeSetOfStringsSize(ss)
	if err == nil {
		t.Errorf("expected *requestTooLongCollectionValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}

	s, err = calcRequestAttributeSetOfStringsSize(newStrTree("one", "two", string(make([]byte, math.MaxUint16+1))))
	if err == nil {
		t.Errorf("expected *requestTooLongStringValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongStringValueError); !ok {
		t.Errorf("expected *requestTooLongStringValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeSetOfNetworksSize(t *testing.T) {
	s, err := calcRequestAttributeSetOfNetworksSize(newIPTree(
		makeTestNetwork("192.0.2.0/24"),
		makeTestNetwork("2001:db8::/32"),
		makeTestNetwork("192.0.2.16/28"),
	))
	if err != nil {
		t.Error(err)
	} else if s != reqBigCounterSize+
		2*(reqIPv4AddressValueSize+reqNetworkCIDRSize)+
		reqIPv6AddressValueSize+reqNetworkCIDRSize {
		t.Errorf("expected %d bytes for value but got %d",
			reqBigCounterSize+
				2*(reqIPv4AddressValueSize+reqNetworkCIDRSize)+
				reqIPv6AddressValueSize+reqNetworkCIDRSize,
			s)
	}

	sn := iptree.NewTree()
	var ip [4]byte
	for i := 0; i < math.MaxUint16+1; i++ {
		binary.BigEndian.PutUint32(ip[:], uint32(i+1))
		ip[0] = 127
		sn.InplaceInsertIP(net.IP(ip[:]), i)
	}

	s, err = calcRequestAttributeSetOfNetworksSize(sn)
	if err == nil {
		t.Errorf("expected *requestTooLongCollectionValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeSetOfDomainsSize(t *testing.T) {
	s, err := calcRequestAttributeSetOfDomainsSize(newDomainTree(
		makeTestDomain("example.com"),
		makeTestDomain("example.gov"),
		makeTestDomain("www.example.com"),
	))
	if err != nil {
		t.Error(err)
	} else if s != 4*reqBigCounterSize+2*11+15 {
		t.Errorf("expected %d bytes for value but got %d", 4*reqBigCounterSize+2*11+15, s)
	}

	sd := new(domaintree.Node)
	for i := 0; i < math.MaxUint16+1; i++ {
		sd.InplaceInsert(makeTestDomain(strconv.Itoa(i)+".com"), i)
	}

	s, err = calcRequestAttributeSetOfDomainsSize(sd)
	if err == nil {
		t.Errorf("expected *requestTooLongCollectionValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}
}

func TestCalcRequestAttributeListOfStringsSize(t *testing.T) {
	s, err := calcRequestAttributeListOfStringsSize([]string{"one", "two", "three"})
	if err != nil {
		t.Error(err)
	} else if s != 4*reqBigCounterSize+3+3+5 {
		t.Errorf("expected %d bytes for value but got %d", 4*reqBigCounterSize+3+3+5, s)
	}

	ls := make([]string, math.MaxUint16+1)
	for i := range ls {
		ls[i] = strconv.Itoa(i)
	}

	s, err = calcRequestAttributeListOfStringsSize(ls)
	if err == nil {
		t.Errorf("expected *requestTooLongCollectionValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongCollectionValueError); !ok {
		t.Errorf("expected *requestTooLongCollectionValueError but got %T (%s)", err, err)
	}

	s, err = calcRequestAttributeListOfStringsSize([]string{"one", "two", string(make([]byte, math.MaxUint16+1))})
	if err == nil {
		t.Errorf("expected *requestTooLongStringValueError but got %d bytes in request", s)
	} else if _, ok := err.(*requestTooLongStringValueError); !ok {
		t.Errorf("expected *requestTooLongStringValueError but got %T (%s)", err, err)
	}
}

func assertRequestBytesBuffer(t *testing.T, desc string, err error, b []byte, n int, e ...byte) {
	if err != nil {
		t.Errorf("expected no error for %s but got: %s", desc, err)
	} else if n != len(b) {
		t.Errorf("expected exactly all buffer used (%d bytes) for %s but got %d bytes", len(b), desc, n)
	} else {
		if bytes.Compare(b[:], e) != 0 {
			t.Errorf("expected [% x] for %s but got [% x]", e, desc, b)
		}
	}
}

func assertRequestBufferOverflow(t *testing.T, desc string, err error, n int) {
	if err == nil {
		t.Errorf("expected no data put to nil buffer for %s but got %d", desc, n)
	} else if _, ok := err.(*requestBufferOverflowError); !ok {
		t.Errorf("expected *requestBufferOverflowError for %s but got %T (%s)", desc, err, err)
	}
}

func assertRequestAssignmentExpressions(t *testing.T, desc string, err error, a []AttributeAssignment, n int, e ...AttributeAssignment) {
	if err != nil {
		t.Errorf("expected no error for %s but got: %s", desc, err)
	} else if n != len(a) {
		t.Errorf("expected exactly all buffer used (%d assignments) for %s but got %d assignments", len(a), desc, n)
	} else {
		aStrs, err := serializeAssignmentExpressions(a)
		if err != nil {
			t.Errorf("can't serialize assignment %d for %s: %s", len(aStrs)+1, desc, err)
			return
		}

		eStrs, err := serializeAssignmentExpressions(e)
		if err != nil {
			t.Errorf("can't serialize expected assignment %d for %s: %s", len(aStrs)+1, desc, err)
			return
		}

		assertStrings(aStrs, eStrs, desc, t)
	}
}

func serializeAssignmentExpressions(a []AttributeAssignment) ([]string, error) {
	out := make([]string, len(a))
	for i, a := range a {
		id, t, v, err := a.Serialize(nil)
		if err != nil {
			return out[:i], err
		}

		out[i] = fmt.Sprintf("id: %q, type: %q, value: %q", id, t, v)
	}

	return out, nil
}

package pdp

import (
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

var (
	testWireAttributes = []byte{
		3, 0,
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger),
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
	}
)

func TestMarshalResponse(t *testing.T) {
	var b [90]byte

	n, err := marshalResponse(b[:], EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError1"),
		fmt.Errorf("testError2"),
	)
	assertRequestBytesBuffer(t, "marshalResponse", err, b[:], n, append(
		[]byte{
			1, 0, 3,
			43, 0, 'm', 'u', 'l', 't', 'i', 'p', 'l', 'e', ' ', 'e', 'r', 'r', 'o', 'r', 's', ':', ' ',
			'"', 't', 'e', 's', 't', 'E', 'r', 'r', 'o', 'r', '1', '"', ',', ' ',
			'"', 't', 'e', 's', 't', 'E', 'r', 'r', 'o', 'r', '2', '"',
		},
		testWireAttributes...)...,
	)

	n, err = marshalResponse([]byte{}, EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError1"),
		fmt.Errorf("testError2"),
	)
	assertRequestBufferOverflow(t, "marshalResponse(version)", err, n)

	n, err = marshalResponse(b[:2], EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError1"),
		fmt.Errorf("testError2"),
	)
	assertRequestBufferOverflow(t, "marshalResponse(effect)", err, n)

	n, err = marshalResponse(b[:5], EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError1"),
		fmt.Errorf("testError2"),
	)
	assertRequestBufferOverflow(t, "marshalResponse(status)", err, n)

	n, err = marshalResponse(b[:22], EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError1"),
		fmt.Errorf("testError2"),
	)
	assertRequestBytesBuffer(t, "marshalResponse(longStatus)", err, b[:22], n,
		1, 0, 3,
		15, 0, 's', 't', 'a', 't', 'u', 's', ' ', 't', 'o', 'o', ' ', 'l', 'o', 'n', 'g',
		0, 0,
	)

	n, err = marshalResponse(b[:27], EffectIndeterminate, testRequestAssignments, fmt.Errorf("testError"))
	assertRequestBytesBuffer(t, "marshalResponse(longObligation)", err, b[:27], n,
		1, 0, 3,
		20, 0, 'o', 'b', 'l', 'i', 'g', 'a', 't', 'i', 'o', 'n', 's', ' ', 't', 'o', 'o', ' ', 'l', 'o', 'n', 'g',
		0, 0,
	)

	n, err = marshalResponse(b[:14], EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError"),
	)
	assertRequestBufferOverflow(t, "marshalResponse(error)", err, n)

	n, err = marshalResponse(b[:20], EffectIndeterminate, testRequestAssignments,
		fmt.Errorf("testError1"),
		fmt.Errorf("testError2"),
	)
	assertRequestBufferOverflow(t, "marshalResponse(multi-error)", err, n)

	n, err = marshalResponse(b[:25], EffectIndeterminate, testRequestAssignments, fmt.Errorf("testError"))
	assertRequestBufferOverflow(t, "marshalResponse(longObligation)", err, n)

	n, err = marshalResponse(b[:], EffectIndeterminate, []AttributeAssignment{
		MakeAddressAssignment("address", net.IP{1, 2, 3, 4, 5, 6}),
	})
	if err == nil {
		t.Errorf("expected no data put to buffer for response with invalid network but got %d", n)
	} else if _, ok := err.(*requestAddressValueError); !ok {
		t.Errorf("expected *requestAddressValueError but got %T (%s)", err, err)
	}
}

func TestMakeIndeterminateResponse(t *testing.T) {
	var b [17]byte

	n, err := MakeIndeterminateResponse(b[:], fmt.Errorf("test error"))
	assertRequestBytesBuffer(t, "MakeIndeterminateResponse", err, b[:], n,
		1, 0, 3,
		10, 0, 't', 'e', 's', 't', ' ', 'e', 'r', 'r', 'o', 'r',
		0, 0,
	)
}

func TestUnmarshalResponse(t *testing.T) {
	var a [3]AttributeAssignment

	effect, n, err := UnmarshalResponse(append([]byte{1, 0, 1, 0, 0}, testWireAttributes...), a[:])
	assertRequestAssignmentExpressions(t, "UnmarshalResponse", err, a[:], n, testRequestAssignments...)
	if effect != EffectPermit {
		t.Errorf("expected %q effect but got %q", EffectNameFromEnum(EffectPermit), EffectNameFromEnum(effect))
	}

	effect, n, err = UnmarshalResponse(append([]byte{
		1, 0, 3,
		9, 0, 't', 'e', 's', 't', 'E', 'r', 'r', 'o', 'r',
	}, testWireAttributes...), a[:])
	if err == nil {
		t.Errorf("expected *ResponseServerError but got no error")
	} else if _, ok := err.(*ResponseServerError); !ok {
		t.Errorf("expected *ResponseServerError but got %T (%s)", err, err)
	}

	assertRequestAssignmentExpressions(t, "UnmarshalResponse", nil, a[:], n, testRequestAssignments...)
	if effect != EffectIndeterminate {
		t.Errorf("expected %q effect but got %q", EffectNameFromEnum(EffectIndeterminate), EffectNameFromEnum(effect))
	}

	effect, n, err = UnmarshalResponse([]byte{}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got effect %q and %d attributes",
			EffectNameFromEnum(effect), n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	effect, n, err = UnmarshalResponse([]byte{
		1, 0,
	}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got effect %q and %d attributes",
			EffectNameFromEnum(effect), n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	effect, n, err = UnmarshalResponse([]byte{
		1, 0, 3,
	}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got effect %q and %d attributes",
			EffectNameFromEnum(effect), n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	effect, n, err = UnmarshalResponse([]byte{
		1, 0, 3, 0, 0,
	}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got effect %q and %d attributes",
			EffectNameFromEnum(effect), n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestUnmarshalResponseToReflection(t *testing.T) {
	var (
		effect int
		s      string
		b      bool
		i64    int64
	)

	if err := UnmarshalResponseToReflection(append([]byte{1, 0, 1, 0, 0}, testWireAttributes...),
		func(id string, t Type) (reflect.Value, error) {
			switch id {
			case ResponseEffectFieldName:
				return reflect.Indirect(reflect.ValueOf(&effect)), nil

			case ResponseStatusFieldName:
				return reflectValueNil, nil

			case "string":
				return reflect.Indirect(reflect.ValueOf(&s)), nil

			case "boolean":
				return reflect.Indirect(reflect.ValueOf(&b)), nil

			case "integer":
				return reflect.Indirect(reflect.ValueOf(&i64)), nil
			}

			return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
		},
	); err != nil {
		t.Error(err)
	} else {
		if effect != EffectPermit {
			t.Errorf("expected %q effect but got %q", EffectNameFromEnum(EffectPermit), EffectNameFromEnum(effect))
		}

		a := []AttributeAssignment{
			MakeStringAssignment("string", s),
			MakeBooleanAssignment("boolean", b),
			MakeIntegerAssignment("integer", i64),
		}
		assertRequestAssignmentExpressions(t, "UnmarshalResponseToReflection", err, a, 3, testRequestAssignments...)
	}

	err := UnmarshalResponseToReflection([]byte{255, 255}, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected *requestVersionError")
	} else if _, ok := err.(*requestVersionError); !ok {
		t.Errorf("expected *requestVersionError but got %T (%s)", err, err)
	}

	err = UnmarshalResponseToReflection([]byte{1, 0, 255}, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected *responseEffectError")
	} else if _, ok := err.(*responseEffectError); !ok {
		t.Errorf("expected *responseEffectError but got %T (%s)", err, err)
	}

	testErr := fmt.Errorf("testError")
	err = UnmarshalResponseToReflection([]byte{1, 0, 1}, func(id string, t Type) (reflect.Value, error) {
		if id == ResponseEffectFieldName {
			return reflectValueNil, testErr
		}

		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected testErr")
	} else if err != testErr {
		t.Errorf("expected testErr but got %T (%s)", err, err)
	}

	err = UnmarshalResponseToReflection([]byte{1, 0, 1}, func(id string, t Type) (reflect.Value, error) {
		if id == ResponseEffectFieldName {
			return reflect.ValueOf(effect), nil
		}

		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected *requestUnmarshalEffectConstError")
	} else if _, ok := err.(*requestUnmarshalEffectConstError); !ok {
		t.Errorf("expected *requestUnmarshalEffectConstError but got %T (%s)", err, err)
	}

	err = UnmarshalResponseToReflection([]byte{
		1, 0, 1, 8, 0, 't', 'e', 's', 't',
	}, func(id string, t Type) (reflect.Value, error) {
		if id == ResponseEffectFieldName {
			return reflect.Indirect(reflect.ValueOf(&effect)), nil
		}

		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = UnmarshalResponseToReflection([]byte{
		1, 0, 1, 4, 0, 't', 'e', 's', 't',
	}, func(id string, t Type) (reflect.Value, error) {
		if id == ResponseEffectFieldName {
			return reflect.Indirect(reflect.ValueOf(&effect)), nil
		}

		if id == ResponseStatusFieldName {
			return reflectValueNil, testErr
		}

		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected testErr")
	} else if err != testErr {
		t.Errorf("expected testErr but got %T (%s)", err, err)
	}

	err = UnmarshalResponseToReflection([]byte{
		1, 0, 1, 4, 0, 't', 'e', 's', 't',
	}, func(id string, t Type) (reflect.Value, error) {
		if id == ResponseEffectFieldName {
			return reflect.Indirect(reflect.ValueOf(&effect)), nil
		}

		if id == ResponseStatusFieldName {
			return reflect.ValueOf(s), nil
		}

		return reflectValueNil, fmt.Errorf("unexpected attribute %s.(%s)", id, t)
	})
	if err == nil {
		t.Error("expected *requestUnmarshalStatusConstError")
	} else if _, ok := err.(*requestUnmarshalStatusConstError); !ok {
		t.Errorf("expected *requestUnmarshalStatusConstError but got %T (%s)", err, err)
	}
}

func TestPutResponseEffect(t *testing.T) {
	var b [1]byte

	n, err := putResponseEffect(b[:], EffectPermit)
	assertRequestBytesBuffer(t, "putResponseEffect", err, b[:], n, 1)

	n, err = putResponseEffect([]byte{}, EffectPermit)
	assertRequestBufferOverflow(t, "putResponseEffect", err, n)

	n, err = putResponseEffect(b[:], -1)
	if err == nil {
		t.Errorf("expected no data put to buffer for invalid effect but got %d", n)
	} else if _, ok := err.(*responseEffectError); !ok {
		t.Errorf("expected *responseEffectError but got %T (%s)", err, err)
	}
}

func TestGetResponseEffect(t *testing.T) {
	effect, n, err := getResponseEffect([]byte{1})
	if err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Errorf("expected one byte consumed but got %d", n)
	} else if effect != EffectPermit {
		t.Errorf("expected %q effect but got %q",
			EffectNameFromEnum(EffectPermit), EffectNameFromEnum(effect),
		)
	}

	effect, n, err = getResponseEffect([]byte{})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but consumed %d bytes and got %q effect",
			n, EffectNameFromEnum(effect),
		)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	effect, n, err = getResponseEffect([]byte{255})
	if err == nil {
		t.Errorf("expected *responseEffectError but consumed %d bytes and got %q effect",
			n, EffectNameFromEnum(effect),
		)
	} else if _, ok := err.(*responseEffectError); !ok {
		t.Errorf("expected *responseEffectError but got %T (%s)", err, err)
	}
}

func TestPutResponseStatus(t *testing.T) {
	var b [65536]byte

	n, err := putResponseStatus(b[:])
	assertRequestBytesBuffer(t, "putResponseStatus", err, b[:2], n,
		0, 0,
	)

	n, err = putResponseStatus(b[:], fmt.Errorf("test"))
	assertRequestBytesBuffer(t, "putResponseStatus(1)", err, b[:6], n,
		4, 0, 't', 'e', 's', 't',
	)

	n, err = putResponseStatus(b[:], fmt.Errorf("test1"), fmt.Errorf("test2"))
	assertRequestBytesBuffer(t, "putResponseStatus(2)", err, b[:35], n,
		33, 0, 'm', 'u', 'l', 't', 'i', 'p', 'l', 'e', ' ', 'e', 'r', 'r', 'o', 'r', 's', ':', ' ',
		'"', 't', 'e', 's', 't', '1', '"', ',', ' ', '"', 't', 'e', 's', 't', '2', '"',
	)

	n, err = putResponseStatus([]byte{})
	assertRequestBufferOverflow(t, "putResponseStatus", err, n)

	n, err = putResponseStatus([]byte{}, fmt.Errorf("test"))
	assertRequestBufferOverflow(t, "putResponseStatus(1)", err, n)

	s := ""
	for i := 0; i < 6553; i++ {
		s += "0123456789"
	}
	s += "0123\u56db56789"

	e := make([]byte, 65536)
	e[0] = 254
	e[1] = 255
	for i := 0; i < 6553; i++ {
		copy(e[10*i+2:], "0123456789")
	}
	e[65532] = '0'
	e[65533] = '1'
	e[65534] = '2'
	e[65535] = '3'

	n, err = putResponseStatus(b[:], fmt.Errorf(s))
	assertRequestBytesBuffer(t, "putResponseStatus(long)", err, b[:], n, e...)
}

func TestPutResponseStatusTooLong(t *testing.T) {
	if len(responseStatusTooLong) > math.MaxUint16 {
		t.Errorf("expected no more than %d bytes for responseStatusTooLong but got %d",
			math.MaxUint16, len(responseStatusTooLong),
		)
	}

	var b [17]byte

	n, err := putResponseStatusTooLong(b[:])
	assertRequestBytesBuffer(t, "putResponseStatusTooLong", err, b[:], n,
		15, 0, 's', 't', 'a', 't', 'u', 's', ' ', 't', 'o', 'o', ' ', 'l', 'o', 'n', 'g',
	)

	n, err = putResponseStatusTooLong([]byte{})
	assertRequestBufferOverflow(t, "putResponseStatusTooLong", err, n)
}

func TestPutResponseObligationsTooLong(t *testing.T) {
	if len(responseStatusObligationsTooLong) > math.MaxUint16 {
		t.Errorf("expected no more than %d bytes for responseStatusObligationsTooLong but got %d",
			math.MaxUint16, len(responseStatusObligationsTooLong),
		)
	}

	var b [22]byte

	n, err := putResponseObligationsTooLong(b[:])
	assertRequestBytesBuffer(t, "putResponseObligationsTooLong", err, b[:], n,
		20, 0, 'o', 'b', 'l', 'i', 'g', 'a', 't', 'i', 'o', 'n', 's', ' ', 't', 'o', 'o', ' ', 'l', 'o', 'n', 'g',
	)

	n, err = putResponseObligationsTooLong([]byte{})
	assertRequestBufferOverflow(t, "putResponseObligationsTooLong", err, n)
}

func TestPutAssignmentExpressions(t *testing.T) {
	var b [42]byte
	n, err := putAssignmentExpressions(b[:], testRequestAssignments)
	assertRequestBytesBuffer(t, "putAssignmentExpressions", err, b[:], n, testWireAttributes...)

	n, err = putAssignmentExpressions([]byte{}, testRequestAssignments)
	assertRequestBufferOverflow(t, "putAssignmentExpressions", err, n)

	n, err = putAssignmentExpressions(b[:], []AttributeAssignment{
		MakeExpressionAssignment("boolean", makeFunctionBooleanNot([]Expression{MakeBooleanValue(true)})),
	})
	if err == nil {
		t.Errorf("expected no data put to buffer for invalid expression but got %d", n)
	} else if _, ok := err.(*requestInvalidExpressionError); !ok {
		t.Errorf("expected *requestInvalidExpressionError but got %T (%s)", err, err)
	}

	n, err = putAssignmentExpressions(b[:12], testRequestAssignments)
	assertRequestBufferOverflow(t, "putAssignmentExpressions(expressions)", err, n)
}

func TestPutAttributesFromReflection(t *testing.T) {
	var b [287]byte

	f := func(i int) (string, Type, reflect.Value, error) {
		switch i {
		case 0:
			return "boolean", TypeBoolean, reflect.ValueOf(true), nil

		case 1:
			return "string", TypeString, reflect.ValueOf("test"), nil

		case 2:
			return "integer", TypeInteger, reflect.ValueOf(int64(9223372036854775807)), nil

		case 3:
			return "float", TypeFloat, reflect.ValueOf(float64(math.Pi)), nil

		case 4:
			return "address", TypeAddress, reflect.ValueOf(net.ParseIP("192.0.2.1")), nil

		case 5:
			return "network", TypeNetwork, reflect.ValueOf(makeTestNetwork("192.0.2.0/24")), nil

		case 6:
			return "domain", TypeDomain, reflect.ValueOf(makeTestDomain("www.example.com")), nil

		case 7:
			return "set of strings", TypeSetOfStrings, reflect.ValueOf(newStrTree("one", "two", "three")), nil

		case 8:
			return "set of networks", TypeSetOfNetworks, reflect.ValueOf(newIPTree(
				makeTestNetwork("192.0.2.0/24"),
				makeTestNetwork("2001:db8::/32"),
				makeTestNetwork("192.0.2.16/28"),
			)), nil

		case 9:
			return "set of domains", TypeSetOfDomains, reflect.ValueOf(newDomainTree(
				makeTestDomain("example.com"),
				makeTestDomain("example.gov"),
				makeTestDomain("www.example.com"),
			)), nil

		case 10:
			return "list of strings", TypeListOfStrings, reflect.ValueOf([]string{"one", "two", "three"}), nil
		}

		return "", TypeUndefined, reflectValueNil, fmt.Errorf("unexpected intex %d", i)
	}
	n, err := putAttributesFromReflection(b[:], 11, f)
	assertRequestBytesBuffer(t, "putAttributesFromReflection", err, b[:], n,
		11, 0,
		7, 'b', 'o', 'o', 'l', 'e', 'a', 'n', byte(requestWireTypeBooleanTrue),
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger),
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84, 251, 33, 9, 64,
		7, 'a', 'd', 'd', 'r', 'e', 's', 's', byte(requestWireTypeIPv4Address), 192, 0, 2, 1,
		7, 'n', 'e', 't', 'w', 'o', 'r', 'k', byte(requestWireTypeIPv4Network), 24, 192, 0, 2, 0,
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain),
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings),
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
		15, 's', 'e', 't', ' ', 'o', 'f', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', 's',
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains),
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		15, 'l', 'i', 's', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's',
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	)

	n, err = putAttributesFromReflection([]byte{}, 1, f)
	assertRequestBufferOverflow(t, "putAttributesFromReflection", err, n)

	testFuncErr := errors.New("test function error")
	n, err = putAttributesFromReflection(b[:], 1, func(i int) (string, Type, reflect.Value, error) {
		return "", TypeUndefined, reflectValueNil, testFuncErr
	})
	if err == nil {
		t.Errorf("expected no data put to buffer for broken function but got %d", n)
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}

	n, err = putAttributesFromReflection(b[:], 1, func(i int) (string, Type, reflect.Value, error) {
		return "undefined", TypeUndefined, reflectValueNil, nil
	})
	if err == nil {
		t.Errorf("expected no data put to buffer for undefined value but got %d", n)
	} else if _, ok := err.(*requestAttributeMarshallingNotImplemented); !ok {
		t.Errorf("expected *requestAttributeMarshallingNotImplemented but got %T (%s)", err, err)
	}

	n, err = putAttributesFromReflection(b[:10], 1, f)
	assertRequestBufferOverflow(t, "putAttributesFromReflection(values)", err, n)
}

func TestGetAssignmentExpressions(t *testing.T) {
	var a [3]AttributeAssignment

	n, err := getAssignmentExpressions(testWireAttributes, a[:])
	assertRequestAssignmentExpressions(t, "getAssignmentExpressions", err, a[:], n, testRequestAssignments...)

	n, err = getAssignmentExpressions([]byte{}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %d bytes", n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	n, err = getAssignmentExpressions([]byte{255, 255}, a[:])
	if err == nil {
		t.Errorf("expected *requestAssignmentsOverflowError but got %d bytes", n)
	} else if _, ok := err.(*requestAssignmentsOverflowError); !ok {
		t.Errorf("expected *requestAssignmentsOverflowError but got %T (%s)", err, err)
	}

	n, err = getAssignmentExpressions([]byte{1, 0}, a[:])
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got got %d bytes", n)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

func TestGetAttributesToReflection(t *testing.T) {
	var (
		names [14]string
	)

	i := 0

	booleanFalse := true
	booleanTrue := false

	var (
		str string
		num int64
		flt float64
		av4 net.IP
		av6 net.IP
		nv4 *net.IPNet
		nv6 *net.IPNet
		dn  domain.Name
		ss  *strtree.Tree
		sn  *iptree.Tree
		sd  *domaintree.Node
		ls  []string
	)

	values := []reflect.Value{
		reflect.Indirect(reflect.ValueOf(&booleanFalse)),
		reflect.Indirect(reflect.ValueOf(&booleanTrue)),
		reflect.Indirect(reflect.ValueOf(&str)),
		reflect.Indirect(reflect.ValueOf(&num)),
		reflect.Indirect(reflect.ValueOf(&flt)),
		reflect.Indirect(reflect.ValueOf(&av4)),
		reflect.Indirect(reflect.ValueOf(&av6)),
		reflect.Indirect(reflect.ValueOf(&nv4)),
		reflect.Indirect(reflect.ValueOf(&nv6)),
		reflect.Indirect(reflect.ValueOf(&dn)),
		reflect.Indirect(reflect.ValueOf(&ss)),
		reflect.Indirect(reflect.ValueOf(&sn)),
		reflect.Indirect(reflect.ValueOf(&sd)),
		reflect.Indirect(reflect.ValueOf(&ls)),
	}

	err := getAttributesToReflection([]byte{
		byte(len(names)), 0,
		12, 'b', 'o', 'o', 'l', 'e', 'a', 'n', 'F', 'a', 'l', 's', 'e', byte(requestWireTypeBooleanFalse),
		11, 'b', 'o', 'o', 'l', 'e', 'a', 'n', 'T', 'r', 'u', 'e', byte(requestWireTypeBooleanTrue),
		6, 's', 't', 'r', 'i', 'n', 'g', byte(requestWireTypeString), 4, 0, 't', 'e', 's', 't',
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 1, 0, 0, 0, 0, 0, 0, 0,
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84, 251, 33, 9, 64,
		8, 'a', 'd', 'd', 'r', 'e', 's', 's', '4', byte(requestWireTypeIPv4Address), 192, 0, 2, 1,
		8, 'a', 'd', 'd', 'r', 'e', 's', 's', '6', byte(requestWireTypeIPv6Address),
		32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		8, 'n', 'e', 't', 'w', 'o', 'r', 'k', '4', byte(requestWireTypeIPv4Network), 24, 192, 0, 2, 0,
		8, 'n', 'e', 't', 'w', 'o', 'r', 'k', '6', byte(requestWireTypeIPv6Network),
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain),
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings),
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
		15, 's', 'e', 't', ' ', 'o', 'f', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', 's',
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		220, 192, 0, 2, 16,
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains),
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		15, 'l', 'i', 's', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's',
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
		5, 0, 't', 'h', 'r', 'e', 'e',
	}, func(id string, t Type) (reflect.Value, error) {
		if i >= len(names) || i >= len(values) || i >= len(builtinTypeByWire) {
			return reflectValueNil, fmt.Errorf("requested invalid value number: %d", i)
		}

		if et := builtinTypeByWire[i]; t != et {
			return reflectValueNil, fmt.Errorf("expected %q for %d but got %q", et, i, t)
		}

		names[i] = id
		v := values[i]
		i++

		return v, nil
	})

	a := []AttributeAssignment{
		MakeBooleanAssignment(names[0], booleanFalse),
		MakeBooleanAssignment(names[1], booleanTrue),
		MakeStringAssignment(names[2], str),
		MakeIntegerAssignment(names[3], num),
		MakeFloatAssignment(names[4], flt),
		MakeAddressAssignment(names[5], av4),
		MakeAddressAssignment(names[6], av6),
		MakeNetworkAssignment(names[7], nv4),
		MakeNetworkAssignment(names[8], nv6),
		MakeDomainAssignment(names[9], dn),
		MakeSetOfStringsAssignment(names[10], ss),
		MakeSetOfNetworksAssignment(names[11], sn),
		MakeSetOfDomainsAssignment(names[12], sd),
		MakeListOfStringsAssignment(names[13], ls),
	}

	assertRequestAssignmentExpressions(t, "getAttributesToReflection", err, a, i,
		MakeBooleanAssignment("booleanFalse", false),
		MakeBooleanAssignment("booleanTrue", true),
		MakeStringAssignment("string", "test"),
		MakeIntegerAssignment("integer", 1),
		MakeFloatAssignment("float", float64(math.Pi)),
		MakeAddressAssignment("address4", net.ParseIP("192.0.2.1")),
		MakeAddressAssignment("address6", net.ParseIP("2001:db8::1")),
		MakeNetworkAssignment("network4", makeTestNetwork("192.0.2.0/24")),
		MakeNetworkAssignment("network6", makeTestNetwork("2001:db8::/32")),
		MakeDomainAssignment("domain", makeTestDomain("www.example.com")),
		MakeSetOfStringsAssignment("set of strings", newStrTree("one", "two", "three")),
		MakeSetOfNetworksAssignment("set of networks", newIPTree(
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		)),
		MakeSetOfDomainsAssignment("set of domains", newDomainTree(
			makeTestDomain("example.com"),
			makeTestDomain("example.gov"),
			makeTestDomain("www.example.com"),
		)),
		MakeListOfStringsAssignment("list of strings", []string{"one", "two", "three"}),
	)

	err = getAttributesToReflection([]byte{}, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, fmt.Errorf("in unreacheable place with id %q and type %q", id, t)
	})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError but got nothing")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, fmt.Errorf("in unreacheable place with id %q and type %q", id, t)
	})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError but got nothing")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		7, 's', 't', 'r', 'i', 'n', 'g', 's',
	}, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, fmt.Errorf("in unreacheable place with id %q and type %q", id, t)
	})
	if err == nil {
		t.Error("expected *requestBufferUnderflowError but got nothing")
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		7, 's', 't', 'r', 'i', 'n', 'g', 's', 255,
	}, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, fmt.Errorf("in unreacheable place with id %q and type %q", id, t)
	})
	if err == nil {
		t.Error("expected *requestAttributeUnmarshallingTypeError but got nothing")
	} else if _, ok := err.(*requestAttributeUnmarshallingTypeError); !ok {
		t.Errorf("expected *requestAttributeUnmarshallingTypeError but got %T (%s)", err, err)
	}

	testFuncErr := errors.New("test function error")
	err = getAttributesToReflection(testWireAttributes, func(id string, t Type) (reflect.Value, error) {
		return reflectValueNil, testFuncErr
	})
	if err == nil {
		t.Error("expected testFuncErr but got nothing")
	} else if err != testFuncErr {
		t.Errorf("expected testFuncErr but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		7, 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeString), 4, 0, 't', 'e',
	}, func(id string, t Type) (reflect.Value, error) {
		return values[2], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", str)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 1, 0, 0, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[3], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %d", num)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	var i8 int8
	v := reflect.Indirect(reflect.ValueOf(&i8))

	err = getAttributesToReflection([]byte{
		1, 0,
		7, 'i', 'n', 't', 'e', 'g', 'e', 'r', byte(requestWireTypeInteger), 128, 0, 0, 0, 0, 0, 0, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return v, nil
	})
	if err == nil {
		t.Errorf("expected *requestUnmarshalIntegerOverflowError but got %d", i8)
	} else if _, ok := err.(*requestUnmarshalIntegerOverflowError); !ok {
		t.Errorf("expected *requestUnmarshalIntegerOverflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		5, 'f', 'l', 'o', 'a', 't', byte(requestWireTypeFloat), 24, 45, 68, 84,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[4], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %g", flt)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		8, 'a', 'd', 'd', 'r', 'e', 's', 's', '4', byte(requestWireTypeIPv4Address), 192, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[5], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", av4)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		8, 'a', 'd', 'd', 'r', 'e', 's', 's', '6', byte(requestWireTypeIPv6Address), 32, 1, 13, 184, 0, 0, 0, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[6], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", av6)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		8, 'n', 'e', 't', 'w', 'o', 'r', 'k', '4', byte(requestWireTypeIPv4Network), 24, 192, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[7], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", nv4)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		8, 'n', 'e', 't', 'w', 'o', 'r', 'k', '6', byte(requestWireTypeIPv6Network), 32, 32, 1, 13, 184, 0, 0, 0, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[8], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", nv6)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		6, 'd', 'o', 'm', 'a', 'i', 'n', byte(requestWireTypeDomain),
		15, 0, 'w', 'w', 'w', '.', 'e', 'x', 'a', 'm', 'p', 'l', 'e',
	}, func(id string, t Type) (reflect.Value, error) {
		return values[9], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %q", dn)
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's', byte(requestWireTypeSetOfStrings),
		3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
	}, func(id string, t Type) (reflect.Value, error) {
		return values[10], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", MakeSetOfStringsValue(ss).describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		15, 's', 'e', 't', ' ', 'o', 'f', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', 's',
		byte(requestWireTypeSetOfNetworks), 3, 0,
		216, 192, 0, 2, 0,
		32, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}, func(id string, t Type) (reflect.Value, error) {
		return values[11], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", MakeSetOfNetworksValue(sn).describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		14, 's', 'e', 't', ' ', 'o', 'f', ' ', 'd', 'o', 'm', 'a', 'i', 'n', 's', byte(requestWireTypeSetOfDomains),
		3, 0,
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'c', 'o', 'm',
		11, 0, 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'g', 'o', 'v',
	}, func(id string, t Type) (reflect.Value, error) {
		return values[12], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", MakeSetOfDomainsValue(sd).describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}

	err = getAttributesToReflection([]byte{
		1, 0,
		15, 'l', 'i', 's', 't', ' ', 'o', 'f', ' ', 's', 't', 'r', 'i', 'n', 'g', 's',
		byte(requestWireTypeListOfStrings), 3, 0,
		3, 0, 'o', 'n', 'e',
		3, 0, 't', 'w', 'o',
	}, func(id string, t Type) (reflect.Value, error) {
		return values[13], nil
	})
	if err == nil {
		t.Errorf("expected *requestBufferUnderflowError but got %s", MakeListOfStringsValue(ls).describe())
	} else if _, ok := err.(*requestBufferUnderflowError); !ok {
		t.Errorf("expected *requestBufferUnderflowError but got %T (%s)", err, err)
	}
}

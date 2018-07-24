package pdp

import (
	"encoding/binary"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
)

const (
	responseStatusTooLong            = "status too long"
	responseStatusObligationsTooLong = "obligations too long"
)

// Names of special response fields.
const (
	// ResponseEffectFieldName holds name of response effect.
	ResponseEffectFieldName = "effect"
	// ResponseStatusFieldName stores name of response status.
	ResponseStatusFieldName = "status"
)

// MinResponseSize represents lower response buffer limit required to return
// error that real error message or set of obligations are too long.
var MinResponseSize uint

const (
	resEffectSize = 1
)

func init() {
	n := len(responseStatusTooLong)
	if len(responseStatusObligationsTooLong) > n {
		n = len(responseStatusObligationsTooLong)
	}

	MinResponseSize = 7 + uint(n)
}

func marshalResponse(effect int, obligations []AttributeAssignment, errs ...error) ([]byte, error) {
	n, err := calcResponseSize(obligations, errs...)
	if err != nil {
		return nil, err
	}

	b := make([]byte, n)
	_, err = marshalResponseToBuffer(b, effect, obligations, errs...)
	return b, err
}

func marshalResponseWithAllocator(f func(n int) ([]byte, error), effect int, obligations []AttributeAssignment, errs ...error) ([]byte, error) {
	n, err := calcResponseSize(obligations, errs...)
	if err != nil {
		return nil, err
	}

	b, err := f(n)
	if err != nil {
		return nil, err
	}

	n, err = marshalResponseToBuffer(b, effect, obligations, errs...)
	if err != nil {
		return nil, err
	}

	return b[:n], nil
}

func marshalResponseToBuffer(b []byte, effect int, obligations []AttributeAssignment, errs ...error) (int, error) {
	off, err := putRequestVersion(b)
	if err != nil {
		return off, err
	}

	n, err := putResponseEffect(b[off:], effect)
	if err != nil {
		return off, err
	}
	off += n

	n, err = putResponseStatus(b[off:], errs...)
	if err != nil {
		n, err = putResponseStatusTooLong(b[off:])
		if err != nil {
			return off, err
		}
		off += n

		n, err = putRequestAttributeCount(b[off:], 0)
		if err != nil {
			return off, err
		}

		return off + n, nil
	}

	off += n

	n, err = putAssignmentExpressions(b[off:], obligations)
	if err != nil {
		if _, ok := err.(*requestBufferOverflowError); ok {
			off, _ := putRequestVersion(b)
			n, _ := putResponseEffect(b[off:], effect)
			off += n

			n, err := putResponseObligationsTooLong(b[off:])
			if err != nil {
				return off, err
			}
			off += n

			n, err = putRequestAttributeCount(b[off:], 0)
			if err != nil {
				return off, err
			}

			return off + n, nil
		}

		return off, err
	}

	return off + n, nil
}

// MakeIndeterminateResponse marshals given error as indenterminate response
// with no obligations as a sequebce of bytes.
func MakeIndeterminateResponse(err error) ([]byte, error) {
	return marshalResponse(EffectIndeterminate, nil, err)
}

// MakeIndeterminateResponseWithAllocator marshals given error as indenterminate
// response with no obligations as a sequebce of bytes. The allocator is
// expected to take number of bytes required and return slice of that length.
func MakeIndeterminateResponseWithAllocator(f func(n int) ([]byte, error), err error) ([]byte, error) {
	return marshalResponseWithAllocator(f, EffectIndeterminate, nil, err)
}

// MakeIndeterminateResponseWithBuffer marshals given error as indenterminate
// response with no obligations to given buffer. Caller need to allocate big
// enough buffer. It should be at least MinResponseSize to put message that
// buffer isn't long enough. The function returns number of bytes written to
// the buffer.
func MakeIndeterminateResponseWithBuffer(b []byte, err error) (int, error) {
	return marshalResponseToBuffer(b, EffectIndeterminate, nil, err)
}

// UnmarshalResponseAssignments unmarshals response from given sequence of
// bytes. Effect is returned as the first result value. The second returned
// value is an array of obligations. Finally, the third value is an error
// occured during unmarshalling or response status if it has type
// *ResponseServerError.
func UnmarshalResponseAssignments(b []byte) (int, []AttributeAssignment, error) {
	off, err := checkRequestVersion(b)
	if err != nil {
		return EffectIndeterminate, nil, err
	}

	effect, n, err := getResponseEffect(b[off:])
	if err != nil {
		return EffectIndeterminate, nil, err
	}
	off += n

	s, n, err := getRequestStringValue(b[off:])
	if err != nil {
		return EffectIndeterminate, nil, err
	}
	off += n

	out, err := getAssignmentExpressions(b[off:])
	if err != nil {
		return EffectIndeterminate, nil, err
	}

	if len(s) > 0 {
		return effect, out, newResponseServerError(s)
	}

	return effect, out, nil
}

// UnmarshalResponseAssignmentsWithAllocator works similarly to
// UnmarshalResponseAssignments but requires custom allocator for obligations.
// The allocator is expected to take number of obligations and return slice of
// assignments of that length.
func UnmarshalResponseAssignmentsWithAllocator(b []byte, f func(n int) ([]AttributeAssignment, error)) (int, []AttributeAssignment, error) {
	off, err := checkRequestVersion(b)
	if err != nil {
		return EffectIndeterminate, nil, err
	}

	effect, n, err := getResponseEffect(b[off:])
	if err != nil {
		return EffectIndeterminate, nil, err
	}
	off += n

	s, n, err := getRequestStringValue(b[off:])
	if err != nil {
		return EffectIndeterminate, nil, err
	}
	off += n

	out, err := getAssignmentExpressionsWithAllocator(b[off:], f)
	if err != nil {
		return EffectIndeterminate, nil, err
	}

	if len(s) > 0 {
		return effect, out, newResponseServerError(s)
	}

	return effect, out, nil
}

// UnmarshalResponseToAssignmentsArray unmarshals response from given
// sequence of bytes. Effect is returned as the first result value. The second
// returned value gives number of obligations put to out parameter. Finally,
// the third value is an error occured during unmarshalling or response status
// if it has type *ResponseServerError. Caller need to allocate and pass big
// enough array to out argument.
func UnmarshalResponseToAssignmentsArray(b []byte, out []AttributeAssignment) (int, int, error) {
	off, err := checkRequestVersion(b)
	if err != nil {
		return EffectIndeterminate, 0, err
	}

	effect, n, err := getResponseEffect(b[off:])
	if err != nil {
		return EffectIndeterminate, 0, err
	}
	off += n

	s, n, err := getRequestStringValue(b[off:])
	if err != nil {
		return EffectIndeterminate, 0, err
	}
	off += n

	n, err = getAssignmentExpressionsToArray(b[off:], out)
	if err != nil {
		return EffectIndeterminate, 0, err
	}

	if len(s) > 0 {
		return effect, n, newResponseServerError(s)
	}

	return effect, n, nil
}

// UnmarshalResponseToReflection unmarshals response from given sequence
// of bytes to a set reflected values. The function extracts a parameter or
// obligation from response and calls f function with its name and type.
// The function should return reflected value to put data to. If f returns
// error unmarshlling stopped with the error. If f don't want to get value of
// attribute or response parameter it can return invalid reflect.Value
// (reflect.Value(nil). For Effect parameter UnmarshalResponseToReflection
// passes to f ResponseEffectFieldName as name and nil type and expectes
// value of bool, string, intX or uintX (for bool true means EffectPermit and
// false all other effects). For Status parameter ResponseStatusFieldName with
// nil type passed to f and string or error expected as reflected value.
// For any obligation its name and Type passed to f. Which value is expected
// depends on attribute type for TypeBoolean - bool, TypeString - string,
// TypeInteger - intX or uintX (note that small int types can be overflowed
// while uint can't take negative value), TypeFloat - float32/64, TypeAddress -
// net.IP, TypeNetwork - net.IPNet or *net.IPNet, TypeDomain - string or
// domain.Name from github.com/infobloxopen/go-trees/domain package,
// TypeSetOfStrings - *strtree.Tree from
// github.com/infobloxopen/go-trees/strtree package, TypeSetOfNetworks -
// *iptree.Tree from github.com/infobloxopen/go-trees/iptree package,
// TypeSetOfDomains - *domaintree.Node from
// github.com/infobloxopen/go-trees/domaintree package, TypeListOfStrings -
// []string.
func UnmarshalResponseToReflection(b []byte, f func(string, Type) (reflect.Value, error)) error {
	off, err := checkRequestVersion(b)
	if err != nil {
		return err
	}

	effect, n, err := getResponseEffect(b[off:])
	if err != nil {
		return err
	}
	off += n

	v, err := f(ResponseEffectFieldName, nil)
	if err != nil {
		return err
	}

	if err := setEffect(v, effect); err != nil {
		return err
	}

	s, n, err := getRequestStringValue(b[off:])
	if err != nil {
		return err
	}
	off += n

	v, err = f(ResponseStatusFieldName, nil)
	if err != nil {
		return err
	}

	if err := setStatus(v, s); err != nil {
		return err
	}

	return getAttributesToReflection(b[off:], f)
}

func putResponseEffect(b []byte, effect int) (int, error) {
	if effect < 0 || effect >= effectTotalCount {
		return 0, newResponseEffectError(effect)
	}

	if len(b) < 1 {
		return 0, newRequestBufferOverflowError()
	}

	b[0] = byte(effect)
	return 1, nil
}

func getResponseEffect(b []byte) (int, int, error) {
	if len(b) < 1 {
		return EffectIndeterminate, 0, newRequestBufferUnderflowError()
	}

	effect := int(b[0])
	if effect < 0 || effect >= effectTotalCount {
		return EffectIndeterminate, 0, newResponseEffectError(effect)
	}

	return effect, 1, nil
}

func putResponseStatus(b []byte, err ...error) (int, error) {
	if len(err) < 1 || len(err) == 1 && err[0] == nil {
		if len(b) < 2 {
			return 0, newRequestBufferOverflowError()
		}

		binary.LittleEndian.PutUint16(b, 0)
		return 2, nil
	}

	var msg string
	if len(err) == 1 {
		msg = err[0].Error()
	} else {
		msgs := make([]string, len(err))
		for i, err := range err {
			msgs[i] = strconv.Quote(err.Error())
		}

		msg = "multiple errors: " + strings.Join(msgs, ", ")
	}

	if len(msg) > math.MaxUint16 {
		i := 0
		for j := range msg {
			if j > math.MaxUint16 {
				break
			}

			i = j
		}

		msg = msg[:i]
	}

	size := len(msg) + 2
	if len(b) < size {
		return 0, newRequestBufferOverflowError()
	}

	binary.LittleEndian.PutUint16(b, uint16(len(msg)))
	copy(b[2:], msg)

	return size, nil
}

func putResponseStatusTooLong(b []byte) (int, error) {
	size := len(responseStatusTooLong) + 2
	if len(b) < size {
		return 0, newRequestBufferOverflowError()
	}

	binary.LittleEndian.PutUint16(b, uint16(len(responseStatusTooLong)))
	copy(b[2:], responseStatusTooLong)

	return size, nil
}

func putResponseObligationsTooLong(b []byte) (int, error) {
	size := len(responseStatusObligationsTooLong) + 2
	if len(b) < size {
		return 0, newRequestBufferOverflowError()
	}

	binary.LittleEndian.PutUint16(b, uint16(len(responseStatusObligationsTooLong)))
	copy(b[2:], responseStatusObligationsTooLong)

	return size, nil
}

func putAssignmentExpressions(b []byte, in []AttributeAssignment) (int, error) {
	off, err := putRequestAttributeCount(b, len(in))
	if err != nil {
		return off, err
	}

	for _, a := range in {
		id := a.a.id
		v, ok := a.e.(AttributeValue)
		if !ok {
			return off, newRequestInvalidExpressionError(a)
		}

		n, err := putRequestAttribute(b[off:], id, v)
		if err != nil {
			return off, bindError(err, id)
		}
		off += n
	}

	return off, nil
}

func putAttributesFromReflection(b []byte, c int, f func(i int) (string, Type, reflect.Value, error)) (int, error) {
	off, err := putRequestAttributeCount(b, c)
	if err != nil {
		return off, err
	}

	for i := 0; i < c; i++ {
		id, t, v, err := f(i)
		if err != nil {
			return off, err
		}

		var n int
		switch t {
		default:
			return off, bindError(newRequestAttributeMarshallingNotImplementedError(t), id)

		case TypeBoolean:
			n, err = putRequestAttributeBoolean(b[off:], id, v.Bool())

		case TypeString:
			n, err = putRequestAttributeString(b[off:], id, v.String())

		case TypeInteger:
			n, err = putRequestAttributeInteger(b[off:], id, v.Int())

		case TypeFloat:
			n, err = putRequestAttributeFloat(b[off:], id, v.Float())

		case TypeAddress:
			n, err = putRequestAttributeAddress(b[off:], id, net.IP(v.Bytes()))

		case TypeNetwork:
			n, err = putRequestAttributeNetwork(b[off:], id, getNetwork(v))

		case TypeDomain:
			n, err = putRequestAttributeDomain(b[off:], id, domain.MakeNameFromReflection(v))

		case TypeSetOfStrings:
			n, err = putRequestAttributeSetOfStrings(b[off:], id, getSetOfStrings(v))

		case TypeSetOfNetworks:
			n, err = putRequestAttributeSetOfNetworks(b[off:], id, getSetOfNetworks(v))

		case TypeSetOfDomains:
			n, err = putRequestAttributeSetOfDomains(b[off:], id, getSetOfDomains(v))

		case TypeListOfStrings:
			n, err = putRequestAttributeListOfStrings(b[off:], id, getListOfStrings(v))
		}

		if err != nil {
			return off, bindError(err, id)
		}
		off += n
	}

	return off, nil
}

func getAssignmentExpressions(b []byte) ([]AttributeAssignment, error) {
	c, n, err := getRequestAttributeCount(b)
	if err != nil {
		return nil, err
	}

	if c > 0 {
		b = b[n:]

		out := make([]AttributeAssignment, c)
		for i := range out {
			id, v, n, err := getRequestAttribute(b)
			if err != nil {
				return nil, bindErrorf(err, "%d", i+1)
			}
			b = b[n:]

			out[i] = MakeExpressionAssignment(id, v)
		}

		return out, nil
	}

	return nil, nil
}

func getAssignmentExpressionsWithAllocator(b []byte, f func(n int) ([]AttributeAssignment, error)) ([]AttributeAssignment, error) {
	c, n, err := getRequestAttributeCount(b)
	if err != nil {
		return nil, err
	}

	if c > 0 {
		b = b[n:]

		out, err := f(c)
		if err != nil {
			return nil, err
		}

		if len(out) < c {
			return nil, newRequestAssignmentsOverflowError(c, len(out))
		}
		out = out[:c]

		for i := range out {
			id, v, n, err := getRequestAttribute(b)
			if err != nil {
				return nil, bindErrorf(err, "%d", i+1)
			}
			b = b[n:]

			out[i] = MakeExpressionAssignment(id, v)
		}

		return out, nil
	}

	return nil, nil
}

func getAssignmentExpressionsToArray(b []byte, out []AttributeAssignment) (int, error) {
	c, n, err := getRequestAttributeCount(b)
	if err != nil {
		return 0, err
	}
	b = b[n:]

	if len(out) < c {
		return 0, newRequestAssignmentsOverflowError(c, len(out))
	}

	for i := 0; i < c; i++ {
		id, v, n, err := getRequestAttribute(b)
		if err != nil {
			return 0, bindErrorf(err, "%d", i+1)
		}
		b = b[n:]

		out[i] = MakeExpressionAssignment(id, v)
	}

	return c, nil
}

func getAttributesToReflection(b []byte, f func(string, Type) (reflect.Value, error)) error {
	c, n, err := getRequestAttributeCount(b)
	if err != nil {
		return err
	}
	b = b[n:]

	for i := 0; i < c; i++ {
		id, n, err := getRequestAttributeName(b)
		if err != nil {
			return bindErrorf(err, "%d", i+1)
		}
		b = b[n:]

		t, n, err := getRequestAttributeType(b)
		if err != nil {
			return bindError(err, id)
		}
		b = b[n:]

		if t < 0 || t >= len(builtinTypeByWire) {
			return bindError(newRequestAttributeUnmarshallingTypeError(t), id)
		}

		v, err := f(id, builtinTypeByWire[t])
		if err != nil {
			return err
		}

		switch t {
		case requestWireTypeBooleanFalse:
			err = setBool(v, false)

		case requestWireTypeBooleanTrue:
			err = setBool(v, true)

		case requestWireTypeString:
			var s string
			s, n, err = getRequestStringValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setString(v, s)

		case requestWireTypeInteger:
			var i int64
			i, n, err = getRequestIntegerValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setInt(v, i)

		case requestWireTypeFloat:
			var f float64
			f, n, err = getRequestFloatValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setFloat(v, f)

		case requestWireTypeIPv4Address:
			var a net.IP
			a, n, err = getRequestIPv4AddressValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setAddress(v, a)

		case requestWireTypeIPv6Address:
			var a net.IP
			a, n, err = getRequestIPv6AddressValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setAddress(v, a)

		case requestWireTypeIPv4Network:
			var ip *net.IPNet
			ip, n, err = getRequestIPv4NetworkValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setNetwork(v, ip)

		case requestWireTypeIPv6Network:
			var ip *net.IPNet
			ip, n, err = getRequestIPv6NetworkValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setNetwork(v, ip)

		case requestWireTypeDomain:
			var d domain.Name
			d, n, err = getRequestDomainValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setDomain(v, d)

		case requestWireTypeSetOfStrings:
			var ss *strtree.Tree
			ss, n, err = getRequestSetOfStringsValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setSetOfStrings(v, ss)

		case requestWireTypeSetOfNetworks:
			var sn *iptree.Tree
			sn, n, err = getRequestSetOfNetworksValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setSetOfNetworks(v, sn)

		case requestWireTypeSetOfDomains:
			var sd *domaintree.Node
			sd, n, err = getRequestSetOfDomainsValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setSetOfDomains(v, sd)

		case requestWireTypeListOfStrings:
			var ls []string
			ls, n, err = getRequestListOfStringsValue(b)
			if err != nil {
				return bindError(err, id)
			}
			b = b[n:]

			err = setListOfStrings(v, ls)
		}

		if err != nil {
			return bindError(err, id)
		}
	}

	return nil
}

func calcResponseSize(obligations []AttributeAssignment, errs ...error) (int, error) {
	s, err := calcAssignmentExpressionsSize(obligations)
	if err != nil {
		return 0, err
	}

	return reqVersionSize + resEffectSize + calcResponseStatus(errs...) + s, nil
}

func calcResponseStatus(err ...error) int {
	if len(err) < 1 || len(err) == 1 && err[0] == nil {
		return reqBigCounterSize
	}

	var msg string
	if len(err) == 1 {
		msg = err[0].Error()
	} else {
		msgs := make([]string, len(err))
		for i, err := range err {
			msgs[i] = strconv.Quote(err.Error())
		}

		msg = "multiple errors: " + strings.Join(msgs, ", ")
	}

	if len(msg) > math.MaxUint16 {
		i := 0
		for j := range msg {
			if j > math.MaxUint16 {
				break
			}

			i = j
		}

		msg = msg[:i]
	}

	return len(msg) + reqBigCounterSize
}
func calcAssignmentExpressionsSize(in []AttributeAssignment) (int, error) {
	s := reqBigCounterSize

	for _, a := range in {
		id := a.a.id
		v, ok := a.e.(AttributeValue)
		if !ok {
			return 0, newRequestInvalidExpressionError(a)
		}

		n, err := calcRequestAttributeNameSize(id)
		if err != nil {
			return 0, bindError(err, id)
		}
		s += n

		n, err = calcRequestAttributeSize(v)
		if err != nil {
			return 0, bindError(err, id)
		}

		s += n
	}

	return s, nil
}

func calcAttributesSizeFromReflection(c int, f func(i int) (string, Type, reflect.Value, error)) (int, error) {
	s := reqBigCounterSize

	for i := 0; i < c; i++ {
		id, t, v, err := f(i)
		if err != nil {
			return 0, err
		}

		n, err := calcRequestAttributeNameSize(id)
		if err != nil {
			return 0, bindError(err, id)
		}
		s += n

		switch t {
		default:
			return 0, bindError(newRequestAttributeMarshallingNotImplementedError(t), id)

		case TypeBoolean:
			n = reqBooleanValueSize

		case TypeString:
			n, err = calcRequestAttributeStringSize(v.String())

		case TypeInteger:
			n, err = calcRequestAttributeIntegerSize(v.Int())

		case TypeFloat:
			n, err = calcRequestAttributeFloatSize(v.Float())

		case TypeAddress:
			n, err = calcRequestAttributeAddressSize(net.IP(v.Bytes()))

		case TypeNetwork:
			n, err = calcRequestAttributeNetworkSize(getNetwork(v))

		case TypeDomain:
			n, err = calcRequestAttributeDomainSize(domain.MakeNameFromReflection(v))

		case TypeSetOfStrings:
			n, err = calcRequestAttributeSetOfStringsSize(getSetOfStrings(v))

		case TypeSetOfNetworks:
			n, err = calcRequestAttributeSetOfNetworksSize(getSetOfNetworks(v))

		case TypeSetOfDomains:
			n, err = calcRequestAttributeSetOfDomainsSize(getSetOfDomains(v))

		case TypeListOfStrings:
			n, err = calcRequestAttributeListOfStringsSize(getListOfStrings(v))
		}

		if err != nil {
			return 0, bindError(err, id)
		}
		s += reqTypeSize + n
	}

	return s, nil
}

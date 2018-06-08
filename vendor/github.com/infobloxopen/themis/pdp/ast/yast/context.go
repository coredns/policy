package yast

import (
	"math"

	"github.com/infobloxopen/themis/pdp"
)

type context struct {
	symbols pdp.Symbols
}

func newContext() *context {
	return &context{
		symbols: pdp.MakeSymbols(),
	}
}

func newContextWithSymbols(symbols pdp.Symbols) *context {
	return &context{
		symbols: symbols,
	}
}

func (ctx context) validateString(v interface{}, desc string) (string, boundError) {
	r, ok := v.(string)
	if !ok {
		return "", newStringError(v, desc)
	}

	return r, nil
}

func (ctx context) extractString(m map[interface{}]interface{}, k string, desc string) (string, boundError) {
	v, ok := m[k]
	if !ok {
		return "", newMissingStringError(desc)
	}

	return ctx.validateString(v, desc)
}

func (ctx context) extractStringOpt(m map[interface{}]interface{}, k string, desc string) (string, bool, boundError) {
	v, ok := m[k]
	if !ok {
		return "", false, nil
	}

	s, err := ctx.validateString(v, desc)
	return s, true, err
}

func (ctx context) validateInteger(v interface{}, desc string) (int64, boundError) {
	switch v := v.(type) {
	case int:
		return int64(v), nil

	case int64:
		return v, nil

	case uint64:
		if v > math.MaxInt64 {
			return 0, newIntegerUint64OverflowError(v, desc)
		}

		return int64(v), nil

	case float64:
		if v < -9007199254740992 || v > 9007199254740992 {
			return 0, newIntegerFloat64OverflowError(v, desc)
		}

		return int64(v), nil
	}

	return 0, newIntegerError(v, desc)
}

func (ctx context) validateFloat(v interface{}, desc string) (float64, boundError) {
	switch v := v.(type) {
	case int:
		return float64(v), nil

	case int64:
		return float64(v), nil

	case uint64:
		return float64(v), nil

	case float64:
		return float64(v), nil
	}

	return 0, newFloatError(v, desc)
}

func (ctx context) validateMap(v interface{}, desc string) (map[interface{}]interface{}, boundError) {
	r, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, newMapError(v, desc)
	}

	return r, nil
}

func (ctx context) extractMap(m map[interface{}]interface{}, k string, desc string) (map[interface{}]interface{}, boundError) {
	v, ok := m[k]
	if !ok {
		return nil, newMissingMapError(desc)
	}

	return ctx.validateMap(v, desc)
}

func (ctx context) extractMapOpt(m map[interface{}]interface{}, k string, desc string) (map[interface{}]interface{}, bool, boundError) {
	v, ok := m[k]
	if !ok {
		return nil, false, nil
	}

	m, err := ctx.validateMap(v, desc)
	return m, true, err
}

func (ctx context) validateList(v interface{}, desc string) ([]interface{}, boundError) {
	r, ok := v.([]interface{})
	if !ok {
		return nil, newListError(v, desc)
	}

	return r, nil
}

func (ctx context) extractList(m map[interface{}]interface{}, k, desc string) ([]interface{}, boundError) {
	v, ok := m[k]
	if !ok {
		return nil, newMissingListError(desc)
	}

	return ctx.validateList(v, desc)
}

func (ctx context) extractListOpt(m map[interface{}]interface{}, k, desc string) ([]interface{}, bool, boundError) {
	v, ok := m[k]
	if !ok {
		return nil, false, nil
	}

	l, err := ctx.validateList(v, desc)
	return l, true, err
}

func (ctx context) getSingleMapPair(m map[interface{}]interface{}, desc string) (interface{}, interface{}, boundError) {
	if len(m) > 1 {
		return nil, nil, newTooManySMPItemsError(desc, len(m))
	}

	for k, v := range m {
		return k, v, nil
	}

	return nil, nil, newNoSMPItemsError(desc, len(m))
}

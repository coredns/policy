// Package jparser provides helper methods to parse JSON from stream.
package jparser

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	// DelimObjectStart contains json object '{' start delimiter.
	DelimObjectStart = "{"
	// DelimObjectEnd contains json object '}' start delimiter.
	DelimObjectEnd = "}"
	// DelimArrayStart contains json array '[' start delimiter.
	DelimArrayStart = "["
	// DelimArrayEnd contains json array '}' end delimiter.
	DelimArrayEnd = "]"
)

// CheckRootObjectStart checks whether JSON byte stream starts with '{'
func CheckRootObjectStart(d *json.Decoder) (bool, error) {
	t, err := d.Token()
	if err == io.EOF {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	delim, ok := t.(json.Delim)
	if !ok {
		return false, newRootObjectStartTokenError(t, DelimObjectStart)
	}

	if delim.String() != DelimObjectStart {
		return false, newRootObjectStartDelimiterError(delim, DelimObjectStart)
	}

	return true, nil
}

// CheckRootArrayStart checks whether JSON byte stream starts with '['
func CheckRootArrayStart(d *json.Decoder) (bool, error) {
	t, err := d.Token()
	if err == io.EOF {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	delim, ok := t.(json.Delim)
	if !ok {
		return false, newRootArrayStartTokenError(t, DelimArrayStart)
	}

	if delim.String() != DelimArrayStart {
		return false, newRootArrayStartDelimiterError(delim, DelimArrayStart)
	}

	return true, nil
}

// CheckObjectStart checks whether next token is '{' in JSON byte stream.
func CheckObjectStart(d *json.Decoder, desc string) error {
	t, err := d.Token()
	if err != nil {
		return err
	}

	delim, ok := t.(json.Delim)
	if !ok {
		return newObjectStartTokenError(t, DelimObjectStart, desc)
	}

	if delim.String() != DelimObjectStart {
		return newObjectStartDelimiterError(delim, DelimObjectStart, desc)
	}

	return nil
}

// CheckArrayStart checks whether next token is '[' in JSON byte stream.
func CheckArrayStart(d *json.Decoder, desc string) error {
	t, err := d.Token()
	if err != nil {
		return err
	}

	delim, ok := t.(json.Delim)
	if !ok {
		return newArrayStartTokenError(t, DelimArrayStart, desc)
	}

	if delim.String() != DelimArrayStart {
		return newArrayStartDelimiterError(delim, DelimArrayStart, desc)
	}

	return nil
}

// CheckObjectArrayStart checks whether next token is '[' or "{" in JSON byte stream.
// It returns true if it is "{" if there wasn't any error.
func CheckObjectArrayStart(d *json.Decoder, desc string) (bool, error) {
	t, err := d.Token()
	if err != nil {
		return false, err
	}

	delim, ok := t.(json.Delim)
	if !ok {
		return false, newObjectArrayStartTokenError(t, DelimObjectStart, DelimArrayStart, desc)
	}

	switch delim.String() {
	case DelimObjectStart:
		return true, nil

	case DelimArrayStart:
		return false, nil
	}

	return false, newObjectArrayStartDelimiterError(delim, DelimObjectStart, DelimArrayStart, desc)
}

// CheckEOF checks EOF of the JSON byte stream.
func CheckEOF(d *json.Decoder) error {
	t, err := d.Token()
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	return newMissingEOFError(t)
}

// SkipValue skips object or array in JSON byte stream.
func SkipValue(d *json.Decoder, desc string) error {
	t, err := d.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); ok {
		s := delim.String()
		switch s {
		default:
			return newUnexpectedDelimiterError(s, desc)

		case DelimObjectStart:
			return SkipObject(d, desc)

		case DelimArrayStart:
			return SkipArray(d, desc)
		}
	}

	return nil
}

// SkipObject skips object in JSON byte stream.
func SkipObject(d *json.Decoder, desc string) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch t := t.(type) {
		default:
			return newObjectTokenError(t, DelimObjectEnd, desc)

		case string:
			err := SkipValue(d, desc)
			if err != nil {
				return bindError(err, t)
			}

		case json.Delim:
			if t.String() != DelimObjectEnd {
				return newObjectEndDelimiterError(t, DelimObjectEnd, desc)
			}

			return nil
		}
	}
}

// SkipArray skips array in JSON byte stream.
func SkipArray(d *json.Decoder, desc string) error {
	i := 1
	for {
		src := fmt.Sprintf("%d", i)

		t, err := d.Token()
		if err != nil {
			return bindError(err, src)
		}

		if delim, ok := t.(json.Delim); ok {
			s := delim.String()
			switch s {
			default:
				return bindError(newUnexpectedDelimiterError(s, desc), src)

			case DelimArrayEnd:
				return nil

			case DelimObjectStart:
				err := SkipObject(d, desc)
				if err != nil {
					return bindError(err, src)
				}

			case DelimArrayStart:
				err := SkipArray(d, desc)
				if err != nil {
					return bindError(err, src)
				}
			}
		}

		i++
	}
}

// Pair represents unmarshalled part of JSON byte stream.
// Value is an array of Pairs or primitive value or []interface{}.
type Pair struct {
	K string
	V interface{}
}

// GetUndefined unmarshals whole part of JSON byte stream.
// Part is an object or an array or primitive value.
// It returns array of Pairs or primitive value or []interface{}.
func GetUndefined(d *json.Decoder, desc string) (interface{}, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}

	switch t := t.(type) {
	case json.Delim:
		s := t.String()
		switch s {
		case DelimObjectStart:
			return GetObject(d, desc)

		case DelimArrayStart:
			return GetArray(d, desc)
		}

		return nil, newUnexpectedDelimiterError(s, desc)

	case bool:
		return t, nil

	case float64:
		return t, nil

	case json.Number:
		return t, nil

	case string:
		return t, nil
	}

	return t, nil
}

// GetObject unmarshals whole object from JSON byte stream to a list of Piars.
func GetObject(d *json.Decoder, desc string) ([]Pair, error) {
	obj := []Pair{}

	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}

		switch t := t.(type) {
		default:
			return nil, newObjectTokenError(t, DelimObjectEnd, desc)

		case string:
			v, err := GetUndefined(d, desc)
			if err != nil {
				return nil, bindError(err, t)
			}

			obj = append(obj, Pair{K: t, V: v})

		case json.Delim:
			if t.String() != DelimObjectEnd {
				return nil, newObjectEndDelimiterError(t, DelimObjectEnd, desc)
			}

			return obj, nil
		}
	}
}

// GetArray unmarshals whole array from JSON byte stream to an []interface{}.
// []interface{} item is an array of Pairs or primitive value or []interface{}.
func GetArray(d *json.Decoder, desc string) ([]interface{}, error) {
	arr := []interface{}{}
	i := 1
	for {
		src := fmt.Sprintf("%d", i)

		t, err := d.Token()
		if err != nil {
			return nil, bindError(err, src)
		}

		if delim, ok := t.(json.Delim); ok {
			s := delim.String()
			switch s {
			default:
				return nil, bindError(newUnexpectedDelimiterError(s, desc), src)

			case DelimArrayEnd:
				return arr, nil

			case DelimObjectStart:
				v, err := GetObject(d, desc)
				if err != nil {
					return nil, bindError(err, src)
				}

				arr = append(arr, v)

			case DelimArrayStart:
				v, err := GetArray(d, desc)
				if err != nil {
					return nil, bindError(err, src)
				}

				arr = append(arr, v)
			}
		} else {
			arr = append(arr, t)
		}

		i++
	}
}

// GetBoolean unmarshals boolean value from JSON byte stream.
func GetBoolean(d *json.Decoder, desc string) (bool, error) {
	t, err := d.Token()
	if err != nil {
		return false, err
	}

	b, ok := t.(bool)
	if !ok {
		return false, newBooleanCastError(t, desc)
	}

	return b, nil
}

// GetString unmarshals string value from JSON byte stream.
func GetString(d *json.Decoder, desc string) (string, error) {
	t, err := d.Token()
	if err != nil {
		return "", err
	}

	s, ok := t.(string)
	if !ok {
		return "", newStringCastError(t, desc)
	}

	return s, nil
}

// GetNumber unmarshals number from JSON byte stream.
func GetNumber(d *json.Decoder, desc string) (float64, error) {
	t, err := d.Token()
	if err != nil {
		return 0, err
	}

	n, ok := t.(float64)
	if !ok {
		return 0, newNumberCastError(t, desc)
	}

	return n, nil
}

// GetStringSequence iterates over object keys or string array items in JSON byte stream.
func GetStringSequence(d *json.Decoder, f func(idx int, s string) error, desc string) error {
	ok, err := CheckObjectArrayStart(d, desc)
	if err != nil {
		return err
	}

	if ok {
		return GetStringSequenceFromObject(d, f, desc)
	}

	return GetStringSequenceFromArray(d, f, desc)
}

// GetStringSequenceFromObject iterates over object keys in JSON byte stream.
func GetStringSequenceFromObject(d *json.Decoder, f func(idx int, s string) error, desc string) error {
	i := 1
	for {
		t, err := d.Token()
		if err != nil {
			return bindErrorf(err, "%d", i)
		}

		switch t := t.(type) {
		default:
			return bindErrorf(newObjectTokenError(t, DelimObjectEnd, desc), "%d", i)

		case string:
			err := f(i, t)
			if err != nil {
				return err
			}

			err = SkipValue(d, desc)
			if err != nil {
				return bindErrorf(err, "%d", i)
			}

		case json.Delim:
			if t.String() != DelimObjectEnd {
				return bindErrorf(newObjectEndDelimiterError(t, DelimObjectEnd, desc), "%d", i)
			}

			return nil
		}

		i++
	}
}

// GetStringSequenceFromArray iterates over string array items in JSON byte stream.
func GetStringSequenceFromArray(d *json.Decoder, f func(idx int, s string) error, desc string) error {
	i := 1
	for {
		t, err := d.Token()
		if err != nil {
			return bindErrorf(err, "%d", i)
		}

		switch t := t.(type) {
		default:
			return bindErrorf(newStringArrayTokenError(t, DelimArrayEnd, desc), "%d", i)

		case string:
			err := f(i, t)
			if err != nil {
				return err
			}

		case json.Delim:
			if t.String() != DelimArrayEnd {
				return bindErrorf(newArrayEndDelimiterError(t, DelimArrayEnd, desc), "%d", i)
			}

			return nil
		}

		i++
	}
}

// UnmarshalObject sequentially unmarshals object from JSON byte stream.
func UnmarshalObject(d *json.Decoder, u func(key string, d *json.Decoder) error, desc string) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch t := t.(type) {
		default:
			return newObjectTokenError(t, DelimObjectEnd, desc)

		case string:
			err = u(t, d)
			if err != nil {
				return err
			}

		case json.Delim:
			if t.String() != DelimObjectEnd {
				return newObjectEndDelimiterError(t, DelimObjectEnd, desc)
			}

			return nil
		}
	}
}

// UnmarshalObjectArray sequentially unmarshals array from JSON byte stream.
func UnmarshalObjectArray(d *json.Decoder, u func(idx int, d *json.Decoder) error, desc string) error {
	i := 1
	for {
		src := fmt.Sprintf("%d", i)

		t, err := d.Token()
		if err != nil {
			return bindError(err, src)
		}

		delim, ok := t.(json.Delim)
		if !ok {
			return bindError(newObjectArrayTokenError(t, DelimArrayEnd, desc), src)
		}

		s := delim.String()
		switch s {
		default:
			return bindError(newUnexpectedObjectArrayDelimiterError(s, desc), src)

		case DelimArrayEnd:
			return nil

		case DelimObjectStart:
			err := u(i, d)
			if err != nil {
				return err
			}
		}

		i++
	}
}

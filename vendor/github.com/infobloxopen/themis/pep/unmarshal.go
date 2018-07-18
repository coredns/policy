package pep

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

const (
	effectFieldName = "Effect"
	reasonFieldName = "Reason"
)

var (
	// ErrorInvalidDestination indicates that output value of validate method is
	// not a structure.
	ErrorInvalidDestination = errors.New("given value is not a pointer to structure")
)

type resFieldsInfo struct {
	fields map[string]string
	err    error
}

var (
	resTypeCache     = map[string]resFieldsInfo{}
	resTypeCacheLock = sync.RWMutex{}

	specialNameByID = map[string]string{
		pdp.ResponseEffectFieldName: effectFieldName,
		pdp.ResponseStatusFieldName: reasonFieldName,
	}
)

func fillResponse(res pb.Msg, v interface{}) error {
	switch v := v.(type) {
	case *pb.Msg:
		*v = res
		return nil

	case *pdp.Response:
		effect, n, err := pdp.UnmarshalResponse(res.Body, v.Obligations)
		if err != nil {
			if _, ok := err.(*pdp.ResponseServerError); !ok {
				return err
			}
		}

		v.Effect = effect
		v.Status = err
		if v.Obligations != nil {
			v.Obligations = v.Obligations[:n]
		}

		return nil
	}

	return unmarshalToValue(res.Body, reflect.ValueOf(v))
}

func unmarshalToValue(res []byte, v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return ErrorInvalidDestination
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrorInvalidDestination
	}

	fields, err := makeFieldMap(v.Type())
	if err != nil {
		return err
	}

	if len(fields) > 0 {
		return unmarshalToTaggedStruct(res, v, fields)
	}

	return unmarshalToUntaggedStruct(res, v)
}

func parseTag(tag string, f reflect.StructField, t reflect.Type) (string, error) {
	items := strings.Split(tag, ",")
	if len(items) > 1 {
		tag = items[0]
		taggedTypeName := items[1]

		if tag == effectFieldName || tag == reasonFieldName {
			return "", fmt.Errorf("don't support type definition for \"%s\" and \"%s\" fields (%s.%s)",
				effectFieldName, reasonFieldName, t.Name(), f.Name)
		}

		taggedTypes, ok := typeByTag[strings.ToLower(taggedTypeName)]
		if !ok {
			return "", fmt.Errorf("unknown type \"%s\" (%s.%s)", taggedTypeName, t.Name(), f.Name)
		}

		if _, ok := taggedTypes[f.Type]; !ok {
			return "", fmt.Errorf("tagged type \"%s\" doesn't match field type \"%s\" (%s.%s)",
				taggedTypeName, f.Type.Name(), t.Name(), f.Name)
		}

		return tag, nil
	}

	return tag, nil
}

func makeFieldMap(t reflect.Type) (map[string]string, error) {
	key := t.PkgPath() + "." + t.Name()
	resTypeCacheLock.RLock()
	if info, ok := resTypeCache[key]; ok {
		resTypeCacheLock.RUnlock()
		return info.fields, info.err
	}
	resTypeCacheLock.RUnlock()

	m := make(map[string]string)
	var err error
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag, ok := getTag(f)
		if !ok {
			continue
		}

		if len(tag) <= 0 {
			tag, ok = getName(f)
			if !ok {
				continue
			}
		}

		tag, err = parseTag(tag, f, t)
		if err != nil {
			break
		}

		m[tag] = f.Name
	}

	resTypeCacheLock.Lock()
	resTypeCache[key] = resFieldsInfo{
		fields: m,
		err:    err,
	}
	resTypeCacheLock.Unlock()

	return m, err
}

func unmarshalToTaggedStruct(res []byte, v reflect.Value, fields map[string]string) error {
	return pdp.UnmarshalResponseToReflection(res, func(id string, t pdp.Type) (reflect.Value, error) {
		if t == nil {
			name, ok := specialNameByID[id]
			if !ok {
				return reflect.ValueOf(nil), fmt.Errorf("unknown id %q", id)
			}

			name, ok = fields[name]
			if !ok {
				return reflect.ValueOf(nil), nil
			}

			return v.FieldByName(name), nil
		}

		name, ok := fields[id]
		if !ok {
			return reflect.ValueOf(nil), nil
		}

		f := v.FieldByName(name)
		if !f.CanSet() {
			return reflect.ValueOf(nil), fmt.Errorf("field %s.%s is tagged but can't be set", v.Type().Name(), name)
		}

		types, ok := typeByAttrType[t]
		if !ok {
			return reflect.ValueOf(nil), fmt.Errorf("can't unmarshal \"%s\" of \"%s\" type", id, t)
		}

		if _, ok := types[f.Type()]; !ok {
			return reflect.ValueOf(nil), fmt.Errorf("can't unmarshal \"%s\" of \"%s\" type to field %s.%s",
				id, t, v.Type().Name(), name)
		}

		return f, nil
	})
}

func unmarshalToUntaggedStruct(res []byte, v reflect.Value) error {
	return pdp.UnmarshalResponseToReflection(res, func(id string, t pdp.Type) (reflect.Value, error) {
		if t == nil {
			name, ok := specialNameByID[id]
			if !ok {
				return reflect.ValueOf(nil), fmt.Errorf("unknown id %q", id)
			}

			f := v.FieldByName(name)
			if !f.CanSet() {
				return reflect.ValueOf(nil), nil
			}

			k := f.Kind()
			switch id {
			case pdp.ResponseEffectFieldName:
				if k != reflect.Bool &&
					k != reflect.Int &&
					k != reflect.Int8 &&
					k != reflect.Int16 &&
					k != reflect.Int32 &&
					k != reflect.Int64 &&
					k != reflect.Uint &&
					k != reflect.Uint8 &&
					k != reflect.Uint16 &&
					k != reflect.Uint32 &&
					k != reflect.Uint64 &&
					k != reflect.String {
					return reflect.ValueOf(nil), nil
				}

			case pdp.ResponseStatusFieldName:
				if k != reflect.String {
					t := f.Type()
					if t.PkgPath() != "" || t.Name() != "error" {
						return reflect.ValueOf(nil), nil
					}
				}
			}

			return f, nil
		}

		f := v.FieldByName(id)
		if !f.CanSet() {
			return reflect.ValueOf(nil), nil
		}

		types, ok := typeByAttrType[t]
		if !ok {
			return reflect.ValueOf(nil), nil
		}

		if _, ok := types[f.Type()]; !ok {
			return reflect.ValueOf(nil), nil
		}

		return f, nil
	})
}

package jast

import (
	"encoding/json"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx *context) unmarshalTypeDeclaration(ID string, d *json.Decoder) error {
	if err := jparser.CheckObjectStart(d, "type declaration"); err != nil {
		return err
	}

	var (
		metaOk bool
		meta   string
		flags  []string
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		switch k {
		default:
			return newUnknownFieldError(k)

		case yastTagMeta:
			s, err := jparser.GetString(d, "meta type name")
			if err != nil {
				return err
			}

			meta = s
			metaOk = true

		case yastTagFlags:
			flags = []string{}
			if err := jparser.GetStringSequence(d, func(i int, s string) error {
				flags = append(flags, s)
				return nil
			}, "list of flag names"); err != nil {
				return err
			}
		}

		return nil
	}, "type declarations"); err != nil {
		return err
	}

	if !metaOk {
		return newMissingMetaTypeNameError()
	}

	switch strings.ToLower(meta) {
	default:
		return newUnknownMetaTypeError(meta)

	case yastTagFlags:
		if flags == nil {
			return newMissingFlagNameListError()
		}

		t, err := pdp.NewFlagsType(ID, flags...)
		if err != nil {
			return err
		}

		if err := ctx.symbols.PutType(t); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *context) unmarshalTypeDeclarations(d *json.Decoder) error {
	if err := jparser.CheckObjectStart(d, "type declarations"); err != nil {
		return bindError(err, yastTagTypes)
	}

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		if err := ctx.unmarshalTypeDeclaration(k, d); err != nil {
			return bindError(err, k)
		}

		return nil
	}, "type declarations"); err != nil {
		return bindError(err, yastTagTypes)
	}

	return nil
}

package jast

import (
	"encoding/json"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

type context struct {
	symbols    pdp.Symbols
	rootPolicy pdp.Evaluable
}

func newContext() *context {
	return &context{
		symbols: pdp.MakeSymbols(),
	}
}

func newContextWithSymbols(s pdp.Symbols) *context {
	return &context{
		symbols: s,
	}
}

func (ctx *context) unmarshal(d *json.Decoder) error {
	ok, err := jparser.CheckRootObjectStart(d)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	if err = jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		switch strings.ToLower(k) {
		case yastTagTypes:
			return ctx.unmarshalTypeDeclarations(d)

		case yastTagAttributes:
			return ctx.unmarshalAttributeDeclarations(d)

		case yastTagPolicies:
			return ctx.unmarshalRootPolicy(d)
		}

		return newUnknownFieldError(k)
	}, "root"); err != nil {
		return err
	}

	return jparser.CheckEOF(d)
}

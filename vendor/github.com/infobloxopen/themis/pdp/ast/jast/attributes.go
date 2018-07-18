package jast

import (
	"encoding/json"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx *context) unmarshalAttributeDeclarations(d *json.Decoder) boundError {
	err := jparser.CheckObjectStart(d, "attribute declarations")
	if err != nil {
		return bindError(err, yastTagAttributes)
	}

	if err = jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		tstr, err := jparser.GetString(d, "attribute data type")
		if err != nil {
			return err
		}

		t := ctx.symbols.GetType(tstr)
		if t == nil {
			return bindError(newUnknownTypeError(tstr), k)
		}

		if err := ctx.symbols.PutAttribute(pdp.MakeAttribute(k, t)); err != nil {
			return bindError(err, k)
		}

		return nil
	}, "attribute declarations"); err != nil {
		return bindError(err, yastTagAttributes)
	}

	return nil
}

package jast

import (
	"encoding/json"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx context) unmarshalAttributeDesignator(d *json.Decoder) (pdp.AttributeDesignator, error) {
	id, err := jparser.GetString(d, "attribute ID")
	if err != nil {
		return pdp.AttributeDesignator{}, err
	}

	a, ok := ctx.symbols.GetAttribute(id)
	if !ok {
		return pdp.AttributeDesignator{}, newUnknownAttributeError(id)
	}

	return pdp.MakeAttributeDesignator(a), nil
}

func (ctx context) unmarshalArguments(d *json.Decoder) ([]pdp.Expression, error) {
	err := jparser.CheckArrayStart(d, "arguments")
	if err != nil {
		return nil, err
	}

	args := []pdp.Expression{}
	if err := jparser.UnmarshalObjectArray(d, func(idx int, d *json.Decoder) error {
		arg, err := ctx.unmarshalExpression(d)
		if err != nil {
			return bindErrorf(err, "%d", idx)
		}

		args = append(args, arg)

		return nil
	}, "arguments"); err != nil {
		return nil, err
	}

	return args, nil
}

func (ctx context) unmarshalExpression(d *json.Decoder) (pdp.Expression, error) {
	var expr pdp.Expression

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var err error

		switch strings.ToLower(k) {
		case yastTagAttribute:
			expr, err = ctx.unmarshalAttributeDesignator(d)
			return err

		case yastTagValue:
			expr, err = ctx.unmarshalValue(d)
			return err

		case yastTagSelector:
			expr, err = ctx.unmarshalSelector(d)
			return err

		default:
			validators, ok := pdp.FunctionArgumentValidators[k]
			if !ok {
				return newUnknownFunctionError(k)
			}

			args, err := ctx.unmarshalArguments(d)
			if err != nil {
				return bindError(err, k)
			}

			for _, validator := range validators {
				if maker := validator(args); maker != nil {
					expr = maker(args)
					return nil
				}
			}

			return newFunctionCastError(k, args)
		}
	}, "expression"); err != nil {
		return nil, err
	}

	return expr, nil
}

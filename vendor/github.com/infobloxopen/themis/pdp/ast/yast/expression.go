package yast

import "github.com/infobloxopen/themis/pdp"

func (ctx context) unmarshalAttributeDesignator(v interface{}) (pdp.AttributeDesignator, boundError) {
	ID, err := ctx.validateString(v, "attribute ID")
	if err != nil {
		return pdp.AttributeDesignator{}, err
	}

	a, ok := ctx.symbols.GetAttribute(ID)
	if !ok {
		return pdp.AttributeDesignator{}, newUnknownAttributeError(ID)
	}

	return pdp.MakeAttributeDesignator(a), nil
}

func (ctx context) unmarshalArguments(v interface{}) ([]pdp.Expression, boundError) {
	items, err := ctx.validateList(v, "arguments")
	if err != nil {
		return nil, err
	}

	args := make([]pdp.Expression, len(items))
	for i, item := range items {
		arg, err := ctx.unmarshalExpression(item)
		if err != nil {
			return nil, bindErrorf(err, "%d", i)
		}

		args[i] = arg
	}

	return args, nil
}

func (ctx context) unmarshalExpression(v interface{}) (pdp.Expression, boundError) {
	m, err := ctx.validateMap(v, "expression")
	if err != nil {
		return nil, err
	}

	k, v, err := ctx.getSingleMapPair(m, "expression map")
	if err != nil {
		return nil, err
	}

	ID, err := ctx.validateString(k, "specificator or function name")
	if err != nil {
		return nil, err
	}

	switch ID {
	case yastTagAttribute:
		return ctx.unmarshalAttributeDesignator(v)

	case yastTagValue:
		return ctx.unmarshalValue(v)

	case yastTagSelector:
		return ctx.unmarshalSelector(v)
	}

	validators, ok := pdp.FunctionArgumentValidators[ID]
	if !ok {
		return nil, newUnknownFunctionError(ID)
	}

	args, err := ctx.unmarshalArguments(v)
	if err != nil {
		return nil, bindError(err, ID)
	}

	for _, validator := range validators {
		if maker := validator(args); maker != nil {
			return maker(args), nil
		}
	}

	return nil, newFunctionCastError(ID, args)
}

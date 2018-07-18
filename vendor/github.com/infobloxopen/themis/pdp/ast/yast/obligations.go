package yast

import "github.com/infobloxopen/themis/pdp"

func (ctx context) unmarshalObligationItem(v interface{}) (pdp.AttributeAssignment, boundError) {
	m, err := ctx.validateMap(v, "obligation")
	if err != nil {
		return pdp.AttributeAssignment{}, err
	}

	k, v, err := ctx.getSingleMapPair(m, "obligation")
	if err != nil {
		return pdp.AttributeAssignment{}, err
	}

	ID, err := ctx.validateString(k, "obligation attribute id")
	if err != nil {
		return pdp.AttributeAssignment{}, err
	}

	a, ok := ctx.symbols.GetAttribute(ID)
	if !ok {
		return pdp.AttributeAssignment{}, newUnknownAttributeError(ID)
	}

	m, err = ctx.validateMap(v, "obligation assignment")
	if err != nil {
		return pdp.AttributeAssignment{}, bindError(err, ID)
	}

	e, err := ctx.unmarshalExpression(v)
	if err != nil {
		return pdp.AttributeAssignment{}, bindError(err, ID)
	}

	return pdp.MakeAttributeAssignment(a, e), nil
}

func (ctx context) unmarshalObligations(m map[interface{}]interface{}) ([]pdp.AttributeAssignment, boundError) {
	items, ok, err := ctx.extractListOpt(m, yastTagObligation, "obligations")
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	var r []pdp.AttributeAssignment
	for i, item := range items {
		o, err := ctx.unmarshalObligationItem(item)
		if err != nil {
			return nil, bindError(bindErrorf(err, "%d", i), "obligations")
		}

		r = append(r, o)
	}

	return r, nil
}

package yast

import "github.com/infobloxopen/themis/pdp"

func (ctx context) unmarshalObligationItem(v interface{}) (pdp.AttributeAssignmentExpression, boundError) {
	m, err := ctx.validateMap(v, "obligation")
	if err != nil {
		return pdp.AttributeAssignmentExpression{}, err
	}

	k, v, err := ctx.getSingleMapPair(m, "obligation")
	if err != nil {
		return pdp.AttributeAssignmentExpression{}, err
	}

	ID, err := ctx.validateString(k, "obligation attribute id")
	if err != nil {
		return pdp.AttributeAssignmentExpression{}, err
	}

	a, ok := ctx.symbols.GetAttribute(ID)
	if !ok {
		return pdp.AttributeAssignmentExpression{}, newUnknownAttributeError(ID)
	}

	m, err = ctx.validateMap(v, "obligation assignment")
	if err != nil {
		return pdp.AttributeAssignmentExpression{}, bindError(err, ID)
	}

	e, err := ctx.unmarshalExpression(v)
	if err != nil {
		return pdp.AttributeAssignmentExpression{}, bindError(err, ID)
	}

	return pdp.MakeAttributeAssignmentExpression(a, e), nil
}

func (ctx context) unmarshalObligations(m map[interface{}]interface{}) ([]pdp.AttributeAssignmentExpression, boundError) {
	items, ok, err := ctx.extractListOpt(m, yastTagObligation, "obligations")
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	var r []pdp.AttributeAssignmentExpression
	for i, item := range items {
		o, err := ctx.unmarshalObligationItem(item)
		if err != nil {
			return nil, bindError(bindErrorf(err, "%d", i), "obligations")
		}

		r = append(r, o)
	}

	return r, nil
}

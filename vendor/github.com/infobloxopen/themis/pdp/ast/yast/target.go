package yast

import (
	"strings"

	"github.com/infobloxopen/themis/pdp"
)

func (ctx context) getTargetCompatibleArgument(v interface{}) (pdp.Expression, int, boundError) {
	e, err := ctx.unmarshalExpression(v)
	if err != nil {
		return nil, 0, err
	}

	TCAID, ok := pdp.CheckExpressionAsTargetArgument(e)
	if !ok {
		return nil, 0, newInvalidMatchFunctionArgError(e)
	}

	return e, TCAID, nil
}

func (ctx context) getAdjustedArgumentPair(v interface{}) (pdp.Expression, pdp.Expression, boundError) {
	args, err := ctx.validateList(v, "function arguments")
	if len(args) != 2 {
		return nil, nil, newMatchFunctionArgsNumberError(len(args))
	}

	first, firstTCAID, err := ctx.getTargetCompatibleArgument(args[0])
	if err != nil {
		return nil, nil, err
	}

	second, secondTCAID, err := ctx.getTargetCompatibleArgument(args[1])
	if err != nil {
		return nil, nil, err
	}

	if firstTCAID == pdp.TargetCompatibleArgumentAttributeValue &&
		secondTCAID == pdp.TargetCompatibleArgumentAttributeValue {
		return nil, nil, newMatchFunctionBothValuesError()
	}

	if firstTCAID == pdp.TargetCompatibleArgumentAttributeDesignator &&
		secondTCAID == pdp.TargetCompatibleArgumentAttributeDesignator {
		return nil, nil, newMatchFunctionBothAttrsError()
	}

	return first, second, nil
}

func (ctx context) unmarshalTargetMatchExpression(ID string, v interface{}) (pdp.Expression, boundError) {
	typeFunctionMap, ok := pdp.TargetCompatibleExpressions[strings.ToLower(ID)]
	if !ok {
		return nil, newUnknownMatchFunctionError(ID)
	}

	first, second, err := ctx.getAdjustedArgumentPair(v)
	if err != nil {
		return nil, bindError(err, ID)
	}

	firstType := first.GetResultType()
	secondType := second.GetResultType()

	subTypeFunctionMap, ok := typeFunctionMap[firstType]
	if !ok {
		return nil, newMatchFunctionCastError(ID, firstType, secondType)
	}

	maker, ok := subTypeFunctionMap[secondType]
	if !ok {
		return nil, newMatchFunctionCastError(ID, firstType, secondType)
	}

	return maker(first, second), nil
}

func (ctx context) unmarshalTargetAllOfItem(v interface{}) (pdp.Match, boundError) {
	m, err := ctx.validateMap(v, "expression")
	if err != nil {
		return pdp.Match{}, err
	}

	k, v, err := ctx.getSingleMapPair(m, "expression")
	if err != nil {
		return pdp.Match{}, err
	}

	ID, err := ctx.validateString(k, "function identifier")
	if err != nil {
		return pdp.Match{}, err
	}

	e, err := ctx.unmarshalTargetMatchExpression(ID, v)
	if err != nil {
		return pdp.Match{}, err
	}

	return pdp.MakeMatch(e), nil
}

func (ctx context) unmarshalTargetAnyOfItem(v interface{}) (pdp.AllOf, boundError) {
	al := pdp.MakeAllOf()

	m, err := ctx.validateMap(v, "expression")
	if err != nil {
		return al, err
	}

	k, v, err := ctx.getSingleMapPair(m, "expression")
	if err != nil {
		return al, err
	}

	ID, err := ctx.validateString(k, "function identifier")
	if err != nil {
		return al, err
	}

	if strings.ToLower(ID) == yastTagAll {
		items, err := ctx.validateList(v, "list of all expressions")
		if err != nil {
			return al, bindError(err, ID)
		}

		for i, item := range items {
			m, err := ctx.unmarshalTargetAllOfItem(item)
			if err != nil {
				return al, bindError(bindErrorf(err, "%d", i+1), ID)
			}

			al.Append(m)
		}
	} else {
		e, err := ctx.unmarshalTargetMatchExpression(ID, v)
		if err != nil {
			return al, err
		}

		m := pdp.MakeMatch(e)
		al.Append(m)
	}

	return al, nil
}

func (ctx context) unmarshalTargetItem(v interface{}) (pdp.AnyOf, boundError) {
	an := pdp.MakeAnyOf()

	m, err := ctx.validateMap(v, "expression")
	if err != nil {
		return an, err
	}

	k, v, err := ctx.getSingleMapPair(m, "expression")
	if err != nil {
		return an, err
	}

	ID, err := ctx.validateString(k, "function identifier")
	if err != nil {
		return an, err
	}

	if strings.ToLower(ID) == yastTagAny {
		items, err := ctx.validateList(v, "list of any expressions")
		if err != nil {
			return an, bindError(err, ID)
		}

		for i, item := range items {
			al, err := ctx.unmarshalTargetAnyOfItem(item)
			if err != nil {
				return an, bindError(bindErrorf(err, "%d", i+1), ID)
			}

			an.Append(al)
		}
	} else {
		e, err := ctx.unmarshalTargetMatchExpression(ID, v)
		if err != nil {
			return an, err
		}

		m := pdp.MakeMatch(e)
		al := pdp.MakeAllOf()
		al.Append(m)
		an.Append(al)
	}

	return an, nil
}

func (ctx context) unmarshalTarget(m map[interface{}]interface{}) (pdp.Target, boundError) {
	t := pdp.MakeTarget()
	v, ok := m[yastTagTarget]
	if !ok {
		return t, nil
	}

	items, err := ctx.validateList(v, "target")
	if err != nil {
		return t, err
	}

	for i, item := range items {
		a, err := ctx.unmarshalTargetItem(item)
		if err != nil {
			return t, bindError(bindErrorf(err, "%d", i+1), "target")
		}

		t.Append(a)
	}

	return t, nil
}

package yast

import "github.com/infobloxopen/themis/pdp"

func (ctx *context) unmarshalAttributeDeclaration(k, v interface{}) boundError {
	ID, err := ctx.validateString(k, "attribute id")
	if err != nil {
		return err
	}

	strT, err := ctx.validateString(v, "attribute data type")
	if err != nil {
		return bindError(err, ID)
	}

	t := ctx.symbols.GetType(strT)
	if t == nil {
		return bindError(newUnknownTypeError(strT), ID)
	}

	if err := ctx.symbols.PutAttribute(pdp.MakeAttribute(ID, t)); err != nil {
		return bindError(err, ID)
	}

	return nil
}

func (ctx *context) unmarshalAttributeDeclarations(m map[interface{}]interface{}) boundError {
	attrs, ok, err := ctx.extractMapOpt(m, yastTagAttributes, "attribute declarations")
	if !ok || err != nil {
		return err
	}

	for k, v := range attrs {
		err = ctx.unmarshalAttributeDeclaration(k, v)
		if err != nil {
			return bindError(err, yastTagAttributes)
		}
	}

	return nil
}

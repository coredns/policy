package yast

import (
	"strings"

	"github.com/infobloxopen/themis/pdp"
)

func (ctx *context) unmarshalFlagsTypeDeclaration(id string, m map[interface{}]interface{}) boundError {
	items, err := ctx.extractList(m, yastTagFlags, "list of flag names")
	if err != nil {
		return err
	}

	flags := make([]string, len(items))
	for i, v := range items {
		s, err := ctx.validateString(v, "flag")
		if err != nil {
			return err
		}

		flags[i] = s
	}

	t, eErr := pdp.NewFlagsType(id, flags...)
	if eErr != nil {
		return newExternalError(eErr)
	}

	if eErr := ctx.symbols.PutType(t); eErr != nil {
		return newExternalError(eErr)
	}

	return nil
}

func (ctx *context) unmarshalTypeDeclarationByMetaType(id, meta string, m map[interface{}]interface{}) boundError {
	switch strings.ToLower(meta) {
	case yastTagFlags:
		return ctx.unmarshalFlagsTypeDeclaration(id, m)
	}

	return newUnknownMetaTypeError(meta)
}

func (ctx *context) unmarshalTypeDeclaration(k, v interface{}) boundError {
	ID, err := ctx.validateString(k, "type id")
	if err != nil {
		return err
	}

	m, err := ctx.validateMap(v, "type declaration")
	if err != nil {
		return bindError(err, ID)
	}

	meta, err := ctx.extractString(m, yastTagMeta, "meta type name")
	if err != nil {
		return bindError(err, ID)
	}

	if err := ctx.unmarshalTypeDeclarationByMetaType(ID, meta, m); err != nil {
		return bindError(err, ID)
	}

	return nil
}

func (ctx *context) unmarshalTypeDeclarations(m map[interface{}]interface{}) boundError {
	types, ok, err := ctx.extractMapOpt(m, yastTagTypes, "type declarations")
	if !ok || err != nil {
		return err
	}

	for k, v := range types {
		err = ctx.unmarshalTypeDeclaration(k, v)
		if err != nil {
			return bindError(err, yastTagTypes)
		}
	}

	return nil
}

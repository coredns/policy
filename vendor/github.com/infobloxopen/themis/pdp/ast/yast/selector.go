package yast

import (
	"net/url"

	"github.com/infobloxopen/themis/pdp"
)

func (ctx context) unmarshalSelector(v interface{}) (pdp.Expression, boundError) {
	m, err := ctx.validateMap(v, "selector attributes")
	if err != nil {
		return nil, err
	}

	uri, err := ctx.extractString(m, yastTagURI, "selector URI")
	if err != nil {
		return nil, err
	}

	id, ierr := url.Parse(uri)
	if ierr != nil {
		return nil, newSelectorURIError(uri, ierr)
	}

	items, err := ctx.extractList(m, yastTagPath, "path")
	if err != nil {
		return nil, bindErrorf(err, "selector(%s)", uri)
	}

	path := make([]pdp.Expression, len(items))
	for i, item := range items {
		e, err := ctx.unmarshalExpression(item)
		if err != nil {
			return nil, bindErrorf(bindErrorf(err, "%d", i), "selector(%s)", uri)
		}

		path[i] = e
	}

	st, err := ctx.extractString(m, yastTagType, "type")
	if err != nil {
		return nil, bindErrorf(err, "selector(%s)", uri)
	}

	t := ctx.symbols.GetType(st)
	if t == nil {
		return nil, bindErrorf(newUnknownTypeError(st), "selector(%s)", uri)
	}

	if t == pdp.TypeUndefined {
		return nil, bindErrorf(newInvalidTypeError(t), "selector(%s)", uri)
	}

	e, eErr := pdp.MakeSelector(id, path, t)
	if eErr != nil {
		return nil, bindErrorf(eErr, "selector(%s)", uri)
	}
	return e, nil
}

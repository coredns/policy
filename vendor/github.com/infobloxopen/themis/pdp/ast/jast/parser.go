// Package jast implements policies JSON AST (JAST) parser.
package jast

import (
	"encoding/json"
	"io"

	"github.com/infobloxopen/themis/pdp"

	"github.com/google/uuid"
)

const (
	yastTagTypes      = "types"
	yastTagMeta       = "meta"
	yastTagFlags      = "flags"
	yastTagAttributes = "attributes"
	yastTagID         = "id"
	yastTagTarget     = "target"
	yastTagPolicies   = "policies"
	yastTagRules      = "rules"
	yastTagCondition  = "condition"
	yastTagAlg        = "alg"
	yastTagMap        = "map"
	yastTagDefault    = "default"
	yastTagError      = "error"
	yastTagOrder      = "order"
	yastTagEffect     = "effect"
	yastTagObligation = "obligations"
	yastTagAny        = "any"
	yastTagAll        = "all"
	yastTagAttribute  = "attr"
	yastTagValue      = "val"
	yastTagSelector   = "selector"
	yastTagType       = "type"
	yastTagContent    = "content"
	yastTagURI        = "uri"
	yastTagPath       = "path"
	yastTagOp         = "op"
	yastTagEntity     = "entity"
)

// Parser is a JAST parser implementation.
type Parser struct{}

// Unmarshal parses policies JSON representation to PDP's internal representation.
func (p Parser) Unmarshal(in io.Reader, tag *uuid.UUID) (*pdp.PolicyStorage, error) {
	ctx := newContext()
	if err := ctx.unmarshal(json.NewDecoder(in)); err != nil {
		return nil, err
	}

	return pdp.NewPolicyStorage(ctx.rootPolicy, ctx.symbols, tag), nil
}

// UnmarshalUpdate parses policies update JSON representation to PDP's internal representation.
func (p Parser) UnmarshalUpdate(in io.Reader, s pdp.Symbols, oldTag, newTag uuid.UUID) (*pdp.PolicyUpdate, error) {
	ctx := newContextWithSymbols(s)
	u := pdp.NewPolicyUpdate(oldTag, newTag)
	if err := ctx.unmarshalCommands(json.NewDecoder(in), u); err != nil {
		return nil, err
	}

	return u, nil
}

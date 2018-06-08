// Package ast implements policies Abstract Syntax Tree (AST) parsers.
package ast

import (
	"io"

	"github.com/infobloxopen/themis/pdp"
	"github.com/infobloxopen/themis/pdp/ast/jast"
	"github.com/infobloxopen/themis/pdp/ast/yast"

	"github.com/google/uuid"
)

// Parser is interface that wraps the AST parser methods.
type Parser interface {
	// Unmarshal parses policies representation to PDP's internal
	// representation and returns pointer to PolicyStorage with the policies.
	// It sets given tag to the policies. Policies with no tag can't be updated
	// incrementally.
	Unmarshal(in io.Reader, tag *uuid.UUID) (*pdp.PolicyStorage, error)

	// UnmarshalUpdate parses policies update representation to PDP's internal
	// representation. Requires attribute symbols table as attrs argument which maps
	// attribute name to its specification. Argument oldTag should match current
	// policies tag to make update applicable. Value of newTag is set to policies
	// when update is applied.
	UnmarshalUpdate(in io.Reader, s pdp.Symbols, oldTag, newTag uuid.UUID) (*pdp.PolicyUpdate, error)
}

// NewJSONParser is a JSON AST parser constructor.
func NewJSONParser() Parser {
	return jast.Parser{}
}

// NewYAMLParser is a YAML AST parser constructor.
func NewYAMLParser() Parser {
	return yast.Parser{}
}

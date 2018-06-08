package jast

/* AUTOMATICALLY GENERATED FROM errors.yaml - DO NOT EDIT */

import (
	"encoding/json"
	"fmt"
	"github.com/infobloxopen/themis/pdp"
	"strconv"
	"strings"
)

const (
	externalErrorID                     = 0
	policyAmbiguityErrorID              = 1
	policyMissingKeyErrorID             = 2
	unknownRCAErrorID                   = 3
	missingRCAErrorID                   = 4
	parseCAErrorID                      = 5
	invalidRCAErrorID                   = 6
	missingDefaultRuleRCAErrorID        = 7
	missingErrorRuleRCAErrorID          = 8
	notImplementedRCAErrorID            = 9
	unknownPCAErrorID                   = 10
	missingPCAErrorID                   = 11
	invalidPCAErrorID                   = 12
	missingDefaultPolicyPCAErrorID      = 13
	missingErrorPolicyPCAErrorID        = 14
	notImplementedPCAErrorID            = 15
	mapperArgumentTypeErrorID           = 16
	conditionTypeErrorID                = 17
	unknownEffectErrorID                = 18
	unknownMatchFunctionErrorID         = 19
	matchFunctionCastErrorID            = 20
	matchFunctionArgsNumberErrorID      = 21
	invalidMatchFunctionArgErrorID      = 22
	matchFunctionBothValuesErrorID      = 23
	matchFunctionBothAttrsErrorID       = 24
	unknownFunctionErrorID              = 25
	functionCastErrorID                 = 26
	unknownAttributeErrorID             = 27
	missingAttributeErrorID             = 28
	unknownMapperCAOrderID              = 29
	unknownTypeErrorID                  = 30
	invalidTypeErrorID                  = 31
	missingContentErrorID               = 32
	notImplementedValueTypeErrorID      = 33
	invalidAddressErrorID               = 34
	integerOverflowErrorID              = 35
	invalidNetworkErrorID               = 36
	invalidDomainErrorID                = 37
	selectorURIErrorID                  = 38
	entityAmbiguityErrorID              = 39
	entityMissingKeyErrorID             = 40
	unknownPolicyUpdateOperationErrorID = 41
	missingContentTypeErrorID           = 42
	unknownFieldErrorID                 = 43
	missingMetaTypeNameErrorID          = 44
	unknownMetaTypeErrorID              = 45
	missingFlagNameListErrorID          = 46
	unknownFlagNameErrorID              = 47
)

type externalError struct {
	errorLink
	err error
}

func newExternalError(err error) *externalError {
	return &externalError{
		errorLink: errorLink{id: externalErrorID},
		err:       err}
}

func (e *externalError) Error() string {
	return e.errorf("%s", e.err)
}

type policyAmbiguityError struct {
	errorLink
}

func newPolicyAmbiguityError() *policyAmbiguityError {
	return &policyAmbiguityError{
		errorLink: errorLink{id: policyAmbiguityErrorID}}
}

func (e *policyAmbiguityError) Error() string {
	return e.errorf("Expected rules (for policy) or policies (for policy set) but got both")
}

type policyMissingKeyError struct {
	errorLink
}

func newPolicyMissingKeyError() *policyMissingKeyError {
	return &policyMissingKeyError{
		errorLink: errorLink{id: policyMissingKeyErrorID}}
}

func (e *policyMissingKeyError) Error() string {
	return e.errorf("Expected rules (for policy) or policies (for policy set) but got nothing")
}

type unknownRCAError struct {
	errorLink
	alg string
}

func newUnknownRCAError(alg string) *unknownRCAError {
	return &unknownRCAError{
		errorLink: errorLink{id: unknownRCAErrorID},
		alg:       alg}
}

func (e *unknownRCAError) Error() string {
	return e.errorf("Unknown rule combinig algorithm \"%s\"", e.alg)
}

type missingRCAError struct {
	errorLink
}

func newMissingRCAError() *missingRCAError {
	return &missingRCAError{
		errorLink: errorLink{id: missingRCAErrorID}}
}

func (e *missingRCAError) Error() string {
	return e.errorf("Missing policy combinig algorithm")
}

type parseCAError struct {
	errorLink
	token json.Token
}

func newParseCAError(token json.Token) *parseCAError {
	return &parseCAError{
		errorLink: errorLink{id: parseCAErrorID},
		token:     token}
}

func (e *parseCAError) Error() string {
	return e.errorf("Expected string or { object delimiter for combinig algorithm but got %T (%#v)", e.token, e.token)
}

type invalidRCAError struct {
	errorLink
	v interface{}
}

func newInvalidRCAError(v interface{}) *invalidRCAError {
	return &invalidRCAError{
		errorLink: errorLink{id: invalidRCAErrorID},
		v:         v}
}

func (e *invalidRCAError) Error() string {
	return e.errorf("Expected string or *caParams as policy combinig algorithm but got %T", e.v)
}

type missingDefaultRuleRCAError struct {
	errorLink
	ID string
}

func newMissingDefaultRuleRCAError(ID string) *missingDefaultRuleRCAError {
	return &missingDefaultRuleRCAError{
		errorLink: errorLink{id: missingDefaultRuleRCAErrorID},
		ID:        ID}
}

func (e *missingDefaultRuleRCAError) Error() string {
	return e.errorf("No rule with ID %q to use as default rule", e.ID)
}

type missingErrorRuleRCAError struct {
	errorLink
	ID string
}

func newMissingErrorRuleRCAError(ID string) *missingErrorRuleRCAError {
	return &missingErrorRuleRCAError{
		errorLink: errorLink{id: missingErrorRuleRCAErrorID},
		ID:        ID}
}

func (e *missingErrorRuleRCAError) Error() string {
	return e.errorf("No rule with ID %q to use as on error rule", e.ID)
}

type notImplementedRCAError struct {
	errorLink
	ID string
}

func newNotImplementedRCAError(ID string) *notImplementedRCAError {
	return &notImplementedRCAError{
		errorLink: errorLink{id: notImplementedRCAErrorID},
		ID:        ID}
}

func (e *notImplementedRCAError) Error() string {
	return e.errorf("Parsing for %q rule combinig algorithm hasn't been implemented yet", e.ID)
}

type unknownPCAError struct {
	errorLink
	alg string
}

func newUnknownPCAError(alg string) *unknownPCAError {
	return &unknownPCAError{
		errorLink: errorLink{id: unknownPCAErrorID},
		alg:       alg}
}

func (e *unknownPCAError) Error() string {
	return e.errorf("Unknown policy combinig algorithm \"%s\"", e.alg)
}

type missingPCAError struct {
	errorLink
}

func newMissingPCAError() *missingPCAError {
	return &missingPCAError{
		errorLink: errorLink{id: missingPCAErrorID}}
}

func (e *missingPCAError) Error() string {
	return e.errorf("Missing policy combinig algorithm")
}

type invalidPCAError struct {
	errorLink
	v interface{}
}

func newInvalidPCAError(v interface{}) *invalidPCAError {
	return &invalidPCAError{
		errorLink: errorLink{id: invalidPCAErrorID},
		v:         v}
}

func (e *invalidPCAError) Error() string {
	return e.errorf("Expected string or *caParams as policy combinig algorithm but got %T", e.v)
}

type missingDefaultPolicyPCAError struct {
	errorLink
	ID string
}

func newMissingDefaultPolicyPCAError(ID string) *missingDefaultPolicyPCAError {
	return &missingDefaultPolicyPCAError{
		errorLink: errorLink{id: missingDefaultPolicyPCAErrorID},
		ID:        ID}
}

func (e *missingDefaultPolicyPCAError) Error() string {
	return e.errorf("No policy with ID %q to use as default policy", e.ID)
}

type missingErrorPolicyPCAError struct {
	errorLink
	ID string
}

func newMissingErrorPolicyPCAError(ID string) *missingErrorPolicyPCAError {
	return &missingErrorPolicyPCAError{
		errorLink: errorLink{id: missingErrorPolicyPCAErrorID},
		ID:        ID}
}

func (e *missingErrorPolicyPCAError) Error() string {
	return e.errorf("No policy with ID %q to use as on error policy", e.ID)
}

type notImplementedPCAError struct {
	errorLink
	ID string
}

func newNotImplementedPCAError(ID string) *notImplementedPCAError {
	return &notImplementedPCAError{
		errorLink: errorLink{id: notImplementedPCAErrorID},
		ID:        ID}
}

func (e *notImplementedPCAError) Error() string {
	return e.errorf("Parsing for %q policy combinig algorithm hasn't been implemented yet", e.ID)
}

type mapperArgumentTypeError struct {
	errorLink
	actual pdp.Type
}

func newMapperArgumentTypeError(actual pdp.Type) *mapperArgumentTypeError {
	return &mapperArgumentTypeError{
		errorLink: errorLink{id: mapperArgumentTypeErrorID},
		actual:    actual}
}

func (e *mapperArgumentTypeError) Error() string {
	return e.errorf("Expected %q, %q, %q or flags as argument but got %q", pdp.TypeString, pdp.TypeSetOfStrings, pdp.TypeListOfStrings, e.actual)
}

type conditionTypeError struct {
	errorLink
	t pdp.Type
}

func newConditionTypeError(t pdp.Type) *conditionTypeError {
	return &conditionTypeError{
		errorLink: errorLink{id: conditionTypeErrorID},
		t:         t}
}

func (e *conditionTypeError) Error() string {
	return e.errorf("Expected %q as condition expression result but got %q", pdp.TypeBoolean, e.t)
}

type unknownEffectError struct {
	errorLink
	e string
}

func newUnknownEffectError(e string) *unknownEffectError {
	return &unknownEffectError{
		errorLink: errorLink{id: unknownEffectErrorID},
		e:         e}
}

func (e *unknownEffectError) Error() string {
	return e.errorf("Unknown rule effect %q", e.e)
}

type unknownMatchFunctionError struct {
	errorLink
	ID string
}

func newUnknownMatchFunctionError(ID string) *unknownMatchFunctionError {
	return &unknownMatchFunctionError{
		errorLink: errorLink{id: unknownMatchFunctionErrorID},
		ID:        ID}
}

func (e *unknownMatchFunctionError) Error() string {
	return e.errorf("Unknown match function %q", e.ID)
}

type matchFunctionCastError struct {
	errorLink
	ID     string
	first  pdp.Type
	second pdp.Type
}

func newMatchFunctionCastError(ID string, first, second pdp.Type) *matchFunctionCastError {
	return &matchFunctionCastError{
		errorLink: errorLink{id: matchFunctionCastErrorID},
		ID:        ID,
		first:     first,
		second:    second}
}

func (e *matchFunctionCastError) Error() string {
	return e.errorf("No function %q for arguments %q and %q", e.ID, e.first, e.second)
}

type matchFunctionArgsNumberError struct {
	errorLink
	n int
}

func newMatchFunctionArgsNumberError(n int) *matchFunctionArgsNumberError {
	return &matchFunctionArgsNumberError{
		errorLink: errorLink{id: matchFunctionArgsNumberErrorID},
		n:         n}
}

func (e *matchFunctionArgsNumberError) Error() string {
	return e.errorf("Expected two arguments got %d", e.n)
}

type invalidMatchFunctionArgError struct {
	errorLink
	expr pdp.Expression
}

func newInvalidMatchFunctionArgError(expr pdp.Expression) *invalidMatchFunctionArgError {
	return &invalidMatchFunctionArgError{
		errorLink: errorLink{id: invalidMatchFunctionArgErrorID},
		expr:      expr}
}

func (e *invalidMatchFunctionArgError) Error() string {
	return e.errorf("Expected one immediate value and one attribute got %T", e.expr)
}

type matchFunctionBothValuesError struct {
	errorLink
}

func newMatchFunctionBothValuesError() *matchFunctionBothValuesError {
	return &matchFunctionBothValuesError{
		errorLink: errorLink{id: matchFunctionBothValuesErrorID}}
}

func (e *matchFunctionBothValuesError) Error() string {
	return e.errorf("Expected one immediate value and one attribute got both immediate values")
}

type matchFunctionBothAttrsError struct {
	errorLink
}

func newMatchFunctionBothAttrsError() *matchFunctionBothAttrsError {
	return &matchFunctionBothAttrsError{
		errorLink: errorLink{id: matchFunctionBothAttrsErrorID}}
}

func (e *matchFunctionBothAttrsError) Error() string {
	return e.errorf("Expected one immediate value and one attribute got both immediate values")
}

type unknownFunctionError struct {
	errorLink
	ID string
}

func newUnknownFunctionError(ID string) *unknownFunctionError {
	return &unknownFunctionError{
		errorLink: errorLink{id: unknownFunctionErrorID},
		ID:        ID}
}

func (e *unknownFunctionError) Error() string {
	return e.errorf("Unknown function %q", e.ID)
}

type functionCastError struct {
	errorLink
	ID    string
	exprs []pdp.Expression
}

func newFunctionCastError(ID string, exprs []pdp.Expression) *functionCastError {
	return &functionCastError{
		errorLink: errorLink{id: functionCastErrorID},
		ID:        ID,
		exprs:     exprs}
}

func (e *functionCastError) Error() string {
	args := ""
	if len(e.exprs) > 1 {
		t := make([]string, len(e.exprs))
		for i, e := range e.exprs {
			t[i] = strconv.Quote(e.GetResultType().String())
		}
		args = fmt.Sprintf("%d arguments of following types %q", len(e.exprs), strings.Join(t, ", "))
	} else if len(e.exprs) > 0 {
		args = fmt.Sprintf("argument of type %q", e.exprs[0].GetResultType())
	} else {
		args = "no arguments"
	}

	return e.errorf("Can't find function %s which takes %s", e.ID, args)
}

type unknownAttributeError struct {
	errorLink
	ID string
}

func newUnknownAttributeError(ID string) *unknownAttributeError {
	return &unknownAttributeError{
		errorLink: errorLink{id: unknownAttributeErrorID},
		ID:        ID}
}

func (e *unknownAttributeError) Error() string {
	return e.errorf("Unknown attribute %q", e.ID)
}

type missingAttributeError struct {
	errorLink
	attr string
	obj  string
}

func newMissingAttributeError(attr, obj string) *missingAttributeError {
	return &missingAttributeError{
		errorLink: errorLink{id: missingAttributeErrorID},
		attr:      attr,
		obj:       obj}
}

func (e *missingAttributeError) Error() string {
	return e.errorf("Missing %q attribute %q", e.obj, e.attr)
}

type unknownMapperCAOrder struct {
	errorLink
	ord string
}

func newUnknownMapperCAOrder(ord string) *unknownMapperCAOrder {
	return &unknownMapperCAOrder{
		errorLink: errorLink{id: unknownMapperCAOrderID},
		ord:       ord}
}

func (e *unknownMapperCAOrder) Error() string {
	return e.errorf("Unknown ordering for mapper \"%s\"", e.ord)
}

type unknownTypeError struct {
	errorLink
	t string
}

func newUnknownTypeError(t string) *unknownTypeError {
	return &unknownTypeError{
		errorLink: errorLink{id: unknownTypeErrorID},
		t:         t}
}

func (e *unknownTypeError) Error() string {
	return e.errorf("Unknown value type %q", e.t)
}

type invalidTypeError struct {
	errorLink
	t pdp.Type
}

func newInvalidTypeError(t pdp.Type) *invalidTypeError {
	return &invalidTypeError{
		errorLink: errorLink{id: invalidTypeErrorID},
		t:         t}
}

func (e *invalidTypeError) Error() string {
	return e.errorf("Can't make value of %q type", e.t)
}

type missingContentError struct {
	errorLink
}

func newMissingContentError() *missingContentError {
	return &missingContentError{
		errorLink: errorLink{id: missingContentErrorID}}
}

func (e *missingContentError) Error() string {
	return e.errorf("Missing value content")
}

type notImplementedValueTypeError struct {
	errorLink
	t pdp.Type
}

func newNotImplementedValueTypeError(t pdp.Type) *notImplementedValueTypeError {
	return &notImplementedValueTypeError{
		errorLink: errorLink{id: notImplementedValueTypeErrorID},
		t:         t}
}

func (e *notImplementedValueTypeError) Error() string {
	return e.errorf("Parsing for type %s hasn't been implemented yet", e.t)
}

type invalidAddressError struct {
	errorLink
	s string
}

func newInvalidAddressError(s string) *invalidAddressError {
	return &invalidAddressError{
		errorLink: errorLink{id: invalidAddressErrorID},
		s:         s}
}

func (e *invalidAddressError) Error() string {
	return e.errorf("Expected value of address type but got %q", e.s)
}

type integerOverflowError struct {
	errorLink
	x float64
}

func newIntegerOverflowError(x float64) *integerOverflowError {
	return &integerOverflowError{
		errorLink: errorLink{id: integerOverflowErrorID},
		x:         x}
}

func (e *integerOverflowError) Error() string {
	return e.errorf("%f overflows integer", e.x)
}

type invalidNetworkError struct {
	errorLink
	s   string
	err error
}

func newInvalidNetworkError(s string, err error) *invalidNetworkError {
	return &invalidNetworkError{
		errorLink: errorLink{id: invalidNetworkErrorID},
		s:         s,
		err:       err}
}

func (e *invalidNetworkError) Error() string {
	return e.errorf("Expected value of network type but got %q (%v)", e.s, e.err)
}

type invalidDomainError struct {
	errorLink
	s   string
	err error
}

func newInvalidDomainError(s string, err error) *invalidDomainError {
	return &invalidDomainError{
		errorLink: errorLink{id: invalidDomainErrorID},
		s:         s,
		err:       err}
}

func (e *invalidDomainError) Error() string {
	return e.errorf("Expected value of domain type but got %q (%v)", e.s, e.err)
}

type selectorURIError struct {
	errorLink
	uri string
	err error
}

func newSelectorURIError(uri string, err error) *selectorURIError {
	return &selectorURIError{
		errorLink: errorLink{id: selectorURIErrorID},
		uri:       uri,
		err:       err}
}

func (e *selectorURIError) Error() string {
	return e.errorf("Expected seletor URI but got %q (%s)", e.uri, e.err)
}

type entityAmbiguityError struct {
	errorLink
	fields []string
}

func newEntityAmbiguityError(fields []string) *entityAmbiguityError {
	return &entityAmbiguityError{
		errorLink: errorLink{id: entityAmbiguityErrorID},
		fields:    fields}
}

func (e *entityAmbiguityError) Error() string {
	return e.errorf("Expected rules (for policy), policies (for policy set) or effect (for rule) but got %s", strings.Join(e.fields, ", "))
}

type entityMissingKeyError struct {
	errorLink
}

func newEntityMissingKeyError() *entityMissingKeyError {
	return &entityMissingKeyError{
		errorLink: errorLink{id: entityMissingKeyErrorID}}
}

func (e *entityMissingKeyError) Error() string {
	return e.errorf("Expected rules (for policy), policies (for policy set) or effect (for rule) but got nothing")
}

type unknownPolicyUpdateOperationError struct {
	errorLink
	op string
}

func newUnknownPolicyUpdateOperationError(op string) *unknownPolicyUpdateOperationError {
	return &unknownPolicyUpdateOperationError{
		errorLink: errorLink{id: unknownPolicyUpdateOperationErrorID},
		op:        op}
}

func (e *unknownPolicyUpdateOperationError) Error() string {
	return e.errorf("Unknown policy update operation %q", e.op)
}

type missingContentTypeError struct {
	errorLink
}

func newMissingContentTypeError() *missingContentTypeError {
	return &missingContentTypeError{
		errorLink: errorLink{id: missingContentTypeErrorID}}
}

func (e *missingContentTypeError) Error() string {
	return e.errorf("Value 'type' attribute is missing or placed after 'content' attribute")
}

type unknownFieldError struct {
	errorLink
	name string
}

func newUnknownFieldError(name string) *unknownFieldError {
	return &unknownFieldError{
		errorLink: errorLink{id: unknownFieldErrorID},
		name:      name}
}

func (e *unknownFieldError) Error() string {
	return e.errorf("Unknown field %q", e.name)
}

type missingMetaTypeNameError struct {
	errorLink
}

func newMissingMetaTypeNameError() *missingMetaTypeNameError {
	return &missingMetaTypeNameError{
		errorLink: errorLink{id: missingMetaTypeNameErrorID}}
}

func (e *missingMetaTypeNameError) Error() string {
	return e.errorf("Missing meta type name")
}

type unknownMetaTypeError struct {
	errorLink
	meta string
}

func newUnknownMetaTypeError(meta string) *unknownMetaTypeError {
	return &unknownMetaTypeError{
		errorLink: errorLink{id: unknownMetaTypeErrorID},
		meta:      meta}
}

func (e *unknownMetaTypeError) Error() string {
	return e.errorf("Unknown meta type %q", e.meta)
}

type missingFlagNameListError struct {
	errorLink
}

func newMissingFlagNameListError() *missingFlagNameListError {
	return &missingFlagNameListError{
		errorLink: errorLink{id: missingFlagNameListErrorID}}
}

func (e *missingFlagNameListError) Error() string {
	return e.errorf("Missing list of flag names")
}

type unknownFlagNameError struct {
	errorLink
	name string
	t    *pdp.FlagsType
}

func newUnknownFlagNameError(name string, t *pdp.FlagsType) *unknownFlagNameError {
	return &unknownFlagNameError{
		errorLink: errorLink{id: unknownFlagNameErrorID},
		name:      name,
		t:         t}
}

func (e *unknownFlagNameError) Error() string {
	return e.errorf("Type %q doesn't have flag %q", e.t, e.name)
}

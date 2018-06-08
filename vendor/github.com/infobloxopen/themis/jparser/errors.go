package jparser

/* AUTOMATICALLY GENERATED FROM errors.yaml - DO NOT EDIT */

import "encoding/json"

const (
	externalErrorID                       = 0
	rootObjectStartTokenErrorID           = 1
	rootObjectStartDelimiterErrorID       = 2
	objectStartTokenErrorID               = 3
	objectStartDelimiterErrorID           = 4
	objectEndDelimiterErrorID             = 5
	objectTokenErrorID                    = 6
	rootArrayStartTokenErrorID            = 7
	rootArrayStartDelimiterErrorID        = 8
	arrayStartTokenErrorID                = 9
	arrayStartDelimiterErrorID            = 10
	arrayEndDelimiterErrorID              = 11
	stringArrayTokenErrorID               = 12
	objectArrayStartTokenErrorID          = 13
	objectArrayStartDelimiterErrorID      = 14
	objectArrayTokenErrorID               = 15
	unexpectedObjectArrayDelimiterErrorID = 16
	unexpectedDelimiterErrorID            = 17
	missingEOFErrorID                     = 18
	booleanCastErrorID                    = 19
	stringCastErrorID                     = 20
	numberCastErrorID                     = 21
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

type rootObjectStartTokenError struct {
	errorLink
	actual   json.Token
	expected string
}

func newRootObjectStartTokenError(actual json.Token, expected string) *rootObjectStartTokenError {
	return &rootObjectStartTokenError{
		errorLink: errorLink{id: rootObjectStartTokenErrorID},
		actual:    actual,
		expected:  expected}
}

func (e *rootObjectStartTokenError) Error() string {
	return e.errorf("Expected root JSON object start %q but got token %T (%#v)", e.expected, e.actual, e.actual)
}

type rootObjectStartDelimiterError struct {
	errorLink
	actual   json.Delim
	expected string
}

func newRootObjectStartDelimiterError(actual json.Delim, expected string) *rootObjectStartDelimiterError {
	return &rootObjectStartDelimiterError{
		errorLink: errorLink{id: rootObjectStartDelimiterErrorID},
		actual:    actual,
		expected:  expected}
}

func (e *rootObjectStartDelimiterError) Error() string {
	return e.errorf("Expected root JSON object start %q but got delimiter %q", e.expected, e.actual)
}

type objectStartTokenError struct {
	errorLink
	actual   json.Token
	expected string
	desc     string
}

func newObjectStartTokenError(actual json.Token, expected, desc string) *objectStartTokenError {
	return &objectStartTokenError{
		errorLink: errorLink{id: objectStartTokenErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *objectStartTokenError) Error() string {
	return e.errorf("Expected %s JSON object start %q but got token %T (%#v)", e.desc, e.expected, e.actual, e.actual)
}

type objectStartDelimiterError struct {
	errorLink
	actual   json.Delim
	expected string
	desc     string
}

func newObjectStartDelimiterError(actual json.Delim, expected, desc string) *objectStartDelimiterError {
	return &objectStartDelimiterError{
		errorLink: errorLink{id: objectStartDelimiterErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *objectStartDelimiterError) Error() string {
	return e.errorf("Expected %s JSON object start %q but got delimiter %q", e.desc, e.expected, e.actual)
}

type objectEndDelimiterError struct {
	errorLink
	actual   json.Delim
	expected string
	desc     string
}

func newObjectEndDelimiterError(actual json.Delim, expected, desc string) *objectEndDelimiterError {
	return &objectEndDelimiterError{
		errorLink: errorLink{id: objectEndDelimiterErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *objectEndDelimiterError) Error() string {
	return e.errorf("Expected %s JSON object end %q but got delimiter %q", e.desc, e.expected, e.actual)
}

type objectTokenError struct {
	errorLink
	actual   json.Token
	expected string
	desc     string
}

func newObjectTokenError(actual json.Token, expected, desc string) *objectTokenError {
	return &objectTokenError{
		errorLink: errorLink{id: objectTokenErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *objectTokenError) Error() string {
	return e.errorf("Expected %s JSON object string key or end %q but got token %T (%#v)", e.desc, e.expected, e.actual, e.actual)
}

type rootArrayStartTokenError struct {
	errorLink
	actual   json.Token
	expected string
}

func newRootArrayStartTokenError(actual json.Token, expected string) *rootArrayStartTokenError {
	return &rootArrayStartTokenError{
		errorLink: errorLink{id: rootArrayStartTokenErrorID},
		actual:    actual,
		expected:  expected}
}

func (e *rootArrayStartTokenError) Error() string {
	return e.errorf("Expected root JSON array start %q but got token %T (%#v)", e.expected, e.actual, e.actual)
}

type rootArrayStartDelimiterError struct {
	errorLink
	actual   json.Delim
	expected string
}

func newRootArrayStartDelimiterError(actual json.Delim, expected string) *rootArrayStartDelimiterError {
	return &rootArrayStartDelimiterError{
		errorLink: errorLink{id: rootArrayStartDelimiterErrorID},
		actual:    actual,
		expected:  expected}
}

func (e *rootArrayStartDelimiterError) Error() string {
	return e.errorf("Expected root JSON array start %q but got delimiter %q", e.expected, e.actual)
}

type arrayStartTokenError struct {
	errorLink
	actual   json.Token
	expected string
	desc     string
}

func newArrayStartTokenError(actual json.Token, expected, desc string) *arrayStartTokenError {
	return &arrayStartTokenError{
		errorLink: errorLink{id: arrayStartTokenErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *arrayStartTokenError) Error() string {
	return e.errorf("Expected %s JSON array start %q but got token %T (%#v)", e.desc, e.expected, e.actual, e.actual)
}

type arrayStartDelimiterError struct {
	errorLink
	actual   json.Delim
	expected string
	desc     string
}

func newArrayStartDelimiterError(actual json.Delim, expected, desc string) *arrayStartDelimiterError {
	return &arrayStartDelimiterError{
		errorLink: errorLink{id: arrayStartDelimiterErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *arrayStartDelimiterError) Error() string {
	return e.errorf("Expected %s JSON array start %q but got delimiter %q", e.desc, e.expected, e.actual)
}

type arrayEndDelimiterError struct {
	errorLink
	actual   json.Delim
	expected string
	desc     string
}

func newArrayEndDelimiterError(actual json.Delim, expected, desc string) *arrayEndDelimiterError {
	return &arrayEndDelimiterError{
		errorLink: errorLink{id: arrayEndDelimiterErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *arrayEndDelimiterError) Error() string {
	return e.errorf("Expected %s JSON array end %q but got delimiter %q", e.desc, e.expected, e.actual)
}

type stringArrayTokenError struct {
	errorLink
	actual   json.Token
	expected string
	desc     string
}

func newStringArrayTokenError(actual json.Token, expected, desc string) *stringArrayTokenError {
	return &stringArrayTokenError{
		errorLink: errorLink{id: stringArrayTokenErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *stringArrayTokenError) Error() string {
	return e.errorf("Expected %s JSON array string value or end %q but got token %T (%#v)", e.desc, e.expected, e.actual, e.actual)
}

type objectArrayStartTokenError struct {
	errorLink
	actual         json.Token
	firstExpected  string
	secondExpected string
	desc           string
}

func newObjectArrayStartTokenError(actual json.Token, firstExpected, secondExpected, desc string) *objectArrayStartTokenError {
	return &objectArrayStartTokenError{
		errorLink:      errorLink{id: objectArrayStartTokenErrorID},
		actual:         actual,
		firstExpected:  firstExpected,
		secondExpected: secondExpected,
		desc:           desc}
}

func (e *objectArrayStartTokenError) Error() string {
	return e.errorf("Expected %s JSON object or array start %q or %q but got token %T (%#v)", e.desc, e.firstExpected, e.secondExpected, e.actual, e.actual)
}

type objectArrayStartDelimiterError struct {
	errorLink
	actual         json.Delim
	firstExpected  string
	secondExpected string
	desc           string
}

func newObjectArrayStartDelimiterError(actual json.Delim, firstExpected, secondExpected, desc string) *objectArrayStartDelimiterError {
	return &objectArrayStartDelimiterError{
		errorLink:      errorLink{id: objectArrayStartDelimiterErrorID},
		actual:         actual,
		firstExpected:  firstExpected,
		secondExpected: secondExpected,
		desc:           desc}
}

func (e *objectArrayStartDelimiterError) Error() string {
	return e.errorf("Expected %s JSON object or array start %q or %q but got delimiter %q", e.desc, e.firstExpected, e.secondExpected, e.actual)
}

type objectArrayTokenError struct {
	errorLink
	actual   json.Token
	expected string
	desc     string
}

func newObjectArrayTokenError(actual json.Token, expected, desc string) *objectArrayTokenError {
	return &objectArrayTokenError{
		errorLink: errorLink{id: objectArrayTokenErrorID},
		actual:    actual,
		expected:  expected,
		desc:      desc}
}

func (e *objectArrayTokenError) Error() string {
	return e.errorf("Expected %s JSON array object or end %q but got token %T (%#v)", e.desc, e.expected, e.actual, e.actual)
}

type unexpectedObjectArrayDelimiterError struct {
	errorLink
	delim string
	desc  string
}

func newUnexpectedObjectArrayDelimiterError(delim, desc string) *unexpectedObjectArrayDelimiterError {
	return &unexpectedObjectArrayDelimiterError{
		errorLink: errorLink{id: unexpectedObjectArrayDelimiterErrorID},
		delim:     delim,
		desc:      desc}
}

func (e *unexpectedObjectArrayDelimiterError) Error() string {
	return e.errorf("Unexpected delimiter %q for %s", e.delim, e.desc)
}

type unexpectedDelimiterError struct {
	errorLink
	delim string
	desc  string
}

func newUnexpectedDelimiterError(delim, desc string) *unexpectedDelimiterError {
	return &unexpectedDelimiterError{
		errorLink: errorLink{id: unexpectedDelimiterErrorID},
		delim:     delim,
		desc:      desc}
}

func (e *unexpectedDelimiterError) Error() string {
	return e.errorf("Unexpected delimiter %q for %s", e.delim, e.desc)
}

type missingEOFError struct {
	errorLink
	token json.Token
}

func newMissingEOFError(token json.Token) *missingEOFError {
	return &missingEOFError{
		errorLink: errorLink{id: missingEOFErrorID},
		token:     token}
}

func (e *missingEOFError) Error() string {
	return e.errorf("Expected expected EOF after root object end but got %T (%#v)", e.token, e.token)
}

type booleanCastError struct {
	errorLink
	token json.Token
	desc  string
}

func newBooleanCastError(token json.Token, desc string) *booleanCastError {
	return &booleanCastError{
		errorLink: errorLink{id: booleanCastErrorID},
		token:     token,
		desc:      desc}
}

func (e *booleanCastError) Error() string {
	return e.errorf("Expected boolean as %s but got %T (%#v)", e.desc, e.token, e.token)
}

type stringCastError struct {
	errorLink
	token json.Token
	desc  string
}

func newStringCastError(token json.Token, desc string) *stringCastError {
	return &stringCastError{
		errorLink: errorLink{id: stringCastErrorID},
		token:     token,
		desc:      desc}
}

func (e *stringCastError) Error() string {
	return e.errorf("Expected string as %s but got %T (%#v)", e.desc, e.token, e.token)
}

type numberCastError struct {
	errorLink
	token json.Token
	desc  string
}

func newNumberCastError(token json.Token, desc string) *numberCastError {
	return &numberCastError{
		errorLink: errorLink{id: numberCastErrorID},
		token:     token,
		desc:      desc}
}

func (e *numberCastError) Error() string {
	return e.errorf("Expected number as %s but got %T (%#v)", e.desc, e.token, e.token)
}

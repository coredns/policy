package server

/* AUTOMATICALLY GENERATED FROM errors.yaml - DO NOT EDIT */

import (
	"fmt"
	"github.com/infobloxopen/themis/pdp-control"
	"strings"
)

const (
	externalErrorID                   = 0
	multiErrorID                      = 1
	tracingTypeErrorID                = 2
	unknownEffectErrorID              = 3
	unknownAttributeTypeErrorID       = 4
	contextCreationErrorID            = 5
	missingPolicyErrorID              = 6
	policyCalculationErrorID          = 7
	effectTranslationErrorID          = 8
	effectCombiningErrorID            = 9
	obligationTranslationErrorID      = 10
	queueOverflowErrorID              = 11
	unknownUploadRequestErrorID       = 12
	invalidFromTagErrorID             = 13
	invalidToTagErrorID               = 14
	invalidTagsErrorID                = 15
	tagCheckErrorID                   = 16
	emptyUploadErrorID                = 17
	unknownUploadErrorID              = 18
	policyUploadParseErrorID          = 19
	policyUploadStoreErrorID          = 20
	contentUploadParseErrorID         = 21
	contentUploadStoreErrorID         = 22
	missingPolicyStorageErrorID       = 23
	policyTransactionCreationErrorID  = 24
	policyUpdateParseErrorID          = 25
	policyUpdateApplicationErrorID    = 26
	policyUpdateUploadStoreErrorID    = 27
	policyTransactionCommitErrorID    = 28
	missingPolicyDataApplyErrorID     = 29
	missingContentDataApplyErrorID    = 30
	contentTransactionCreationErrorID = 31
	contentUpdateParseErrorID         = 32
	contentUpdateApplicationErrorID   = 33
	contentUpdateUploadStoreErrorID   = 34
	contentTransactionCommitErrorID   = 35
	unknownUploadedRequestErrorID     = 36
	unsupportedPolicyFromatErrorID    = 37
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

type multiError struct {
	errorLink
	errs []error
}

func newMultiError(errs []error) *multiError {
	return &multiError{
		errorLink: errorLink{id: multiErrorID},
		errs:      errs}
}

func (e *multiError) Error() string {
	msgs := make([]string, len(e.errs))
	for i, err := range e.errs {
		msgs[i] = fmt.Sprintf("%q", err.Error())
	}
	msg := strings.Join(msgs, ", ")

	return e.errorf("multiple errors: %s", msg)
}

type tracingTypeError struct {
	errorLink
	t string
}

func newTracingTypeError(t string) *tracingTypeError {
	return &tracingTypeError{
		errorLink: errorLink{id: tracingTypeErrorID},
		t:         t}
}

func (e *tracingTypeError) Error() string {
	return e.errorf("Unknown tracing type %q", e.t)
}

type unknownEffectError struct {
	errorLink
	effect int
}

func newUnknownEffectError(effect int) *unknownEffectError {
	return &unknownEffectError{
		errorLink: errorLink{id: unknownEffectErrorID},
		effect:    effect}
}

func (e *unknownEffectError) Error() string {
	return e.errorf("Unknown policy effect %d", e.effect)
}

type unknownAttributeTypeError struct {
	errorLink
	t string
}

func newUnknownAttributeTypeError(t string) *unknownAttributeTypeError {
	return &unknownAttributeTypeError{
		errorLink: errorLink{id: unknownAttributeTypeErrorID},
		t:         t}
}

func (e *unknownAttributeTypeError) Error() string {
	return e.errorf("Unknown attribute type %q", e.t)
}

type contextCreationError struct {
	errorLink
	err error
}

func newContextCreationError(err error) *contextCreationError {
	return &contextCreationError{
		errorLink: errorLink{id: contextCreationErrorID},
		err:       err}
}

func (e *contextCreationError) Error() string {
	return e.errorf("Failed to create request context: %s", e.err)
}

type missingPolicyError struct {
	errorLink
}

func newMissingPolicyError() *missingPolicyError {
	return &missingPolicyError{
		errorLink: errorLink{id: missingPolicyErrorID}}
}

func (e *missingPolicyError) Error() string {
	return e.errorf("There is no any policy to process request")
}

type policyCalculationError struct {
	errorLink
	err error
}

func newPolicyCalculationError(err error) *policyCalculationError {
	return &policyCalculationError{
		errorLink: errorLink{id: policyCalculationErrorID},
		err:       err}
}

func (e *policyCalculationError) Error() string {
	return e.errorf("Failed to process request: %s", e.err)
}

type effectTranslationError struct {
	errorLink
	err error
}

func newEffectTranslationError(err error) *effectTranslationError {
	return &effectTranslationError{
		errorLink: errorLink{id: effectTranslationErrorID},
		err:       err}
}

func (e *effectTranslationError) Error() string {
	return e.errorf("Failed to translate effect: %s", e.err)
}

type effectCombiningError struct {
	errorLink
	err error
}

func newEffectCombiningError(err error) *effectCombiningError {
	return &effectCombiningError{
		errorLink: errorLink{id: effectCombiningErrorID},
		err:       err}
}

func (e *effectCombiningError) Error() string {
	return e.errorf("Failed to make failure effect: %s", e.err)
}

type obligationTranslationError struct {
	errorLink
	err error
}

func newObligationTranslationError(err error) *obligationTranslationError {
	return &obligationTranslationError{
		errorLink: errorLink{id: obligationTranslationErrorID},
		err:       err}
}

func (e *obligationTranslationError) Error() string {
	return e.errorf("Failed to translate obligations: %s", e.err)
}

type queueOverflowError struct {
	errorLink
	idx int32
}

func newQueueOverflowError(idx int32) *queueOverflowError {
	return &queueOverflowError{
		errorLink: errorLink{id: queueOverflowErrorID},
		idx:       idx}
}

func (e *queueOverflowError) Error() string {
	return e.errorf("Can't enqueue more than %d items", e.idx)
}

type unknownUploadRequestError struct {
	errorLink
	t control.Item_DataType
}

func newUnknownUploadRequestError(t control.Item_DataType) *unknownUploadRequestError {
	return &unknownUploadRequestError{
		errorLink: errorLink{id: unknownUploadRequestErrorID},
		t:         t}
}

func (e *unknownUploadRequestError) Error() string {
	return e.errorf("Unknown upload request type: %d", e.t)
}

type invalidFromTagError struct {
	errorLink
	tag string
	err error
}

func newInvalidFromTagError(tag string, err error) *invalidFromTagError {
	return &invalidFromTagError{
		errorLink: errorLink{id: invalidFromTagErrorID},
		tag:       tag,
		err:       err}
}

func (e *invalidFromTagError) Error() string {
	return e.errorf("Can't treat %q as current tag: %s", e.tag, e.err)
}

type invalidToTagError struct {
	errorLink
	tag string
	err error
}

func newInvalidToTagError(tag string, err error) *invalidToTagError {
	return &invalidToTagError{
		errorLink: errorLink{id: invalidToTagErrorID},
		tag:       tag,
		err:       err}
}

func (e *invalidToTagError) Error() string {
	return e.errorf("Can't treat %q as new tag: %s", e.tag, e.err)
}

type invalidTagsError struct {
	errorLink
	tag string
}

func newInvalidTagsError(tag string) *invalidTagsError {
	return &invalidTagsError{
		errorLink: errorLink{id: invalidTagsErrorID},
		tag:       tag}
}

func (e *invalidTagsError) Error() string {
	return e.errorf("Can't update from %q tag to no tag", e.tag)
}

type tagCheckError struct {
	errorLink
	err error
}

func newTagCheckError(err error) *tagCheckError {
	return &tagCheckError{
		errorLink: errorLink{id: tagCheckErrorID},
		err:       err}
}

func (e *tagCheckError) Error() string {
	return e.errorf("Failed tag check: %s", e.err)
}

type emptyUploadError struct {
	errorLink
}

func newEmptyUploadError() *emptyUploadError {
	return &emptyUploadError{
		errorLink: errorLink{id: emptyUploadErrorID}}
}

func (e *emptyUploadError) Error() string {
	return e.errorf("Empty upload")
}

type unknownUploadError struct {
	errorLink
	id int32
}

func newUnknownUploadError(id int32) *unknownUploadError {
	return &unknownUploadError{
		errorLink: errorLink{id: unknownUploadErrorID},
		id:        id}
}

func (e *unknownUploadError) Error() string {
	return e.errorf("Can't find upload request with id %d", e.id)
}

type policyUploadParseError struct {
	errorLink
	id  int32
	err error
}

func newPolicyUploadParseError(id int32, err error) *policyUploadParseError {
	return &policyUploadParseError{
		errorLink: errorLink{id: policyUploadParseErrorID},
		id:        id,
		err:       err}
}

func (e *policyUploadParseError) Error() string {
	return e.errorf("Failed to parse policy %d: %s", e.id, e.err)
}

type policyUploadStoreError struct {
	errorLink
	id  int32
	err error
}

func newPolicyUploadStoreError(id int32, err error) *policyUploadStoreError {
	return &policyUploadStoreError{
		errorLink: errorLink{id: policyUploadStoreErrorID},
		id:        id,
		err:       err}
}

func (e *policyUploadStoreError) Error() string {
	return e.errorf("Failed to store parsed policy %d: %s", e.id, e.err)
}

type contentUploadParseError struct {
	errorLink
	id  int32
	err error
}

func newContentUploadParseError(id int32, err error) *contentUploadParseError {
	return &contentUploadParseError{
		errorLink: errorLink{id: contentUploadParseErrorID},
		id:        id,
		err:       err}
}

func (e *contentUploadParseError) Error() string {
	return e.errorf("Failed to parse content %d: %s", e.id, e.err)
}

type contentUploadStoreError struct {
	errorLink
	id  int32
	err error
}

func newContentUploadStoreError(id int32, err error) *contentUploadStoreError {
	return &contentUploadStoreError{
		errorLink: errorLink{id: contentUploadStoreErrorID},
		id:        id,
		err:       err}
}

func (e *contentUploadStoreError) Error() string {
	return e.errorf("Failed to store parsed content %d: %s", e.id, e.err)
}

type missingPolicyStorageError struct {
	errorLink
}

func newMissingPolicyStorageError() *missingPolicyStorageError {
	return &missingPolicyStorageError{
		errorLink: errorLink{id: missingPolicyStorageErrorID}}
}

func (e *missingPolicyStorageError) Error() string {
	return e.errorf("No any policy to update")
}

type policyTransactionCreationError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newPolicyTransactionCreationError(id int32, v *item, err error) *policyTransactionCreationError {
	return &policyTransactionCreationError{
		errorLink: errorLink{id: policyTransactionCreationErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *policyTransactionCreationError) Error() string {
	return e.errorf("Can't create transaction for policy update %d from tag %q to %q: %s", e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type policyUpdateParseError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newPolicyUpdateParseError(id int32, v *item, err error) *policyUpdateParseError {
	return &policyUpdateParseError{
		errorLink: errorLink{id: policyUpdateParseErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *policyUpdateParseError) Error() string {
	return e.errorf("Failed to parse update %d from tag %q to %q: %s", e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type policyUpdateApplicationError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newPolicyUpdateApplicationError(id int32, v *item, err error) *policyUpdateApplicationError {
	return &policyUpdateApplicationError{
		errorLink: errorLink{id: policyUpdateApplicationErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *policyUpdateApplicationError) Error() string {
	return e.errorf("Failed to apply update %d from tag %q to %q: %s", e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type policyUpdateUploadStoreError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newPolicyUpdateUploadStoreError(id int32, v *item, err error) *policyUpdateUploadStoreError {
	return &policyUpdateUploadStoreError{
		errorLink: errorLink{id: policyUpdateUploadStoreErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *policyUpdateUploadStoreError) Error() string {
	return e.errorf("Failed to store parsed policy update %d from tag %q to %q: %s", e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type policyTransactionCommitError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newPolicyTransactionCommitError(id int32, v *item, err error) *policyTransactionCommitError {
	return &policyTransactionCommitError{
		errorLink: errorLink{id: policyTransactionCommitErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *policyTransactionCommitError) Error() string {
	return e.errorf("Failed to commit transaction %d from tag %q to %q: %s", e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type missingPolicyDataApplyError struct {
	errorLink
	id int32
}

func newMissingPolicyDataApplyError(id int32) *missingPolicyDataApplyError {
	return &missingPolicyDataApplyError{
		errorLink: errorLink{id: missingPolicyDataApplyErrorID},
		id:        id}
}

func (e *missingPolicyDataApplyError) Error() string {
	return e.errorf("Request %d doesn't contain parsed policy or parsed policy update", e.id)
}

type missingContentDataApplyError struct {
	errorLink
	id  int32
	cid string
}

func newMissingContentDataApplyError(id int32, cid string) *missingContentDataApplyError {
	return &missingContentDataApplyError{
		errorLink: errorLink{id: missingContentDataApplyErrorID},
		id:        id,
		cid:       cid}
}

func (e *missingContentDataApplyError) Error() string {
	return e.errorf("Request %d doesn't contain parsed content %q or parsed content update", e.id, e.cid)
}

type contentTransactionCreationError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newContentTransactionCreationError(id int32, v *item, err error) *contentTransactionCreationError {
	return &contentTransactionCreationError{
		errorLink: errorLink{id: contentTransactionCreationErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *contentTransactionCreationError) Error() string {
	return e.errorf("Can't create transaction for content %q update %d from tag %q to %q: %s", e.v.id, e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type contentUpdateParseError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newContentUpdateParseError(id int32, v *item, err error) *contentUpdateParseError {
	return &contentUpdateParseError{
		errorLink: errorLink{id: contentUpdateParseErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *contentUpdateParseError) Error() string {
	return e.errorf("Failed to parse content %q update %d from tag %q to %q: %s", e.v.id, e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type contentUpdateApplicationError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newContentUpdateApplicationError(id int32, v *item, err error) *contentUpdateApplicationError {
	return &contentUpdateApplicationError{
		errorLink: errorLink{id: contentUpdateApplicationErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *contentUpdateApplicationError) Error() string {
	return e.errorf("Failed to apply content %q update %d from tag %q to %q: %s", e.v.id, e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type contentUpdateUploadStoreError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newContentUpdateUploadStoreError(id int32, v *item, err error) *contentUpdateUploadStoreError {
	return &contentUpdateUploadStoreError{
		errorLink: errorLink{id: contentUpdateUploadStoreErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *contentUpdateUploadStoreError) Error() string {
	return e.errorf("Failed to store parsed content %q update %d from tag %q to %q: %s", e.v.id, e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type contentTransactionCommitError struct {
	errorLink
	id  int32
	v   *item
	err error
}

func newContentTransactionCommitError(id int32, v *item, err error) *contentTransactionCommitError {
	return &contentTransactionCommitError{
		errorLink: errorLink{id: contentTransactionCommitErrorID},
		id:        id,
		v:         v,
		err:       err}
}

func (e *contentTransactionCommitError) Error() string {
	return e.errorf("Failed to commit content %q transaction %d from tag %q to %q: %s", e.v.id, e.id, e.v.fromTag.String(), e.v.toTag.String(), e.err)
}

type unknownUploadedRequestError struct {
	errorLink
	id int32
}

func newUnknownUploadedRequestError(id int32) *unknownUploadedRequestError {
	return &unknownUploadedRequestError{
		errorLink: errorLink{id: unknownUploadedRequestErrorID},
		id:        id}
}

func (e *unknownUploadedRequestError) Error() string {
	return e.errorf("Can't find parsed policy or content with id %d", e.id)
}

type unsupportedPolicyFromatError struct {
	errorLink
	format string
}

func newUnsupportedPolicyFromatError(format string) *unsupportedPolicyFromatError {
	return &unsupportedPolicyFromatError{
		errorLink: errorLink{id: unsupportedPolicyFromatErrorID},
		format:    format}
}

func (e *unsupportedPolicyFromatError) Error() string {
	return e.errorf("The %s policy format is unsupported. Must be YAML or JSON", e.format)
}

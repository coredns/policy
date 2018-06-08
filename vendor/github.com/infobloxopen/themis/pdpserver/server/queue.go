package server

import (
	"math"
	"sync"

	"github.com/google/uuid"

	"github.com/infobloxopen/themis/pdp"
)

type item struct {
	policy bool
	id     string

	fromTag *uuid.UUID
	toTag   *uuid.UUID

	p  *pdp.PolicyStorage
	pt *pdp.PolicyStorageTransaction

	c  *pdp.LocalContent
	ct *pdp.LocalContentStorageTransaction
}

type queue struct {
	sync.Mutex

	idx   int32
	items map[int32]*item
}

func newQueue() *queue {
	return &queue{
		idx:   -1,
		items: make(map[int32]*item)}
}

func newPolicyItem(fromTag, toTag *uuid.UUID) *item {
	return &item{
		policy:  true,
		fromTag: fromTag,
		toTag:   toTag}
}

func newContentItem(id string, fromTag, toTag *uuid.UUID) *item {
	return &item{
		policy:  false,
		id:      id,
		fromTag: fromTag,
		toTag:   toTag}
}

func (q *queue) push(v *item) (int32, error) {
	q.Lock()
	defer q.Unlock()

	if q.idx >= math.MaxInt32 {
		return q.idx, newQueueOverflowError(q.idx)
	}

	q.idx++
	q.items[q.idx] = v

	return q.idx, nil
}

func (q *queue) pop(idx int32) (*item, bool) {
	q.Lock()
	defer q.Unlock()

	v, ok := q.items[idx]
	if ok {
		delete(q.items, idx)
	}

	return v, ok
}

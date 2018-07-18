package perf

import (
	"sync"

	"github.com/infobloxopen/themis/pdp"
)

type obligationsPool struct {
	p *sync.Pool
}

func makeObligationsPool(count uint32) obligationsPool {
	return obligationsPool{
		p: &sync.Pool{
			New: func() interface{} {
				return make([]pdp.AttributeAssignment, count)
			},
		},
	}
}

func (p obligationsPool) get() []pdp.AttributeAssignment {
	return p.p.Get().([]pdp.AttributeAssignment)
}

func (p obligationsPool) put(o []pdp.AttributeAssignment) {
	p.p.Put(o)
}

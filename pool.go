package policy

import (
	"sync"

	"github.com/infobloxopen/themis/pdp"
)

type attrPool struct {
	s int
	a *sync.Pool
}

func makeAttrPool(size int, dummy bool) attrPool {
	p := attrPool{s: size}
	if !dummy {
		p.a = &sync.Pool{
			New: func() interface{} {
				return p.newAttrs()
			},
		}
	}

	return p
}

func (p attrPool) newAttrs() []pdp.AttributeAssignment {
	return make([]pdp.AttributeAssignment, p.s)
}

func (p attrPool) Get() []pdp.AttributeAssignment {
	if p.a != nil {
		return p.a.Get().([]pdp.AttributeAssignment)
	}

	return p.newAttrs()
}

func (p attrPool) Put(a []pdp.AttributeAssignment) {
	if p.a != nil {
		p.a.Put(a)
	}
}

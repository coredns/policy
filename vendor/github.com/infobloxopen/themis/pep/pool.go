package pep

import "sync"

type bytePool struct {
	s int
	b *sync.Pool
}

func makeBytePool(size int, dummy bool) bytePool {
	p := bytePool{s: size}
	if !dummy {
		p.b = &sync.Pool{
			New: func() interface{} {
				return p.newBytes()
			},
		}
	}

	return p
}

func (p bytePool) newBytes() []byte {
	return make([]byte, p.s)
}

func (p bytePool) Get() []byte {
	if p.b != nil {
		return p.b.Get().([]byte)
	}

	return p.newBytes()
}

func (p bytePool) Put(b []byte) {
	if p.b != nil {
		p.b.Put(b)
	}
}

package pep

import (
	"golang.org/x/net/context"

	"github.com/allegro/bigcache"
	pb "github.com/infobloxopen/themis/pdp-service"
	ps "github.com/infobloxopen/themis/pdpserver/server"
)

type builtinClient struct {
	s *ps.PDPService

	pool bytePool

	cache *bigcache.BigCache

	opts options
}

func newBuiltinClient(opts options) *builtinClient {
	s := ps.NewBuiltinPDPService(opts.policyFile, opts.contentFiles)
	c := &builtinClient{
		s:    s,
		opts: opts,
	}

	if !opts.autoRequestSize {
		c.pool = makeBytePool(int(opts.maxRequestSize), opts.noPool)
	}

	return c
}

func (c *builtinClient) Connect(addr string) error {
	cache, err := newCacheFromOptions(c.opts)
	if err != nil {
		return err
	}
	c.cache = cache

	return nil
}

func (c *builtinClient) Close() {
	if c.cache != nil {
		c.cache.Reset()
		c.cache = nil
	}

	c.s = nil
}

func (c *builtinClient) Validate(in, out interface{}) error {
	if c.s == nil {
		return ErrorNotConnected
	}

	var (
		req pb.Msg
		err error
	)

	if c.opts.autoRequestSize {
		req, err = makeRequest(in)
	} else {
		var b []byte
		switch in.(type) {
		default:
			b = c.pool.Get()
			defer c.pool.Put(b)

		case []byte, pb.Msg, *pb.Msg:
		}

		req, err = makeRequestWithBuffer(in, b)
	}
	if err != nil {
		return err
	}

	if c.cache != nil {
		if b, err := c.cache.Get(string(req.Body)); err == nil {
			return fillResponse(pb.Msg{Body: b}, out)
		}
	}

	res, err := c.s.Validate(context.Background(), &req)
	if err != nil {
		return err
	}

	if c.cache != nil {
		c.cache.Set(string(req.Body), res.Body)
	}

	return fillResponse(*res, out)
}

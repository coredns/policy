package pep

import (
	"fmt"
	"sync"

	"github.com/allegro/bigcache"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	ot "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/infobloxopen/themis/pdp-service"
)

type unaryClient struct {
	lock   *sync.RWMutex
	conn   *grpc.ClientConn
	client *pb.PDPClient

	pool bytePool

	cache *bigcache.BigCache

	opts options
}

func newUnaryClient(opts options) *unaryClient {
	return &unaryClient{
		lock: &sync.RWMutex{},
		opts: opts,
		pool: makeBytePool(int(opts.maxRequestSize), opts.noPool),
	}
}

func (c *unaryClient) Connect(addr string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn != nil {
		return ErrorConnected
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if len(c.opts.addresses) > 0 {
		addr = virtualServerAddress
		switch c.opts.balancer {
		default:
			panic(fmt.Errorf("invalid balancer %d", c.opts.balancer))

		case roundRobinBalancer:
			opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(newStaticResolver(addr, c.opts.addresses...))))

		case hotSpotBalancer:
			return ErrorHotSpotBalancerUnsupported
		}
	}

	if c.opts.tracer != nil {
		opts = append(opts,
			grpc.WithUnaryInterceptor(
				otgrpc.OpenTracingClientInterceptor(
					c.opts.tracer,
					otgrpc.IncludingSpans(
						func(parentSpanCtx ot.SpanContext, method string, req, resp interface{}) bool {
							return parentSpanCtx != nil
						},
					),
				),
			),
		)
	}

	cache, err := newCacheFromOptions(c.opts)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}

	c.conn = conn
	c.cache = cache

	client := pb.NewPDPClient(c.conn)
	c.client = &client

	return nil
}

func (c *unaryClient) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	if c.cache != nil {
		c.cache.Reset()
		c.cache = nil
	}

	c.client = nil
}

func (c *unaryClient) Validate(in, out interface{}) error {
	c.lock.RLock()
	uc := c.client
	c.lock.RUnlock()

	if uc == nil {
		return ErrorNotConnected
	}

	var b []byte
	switch in.(type) {
	default:
		b = c.pool.Get()
		defer c.pool.Put(b)

	case []byte, pb.Msg, *pb.Msg:
	}

	req, err := makeRequest(in, b)
	if err != nil {
		return err
	}

	if c.cache != nil {
		if b, err := c.cache.Get(string(req.Body)); err == nil {
			return fillResponse(pb.Msg{Body: b}, out)
		}
	}

	res, err := (*uc).Validate(context.Background(), &req, grpc.FailFast(false))
	if err != nil {
		return err
	}

	if c.cache != nil {
		c.cache.Set(string(req.Body), res.Body)
	}

	return fillResponse(*res, out)
}

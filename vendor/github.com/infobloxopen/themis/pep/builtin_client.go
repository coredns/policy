package pep

import (
	"fmt"
	"golang.org/x/net/context"

	pb "github.com/infobloxopen/themis/pdp-service"
	ps "github.com/infobloxopen/themis/pdpserver/server"
)

type builtinClient struct {
	s    *ps.PDPService
	pool bytePool
}

func NewBuiltinClient(policyFile string, contentFiles []string) *builtinClient {
	fmt.Printf("pep client NewBuiltinClient() called..........\n")
	s := ps.NewPDPService(policyFile, contentFiles)
	return &builtinClient{
		s: s,
	}
}

func (c *builtinClient) Connect(addr string) error {
	fmt.Printf("pep client Connect() called..........\n")
	return nil
}

func (c *builtinClient) Close() {
	fmt.Printf("pep client Close() called..........\n")
}

func (c *builtinClient) Validate(in, out interface{}) error {
	fmt.Printf("pep client Validate() called..........\n")
	if c.s == nil {
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

	res, err := c.s.Validate(context.Background(), &req)
	if err != nil {
		return err
	}

	return fillResponse(*res, out)
}


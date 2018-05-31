package policy

import (
	"fmt"

	pdp "github.com/infobloxopen/themis/pdp-service"
)

type testClient struct {
	nextResponse   *pdp.Response
	nextResponseIP *pdp.Response
	errResponse    error
	errResponseIP  error
}

func newTestClientInit(nextResponse *pdp.Response, nextResponseIP *pdp.Response,
	errResponse error, errResponseIP error) *testClient {
	return &testClient{
		nextResponse:   nextResponse,
		nextResponseIP: nextResponseIP,
		errResponse:    errResponse,
		errResponseIP:  errResponseIP,
	}
}

func (c *testClient) Connect(addr string) error { return nil }
func (c *testClient) Close()                    {}
func (c *testClient) Validate(in, out interface{}) error {
	if in != nil {
		p := in.(pdp.Request)
		for _, a := range p.Attributes {
			if a.Id == attrNameAddress {
				if c.errResponseIP != nil {
					return c.errResponseIP
				}
				if c.nextResponseIP != nil {
					return fillResponse(c.nextResponseIP, out)
				}
				continue
			}
		}
	}
	if c.errResponse != nil {
		return c.errResponse
	}
	return fillResponse(c.nextResponse, out)
}

func fillResponse(in *pdp.Response, out interface{}) error {
	r, ok := out.(*pdp.Response)
	if !ok {
		return fmt.Errorf("testClient can only translate response to *Response type but got %T", out)
	}
	r.Effect = in.Effect
	r.Obligation = in.Obligation
	return nil
}

package pep

import (
	"fmt"
	"golang.org/x/net/context"

	//pb "github.com/infobloxopen/themis/pdp-service"
	ps "github.com/infobloxopen/themis/pdpserver/server"
)

type integratedClient struct {
	s      *ps.Server
}

func NewIntegratedClient(policyFile string, contentFiles []string) *integratedClient {
	fmt.Printf("pep client NewIntegratedClient() called..........\n")
	s := ps.NewIntegratedServer(policyFile, contentFiles)
	return &integratedClient{
		s: s,
	}
}

func (c *integratedClient) Connect(addr string) error {
	fmt.Printf("pep client Connect() called..........\n")
	return nil
}

func (c *integratedClient) Close() {
	fmt.Printf("pep client Close() called..........\n")
}

func (c *integratedClient) Validate(in, out interface{}) error {
	fmt.Printf("pep client Validate() called..........\n")
	if c.s == nil {
		return ErrorNotConnected
	}

	req, err := makeRequest(in)
	if err != nil {
		return err
	}

	//func (s *Server) Validate(ctx context.Context, in *pb.Request) (*pb.Response, error) {

	res, err := c.s.Validate(context.Background(), &req)
	if err != nil {
		return err
	}

	return fillResponse(res, out)
}

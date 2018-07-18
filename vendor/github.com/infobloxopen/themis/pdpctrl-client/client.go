// Package pdpcc implements gRPC client to control Policy Decision Point (PDP)
// server. It wraps control part of golang gRPC protocol implementation for
// PDP. The protocol is defined by
// github.com/infobloxopen/themis/proto/control.proto. Its golang implementation
// can be found at github.com/infobloxopen/themis/pdp-control.
package pdpcc

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/infobloxopen/themis/pdp-control && protoc -I $GOPATH/src/github.com/infobloxopen/themis/proto/ $GOPATH/src/github.com/infobloxopen/themis/proto/control.proto --go_out=plugins=grpc:$GOPATH/src/github.com/infobloxopen/themis/pdp-control && ls $GOPATH/src/github.com/infobloxopen/themis/pdp-control"

import (
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/infobloxopen/themis/pdp-control"
)

// TagError is the error returned by upload request if policies or content
// update tag doesn't match to the tag PDP server has.
type TagError struct {
	tag string
}

// Error implements method of error interface.
func (e *TagError) Error() string {
	return e.tag
}

// Client structure represents client side of PDP control protocol. It's
// responsible for establishing connection and uploading data to PDP server.
type Client struct {
	address   string
	chunkSize int

	conn   *grpc.ClientConn
	client pb.PDPControlClient
}

// NewClient function creates new instance of Client structure.
func NewClient(addr string, chunkSize int) *Client {
	return &Client{
		address:   addr,
		chunkSize: chunkSize,
	}
}

// Connect establishes connection to PDP server.
func (c *Client) Connect(timeout time.Duration) error {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout))
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = pb.NewPDPControlClient(c.conn)

	return nil
}

// Close terminates established connection. It does nothing if there is no
// connection.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}

	c.client = nil
}

// RequestPoliciesUpload makes request to upload policies. On success returns
// id of the request. The id should be used to make actual upload. Arguments
// fromTag and toTag should be valid text representaions of UUID or empty
// strings. If fromTag is emply, server expects full policy on upload otherwise
// it expects update. If toTag is empty, policy becomes not tagged and can't be
// updated incrementally. If incremental update is not possible because fromTag
// doesn't match to current server policies tag, the method returns TagError.
func (c *Client) RequestPoliciesUpload(fromTag, toTag string) (int32, error) {
	return c.request(&pb.Item{
		Type:    pb.Item_POLICIES,
		FromTag: fromTag,
		ToTag:   toTag})
}

// RequestContentUpload requests content upload. The method returns request's
// id which should be used on upload call. Argument id is content identifier.
// It must be equal to id field of full content representation. As for policies
// fromTag and toTag should be valid text representaion of UUID or empty string.
// If fromTag is emply server expects full content on upload otherwise it
// expects update. If toTag is empty content becomes not tagged and can't be
// updated incrementally. TagError is returned to indicate that fromTag doesn't
// match to current server's content with the same content id.
func (c *Client) RequestContentUpload(id, fromTag, toTag string) (int32, error) {
	return c.request(&pb.Item{
		Type:    pb.Item_CONTENT,
		FromTag: fromTag,
		ToTag:   toTag,
		Id:      id})
}

// Upload implements actual data transfer to PDP server. It streams all data
// from given reader. Argument id should be valid request id obtained previously
// from RequestPoliciesUpload or RequestContentUpload call. On success
// Upload method returns upload id (which is not equal to request id).
// Obtained id shuld be used in next Apply call.
func (c *Client) Upload(id int32, r io.Reader) (int32, error) {
	u, err := c.client.Upload(context.Background())
	if err != nil {
		return -1, err
	}

	p := make([]byte, c.chunkSize)
	for {
		n, err := r.Read(p)
		if n > 0 {
			chunk := &pb.Chunk{
				Id:   id,
				Data: string(p[:n])}
			if err := u.Send(chunk); err != nil {
				return -1, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			c.closeUpload(u)
			return -1, err
		}
	}

	return c.closeUpload(u)
}

// Apply requests server to switch to recently apploaded policies or content.
// Its id argument should be upload id obtained on previous Upload call.
func (c *Client) Apply(id int32) error {
	r, err := c.client.Apply(context.Background(), &pb.Update{Id: id})
	if err != nil {
		return err
	}

	if r.Status != pb.Response_ACK {
		return errors.New(r.Details)
	}

	return nil
}

// NotifyReady set server to 'ready' state -
// after that server will open service port for serve decision requests
func (c *Client) NotifyReady() error {
	r, err := c.client.NotifyReady(context.Background(), &pb.Empty{})
	if err != nil {
		return err
	}

	if r.Status != pb.Response_ACK {
		return errors.New(r.Details)
	}

	return nil
}

func (c *Client) request(item *pb.Item) (int32, error) {
	r, err := c.client.Request(context.Background(), item)
	if err != nil {
		return -1, err
	}

	switch r.Status {
	case pb.Response_ACK:
		return r.Id, nil

	case pb.Response_ERROR:
		return -1, errors.New(r.Details)

	case pb.Response_TAG_ERROR:
		return -1, &TagError{tag: r.Details}
	}

	return -1, fmt.Errorf("unknown response statue: %d", r.Status)
}

func (c *Client) closeUpload(u pb.PDPControl_UploadClient) (int32, error) {
	r, err := u.CloseAndRecv()
	if err != nil {
		return -1, err
	}

	if r.Status != pb.Response_ACK {
		return -1, errors.New(r.Details)
	}

	return r.Id, nil
}

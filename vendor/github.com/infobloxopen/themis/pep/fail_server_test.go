package pep

import (
	"errors"
	"io"
	"net"
	"reflect"
	"strings"
	"sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

const (
	failID = "fail"
	IDID   = "id"

	thisRequest = "this"
)

var errRequested = errors.New("failed as requested by client")

type failServer struct {
	ID       uint64
	failNext int32
	s        *grpc.Server
}

func newFailServer(addr string) (*failServer, error) {
	if err := waitForPortClosed(addr); err != nil {
		return nil, err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &failServer{s: grpc.NewServer()}
	pb.RegisterPDPServer(s.s, s)
	go s.s.Serve(ln)

	if err := waitForPortOpened(addr); err != nil {
		s.Stop()
		return nil, err
	}

	return s, nil
}

func (s *failServer) Stop() {
	s.s.Stop()
}

func (s *failServer) Validate(ctx context.Context, in *pb.Msg) (*pb.Msg, error) {
	reqID := atomic.AddUint64(&s.ID, 1)

	targetID, fail := parseFailRequest(in)
	if fail == thisRequest && reqID == targetID {
		return nil, errRequested
	}

	return &pb.Msg{
		Body: append(
			[]byte{1, 0, 1, 0, 0},
			in.Body[2:]...,
		),
	}, nil
}

func (s *failServer) NewValidationStream(stream pb.PDP_NewValidationStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		reqID := atomic.AddUint64(&s.ID, 1)
		targetID, fail := parseFailRequest(in)
		if fail == thisRequest && reqID == targetID {
			return errRequested
		}

		err = stream.Send(&pb.Msg{
			Body: append(
				[]byte{1, 0, 1, 0, 0},
				in.Body[2:]...,
			),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func parseFailRequest(in *pb.Msg) (uint64, string) {
	var (
		targetID uint64
		fail     string
	)

	err := pdp.UnmarshalRequestReflection(in.Body, func(id string, t pdp.Type) (reflect.Value, error) {
		switch strings.ToLower(id) {
		case IDID:
			if t == pdp.TypeInteger {
				return reflect.ValueOf(&targetID).Elem(), nil
			}

		case failID:
			if t == pdp.TypeString {
				return reflect.ValueOf(&fail).Elem(), nil
			}
		}

		return reflect.ValueOf(nil), nil
	})
	if err != nil {
		panic(err)
	}

	return targetID, fail
}

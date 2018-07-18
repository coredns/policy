package server

import (
	"io"

	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-control"
)

func controlFail(err error) *pb.Response {
	status := pb.Response_ERROR
	switch e := err.(type) {
	case *tagCheckError:
		switch e.err.(type) {
		case *pdp.UntaggedPolicyModificationError, *pdp.MissingPolicyTagError, *pdp.PolicyTagsNotMatchError, *pdp.UntaggedContentModificationError, *pdp.MissingContentTagError, *pdp.ContentTagsNotMatchError:
			status = pb.Response_TAG_ERROR
		}

	case *policyTransactionCreationError:
		switch e.err.(type) {
		case *pdp.UntaggedPolicyModificationError, *pdp.MissingPolicyTagError, *pdp.PolicyTagsNotMatchError:
			status = pb.Response_TAG_ERROR
		}
	case *contentTransactionCreationError:
		switch e.err.(type) {
		case *pdp.UntaggedContentModificationError, *pdp.MissingContentTagError, *pdp.ContentTagsNotMatchError:
			status = pb.Response_TAG_ERROR
		}
	}

	return &pb.Response{
		Status:  status,
		Id:      -1,
		Details: err.Error()}
}

func newTag(s string) (*uuid.UUID, error) {
	if len(s) <= 0 {
		return nil, nil
	}

	t, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Request is a server handler for gRPC call
// It serves PAP control requests
func (s *Server) Request(ctx context.Context, in *pb.Item) (*pb.Response, error) {
	s.opts.logger.Info("Got new control request")

	fromTag, err := newTag(in.FromTag)
	if err != nil {
		return controlFail(newInvalidFromTagError(in.FromTag, err)), nil
	}

	toTag, err := newTag(in.ToTag)
	if err != nil {
		return controlFail(newInvalidToTagError(in.ToTag, err)), nil
	}

	if fromTag != nil && toTag == nil {
		return controlFail(newInvalidTagsError(in.FromTag)), nil
	}

	var id int32
	switch in.Type {
	default:
		return controlFail(newUnknownUploadRequestError(in.Type)), nil

	case pb.Item_POLICIES:
		id, err = s.policyRequest(fromTag, toTag)

	case pb.Item_CONTENT:
		id, err = s.contentRequest(in.Id, fromTag, toTag)
	}

	if err != nil {
		return controlFail(err), nil
	}

	return &pb.Response{Status: pb.Response_ACK, Id: id}, nil
}

func (s *Server) getHead(stream pb.PDPControl_UploadServer) (int32, *streamReader, error) {
	chunk, err := stream.Recv()
	if err == io.EOF {
		return 0, nil, stream.SendAndClose(controlFail(newEmptyUploadError()))
	}

	if err != nil {
		return 0, nil, err
	}

	return chunk.Id, newStreamReader(chunk.Id, chunk.Data, stream, s.opts.logger), nil
}

// Upload is a server handler for gRPC call
// It uploads data from PAP and save it to PDP
func (s *Server) Upload(stream pb.PDPControl_UploadServer) error {
	s.opts.logger.Info("Got new data stream")

	id, r, err := s.getHead(stream)
	if r == nil {
		return err
	}

	req, ok := s.q.pop(id)
	if !ok {
		s.opts.logger.WithField("id", id).Error("no such request")
		err := r.skip()
		if err != nil {
			return err
		}

		return stream.SendAndClose(controlFail(newUnknownUploadError(id)))
	}

	if req.fromTag == nil {
		if req.policy {
			err = s.uploadPolicy(id, r, req, stream)
		} else {
			err = s.uploadContent(id, r, req, stream)
		}
	} else {
		if req.policy {
			err = s.uploadPolicyUpdate(id, r, req, stream)
		} else {
			err = s.uploadContentUpdate(id, r, req, stream)
		}
	}

	return err
}

// Apply is a server handler for gRPC call
// It applies data previously saved in PDP
func (s *Server) Apply(ctx context.Context, in *pb.Update) (*pb.Response, error) {
	s.opts.logger.Info("Got apply command")

	req, ok := s.q.pop(in.Id)
	if !ok {
		s.opts.logger.WithField("id", in.Id).Error("no such request")
		return controlFail(newUnknownUploadedRequestError(in.Id)), nil
	}

	var (
		res *pb.Response
		err error
	)
	if req.policy {
		res, err = s.applyPolicy(in.Id, req)
	} else {
		res, err = s.applyContent(in.Id, req)
	}

	return res, err
}

// NotifyReady is a server handler for gRPC call
// It starts handling decision requests
func (s *Server) NotifyReady(ctx context.Context, m *pb.Empty) (*pb.Response, error) {
	s.opts.logger.Info("Got notified about readiness")

	go s.startOnce.Do(func() {
		s.errCh <- s.serveRequests()
	})

	return &pb.Response{Status: pb.Response_ACK}, nil
}

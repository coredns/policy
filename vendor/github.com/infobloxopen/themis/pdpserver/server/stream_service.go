package server

import (
	"io"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/infobloxopen/themis/pdp-service"
)

var streamAutoIncrement uint64

// NewValidationStream is a server handler for gRPC call
// It creates new gRPC stream and handles PDP decision requests using it
func (s *Server) NewValidationStream(stream pb.PDP_NewValidationStreamServer) error {
	ctx := stream.Context()

	sID := atomic.AddUint64(&streamAutoIncrement, 1)
	s.opts.logger.WithField("id", sID).Debug("Got new stream")

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			if err := ctx.Err(); err != nil && (err == context.Canceled || err == context.DeadlineExceeded) {
				break
			}

			s.opts.logger.WithFields(log.Fields{
				"id":  sID,
				"err": err,
			}).Error("Failed to read next request from stream. Dropping stream...")

			return err
		}

		out, err := s.Validate(context.Background(), in)
		if err != nil {
			s.opts.logger.WithFields(log.Fields{
				"id":  sID,
				"err": err,
			}).Panic("Failed to validate request")
		}

		err = stream.Send(out)
		if err != nil {
			s.opts.logger.WithFields(log.Fields{
				"id":  sID,
				"err": err,
			}).Error("Failed to send response. Dropping stream...")

			return err
		}
	}

	s.opts.logger.WithField("id", sID).Debug("Stream deleted")
	return nil
}

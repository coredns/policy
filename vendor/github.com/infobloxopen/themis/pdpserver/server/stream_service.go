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

	buffer := make([]byte, s.opts.maxResponseSize)

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

		s.RLock()
		p := s.p
		c := s.c
		s.RUnlock()

		if s.opts.autoResponseSize {
			err = stream.Send(&pb.Msg{Body: s.rawValidateWithAllocator(p, c, in.Body, func(n int) ([]byte, error) {
				if len(buffer) < n {
					buffer = make([]byte, n)
				}

				return buffer, nil
			})})
		} else {
			err = stream.Send(&pb.Msg{Body: s.rawValidateToBuffer(p, c, in.Body, buffer)})
		}
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

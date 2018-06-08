package server

import (
	"io"

	log "github.com/sirupsen/logrus"

	pb "github.com/infobloxopen/themis/pdp-control"
)

type streamReader struct {
	id     int32
	stream pb.PDPControl_UploadServer
	chunk  []byte
	offset int
	eof    bool
	logger *log.Logger
}

func newStreamReader(id int32, head string, stream pb.PDPControl_UploadServer, logger *log.Logger) *streamReader {
	return &streamReader{
		id:     id,
		stream: stream,
		chunk:  []byte(head),
		logger: logger}
}

func (r *streamReader) skip() error {
	if r.eof {
		return nil
	}

	for {
		_, err := r.stream.Recv()
		if err == io.EOF {
			r.eof = true
			break
		}

		if err != nil {
			r.logger.WithFields(log.Fields{
				"id":    r.id,
				"error": err}).Error("failed to read data stream")
			return err
		}
	}

	return nil
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	if r.eof {
		return 0, io.EOF
	}

	if len(p) <= 0 {
		return 0, nil
	}

	offset := 0
	req := len(p) - offset
	rem := len(r.chunk) - r.offset
	for req > rem {
		for i := 0; i < rem; i++ {
			p[offset+i] = r.chunk[r.offset+i]
		}

		offset += rem
		req -= rem
		r.offset = 0

		chunk, err := r.stream.Recv()
		if err == io.EOF {
			r.eof = true
			return offset, io.EOF
		}

		if err != nil {
			r.logger.WithFields(log.Fields{
				"id":    r.id,
				"error": err}).Error("failed to read data stream")
			return offset, err
		}

		r.chunk = []byte(chunk.Data)

		rem = len(r.chunk)
	}

	for i := 0; i < req; i++ {
		p[offset+i] = r.chunk[r.offset+i]
	}

	r.offset += req

	return offset + req, nil
}

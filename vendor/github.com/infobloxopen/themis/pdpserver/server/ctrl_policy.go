package server

import (
	pb "github.com/infobloxopen/themis/pdp-control"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *Server) policyRequest(fromTag, toTag *uuid.UUID) (int32, error) {
	if fromTag != nil {
		s.RLock()
		p := s.p
		s.RUnlock()

		err := p.CheckTag(fromTag)
		if err != nil {
			return 0, newTagCheckError(err)
		}
	}

	return s.q.push(newPolicyItem(fromTag, toTag))
}

func (s *Server) uploadPolicy(id int32, r *streamReader, req *item, stream pb.PDPControl_UploadServer) error {
	p, err := s.opts.parser.Unmarshal(r, req.toTag)
	if err != nil {
		return stream.SendAndClose(controlFail(newPolicyUploadParseError(id, err)))
	}

	req.p = p
	nid, err := s.q.push(req)
	if err != nil {
		return stream.SendAndClose(controlFail(newPolicyUploadStoreError(id, err)))
	}

	return stream.SendAndClose(&pb.Response{Status: pb.Response_ACK, Id: nid})
}

func (s *Server) uploadPolicyUpdate(id int32, r *streamReader, req *item, stream pb.PDPControl_UploadServer) error {
	s.RLock()
	if s.p == nil {
		s.RUnlock()
		r.skip()
		return stream.SendAndClose(controlFail(newMissingPolicyStorageError()))
	}

	t, err := s.p.NewTransaction(req.fromTag)
	if err != nil {
		s.RUnlock()
		r.skip()
		return stream.SendAndClose(controlFail(newPolicyTransactionCreationError(id, req, err)))
	}
	s.RUnlock()

	u, err := s.opts.parser.UnmarshalUpdate(r, t.Symbols(), *req.fromTag, *req.toTag)
	if err != nil {
		return stream.SendAndClose(controlFail(newPolicyUpdateParseError(id, req, err)))
	}

	s.opts.logger.WithField("update", u).Debug("Policy update")

	err = t.Apply(u)
	if err != nil {
		return stream.SendAndClose(controlFail(newPolicyUpdateApplicationError(id, req, err)))
	}

	req.pt = t
	nid, err := s.q.push(req)
	if err != nil {
		return stream.SendAndClose(controlFail(newPolicyUpdateUploadStoreError(id, req, err)))
	}

	return stream.SendAndClose(&pb.Response{Status: pb.Response_ACK, Id: nid})
}

func (s *Server) applyPolicy(id int32, req *item) (*pb.Response, error) {
	if req.p != nil {
		s.Lock()
		s.p = req.p
		s.Unlock()

		if req.toTag == nil {
			s.opts.logger.WithField("id", id).Info("New policy has been applied")
		} else {
			s.opts.logger.WithFields(log.Fields{
				"id":  id,
				"tag": req.toTag.String()}).Info("New policy has been applied")
		}

		return &pb.Response{Status: pb.Response_ACK, Id: id}, nil
	}

	if req.pt != nil {
		p, err := req.pt.Commit()
		if err != nil {
			return controlFail(newPolicyTransactionCommitError(id, req, err)), nil
		}

		s.Lock()
		s.p = p
		s.Unlock()

		s.opts.logger.WithFields(log.Fields{
			"id":       id,
			"prev-tag": req.fromTag,
			"curr-tag": req.toTag}).Info("Policy update has been applied")

		return &pb.Response{Status: pb.Response_ACK, Id: id}, nil
	}

	return controlFail(newMissingPolicyDataApplyError(id)), nil
}

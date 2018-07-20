package server

import (
	"fmt"
	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
	"github.com/infobloxopen/themis/pdp/ast"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type PDPService struct {
	sync.RWMutex
	opts options

	p *pdp.PolicyStorage
	c *pdp.LocalContentStorage
}

func NewBuiltinPDPService(policyFile string, contentFiles []string) *PDPService {
	o := options{
		logger:              log.StandardLogger(),
		service:             ":5555",
		memStatsLogInterval: -1 * time.Second,
		maxResponseSize:     10240,
	}

	ext := filepath.Ext(policyFile)
	switch ext {
	case ".json":
		o.parser = ast.NewJSONParser()
	case ".yaml":
		o.parser = ast.NewYAMLParser()
	}

	if o.parser == nil {
		o.parser = ast.NewYAMLParser()
	}

	s := &PDPService{
		c:    pdp.NewLocalContentStorage(nil),
		opts: o,
	}

	log.SetLevel(log.DebugLevel)
	err := s.LoadPolicies(policyFile)
	if err != nil {
		return nil
	}

	if contentFiles != nil && len(contentFiles) > 0 {
		err = s.LoadContent(contentFiles)
		if err != nil {
			return nil
		}
	}

	return s
}

func (s *PDPService) newContext(c *pdp.LocalContentStorage, in []byte) (*pdp.Context, error) {
	ctx, err := pdp.NewContextFromBytes(c, in)
	if err != nil {
		return nil, newContextCreationError(err)
	}

	return ctx, nil
}

func makeFailureResponse(b []byte, err error) []byte {
	n, err := pdp.MakeIndeterminateResponse(b, err)
	if err != nil {
		panic(err)
	}

	return b[:n]
}

func (s *PDPService) rawValidate(p *pdp.PolicyStorage, c *pdp.LocalContentStorage, in []byte, out []byte) []byte {
	if out == nil {
		out = make([]byte, s.opts.maxResponseSize)
	}

	if p == nil {
		return makeFailureResponse(out, newMissingPolicyError())
	}

	ctx, err := s.newContext(c, in)
	if err != nil {
		return makeFailureResponse(out, err)
	}

	if s.opts.logger.Level >= log.DebugLevel {
		s.opts.logger.WithField("context", ctx).Debug("Request context")
	}

	r := p.Root().Calculate(ctx)

	if s.opts.logger.Level >= log.DebugLevel {
		s.opts.logger.WithFields(log.Fields{
			"effect": pdp.EffectNameFromEnum(r.Effect),
			"reason": r.Status,
			"obligations": obligations{
				ctx: ctx,
				o:   r.Obligations,
			},
		}).Debug("Response")
	}

	n, err := r.Marshal(out, ctx)
	if err != nil {
		panic(err)
	}

	return out[:n]
}

// Validate is a server handler for gRPC call
// It handles PDP decision requests
func (s *PDPService) Validate(ctx context.Context, in *pb.Msg) (*pb.Msg, error) {
	s.RLock()
	p := s.p
	c := s.c
	s.RUnlock()

	return &pb.Msg{
		Body: s.rawValidate(p, c, in.Body, nil),
	}, nil
}

type obligations struct {
	ctx *pdp.Context
	o   []pdp.AttributeAssignment
}

func (o obligations) String() string {
	if len(o.o) <= 0 {
		return "no attributes"
	}

	lines := make([]string, len(o.o)+1)
	lines[0] = "attributes:"
	for i, e := range o.o {
		id, t, v, err := e.Serialize(o.ctx)
		if err != nil {
			lines[i+1] = fmt.Sprintf("- %d: %s", i+1, err)
		} else {
			lines[i+1] = fmt.Sprintf("- %s.(%s): %q", id, t, v)
		}
	}

	return strings.Join(lines, "\n")
}

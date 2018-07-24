package server

import (
	"fmt"
	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
	"github.com/infobloxopen/themis/pdp/ast"
	"github.com/infobloxopen/themis/pdp/jcon"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type PDPService struct {
	sync.RWMutex

	opts options

	p *pdp.PolicyStorage
	c *pdp.LocalContentStorage

	pool bytePool
}

func NewBuiltinPDPService(policyFile string, contentFiles []string) *PDPService {
	var opts options
	return NewPDPService(opts,
		WithPolicyFile(policyFile),
		WithContentFiles(contentFiles),
		WithLogger(log.StandardLogger()),
		WithMaxResponseSize(10240))
}

//func NewBuiltinPDPService(policyFile string, contentFiles []string, logger *log.Logger) *PDPService {
func NewPDPService(options options, addOpts ...Option) *PDPService {
	s := &PDPService{
		opts: options,
		c:    pdp.NewLocalContentStorage(nil),
	}

	for _, opt := range addOpts {
		opt(&s.opts)
	}

	if s.opts.policyFile != "" {
		ext := filepath.Ext(s.opts.policyFile)
		switch ext {
		case ".json":
			s.opts.parser = ast.NewJSONParser()
		case ".yaml":
			s.opts.parser = ast.NewYAMLParser()
		}
	}

	if s.opts.parser == nil {
		s.opts.parser = ast.NewYAMLParser()
	}

	log.SetLevel(log.DebugLevel)
	err := s.LoadPolicies(s.opts.policyFile)
	if err != nil {
		return nil
	}

	if s.opts.contentFiles != nil && len(s.opts.contentFiles) > 0 {
		err = s.LoadContent(s.opts.contentFiles)
		if err != nil {
			return nil
		}
	} else {
		s.c = pdp.NewLocalContentStorage(nil)
	}

	if !s.opts.autoResponseSize {
		s.pool = makeBytePool(int(s.opts.maxResponseSize), false)
	}

	return s
}

// LoadPolicies loads policies from file
func (s *PDPService) LoadPolicies(path string) error {
	if len(path) <= 0 {
		return nil
	}

	s.opts.logger.WithField("policy", path).Info("Loading policy")
	pf, err := os.Open(path)
	if err != nil {
		s.opts.logger.WithFields(log.Fields{"policy": path, "error": err}).Error("Failed load policy")
		return err
	}

	s.opts.logger.WithField("policy", path).Info("Parsing policy")
	p, err := s.opts.parser.Unmarshal(pf, nil)
	if err != nil {
		s.opts.logger.WithFields(log.Fields{"policy": path, "error": err}).Error("Failed parse policy")
		return err
	}

	s.p = p

	return nil
}

// LoadContent loads content from files
func (s *PDPService) LoadContent(paths []string) error {
	items := []*pdp.LocalContent{}
	for _, path := range paths {
		err := func() error {
			s.opts.logger.WithField("content", path).Info("Opening content")
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			defer f.Close()

			s.opts.logger.WithField("content", path).Info("Parsing content")
			item, err := jcon.Unmarshal(f, nil)
			if err != nil {
				return err
			}

			items = append(items, item)
			return nil
		}()
		if err != nil {
			return err
		}
	}

	s.c = pdp.NewLocalContentStorage(items)

	return nil
}

func (s *PDPService) newContext(c *pdp.LocalContentStorage, in []byte) (*pdp.Context, error) {
	ctx, err := pdp.NewContextFromBytes(c, in)
	if err != nil {
		return nil, newContextCreationError(err)
	}

	return ctx, nil
}

func makeFailureResponse(err error) []byte {
	b, err := pdp.MakeIndeterminateResponse(err)
	if err != nil {
		panic(err)
	}

	return b
}

func makeFailureResponseWithAllocator(f func(n int) ([]byte, error), err error) []byte {
	b, err := pdp.MakeIndeterminateResponseWithAllocator(f, err)
	if err != nil {
		panic(err)
	}

	return b
}

func makeFailureResponseWithBuffer(b []byte, err error) []byte {
	n, err := pdp.MakeIndeterminateResponseWithBuffer(b, err)
	if err != nil {
		panic(err)
	}

	return b[:n]
}

func (s *PDPService) rawValidate(p *pdp.PolicyStorage, c *pdp.LocalContentStorage, in []byte) []byte {
	if p == nil {
		return makeFailureResponse(newMissingPolicyError())
	}

	ctx, err := s.newContext(c, in)
	if err != nil {
		return makeFailureResponse(err)
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

	out, err := r.Marshal(ctx)
	if err != nil {
		panic(err)
	}

	return out
}

func (s *PDPService) rawValidateWithAllocator(p *pdp.PolicyStorage, c *pdp.LocalContentStorage, in []byte, f func(n int) ([]byte, error)) []byte {
	if p == nil {
		return makeFailureResponseWithAllocator(f, newMissingPolicyError())
	}

	ctx, err := s.newContext(c, in)
	if err != nil {
		return makeFailureResponseWithAllocator(f, err)
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

	out, err := r.MarshalWithAllocator(f, ctx)
	if err != nil {
		panic(err)
	}

	return out
}

func (s *PDPService) rawValidateToBuffer(p *pdp.PolicyStorage, c *pdp.LocalContentStorage, in []byte, out []byte) []byte {
	if p == nil {
		return makeFailureResponseWithBuffer(out, newMissingPolicyError())
	}

	ctx, err := s.newContext(c, in)
	if err != nil {
		return makeFailureResponseWithBuffer(out, err)
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

	n, err := r.MarshalToBuffer(out, ctx)
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

	if s.opts.autoResponseSize {
		return &pb.Msg{
			Body: s.rawValidate(p, c, in.Body),
		}, nil
	}

	b := s.pool.Get()
	defer s.pool.Put(b)

	return &pb.Msg{
		Body: s.rawValidateToBuffer(p, c, in.Body, b),
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

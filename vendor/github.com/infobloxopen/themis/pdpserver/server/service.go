package server

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
)

func makeEffect(effect int) (pb.Response_Effect, error) {
	switch effect {
	case pdp.EffectDeny:
		return pb.Response_DENY, nil

	case pdp.EffectPermit:
		return pb.Response_PERMIT, nil

	case pdp.EffectNotApplicable:
		return pb.Response_NOTAPPLICABLE, nil

	case pdp.EffectIndeterminate:
		return pb.Response_INDETERMINATE, nil

	case pdp.EffectIndeterminateD:
		return pb.Response_INDETERMINATED, nil

	case pdp.EffectIndeterminateP:
		return pb.Response_INDETERMINATEP, nil

	case pdp.EffectIndeterminateDP:
		return pb.Response_INDETERMINATEDP, nil
	}

	return pb.Response_INDETERMINATE, newUnknownEffectError(effect)
}

func makeFailEffect(effect pb.Response_Effect) (pb.Response_Effect, error) {
	switch effect {
	case pb.Response_DENY:
		return pb.Response_INDETERMINATED, nil

	case pb.Response_PERMIT:
		return pb.Response_INDETERMINATEP, nil

	case pb.Response_NOTAPPLICABLE, pb.Response_INDETERMINATE, pb.Response_INDETERMINATED, pb.Response_INDETERMINATEP, pb.Response_INDETERMINATEDP:
		return effect, nil
	}

	return pb.Response_INDETERMINATE, newUnknownEffectError(int(effect))
}

type obligation []*pb.Attribute

func (o obligation) String() string {
	if len(o) <= 0 {
		return "no attributes"
	}

	lines := []string{"attributes:"}
	for _, attr := range o {
		lines = append(lines, fmt.Sprintf("- %s.(%s): %q", attr.Id, attr.Type, attr.Value))
	}

	return strings.Join(lines, "\n")
}

func (s *Server) newContext(c *pdp.LocalContentStorage, in *pb.Request) (*pdp.Context, error) {
	ctx, err := pdp.NewContext(c, len(in.Attributes), func(i int) (string, pdp.AttributeValue, error) {
		a := in.Attributes[i]

		t, ok := pdp.BuiltinTypes[strings.ToLower(a.Type)]
		if !ok {
			return "", pdp.UndefinedValue, bindError(newUnknownAttributeTypeError(a.Type), a.Id)
		}

		v, err := pdp.MakeValueFromString(t, a.Value)
		if err != nil {
			return "", pdp.UndefinedValue, bindError(err, a.Id)
		}

		return a.Id, v, nil
	})
	if err != nil {
		return nil, newContextCreationError(err)
	}

	return ctx, nil
}

func (s *Server) newAttributes(obligations []pdp.AttributeAssignmentExpression, ctx *pdp.Context) ([]*pb.Attribute, error) {
	attrs := make([]*pb.Attribute, len(obligations))
	for i, e := range obligations {
		ID, t, s, err := e.Serialize(ctx)
		if err != nil {
			return attrs[:i], err
		}

		attrs[i] = &pb.Attribute{
			Id:    ID,
			Type:  t,
			Value: s}
	}

	return attrs, nil
}

func (s *Server) rawValidate(p *pdp.PolicyStorage, c *pdp.LocalContentStorage, in *pb.Request) (pb.Response_Effect, []error, []*pb.Attribute) {
	if p == nil {
		return pb.Response_INDETERMINATE, []error{newMissingPolicyError()}, nil
	}

	ctx, err := s.newContext(c, in)
	if err != nil {
		return pb.Response_INDETERMINATE, []error{err}, nil
	}

	if s.opts.logger.Level >= log.DebugLevel {
		s.opts.logger.WithField("context", ctx).Debug("Request context")
	}

	errs := []error{}

	r := p.Root().Calculate(ctx)
	effect, obligations, err := r.Status()
	if err != nil {
		errs = append(errs, newPolicyCalculationError(err))
	}

	re, err := makeEffect(effect)
	if err != nil {
		errs = append(errs, newEffectTranslationError(err))
	}

	if len(errs) > 0 {
		re, err = makeFailEffect(re)
		if err != nil {
			errs = append(errs, newEffectCombiningError(err))
		}
	}

	attrs, err := s.newAttributes(obligations, ctx)
	if err != nil {
		errs = append(errs, newObligationTranslationError(err))
	}

	return re, errs, attrs
}

// Validate is a server handler for gRPC call
// It handles PDP decision requests
func (s *Server) Validate(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Validate() called: in='%v'", in)
	s.RLock()
	p := s.p
	c := s.c
	s.RUnlock()

	effect, errs, attrs := s.rawValidate(p, c, in)

	status := "Ok"
	if len(errs) > 1 {
		status = newMultiError(errs).Error()
	} else if len(errs) > 0 {
		status = errs[0].Error()
	}

	if s.opts.logger.Level >= log.DebugLevel {
		s.opts.logger.WithFields(log.Fields{
			"effect":     pb.Response_Effect_name[int32(effect)],
			"reason":     status,
			"obligation": obligation(attrs),
		}).Debug("Response")
	}

	return &pb.Response{
		Effect:     effect,
		Reason:     status,
		Obligation: attrs}, nil
}

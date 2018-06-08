package jast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func makeSource(desc string, id string, hidden bool) string {
	if hidden {
		return fmt.Sprintf("hidden %s", desc)
	}
	return fmt.Sprintf("%s \"%s\"", desc, id)
}

func (ctx *context) unmarshalPolicies(d *json.Decoder) ([]pdp.Evaluable, error) {
	if err := jparser.CheckArrayStart(d, "policy set"); err != nil {
		return nil, err
	}

	policies := []pdp.Evaluable{}
	if err := jparser.UnmarshalObjectArray(d, func(idx int, d *json.Decoder) error {
		e, err := ctx.unmarshalEvaluable(d)
		if err != nil {
			return bindErrorf(err, "%d", idx)
		}

		policies = append(policies, e)

		return nil
	}, "policy set"); err != nil {
		return nil, err
	}

	return policies, nil
}

func (ctx *context) unmarshalEvaluable(d *json.Decoder) (pdp.Evaluable, error) {
	var (
		hidden      = true
		isPolicy    bool
		isPolicySet bool

		pid      string
		policies []pdp.Evaluable
		rules    []*pdp.Rule
		target   pdp.Target
		obligs   []pdp.AttributeAssignmentExpression
		alg      interface{}
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var err error

		switch strings.ToLower(k) {
		case yastTagID:
			hidden = false
			pid, err = jparser.GetString(d, "policy or policy set id")
			return err

		case yastTagAlg:
			alg, err = ctx.unmarshalCombiningAlg(d)
			if err != nil {
				return bindError(err, makeSource("policy or policy set", pid, hidden))
			}
			return err

		case yastTagTarget:
			target, err = ctx.unmarshalTarget(d)
			if err != nil {
				return bindError(err, makeSource("policy or policy set", pid, hidden))
			}
			return nil

		case yastTagObligation:
			obligs, err = ctx.unmarshalObligations(d)
			if err != nil {
				return bindError(err, makeSource("policy or policy set", pid, hidden))
			}
			return nil

		case yastTagPolicies:
			isPolicySet = true
			policies, err = ctx.unmarshalPolicies(d)
			if err != nil {
				return bindError(err, makeSource("policy set", pid, hidden))
			}
			return nil

		case yastTagRules:
			isPolicy = true
			rules, err = ctx.unmarshalRules(d)
			if err != nil {
				return bindError(err, makeSource("policy", pid, hidden))
			}
			return nil
		}

		return newUnknownFieldError(k)
	}, "policy or policy set"); err != nil {
		return nil, err
	}

	if isPolicy && isPolicySet {
		return nil, newPolicyAmbiguityError()
	}

	if isPolicySet {
		src := makeSource("policy set", pid, hidden)

		if alg == nil {
			return nil, bindError(newMissingPCAError(), src)
		}

		maker, params, err := ctx.buildPolicyCombiningAlg(alg, policies)
		if err != nil {
			return nil, bindError(err, src)
		}

		return pdp.NewPolicySet(pid, hidden, target, policies, maker, params, obligs), nil
	}

	if isPolicy {
		src := makeSource("policy", pid, hidden)

		if alg == nil {
			return nil, bindError(newMissingRCAError(), src)
		}

		maker, params, err := ctx.buildRuleCombiningAlg(alg, rules)
		if err != nil {
			return nil, bindError(err, src)
		}

		return pdp.NewPolicy(pid, hidden, target, rules, maker, params, obligs), nil
	}

	return nil, newPolicyMissingKeyError()
}

func (ctx *context) unmarshalRootPolicy(d *json.Decoder) error {
	if err := jparser.CheckObjectStart(d, "root policy or policy set"); err != nil {
		return err
	}

	e, err := ctx.unmarshalEvaluable(d)
	if err != nil {
		return err
	}

	ctx.rootPolicy = e
	return nil
}

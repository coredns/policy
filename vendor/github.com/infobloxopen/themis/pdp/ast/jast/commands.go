package jast

import (
	"encoding/json"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx *context) unmarshalEntity(d *json.Decoder) (interface{}, error) {
	if err := jparser.CheckObjectStart(d, "entity"); err != nil {
		return nil, err
	}

	var (
		hidden      = true
		isPolicy    bool
		isPolicySet bool
		isRule      bool

		id       string
		effect   = -1
		policies []pdp.Evaluable
		rules    []*pdp.Rule
		target   pdp.Target
		cond     pdp.Expression
		obligs   []pdp.AttributeAssignmentExpression
		alg      interface{}
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var err error

		switch strings.ToLower(k) {
		case yastTagID:
			hidden = false
			id, err = jparser.GetString(d, "policy or set or rule id")
			return err

		case yastTagAlg:
			alg, err = ctx.unmarshalCombiningAlg(d)
			return err

		case yastTagTarget:
			target, err = ctx.unmarshalTarget(d)
			return err

		case yastTagObligation:
			obligs, err = ctx.unmarshalObligations(d)
			return err

		case yastTagPolicies:
			isPolicySet = true
			policies, err = ctx.unmarshalPolicies(d)
			if err != nil {
				return bindError(err, makeSource("policy set", id, hidden))
			}
			return nil

		case yastTagRules:
			isPolicy = true
			rules, err = ctx.unmarshalRules(d)
			if err != nil {
				return bindError(err, makeSource("policy", id, hidden))
			}
			return nil

		case yastTagEffect:
			isRule = true
			var s string
			src := makeSource("rule", id, hidden)
			s, err = jparser.GetString(d, "effect")
			if err != nil {
				return bindError(err, src)
			}

			var ok bool
			effect, ok = pdp.EffectIDs[strings.ToLower(s)]
			if !ok {
				return bindError(newUnknownEffectError(s), src)
			}
			return nil

		case yastTagCondition:
			cond, err = ctx.unmarshalCondition(d)
			return err
		}

		return newUnknownFieldError(k)
	}, "entity"); err != nil {
		return nil, err
	}

	if isRule && isPolicy || isRule && isPolicySet || isPolicy && isPolicySet {
		tags := []string{}
		if isPolicy {
			tags = append(tags, yastTagRules)
		}

		if isPolicySet {
			tags = append(tags, yastTagPolicies)
		}

		if isRule {
			tags = append(tags, yastTagEffect)
		}

		return nil, newEntityAmbiguityError(tags)
	}

	if isPolicySet {
		src := makeSource("policy set", id, hidden)

		if alg == nil {
			return nil, bindError(newMissingPCAError(), src)
		}

		maker, params, err := ctx.buildPolicyCombiningAlg(alg, policies)
		if err != nil {
			return nil, bindError(err, src)
		}

		return pdp.NewPolicySet(id, hidden, target, policies, maker, params, obligs), nil
	}

	if isPolicy {
		src := makeSource("policy", id, hidden)

		if alg == nil {
			return nil, bindError(newMissingRCAError(), src)
		}

		maker, params, err := ctx.buildRuleCombiningAlg(alg, rules)
		if err != nil {
			return nil, bindError(err, src)
		}

		return pdp.NewPolicy(id, hidden, target, rules, maker, params, obligs), nil
	}

	if isRule {
		if effect == -1 {
			return nil, bindError(newMissingAttributeError(yastTagEffect, "rule"), makeSource("rule", id, hidden))
		}

		return pdp.NewRule(id, hidden, target, cond, effect, obligs), nil
	}

	return nil, newEntityMissingKeyError()
}

func (ctx *context) unmarshalCommand(d *json.Decoder, u *pdp.PolicyUpdate) error {
	var (
		op     int
		path   []string
		entity interface{}
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var err error

		switch strings.ToLower(k) {
		case yastTagOp:
			var s string
			s, err = jparser.GetString(d, "operation")
			if err != nil {
				return err
			}

			var ok bool
			op, ok = pdp.UpdateOpIDs[strings.ToLower(s)]
			if !ok {
				return newUnknownPolicyUpdateOperationError(s)
			}

			return nil

		case yastTagPath:
			path = []string{}
			err = jparser.GetStringSequence(d, func(idx int, s string) error {
				path = append(path, s)
				return nil
			}, "path")

			return err

		case yastTagEntity:
			if op == pdp.UOAdd {
				entity, err = ctx.unmarshalEntity(d)
			}

			return err
		}

		return newUnknownFieldError(k)
	}, "command"); err != nil {
		return err
	}

	u.Append(op, path, entity)

	return nil
}

func (ctx *context) unmarshalCommands(d *json.Decoder, u *pdp.PolicyUpdate) error {
	if err := jparser.CheckArrayStart(d, "commands"); err != nil {
		return err
	}

	return jparser.UnmarshalObjectArray(d, func(idx int, d *json.Decoder) error {
		if err := ctx.unmarshalCommand(d, u); err != nil {
			return bindErrorf(err, "%d", idx)
		}

		return nil
	}, "commands")
}

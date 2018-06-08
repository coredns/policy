package jast

import (
	"encoding/json"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx context) unmarshalCondition(d *json.Decoder) (pdp.Expression, error) {
	if err := jparser.CheckObjectStart(d, "condition"); err != nil {
		return nil, err
	}

	e, err := ctx.unmarshalExpression(d)
	if err != nil {
		return nil, err
	}

	t := e.GetResultType()
	if t != pdp.TypeBoolean {
		return nil, newConditionTypeError(t)
	}

	return e, nil
}

func (ctx context) unmarshalRule(d *json.Decoder) (*pdp.Rule, error) {
	var (
		hidden = true
		id     string
		effect = -1
		target pdp.Target
		cond   pdp.Expression
		obligs []pdp.AttributeAssignmentExpression
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var err error

		switch strings.ToLower(k) {
		case yastTagID:
			hidden = false
			id, err = jparser.GetString(d, "rule id")
			return err

		case yastTagTarget:
			target, err = ctx.unmarshalTarget(d)
			return err

		case yastTagObligation:
			obligs, err = ctx.unmarshalObligations(d)
			return err

		case yastTagEffect:
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
	}, "rule"); err != nil {
		return nil, err
	}

	if effect == -1 {
		return nil, bindError(newMissingAttributeError(yastTagEffect, "rule"), makeSource("rule", id, hidden))
	}

	return pdp.NewRule(id, hidden, target, cond, effect, obligs), nil
}

func (ctx *context) unmarshalRules(d *json.Decoder) ([]*pdp.Rule, error) {
	err := jparser.CheckArrayStart(d, "rules")
	if err != nil {
		return nil, err
	}

	rules := []*pdp.Rule{}
	if err = jparser.UnmarshalObjectArray(d, func(idx int, d *json.Decoder) error {
		e, err := ctx.unmarshalRule(d)
		if err != nil {
			return bindErrorf(err, "%d", idx)
		}

		rules = append(rules, e)

		return nil
	}, "rules"); err != nil {
		return nil, err
	}

	return rules, nil
}

package jast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

type (
	policyCombiningAlgParamBuilder func(ctx context, alg *caParams, policies []pdp.Evaluable) (interface{}, boundError)
	ruleCombiningAlgParamBuilder   func(ctx context, alg *caParams, rules []*pdp.Rule) (interface{}, boundError)
)

var (
	policyCombiningAlgParamBuilders = map[string]policyCombiningAlgParamBuilder{}
	ruleCombiningAlgParamBuilders   = map[string]ruleCombiningAlgParamBuilder{}
)

func init() {
	policyCombiningAlgParamBuilders["mapper"] = buildMapperPolicyCombiningAlgParams
	ruleCombiningAlgParamBuilders["mapper"] = buildMapperRuleCombiningAlgParams
}

type caParams struct {
	id     string
	defOk  bool
	defID  string
	errOk  bool
	errID  string
	arg    pdp.Expression
	order  int
	subAlg interface{}
}

func checkPolicyID(ID string, policies []pdp.Evaluable) bool {
	for _, p := range policies {
		if pid, ok := p.GetID(); ok && ID == pid {
			return true
		}
	}

	return false
}

func buildMapperPolicyCombiningAlgParams(ctx context, alg *caParams, policies []pdp.Evaluable) (interface{}, boundError) {
	if alg.defOk {
		if !checkPolicyID(alg.defID, policies) {
			return nil, newMissingDefaultPolicyPCAError(alg.defID)
		}
	}

	if alg.errOk {
		if !checkPolicyID(alg.errID, policies) {
			return nil, newMissingErrorPolicyPCAError(alg.errID)
		}
	}

	var subAlg pdp.PolicyCombiningAlg
	t := alg.arg.GetResultType()
	if _, ok := t.(*pdp.FlagsType); ok || t == pdp.TypeSetOfStrings || t == pdp.TypeListOfStrings {
		if alg.subAlg != nil {
			maker, params, err := ctx.buildPolicyCombiningAlg(alg.subAlg, policies)
			if err != nil {
				return nil, err
			}

			subAlg = maker(nil, params)
		}

		if subAlg == nil {
			return nil, newMissingPCAError()
		}
	}

	return pdp.MapperPCAParams{
		Argument:  alg.arg,
		DefOk:     alg.defOk,
		Def:       alg.defID,
		ErrOk:     alg.errOk,
		Err:       alg.errID,
		Order:     alg.order,
		Algorithm: subAlg}, nil
}

func checkRuleID(ID string, rules []*pdp.Rule) bool {
	for _, r := range rules {
		if rid, ok := r.GetID(); ok && ID == rid {
			return true
		}
	}

	return false
}

func buildMapperRuleCombiningAlgParams(ctx context, alg *caParams, rules []*pdp.Rule) (interface{}, boundError) {
	if alg.defOk {
		if !checkRuleID(alg.defID, rules) {
			return nil, newMissingDefaultRuleRCAError(alg.defID)
		}
	}

	if alg.errOk {
		if !checkRuleID(alg.errID, rules) {
			return nil, newMissingErrorRuleRCAError(alg.errID)
		}
	}

	var subAlg pdp.RuleCombiningAlg
	t := alg.arg.GetResultType()
	if _, ok := t.(*pdp.FlagsType); ok || t == pdp.TypeSetOfStrings || t == pdp.TypeListOfStrings {
		if alg.subAlg != nil {
			maker, params, err := ctx.buildRuleCombiningAlg(alg.subAlg, rules)
			if err != nil {
				return nil, err
			}

			subAlg = maker(nil, params)
		}

		if subAlg == nil {
			return nil, newMissingRCAError()
		}
	}

	return pdp.MapperRCAParams{
		Argument:  alg.arg,
		DefOk:     alg.defOk,
		Def:       alg.defID,
		ErrOk:     alg.errOk,
		Err:       alg.errID,
		Order:     alg.order,
		Algorithm: subAlg}, nil
}

func (ctx context) buildRuleCombiningAlg(alg interface{}, rules []*pdp.Rule) (pdp.RuleCombiningAlgMaker, interface{}, boundError) {
	switch alg := alg.(type) {
	case *caParams:
		id := strings.ToLower(alg.id)
		maker, ok := pdp.RuleCombiningParamAlgs[id]
		if !ok {
			return nil, nil, newUnknownRCAError(alg.id)
		}

		paramBuilder, ok := ruleCombiningAlgParamBuilders[id]
		if !ok {
			return nil, nil, newNotImplementedRCAError(alg.id)
		}

		params, err := paramBuilder(ctx, alg, rules)
		if err != nil {
			return nil, nil, bindError(err, alg.id)
		}

		return maker, params, nil
	case string:
		maker, ok := pdp.RuleCombiningAlgs[strings.ToLower(alg)]
		if !ok {
			return nil, nil, newUnknownRCAError(alg)
		}

		return maker, nil, nil
	}

	return nil, nil, newInvalidRCAError(alg)
}

func (ctx context) buildPolicyCombiningAlg(alg interface{}, policies []pdp.Evaluable) (pdp.PolicyCombiningAlgMaker, interface{}, boundError) {
	switch alg := alg.(type) {
	case *caParams:
		id := strings.ToLower(alg.id)
		maker, ok := pdp.PolicyCombiningParamAlgs[id]
		if !ok {
			return nil, nil, newUnknownPCAError(alg.id)
		}

		paramBuilder, ok := policyCombiningAlgParamBuilders[id]
		if !ok {
			return nil, nil, newNotImplementedPCAError(alg.id)
		}

		params, err := paramBuilder(ctx, alg, policies)
		if err != nil {
			return nil, nil, bindError(err, alg.id)
		}

		return maker, params, nil
	case string:
		maker, ok := pdp.PolicyCombiningAlgs[strings.ToLower(alg)]
		if !ok {
			return nil, nil, newUnknownPCAError(alg)
		}

		return maker, nil, nil
	}

	return nil, nil, newInvalidPCAError(alg)
}

func (ctx context) unmarshalCombiningAlgObj(d *json.Decoder) (*caParams, error) {
	var (
		params caParams
		idOk   bool
		mapOk  bool
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var (
			err error
			ok  bool
		)

		switch strings.ToLower(k) {
		case yastTagID:
			params.id, err = jparser.GetString(d, "id")
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			idOk = true
			return nil

		case yastTagMap:
			err = jparser.CheckObjectStart(d, "expression")
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			params.arg, err = ctx.unmarshalExpression(d)
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			t := params.arg.GetResultType()
			if _, ok := t.(*pdp.FlagsType); !ok && t != pdp.TypeString && t != pdp.TypeSetOfStrings && t != pdp.TypeListOfStrings {
				err = newMapperArgumentTypeError(t)
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			mapOk = true
			return nil

		case yastTagDefault:
			params.defID, err = jparser.GetString(d, "algorithm default id")
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			params.defOk = true
			return nil

		case yastTagError:
			params.errID, err = jparser.GetString(d, "algorithm error id")
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			params.errOk = true
			return nil

		case yastTagOrder:
			s, err := jparser.GetString(d, "ordering option")
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			params.order, ok = pdp.MapperPCAOrderIDs[strings.ToLower(s)]
			if !ok {
				err = newUnknownMapperCAOrder(s)
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			return nil

		case yastTagAlg:
			params.subAlg, err = ctx.unmarshalCombiningAlg(d)
			if err != nil {
				if idOk {
					err = bindErrorf(err, "%q", params.id)
				}

				return err
			}

			return nil
		}

		return newUnknownFieldError(k)
	}, "algorithm"); err != nil {
		return nil, err
	}

	if !idOk {
		return nil, newMissingAttributeError(yastTagID, "algorithm")
	}

	if !mapOk {
		if idOk {
			return nil, bindError(newMissingAttributeError(yastTagMap, fmt.Sprintf("%q", params.id)), "algorithm")
		}

		return nil, newMissingAttributeError(yastTagMap, "algorithm")
	}

	return &params, nil
}

func (ctx context) unmarshalCombiningAlg(d *json.Decoder) (interface{}, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}

	switch t := t.(type) {
	case json.Delim:
		if t.String() == jparser.DelimObjectStart {
			return ctx.unmarshalCombiningAlgObj(d)
		}

		return nil, newParseCAError(t)
	case string:
		return t, nil
	default:
		return nil, newParseCAError(t)
	}
}

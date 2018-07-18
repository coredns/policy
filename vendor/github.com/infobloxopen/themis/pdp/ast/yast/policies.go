package yast

import (
	"fmt"
	"strings"

	"github.com/infobloxopen/themis/pdp"
)

type (
	policyCombiningAlgParamUnmarshaler func(ctx context, m map[interface{}]interface{}, policies []pdp.Evaluable) (interface{}, boundError)
	ruleCombiningAlgParamUnmarshaler   func(ctx context, m map[interface{}]interface{}, rules []*pdp.Rule) (interface{}, boundError)
)

var (
	policyCombiningAlgParamUnmarshalers = map[string]policyCombiningAlgParamUnmarshaler{}
	ruleCombiningAlgParamUnmarshalers   = map[string]ruleCombiningAlgParamUnmarshaler{}
)

func init() {
	policyCombiningAlgParamUnmarshalers["mapper"] = unmarshalMapperPolicyCombiningAlgParams
	ruleCombiningAlgParamUnmarshalers["mapper"] = unmarshalMapperRuleCombiningAlgParams
}

func checkRuleID(ID string, rules []*pdp.Rule) bool {
	for _, r := range rules {
		if rid, ok := r.GetID(); ok && ID == rid {
			return true
		}
	}

	return false
}

func unmarshalMapperRuleCombiningAlgParams(ctx context, m map[interface{}]interface{}, rules []*pdp.Rule) (interface{}, boundError) {
	v, ok := m[yastTagMap]
	if !ok {
		return nil, newMissingMapRCAParamError()
	}

	arg, err := ctx.unmarshalExpression(v)
	if err != nil {
		return nil, err
	}

	t := arg.GetResultType()
	_, flagsOk := t.(*pdp.FlagsType)
	if !flagsOk && t != pdp.TypeString && t != pdp.TypeSetOfStrings && t != pdp.TypeListOfStrings {
		return nil, newMapperArgumentTypeError(t)
	}

	defID, defOk, err := ctx.extractStringOpt(m, yastTagDefault, "default rule id")
	if err != nil {
		return nil, err
	}

	if defOk {
		if !checkRuleID(defID, rules) {
			return nil, newMissingDefaultRuleRCAError(defID)
		}
	}

	errID, errOk, err := ctx.extractStringOpt(m, yastTagError, "on error rule id")
	if err != nil {
		return nil, err
	}

	if errOk {
		if !checkRuleID(errID, rules) {
			return nil, newMissingErrorRuleRCAError(errID)
		}
	}

	var subAlg pdp.RuleCombiningAlg
	order := pdp.MapperPCAExternalOrder
	if flagsOk || t == pdp.TypeSetOfStrings || t == pdp.TypeListOfStrings {
		s, ok, err := ctx.extractStringOpt(m, yastTagOrder, "ordering option")
		if err != nil {
			return nil, err
		}

		if ok {
			order, ok = pdp.MapperRCAOrderIDs[strings.ToLower(s)]
			if !ok {
				return nil, newUnknownMapperRCAOrder(s)
			}
		}

		maker, params, err := ctx.unmarshalRuleCombiningAlg(m, nil)
		if err != nil {
			return nil, err
		}
		subAlg = maker(nil, params)
	}

	return pdp.MapperRCAParams{
		Argument:  arg,
		DefOk:     defOk,
		Def:       defID,
		ErrOk:     errOk,
		Err:       errID,
		Order:     order,
		Algorithm: subAlg}, nil
}

func (ctx context) unmarshalRuleCombiningAlgObj(m map[interface{}]interface{}, rules []*pdp.Rule) (pdp.RuleCombiningAlgMaker, interface{}, boundError) {
	ID, err := ctx.extractString(m, yastTagID, "algorithm id")
	if err != nil {
		return nil, nil, err
	}

	s := strings.ToLower(ID)
	maker, ok := pdp.RuleCombiningParamAlgs[s]
	if !ok {
		return nil, nil, newUnknownRCAError(ID)
	}

	paramUnmarshaler, ok := ruleCombiningAlgParamUnmarshalers[s]
	if !ok {
		return nil, nil, newNotImplementedRCAError(ID)
	}

	params, err := paramUnmarshaler(ctx, m, rules)
	if err != nil {
		return nil, nil, bindError(err, ID)
	}

	return maker, params, nil
}

func (ctx context) unmarshalRuleCombiningAlg(m map[interface{}]interface{}, rules []*pdp.Rule) (pdp.RuleCombiningAlgMaker, interface{}, boundError) {
	v, ok := m[yastTagAlg]
	if !ok {
		return nil, nil, newMissingRCAError()
	}

	switch alg := v.(type) {
	case string:
		maker, ok := pdp.RuleCombiningAlgs[strings.ToLower(alg)]
		if !ok {
			return nil, nil, newUnknownRCAError(alg)
		}

		return maker, nil, nil

	case map[interface{}]interface{}:
		return ctx.unmarshalRuleCombiningAlgObj(alg, rules)
	}

	return nil, nil, newInvalidRCAError(v)
}

func (ctx context) unmarshalPolicy(m map[interface{}]interface{}, i int, ID string, hidden bool, rules interface{}) (pdp.Evaluable, boundError) {
	src := makeSource("policy", ID, hidden, i)

	target, err := ctx.unmarshalTarget(m)
	if err != nil {
		return nil, bindError(err, src)
	}

	rls, err := ctx.unmarshalRules(rules)
	if err != nil {
		return nil, bindError(err, src)
	}

	alg, params, err := ctx.unmarshalRuleCombiningAlg(m, rls)
	if err != nil {
		return nil, bindError(err, src)
	}

	obls, err := ctx.unmarshalObligations(m)
	if err != nil {
		return nil, bindError(err, src)
	}

	return pdp.NewPolicy(ID, hidden, target, rls, alg, params, obls), nil
}

func checkPolicyID(ID string, policies []pdp.Evaluable) bool {
	for _, p := range policies {
		if pid, ok := p.GetID(); ok && ID == pid {
			return true
		}
	}

	return false
}

func unmarshalMapperPolicyCombiningAlgParams(ctx context, m map[interface{}]interface{}, policies []pdp.Evaluable) (interface{}, boundError) {
	v, ok := m[yastTagMap]
	if !ok {
		return nil, newMissingMapPCAParamError()
	}

	arg, err := ctx.unmarshalExpression(v)
	if err != nil {
		return nil, err
	}

	t := arg.GetResultType()
	_, flagsOk := t.(*pdp.FlagsType)
	if !flagsOk && t != pdp.TypeString && t != pdp.TypeSetOfStrings && t != pdp.TypeListOfStrings {
		return nil, newMapperArgumentTypeError(t)
	}

	defID, defOk, err := ctx.extractStringOpt(m, yastTagDefault, "default policy id")
	if err != nil {
		return nil, err
	}

	if defOk {
		if !checkPolicyID(defID, policies) {
			return nil, newMissingDefaultPolicyPCAError(defID)
		}
	}

	errID, errOk, err := ctx.extractStringOpt(m, yastTagError, "on error policy id")
	if err != nil {
		return nil, err
	}

	if errOk {
		if !checkPolicyID(errID, policies) {
			return nil, newMissingErrorPolicyPCAError(errID)
		}
	}

	var subAlg pdp.PolicyCombiningAlg
	order := pdp.MapperPCAExternalOrder
	if flagsOk || t == pdp.TypeSetOfStrings || t == pdp.TypeListOfStrings {
		s, ok, err := ctx.extractStringOpt(m, yastTagOrder, "ordering option")
		if err != nil {
			return nil, err
		}

		if ok {
			order, ok = pdp.MapperPCAOrderIDs[strings.ToLower(s)]
			if !ok {
				return nil, newUnknownMapperPCAOrder(s)
			}
		}

		maker, params, err := ctx.unmarshalPolicyCombiningAlg(m, nil)
		if err != nil {
			return nil, err
		}
		subAlg = maker(nil, params)
	}

	return pdp.MapperPCAParams{
		Argument:  arg,
		DefOk:     defOk,
		Def:       defID,
		ErrOk:     errOk,
		Err:       errID,
		Order:     order,
		Algorithm: subAlg}, nil
}

func (ctx context) unmarshalPolicyCombiningAlgObj(m map[interface{}]interface{}, policies []pdp.Evaluable) (pdp.PolicyCombiningAlgMaker, interface{}, boundError) {
	ID, err := ctx.extractString(m, yastTagID, "algorithm id")
	if err != nil {
		return nil, nil, err
	}

	s := strings.ToLower(ID)
	maker, ok := pdp.PolicyCombiningParamAlgs[s]
	if !ok {
		return nil, nil, newUnknownPCAError(ID)
	}

	paramUnmarshaler, ok := policyCombiningAlgParamUnmarshalers[s]
	if !ok {
		return nil, nil, newNotImplementedPCAError(ID)
	}

	params, err := paramUnmarshaler(ctx, m, policies)
	if err != nil {
		return nil, nil, bindError(err, ID)
	}

	return maker, params, nil
}

func (ctx context) unmarshalPolicyCombiningAlg(m map[interface{}]interface{}, policies []pdp.Evaluable) (pdp.PolicyCombiningAlgMaker, interface{}, boundError) {
	v, ok := m[yastTagAlg]
	if !ok {
		return nil, nil, newMissingPCAError()
	}

	switch alg := v.(type) {
	case string:
		maker, ok := pdp.PolicyCombiningAlgs[strings.ToLower(alg)]
		if !ok {
			return nil, nil, newUnknownPCAError(alg)
		}

		return maker, nil, nil

	case map[interface{}]interface{}:
		return ctx.unmarshalPolicyCombiningAlgObj(alg, policies)
	}

	return nil, nil, newInvalidPCAError(v)
}

func (ctx context) unmarshalPolicies(v interface{}) ([]pdp.Evaluable, boundError) {
	items, err := ctx.validateList(v, "list of policies")
	if err != nil {
		return nil, err
	}

	pols := []pdp.Evaluable{}
	for i, v := range items {
		p, err := ctx.unmarshalItem(v, i+1)
		if err != nil {
			return nil, err
		}

		pols = append(pols, p)
	}

	return pols, nil
}

func (ctx context) unmarshalPolicySet(m map[interface{}]interface{}, i int, ID string, hidden bool, policies interface{}) (pdp.Evaluable, boundError) {
	src := makeSource("policy set", ID, hidden, i)

	target, err := ctx.unmarshalTarget(m)
	if err != nil {
		return nil, bindError(err, src)
	}

	pols, err := ctx.unmarshalPolicies(policies)
	if err != nil {
		return nil, bindError(err, src)
	}

	alg, params, err := ctx.unmarshalPolicyCombiningAlg(m, pols)
	if err != nil {
		return nil, bindError(err, src)
	}

	obls, err := ctx.unmarshalObligations(m)
	if err != nil {
		return nil, bindError(err, src)
	}

	return pdp.NewPolicySet(ID, hidden, target, pols, alg, params, obls), nil
}

func makeSource(desc string, ID string, hidden bool, idx int) string {
	src := fmt.Sprintf("hidden %s", desc)
	if !hidden {
		src = fmt.Sprintf("%s \"%s\"", desc, ID)
	}

	if idx > 0 {
		src = fmt.Sprintf("(%d) %s", idx, src)
	}

	return src
}

func (ctx context) unmarshalItem(v interface{}, i int) (pdp.Evaluable, boundError) {
	m, err := ctx.validateMap(v, "policy or policy set")
	if err != nil {
		if i > 0 {
			err = bindErrorf(err, "%d", i)
		}

		return nil, err
	}

	ID, ok, err := ctx.extractStringOpt(m, yastTagID, "policy or policy set id")
	if err != nil {
		if i > 0 {
			err = bindErrorf(err, "%d", i)
		}

		return nil, err
	}

	src := makeSource("policy or policy set", ID, !ok, i)

	rules, rOk := m[yastTagRules]
	policies, pOk := m[yastTagPolicies]

	if rOk && pOk {
		return nil, bindError(newPolicyAmbiguityError(), src)
	}

	if rOk {
		return ctx.unmarshalPolicy(m, i, ID, !ok, rules)
	}

	if pOk {
		return ctx.unmarshalPolicySet(m, i, ID, !ok, policies)
	}

	return nil, bindError(newPolicyMissingKeyError(), src)
}

func (ctx context) unmarshalRootPolicy(m map[interface{}]interface{}) (pdp.Evaluable, boundError) {
	m, ok, err := ctx.extractMapOpt(m, yastTagPolicies, "root policy or policy set")
	if !ok || err != nil {
		return nil, err
	}

	return ctx.unmarshalItem(m, 0)
}

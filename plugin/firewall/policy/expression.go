package policy

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/pkg/rqdata"

	expr "github.com/Knetic/govaluate"
)

type ruleExpr struct {
	action        int
	actionIfError int
	expression    *expr.EvaluableExpression
}

// ExprEngine implement interface Engine for Firewall plugin
// it evaluate the rues using an the lib Knetic/govaluate
type ExprEngine struct {
	actionIfErrorEvaluation int
	dataFromReq             *rqdata.Mapping
}

type dataAsParam struct {
	ctx         context.Context
	dataFromReq *rqdata.Extractor
}

// NewExprEngine create a new Engine with default configuration
func NewExprEngine() *ExprEngine {
	return &ExprEngine{TypeRefuse, rqdata.NewMapping("")}
}

//BuildQueryData here return a dataAsParam that can be used by to evaluate the variables of the expression
func (x *ExprEngine) BuildQueryData(ctx context.Context, state request.Request) (interface{}, error) {
	return &dataAsParam{ctx, rqdata.NewExtractor(state, x.dataFromReq)}, nil
}

//BuildReplyData here return a dataAsParam that can be used by to evaluate the variables of the expression
func (x *ExprEngine) BuildReplyData(ctx context.Context, state request.Request, query interface{}) (interface{}, error) {
	return &dataAsParam{ctx, rqdata.NewExtractor(state, x.dataFromReq)}, nil
}

//BuildRule create a rule for Expression Engine:
// - first param is one of the action to return
// - second and following param is a sentence the represent an Expression
func (x *ExprEngine) BuildRule(args []string) (Rule, error) {
	keyword := args[0]
	exp := args[1:]
	e, err := expr.NewEvaluableExpressionWithFunctions(strings.Join(exp, " "), map[string]expr.ExpressionFunction{
		"atoi":   atoi,
		"incidr": incidr,
		"random": random,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create a valid expression : %s", err)
	}

	var kind = TypeNone
	for k, n := range NameTypes {
		if keyword == n {
			kind = k
		}
	}
	if kind == TypeNone {
		return nil, fmt.Errorf("invalid keyword %s for a policy rule", keyword)
	}
	return &ruleExpr{kind, x.actionIfErrorEvaluation, e}, nil
}

func random(arguments ...interface{}) (interface{}, error) {
	if len(arguments) != 0 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	return rand.Float64(), nil
}

func atoi(arguments ...interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("atoi requires exactly one string argument")
	}
	s, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("atoi requires exactly one string argument")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return float64(n), nil
}

func incidr(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments")
	}
	ip := net.ParseIP(args[0].(string))
	if ip == nil {
		return nil, fmt.Errorf("first argument is not an IP address")
	}
	_, cidr, err := net.ParseCIDR(args[1].(string))
	if err != nil {
		return nil, err
	}
	return cidr.Contains(ip), nil
}

func toBoolean(v interface{}) (bool, error) {
	if s, ok := v.(string); ok {
		return strings.ToLower(s) == "true", nil
	}
	if b, ok := v.(bool); ok {
		return b, nil
	}
	if i, ok := v.(int); ok {
		return i != 0, nil
	}
	return false, fmt.Errorf("cannot extract a boolean value from result of expression")
}

//Evaluate the current expression, using data as a variable resolver for Expression
func (r *ruleExpr) Evaluate(data interface{}) (int, error) {

	params, ok := data.(*dataAsParam)
	if !ok {
		return r.actionIfError, fmt.Errorf("evaluation of expression '%s' - params provided are of wrong type, expect a go Context", r.expression.String())
	}
	res, err := r.expression.Eval(params)
	if err != nil {
		return r.actionIfError, fmt.Errorf("evaluation of expression '%s' return an error : %s", r.expression.String(), err)
	}
	result, err := toBoolean(res)
	if err != nil {
		return r.actionIfError, fmt.Errorf("evaluation of expression '%s' return an non boolean value : %s", r.expression.String(), err)
	}

	if result {
		return r.action, nil
	}
	return TypeNone, nil
}

// Get return the value associated with the variable
// required by the interface of Knetic/govaluate for evaluation of the 'variables' in the expression
// DataRequestExtractor is evaluated first, and if the name does not match then metadata is evaluated
func (p *dataAsParam) Get(name string) (interface{}, error) {
	v, exist := p.dataFromReq.Value(name)
	if exist {
		return v, nil
	}
	f := metadata.ValueFunc(p.ctx, name)
	if f == nil {
		return "", nil
	}
	return f(), nil
}

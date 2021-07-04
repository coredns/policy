package casbin

import (
	"context"
	"errors"
	"fmt"
	casbin2 "github.com/casbin/casbin/v2"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/rqdata"
	"github.com/miekg/dns"
	"reflect"
	"strings"
)

type casbin struct {
	engines map[string]*engine
	next    plugin.Handler
}

type engine struct {
	modelPath   string
	policyPath  string
	enforcer    *casbin2.Enforcer
	fields      []string
	actionIndex int
	mapping     *rqdata.Mapping
}

func newCasbin() *casbin {
	return &casbin{
		engines: make(map[string]*engine),
	}
}

func newEngine(m *rqdata.Mapping) *engine {
	return &engine{
		mapping: m,
	}
}

func (c *casbin) Name() string {
	return "casbin"
}

func (c *casbin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return plugin.NextOrFailure(c.Name(), c.next, ctx, w, r)
}

func (c *casbin) Engine(name string) policy.Engine {
	return c.engines[name]
}

func (e *engine) BuildQueryData(ctx context.Context, state request.Request) (interface{}, error) {
	return e.buildData(ctx, state, make(map[string]string)), nil
}

func (e *engine) BuildReplyData(ctx context.Context, state request.Request, queryData interface{}) (interface{}, error) {
	return e.buildData(ctx, state, queryData.(map[string]string)), nil
}

func (e *engine) buildData(ctx context.Context, state request.Request, data map[string]string) map[string]string {
	extractor := rqdata.NewExtractor(state, e.mapping)
	for _, f := range e.fields {
		if _, ok := data[f]; ok {
			continue
		}
		var (
			v  string
			ok bool
		)
		if e.mapping.ValidField(f) {
			v, ok = extractor.Value(f)
			if !ok {
				continue
			}
		} else {
			mdf := metadata.ValueFunc(ctx, f)
			v := mdf()
			if v == "" {
				continue
			}
		}
		if len(v) > 0 && v[0] == '[' && strings.HasSuffix(f, "_ip") {
			v = v[1 : len(v)-1]
		}
		data[f] = v
	}
	return data
}

func (e *engine) getFields() {
	m := e.enforcer.GetModel()
	fields := make([]string, 0)
	ast := m["r"]["r"]
	for _, token := range ast.Tokens {
		field := strings.TrimPrefix(token, "r_")
		fields = append(fields, field)
	}
	e.fields = fields
}

func (e *engine) getActionIndex() error {
	m := e.enforcer.GetModel()
	index := -1
	key := "p_action"
	for i, k := range m["p"]["p"].Tokens {
		if k == key {
			index = i
		}
	}
	if index == -1 {
		return errors.New("could not get action column")
	}
	e.actionIndex = index
	return nil
}

func (e *engine) BuildRule(args []string) (policy.Rule, error) {
	return e, nil
}

func (e *engine) Evaluate(data interface{}) (int, error) {
	pdata, ok := data.(map[string]string)
	if !ok {
		return 0, fmt.Errorf("input should be map[string]string instead of %v", reflect.TypeOf(data))
	}
	params := make([]interface{}, 0, len(pdata))
	for _, v := range pdata {
		params = append(params, v)
	}

	ok, p, err := e.enforcer.EnforceEx(params...)
	if err != nil {
		return 0, err
	}
	if len(p) == 0 {
		return policy.TypeNone, nil
	}
	if ok {
		switch p[e.actionIndex] {
		case "allow":
			return policy.TypeAllow, nil
		case "refuse":
			return policy.TypeRefuse, nil
		case "block":
			return policy.TypeBlock, nil
		case "drop":
			return policy.TypeDrop, nil
		default:
			return 0, fmt.Errorf("unknown action: '%s'", p[3])
		}
	}
	return policy.TypeNone, nil
}

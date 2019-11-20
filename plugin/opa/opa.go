package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/rqdata"
	"github.com/miekg/dns"
)

// opa is a policy engine plugin for the firewall plugin that can validate DNS requests and
// replies against OPA servers.
type opa struct {
	engines map[string]*engine
	next    plugin.Handler
}

// engine can validate DNS requests and replies against an OPA server.
type engine struct {
	endpoint string // url to opa server api package e.g. http://example.com/v1/data/dns
	client   *http.Client
	fields   []string        // fields to send as input to opa
	mapping  *rqdata.Mapping // store this so we dont have to rebuild it for every request
}

type input map[string]string

func newOpa() *opa {
	return &opa{engines: make(map[string]*engine)}
}

func newEngine(m *rqdata.Mapping) *engine {
	return &engine{
		mapping: m,
		fields:  []string{"client_ip", "name", "rcode", "response_ip"},
	}
}

// Name implements the Handler interface
func (p *opa) Name() string { return "opa" }

// ServeDNS implements the Handler interface
func (p *opa) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// do nothing
	return plugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

// Engine implements the policy.Engineer interface
func (p *opa) Engine(name string) policy.Engine {
	return p.engines[name]
}

// BuildQueryData implements the policy.Engine interface
func (e *engine) BuildQueryData(ctx context.Context, state request.Request) (interface{}, error) {
	return e.buildData(ctx, state, make(input)), nil
}

// BuildReplyData implements the policy.Engine interface
func (e *engine) BuildReplyData(ctx context.Context, state request.Request, queryData interface{}) (interface{}, error) {
	return e.buildData(ctx, state, queryData.(input)), nil
}

// BuildRule implements the policy.Engine interface
func (e *engine) BuildRule(args []string) (policy.Rule, error) { return e, nil }

// Evaluate implements the policy.Rule interface
func (e *engine) Evaluate(data interface{}) (int, error) {
	// put all query/response data in "input" field, and marshal to json
	bdata, err := json.Marshal(map[string]interface{}{"input": data})
	if err != nil {
		return 0, err
	}

	// send to opa api
	resp, err := e.client.Post(e.endpoint, "application/json", bytes.NewBuffer(bdata))
	if err != nil {
		return 0, err
	}

	// decode response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	action, ok := result["result"]
	if !ok {
		return policy.TypeNone, nil
	}
	switch action {
	case "refuse":
		return policy.TypeRefuse, nil
	case "allow":
		return policy.TypeAllow, nil
	case "block":
		return policy.TypeBlock, nil
	case "drop":
		return policy.TypeDrop, nil
	default:
		return 0, fmt.Errorf("unknown action: '%s'", action)
	}
}

// buildData fills the map of values for policy input
func (e *engine) buildData(ctx context.Context, state request.Request, data input) input {
	extractor := rqdata.NewExtractor(state, e.mapping)
	for _, f := range e.fields {
		if _, ok := data[f]; ok {
			// skip if already defined
			continue
		}
		var v string
		var ok bool
		if e.mapping.ValidField(f) {
			v, ok = extractor.Value(f)
			if !ok {
				continue
			}
		} else {
			mdf := metadata.ValueFunc(ctx, f)
			v = mdf()
			if v == "" {
				continue
			}
		}
		// strip brackets from ipv6 addresses in *_ip fields
		if len(v) > 0 && v[0] == '[' && strings.HasSuffix(f, "_ip") {
			v = v[1 : len(v)-1]
		}
		data[f] = v
	}
	return data
}

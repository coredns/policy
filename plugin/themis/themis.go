package themisplugin

import (
	"errors"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/rqdata"
	"sync"

	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"

	"github.com/infobloxopen/themis/pdp"
	"github.com/infobloxopen/themis/pep"
)

var errInvalidAction = errors.New("invalid action")

const ThemisPluginName = "themis"

// ThemisPlugin represents a plugin instance that can validate DNS
// requests and replies using PDP server.

type ThemisEngine struct {
	conf            config
	trace           plugin.Handler
	next            plugin.Handler
	pdp             pep.Client
	attrPool        attrPool
	attrGauges      *AttrGauge
	connAttempts    map[string]*uint32
	unkConnAttempts *uint32
	mapping *rqdata.Mapping
	wg              sync.WaitGroup
}

func newThemisEngine() *ThemisEngine {
	return &ThemisEngine{
		conf: config{
			options:     make([]*attrSetting, 0),
			custAttrs:   make(map[string]custAttr),
			connTimeout: -1,
			maxReqSize:  -1,
			maxResAttrs: 64,
		},
		connAttempts:    make(map[string]*uint32),
		unkConnAttempts: new(uint32),
		mapping:rqdata.NewMapping(""),
	}
}

func (p *ThemisEngine) BuildQueryData(ctx context.Context, state request.Request) (interface{}, error) {
	ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, p.mapping), p.conf.options, p.attrGauges)
	return ah, nil
}

func (p *ThemisEngine) BuildReplyData(ctx context.Context, state request.Request, queryData interface{}) (interface{}, error) {
	ah := queryData.(*attrHolder)
	ah.prepareResponseFromContext(ctx, rqdata.NewExtractor(state, p.mapping))
	return ah, nil
}

func (p *ThemisEngine) BuildRule(args []string) (policy.Rule, error) {
	return p, nil
}

func (p *ThemisEngine) Evaluate(data interface{}) (int, error) {
	ah := data.(*attrHolder)
	var attrsRequest []pdp.AttributeAssignment
	if !p.conf.autoResAttrs {
		attrsRequest = p.attrPool.Get()
		defer p.attrPool.Put(attrsRequest)
	}
	// validate domain name (validation #1)
	if err := p.validate(ah, attrsRequest); err != nil {
		return dns.RcodeSuccess, err
	}
	return int(ah.action), nil
}

type ThemisPlugin struct {
	engines map[string]*ThemisEngine
	next    plugin.Handler
}

func newThemisPlugin() *ThemisPlugin {
	return &ThemisPlugin{engines: make(map[string]*ThemisEngine)}
}

// Name implements the Handler interface
func (p *ThemisPlugin) Name() string { return ThemisPluginName }

// ServeDNS implements the Handler interface.
func (p *ThemisPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// do nothing
	return plugin.NextOrFailure(p.Name(), p.next, ctx, w, r)
}

func (p *ThemisPlugin) GetEngine(name string) policy.Engine {
	return p.engines[name]
}

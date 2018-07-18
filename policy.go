package policy

import (
	"errors"
	"strings"
	"sync"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/miekg/dns"
	"golang.org/x/net/context"

	"github.com/infobloxopen/themis/pep"
)

var errInvalidAction = errors.New("invalid action")

// policyPlugin represents a plugin instance that can validate DNS
// requests and replies using PDP server.
type policyPlugin struct {
	conf            config
	tapIO           dnstapSender
	trace           plugin.Handler
	next            plugin.Handler
	pdp             pep.Client
	attrPool        attrPool
	attrGauges      *AttrGauge
	connAttempts    map[string]*uint32
	unkConnAttempts *uint32
	wg              sync.WaitGroup
}

func newPolicyPlugin() *policyPlugin {
	return &policyPlugin{
		conf: config{
			options:     make(map[uint16][]*edns0Opt),
			custAttrs:   make(map[string]custAttr),
			connTimeout: -1,
			maxReqSize:  -1,
			maxResAttrs: 64,
		},
		connAttempts:    make(map[string]*uint32),
		unkConnAttempts: new(uint32),
	}
}

// Name implements the Handler interface
func (p *policyPlugin) Name() string { return "policy" }

// ServeDNS implements the Handler interface.
func (p *policyPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var (
		status        = -1
		respMsg       *dns.Msg
		resolveFailed bool
	)
	p.wg.Add(1)
	defer p.wg.Done()

	debug := p.patchDebugMsg(r)

	// turn off default Cq and Cr dnstap messages
	resetCqCr(ctx)

	ah := newAttrHolderWithDnReq(w, r, p.conf.options, p.attrGauges)
	defer func() {
		if ah.action == actionDrop {
			return
		}

		if r != nil {
			r.Rcode = status
			r.Response = true
			clearECS(r)

			if debug && len(r.Question) > 0 {
				q := r.Question[0]

				q.Name += p.conf.debugSuffix
				q.Qtype = dns.TypeTXT
				q.Qclass = dns.ClassCHAOS

				r.Question[0] = q
			}

			if status != dns.RcodeServerFailure || resolveFailed {
				w.WriteMsg(r)
			}
		}

		if p.tapIO != nil && !debug {
			p.tapIO.sendCRExtraMsg(w, r, ah)
		}
	}()

	for _, s := range p.conf.passthrough {
		if strings.HasSuffix(ah.dn, s) {
			nw := nonwriter.New(w)
			_, err := plugin.NextOrFailure(p.Name(), p.next, ctx, nw, r)
			r = nw.Msg
			if r != nil {
				status = r.Rcode

				if debug {
					p.setDebugQueryPassthroughAnswer(ah, r)
					status = dns.RcodeSuccess
				}
			}

			return dns.RcodeSuccess, err
		}
	}

	attrsRequest := p.attrPool.Get()
	defer p.attrPool.Put(attrsRequest)
	// validate domain name (validation #1)
	if err := p.validate(ah, attrsRequest); err != nil {
		status = dns.RcodeServerFailure
		return dns.RcodeSuccess, err
	}

	if ah.action == actionAllow || ah.action == actionLog {
		// resolve domain name to IP
		nw := nonwriter.New(w)
		_, err := plugin.NextOrFailure(p.Name(), p.next, ctx, nw, r)
		if err != nil {
			resolveFailed = true
			status = dns.RcodeServerFailure

			if debug {
				p.setDebugQueryAnswer(ah, r, status)
				status = dns.RcodeSuccess
				return dns.RcodeSuccess, nil
			}

			return dns.RcodeSuccess, err
		}

		respMsg = nw.Msg
		if respMsg == nil {
			r = nil
			return dns.RcodeSuccess, nil
		}

		status = respMsg.Rcode
		if status == dns.RcodeServerFailure {
			resolveFailed = true
		}

		address := getRespIP(respMsg)
		// if external resolver ret code is not RcodeSuccess
		// address is not filled from the answer
		// in this case just pass through answer w/o validation
		if address != nil {
			ah.addIPReq(address)

			attrsResponse := p.attrPool.Get()
			defer p.attrPool.Put(attrsResponse)
			// validate response IP (validation #2)
			if err := p.validate(ah, attrsResponse); err != nil {
				status = dns.RcodeServerFailure
				return dns.RcodeSuccess, err
			}
		}
	}

	if debug && ah.action != actionRefuse {
		p.setDebugQueryAnswer(ah, r, status)
		status = dns.RcodeSuccess

		return dns.RcodeSuccess, nil
	}

	switch ah.action {
	case actionAllow, actionLog:
		r = respMsg
		return dns.RcodeSuccess, nil

	case actionRedirect:
		var err error
		status, err = p.setRedirectQueryAnswer(ctx, w, r, ah.dst)
		r.AuthenticatedData = false
		return dns.RcodeSuccess, err

	case actionBlock:
		status = dns.RcodeNameError
		r.AuthenticatedData = false
		return dns.RcodeSuccess, nil

	case actionRefuse:
		status = dns.RcodeRefused
		return dns.RcodeSuccess, nil

	case actionDrop:
		return dns.RcodeSuccess, nil
	}

	status = dns.RcodeServerFailure
	return dns.RcodeSuccess, errInvalidAction
}

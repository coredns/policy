package policy

import (
	"log"
	"sync/atomic"

	"github.com/coredns/coredns/plugin/pkg/trace"
	"github.com/infobloxopen/themis/pdp"
	"github.com/infobloxopen/themis/pep"
)

// connect establishes connection to PDP server.
func (p *policyPlugin) connect() error {
	log.Printf("[DEBUG] Connecting %v", p)

	for _, addr := range p.conf.endpoints {
		p.connAttempts[addr] = new(uint32)
	}

	opts := []pep.Option{
		pep.WithConnectionTimeout(p.conf.connTimeout),
		pep.WithConnectionStateNotification(p.connStateCb),
	}

	if p.conf.cacheTTL > 0 {
		if p.conf.cacheLimit > 0 {
			opts = append(opts, pep.WithCacheTTLAndMaxSize(p.conf.cacheTTL, p.conf.cacheLimit))
		} else {
			opts = append(opts, pep.WithCacheTTL(p.conf.cacheTTL))
		}
	}

	if p.conf.streams <= 0 || !p.conf.hotSpot {
		opts = append(opts, pep.WithRoundRobinBalancer(p.conf.endpoints...))
	}

	if p.conf.streams > 0 {
		opts = append(opts, pep.WithStreams(p.conf.streams))
		if p.conf.hotSpot {
			opts = append(opts, pep.WithHotSpotBalancer(p.conf.endpoints...))
		}
	}

	opts = append(opts, pep.WithAutoRequestSize(p.conf.autoReqSize))
	if p.conf.maxReqSize > 0 {
		opts = append(opts, pep.WithMaxRequestSize(uint32(p.conf.maxReqSize)))
	}

	p.attrPool = makeAttrPool(p.conf.maxResAttrs, false)

	if p.trace != nil {
		if t, ok := p.trace.(trace.Trace); ok {
			opts = append(opts, pep.WithTracer(t.Tracer()))
		}
	}

	p.pdp = pep.NewClient(opts...)
	return p.pdp.Connect("")
}

// closeConn terminates previously established connection.
func (p *policyPlugin) closeConn() {
	if p.pdp != nil {
		go func() {
			p.wg.Wait()
			p.pdp.Close()
		}()
	}
}

func (p *policyPlugin) validate(ah *attrHolder, a []pdp.AttributeAssignment) error {
	var req []pdp.AttributeAssignment
	if len(ah.ipReq) > 0 {
		req = ah.ipReq
	} else {
		req = ah.dnReq
	}

	if p.conf.log {
		log.Printf("[INFO] PDP request: %+v", req)
	}

	res := pdp.Response{Obligations: a}
	err := p.pdp.Validate(req, &res)
	if err != nil {
		log.Printf("[ERROR] Policy validation failed due to error %s", err)
		return err
	}

	if p.conf.log {
		log.Printf("[INFO] PDP response: %+v", res)
	}

	if len(ah.ipReq) > 0 {
		ah.addIPRes(&res)
	} else {
		ah.addDnRes(&res, p.conf.custAttrs)
	}

	return nil
}

func (p *policyPlugin) connStateCb(addr string, state int, err error) {
	switch state {
	default:
		if err != nil {
			log.Printf("[DEBUG] Unknown connection notification %s (%s)", addr, err)
		} else {
			log.Printf("[DEBUG] Unknown connection notification %s", addr)
		}

	case pep.StreamingConnectionEstablished:
		ptr, ok := p.connAttempts[addr]
		if !ok {
			ptr = p.unkConnAttempts
		}
		atomic.StoreUint32(ptr, 0)

		log.Printf("[INFO] Connected to %s", addr)

	case pep.StreamingConnectionBroken:
		log.Printf("[ERROR] Connection to %s has been broken", addr)

	case pep.StreamingConnectionConnecting:
		ptr, ok := p.connAttempts[addr]
		if !ok {
			ptr = p.unkConnAttempts
		}
		count := atomic.AddUint32(ptr, 1)

		if count <= 1 {
			log.Printf("[INFO] Connecting to %s", addr)
		}

		if count > 100 {
			log.Printf("[ERROR] Connecting to %s", addr)
			atomic.StoreUint32(ptr, 1)
		}

	case pep.StreamingConnectionFailure:
		ptr, ok := p.connAttempts[addr]
		if !ok {
			ptr = p.unkConnAttempts
		}
		if atomic.LoadUint32(ptr) <= 1 {
			log.Printf("[ERROR] Failed to connect to %s (%s)", addr, err)
		}
	}
}

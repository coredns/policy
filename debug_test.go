package policy

import (
	"net"
	"testing"

	"github.com/infobloxopen/themis/pdp"
	"github.com/miekg/dns"
)

func TestPatchDebugMsg(t *testing.T) {
	p := newPolicyPlugin()
	p.conf.debugSuffix = "debug.local."

	m := makeTestDNSMsg("example.com.debug.local", dns.TypeTXT, dns.ClassCHAOS)
	patched := p.patchDebugMsg(m)
	if !patched {
		t.Error("expected patched message")
	}

	assertDNSMessage(t, "patchDebugMsg", 0, m, 0,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n"+
			";example.com.\tIN\t A\n",
	)

	m = new(dns.Msg)
	patched = p.patchDebugMsg(m)
	if patched {
		t.Error("expected NOT patched message")
	}

	assertDNSMessage(t, "patchDebugMsg(nil)", 0, m, 0,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 0, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n",
	)

	m = makeTestDNSMsg("example.com", dns.TypeTXT, dns.ClassCHAOS)
	patched = p.patchDebugMsg(m)
	if patched {
		t.Error("expected NOT patched message")
	}

	assertDNSMessage(t, "patchDebugMsg(nil)", 0, m, 0,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n"+
			";example.com.\tCH\t TXT\n",
	)
}

func TestSetDebugQueryPassthroughAnswer(t *testing.T) {
	p := newPolicyPlugin()
	p.conf.debugSuffix = "debug.local."

	m := makeTestDNSMsg("example.passthrough.local", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")

	ah := newAttrHolderWithDnReq(w, m, nil, nil)

	p.setDebugQueryPassthroughAnswer(ah, m)
	assertDNSMessage(t, "setDebugQueryPassthroughAnswer", 0, m, 0,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n"+
			";example.passthrough.local.\tIN\t A\n\n"+
			";; ANSWER SECTION:\n"+
			"example.passthrough.local.debug.local.\t0\tCH\tTXT\t\"action:passthrough\"\n",
	)
}

func TestSetDebugQueryAnswer(t *testing.T) {
	p := newPolicyPlugin()
	p.conf.debugID = "<DEBUG>"
	p.conf.debugSuffix = "debug.local."

	t.Run("OnDomainResponse(RcodeSuccess)", func(t *testing.T) {
		m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
		w := newTestAddressedNonwriter("192.0.2.1")

		ah := newAttrHolderWithDnReq(w, m, nil, nil)
		ah.addDnRes(&pdp.Response{
			Effect: pdp.EffectDeny,
			Obligations: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment(attrNameRedirectTo, "192.0.2.54"),
			},
		}, nil)

		p.setDebugQueryAnswer(ah, m, dns.RcodeSuccess)
		assertDNSMessage(t, "setDebugQueryAnswer", 0, m, 0,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t"+
				"\"resolve:yes,query:'redirect',redirect_to:'192.0.2.54',ident:'<DEBUG>'\"\n",
		)
	})

	t.Run("OnDomainResponse(RcodeNameError)", func(t *testing.T) {
		m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
		w := newTestAddressedNonwriter("192.0.2.1")

		ah := newAttrHolderWithDnReq(w, m, nil, nil)
		ah.addDnRes(&pdp.Response{
			Effect: pdp.EffectPermit,
		}, nil)

		p.setDebugQueryAnswer(ah, m, dns.RcodeNameError)
		assertDNSMessage(t, "setDebugQueryAnswer", 0, m, 0,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t"+
				"\"resolve:no,query:'allow',ident:'<DEBUG>'\"\n",
		)
	})

	t.Run("OnDomainResponse(RcodeServerFailure)", func(t *testing.T) {
		m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
		w := newTestAddressedNonwriter("192.0.2.1")

		ah := newAttrHolderWithDnReq(w, m, nil, nil)
		ah.addDnRes(&pdp.Response{
			Effect: pdp.EffectPermit,
		}, nil)

		p.setDebugQueryAnswer(ah, m, dns.RcodeServerFailure)
		assertDNSMessage(t, "setDebugQueryAnswer", 0, m, 0,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t"+
				"\"resolve:failed,query:'allow',ident:'<DEBUG>'\"\n",
		)
	})

	t.Run("OnDomainResponse(-1)", func(t *testing.T) {
		m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
		w := newTestAddressedNonwriter("192.0.2.1")

		ah := newAttrHolderWithDnReq(w, m, nil, nil)
		ah.addDnRes(&pdp.Response{
			Effect: pdp.EffectPermit,
		}, nil)

		p.setDebugQueryAnswer(ah, m, -1)
		assertDNSMessage(t, "setDebugQueryAnswer", 0, m, 0,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t"+
				"\"resolve:skip,query:'allow',ident:'<DEBUG>'\"\n",
		)
	})

	t.Run("OnIPResponse", func(t *testing.T) {
		m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
		w := newTestAddressedNonwriter("192.0.2.1")

		ah := newAttrHolderWithDnReq(w, m, nil, nil)
		ah.addDnRes(&pdp.Response{
			Effect: pdp.EffectPermit,
		}, nil)

		ah.addIPReq(net.ParseIP("192.0.2.53"))
		ah.addIPRes(&pdp.Response{
			Effect: pdp.EffectDeny,
			Obligations: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment(attrNameRedirectTo, "192.0.2.54"),
			},
		})

		p.setDebugQueryAnswer(ah, m, dns.RcodeSuccess)
		assertDNSMessage(t, "setDebugQueryAnswer", 0, m, 0,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t"+
				"\"resolve:yes,query:'pass',response:'redirect',redirect_to:'192.0.2.54',ident:'<DEBUG>'\"\n",
		)
	})
}

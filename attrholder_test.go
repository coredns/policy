package policy

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/miekg/dns"

	pb "github.com/infobloxopen/themis/contrib/coredns/policy/dnstap"
	"github.com/infobloxopen/themis/pdp"
)

func TestNewAttrHolderWithDnReq(t *testing.T) {
	optsMap := map[uint16][]*edns0Opt{
		0xfffc: {
			{
				name:     attrNameSourceIP,
				dataType: typeEDNS0IP,
			},
		},
		0xfffd: {
			{
				name:     "low",
				dataType: typeEDNS0Hex,
				size:     16,
				end:      8,
			},
			{
				name:     "high",
				dataType: typeEDNS0Hex,
				size:     16,
				start:    8,
			},
		},
		0xfffe: {
			{
				name:     "byte",
				dataType: typeEDNS0Bytes,
			},
		},
	}

	m := makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET,
		newEdns0(
			newEdns0Local(0xfffd,
				[]byte{
					0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7,
					0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
				},
			),
			newEdns0Local(0xfffd, []byte{}),
			newEdns0Local(0xfffe, []byte("test")),
			newEdns0Local(0xfffc, []byte(net.ParseIP("2001:db8::1"))),
		),
	)
	w := newTestAddressedNonwriter("192.0.2.1")

	ah := newAttrHolderWithDnReq(w, m, optsMap, nil)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq", ah.dnReq,
		pdp.MakeStringAssignment(attrNameType, typeValueQuery),
		pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain(dns.Fqdn("example.com"))),
		pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.TypeA), 16)),
		pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("2001:db8::1")),
		pdp.MakeStringAssignment("low", "0001020304050607"),
		pdp.MakeStringAssignment("high", "08090a0b0c0d0e0f"),
		pdp.MakeStringAssignment("byte", "test"),
	)

	m = makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET,
		newEdns0(
			newEdns0Local(0xfffd,
				[]byte{
					0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7,
					0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
				},
			),
			newEdns0Local(0xfffd, []byte{}),
			newEdns0Local(0xfffe, []byte("test")),
			newEdns0Local(0xfffc, []byte(net.ParseIP("2001:db8::1"))),
		),
	)
	w = newTestAddressedNonwriter("example.com:53")

	ah = newAttrHolderWithDnReq(w, m, optsMap, nil)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq(notIPRemoteAddr)", ah.dnReq,
		pdp.MakeStringAssignment(attrNameType, typeValueQuery),
		pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain(dns.Fqdn("example.com"))),
		pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.TypeA), 16)),
		pdp.MakeStringAssignment("low", "0001020304050607"),
		pdp.MakeStringAssignment("high", "08090a0b0c0d0e0f"),
		pdp.MakeStringAssignment("byte", "test"),
		pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("2001:db8::1")),
	)

	m = makeTestDNSMsg("...", dns.TypeA, dns.ClassINET)
	assertPanicWithError(t, "newAttrHolderWithDnReq(invalidDomainName)", func() {
		newAttrHolderWithDnReq(w, m, nil, nil)
	}, "Can't treat %q as domain name: %s", "...", domain.ErrEmptyLabel)
}

func TestAddIpReq(t *testing.T) {
	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")

	custAttrs := map[string]custAttr{
		"trans": custAttrTransfer,
	}

	ah := newAttrHolderWithDnReq(w, m, nil, nil)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq - dnReq", ah.dnReq,
		pdp.MakeStringAssignment(attrNameType, typeValueQuery),
		pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain(dns.Fqdn("example.com"))),
		pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.TypeA), 16)),
		pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.1")),
	)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq - dnRes", ah.dnRes)

	ah.addDnRes(
		&pdp.Response{
			Effect: pdp.EffectPermit,
			Obligations: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans", "val"),
			},
		},
		custAttrs,
	)
	pdp.AssertAttributeAssignments(t, "addDnRes - dnRes", ah.dnRes,
		pdp.MakeStringAssignment("trans", "val"),
	)
	pdp.AssertAttributeAssignments(t, "addDnRes - transfer", ah.transfer,
		pdp.MakeStringAssignment("trans", "val"),
	)

	ah.addIPReq(net.ParseIP("2001:db8::1"))
	pdp.AssertAttributeAssignments(t, "addIPReq - ipReq", ah.ipReq,
		pdp.MakeStringAssignment(attrNameType, typeValueResponse),
		pdp.MakeAddressAssignment(attrNameAddress, net.ParseIP("2001:db8::1")),
		pdp.MakeStringAssignment("trans", "val"),
	)
}

func TestActionDomainResponse(t *testing.T) {
	tests := []struct {
		res    *pdp.Response
		action byte
		dst    string
	}{
		{
			res: &pdp.Response{
				Effect: pdp.EffectPermit,
			},
			action: actionAllow,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectIndeterminate,
				Status: fmt.Errorf("example of pdp failure on domain validation"),
			},
			action: actionInvalid,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
			},
			action: actionBlock,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment(attrNameRedirectTo, "192.0.2.1"),
				},
			},
			action: actionRedirect,
			dst:    "192.0.2.1",
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeIntegerAssignment(attrNameRedirectTo, 0),
				},
			},
			action: actionInvalid,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameRefuse, true),
				},
			},
			action: actionRefuse,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameLog, true),
				},
			},
			action: actionLog,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameDrop, true),
				},
			},
			action: actionDrop,
		},
	}

	r := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ah := newAttrHolderWithDnReq(w, r, nil, nil)

			g := newLogGrabber()
			ah.addDnRes(test.res, nil)
			logs := g.Release()

			if ah.action != test.action {
				t.Errorf("unexpected action in TC #%d: expected=%d, actual=%d", i, test.action, ah.action)
				t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
			}
			if ah.dst != test.dst {
				t.Errorf("unexpected redirect destination in TC #%d: expected=%q, actual=%q", i, test.dst, ah.dst)
				t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
			}
		})
	}
}

func TestActionIpResponse(t *testing.T) {
	tests := []struct {
		res        *pdp.Response
		initAction byte
		action     byte
		dst        string
	}{
		{
			res: &pdp.Response{
				Effect: pdp.EffectPermit,
			},
			initAction: actionAllow,
			action:     actionAllow,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectPermit,
			},
			initAction: actionLog,
			action:     actionLog,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameLog, true),
				},
			},
			initAction: actionAllow,
			action:     actionLog,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
			},
			initAction: actionAllow,
			action:     actionBlock,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment(attrNameRedirectTo, "192.0.2.1"),
				},
			},
			initAction: actionAllow,
			action:     actionRedirect,
			dst:        "192.0.2.1",
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameRefuse, true),
				},
			},
			initAction: actionAllow,
			action:     actionRefuse,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameDrop, true),
				},
			},
			initAction: actionAllow,
			action:     actionDrop,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectIndeterminate,
				Status: fmt.Errorf("example of pdp failure on IP validation"),
			},
			initAction: actionAllow,
			action:     actionInvalid,
		},
	}

	r := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ah := newAttrHolderWithDnReq(w, r, nil, nil)
			ah.action = test.initAction

			g := newLogGrabber()
			ah.addIPRes(test.res)
			logs := g.Release()

			if ah.action != test.action {
				t.Errorf("unexpected action in TC #%d: expected=%d, actual=%d", i, test.action, ah.action)
				t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
			}
			if ah.dst != test.dst {
				t.Errorf("unexpected redirect destination in TC #%d: expected=%q, actual=%q", i, test.dst, ah.dst)
				t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
			}
		})
	}
}

func TestAddResponse(t *testing.T) {
	optsMap := map[uint16][]*edns0Opt{
		0xfffe: {
			{
				name:     "edns1",
				dataType: typeEDNS0Bytes,
			},
		},
	}

	opts := []*dns.OPT{
		newEdns0(
			newEdns0Local(0xfffe, []byte("edns1Val")),
		),
	}

	custAttrs := map[string]custAttr{
		"edns1":       custAttrEdns,
		"trans1":      custAttrTransfer,
		"transdnstap": custAttrTransfer | custAttrDnstap,
		"dnstap":      custAttrDnstap,
	}

	tests := []struct {
		opts   []*dns.OPT
		resp   *pdp.Response
		dnRes  []pdp.AttributeAssignment
		ipRes  []pdp.AttributeAssignment
		edns0  []pdp.AttributeAssignment
		trans  []pdp.AttributeAssignment
		dnstap []pdp.AttributeAssignment
	}{
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
			},
			dnRes: []pdp.AttributeAssignment{},
			ipRes: []pdp.AttributeAssignment{},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.1")),
			},
			trans:  []pdp.AttributeAssignment{},
			dnstap: []pdp.AttributeAssignment{},
		},
		{
			opts: opts,
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
			},
			dnRes: []pdp.AttributeAssignment{},
			ipRes: []pdp.AttributeAssignment{},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.2")),
				pdp.MakeStringAssignment("edns1", "edns1Val"),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "edns1Val"),
				},
			},
			dnRes: []pdp.AttributeAssignment{},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.3")),
				pdp.MakeStringAssignment("edns1", "edns1Val"),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "edns1Val"),
				},
			},
			dnRes: []pdp.AttributeAssignment{},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.4")),
				pdp.MakeStringAssignment("edns1", "edns1Val"),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "edns1Val"),
				},
			},
			ipRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("edns1", "edns1Val"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.5")),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			opts: opts,
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "edns1Val2"),
				},
			},
			dnRes: []pdp.AttributeAssignment{},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.6")),
				pdp.MakeStringAssignment("edns1", "edns1Val"),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "edns1Val2"),
				},
			},
			dnRes: []pdp.AttributeAssignment{},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.7")),
				pdp.MakeStringAssignment("edns1", "edns1Val2"),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
				},
			},
			ipRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.8")),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
				},
			},
			dnRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.9")),
			},
			trans: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
				},
			},
			dnRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.10")),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "ends1Val"),
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
					pdp.MakeStringAssignment("other1", "other1Val1"),
				},
			},
			dnRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("other1", "other1Val1"),
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.11")),
				pdp.MakeStringAssignment("edns1", "ends1Val"),
			},
			trans: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "ends1Val"),
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
					pdp.MakeStringAssignment("other1", "other1Val1"),
				},
			},
			ipRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("edns1", "ends1Val"),
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
				pdp.MakeStringAssignment("other1", "other1Val1"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.12")),
			},
			trans: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "ends1Val"),
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
					pdp.MakeStringAssignment("other1", "other1Val1"),
					pdp.MakeStringAssignment("transdnstap", "val"),
				},
			},
			dnRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("transdnstap", "val"),
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
				pdp.MakeStringAssignment("other1", "other1Val1"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.13")),
				pdp.MakeStringAssignment("edns1", "ends1Val"),
			},
			trans: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("transdnstap", "val"),
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
			},
			dnstap: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("transdnstap", "val"),
			},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("edns1", "ends1Val"),
					pdp.MakeStringAssignment("trans1", "trans1Val1"),
					pdp.MakeStringAssignment("other1", "other1Val1"),
					pdp.MakeStringAssignment("transdnstap", "val"),
				},
			},
			ipRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("edns1", "ends1Val"),
				pdp.MakeStringAssignment("trans1", "trans1Val1"),
				pdp.MakeStringAssignment("other1", "other1Val1"),
				pdp.MakeStringAssignment("transdnstap", "val"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.14")),
			},
			trans:  []pdp.AttributeAssignment{},
			dnstap: []pdp.AttributeAssignment{},
		},
		{
			resp: &pdp.Response{
				Effect: pdp.EffectPermit,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeStringAssignment("other1", "other1Val1"),
					pdp.MakeStringAssignment("dnstap", "val"),
				},
			},
			dnRes: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("other1", "other1Val1"),
				pdp.MakeStringAssignment("dnstap", "val"),
			},
			edns0: []pdp.AttributeAssignment{
				pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.15")),
			},
			trans: []pdp.AttributeAssignment{},
			dnstap: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("dnstap", "val"),
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := newTestAddressedNonwriter(fmt.Sprintf("192.0.2.%d", i+1))
			if test.dnRes != nil {
				r := makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET, copyEdns0(test.opts...)...)
				ah := newAttrHolderWithDnReq(w, r, optsMap, nil)
				ah.addDnRes(test.resp, custAttrs)

				pdp.AssertAttributeAssignments(t,
					fmt.Sprintf("TestAddResponse test %d dnRes", i+1),
					ah.dnRes, test.dnRes...,
				)
				if test.edns0 != nil {
					pdp.AssertAttributeAssignments(t,
						fmt.Sprintf("TestAddResponse test %d (dnRes) - edns0", i+1),
						ah.dnReq[ednsAttrsStart:], test.edns0...,
					)
				}
				if test.trans != nil {
					pdp.AssertAttributeAssignments(t,
						fmt.Sprintf("TestAddResponse test %d (dnRes) - trans", i+1),
						ah.transfer, test.trans...,
					)
				}
				if test.dnstap != nil {
					pdp.AssertAttributeAssignments(t,
						fmt.Sprintf("TestAddResponse test %d (dnRes) - dnstap", i+1),
						ah.dnstap, test.dnstap...,
					)
				}
			}

			if test.ipRes != nil {
				r := makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET, copyEdns0(test.opts...)...)
				ah := newAttrHolderWithDnReq(w, r, optsMap, nil)
				ah.addIPRes(test.resp)

				pdp.AssertAttributeAssignments(t,
					fmt.Sprintf("TestAddResponse test %d ipRes", i+1),
					ah.dnRes, test.dnRes...,
				)
				if test.edns0 != nil {
					pdp.AssertAttributeAssignments(t,
						fmt.Sprintf("TestAddResponse test %d (ipRes) - edns0", i+1),
						ah.dnReq[ednsAttrsStart:], test.edns0...,
					)
				}
				if test.trans != nil {
					pdp.AssertAttributeAssignments(t,
						fmt.Sprintf("TestAddResponse test %d (ipRes) - trans", i+1),
						ah.transfer, test.trans...,
					)
				}
				if test.dnstap != nil {
					pdp.AssertAttributeAssignments(t,
						fmt.Sprintf("TestAddResponse test %d (ipRes) - dnstap", i+1),
						ah.dnstap, test.dnstap...,
					)
				}
			}
		})
	}
}

func TestMakeDnstapReport(t *testing.T) {
	optsMap := map[uint16][]*edns0Opt{
		0xfffe: {
			{
				name:     "edns",
				dataType: typeEDNS0Bytes,
			},
		},
	}

	m := makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET,
		newEdns0(
			newEdns0Local(0xfffe, []byte("ednsVal")),
		),
	)
	w := newTestAddressedNonwriter("192.0.2.1")

	custAttrs := map[string]custAttr{
		"edns":   custAttrEdns,
		"trans":  custAttrTransfer,
		"dnstap": custAttrDnstap,
	}

	ah := newAttrHolderWithDnReq(w, m, optsMap, nil)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq - dnReq", ah.dnReq,
		pdp.MakeStringAssignment(attrNameType, typeValueQuery),
		pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain(dns.Fqdn("example.com"))),
		pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.TypeA), 16)),
		pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("192.0.2.1")),
		pdp.MakeStringAssignment("edns", "ednsVal"),
	)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq - dnRes", ah.dnRes)

	ah.addDnRes(
		&pdp.Response{
			Effect: pdp.EffectPermit,
			Obligations: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("trans", "transVal"),
				pdp.MakeStringAssignment("dnstap", "dnstapVal"),
			},
		},
		custAttrs,
	)
	pdp.AssertAttributeAssignments(t, "addDnRes - dnRes", ah.dnRes,
		pdp.MakeStringAssignment("trans", "transVal"),
		pdp.MakeStringAssignment("dnstap", "dnstapVal"),
	)
	pdp.AssertAttributeAssignments(t, "addDnRes - transfer", ah.transfer,
		pdp.MakeStringAssignment("trans", "transVal"),
	)
	pdp.AssertAttributeAssignments(t, "addDnRes - dnstap", ah.dnstap,
		pdp.MakeStringAssignment("dnstap", "dnstapVal"),
	)

	ah.addIPReq(net.ParseIP("2001:db8::1"))
	pdp.AssertAttributeAssignments(t, "addIPReq - ipReq", ah.ipReq,
		pdp.MakeStringAssignment(attrNameType, typeValueResponse),
		pdp.MakeAddressAssignment(attrNameAddress, net.ParseIP("2001:db8::1")),
		pdp.MakeStringAssignment("trans", "transVal"),
	)

	ah.addIPRes(
		&pdp.Response{
			Effect: pdp.EffectPermit,
			Obligations: []pdp.AttributeAssignment{
				pdp.MakeStringAssignment("other", "otherVal"),
			},
		},
	)

	pdp.AssertAttributeAssignments(t, "addIPRes - dnRes", ah.dnRes,
		pdp.MakeStringAssignment("trans", "transVal"),
		pdp.MakeStringAssignment("dnstap", "dnstapVal"),
	)
	pdp.AssertAttributeAssignments(t, "addIPRes - ipRes", ah.ipRes,
		pdp.MakeStringAssignment("other", "otherVal"),
	)
	pdp.AssertAttributeAssignments(t, "addIPRes - transfer", ah.transfer,
		pdp.MakeStringAssignment("trans", "transVal"),
	)
	pdp.AssertAttributeAssignments(t, "addIPRes - dnstap", ah.dnstap,
		pdp.MakeStringAssignment("dnstap", "dnstapVal"),
	)

	assertDnstapAttributes(t, "makeDnstapReport", ah.makeDnstapReport(),
		&pb.DnstapAttribute{Id: attrNameSourceIP, Value: "192.0.2.1"},
		&pb.DnstapAttribute{Id: "edns", Value: "ednsVal"},
		&pb.DnstapAttribute{Id: "dnstap", Value: "dnstapVal"},
	)

	ah.action = actionLog
	assertDnstapAttributes(t, "makeDnstapReport(full)", ah.makeDnstapReport(),
		&pb.DnstapAttribute{Id: attrNameDomainName, Value: dns.Fqdn("example.com")},
		&pb.DnstapAttribute{Id: attrNameDNSQtype, Value: strconv.FormatUint(uint64(dns.TypeA), 16)},
		&pb.DnstapAttribute{Id: attrNameSourceIP, Value: "192.0.2.1"},
		&pb.DnstapAttribute{Id: "edns", Value: "ednsVal"},
		&pb.DnstapAttribute{Id: "trans", Value: "transVal"},
		&pb.DnstapAttribute{Id: "dnstap", Value: "dnstapVal"},
		&pb.DnstapAttribute{Id: attrNameAddress, Value: "2001:db8::1"},
		&pb.DnstapAttribute{Id: "other", Value: "otherVal"},
		&pb.DnstapAttribute{Id: attrNamePolicyAction, Value: dnstapActionValues[actionLog]},
		&pb.DnstapAttribute{Id: attrNameType, Value: typeValueResponse},
	)
}

func makeTestDomain(s string) domain.Name {
	dn, err := domain.MakeNameFromString(s)
	if err != nil {
		panic(err)
	}

	return dn
}

func assertPanicWithError(t *testing.T, desc string, f func(), format string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			e := fmt.Sprintf(format, args...)
			err, ok := r.(error)
			if !ok {
				t.Errorf("excpected error %q on panic for %q but got %T (%#v)", e, desc, r, r)
			} else if err.Error() != e {
				t.Errorf("excpected error %q on panic for %q but got %q", e, desc, r)
			}
		} else {
			t.Errorf("expected panic %q for %q", fmt.Sprintf(format, args...), desc)
		}
	}()

	f()
}

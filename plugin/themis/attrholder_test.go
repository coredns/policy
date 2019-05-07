package themisplugin

import (
	"fmt"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/rqdata"
	"net"
	"strconv"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/miekg/dns"
	"github.com/coredns/coredns/request"

	"context"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/infobloxopen/themis/pdp"
)

type fakeWriter struct {
	dns.ResponseWriter
	clientIP string
	serverIP string
}

func (w *fakeWriter) LocalAddr() net.Addr {
	local := net.ParseIP(w.clientIP)
	return &net.UDPAddr{IP: local, Port: 53} // Port is not used here
}

func (w *fakeWriter) RemoteAddr() net.Addr {
	remote := net.ParseIP(w.serverIP)
	return &net.UDPAddr{IP: remote, Port: 53} // Port is not used here
}

func buildFunc(v string) metadata.Func {
	return func() string { return v }
}
func buildContext(ctx context.Context, data map[string]string) context.Context {
//	ctx = metadata.ForMetadata(ctx)
	for k, v := range data {
		metadata.SetValueFunc(ctx, k, buildFunc(v))
	}
	return ctx
}

func buildState(name string, qtype uint16, clientIp string) request.Request {
	//create a msg
	req := new(dns.Msg)
	req.SetQuestion(name, qtype)
	return request.Request{Req:req, W:&fakeWriter{clientIP:clientIp}}
}

func TestNewAttrHolderWithDnReq(t *testing.T) {
	optsMap := []*attrSetting{
		{"low", "request/low", "String", false},
		{"high", "request/high", "String", false},
		{"byte", "request/byte", "String", false},
		{attrNameSourceIP, "request/source_ip", "Address", false},
	}

	mapping := rqdata.NewMapping("")
	state := buildState("example.com.", dns.TypeA, "192.0.2.1")
	mdata := map[string]string{
		"request/source_ip":            "2001:db8::1",
		"request/low":                  "0001020304050607",
		"request/high":                 "08090a0b0c0d0e0f",
		"request/byte":                 "test",
	}

	ctx := context.TODO()
	ctx = buildContext(ctx, mdata)

	ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping), optsMap, nil)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq", ah.dnReq,
		pdp.MakeStringAssignment(attrNameType, typeValueQuery),
		pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain(dns.Fqdn("example.com"))),
		pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.TypeA), 16)),
		pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("2001:db8::1")),
		pdp.MakeStringAssignment("low", "0001020304050607"),
		pdp.MakeStringAssignment("high", "08090a0b0c0d0e0f"),
		pdp.MakeStringAssignment("byte", "test"),
	)

	state = buildState("example.com.", dns.TypeA, "example.com:53")
	mdata = map[string]string{
		"request/source_ip":            "2001:db8::1",
		"request/low":                  "0001020304050607",
		"request/high":                 "08090a0b0c0d0e0f",
		"request/byte":                 "test",
	}

	ctx = buildContext(context.TODO(), mdata)
	ah = newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping), optsMap, nil)
	pdp.AssertAttributeAssignments(t, "newAttrHolderWithDnReq(notIPRemoteAddr)", ah.dnReq,
		pdp.MakeStringAssignment(attrNameType, typeValueQuery),
		pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain(dns.Fqdn("example.com"))),
		pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.TypeA), 16)),
		pdp.MakeStringAssignment("low", "0001020304050607"),
		pdp.MakeStringAssignment("high", "08090a0b0c0d0e0f"),
		pdp.MakeStringAssignment("byte", "test"),
		pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("2001:db8::1")),
	)

	state = buildState("...", dns.TypeA, "example.com:53")
	ctx = buildContext(context.TODO(), mdata)
	assertPanicWithError(t, "newAttrHolderWithDnReq(invalidDomainName)", func() {
		newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping), optsMap, nil)
	}, "Can't treat %q as domain name: %s", "...", domain.ErrEmptyLabel)
}

func TestAddIpReq(t *testing.T) {

	optsMap := []*attrSetting{
		{attrNameSourceIP, "request/source_ip", "Address", false},
	}
	mdata := map[string]string{}
	mapping := rqdata.NewMapping("")
	state := buildState("example.com.", dns.TypeA, "192.0.2.1")

	ctx := buildContext(context.TODO(), mdata)

	custAttrs := map[string]custAttr{
		"trans": custAttrTransfer,
	}

	ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping),optsMap, nil)
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
			action: policy.TypeAllow,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectIndeterminate,
				Status: fmt.Errorf("example of pdp failure on domain validation"),
			},
			action: policy.TypeNone,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
			},
			action: policy.TypeBlock,
		},
		//{
		//	res: &pdp.Response{
		//		Effect: pdp.EffectDeny,
		//		Obligations: []pdp.AttributeAssignment{
		//			pdp.MakeStringAssignment(attrNameRedirectTo, "192.0.2.1"),
		//		},
		//	},
		//	action: firewall.TypeRedirect,
		//	dst:    "192.0.2.1",
		//},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeIntegerAssignment(attrNameRedirectTo, 0),
				},
			},
			action: policy.TypeNone,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameRefuse, true),
				},
			},
			action: policy.TypeRefuse,
		},
		//{
		//	res: &pdp.Response{
		//		Effect: pdp.EffectPermit,
		//		Obligations: []pdp.AttributeAssignment{
		//			pdp.MakeBooleanAssignment(attrNameLog, true),
		//		},
		//	},
		//	action: firewall.TypeLog,
		//},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameDrop, true),
				},
			},
			action: policy.TypeDrop,
		},
	}

	mdata := map[string]string{}
	mapping := rqdata.NewMapping("")
	state := buildState("example.com.", dns.TypeA, "192.0.2.1")

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			optMap := make([]*attrSetting, 0)
			ctx := buildContext(context.TODO(), mdata)
			ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping), optMap, nil)

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
			initAction: policy.TypeAllow,
			action:     policy.TypeAllow,
		},
		//{
		//	res: &pdp.Response{
		//		Effect: pdp.EffectPermit,
		//	},
		//	initAction: policy.TypeLog,
		//	action:     policy.TypeAllow,
		//},
		//{
		//	res: &pdp.Response{
		//		Effect: pdp.EffectPermit,
		//		Obligations: []pdp.AttributeAssignment{
		//			pdp.MakeBooleanAssignment(attrNameLog, true),
		//		},
		//	},
		//	initAction: policy.TypeAllow,
		//	action:     policy.TypeLog,
		//},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
			},
			initAction: policy.TypeAllow,
			action:     policy.TypeBlock,
		},
		//{
		//	res: &pdp.Response{
		//		Effect: pdp.EffectDeny,
		//		Obligations: []pdp.AttributeAssignment{
		//			pdp.MakeStringAssignment(attrNameRedirectTo, "192.0.2.1"),
		//		},
		//	},
		//	initAction: policy.TypeAllow,
		//	action:     policy.TypeRedirect,
		//	dst:        "192.0.2.1",
		//},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameRefuse, true),
				},
			},
			initAction: policy.TypeAllow,
			action:     policy.TypeRefuse,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectDeny,
				Obligations: []pdp.AttributeAssignment{
					pdp.MakeBooleanAssignment(attrNameDrop, true),
				},
			},
			initAction: policy.TypeAllow,
			action:     policy.TypeDrop,
		},
		{
			res: &pdp.Response{
				Effect: pdp.EffectIndeterminate,
				Status: fmt.Errorf("example of pdp failure on IP validation"),
			},
			initAction: policy.TypeAllow,
			action:     policy.TypeNone,
		},
	}
	mdata := map[string]string{}
	mapping := rqdata.NewMapping("")
	state := buildState("example.com.", dns.TypeA, "192.0.2.1")

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			optMap := make([]*attrSetting, 0)
			ctx := buildContext(context.TODO(), mdata)
			ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping), optMap, nil)

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
	mdata := map[string]string{
		"request/edns1":                "edns1Val",
	}
	mapping := rqdata.NewMapping("")
	state := buildState("example.com.", dns.TypeA, "192.0.2.1")

	optsMap := []*attrSetting{
		{"edns1", "request/edns1", "String", false},
	}

	custAttrs := map[string]custAttr{
		"edns1":       custAttrEdns,
		"trans1":      custAttrTransfer,
		"transdnstap": custAttrTransfer | custAttrDnstap,
		"dnstap":      custAttrDnstap,
	}

	tests := []struct {
		opts   []*attrSetting
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
			opts: optsMap,
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
			opts: optsMap,
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
			state = buildState("example.com.", dns.TypeA, fmt.Sprintf("192.0.2.%d", i+1))
			if test.dnRes != nil {

				optMap := make([]*attrSetting, 0)
				if test.opts != nil {
					optMap = test.opts
				}
				ctx := buildContext(context.TODO(), mdata)
				ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping),optMap, nil)

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
				optMap := make([]*attrSetting, 0)
				if test.opts != nil {
					optMap = test.opts
				}
				ctx := buildContext(context.TODO(), mdata)
				ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping),optMap, nil)
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

	mdata := map[string]string{
		"request/edns1":                "edns1Val",
	}
	mapping := rqdata.NewMapping("")
	state := buildState("example.com.", dns.TypeA, "192.0.2.1")

	optsMap := []*attrSetting{
		{"edns", "request/edns", "String", false},
	}

	custAttrs := map[string]custAttr{
		"edns":   custAttrEdns,
		"trans":  custAttrTransfer,
		"dnstap": custAttrDnstap,
	}

	ctx := buildContext(context.TODO(), mdata)
	ah := newAttrHolderWithContext(ctx, rqdata.NewExtractor(state, mapping), optsMap, nil)
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

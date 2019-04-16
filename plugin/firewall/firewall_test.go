package firewall

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/firewall/policy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

// NextHandler returns a Handler that returns rcode and err.
func ProcessHandler(rcode int, err error) plugin.Handler {
	return plugin.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		if err != nil {
			return dns.RcodeServerFailure, nil
		}

		answer := new(dns.Msg)
		answer.SetRcode(r, rcode)

		w.WriteMsg(answer)

		return rcode, nil
	})
}

func TestFirewallResolution(t *testing.T) {

	tests := []struct {
		queryFilter int
		replyFilter int
		next        int
		resultCode  int
		msgCode     int
		msgNil      bool
	}{
		// This all works because 1 bucket (1 zone, 1 type)
		{policy.TypeDrop, policy.TypeAllow, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeSuccess, true},
		{policy.TypeRefuse, policy.TypeAllow, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeRefused, false},
		{policy.TypeBlock, policy.TypeAllow, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeNameError, false},
		{policy.TypeAllow, policy.TypeAllow, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeSuccess, false},
		{policy.TypeAllow, policy.TypeRefuse, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeRefused, false},
		{policy.TypeAllow, policy.TypeBlock, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeNameError, false},
		{policy.TypeAllow, policy.TypeDrop, dns.RcodeSuccess, dns.RcodeSuccess, dns.RcodeSuccess, true},
	}

	ctx := context.TODO()
	for i, tc := range tests {

		// prepare firewall parameters
		fw, _ := New()
		fw.query.DefaultPolicy = tc.queryFilter
		fw.reply.DefaultPolicy = tc.replyFilter
		fw.next = ProcessHandler(tc.next, nil)

		//create a msg
		req := new(dns.Msg)
		req.SetQuestion("example.com", dns.TypeA)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		rcode, err := fw.ServeDNS(ctx, rec, req)
		if err != nil {
			t.Fatalf("Test %d: Expected no error, but got %s", i, err)
		}

		// now check expectation

		if rcode != tc.resultCode {
			t.Errorf("Test %d: Expected value %s as return code, but got %s", i, dns.RcodeToString[tc.resultCode], dns.RcodeToString[rcode])
		}

		if rec.Rcode != tc.msgCode {
			t.Errorf("Test %d: Expected value %s as DNS reply code, but got %s", i, dns.RcodeToString[tc.msgCode], dns.RcodeToString[rec.Rcode])
		}

		if (rec.Msg == nil) != tc.msgNil {
			t.Errorf("Test %d: Expected MSG to be return as NIL : %v, but got a NIL : %v", i, tc.msgNil, rec.Msg == nil)
		}

	}
}

package whoami

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func TestWhoami(t *testing.T) {
	wh := Whoami{}

	tests := []struct {
		qname         string
		qtype         uint16
		expectedCode  int
		expectedReply []string // ownernames for the records in the additional section.
		expectedErr   error
	}{
		{
			qname:         "example.org",
			qtype:         dns.TypeA,
			expectedCode:  dns.RcodeSuccess,
			expectedReply: []string{"example.org.", "_udp.example.org."},
			expectedErr:   nil,
		},
		// Case insensitive and case preserving
		{
			qname:         "Example.ORG",
			qtype:         dns.TypeA,
			expectedCode:  dns.RcodeSuccess,
			expectedReply: []string{"Example.ORG.", "_udp.Example.ORG."},
			expectedErr:   nil,
		},
	}

	ctx := context.TODO()

	for i, tc := range tests {
		req := new(dns.Msg)
		req.SetQuestion(dns.Fqdn(tc.qname), tc.qtype)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		code, err := wh.ServeDNS(ctx, rec, req)

		if err != tc.expectedErr {
			t.Errorf("Test %d: Expected error %v, but got %v", i, tc.expectedErr, err)
		}
		if code != int(tc.expectedCode) {
			t.Errorf("Test %d: Expected status code %d, but got %d", i, tc.expectedCode, code)
		}
		if len(tc.expectedReply) != 0 {
			for i, expected := range tc.expectedReply {
				actual := rec.Msg.Extra[i].Header().Name
				if actual != expected {
					t.Errorf("Test %d: Expected answer %s, but got %s", i, expected, actual)
				}
			}
		}
	}
}

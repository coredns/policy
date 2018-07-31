package errors

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func TestErrors(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetOutput(&buf)
	em := errorHandler{}

	testErr := errors.New("test error")
	tests := []struct {
		next         plugin.Handler
		expectedCode int
		expectedLog  string
		expectedErr  error
	}{
		{
			next:         genErrorHandler(dns.RcodeSuccess, nil),
			expectedCode: dns.RcodeSuccess,
			expectedLog:  "",
			expectedErr:  nil,
		},
		{
			next:         genErrorHandler(dns.RcodeNotAuth, testErr),
			expectedCode: dns.RcodeNotAuth,
			expectedLog:  fmt.Sprintf("[ERROR] %d %s: %v\n", dns.RcodeNotAuth, "example.org. A", testErr),
			expectedErr:  testErr,
		},
	}

	ctx := context.TODO()
	req := new(dns.Msg)
	req.SetQuestion("example.org.", dns.TypeA)

	for i, tc := range tests {
		em.Next = tc.next
		buf.Reset()
		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		code, err := em.ServeDNS(ctx, rec, req)

		if err != tc.expectedErr {
			t.Errorf("Test %d: Expected error %v, but got %v",
				i, tc.expectedErr, err)
		}
		if code != tc.expectedCode {
			t.Errorf("Test %d: Expected status code %d, but got %d",
				i, tc.expectedCode, code)
		}
		if log := buf.String(); !strings.Contains(log, tc.expectedLog) {
			t.Errorf("Test %d: Expected log %q, but got %q",
				i, tc.expectedLog, log)
		}
	}
}

func genErrorHandler(rcode int, err error) plugin.Handler {
	return plugin.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		return rcode, err
	})
}

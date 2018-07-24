package pep

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
	"github.com/infobloxopen/themis/pdpserver/server"
)

const allPermitPolicy = `# Policy for client tests
attributes:
  x: string

policies:
  alg: FirstApplicableEffect
  rules:
  - effect: Permit
    obligations:
    - x:
       val:
         type: string
         content: AllPermitRule
`

func TestUnaryClientValidation(t *testing.T) {
	pdpServer := startTestPDPServer(allPermitPolicy, 5555, t)
	defer func() {
		if logs := pdpServer.Stop(); len(logs) > 0 {
			t.Logf("server logs:\n%s", logs)
		}
	}()

	t.Run("fixed-buffer", testSingleRequest())
	t.Run("auto-buffer", testSingleRequest(WithAutoRequestSize(true)))
}

func testSingleRequest(opt ...Option) func(t *testing.T) {
	return func(t *testing.T) {
		c := NewClient(opt...)
		err := c.Connect("127.0.0.1:5555")
		if err != nil {
			t.Fatalf("expected no error but got %s", err)
		}
		defer c.Close()

		in := decisionRequest{
			Direction: "Any",
			Policy:    "AllPermitPolicy",
			Domain:    "example.com",
		}
		var out decisionResponse
		err = c.Validate(in, &out)
		if err != nil {
			t.Errorf("expected no error but got %s", err)
		}

		if out.Effect != pdp.EffectPermit || out.Reason != nil || out.X != "AllPermitRule" {
			t.Errorf("got unexpected response: %s", out)
		}
	}
}

func TestUnaryClientValidationWithCache(t *testing.T) {
	pdpServer := startTestPDPServer(allPermitPolicy, 5555, t)
	defer func() {
		if logs := pdpServer.Stop(); len(logs) > 0 {
			t.Logf("server logs:\n%s", logs)
		}
	}()

	c := NewClient(
		WithMaxRequestSize(128),
		WithCacheTTL(15*time.Minute),
	)
	err := c.Connect("127.0.0.1:5555")
	if err != nil {
		t.Fatalf("expected no error but got %s", err)
	}
	defer c.Close()

	uc, ok := c.(*unaryClient)
	if !ok {
		t.Fatalf("expected *unaryClient but got %#v", c)
	}
	bc := uc.cache
	if bc == nil {
		t.Fatal("expected cache")
	}

	in := decisionRequest{
		Direction: "Any",
		Policy:    "AllPermitPolicy",
		Domain:    "example.com",
	}
	var out decisionResponse
	err = c.Validate(in, &out)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	if out.Effect != pdp.EffectPermit || out.Reason != nil || out.X != "AllPermitRule" {
		t.Errorf("got unexpected response: %s", out)
	}

	if bc.Len() == 1 {
		if it := bc.Iterator(); it.SetNext() {
			ei, err := it.Value()
			if err != nil {
				t.Errorf("can't get value from cache: %s", err)
			} else if err := fillResponse(pb.Msg{Body: ei.Value()}, &out); err != nil {
				t.Errorf("can't unmarshal response from cache: %s", err)
			} else if out.Effect != pdp.EffectPermit || out.Reason != nil || out.X != "AllPermitRule" {
				t.Errorf("got unexpected response from cache: %s", out)
			}
		} else {
			t.Error("can't set cache iterator to the first value")
		}
	} else {
		t.Errorf("expected the only record in cache but got %d", bc.Len())
	}

	err = c.Validate(in, &out)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	if out.Effect != pdp.EffectPermit || out.Reason != nil || out.X != "AllPermitRule" {
		t.Errorf("got unexpected response: %s", out)
	}
}

func startTestPDPServer(p string, s uint16, t *testing.T) *loggedServer {
	service := fmt.Sprintf("127.0.0.1:%d", s)
	primary := newServer(
		server.WithServiceAt(service),
	)

	if err := primary.s.ReadPolicies(strings.NewReader(p)); err != nil {
		t.Fatalf("can't read policies: %s", err)
	}

	if err := waitForPortClosed(service); err != nil {
		t.Fatalf("port still in use: %s", err)
	}
	go func() {
		if err := primary.s.Serve(); err != nil {
			t.Fatalf("server failed: %s", err)
		}
	}()

	if err := waitForPortOpened(service); err != nil {
		if logs := primary.Stop(); len(logs) > 0 {
			t.Logf("server logs:\n%s", logs)
		}

		t.Fatalf("can't connect to PDP server: %s", err)
	}
	return primary
}

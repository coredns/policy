package casbin

import (
	"context"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/response"
	"github.com/miekg/dns"
	"testing"
)

func TestEvaluate(t *testing.T) {
	o, err := parse(caddy.NewTestController("dns",
		`casbin myengine {
					model ./examples/model.conf
					policy ./examples/policy.csv
				}`))

	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{
		"client_ip": "10.240.0.1",
		"name":      "example.org.",
	}

	result, err := o.engines["myengine"].Evaluate(data)

	if err != nil {
		t.Fatal(err)
	}

	if result != policy.TypeAllow {
		t.Errorf("Expected %d, got %d.", policy.TypeAllow, result)
	}
}

func TestBuildQueryData(t *testing.T) {
	w := response.NewReader(&test.ResponseWriter{})
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeA)
	state := request.Request{W: w, Req: r}

	o, err := parse(caddy.NewTestController("dns",
		`casbin myengine {
					model ./examples/model.conf
					policy ./examples/policy.csv
				}`))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	d, err := o.Engine("myengine").BuildQueryData(ctx, state)
	if err != nil {
		t.Error(err)
	}
	data := d.(map[string]string)

	if data["client_ip"] != "10.240.0.1" {
		t.Errorf("expected client_ip == '10.240.0.1'. Got '%v'", data["client_ip"])
	}
	if data["name"] != "example.org." {
		t.Errorf("expected name == 'example.org.'. Got '%v'", data["name"])
	}
}

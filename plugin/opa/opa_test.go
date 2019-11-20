package opa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/rqdata"
	"github.com/miekg/dns"
)

func TestEvaluate(t *testing.T) {

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var result map[string]map[string]string
		err := decoder.Decode(&result)
		if err != nil {
			w.Write([]byte("{\"result\":\"json decode error\"}"))
		}
		if _, ok := result["input"]; !ok {
			w.Write([]byte("{\"result\":\"request did not contain input\"}"))
		}
		if result["input"]["a"] != "1" {
			w.Write([]byte("{\"result\":\"expected a -> 1\"}"))
		}
		if result["input"]["b"] != "2" {
			w.Write([]byte("{\"result\":\"expected b -> 2\"}"))
		}

		w.Write([]byte("{\"result\":\"allow\"}"))
	}))

	o, err := parse(caddy.NewTestController("dns",
		`opa myengine {
                 endpoint `+apiStub.URL+`
               }`,
	))

	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{"a": "1", "b": "2"}

	result, err := o.engines["myengine"].Evaluate(data)

	if err != nil {
		t.Fatal(err)
	}

	if result != policy.TypeAllow {
		t.Errorf("Expected %d, got %d.", policy.TypeAllow, result)
	}
}

func TestBuildQueryData(t *testing.T) {
	w := &test.ResponseWriter{}
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeHINFO)
	state := request.Request{W: w, Req: r}

	e := newEngine(rqdata.NewMapping(""))
	ctx := context.TODO()

	d, err := e.BuildQueryData(ctx, state)
	if err != nil {
		t.Error(err)
	}
	data := d.(input)

	if data["client_ip"] != "10.240.0.1" {
		t.Errorf("expected client_ip == '10.240.0.1'. Got '%v'", data["client_ip"])
	}
	if data["name"] != "example.org." {
		t.Errorf("expected name == 'example.org.'. Got '%v'", data["name"])
	}
}

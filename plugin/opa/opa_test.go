package opa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/firewall/policy"
	"github.com/coredns/policy/plugin/pkg/response"
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
	w := response.NewReader(&test.ResponseWriter{})
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeA)
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

func TestBuildReplyData(t *testing.T) {
	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeA)
	m := new(dns.Msg)
	m.SetReply(r)
	m.Rcode = dns.RcodeSuccess
	m.Answer = []dns.RR{test.A("example.org.  5  IN  A  1.2.3.4")}

	w := &response.Reader{Msg: m}
	state := request.Request{W: w, Req: r}

	e := newEngine(rqdata.NewMapping(""))
	ctx := context.TODO()

	indata := input{"client_ip": "10.240.0.1", "name": "test.data.exists."}
	d, err := e.BuildReplyData(ctx, state, indata)
	if err != nil {
		t.Error(err)
	}
	data := d.(input)

	if data["name"] != "test.data.exists." {
		t.Errorf("expected name == 'test.data.exists.'. Got '%v'", data["name"])
	}

	if data["rcode"] != "NOERROR" {
		t.Errorf("expected rcode == 'NOERROR'. Got '%v'", data["rcode"])
	}

	if data["response_ip"] != "1.2.3.4" {
		t.Errorf("expected response_ip == '1.2.3.4'. Got '%v'", data["response_ip"])
	}
}

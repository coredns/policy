package opa

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/coredns/policy/plugin/firewall/policy"
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
                 endpoint `+ apiStub.URL +`
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

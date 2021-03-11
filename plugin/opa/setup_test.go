package opa

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input     string
		expected  *opa
		shouldErr bool
	}{
		{"opa testengine", nil, true},

		{`opa test engine {
                  endpoint test
                }`,
			nil,
			true,
		},

		{`opa testengine {
                  endpoint test
                  fields 1 2 3
                }`,
			&opa{engines: map[string]*engine{
				"testengine": {endpoint: "test", fields: []string{"1", "2", "3"}},
			}},
			false,
		},

		{`opa testengine {
                  endpoint test
                  fields 1 2 3
                }
                opa testengine2 {
                  endpoint test2
                  fields 4
                }`,
			&opa{engines: map[string]*engine{
				"testengine":  {endpoint: "test", fields: []string{"1", "2", "3"}},
				"testengine2": {endpoint: "test2", fields: []string{"4"}},
			}},
			false,
		},
	}

	for i, test := range cases {
		c := caddy.NewTestController("dns", test.input)
		o, err := parse(c)

		if test.shouldErr && err != nil {
			continue
		}

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: expected error but didn't get one for input %s", i, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: expected no error but got one for input %s, got: %v", i, test.input, err)
			}
		}

		if o == nil {
			t.Errorf("Test %d: got nil result for input %s", i, test.input)
		}

		if o.engines == nil {
			t.Errorf("Test %d: got nil engines result for input %s", i, test.input)
			continue
		}

		for name, e := range o.engines {
			if e.endpoint != test.expected.engines[name].endpoint {
				t.Errorf("Test %d: engine '%s' expected endpoint %s, got %s", i, name, test.expected.engines[name].endpoint, e.endpoint)
			}

			if !equal(e.fields, test.expected.engines[name].fields) {
				t.Errorf("Test %d: engine '%s' expected fields %v, got %v", i, name, test.expected.engines[name].fields, e.fields)
			}

		}
	}

}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

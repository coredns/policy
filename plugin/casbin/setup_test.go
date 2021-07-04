package casbin

import (
	"github.com/coredns/caddy"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input     string
		expected  *casbin
		shouldErr bool
	}{
		{"casbin testengine", nil, true},

		{`casbin test engine {
					model path/to/model
					policy path/to/policy
                }`,
			nil,
			true,
		},

		{`casbin testengine {
					model path/to/model
					policy path/to/policy
                }`,
			&casbin{
				engines: map[string]*engine{
					"testengine": {
						modelPath:  "path/to/model",
						policyPath: "path/to/policy",
					},
				},
			},
			false,
		},

		{`casbin testengine {
					model path/to/model
					policy path/to/policy
                }
                casbin testengine2 {
					model path/to/model2
					policy path/to/policy2
                }`,
			&casbin{
				engines: map[string]*engine{
					"testengine": {
						modelPath:  "path/to/model",
						policyPath: "path/to/policy",
					},
					"testengine2": {
						modelPath:  "path/to/model2",
						policyPath: "path/to/policy2",
					},
				},
			},
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
			if e.modelPath != test.expected.engines[name].modelPath {
				t.Errorf("Test %d: engine '%s' expected model path %s, got %s", i, name, test.expected.engines[name].modelPath, e.modelPath)
			}

			if e.policyPath != test.expected.engines[name].policyPath {
				t.Errorf("Test %d: engine '%s' expected policy path %s, got %s", i, name, test.expected.engines[name].modelPath, e.policyPath)
			}
		}
	}

}

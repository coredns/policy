package firewall

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input        string
		shouldErr    bool
		queryNbRules int
		replyNbRules int
	}{
		{`firewall`, true, 0, 0},
		{`firewall {}`, true, 0, 0},
		{`firewall query {
                 }`, false, 0, 0},
		{`firewall query {
				allow true
			}`, false, 1, 0},
		{`firewall query {
				allow true
			}
			firewall response {
                allow true
            }`, false, 1, 1},

		{`firewall query {
				allow true
				drop 1
				refuse 2+1
				block false
			}`, false, 4, 0},
		{`firewall query {
				allow
			}`, true, 1, 0},
		{`firewall query {
				allow invalid expression for our / engine
			}`, true, 1, 0},
		{`firewall query {
				allow true
				opa policy parameter
				themis whatever parameter to themis
 				name-of-plugin name-of-policy paramA paramB paramC
			}`, false, 4, 0},
		{`firewall query {
				allow true
				opa policy parameter
				themis policy parameter to themis
			}`, true, 3, 0},
		{`firewall query {
 				name-of-plugin-error-if-no-policy-name
			}`, true, 1, 0},
	}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		fw, err := parse(c)
		if test.shouldErr && err == nil {
			t.Errorf("Test %v: Expected error but found nil", i)
			continue
		} else if !test.shouldErr && err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
			continue
		}
		if test.shouldErr && err != nil {
			continue
		}

		if len(fw.query.Rules) != test.queryNbRules {
			t.Errorf("Test %v: Expected %v query rules but got %v", i, test.queryNbRules, len(fw.query.Rules))
			continue
		}
		if len(fw.reply.Rules) != test.replyNbRules {
			t.Errorf("Test %v: Expected %v reply rules but got %v", i, test.replyNbRules, len(fw.reply.Rules))
			continue
		}

	}
}

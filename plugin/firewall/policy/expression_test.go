package policy

import (
	"context"
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/pkg/rqdata"
	tst "github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

func TestBuildRule(t *testing.T) {
	tests := []struct {
		expression string
		errorBuild bool
	}{
		{"allow true", false},
		{"block 1 + 1", false},
		{"drop [my/variable] / 20", false},
		{"unknown", true},
		{"untype 'expression'", true},
		{"drop [my/variable / 20", true},
	}
	for i, test := range tests {
		engine := &ExprEngine{TypeDrop, rqdata.NewMapping("-")}
		_, err := engine.BuildRule(strings.Split(test.expression, " "))
		if err != nil {
			if !test.errorBuild {
				t.Errorf("Test %d : unexpected error at build rule : %s", i, err)
			}
			continue
		}
		if test.errorBuild {
			t.Errorf("Test %d : no error at BuilRule returned, when one was expected", i)
		}
	}
}

func TestToBoolean(t *testing.T) {
	tests := []struct {
		data  interface{}
		value bool
		error bool
	}{
		{"", false, false},
		{false, false, false},
		{"false", false, false},
		{0, false, false},
		{[]string{}, false, true},
		{"whatever", false, false},
		{true, true, false},
		{"TRue", true, false},
		{3, true, false},
		{[]string{"whatever"}, true, true},
	}

	for i, test := range tests {
		v, err := toBoolean(test.data)
		if err != nil {
			if !test.error {
				t.Errorf("Test %d : unexpected error at boolean evaluation : %s", i, err)
			}
			continue
		}
		if test.error {
			t.Errorf("Test %d : no error at boolean evaluation, when one was expected", i)
			continue
		}

		if v != test.value {
			t.Errorf("Test %d : value return is not the one expected - expected : %v, got : %v", i, test.value, v)
		}

	}
}

func TestRuleEvaluate(t *testing.T) {
	tests := []struct {
		expression string
		value      bool
		errorExec  bool
	}{
		{"true", true, false},
		{"type == 'HINFO'", true, false},
		{"type == 'AAAA'", false, false},
		{"name =~ 'org'", true, false},
	}
	for i, test := range tests {

		engine := &ExprEngine{TypeDrop, rqdata.NewMapping("-")}
		rule, err := engine.BuildRule(append([]string{NameTypes[TypeAllow]}, strings.Split(test.expression, " ")...))
		if err != nil {
			t.Errorf("Test %d, expr : %s - unexpected error at build rule : %s", i, test.expression, err)
			continue
		}

		ctx := context.TODO()

		// build a Request
		r := new(dns.Msg)
		r.SetQuestion("example.org.", dns.TypeHINFO)
		r.MsgHdr.AuthenticatedData = true
		w := dnstest.NewRecorder(&tst.ResponseWriter{})
		state := request.Request{Req: r, W: w}

		data, err := engine.BuildQueryData(ctx, state)
		result, err := rule.Evaluate(data)
		if err != nil {
			if !test.errorExec {
				t.Errorf("Test %d, expr : %s - unexpected error at evaluate  : %s", i, test.expression, err)
			}
			continue
		}

		if (result == TypeAllow) != test.value {
			t.Errorf("Test %d, expr : %v -  value return is not the one expected - expected : %v, got : %v", i, test.expression, test.value, (result == TypeAllow))
		}

	}
}

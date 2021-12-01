package policy

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	tst "github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"

	"github.com/coredns/policy/plugin/pkg/response"
	"github.com/coredns/policy/plugin/pkg/rqdata"

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
		{"atoi('4') == 4.0", true, false},
		{"incidr('1.2.3.4','1.2.3.0/24')", true, false},
		{"incidr('1:2:3:4::1','1:2:3:4::/32')", true, false},
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
		w := response.NewReader(&tst.ResponseWriter{})
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

func TestAtoi(t *testing.T) {
	tests := []struct {
		args        []interface{}
		expected    interface{}
		expectedErr error
	}{
		{
			args:     []interface{}{"42"},
			expected: float64(42),
		},
		{
			args:        []interface{}{"42", "100"},
			expectedErr: fmt.Errorf("atoi requires exactly one string argument"),
		},
		{
			args:        []interface{}{},
			expectedErr: fmt.Errorf("atoi requires exactly one string argument"),
		},
		{
			args:        []interface{}{42},
			expectedErr: fmt.Errorf("atoi requires exactly one string argument"),
		},
		{
			args:        []interface{}{"foo"},
			expectedErr: fmt.Errorf("strconv.Atoi: parsing \"foo\": invalid syntax"),
		},
	}
	for i, test := range tests {
		v, err := atoi(test.args...)
		if test.expectedErr != nil {
			if err == nil {
				t.Errorf("Test %d, args : %v - expected error - expected : %v, got : %v", i, test.args, test.expectedErr, nil)
			} else if err.Error() != test.expectedErr.Error() {
				t.Errorf("Test %d, args : %v - expected error - expected : %v, got : %v", i, test.args, test.expectedErr, err)
			}
			continue
		}

		if !reflect.DeepEqual(v, test.expected) {
			t.Errorf("Test %d, args : %v -  value return is not the one expected - expected : %v, got : %v", i, test.args, test.expected, v)
		}
	}
}

func TestInCidr(t *testing.T) {
	tests := []struct {
		args        []interface{}
		expected    interface{}
		expectedErr error
	}{
		{
			args:     []interface{}{"1.2.3.4", "1.2.3.0/24"},
			expected: true,
		},
		{
			args:     []interface{}{"1.2.3.4", "5.6.7.0/24"},
			expected: false,
		},
		{
			args:     []interface{}{"1:2:3:4::1", "1:2:3:4::/32"},
			expected: true,
		},
		{
			args:     []interface{}{"1:2:3:4::1", "5:6:7:8::/32"},
			expected: false,
		},
		{
			args:        []interface{}{"1.2.3.4"},
			expectedErr: fmt.Errorf("invalid number of arguments"),
		},
		{
			args:        []interface{}{"foo", "5.6.7.0/24"},
			expectedErr: fmt.Errorf("first argument is not an IP address"),
		},
	}
	for i, test := range tests {
		v, err := incidr(test.args...)
		if test.expectedErr != nil {
			if err == nil {
				t.Errorf("Test %d, args : %v - expected error - expected : %v, got : %v", i, test.args, test.expectedErr, nil)
			} else if err.Error() != test.expectedErr.Error() {
				t.Errorf("Test %d, args : %v - expected error - expected : %v, got : %v", i, test.args, test.expectedErr, err)
			}
			continue
		}

		if !reflect.DeepEqual(v, test.expected) {
			t.Errorf("Test %d, args : %v -  value return is not the one expected - expected : %v, got : %v", i, test.args, test.expected, v)
		}
	}
}
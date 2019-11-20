package rqdata

import (
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/pkg/response"

	"github.com/miekg/dns"
)

func buildExtractorOnSimpleMsg(mapping *Mapping) *Extractor {

	w := response.NewReader(&test.ResponseWriter{})

	r := new(dns.Msg)
	r.SetQuestion("example.org.", dns.TypeHINFO)
	r.MsgHdr.AuthenticatedData = true
	state := request.Request{Req: r, W: w}

	return &Extractor{state, mapping}
}

func TestNewRequestData(t *testing.T) {

	mapping := NewMapping("")
	extractFromQuery := buildExtractorOnSimpleMsg(mapping)
	tests := []struct {
		extractor *Extractor
		name      string
		value     string
		subValue  string
		error     bool
	}{
		{extractFromQuery, "type", "HINFO", "", false},
		{extractFromQuery, "name", "example.org.", "", false},
		{extractFromQuery, "size", "29", "", false},
		{extractFromQuery, "invalid", "", "", true},
	}

	for i, tst := range tests {
		d, ok := tst.extractor.Value(tst.name)
		if !ok {
			if !tst.error {
				t.Errorf("Test %d, name : %s : unexpected invalid name returned", i, tst.name)
			}
			continue
		}
		if tst.error {
			t.Errorf("Test %d, name : %s : unexpected valid name returned with value %s", i, tst.name, tst.value)
		}
		if len(tst.subValue) > 0 {
			if !strings.Contains(d, tst.subValue) {
				t.Errorf("Test %d, name %s : valued returned : %s, expected to include : %s", i, tst.name, d, tst.subValue)
			}
			continue
		}
		if d != tst.value {
			t.Errorf("Test %d, name %s : valued returned : %s, expected : %s", i, tst.name, d, tst.value)
		}
	}
}

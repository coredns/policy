package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/dnstap/msg"
	"github.com/coredns/coredns/plugin/dnstap/test"
	mwtest "github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"

	tap "github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
)

func testCase(t *testing.T, ex Exchanger, q, r *dns.Msg, datq, datr *msg.Builder) {
	tapq, _ := datq.ToOutsideQuery(tap.Message_FORWARDER_QUERY)
	tapr, _ := datr.ToOutsideResponse(tap.Message_FORWARDER_RESPONSE)
	ctx := test.Context{}
	err := toDnstap(&ctx, "10.240.0.1:40212", ex,
		request.Request{W: &mwtest.ResponseWriter{}, Req: q}, r, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(ctx.Trap) != 2 {
		t.Fatalf("Messages: %d", len(ctx.Trap))
	}
	if !test.MsgEqual(ctx.Trap[0], tapq) {
		t.Errorf("Want: %v\nhave: %v", tapq, ctx.Trap[0])
	}
	if !test.MsgEqual(ctx.Trap[1], tapr) {
		t.Errorf("Want: %v\nhave: %v", tapr, ctx.Trap[1])
	}
}

func TestDnstap(t *testing.T) {
	q := mwtest.Case{Qname: "example.org", Qtype: dns.TypeA}.Msg()
	r := mwtest.Case{
		Qname: "example.org.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			mwtest.A("example.org. 3600	IN	A 10.0.0.1"),
		},
	}.Msg()
	tapq, tapr := test.TestingData(), test.TestingData()
	testCase(t, newDNSEx(), q, r, tapq, tapr)
	tapq.SocketProto = tap.SocketProtocol_TCP
	tapr.SocketProto = tap.SocketProtocol_TCP
	testCase(t, newDNSExWithOption(Options{ForceTCP: true}), q, r, tapq, tapr)
}

func TestNoDnstap(t *testing.T) {
	err := toDnstap(context.TODO(), "", nil, request.Request{}, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
}

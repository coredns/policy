package policy

import (
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/dnstap"
	"github.com/coredns/coredns/plugin/dnstap/taprw"
	dtest "github.com/coredns/coredns/plugin/dnstap/test"
	"github.com/coredns/coredns/plugin/test"
	tap "github.com/dnstap/golang-dnstap"
	"github.com/golang/protobuf/proto"
	pb "github.com/coredns/policy/dnstap"
	pdp "github.com/infobloxopen/themis/pdp-service"
	"github.com/miekg/dns"
	context "golang.org/x/net/context"
)

type testIORoutine struct {
	dnstapChan chan tap.Dnstap
}

func newIORoutine(timeout time.Duration) testIORoutine {
	ch := make(chan tap.Dnstap, 1)
	tapIO := testIORoutine{dnstapChan: ch}
	// close channel by timeout to prevent checker from waiting forever
	go func() {
		time.Sleep(timeout)
		close(ch)
	}()
	return tapIO
}

func (tapIO testIORoutine) Dnstap(msg tap.Dnstap) {
	tapIO.dnstapChan <- msg
}

func TestSendCRExtraNoMsg(t *testing.T) {
	trapper := dtest.TrapTapper{Full: true}
	tapRW := &taprw.ResponseWriter{
		Query:          new(dns.Msg),
		ResponseWriter: &test.ResponseWriter{},
		Tapper:         &trapper,
	}

	io := newIORoutine(100 * time.Millisecond)
	tapIO := newPolicyDnstapSender(io)
	tapIO.sendCRExtraMsg(tapRW, nil, nil)
	_, ok := <-io.dnstapChan
	if ok {
		t.Errorf("Unexpected msg received")
		return
	}
}

func TestSendCRExtraInvalidMsg(t *testing.T) {
	msg := dns.Msg{}
	msg.SetQuestion("test.com.", dns.TypeA)
	msg.Answer = []dns.RR{
		test.A("test.com.		600	IN	A			10.240.0.1"),
	}
	msg.Rcode = -1

	trapper := dtest.TrapTapper{Full: true}
	tapRW := &taprw.ResponseWriter{
		Query:          new(dns.Msg),
		ResponseWriter: &test.ResponseWriter{},
		Tapper:         &trapper,
	}
	tapRW.WriteMsg(&msg)

	io := newIORoutine(100 * time.Millisecond)
	tapIO := newPolicyDnstapSender(io)
	tapIO.sendCRExtraMsg(tapRW, &msg, nil)
	_, ok := <-io.dnstapChan
	if ok {
		t.Errorf("Unexpected msg received")
		return
	}
}

func TestSendCRExtraMsg(t *testing.T) {
	msg := dns.Msg{}
	msg.SetQuestion("test.com.", dns.TypeA)
	msg.Answer = []dns.RR{
		test.A("test.com.		600	IN	A			10.240.0.1"),
	}

	trapper := dtest.TrapTapper{Full: true}
	tapRW := &taprw.ResponseWriter{
		Query:          new(dns.Msg),
		ResponseWriter: &test.ResponseWriter{},
		Tapper:         &trapper,
		Send:           &taprw.SendOption{Cq: false, Cr: false},
	}
	tapRW.WriteMsg(&msg)

	io := newIORoutine(5000 * time.Millisecond)
	tapIO := newPolicyDnstapSender(io)

	testAttrHolder := &attrHolder{
		attrsReqDomain: []*pdp.Attribute{
			{Id: attrNameType, Value: typeValueQuery},
			{Id: attrNameDomainName, Value: "test.com"},
			{Id: attrNameDNSQtype, Value: "1"},
			{Id: attrNameSourceIP, Value: "10.0.0.7"},
			{Id: "option", Value: "option"},
		},
		attrsDnstap: []*pdp.Attribute{
			{Id: "dnstap", Value: "val"},
		},
		attrsEdnsStart: 4,
		action:         2,
	}

	tapIO.sendCRExtraMsg(tapRW, &msg, testAttrHolder)

	expectedAttrs := []*pdp.Attribute{
		{Id: "option", Value: "option"},
		{Id: "dnstap", Value: "val"},
	}
	checkCRExtraResult(t, io, &msg, expectedAttrs)

	if l := len(trapper.Trap); l != 0 {
		t.Errorf("Dnstap unexpectedly sent %d messages", l)
		return
	}

	testAttrHolder.action = 4

	tapIO.sendCRExtraMsg(tapRW, &msg, testAttrHolder)

	expectedAttrs = []*pdp.Attribute{
		{Id: attrNameDomainName, Value: "test.com"},
		{Id: attrNameDNSQtype, Value: "1"},
		{Id: attrNameSourceIP, Value: "10.0.0.7"},
		{Id: attrNamePolicyAction, Value: "3"},
		{Id: "option", Value: "option"},
		{Id: attrNameType, Value: typeValueQuery},
	}
	checkCRExtraResult(t, io, &msg, expectedAttrs)

	if l := len(trapper.Trap); l != 0 {
		t.Errorf("Dnstap unexpectedly sent %d messages", l)
		return
	}
}

func checkCRExtraResult(t *testing.T, io testIORoutine, orgMsg *dns.Msg, attrs []*pdp.Attribute) {
	dnstapMsg, ok := <-io.dnstapChan
	if !ok {
		t.Errorf("Receiving Dnstap message was timed out")
		return
	}
	extra := &pb.Extra{}
	err := proto.Unmarshal(dnstapMsg.Extra, extra)
	if err != nil {
		t.Errorf("Failed to unmarshal Extra (%v)", err)
		return
	}

	checkExtraAttrs(t, extra.GetAttrs(), attrs)
	checkCRMessage(t, dnstapMsg.Message, orgMsg)
}

func checkExtraAttrs(t *testing.T, actual []*pb.DnstapAttribute, expected []*pdp.Attribute) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %d attributes, found %d", len(expected), len(actual))
	}

checkAttr:
	for _, a := range actual {
		for _, e := range expected {
			if e.Id == a.Id {
				if a.Value != e.Value {
					t.Errorf("Attribute %s: expected %v , found %v", e.Id, e, a)
					return
				}
				continue checkAttr
			}
		}
		t.Errorf("Unexpected attribute found %v", a)
	}
}

func checkCRMessage(t *testing.T, msg *tap.Message, orgMsg *dns.Msg) {
	if msg == nil {
		t.Errorf("CR message not found")
		return
	}

	d := dtest.TestingData()
	bin, err := orgMsg.Pack()
	if err != nil {
		t.Errorf("Failed to pack message (%v)", err)
		return
	}
	d.Packed = bin
	expMsg, _ := d.ToClientResponse()
	if !dtest.MsgEqual(expMsg, msg) {
		t.Errorf("Unexpected message: expected: %v\nactual: %v", expMsg, msg)
	}
}

func TestRestCqCr(t *testing.T) {
	so := &taprw.SendOption{Cq: true, Cr: true}
	ctx := context.WithValue(context.Background(), dnstap.DnstapSendOption, so)
	resetCqCr(ctx)
	if so.Cq || so.Cr {
		t.Errorf("Failed to reset Cq/Cr flags")
	}
}

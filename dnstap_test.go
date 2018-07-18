package policy

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/dnstap"
	"github.com/coredns/coredns/plugin/dnstap/taprw"
	dtest "github.com/coredns/coredns/plugin/dnstap/test"
	"github.com/coredns/coredns/plugin/test"
	tap "github.com/dnstap/golang-dnstap"
	"github.com/golang/protobuf/proto"
	"github.com/infobloxopen/themis/pdp"
	"github.com/miekg/dns"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/context"

	pb "github.com/infobloxopen/themis/contrib/coredns/policy/dnstap"
)

func TestSendCRExtraNoMsg(t *testing.T) {
	ok := false
	g := newLogGrabber()
	defer func() {
		logs := g.Release()
		if ok {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}()

	trapper := dtest.TrapTapper{Full: true}
	tapRW := &taprw.ResponseWriter{
		Query:          new(dns.Msg),
		ResponseWriter: &test.ResponseWriter{},
		Tapper:         &trapper,
	}

	io := newIORoutine(100 * time.Millisecond)
	tapIO := newPolicyDnstapSender(io)
	tapIO.sendCRExtraMsg(tapRW, nil, nil)
	_, ok = <-io.dnstapChan
	if ok {
		t.Errorf("Unexpected msg received")
	}
}

func TestSendCRExtraInvalidMsg(t *testing.T) {
	ok := false
	g := newLogGrabber()
	defer func() {
		logs := g.Release()
		if ok {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}()

	msg := dns.Msg{}
	msg.SetQuestion("test.com.", dns.TypeA)
	msg.Answer = []dns.RR{
		test.A("test.com.       600 IN  A           10.240.0.1"),
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
	_, ok = <-io.dnstapChan
	if ok {
		t.Errorf("Unexpected msg received")
	}
}

func TestSendCRExtraMsg(t *testing.T) {
	ok := true
	g := newLogGrabber()
	defer func() {
		logs := g.Release()
		if !ok {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}()

	msg := dns.Msg{}
	msg.SetQuestion("test.com.", dns.TypeA)
	msg.Answer = []dns.RR{
		test.A("test.com.       600 IN  A           10.240.0.1"),
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
		dnReq: []pdp.AttributeAssignment{
			pdp.MakeStringAssignment(attrNameType, typeValueQuery),
			pdp.MakeDomainAssignment(attrNameDomainName, makeTestDomain("test.com")),
			pdp.MakeStringAssignment(attrNameDNSQtype, "1"),
			pdp.MakeAddressAssignment(attrNameSourceIP, net.ParseIP("10.0.0.7")),
			pdp.MakeStringAssignment("option", "option"),
		},
		dnstap: []pdp.AttributeAssignment{
			pdp.MakeStringAssignment("dnstap", "val"),
		},
		action: actionAllow,
	}

	tapIO.sendCRExtraMsg(tapRW, &msg, testAttrHolder)

	ok = assertCRExtraResult(t, "sendCRExtraMsg(actionAllow)", io, &msg,
		&pb.DnstapAttribute{Id: attrNameSourceIP, Value: "10.0.0.7"},
		&pb.DnstapAttribute{Id: "option", Value: "option"},
		&pb.DnstapAttribute{Id: "dnstap", Value: "val"},
	)

	if l := len(trapper.Trap); l != 0 {
		t.Fatalf("Dnstap unexpectedly sent %d messages", l)
		ok = false
	}

	testAttrHolder.action = actionBlock

	tapIO.sendCRExtraMsg(tapRW, &msg, testAttrHolder)

	ok = assertCRExtraResult(t, "sendCRExtraMsg(actionBlock)", io, &msg,
		&pb.DnstapAttribute{Id: attrNameDomainName, Value: "test.com"},
		&pb.DnstapAttribute{Id: attrNameDNSQtype, Value: "1"},
		&pb.DnstapAttribute{Id: attrNameSourceIP, Value: "10.0.0.7"},
		&pb.DnstapAttribute{Id: "option", Value: "option"},
		&pb.DnstapAttribute{Id: attrNamePolicyAction, Value: "3"},
		&pb.DnstapAttribute{Id: attrNameType, Value: typeValueQuery},
	) && ok

	if l := len(trapper.Trap); l != 0 {
		t.Fatalf("Dnstap unexpectedly sent %d messages", l)
		ok = false
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

func assertCRExtraResult(t *testing.T, desc string, io testIORoutine, eMsg *dns.Msg, e ...*pb.DnstapAttribute) bool {
	dnstapMsg, ok := <-io.dnstapChan
	if !ok {
		t.Errorf("Receiving Dnstap for %q message was timed out", desc)
		return false
	}

	extra := &pb.Extra{}
	err := proto.Unmarshal(dnstapMsg.Extra, extra)
	if err != nil {
		t.Errorf("Failed to unmarshal Extra for %q (%v)", desc, err)
		return false
	}

	ok = assertDnstapAttributes(t, desc, extra.GetAttrs(), e...)
	ok = assertCRMessage(t, desc, dnstapMsg.Message, eMsg) && ok
	return ok
}

func assertCRMessage(t *testing.T, desc string, msg *tap.Message, e *dns.Msg) bool {
	if msg == nil {
		t.Errorf("CR message for %q not found", desc)
		return false
	}

	bin, err := e.Pack()
	if err != nil {
		t.Errorf("Failed to pack message for %q (%v)", desc, err)
		return false
	}

	d := dtest.TestingData()
	d.Packed = bin
	eMsg, _ := d.ToClientResponse()
	if !dtest.MsgEqual(eMsg, msg) {
		t.Errorf("Unexpected message for %q: expected:\n%v\nactual:\n%v", desc, eMsg, msg)
		return false
	}

	return true
}

func serializeDnstapAttributesForAssert(a []*pb.DnstapAttribute) []string {
	out := make([]string, len(a))
	for i, a := range a {
		out[i] = fmt.Sprintf("%q = %q\n", a.Id, a.Value)
	}

	return out
}

func assertDnstapAttributes(t *testing.T, desc string, a []*pb.DnstapAttribute, e ...*pb.DnstapAttribute) bool {
	ctx := difflib.ContextDiff{
		A:        serializeDnstapAttributesForAssert(a),
		B:        serializeDnstapAttributesForAssert(e),
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("can't compare \"%s\": %s", desc, err))
	}

	if len(diff) > 0 {
		t.Errorf("\"%s\" doesn't match:\n%s", desc, diff)
		return false
	}

	return true
}

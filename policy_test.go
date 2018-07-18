package policy

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/dnstap/taprw"
	dtest "github.com/coredns/coredns/plugin/dnstap/test"
	"github.com/infobloxopen/themis/pep"
	"github.com/miekg/dns"
	"golang.org/x/net/context"

	pb "github.com/infobloxopen/themis/contrib/coredns/policy/dnstap"
)

func TestNewPolicyPlugin(t *testing.T) {
	p := newPolicyPlugin()
	if p == nil {
		t.Error("can't create new policy plugin instance")
	}
}

func TestPolicyPluginName(t *testing.T) {
	p := newPolicyPlugin()

	n := p.Name()
	if n != "policy" {
		t.Errorf("expected %q as plugin name but got %q", "policy", n)
	}
}

func TestPolicyPluginServeDNS(t *testing.T) {
	endpoint := "127.0.0.1:5555"
	srv := startPDPServer(t, serveDNSTestPolicy, endpoint)
	defer func() {
		if logs := srv.Stop(); len(logs) > 0 {
			t.Logf("server logs:\n%s", logs)
		}
	}()

	if err := waitForPortOpened(endpoint); err != nil {
		t.Fatalf("can't connect to PDP server: %s", err)
	}

	p := newPolicyPlugin()
	p.conf.endpoints = []string{endpoint}
	p.conf.connTimeout = time.Second
	p.conf.streams = 1
	p.conf.log = true
	p.conf.debugID = "<DEBUG>"
	p.conf.debugSuffix = "debug.local."

	mp := &mockPlugin{
		ip: net.ParseIP("192.0.2.53"),
		rc: dns.RcodeSuccess,
	}
	p.next = mp

	g := newLogGrabber()
	if err := p.connect(); err != nil {
		logs := g.Release()
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		t.Fatal(err)
	}
	defer p.closeConn()

	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err := p.ServeDNS(context.TODO(), w, m)
	logs := g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr aa; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.\t0\tIN\tA\t192.0.2.53\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.com.debug.local", dns.TypeTXT, dns.ClassCHAOS)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(debug)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.debug.local.\tCH\t TXT\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t\"resolve:yes,query:'allow',ident:'<DEBUG>'\"\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.redirect", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(domain redirect)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.redirect.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.redirect.\t0\tIN\tA\t192.0.2.254\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	mp.ip = net.ParseIP("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(address redirect)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.\t0\tIN\tA\t192.0.2.253\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.ip = net.ParseIP("192.0.2.53")

	m = makeTestDNSMsg("example.block", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(domain block)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NXDOMAIN, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.block.\tIN\t A\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	mp.ip = net.ParseIP("192.0.2.17")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(address block)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NXDOMAIN, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.ip = net.ParseIP("192.0.2.53")

	m = makeTestDNSMsg("example.refuse", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(domain refuse)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: REFUSED, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.refuse.\tIN\t A\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	mp.ip = net.ParseIP("192.0.2.33")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(address refuse)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: REFUSED, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.ip = net.ParseIP("192.0.2.53")

	m = makeTestDNSMsg("example.drop", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(domain drop)", rc, w.Msg, dns.RcodeSuccess,
			"<nil> MsgHdr",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	mp.ip = net.ParseIP("192.0.2.65")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(address drop)", rc, w.Msg, dns.RcodeSuccess,
			"<nil> MsgHdr",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.ip = net.ParseIP("192.0.2.53")

	m = makeTestDNSMsg("example.missing", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err == nil {
		t.Errorf("expected errInvalidAction but got rc: %d, msg:\n%q", rc, w.Msg)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else if err != errInvalidAction {
		t.Errorf("exepcted errInvalidAction but got %T: %s", err, err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	mp.ip = net.ParseIP("192.0.2.81")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err == nil {
		t.Errorf("expected errInvalidAction but got rc: %d, msg:\n%q", rc, w.Msg)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else if err != errInvalidAction {
		t.Errorf("exepcted errInvalidAction but got %T: %s", err, err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}

	mp.err = fmt.Errorf("test next plugin error")

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != mp.err {
		t.Errorf("expected %q but got %q", mp.err, err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}
	if !assertDNSMessage(t, "ServeDNS(next plugin error)", rc, w.Msg, dns.RcodeSuccess,
		";; opcode: QUERY, status: SERVFAIL, id: 0\n"+
			";; flags: qr; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n;example.com.\tIN\t A\n",
	) {
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}

	m = makeTestDNSMsg("example.com.debug.local", dns.TypeTXT, dns.ClassCHAOS)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(next plugin error with debug)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.debug.local.\tCH\t TXT\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.debug.local.\t0\tCH\tTXT\t\"resolve:failed,query:'allow',ident:'<DEBUG>'\"\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.err = nil
	mp.ip = nil

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(dropped by resolver)", rc, w.Msg, dns.RcodeSuccess,
			"<nil> MsgHdr",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.ip = net.ParseIP("192.0.2.53")
	mp.rc = dns.RcodeServerFailure

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(resolver failed)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: SERVFAIL, id: 0\n"+
				";; flags: qr aa; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	mp.rc = dns.RcodeSuccess

	client := p.pdp

	dnErr := fmt.Errorf("test error on domain validation")
	errPep := newErraticPep(client, dnErr, nil)
	p.pdp = errPep

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != dnErr {
		t.Errorf("expected %q but got %q", dnErr, err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}
	if !assertDNSMessage(t, "ServeDNS(error on domain validation)", rc, w.Msg, dns.RcodeSuccess,
		"<nil> MsgHdr",
	) {
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}

	ipErr := fmt.Errorf("test error on address validation")
	errPep = newErraticPep(client, nil, ipErr)
	p.pdp = errPep

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != ipErr {
		t.Errorf("expected %q but got %q", ipErr, err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}
	if !assertDNSMessage(t, "ServeDNS(error on domain validation)", rc, w.Msg, dns.RcodeSuccess,
		"<nil> MsgHdr",
	) {
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	}
}

func TestPolicyPluginServeDNSPassthrough(t *testing.T) {
	endpoint := "127.0.0.1:5555"
	if err := waitForPortClosed(endpoint); err != nil {
		t.Fatalf("port still in use: %s", err)
	}

	p := newPolicyPlugin()

	p.conf.endpoints = []string{endpoint}
	p.conf.connTimeout = time.Second
	p.conf.streams = 1
	p.conf.log = true
	p.conf.passthrough = []string{"passthrough.local."}
	p.conf.debugID = "<DEBUG>"
	p.conf.debugSuffix = "debug.local."

	p.next = &mockPlugin{
		ip: net.ParseIP("192.0.2.53"),
		rc: dns.RcodeSuccess,
	}

	g := newLogGrabber()
	if err := p.connect(); err != nil {
		logs := g.Release()
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		t.Fatal(err)
	}
	defer p.closeConn()

	m := makeTestDNSMsg("example.passthrough.local", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err := p.ServeDNS(context.TODO(), w, m)
	logs := g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(passthrough)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr aa; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.passthrough.local.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.passthrough.local.\t0\tIN\tA\t192.0.2.53\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}

	m = makeTestDNSMsg("example.passthrough.local.debug.local.", dns.TypeTXT, dns.ClassCHAOS)
	w = newTestAddressedNonwriter("192.0.2.1")

	g = newLogGrabber()
	rc, err = p.ServeDNS(context.TODO(), w, m)
	logs = g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS(passthrough+debug)", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr aa; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.passthrough.local.debug.local.\tCH\t TXT\n\n"+
				";; ANSWER SECTION:\n"+
				"example.passthrough.local.debug.local.\t0\tCH\tTXT\t\"action:passthrough\"\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}
}

func TestPolicyPluginServeDNSWithDnstap(t *testing.T) {
	endpoint := "127.0.0.1:5555"
	srv := startPDPServer(t, serveDNSTestPolicy, endpoint)
	defer func() {
		if logs := srv.Stop(); len(logs) > 0 {
			t.Logf("server logs:\n%s", logs)
		}
	}()

	if err := waitForPortOpened(endpoint); err != nil {
		t.Fatalf("can't connect to PDP server: %s", err)
	}

	p := newPolicyPlugin()
	p.conf.endpoints = []string{endpoint}
	p.conf.connTimeout = time.Second
	p.conf.streams = 1
	p.conf.log = true
	p.conf.debugID = "<DEBUG>"
	p.conf.debugSuffix = "debug.local."

	mp := &mockPlugin{
		ip: net.ParseIP("192.0.2.53"),
		rc: dns.RcodeSuccess,
	}
	p.next = mp

	io := newIORoutine(5000 * time.Millisecond)
	p.tapIO = newPolicyDnstapSender(io)

	g := newLogGrabber()
	if err := p.connect(); err != nil {
		logs := g.Release()
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		t.Fatal(err)
	}
	defer p.closeConn()

	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriterWithAddr(&net.UDPAddr{
		IP:   net.ParseIP("10.240.0.1"),
		Port: 40212,
		Zone: "",
	})

	tapRW := &taprw.ResponseWriter{
		Query:          new(dns.Msg),
		ResponseWriter: w,
		Tapper:         &dtest.TrapTapper{Full: true},
		Send:           &taprw.SendOption{Cq: false, Cr: false},
	}

	g = newLogGrabber()
	rc, err := p.ServeDNS(context.TODO(), tapRW, m)
	logs := g.Release()
	if err != nil {
		t.Error(err)
		t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
	} else {
		if !assertDNSMessage(t, "ServeDNS", rc, w.Msg, dns.RcodeSuccess,
			";; opcode: QUERY, status: NOERROR, id: 0\n"+
				";; flags: qr aa; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
				";; QUESTION SECTION:\n"+
				";example.com.\tIN\t A\n\n"+
				";; ANSWER SECTION:\n"+
				"example.com.\t0\tIN\tA\t192.0.2.53\n",
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}

		if !assertCRExtraResult(t, "sendCRExtraMsg(actionAllow)", io, w.Msg,
			&pb.DnstapAttribute{Id: attrNameSourceIP, Value: "10.240.0.1"},
		) {
			t.Logf("=== plugin logs ===\n%s--- plugin logs ---", logs)
		}
	}
}

const serveDNSTestPolicy = `# All Permit Policy
attributes:
  type: string
  domain_name: domain
  address: address
  redirect_to: string
  refuse: string
  drop: string
  missing: string
policies:
  alg: FirstApplicableEffect
  policies:
  - id: "Query rules"
    target:
    - equal:
      - attr: type
      - val:
          type: string
          content: query
    alg: FirstApplicableEffect
    rules:
    - id: "Permit example.com"
      target:
      - contains:
        - val:
            type: set of domains
            content:
            - example.com
        - attr: domain_name
      effect: Permit
    - id: "Redirect example.redirect"
      target:
      - contains:
        - val:
            type: set of domains
            content:
            - example.redirect
        - attr: domain_name
      effect: Deny
      obligations:
      - redirect_to:
          val:
            type: string
            content: "192.0.2.254"
    - id: "Block example.block"
      target:
      - contains:
        - val:
            type: set of domains
            content:
            - example.block
        - attr: domain_name
      effect: Deny
    - id: "Refuse example.refuse"
      target:
      - contains:
        - val:
            type: set of domains
            content:
            - example.refuse
        - attr: domain_name
      effect: Deny
      obligations:
      - refuse:
          val:
            type: string
            content: ""
    - id: "Drop example.drop"
      target:
      - contains:
        - val:
            type: set of domains
            content:
            - example.drop
        - attr: domain_name
      effect: Deny
      obligations:
      - drop:
          val:
            type: string
            content: ""
    - id: "Missing attribute example.missing"
      target:
      - contains:
        - val:
            type: set of domains
            content:
            - example.missing
        - attr: domain_name
      condition:
        equal:
        - attr: missing
        - val:
            type: string
            content: missing
      effect: Permit
  - id: "Response rules"
    target:
    - equal:
      - attr: type
      - val:
          type: string
          content: response
    alg: FirstApplicableEffect
    rules:
    - id: "Permit 192.0.2.48/28"
      target:
      - contains:
        - val:
            type: set of networks
            content:
            - 192.0.2.48/28
        - attr: address
      effect: Permit
    - id: "Redirect 192.0.2.0/28"
      target:
      - contains:
        - val:
            type: set of networks
            content:
            - 192.0.2.0/28
        - attr: address
      effect: Deny
      obligations:
      - redirect_to:
          val:
            type: string
            content: "192.0.2.253"
    - id: "Block 192.0.2.16/28"
      target:
      - contains:
        - val:
            type: set of networks
            content:
            - 192.0.2.16/28
        - attr: address
      effect: Deny
    - id: "Refuse 192.0.2.32/28"
      target:
      - contains:
        - val:
            type: set of networks
            content:
            - 192.0.2.32/28
        - attr: address
      effect: Deny
      obligations:
      - refuse:
          val:
            type: string
            content: ""
    - id: "Drop 192.0.2.64/28"
      target:
      - contains:
        - val:
            type: set of networks
            content:
            - 192.0.2.64/28
        - attr: address
      effect: Deny
      obligations:
      - drop:
          val:
            type: string
            content: ""
    - id: "Missing attribute 192.0.2.80/28"
      target:
      - contains:
        - val:
            type: set of networks
            content:
            - 192.0.2.80/28
        - attr: address
      condition:
        equal:
        - attr: missing
        - val:
            type: string
            content: missing
      effect: Permit
`

func assertDNSMessage(t *testing.T, desc string, rc int, m *dns.Msg, erc int, eMsg string) bool {
	ok := true

	if rc != erc {
		t.Errorf("expected %d rcode for %q but got %d", erc, desc, rc)
		ok = false
	}

	if m.String() != eMsg {
		t.Errorf("expected response for %q:\n%q\nbut got:\n%q", desc, eMsg, m)
		ok = false
	}

	return ok
}

type logGrabber struct {
	b *bytes.Buffer
}

func newLogGrabber() *logGrabber {
	b := new(bytes.Buffer)
	log.SetOutput(b)

	return &logGrabber{
		b: b,
	}
}

func (g *logGrabber) Release() string {
	log.SetOutput(os.Stderr)

	return g.b.String()
}

const (
	mpModeConst = iota
	mpModeInc
	mpModeHalfInc
)

type mockPlugin struct {
	ip   net.IP
	err  error
	rc   int
	mode int
	cnt  *uint32
}

// Name implements the plugin.Handler interface.
func (p *mockPlugin) Name() string {
	return "mockPlugin"
}

// ServeDNS implements the plugin.Handler interface.
func (p *mockPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if p.err != nil {
		return dns.RcodeServerFailure, p.err
	}

	if r == nil || len(r.Question) <= 0 {
		return dns.RcodeServerFailure, nil
	}

	ip := p.ip
	if p.mode != mpModeConst && p.cnt != nil {
		i := atomic.AddUint32(p.cnt, 1)

		if p.mode != mpModeHalfInc || i&1 == 0 {
			ip = addToIP(ip, i)
		}
	}

	q := r.Question[0]
	hdr := dns.RR_Header{
		Name:   q.Name,
		Rrtype: q.Qtype,
		Class:  q.Qclass,
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		if q.Qtype != dns.TypeA {
			return dns.RcodeSuccess, nil
		}

		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true
		m.Rcode = p.rc

		if m.Rcode == dns.RcodeSuccess {
			m.Answer = append(m.Answer,
				&dns.A{
					Hdr: hdr,
					A:   ipv4,
				},
			)
		}

		w.WriteMsg(m)
	} else if ipv6 := ip.To16(); ipv6 != nil {
		if q.Qtype != dns.TypeAAAA {
			return dns.RcodeSuccess, nil
		}

		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true
		m.Rcode = p.rc

		if m.Rcode == dns.RcodeSuccess {
			m.Answer = append(m.Answer,
				&dns.AAAA{
					Hdr:  hdr,
					AAAA: ipv6,
				},
			)
		}

		w.WriteMsg(m)
	}

	return p.rc, nil
}

func addToIP(ip net.IP, n uint32) net.IP {
	if n == 0 {
		return ip
	}

	out := net.IP(make([]byte, len(ip)))
	copy(out, ip)

	d := uint(n % 256)
	n /= 256

	c := uint(n % 256)
	n /= 256

	b := uint(n % 256)
	n /= 256

	a := uint(n)

	t := uint(out[len(out)-1])
	t += d
	out[len(out)-1] = byte(t % 256)

	z := uint(out[len(out)-2])
	if t > 255 {
		z++
	}

	z += c
	out[len(out)-2] = byte(z % 256)

	y := uint(out[len(out)-3])
	if z > 255 {
		y++
	}

	y += b
	out[len(out)-3] = byte(y % 256)

	x := uint(out[len(out)-4])
	if y > 255 {
		x++
	}

	x += a
	out[len(out)-4] = byte(x % 256)

	return out
}

type erraticPep struct {
	counter int
	err     []error
	client  pep.Client
}

func newErraticPep(c pep.Client, err ...error) *erraticPep {
	return &erraticPep{
		err:    err,
		client: c,
	}
}

func (c *erraticPep) Connect(addr string) error {
	return c.client.Connect(addr)
}

func (c *erraticPep) Close() {
	c.client.Close()
}

func (c *erraticPep) Validate(in, out interface{}) error {
	if len(c.err) > 0 {
		n := c.counter % len(c.err)
		c.counter++

		err := c.err[n]
		if err != nil {
			return err
		}
	}

	return c.client.Validate(in, out)
}

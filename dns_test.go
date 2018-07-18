package policy

import (
	"errors"
	"net"
	"testing"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

func TestGetNameAndType(t *testing.T) {
	fqdn := dns.Fqdn("example.com")
	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)

	qName, qType := getNameAndType(m)
	if qName != fqdn {
		t.Errorf("expected %q as query name but got %q", fqdn, qName)
	}

	if qType != dns.TypeA {
		t.Errorf("expected %d as query type but got %d", dns.TypeA, qType)
	}

	fqdn = "."
	qName, qType = getNameAndType(nil)
	if qName != fqdn {
		t.Errorf("expected %q as query name but got %q", fqdn, qName)
	}

	if qType != dns.TypeNone {
		t.Errorf("expected %d as query type but got %d", dns.TypeNone, qType)
	}
}

func TestGetNameAndClass(t *testing.T) {
	fqdn := dns.Fqdn("example.com")
	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)

	qName, qClass := getNameAndClass(m)
	if qName != fqdn {
		t.Errorf("expected %q as query name but got %q", fqdn, qName)
	}

	if qClass != dns.ClassINET {
		t.Errorf("expected %d as query class but got %d", dns.ClassINET, qClass)
	}

	fqdn = "."
	qName, qClass = getNameAndClass(nil)
	if qName != fqdn {
		t.Errorf("expected %q as query name but got %q", fqdn, qName)
	}

	if qClass != dns.ClassNONE {
		t.Errorf("expected %d as query class but got %d", dns.ClassNONE, qClass)
	}
}

func TestGetRemoveIP(t *testing.T) {
	w := newTestAddressedNonwriter("192.0.2.1")
	a := getRemoteIP(w)
	if !a.Equal(net.ParseIP("192.0.2.1")) {
		t.Errorf("expected %s as remote address but got %s", "192.0.2.1", a)
	}

	w = newTestAddressedNonwriter("192.0.2.1:53")
	a = getRemoteIP(w)
	if !a.Equal(net.ParseIP("192.0.2.1")) {
		t.Errorf("expected %s as remote address but got %s", "192.0.2.1", a)
	}
}

func TestGetRespIp(t *testing.T) {
	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	appendAnswer(m, newA(net.ParseIP("192.0.2.1")))

	a := getRespIP(m)
	if !a.Equal(net.ParseIP("192.0.2.1")) {
		t.Errorf("expected %s as response address but got %s", "192.0.2.1", a)
	}

	m = makeTestDNSMsg("example.com", dns.TypeAAAA, dns.ClassINET)
	appendAnswer(m, newAAAA(net.ParseIP("2001:db8::1")))

	a = getRespIP(m)
	if !a.Equal(net.ParseIP("2001:db8::1")) {
		t.Errorf("expected %s as response address but got %s", "2001:db8::1", a)
	}

	m = makeTestDNSMsg("www.example.com", dns.TypeCNAME, dns.ClassINET)
	appendAnswer(m, newCNAME("example.com"))

	a = getRespIP(m)
	if a != nil {
		t.Errorf("expected no response address but got %s", a)
	}

	a = getRespIP(nil)
	if a != nil {
		t.Errorf("expected no response address but got %s", a)
	}
}

func TestExtractOptionsFromEDNS0(t *testing.T) {
	optsMap := map[uint16][]*edns0Opt{
		0xfffe: {
			{
				name:     "test",
				dataType: typeEDNS0Bytes,
				size:     4,
			},
		},
	}

	m := makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET,
		newEdns0(
			newEdns0Cookie("badc0de."),
			newEdns0Local(0xfffd, []byte{0xde, 0xc0, 0xad, 0xb}),
			newEdns0Local(0xfffe, []byte{0xef, 0xbe, 0xad, 0xde}),
		),
	)

	n := 0
	extractOptionsFromEDNS0(m, optsMap, func(b []byte, opts []*edns0Opt) {
		n++

		if string(b) != string([]byte{0xef, 0xbe, 0xad, 0xde}) {
			t.Errorf("expected [% x] as EDNS0 data for option %d but got [% x]", []byte{0xef, 0xbe, 0xad, 0xde}, n, b)
		}

		if len(opts) != 1 || opts[0].name != "test" {
			t.Errorf("expected %q ENDS0 for option %d but got %+v", "test", n, opts)
		}
	})

	if n != 1 {
		t.Errorf("expected exactly one EDNS0 option but got %d", n)
	}

	o := m.IsEdns0()
	if o == nil {
		t.Error("expected ENDS0 options in DNS message")
	} else if len(o.Option) != 2 {
		t.Errorf("expected exactly %d options remaining but got %d", 2, len(o.Option))
	}
}

func TestClearECS(t *testing.T) {
	m := makeTestDNSMsgWithEdns0("example.com", dns.TypeA, dns.ClassINET,
		newEdns0(
			newEdns0Cookie("badc0de."),
			newEdns0Subnet(net.ParseIP("192.0.2.1")),
			newEdns0Local(0xfffe, []byte{0xb, 0xad, 0xc0, 0xde}),
			newEdns0Subnet(net.ParseIP("2001:db8::1")),
		),
	)

	clearECS(m)
	assertDNSMessage(t, "clearECS", 0, m, 0,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1\n\n"+
			";; QUESTION SECTION:\n"+
			";example.com.\tIN\t A\n\n"+
			";; ADDITIONAL SECTION:\n\n"+
			";; OPT PSEUDOSECTION:\n"+
			"; EDNS: version 0; flags: ; udp: 0\n"+
			"; COOKIE: badc0de.\n"+
			"; LOCAL OPT: 65534:0x0badc0de\n",
	)
}

func TestSetRedirectQueryAnswer(t *testing.T) {
	p := newPolicyPlugin()

	mp := &mockPlugin{
		ip: net.ParseIP("192.0.2.153"),
		rc: dns.RcodeSuccess,
	}
	p.next = mp

	m := makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w := newTestAddressedNonwriter("192.0.2.1")

	rc, err := p.setRedirectQueryAnswer(context.TODO(), w, m, "192.0.2.53")
	if err != nil {
		t.Error(err)
	}
	assertDNSMessage(t, "setRedirectQueryAnswer(192.0.2.53)", rc, m, dns.RcodeSuccess,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n"+
			";example.com.\tIN\t A\n\n"+
			";; ANSWER SECTION:\n"+
			"example.com.\t0\tIN\tA\t192.0.2.53\n",
	)

	m = makeTestDNSMsg("example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	rc, err = p.setRedirectQueryAnswer(context.TODO(), w, m, "2001:db8::53")
	if err != nil {
		t.Error(err)
	}
	assertDNSMessage(t, "setRedirectQueryAnswer(2001:db8::53)", rc, m, dns.RcodeSuccess,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n"+
			";example.com.\tIN\t A\n\n"+
			";; ANSWER SECTION:\n"+
			"example.com.\t0\tIN\tAAAA\t2001:db8::53\n",
	)

	m = makeTestDNSMsg("redirect.example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	rc, err = p.setRedirectQueryAnswer(context.TODO(), w, m, "example.com")
	if err != nil {
		t.Error(err)
	}
	assertDNSMessage(t, "setRedirectQueryAnswer(redirect.example.com->example.com)", rc, m, dns.RcodeSuccess,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags: qr aa; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n"+
			";redirect.example.com.\tIN\t A\n\n"+
			";; ANSWER SECTION:\n"+
			"redirect.example.com.\t0\tIN\tCNAME\texample.com.\n"+
			"example.com.\t0\tIN\tA\t192.0.2.153\n",
	)

	m = new(dns.Msg)
	w = newTestAddressedNonwriter("192.0.2.1")

	rc, err = p.setRedirectQueryAnswer(context.TODO(), w, m, "example.com")
	if err == nil {
		t.Errorf("expected errInvalidDNSMessage")
	} else if err != errInvalidDNSMessage {
		t.Errorf("expected errInvalidDNSMessage but got %T: %s", err, err)
	}

	assertDNSMessage(t, "setRedirectQueryAnswer(empty)", rc, m, dns.RcodeServerFailure,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 0, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n",
	)

	m = makeTestDNSMsg("redirect.example.com", dns.TypeA, dns.ClassINET)
	w = newTestAddressedNonwriter("192.0.2.1")

	errTest := errors.New("testError")
	mp.err = errTest

	rc, err = p.setRedirectQueryAnswer(context.TODO(), w, m, "example.com")
	if err == nil {
		t.Errorf("expected errTest")
	} else if err != errTest {
		t.Errorf("expected errTest but got %T: %s", err, err)
	}

	assertDNSMessage(t, "setRedirectQueryAnswer(redirect.example.com->error)", rc, m, dns.RcodeServerFailure,
		";; opcode: QUERY, status: NOERROR, id: 0\n"+
			";; flags:; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0\n\n"+
			";; QUESTION SECTION:\n;"+
			"redirect.example.com.\tIN\t A\n",
	)
}

func makeTestDNSMsg(n string, t uint16, c uint16) *dns.Msg {
	out := new(dns.Msg)
	out.Question = make([]dns.Question, 1)
	out.Question[0] = dns.Question{
		Name:   dns.Fqdn(n),
		Qtype:  t,
		Qclass: c,
	}
	return out
}

func appendAnswer(m *dns.Msg, rr ...dns.RR) {
	if m.Answer == nil {
		m.Answer = []dns.RR{}
	}

	m.Answer = append(m.Answer, rr...)
}

func newA(a net.IP) dns.RR {
	out := new(dns.A)
	out.Hdr.Name = "."
	out.Hdr.Rrtype = dns.TypeA
	out.A = a

	return out
}

func newAAAA(a net.IP) dns.RR {
	out := new(dns.AAAA)
	out.Hdr.Name = "."
	out.Hdr.Rrtype = dns.TypeAAAA
	out.AAAA = a

	return out
}

func newCNAME(s string) dns.RR {
	out := new(dns.CNAME)
	out.Hdr.Name = "."
	out.Hdr.Rrtype = dns.TypeCNAME
	out.Target = dns.Fqdn(s)

	return out
}

func makeTestDNSMsgWithEdns0(n string, t uint16, c uint16, o ...*dns.OPT) *dns.Msg {
	out := makeTestDNSMsg(n, t, c)

	extra := make([]dns.RR, len(o))
	for i, o := range o {
		extra[i] = o
	}

	out.Extra = extra
	return out
}

func newEdns0(o ...dns.EDNS0) *dns.OPT {
	out := new(dns.OPT)
	out.Hdr.Name = "."
	out.Hdr.Rrtype = dns.TypeOPT
	out.Option = o

	return out
}

func copyEdns0(in ...*dns.OPT) []*dns.OPT {
	out := make([]*dns.OPT, len(in))
	for i, o := range in {
		out[i] = new(dns.OPT)
		out[i].Hdr = o.Hdr
		out[i].Option = make([]dns.EDNS0, len(o.Option))
		copy(out[i].Option, o.Option)
	}

	return out
}

func newEdns0Cookie(s string) dns.EDNS0 {
	out := new(dns.EDNS0_COOKIE)
	out.Code = dns.EDNS0COOKIE
	out.Cookie = s

	return out
}

func newEdns0Local(c uint16, b []byte) dns.EDNS0 {
	out := new(dns.EDNS0_LOCAL)
	out.Code = c
	out.Data = b

	return out
}

func newEdns0Subnet(ip net.IP) dns.EDNS0 {
	out := new(dns.EDNS0_SUBNET)
	out.Code = dns.EDNS0SUBNET
	if ipv4 := ip.To4(); ipv4 != nil {
		out.Family = 1
		out.SourceNetmask = 32
		out.Address = ipv4
	} else if ipv6 := ip.To16(); ipv6 != nil {
		out.Family = 2
		out.SourceNetmask = 128
		out.Address = ipv6
	}
	out.SourceScope = 0

	return out
}

type testAddressedNonwriter struct {
	dns.ResponseWriter
	ra  net.Addr
	Msg *dns.Msg
}

type testUDPAddr struct {
	addr string
}

func newTestAddressedNonwriter(ra string) *testAddressedNonwriter {
	return &testAddressedNonwriter{
		ResponseWriter: nil,
		ra:             newUDPAddr(ra),
	}
}

func newTestAddressedNonwriterWithAddr(ra net.Addr) *testAddressedNonwriter {
	return &testAddressedNonwriter{
		ResponseWriter: nil,
		ra:             ra,
	}
}

func (w *testAddressedNonwriter) RemoteAddr() net.Addr {
	return w.ra
}

func (w *testAddressedNonwriter) WriteMsg(res *dns.Msg) error {
	w.Msg = res
	return nil
}

func newUDPAddr(addr string) *testUDPAddr {
	return &testUDPAddr{
		addr: addr,
	}
}

func (a *testUDPAddr) String() string {
	return a.addr
}

func (a *testUDPAddr) Network() string {
	return "udp"
}

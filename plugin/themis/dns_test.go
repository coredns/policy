package themis

import (
	"net"

	"github.com/miekg/dns"
)

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

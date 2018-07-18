package policy

import (
	"errors"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var errInvalidDNSMessage = errors.New("invalid DNS message")

func getNameAndType(r *dns.Msg) (string, uint16) {
	if r == nil || len(r.Question) <= 0 {
		return ".", dns.TypeNone
	}

	q := r.Question[0]
	return q.Name, q.Qtype
}

func getNameAndClass(r *dns.Msg) (string, uint16) {
	if r == nil || len(r.Question) <= 0 {
		return ".", dns.ClassNONE
	}

	q := r.Question[0]
	return q.Name, q.Qclass
}

func getRemoteIP(w dns.ResponseWriter) net.IP {
	addrPort := w.RemoteAddr().String()
	addr, _, err := net.SplitHostPort(w.RemoteAddr().String())
	if err != nil {
		addr = addrPort
	}

	return net.ParseIP(addr)
}

func getRespIP(r *dns.Msg) net.IP {
	if r == nil {
		return nil
	}

	var ip net.IP
	for _, rr := range r.Answer {
		switch rr := rr.(type) {
		case *dns.A:
			ip = rr.A

		case *dns.AAAA:
			ip = rr.AAAA
		}
	}

	return ip
}

func extractOptionsFromEDNS0(r *dns.Msg, optsMap map[uint16][]*edns0Opt, f func([]byte, []*edns0Opt)) {
	o := r.IsEdns0()
	if o == nil {
		return
	}

	var option []dns.EDNS0
	for _, o := range o.Option {
		if local, ok := o.(*dns.EDNS0_LOCAL); ok {
			if m, ok := optsMap[local.Code]; ok {
				f(local.Data, m)
				continue
			}
		}

		option = append(option, o)
	}

	o.Option = option
}

func clearECS(r *dns.Msg) {
	o := r.IsEdns0()
	if o == nil {
		return
	}

	option := make([]dns.EDNS0, 0, len(o.Option))
	for _, opt := range o.Option {
		if _, ok := opt.(*dns.EDNS0_SUBNET); !ok {
			option = append(option, opt)
		}
	}

	o.Option = option
}

func (p *policyPlugin) setRedirectQueryAnswer(ctx context.Context, w dns.ResponseWriter, r *dns.Msg, dst string) (int, error) {
	var rr dns.RR

	qName, qClass := getNameAndClass(r)

	ip := net.ParseIP(dst)
	if ipv4 := ip.To4(); ipv4 != nil {
		rr = &dns.A{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeA,
				Class:  qClass,
			},
			A: ipv4,
		}
	} else if ipv6 := ip.To16(); ipv6 != nil {
		rr = &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeAAAA,
				Class:  qClass,
			},
			AAAA: ipv6,
		}
	} else {
		dst = dns.Fqdn(dst)
		rr = &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeCNAME,
				Class:  qClass,
			},
			Target: dst,
		}

		if r == nil || len(r.Question) <= 0 {
			return dns.RcodeServerFailure, errInvalidDNSMessage
		}

		origName := qName
		r.Question[0].Name = dst

		nw := nonwriter.New(w)
		if _, err := plugin.NextOrFailure(p.Name(), p.next, ctx, nw, r); err != nil {
			r.Question[0].Name = origName
			return dns.RcodeServerFailure, err
		}

		nw.Msg.CopyTo(r)
		r.Question[0].Name = origName

		r.Answer = append([]dns.RR{rr}, r.Answer...)
		r.Authoritative = true
		return r.Rcode, nil
	}

	r.Answer = []dns.RR{rr}
	r.Rcode = dns.RcodeSuccess
	return r.Rcode, nil
}

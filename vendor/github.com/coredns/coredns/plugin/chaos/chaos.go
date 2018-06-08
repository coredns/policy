// Package chaos implements a plugin that answer to 'CH version.bind TXT' type queries.
package chaos

import (
	"context"
	"os"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Chaos allows CoreDNS to reply to CH TXT queries and return author or
// version information.
type Chaos struct {
	Next    plugin.Handler
	Version string
	Authors map[string]bool
}

// ServeDNS implements the plugin.Handler interface.
func (c Chaos) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	if state.QClass() != dns.ClassCHAOS || state.QType() != dns.TypeTXT {
		return plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)

	hdr := dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
	switch state.Name() {
	default:
		return c.Next.ServeDNS(ctx, w, r)
	case "authors.bind.":
		for a := range c.Authors {
			m.Answer = append(m.Answer, &dns.TXT{Hdr: hdr, Txt: []string{trim(a)}})
		}
	case "version.bind.", "version.server.":
		m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{trim(c.Version)}}}
	case "hostname.bind.", "id.server.":
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
		m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{trim(hostname)}}}
	}
	state.SizeAndDo(m)
	w.WriteMsg(m)
	return 0, nil
}

// Name implements the Handler interface.
func (c Chaos) Name() string { return "chaos" }

func trim(s string) string {
	if len(s) < 256 {
		return s
	}
	return s[:255]
}

// +build etcd

package etcd

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/pkg/tls"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func init() {
	ctxt = context.TODO()
}

// Note the key is encoded as DNS name, while in "reality" it is a etcd path.
var services = []*msg.Service{
	{Host: "dev.server1", Port: 8080, Key: "a.server1.dev.region1.skydns.test."},
	{Host: "10.0.0.1", Port: 8080, Key: "a.server1.prod.region1.skydns.test."},
	{Host: "10.0.0.2", Port: 8080, Key: "b.server1.prod.region1.skydns.test."},
	{Host: "::1", Port: 8080, Key: "b.server6.prod.region1.skydns.test."},
	// Unresolvable internal name.
	{Host: "unresolvable.skydns.test", Key: "cname.prod.region1.skydns.test."},
	// Priority.
	{Host: "priority.server1", Priority: 333, Port: 8080, Key: "priority.skydns.test."},
	// Subdomain.
	{Host: "sub.server1", Port: 0, Key: "a.sub.region1.skydns.test."},
	{Host: "sub.server2", Port: 80, Key: "b.sub.region1.skydns.test."},
	{Host: "10.0.0.1", Port: 8080, Key: "c.sub.region1.skydns.test."},
	// Cname loop.
	{Host: "a.cname.skydns.test", Key: "b.cname.skydns.test."},
	{Host: "b.cname.skydns.test", Key: "a.cname.skydns.test."},
	// Nameservers.
	{Host: "10.0.0.2", Key: "a.ns.dns.skydns.test."},
	{Host: "10.0.0.3", Key: "b.ns.dns.skydns.test."},
	// Zone name as A record (basic, return all)
	{Host: "10.0.0.2", Key: "x.skydns_zonea.test."},
	{Host: "10.0.0.3", Key: "y.skydns_zonea.test."},
	// Zone name as A (single entry).
	{Host: "10.0.0.2", Key: "x.skydns_zoneb.test."},
	{Host: "10.0.0.3", Key: "y.skydns_zoneb.test."},
	{Host: "10.0.0.4", Key: "apex.dns.skydns_zoneb.test."},
	// A zone record (rr multiple entries).
	{Host: "10.0.0.2", Key: "x.skydns_zonec.test."},
	{Host: "10.0.0.3", Key: "y.skydns_zonec.test."},
	{Host: "10.0.0.4", Key: "a1.apex.dns.skydns_zonec.test."},
	{Host: "10.0.0.5", Key: "a2.apex.dns.skydns_zonec.test."},
	// AAAA zone record (rr multiple entries mixed with A).
	{Host: "10.0.0.2", Key: "x.skydns_zoned.test."},
	{Host: "10.0.0.3", Key: "y.skydns_zoned.test."},
	{Host: "10.0.0.4", Key: "a1.apex.dns.skydns_zoned.test."},
	{Host: "10.0.0.5", Key: "a2.apex.dns.skydns_zoned.test."},
	{Host: "2003::8:1", Key: "a3.apex.dns.skydns_zoned.test."},
	{Host: "2003::8:2", Key: "a4.apex.dns.skydns_zoned.test."},
	// Reverse.
	{Host: "reverse.example.com", Key: "1.0.0.10.in-addr.arpa."}, // 10.0.0.1
}

var dnsTestCases = []test.Case{
	// SRV Test
	{
		Qname: "a.server1.dev.region1.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{test.SRV("a.server1.dev.region1.skydns.test. 300 SRV 10 100 8080 dev.server1.")},
	},
	// SRV Test (case test)
	{
		Qname: "a.SERVer1.dEv.region1.skydns.tEst.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{test.SRV("a.SERVer1.dEv.region1.skydns.tEst. 300 SRV 10 100 8080 dev.server1.")},
	},
	// NXDOMAIN Test
	{
		Qname: "doesnotexist.skydns.test.", Qtype: dns.TypeA,
		Rcode: dns.RcodeNameError,
		Ns: []dns.RR{
			test.SOA("skydns.test. 300 SOA ns.dns.skydns.test. hostmaster.skydns.test. 0 0 0 0 0"),
		},
	},
	// A Test
	{
		Qname: "a.server1.prod.region1.skydns.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{test.A("a.server1.prod.region1.skydns.test. 300 A 10.0.0.1")},
	},
	// SRV Test where target is IP address
	{
		Qname: "a.server1.prod.region1.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{test.SRV("a.server1.prod.region1.skydns.test. 300 SRV 10 100 8080 a.server1.prod.region1.skydns.test.")},
		Extra:  []dns.RR{test.A("a.server1.prod.region1.skydns.test. 300 A 10.0.0.1")},
	},
	// AAAA Test
	{
		Qname: "b.server6.prod.region1.skydns.test.", Qtype: dns.TypeAAAA,
		Answer: []dns.RR{test.AAAA("b.server6.prod.region1.skydns.test. 300 AAAA ::1")},
	},
	// Multiple A Record Test
	{
		Qname: "server1.prod.region1.skydns.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("server1.prod.region1.skydns.test. 300 A 10.0.0.1"),
			test.A("server1.prod.region1.skydns.test. 300 A 10.0.0.2"),
		},
	},
	// Priority Test
	{
		Qname: "priority.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{test.SRV("priority.skydns.test. 300 SRV 333 100 8080 priority.server1.")},
	},
	// Subdomain Test
	{
		Qname: "sub.region1.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{
			test.SRV("sub.region1.skydns.test. 300 IN SRV 10 33 0 sub.server1."),
			test.SRV("sub.region1.skydns.test. 300 IN SRV 10 33 80 sub.server2."),
			test.SRV("sub.region1.skydns.test. 300 IN SRV 10 33 8080 c.sub.region1.skydns.test."),
		},
		Extra: []dns.RR{test.A("c.sub.region1.skydns.test. 300 IN A 10.0.0.1")},
	},
	// CNAME (unresolvable internal name)
	{
		Qname: "cname.prod.region1.skydns.test.", Qtype: dns.TypeA,
		Ns: []dns.RR{test.SOA("skydns.test. 300 SOA ns.dns.skydns.test. hostmaster.skydns.test. 0 0 0 0 0")},
	},
	// Wildcard Test
	{
		Qname: "*.region1.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 0 sub.server1."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 0 unresolvable.skydns.test."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 80 sub.server2."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 8080 a.server1.prod.region1.skydns.test."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 8080 b.server1.prod.region1.skydns.test."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 8080 b.server6.prod.region1.skydns.test."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 8080 c.sub.region1.skydns.test."),
			test.SRV("*.region1.skydns.test.	300	IN	SRV	10 12 8080 dev.server1."),
		},
		Extra: []dns.RR{
			test.A("a.server1.prod.region1.skydns.test.	300	IN	A	10.0.0.1"),
			test.A("b.server1.prod.region1.skydns.test.	300	IN	A	10.0.0.2"),
			test.AAAA("b.server6.prod.region1.skydns.test.	300	IN	AAAA	::1"),
			test.A("c.sub.region1.skydns.test.	300	IN	A	10.0.0.1"),
		},
	},
	// Wildcard Test
	{
		Qname: "prod.*.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{

			test.SRV("prod.*.skydns.test.	300	IN	SRV	10 25 0 unresolvable.skydns.test."),
			test.SRV("prod.*.skydns.test.	300	IN	SRV	10 25 8080 a.server1.prod.region1.skydns.test."),
			test.SRV("prod.*.skydns.test.	300	IN	SRV	10 25 8080 b.server1.prod.region1.skydns.test."),
			test.SRV("prod.*.skydns.test.	300	IN	SRV	10 25 8080 b.server6.prod.region1.skydns.test."),
		},
		Extra: []dns.RR{
			test.A("a.server1.prod.region1.skydns.test.	300	IN	A	10.0.0.1"),
			test.A("b.server1.prod.region1.skydns.test.	300	IN	A	10.0.0.2"),
			test.AAAA("b.server6.prod.region1.skydns.test.	300	IN	AAAA	::1"),
		},
	},
	// Wildcard Test
	{
		Qname: "prod.any.skydns.test.", Qtype: dns.TypeSRV,
		Answer: []dns.RR{
			test.SRV("prod.any.skydns.test.	300	IN	SRV	10 25 0 unresolvable.skydns.test."),
			test.SRV("prod.any.skydns.test.	300	IN	SRV	10 25 8080 a.server1.prod.region1.skydns.test."),
			test.SRV("prod.any.skydns.test.	300	IN	SRV	10 25 8080 b.server1.prod.region1.skydns.test."),
			test.SRV("prod.any.skydns.test.	300	IN	SRV	10 25 8080 b.server6.prod.region1.skydns.test."),
		},
		Extra: []dns.RR{
			test.A("a.server1.prod.region1.skydns.test.	300	IN	A	10.0.0.1"),
			test.A("b.server1.prod.region1.skydns.test.	300	IN	A	10.0.0.2"),
			test.AAAA("b.server6.prod.region1.skydns.test.	300	IN	AAAA	::1"),
		},
	},
	// CNAME loop detection
	{
		Qname: "a.cname.skydns.test.", Qtype: dns.TypeA,
		Ns: []dns.RR{test.SOA("skydns.test. 300 SOA ns.dns.skydns.test. hostmaster.skydns.test. 1407441600 28800 7200 604800 60")},
	},
	// NODATA Test
	{
		Qname: "a.server1.dev.region1.skydns.test.", Qtype: dns.TypeTXT,
		Ns: []dns.RR{test.SOA("skydns.test. 300 SOA ns.dns.skydns.test. hostmaster.skydns.test. 0 0 0 0 0")},
	},
	// NODATA Test
	{
		Qname: "a.server1.dev.region1.skydns.test.", Qtype: dns.TypeHINFO,
		Ns: []dns.RR{test.SOA("skydns.test. 300 SOA ns.dns.skydns.test. hostmaster.skydns.test. 0 0 0 0 0")},
	},
	// NXDOMAIN Test
	{
		Qname: "a.server1.nonexistent.region1.skydns.test.", Qtype: dns.TypeHINFO, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{test.SOA("skydns.test. 300 SOA ns.dns.skydns.test. hostmaster.skydns.test. 0 0 0 0 0")},
	},
	{
		Qname: "skydns.test.", Qtype: dns.TypeSOA,
		Answer: []dns.RR{test.SOA("skydns.test.	300	IN	SOA	ns.dns.skydns.test. hostmaster.skydns.test. 1460498836 14400 3600 604800 60")},
	},
	// NS Record Test
	{
		Qname: "skydns.test.", Qtype: dns.TypeNS,
		Answer: []dns.RR{
			test.NS("skydns.test. 300 NS a.ns.dns.skydns.test."),
			test.NS("skydns.test. 300 NS b.ns.dns.skydns.test."),
		},
		Extra: []dns.RR{
			test.A("a.ns.dns.skydns.test. 300 A 10.0.0.2"),
			test.A("b.ns.dns.skydns.test. 300 A 10.0.0.3"),
		},
	},
	// NS Record Test
	{
		Qname: "a.skydns.test.", Qtype: dns.TypeNS, Rcode: dns.RcodeNameError,
		Ns: []dns.RR{test.SOA("skydns.test.	300	IN	SOA	ns.dns.skydns.test. hostmaster.skydns.test. 1460498836 14400 3600 604800 60")},
	},
	// A Record For NS Record Test
	{
		Qname: "ns.dns.skydns.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("ns.dns.skydns.test. 300 A 10.0.0.2"),
			test.A("ns.dns.skydns.test. 300 A 10.0.0.3"),
		},
	},
	{
		Qname: "skydns_extra.test.", Qtype: dns.TypeSOA,
		Answer: []dns.RR{test.SOA("skydns_extra.test. 300 IN SOA ns.dns.skydns_extra.test. hostmaster.skydns_extra.test. 1460498836 14400 3600 604800 60")},
	},
	// A Record Test for backward compatibility for zone records
	{
		Qname: "skydns_zonea.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("skydns_zonea.test. 300 A 10.0.0.2"),
			test.A("skydns_zonea.test. 300 A 10.0.0.3"),
		},
	},
	// A Record Test for single A zone record
	{
		Qname: "skydns_zoneb.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{test.A("skydns_zoneb.test. 300 A 10.0.0.4")},
	},
	// A Record Test for multiple A zone records
	{
		Qname: "skydns_zonec.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("skydns_zonec.test. 300 A 10.0.0.4"),
			test.A("skydns_zonec.test. 300 A 10.0.0.5"),
		},
	},
	// A Record Test for multiple mixed A and AAAA records
	{
		Qname: "skydns_zoned.test.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("skydns_zoned.test. 300 A 10.0.0.4"),
			test.A("skydns_zoned.test. 300 A 10.0.0.5"),
		},
	},
	// AAAA Record Test for multiple mixed A and AAAA records
	{
		Qname: "skydns_zoned.test.", Qtype: dns.TypeAAAA,
		Answer: []dns.RR{
			test.AAAA("skydns_zoned.test. 300 AAAA 2003::8:1"),
			test.AAAA("skydns_zoned.test. 300 AAAA 2003::8:2"),
		},
	},
	// Reverse lookup
	{
		Qname: "1.0.0.10.in-addr.arpa.", Qtype: dns.TypePTR,
		Answer: []dns.RR{test.PTR("1.0.0.10.in-addr.arpa. 300 PTR reverse.example.com.")},
	},
}

func newEtcdPlugin() *Etcd {
	ctxt = context.TODO()

	endpoints := []string{"http://localhost:2379"}
	tlsc, _ := tls.NewTLSConfigFromArgs()
	client, _ := newEtcdClient(endpoints, tlsc)

	p := proxy.NewLookup([]string{"8.8.8.8:53"})
	return &Etcd{
		Upstream:   upstream.Upstream{Forward: &p},
		PathPrefix: "skydns",
		Ctx:        context.Background(),
		Zones:      []string{"skydns.test.", "skydns_extra.test.", "skydns_zonea.test.", "skydns_zoneb.test.", "skydns_zonec.test.", "skydns_zoned.test.", "in-addr.arpa."},
		Client:     client,
	}
}

func set(t *testing.T, e *Etcd, k string, ttl time.Duration, m *msg.Service) {
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	path, _ := msg.PathWithWildcard(k, e.PathPrefix)
	e.Client.KV.Put(ctxt, path, string(b))
}

func delete(t *testing.T, e *Etcd, k string) {
	path, _ := msg.PathWithWildcard(k, e.PathPrefix)
	e.Client.Delete(ctxt, path)
}

func TestLookup(t *testing.T) {
	etc := newEtcdPlugin()
	for _, serv := range services {
		set(t, etc, serv.Key, 0, serv)
		defer delete(t, etc, serv.Key)
	}

	for _, tc := range dnsTestCases {
		m := tc.Msg()

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		etc.ServeDNS(ctxt, rec, m)

		resp := rec.Msg
		test.SortAndCheck(t, resp, tc)
	}
}

var ctxt context.Context

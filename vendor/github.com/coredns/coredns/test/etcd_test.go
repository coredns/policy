// +build etcd

package test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/etcd"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"

	etcdc "github.com/coreos/etcd/client"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

func etcdPlugin() *etcd.Etcd {
	etcdCfg := etcdc.Config{
		Endpoints: []string{"http://localhost:2379"},
	}
	cli, _ := etcdc.New(etcdCfg)
	client := etcdc.NewKeysAPI(cli)
	return &etcd.Etcd{Client: client, PathPrefix: "/skydns"}
}

// This test starts two coredns servers (and needs etcd). Configure a stubzones in both (that will loop) and
// will then test if we detect this loop.
func TestEtcdStubLoop(t *testing.T) {
	// TODO(miek)
}

func TestEtcdStubAndProxyLookup(t *testing.T) {
	corefile := `.:0 {
    etcd skydns.local {
        stubzones
        path /skydns
        endpoint http://localhost:2379
        upstream 8.8.8.8:53 8.8.4.4:53
	fallthrough
    }
    proxy . 8.8.8.8:53
}`

	ex, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer ex.Stop()

	etc := etcdPlugin()
	log.SetOutput(ioutil.Discard)

	var ctx = context.TODO()
	for _, serv := range servicesStub { // adds example.{net,org} as stubs
		set(ctx, t, etc, serv.Key, 0, serv)
		defer delete(ctx, t, etc, serv.Key)
	}

	p := proxy.NewLookup([]string{udp}) // use udp port from the server
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}
	resp, err := p.Lookup(state, "example.com.", dns.TypeA)
	if err != nil {
		t.Fatalf("Expected to receive reply, but didn't: %v", err)
	}
	if len(resp.Answer) == 0 {
		t.Fatalf("Expected to at least one RR in the answer section, got none")
	}
	if resp.Answer[0].Header().Rrtype != dns.TypeA {
		t.Errorf("Expected RR to A, got: %d", resp.Answer[0].Header().Rrtype)
	}
	if resp.Answer[0].(*dns.A).A.String() != "93.184.216.34" {
		t.Errorf("Expected 93.184.216.34, got: %s", resp.Answer[0].(*dns.A).A.String())
	}
}

var servicesStub = []*msg.Service{
	// Two tests, ask a question that should return servfail because remote it no accessible
	// and one with edns0 option added, that should return refused.
	{Host: "127.0.0.1", Port: 666, Key: "b.example.org.stub.dns.skydns.test."},
	// Actual test that goes out to the internet.
	{Host: "199.43.132.53", Key: "a.example.net.stub.dns.skydns.test."},
}

// Copied from plugin/etcd/setup_test.go
func set(ctx context.Context, t *testing.T, e *etcd.Etcd, k string, ttl time.Duration, m *msg.Service) {
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	path, _ := msg.PathWithWildcard(k, e.PathPrefix)
	e.Client.Set(ctx, path, string(b), &etcdc.SetOptions{TTL: ttl})
}

// Copied from plugin/etcd/setup_test.go
func delete(ctx context.Context, t *testing.T, e *etcd.Etcd, k string) {
	path, _ := msg.PathWithWildcard(k, e.PathPrefix)
	e.Client.Delete(ctx, path, &etcdc.DeleteOptions{Recursive: false})
}

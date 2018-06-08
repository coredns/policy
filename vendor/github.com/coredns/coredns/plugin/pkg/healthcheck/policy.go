package healthcheck

import (
	"math/rand"
	"sync/atomic"

	"github.com/coredns/coredns/plugin/pkg/log"
)

var (
	// SupportedPolicies is the collection of policies registered
	SupportedPolicies = make(map[string]func() Policy)
)

// RegisterPolicy adds a custom policy to the proxy.
func RegisterPolicy(name string, policy func() Policy) {
	SupportedPolicies[name] = policy
}

// Policy decides how a host will be selected from a pool. When all hosts are unhealthy, it is assumed the
// healthchecking failed. In this case each policy will *randomly* return a host from the pool to prevent
// no traffic to go through at all.
type Policy interface {
	Select(pool HostPool) *UpstreamHost
}

func init() {
	RegisterPolicy("random", func() Policy { return &Random{} })
	RegisterPolicy("least_conn", func() Policy { return &LeastConn{} })
	RegisterPolicy("round_robin", func() Policy { return &RoundRobin{} })
	RegisterPolicy("first", func() Policy { return &First{} })
	// 'sequential' is an alias to 'first' to maintain consistency with the forward plugin
	// should probably remove 'first' in a future release
	RegisterPolicy("sequential", func() Policy { return &First{} })
}

// Random is a policy that selects up hosts from a pool at random.
type Random struct{}

// Select selects an up host at random from the specified pool.
func (r *Random) Select(pool HostPool) *UpstreamHost {
	// instead of just generating a random index
	// this is done to prevent selecting a down host
	var randHost *UpstreamHost
	count := 0
	for _, host := range pool {
		if host.Down() {
			continue
		}
		count++
		if count == 1 {
			randHost = host
		} else {
			r := rand.Int() % count
			if r == (count - 1) {
				randHost = host
			}
		}
	}
	return randHost
}

// Spray is a policy that selects a host from a pool at random. This should be used as a last ditch
// attempt to get a host when all hosts are reporting unhealthy.
type Spray struct{}

// Select selects an up host at random from the specified pool.
func (r *Spray) Select(pool HostPool) *UpstreamHost {
	rnd := rand.Int() % len(pool)
	randHost := pool[rnd]
	log.Warningf("All hosts reported as down, spraying to target: %s", randHost.Name)
	return randHost
}

// LeastConn is a policy that selects the host with the least connections.
type LeastConn struct{}

// Select selects the up host with the least number of connections in the
// pool.  If more than one host has the same least number of connections,
// one of the hosts is chosen at random.
func (r *LeastConn) Select(pool HostPool) *UpstreamHost {
	var bestHost *UpstreamHost
	count := 0
	leastConn := int64(1<<63 - 1)
	for _, host := range pool {
		if host.Down() {
			continue
		}
		hostConns := host.Conns
		if hostConns < leastConn {
			bestHost = host
			leastConn = hostConns
			count = 1
		} else if hostConns == leastConn {
			// randomly select host among hosts with least connections
			count++
			if count == 1 {
				bestHost = host
			} else {
				r := rand.Int() % count
				if r == (count - 1) {
					bestHost = host
				}
			}
		}
	}
	return bestHost
}

// RoundRobin is a policy that selects hosts based on round robin ordering.
type RoundRobin struct {
	Robin uint32
}

// Select selects an up host from the pool using a round robin ordering scheme.
func (r *RoundRobin) Select(pool HostPool) *UpstreamHost {
	poolLen := uint32(len(pool))
	selection := atomic.AddUint32(&r.Robin, 1) % poolLen
	host := pool[selection]
	// if the currently selected host is down, just ffwd to up host
	for i := uint32(1); host.Down() && i < poolLen; i++ {
		host = pool[(selection+i)%poolLen]
	}
	return host
}

// First is a policy that selects always the first healthy host in the list order.
type First struct{}

// Select always the first that is not Down.
func (r *First) Select(pool HostPool) *UpstreamHost {
	for i := 0; i < len(pool); i++ {
		host := pool[i]
		if host.Down() {
			continue
		}
		return host
	}
	// return the first one, anyway none is correct
	return nil
}

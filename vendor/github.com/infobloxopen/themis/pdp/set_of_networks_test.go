package pdp

import (
	"net"
	"testing"
)

func TestSortSetOfNetworks(t *testing.T) {
	assertNetworks(
		SortSetOfNetworks(newIPTree(
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		)),
		[]*net.IPNet{
			makeTestNetwork("192.0.2.0/24"),
			makeTestNetwork("2001:db8::/32"),
			makeTestNetwork("192.0.2.16/28"),
		},
		"SortSetOfNetworks",
		t,
	)
}

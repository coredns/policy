package pdp

import (
	"net"
	"sort"

	"github.com/infobloxopen/go-trees/iptree"
)

// SortSetOfNetworks converts set of networks to a slice ordered by assigned
// integer values. It panics if given tree contains not int value.
func SortSetOfNetworks(v *iptree.Tree) []*net.IPNet {
	pairs := newNetworkPairList(v)
	sort.Sort(pairs)

	list := make([]*net.IPNet, len(pairs))
	for i, pair := range pairs {
		list[i] = pair.value
	}

	return list
}

type networkPair struct {
	value *net.IPNet
	order int
}

type networkPairList []networkPair

func newNetworkPairList(v *iptree.Tree) networkPairList {
	pairs := make(networkPairList, 0)
	for p := range v.Enumerate() {
		pairs = append(pairs, networkPair{p.Key, p.Value.(int)})
	}

	return pairs
}

func (p networkPairList) Len() int {
	return len(p)
}

func (p networkPairList) Less(i, j int) bool {
	return p[i].order < p[j].order
}

func (p networkPairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

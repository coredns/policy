package pdp

import (
	"sort"

	"github.com/infobloxopen/go-trees/domaintree"
)

// SortSetOfDomains converts set of domains to a slice of strings ordered by
// assigned integer values. Strings represent human-readable domain names.
// It panics if given tree contains not int value.
func SortSetOfDomains(v *domaintree.Node) []string {
	pairs := newDomainPairList(v)
	sort.Sort(pairs)

	list := make([]string, len(pairs))
	for i, pair := range pairs {
		list[i] = pair.value
	}

	return list
}

type domainPair struct {
	value string
	order int
}

type domainPairList []domainPair

func newDomainPairList(v *domaintree.Node) domainPairList {
	pairs := make(domainPairList, 0)
	for p := range v.Enumerate() {
		pairs = append(pairs, domainPair{p.Key, p.Value.(int)})
	}

	return pairs
}

func (p domainPairList) Len() int {
	return len(p)
}

func (p domainPairList) Less(i, j int) bool {
	return p[i].order < p[j].order
}

func (p domainPairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

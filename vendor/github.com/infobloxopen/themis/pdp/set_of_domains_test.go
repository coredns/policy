package pdp

import "testing"

func TestSortSetOfDomains(t *testing.T) {
	assertStrings(
		SortSetOfDomains(newDomainTree(
			makeTestDomain("example.com"),
			makeTestDomain("example.gov"),
			makeTestDomain("www.example.com"),
		)),
		[]string{"example.com", "example.gov", "www.example.com"},
		"SortSetOfDomains",
		t,
	)
}

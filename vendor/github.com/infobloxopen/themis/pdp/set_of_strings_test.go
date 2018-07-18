package pdp

import "testing"

func TestSortSetOfStrings(t *testing.T) {
	assertStrings(
		SortSetOfStrings(newStrTree("First", "Second", "Third")),
		[]string{"First", "Second", "Third"},
		"SortSetOfStrings",
		t,
	)
}

package strtree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestNewTree(t *testing.T) {
	r := NewTree()

	assertTree(r, TestEmptyTree, "empty tree", t)
}

func TestNewTreeWithCustomComparison(t *testing.T) {
	r := NewTreeWithCustomComparison(strings.Compare)

	assertTree(r, TestEmptyTree, "empty tree", t)
}

func TestInsert(t *testing.T) {
	var r *Tree

	r = r.Insert("k", "v")
	assertTree(r, TestSingleNodeTree, "single node tree", t)

	r = NewTree()
	r = r.Insert("k", "v")
	assertTree(r, TestSingleNodeTree, "single node tree", t)

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("1", nil)
	r = r.Insert("2", nil)
	assertTree(r, TestThreeNodeTreeRed, "tree 012", t)

	r = NewTree()
	r = r.Insert("1", nil)
	r = r.Insert("2", nil)
	r = r.Insert("0", nil)
	assertTree(r, TestThreeNodeTreeRed, "tree 120", t)

	r = NewTree()
	r = r.Insert("2", nil)
	r = r.Insert("0", nil)
	r = r.Insert("1", nil)
	assertTree(r, TestThreeNodeTreeRed, "tree 201", t)

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("2", nil)
	r = r.Insert("1", nil)
	assertTree(r, TestThreeNodeTreeRed, "tree 021", t)

	r = NewTree()
	r = r.Insert("1", nil)
	r = r.Insert("0", nil)
	r = r.Insert("2", nil)
	assertTree(r, TestThreeNodeTreeRed, "tree 102", t)

	r = NewTree()
	r = r.Insert("2", nil)
	r = r.Insert("1", nil)
	r = r.Insert("0", nil)
	assertTree(r, TestThreeNodeTreeRed, "tree 210", t)

	r = NewTree()
	r = r.Insert("1", nil)
	r = r.Insert("0", nil)
	r = r.Insert("4", nil)
	r = r.Insert("2", nil)
	r = r.Insert("3", nil)
	assertTree(r, TestFiveNodeTree, "tree 10423", t)

	r = NewTree()
	r = r.Insert("F", nil)
	r = r.Insert("E", nil)
	r = r.Insert("D", nil)
	r = r.Insert("C", nil)
	r = r.Insert("B", nil)
	r = r.Insert("A", nil)
	r = r.Insert("9", nil)
	r = r.Insert("8", nil)
	r = r.Insert("7", nil)
	r = r.Insert("6", nil)
	r = r.Insert("5", nil)
	r = r.Insert("4", nil)
	r = r.Insert("3", nil)
	r = r.Insert("2", nil)
	r = r.Insert("1", nil)
	r = r.Insert("0", nil)
	assertTree(r, Test16InversedNodeTree, "tree inversed 16 nodes", t)

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("1", nil)
	r = r.Insert("2", nil)
	r = r.Insert("3", nil)
	r = r.Insert("4", nil)
	r = r.Insert("5", nil)
	r = r.Insert("6", nil)
	r = r.Insert("7", nil)
	r = r.Insert("8", nil)
	r = r.Insert("9", nil)
	r = r.Insert("A", nil)
	r = r.Insert("B", nil)
	r = r.Insert("C", nil)
	r = r.Insert("D", nil)
	r = r.Insert("E", nil)
	r = r.Insert("F", nil)
	assertTree(r, Test16DirectNodeTree, "tree direct 16 nodes", t)

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("2", nil)
	r = r.Insert("4", nil)
	r = r.Insert("6", nil)
	r = r.Insert("8", nil)
	r = r.Insert("A", nil)
	r = r.Insert("C", nil)
	r = r.Insert("E", nil)
	r = r.Insert("1", nil)
	r = r.Insert("3", nil)
	r = r.Insert("5", nil)
	r = r.Insert("7", nil)
	r = r.Insert("9", nil)
	r = r.Insert("B", nil)
	r = r.Insert("D", nil)
	r = r.Insert("F", nil)
	assertTree(r, Test16AlternatingNodeTree, "tree alternating 16 nodes", t)

	r = NewTree()
	r = r.Insert("F", nil)
	r = r.Insert("D", nil)
	r = r.Insert("B", nil)
	r = r.Insert("9", nil)
	r = r.Insert("7", nil)
	r = r.Insert("5", nil)
	r = r.Insert("3", nil)
	r = r.Insert("1", nil)
	r = r.Insert("E", nil)
	r = r.Insert("C", nil)
	r = r.Insert("A", nil)
	r = r.Insert("8", nil)
	r = r.Insert("6", nil)
	r = r.Insert("4", nil)
	r = r.Insert("2", nil)
	r = r.Insert("0", nil)
	assertTree(r, Test16AlternatingInversedNodeTree, "tree alternating inversed 16 nodes", t)

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("3", nil)
	r = r.Insert("6", nil)
	r = r.Insert("9", nil)
	r = r.Insert("C", nil)
	r = r.Insert("F", nil)
	r = r.Insert("1", nil)
	r = r.Insert("2", nil)
	r = r.Insert("4", nil)
	r = r.Insert("5", nil)
	r = r.Insert("7", nil)
	r = r.Insert("8", nil)
	r = r.Insert("A", nil)
	r = r.Insert("B", nil)
	r = r.Insert("D", nil)
	r = r.Insert("E", nil)
	assertTree(r, Test16_3AltNodeTree, "tree alternating by 3 16 nodes", t)

	r = NewTree()
	r = r.Insert("00", nil)
	r = r.Insert("02", nil)
	r = r.Insert("04", nil)
	r = r.Insert("06", nil)
	r = r.Insert("08", nil)
	r = r.Insert("0A", nil)
	r = r.Insert("0C", nil)
	r = r.Insert("0E", nil)
	r = r.Insert("10", nil)
	r = r.Insert("12", nil)
	r = r.Insert("14", nil)
	r = r.Insert("16", nil)
	r = r.Insert("18", nil)
	r = r.Insert("1A", nil)
	r = r.Insert("1C", nil)
	r = r.Insert("1E", nil)
	r = r.Insert("01", nil)
	r = r.Insert("03", nil)
	r = r.Insert("05", nil)
	r = r.Insert("07", nil)
	r = r.Insert("09", nil)
	r = r.Insert("0B", nil)
	r = r.Insert("0D", nil)
	r = r.Insert("0F", nil)
	r = r.Insert("11", nil)
	r = r.Insert("13", nil)
	r = r.Insert("15", nil)
	r = r.Insert("17", nil)
	r = r.Insert("19", nil)
	r = r.Insert("1B", nil)
	r = r.Insert("1D", nil)
	r = r.Insert("1F", nil)
	assertTree(r, Test32AlternatingNodeTree, "tree with alternating 32 nodes", t)

	r = nil
	r = r.Insert("00", "test-1")
	assertTree(r, TestTreeSameNodeOnce, "tree with same node first insertion", t)
	r = r.Insert("00", "test-2")
	assertTree(r, TestTreeSameNodeTwice, "tree with same node second insertion", t)
}

func TestInplaceInsert(t *testing.T) {
	r := NewTree()

	r.InplaceInsert("k", "v")
	assertTree(r, TestSingleNodeTree, "single node inplace tree", t)

	r = NewTree()
	r.InplaceInsert("k", "v")
	assertTree(r, TestSingleNodeTree, "single node inplace tree", t)

	r = NewTree()
	r.InplaceInsert("0", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("2", nil)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 012", t)

	r = NewTree()
	r.InplaceInsert("1", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("0", nil)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 120", t)

	r = NewTree()
	r.InplaceInsert("2", nil)
	r.InplaceInsert("0", nil)
	r.InplaceInsert("1", nil)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 201", t)

	r = NewTree()
	r.InplaceInsert("0", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("1", nil)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 021", t)

	r = NewTree()
	r.InplaceInsert("1", nil)
	r.InplaceInsert("0", nil)
	r.InplaceInsert("2", nil)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 102", t)

	r = NewTree()
	r.InplaceInsert("2", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("0", nil)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 210", t)

	r = NewTree()
	r.InplaceInsert("1", nil)
	r.InplaceInsert("0", nil)
	r.InplaceInsert("4", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("3", nil)
	assertTree(r, TestFiveNodeTree, "inplace tree 10423", t)

	r = NewTree()
	r.InplaceInsert("F", nil)
	r.InplaceInsert("E", nil)
	r.InplaceInsert("D", nil)
	r.InplaceInsert("C", nil)
	r.InplaceInsert("B", nil)
	r.InplaceInsert("A", nil)
	r.InplaceInsert("9", nil)
	r.InplaceInsert("8", nil)
	r.InplaceInsert("7", nil)
	r.InplaceInsert("6", nil)
	r.InplaceInsert("5", nil)
	r.InplaceInsert("4", nil)
	r.InplaceInsert("3", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("0", nil)
	assertTree(r, Test16InversedNodeTree, "inplace tree inversed 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("0", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("3", nil)
	r.InplaceInsert("4", nil)
	r.InplaceInsert("5", nil)
	r.InplaceInsert("6", nil)
	r.InplaceInsert("7", nil)
	r.InplaceInsert("8", nil)
	r.InplaceInsert("9", nil)
	r.InplaceInsert("A", nil)
	r.InplaceInsert("B", nil)
	r.InplaceInsert("C", nil)
	r.InplaceInsert("D", nil)
	r.InplaceInsert("E", nil)
	r.InplaceInsert("F", nil)
	assertTree(r, Test16DirectNodeTree, "inplace tree direct 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("0", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("4", nil)
	r.InplaceInsert("6", nil)
	r.InplaceInsert("8", nil)
	r.InplaceInsert("A", nil)
	r.InplaceInsert("C", nil)
	r.InplaceInsert("E", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("3", nil)
	r.InplaceInsert("5", nil)
	r.InplaceInsert("7", nil)
	r.InplaceInsert("9", nil)
	r.InplaceInsert("B", nil)
	r.InplaceInsert("D", nil)
	r.InplaceInsert("F", nil)
	assertTree(r, Test16AlternatingNodeTree, "inplace tree alternating 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("F", nil)
	r.InplaceInsert("D", nil)
	r.InplaceInsert("B", nil)
	r.InplaceInsert("9", nil)
	r.InplaceInsert("7", nil)
	r.InplaceInsert("5", nil)
	r.InplaceInsert("3", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("E", nil)
	r.InplaceInsert("C", nil)
	r.InplaceInsert("A", nil)
	r.InplaceInsert("8", nil)
	r.InplaceInsert("6", nil)
	r.InplaceInsert("4", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("0", nil)
	assertTree(r, Test16AlternatingInversedNodeTree, "inplace tree alternating inversed 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("0", nil)
	r.InplaceInsert("3", nil)
	r.InplaceInsert("6", nil)
	r.InplaceInsert("9", nil)
	r.InplaceInsert("C", nil)
	r.InplaceInsert("F", nil)
	r.InplaceInsert("1", nil)
	r.InplaceInsert("2", nil)
	r.InplaceInsert("4", nil)
	r.InplaceInsert("5", nil)
	r.InplaceInsert("7", nil)
	r.InplaceInsert("8", nil)
	r.InplaceInsert("A", nil)
	r.InplaceInsert("B", nil)
	r.InplaceInsert("D", nil)
	r.InplaceInsert("E", nil)
	assertTree(r, Test16_3AltNodeTree, "inplace tree alternating by 3 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("00", nil)
	r.InplaceInsert("02", nil)
	r.InplaceInsert("04", nil)
	r.InplaceInsert("06", nil)
	r.InplaceInsert("08", nil)
	r.InplaceInsert("0A", nil)
	r.InplaceInsert("0C", nil)
	r.InplaceInsert("0E", nil)
	r.InplaceInsert("10", nil)
	r.InplaceInsert("12", nil)
	r.InplaceInsert("14", nil)
	r.InplaceInsert("16", nil)
	r.InplaceInsert("18", nil)
	r.InplaceInsert("1A", nil)
	r.InplaceInsert("1C", nil)
	r.InplaceInsert("1E", nil)
	r.InplaceInsert("01", nil)
	r.InplaceInsert("03", nil)
	r.InplaceInsert("05", nil)
	r.InplaceInsert("07", nil)
	r.InplaceInsert("09", nil)
	r.InplaceInsert("0B", nil)
	r.InplaceInsert("0D", nil)
	r.InplaceInsert("0F", nil)
	r.InplaceInsert("11", nil)
	r.InplaceInsert("13", nil)
	r.InplaceInsert("15", nil)
	r.InplaceInsert("17", nil)
	r.InplaceInsert("19", nil)
	r.InplaceInsert("1B", nil)
	r.InplaceInsert("1D", nil)
	r.InplaceInsert("1F", nil)
	assertTree(r, Test32AlternatingNodeTree, "inplace tree with alternating 32 nodes", t)

	r = nil
	assertPanic(func() { r.InplaceInsert("00", "test") }, "nil tree inplace insertion", t)
}

func TestGet(t *testing.T) {
	var r *Tree

	v, ok := r.Get("0")
	if ok {
		t.Errorf("Expected nothing but got %T (%#v)", v, v)
	}

	r = NewTree()
	r = r.Insert("1", "test-1")
	r = r.Insert("0", "test-0")
	r = r.Insert("4", "test-4")
	r = r.Insert("2", "test-2")
	r = r.Insert("3", "test-3")

	v, ok = r.Get("3")
	if !ok {
		t.Errorf("Expected \"test-3\" but got nothing")
	} else {
		s, ok := v.(string)
		if !ok {
			t.Errorf("Expected \"test-3\" but got %T (%#v)", v, v)
		} else if s != "test-3" {
			t.Errorf("Expected \"test-3\" but got %q", s)
		}
	}

	v, ok = r.Get("F")
	if ok {
		t.Errorf("Expected nothing but got %T (%#v)", v, v)
	}
}

func TestEnumerate(t *testing.T) {
	var r *Tree

	assertEnumerate(r.Enumerate(), "empty tree", t)

	r = NewTree()
	r = r.Insert("1", "test-1")
	r = r.Insert("0", "test-0")
	r = r.Insert("4", "test-4")
	r = r.Insert("2", "test-2")
	r = r.Insert("3", "test-3")
	assertEnumerate(r.Enumerate(), "enumeration of tree 10423", t,
		"\"0\": \"test-0\"\n",
		"\"1\": \"test-1\"\n",
		"\"2\": \"test-2\"\n",
		"\"3\": \"test-3\"\n",
		"\"4\": \"test-4\"\n")
}

func TestDelete(t *testing.T) {
	var r *Tree

	r, ok := r.Delete("test")
	if ok {
		t.Errorf("Expected nothing to be deleted from empty tree but something has been deleted:\n%s", r.Dot())
	}

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("3", nil)
	r = r.Insert("6", nil)
	r = r.Insert("9", nil)
	r = r.Insert("C", nil)
	r = r.Insert("F", nil)
	r = r.Insert("1", nil)
	r = r.Insert("2", nil)
	r = r.Insert("4", nil)
	r = r.Insert("5", nil)
	r = r.Insert("7", nil)
	r = r.Insert("8", nil)
	r = r.Insert("A", nil)
	r = r.Insert("B", nil)
	r = r.Insert("D", nil)
	r = r.Insert("E", nil)

	r, ok = r.Delete("81")
	if ok {
		t.Errorf("Expected nothing to be deleted by key \"81\" but something has been deleted")
	}
	assertTree(r, TestTreeAfterNonExistingNodeDel, "tree after non-existing node deletion", t)

	r, ok = r.Delete("6")
	if !ok {
		t.Errorf("Expected node \"6\" to be deleted but got nothing")
	}
	assertTree(r, TestTreeAfterNode6Deletion, "tree after node 6 deletion", t)

	r, ok = r.Delete("7")
	if !ok {
		t.Errorf("Expected node \"7\" to be deleted but got nothing")
	}
	r, ok = r.Delete("8")
	if !ok {
		t.Errorf("Expected node \"8\" to be deleted but got nothing")
	}
	r, ok = r.Delete("5")
	if !ok {
		t.Errorf("Expected node \"5\" to be deleted but got nothing")
	}
	r, ok = r.Delete("9")
	if !ok {
		t.Errorf("Expected node \"9\" to be deleted but got nothing")
	}
	assertTree(r, TestTreeAfterNodes7859Deletion, "tree after nodes 7, 8, 5 and 9 deletion", t)

	r, ok = r.Delete("C")
	if !ok {
		t.Errorf("Expected node \"C\" to be deleted but got nothing")
	}
	r, ok = r.Delete("E")
	if !ok {
		t.Errorf("Expected node \"E\" to be deleted but got nothing")
	}
	r, ok = r.Delete("D")
	if !ok {
		t.Errorf("Expected node \"D\" to be deleted but got nothing")
	}
	r, ok = r.Delete("A")
	if !ok {
		t.Errorf("Expected node \"A\" to be deleted but got nothing")
	}
	r, ok = r.Delete("B")
	if !ok {
		t.Errorf("Expected node \"B\" to be deleted but got nothing")
	}
	r, ok = r.Delete("4")
	if !ok {
		t.Errorf("Expected node \"4\" to be deleted but got nothing")
	}
	r, ok = r.Delete("F")
	if !ok {
		t.Errorf("Expected node \"F\" to be deleted but got nothing")
	}
	r, ok = r.Delete("0")
	if !ok {
		t.Errorf("Expected node \"0\" to be deleted but got nothing")
	}
	r, ok = r.Delete("3")
	if !ok {
		t.Errorf("Expected node \"3\" to be deleted but got nothing")
	}
	r, ok = r.Delete("1")
	if !ok {
		t.Errorf("Expected node \"1\" to be deleted but got nothing")
	}
	r, ok = r.Delete("2")
	if !ok {
		t.Errorf("Expected node \"2\" to be deleted but got nothing")
	}
	assertTree(r, TestEmptyTree, "tree after rest nodes deletion", t)
}

func TestIsEmpty(t *testing.T) {
	var r *Tree

	if !r.IsEmpty() {
		t.Errorf("Expected nil tree to be empty")
	}

	r = NewTree()
	r = r.Insert("0", nil)
	r = r.Insert("3", nil)
	r = r.Insert("6", nil)
	if r.IsEmpty() {
		t.Errorf("Expected three nodes tree to be not empty")
	}

	r, ok := r.Delete("3")
	if !ok {
		t.Errorf("Expected element \"3\" to be deleted")
	}

	if r.IsEmpty() {
		t.Errorf("Expected two nodes tree to be not empty")
	}

	r, _ = r.Delete("0")
	r, _ = r.Delete("6")

	if !r.IsEmpty() {
		t.Errorf("Expected empty non-nil tree to be empty")
	}
}

const (
	TestEmptyTree = `digraph d {
N0 [label="nil" style=filled fontcolor=white fillcolor=black]
}
`

	TestSingleNodeTree = `digraph d {
N0 [label="k: \"k\" v: \"\"v\"\"" style=filled fontcolor=white fillcolor=black]
}
`

	TestThreeNodeTreeRed = `digraph d {
N0 [label="1" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="0" style=filled fillcolor=red]
N2 [label="2" style=filled fillcolor=red]
}
`

	TestFiveNodeTree = `digraph d {
N0 [label="1" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="0" style=filled fontcolor=white fillcolor=black]
N2 [label="3" style=filled fontcolor=white fillcolor=black]
N2 -> { N3 N4 }
N3 [label="2" style=filled fillcolor=red]
N4 [label="4" style=filled fillcolor=red]
}
`

	Test16InversedNodeTree = `digraph d {
N0 [label="C" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="8" style=filled fillcolor=red]
N1 -> { N3 N4 }
N2 [label="E" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="4" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="A" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="D" style=filled fontcolor=white fillcolor=black]
N6 [label="F" style=filled fontcolor=white fillcolor=black]
N7 [label="2" style=filled fillcolor=red]
N7 -> { N11 N12 }
N8 [label="6" style=filled fillcolor=red]
N8 -> { N13 N14 }
N9 [label="9" style=filled fontcolor=white fillcolor=black]
N10 [label="B" style=filled fontcolor=white fillcolor=black]
N11 [label="1" style=filled fontcolor=white fillcolor=black]
N11 -> { N15 N16 }
N12 [label="3" style=filled fontcolor=white fillcolor=black]
N13 [label="5" style=filled fontcolor=white fillcolor=black]
N14 [label="7" style=filled fontcolor=white fillcolor=black]
N15 [label="0" style=filled fillcolor=red]
N16 [label="nil" style=filled fontcolor=white fillcolor=black]
}
`

	Test16DirectNodeTree = `digraph d {
N0 [label="3" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="1" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="7" style=filled fillcolor=red]
N2 -> { N5 N6 }
N3 [label="0" style=filled fontcolor=white fillcolor=black]
N4 [label="2" style=filled fontcolor=white fillcolor=black]
N5 [label="5" style=filled fontcolor=white fillcolor=black]
N5 -> { N7 N8 }
N6 [label="B" style=filled fontcolor=white fillcolor=black]
N6 -> { N9 N10 }
N7 [label="4" style=filled fontcolor=white fillcolor=black]
N8 [label="6" style=filled fontcolor=white fillcolor=black]
N9 [label="9" style=filled fillcolor=red]
N9 -> { N11 N12 }
N10 [label="D" style=filled fillcolor=red]
N10 -> { N13 N14 }
N11 [label="8" style=filled fontcolor=white fillcolor=black]
N12 [label="A" style=filled fontcolor=white fillcolor=black]
N13 [label="C" style=filled fontcolor=white fillcolor=black]
N14 [label="E" style=filled fontcolor=white fillcolor=black]
N14 -> { N15 N16 }
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="F" style=filled fillcolor=red]
}
`

	Test16AlternatingNodeTree = `digraph d {
N0 [label="6" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="2" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="A" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="0" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="4" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="8" style=filled fontcolor=white fillcolor=black]
N5 -> { N11 N12 }
N6 [label="C" style=filled fillcolor=red]
N6 -> { N13 N14 }
N7 [label="nil" style=filled fontcolor=white fillcolor=black]
N8 [label="1" style=filled fillcolor=red]
N9 [label="3" style=filled fillcolor=red]
N10 [label="5" style=filled fillcolor=red]
N11 [label="7" style=filled fillcolor=red]
N12 [label="9" style=filled fillcolor=red]
N13 [label="B" style=filled fontcolor=white fillcolor=black]
N14 [label="E" style=filled fontcolor=white fillcolor=black]
N14 -> { N15 N16 }
N15 [label="D" style=filled fillcolor=red]
N16 [label="F" style=filled fillcolor=red]
}
`

	Test16AlternatingInversedNodeTree = `digraph d {
N0 [label="9" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="5" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="D" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="3" style=filled fillcolor=red]
N3 -> { N7 N8 }
N4 [label="7" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="B" style=filled fontcolor=white fillcolor=black]
N5 -> { N11 N12 }
N6 [label="F" style=filled fontcolor=white fillcolor=black]
N6 -> { N13 N14 }
N7 [label="1" style=filled fontcolor=white fillcolor=black]
N7 -> { N15 N16 }
N8 [label="4" style=filled fontcolor=white fillcolor=black]
N9 [label="6" style=filled fillcolor=red]
N10 [label="8" style=filled fillcolor=red]
N11 [label="A" style=filled fillcolor=red]
N12 [label="C" style=filled fillcolor=red]
N13 [label="E" style=filled fillcolor=red]
N14 [label="nil" style=filled fontcolor=white fillcolor=black]
N15 [label="0" style=filled fillcolor=red]
N16 [label="2" style=filled fillcolor=red]
}
`

	Test16_3AltNodeTree = `digraph d {
N0 [label="5" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="3" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="9" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="1" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="4" style=filled fontcolor=white fillcolor=black]
N5 [label="7" style=filled fontcolor=white fillcolor=black]
N5 -> { N9 N10 }
N6 [label="C" style=filled fillcolor=red]
N6 -> { N11 N12 }
N7 [label="0" style=filled fillcolor=red]
N8 [label="2" style=filled fillcolor=red]
N9 [label="6" style=filled fillcolor=red]
N10 [label="8" style=filled fillcolor=red]
N11 [label="A" style=filled fontcolor=white fillcolor=black]
N11 -> { N13 N14 }
N12 [label="E" style=filled fontcolor=white fillcolor=black]
N12 -> { N15 N16 }
N13 [label="nil" style=filled fontcolor=white fillcolor=black]
N14 [label="B" style=filled fillcolor=red]
N15 [label="D" style=filled fillcolor=red]
N16 [label="F" style=filled fillcolor=red]
}
`

	Test32AlternatingNodeTree = `digraph d {
N0 [label="0E" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="06" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="16" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="02" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="0A" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="12" style=filled fontcolor=white fillcolor=black]
N5 -> { N11 N12 }
N6 [label="1A" style=filled fontcolor=white fillcolor=black]
N6 -> { N13 N14 }
N7 [label="00" style=filled fontcolor=white fillcolor=black]
N7 -> { N15 N16 }
N8 [label="04" style=filled fontcolor=white fillcolor=black]
N8 -> { N17 N18 }
N9 [label="08" style=filled fontcolor=white fillcolor=black]
N9 -> { N19 N20 }
N10 [label="0C" style=filled fontcolor=white fillcolor=black]
N10 -> { N21 N22 }
N11 [label="10" style=filled fontcolor=white fillcolor=black]
N11 -> { N23 N24 }
N12 [label="14" style=filled fontcolor=white fillcolor=black]
N12 -> { N25 N26 }
N13 [label="18" style=filled fontcolor=white fillcolor=black]
N13 -> { N27 N28 }
N14 [label="1C" style=filled fillcolor=red]
N14 -> { N29 N30 }
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="01" style=filled fillcolor=red]
N17 [label="03" style=filled fillcolor=red]
N18 [label="05" style=filled fillcolor=red]
N19 [label="07" style=filled fillcolor=red]
N20 [label="09" style=filled fillcolor=red]
N21 [label="0B" style=filled fillcolor=red]
N22 [label="0D" style=filled fillcolor=red]
N23 [label="0F" style=filled fillcolor=red]
N24 [label="11" style=filled fillcolor=red]
N25 [label="13" style=filled fillcolor=red]
N26 [label="15" style=filled fillcolor=red]
N27 [label="17" style=filled fillcolor=red]
N28 [label="19" style=filled fillcolor=red]
N29 [label="1B" style=filled fontcolor=white fillcolor=black]
N30 [label="1E" style=filled fontcolor=white fillcolor=black]
N30 -> { N31 N32 }
N31 [label="1D" style=filled fillcolor=red]
N32 [label="1F" style=filled fillcolor=red]
}
`

	TestTreeAfterNonExistingNodeDel = `digraph d {
N0 [label="5" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="3" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="C" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="1" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="4" style=filled fontcolor=white fillcolor=black]
N5 [label="9" style=filled fillcolor=red]
N5 -> { N9 N10 }
N6 [label="E" style=filled fontcolor=white fillcolor=black]
N6 -> { N11 N12 }
N7 [label="0" style=filled fillcolor=red]
N8 [label="2" style=filled fillcolor=red]
N9 [label="7" style=filled fontcolor=white fillcolor=black]
N9 -> { N13 N14 }
N10 [label="A" style=filled fontcolor=white fillcolor=black]
N10 -> { N15 N16 }
N11 [label="D" style=filled fillcolor=red]
N12 [label="F" style=filled fillcolor=red]
N13 [label="6" style=filled fillcolor=red]
N14 [label="8" style=filled fillcolor=red]
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="B" style=filled fillcolor=red]
}
`

	TestTreeAfterNode6Deletion = `digraph d {
N0 [label="5" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="3" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="C" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="1" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="4" style=filled fontcolor=white fillcolor=black]
N5 [label="9" style=filled fillcolor=red]
N5 -> { N9 N10 }
N6 [label="E" style=filled fontcolor=white fillcolor=black]
N6 -> { N11 N12 }
N7 [label="0" style=filled fillcolor=red]
N8 [label="2" style=filled fillcolor=red]
N9 [label="7" style=filled fontcolor=white fillcolor=black]
N9 -> { N13 N14 }
N10 [label="A" style=filled fontcolor=white fillcolor=black]
N10 -> { N15 N16 }
N11 [label="D" style=filled fillcolor=red]
N12 [label="F" style=filled fillcolor=red]
N13 [label="nil" style=filled fontcolor=white fillcolor=black]
N14 [label="8" style=filled fillcolor=red]
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="B" style=filled fillcolor=red]
}
`

	TestTreeAfterNodes7859Deletion = `digraph d {
N0 [label="A" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="2" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="C" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="1" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="4" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="B" style=filled fontcolor=white fillcolor=black]
N6 [label="E" style=filled fontcolor=white fillcolor=black]
N6 -> { N11 N12 }
N7 [label="0" style=filled fillcolor=red]
N8 [label="nil" style=filled fontcolor=white fillcolor=black]
N9 [label="3" style=filled fillcolor=red]
N10 [label="nil" style=filled fontcolor=white fillcolor=black]
N11 [label="D" style=filled fillcolor=red]
N12 [label="F" style=filled fillcolor=red]
}
`

	TestTreeSameNodeOnce = `digraph d {
N0 [label="k: \"00\" v: \"\"test-1\"\"" style=filled fontcolor=white fillcolor=black]
}
`

	TestTreeSameNodeTwice = `digraph d {
N0 [label="k: \"00\" v: \"\"test-2\"\"" style=filled fontcolor=white fillcolor=black]
}
`
)

func assertTree(r *Tree, e, desc string, t *testing.T) {
	assertStringLists(difflib.SplitLines(r.Dot()), difflib.SplitLines(e), desc, t)
}

func assertEnumerate(ch chan Pair, desc string, t *testing.T, e ...string) {
	pairs := []string{}
	for p := range ch {
		s, ok := p.Value.(string)
		if ok {
			pairs = append(pairs, fmt.Sprintf("%q: %q\n", p.Key, s))
		} else {
			pairs = append(pairs, fmt.Sprintf("%q: %T (%#v)\n", p.Key, p.Value, p.Value))
		}
	}

	assertStringLists(pairs, e, desc, t)
}

func assertStringLists(v, e []string, desc string, t *testing.T) {
	ctx := difflib.ContextDiff{
		A:        e,
		B:        v,
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("Can't compare \"%s\": %s", desc, err))
	}

	if len(diff) > 0 {
		t.Errorf("\"%s\" doesn't match:\n%s", desc, diff)
	}
}

func assertPanic(f func(), desc string, t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic from %s but got nothing", desc)
		}
	}()

	f()
}

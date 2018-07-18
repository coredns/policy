package numtree

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestInsert32(t *testing.T) {
	var r *Node32

	r = r.Insert(0, 32, "test")
	assertTree32(r, TestTree32WithSingleNodeInserted,
		"32-tree with single node inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 18, "bottom")
	r = r.Insert(0xAAAAAAAA, 9, "top")
	assertTree32(r, TestTree32WithTopAfterBottomToLeftNodesInserted,
		"32-tree with top after bottom to left nodes inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 18, "bottom")
	r = r.Insert(0xAAAAAAAA, 10, "top")
	assertTree32(r, TestTree32WithTopAfterBottomToRightNodesInserted,
		"32-tree with top after bottom to right nodes inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 18, "bottom")
	r = r.Insert(0xABAAAAAA, 10, "top")
	assertTree32(r, TestTree32WithTopAfterBottomAndAdditionalNotLeafNodesInserted,
		"32-tree with top after bottom and additional not leaf nodes inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 18, "bottom")
	oldR := r.Insert(0xABAAAAAA, 10, "top")
	newR := oldR.Insert(0xABAAAAAA, 7, "root")
	assertTree32(oldR, TestTree32WithOldTopReplacingTopAfterBottomNodesInserted,
		"32-tree with old top replacing top after bottom nodes inserted", t)
	assertTree32(newR, TestTree32WithNewTopReplacingTopAfterBottomNodesInserted,
		"32-tree with new top replacing top after bottom nodes inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 9, "top")
	r = r.Insert(0xAAAAAAAA, 18, "bottom")
	assertTree32(r, TestTree32WithTopBeforeBottomToLeftNodesInserted,
		"32-tree with top before bottom to left nodes inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 10, "top")
	r = r.Insert(0xAAAAAAAA, 18, "bottom")
	assertTree32(r, TestTree32WithTopBeforeBottomToRightNodesInserted,
		"32-tree with top before bottom to right nodes inserted", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 7, "L1")
	r = r.Insert(0xABAAAAAA, 9, "L2")
	r = r.Insert(0xAAAAAAAA, 18, "L3")
	r = r.Insert(0xAABAAAAA, 19, "L4")
	assertTree32(r, TestTree32WithTopBeforeBottomSeveralLevelNodesInserted,
		"32-tree with top before bottom several level nodes inserted", t)

	r = nil
	r = r.Insert(0, -10, nil)
	assertTree32(r, TestTree32WithNegativeNumberOfBits,
		"32-tree with negative number of significant bits", t)

	r = nil
	r = r.Insert(0, 33, nil)
	assertTree32(r, TestTree32WithTooBigNumberOfBits,
		"32-tree with too big number of significant bits", t)

	r = nil
	for i := uint32(0); i < 256; i++ {
		r = r.Insert(i, 32, fmt.Sprintf("%02x", i))
	}
	assertTree32(r, TestTree32BigTreeInsertions,
		"32-tree big tree", t)

	r = nil
	for i := uint32(0); i < 256; i++ {
		r = r.Insert(inv32[i]<<24, 32, fmt.Sprintf("%02x", inv32[i]))
	}
	assertTree32(r, TestTree32BigTreeInvertedInsertions,
		"32-tree big tree with inverted keys", t)
}

func TestInplaceInsert32(t *testing.T) {
	var r *Node32

	r = r.InplaceInsert(0, 32, "test")
	assertTree32(r, TestTree32WithSingleNodeInserted,
		"32-tree with single node inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 18, "bottom")
	r = r.InplaceInsert(0xAAAAAAAA, 9, "top")
	assertTree32(r, TestTree32WithTopAfterBottomToLeftNodesInserted,
		"32-tree with top after bottom to left nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 18, "bottom")
	r = r.InplaceInsert(0xAAAAAAAA, 10, "top")
	assertTree32(r, TestTree32WithTopAfterBottomToRightNodesInserted,
		"32-tree with top after bottom to right nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 18, "bottom")
	r = r.InplaceInsert(0xABAAAAAA, 10, "top")
	assertTree32(r, TestTree32WithTopAfterBottomAndAdditionalNotLeafNodesInserted,
		"32-tree with top after bottom and additional not leaf nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 18, "bottom")
	r = r.InplaceInsert(0xABAAAAAA, 10, "top")
	assertTree32(r, TestTree32WithOldTopReplacingTopAfterBottomNodesInserted,
		"32-tree with old top replacing top after bottom nodes inplace inserted", t)
	r = r.InplaceInsert(0xABAAAAAA, 7, "root")
	assertTree32(r, TestTree32WithNewTopReplacingTopAfterBottomNodesInserted,
		"32-tree with new top replacing top after bottom nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 9, "top")
	r = r.InplaceInsert(0xAAAAAAAA, 18, "bottom")
	assertTree32(r, TestTree32WithTopBeforeBottomToLeftNodesInserted,
		"32-tree with top before bottom to left nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 10, "top")
	r = r.InplaceInsert(0xAAAAAAAA, 18, "bottom")
	assertTree32(r, TestTree32WithTopBeforeBottomToRightNodesInserted,
		"32-tree with top before bottom to right nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0xAAAAAAAA, 7, "L1")
	r = r.InplaceInsert(0xABAAAAAA, 9, "L2")
	r = r.InplaceInsert(0xAAAAAAAA, 18, "L3")
	r = r.InplaceInsert(0xAABAAAAA, 19, "L4")
	assertTree32(r, TestTree32WithTopBeforeBottomSeveralLevelNodesInserted,
		"32-tree with top before bottom several level nodes inplace inserted", t)

	r = nil
	r = r.InplaceInsert(0, -10, nil)
	assertTree32(r, TestTree32WithNegativeNumberOfBits,
		"32-tree with negative number of significant bits (inplace)", t)

	r = nil
	r = r.InplaceInsert(0, 33, nil)
	assertTree32(r, TestTree32WithTooBigNumberOfBits,
		"32-tree with too big number of significant bits (inplace)", t)

	r = nil
	for i := uint32(0); i < 256; i++ {
		r = r.InplaceInsert(i, 32, fmt.Sprintf("%02x", i))
	}
	assertTree32(r, TestTree32BigTreeInsertions,
		"32-tree big tree (inplace)", t)

	r = nil
	for i := uint32(0); i < 256; i++ {
		r = r.InplaceInsert(inv32[i]<<24, 32, fmt.Sprintf("%02x", inv32[i]))
	}
	assertTree32(r, TestTree32BigTreeInvertedInsertions,
		"32-tree big tree with inverted keys (inplace)", t)
}

func TestEnumerate32(t *testing.T) {
	var r *Node32

	ch := r.Enumerate()
	assertSequence32(ch, t, "32-tree empty tree")

	r = r.Insert(0xAAAAAAAA, 7, "L1")
	r = r.Insert(0xA8AAAAAA, 9, "L2.1")
	r = r.Insert(0xABAAAAAA, 9, "L2.2")
	r = r.Insert(0xAAAAAAAA, 18, "L3")
	r = r.Insert(0xAAABAAAA, 24, "L5")
	r = r.Insert(0xAABAAAAA, 19, "L4")
	ch = r.Enumerate()
	assertSequence32(ch, t, "32-tree for enumeration",
		"0xa8aaaaaa/9: \"L2.1\"",
		"0xaaaaaaaa/7: \"L1\"",
		"0xaaaaaaaa/18: \"L3\"",
		"0xaaabaaaa/24: \"L5\"",
		"0xaabaaaaa/19: \"L4\"",
		"0xabaaaaaa/9: \"L2.2\"")
}

func TestMatch32(t *testing.T) {
	var r *Node32

	v, ok := r.Match(0, 0)
	assertTreeMatch(v, ok, nil,
		"32-bit empty tree", t)

	r = r.Insert(0xAAAAAAAA, 7, "L1")
	r = r.Insert(0xA8AAAAAA, 9, "L2.1")
	r = r.Insert(0xABAAAAAA, 9, "L2.2")
	r = r.Insert(0xAAAAAAAA, 18, "L3")
	r = r.Insert(0xAABAAAAA, 19, "L4")

	v, ok = r.Match(0, -1)
	assertTreeMatch(v, ok, nil,
		"32-tree match with negative significant bits", t)

	v, ok = r.Match(0xAAAAAAAA, 35)
	assertTreeMatch(v, ok, wrapStr("L3"),
		"32-tree match with overflow significant bits number", t)

	v, ok = r.Match(0xAAAAAAAA, 5)
	assertTreeMatch(v, ok, nil,
		"32-tree match with small significant bits number", t)

	v, ok = r.Match(0xA8AAAAAA, 9)
	assertTreeMatch(v, ok, wrapStr("L2.1"),
		"32-tree match with exact match to a node", t)

	v, ok = r.Match(0xA9AAAAAA, 9)
	assertTreeMatch(v, ok, nil,
		"32-tree match with exact not match to a node", t)

	v, ok = r.Match(0xAABAAACA, 32)
	assertTreeMatch(v, ok, wrapStr("L4"),
		"32-tree match with contains match to child node", t)

	v, ok = r.Match(0xABAAAAAA, 9)
	assertTreeMatch(v, ok, wrapStr("L2.2"),
		"32-tree match with exact match to child node", t)

	v, ok = r.Match(0xA80AAAAA, 11)
	assertTreeMatch(v, ok, nil,
		"32-tree match with contains match to non-leaf node", t)
}

func TestExactMatch32(t *testing.T) {
	var r *Node32

	v, ok := r.ExactMatch(0, 0)
	assertTreeMatch(v, ok, nil,
		"32-bit empty tree", t)

	r = r.Insert(0xAAAAAAAA, 7, "L1")
	r = r.Insert(0xA8AAAAAA, 9, "L2.1")
	r = r.Insert(0xABAAAAAA, 9, "L2.2")
	r = r.Insert(0xAAAAAAAA, 18, "L3")
	r = r.Insert(0xAABAAAAA, 19, "L4")

	v, ok = r.ExactMatch(0, -1)
	assertTreeMatch(v, ok, nil,
		"32-tree exact match with negative significant bits", t)

	v, ok = r.ExactMatch(0xAAAAAAAA, 35)
	assertTreeMatch(v, ok, nil,
		"32-tree exact match with overflow significant bits number", t)

	v, ok = r.ExactMatch(0xAAAAAAAA, 5)
	assertTreeMatch(v, ok, nil,
		"32-tree exact match with small significant bits number", t)

	v, ok = r.ExactMatch(0xA8AAAAAA, 9)
	assertTreeMatch(v, ok, wrapStr("L2.1"),
		"32-tree exact match with exact match to a node", t)

	v, ok = r.ExactMatch(0xA9AAAAAA, 9)
	assertTreeMatch(v, ok, nil,
		"32-tree exact match with exact not match to a node", t)

	v, ok = r.ExactMatch(0xAABAAACA, 32)
	assertTreeMatch(v, ok, nil,
		"32-tree exact match with contains not match to child node", t)

	v, ok = r.ExactMatch(0xABAAAAAA, 9)
	assertTreeMatch(v, ok, wrapStr("L2.2"),
		"32-tree match with exact match to child node", t)

	v, ok = r.ExactMatch(0xA80AAAAA, 11)
	assertTreeMatch(v, ok, nil,
		"32-tree match with contains match to non-leaf node", t)
}

func TestDelete32(t *testing.T) {
	var (
		r  *Node32
		ok bool
	)

	r, ok = r.Delete(0, 32)
	assertTree32Delete(r, ok, "", "32-bit empty tree", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 18, "test")
	r, ok = r.Delete(0xAAAAAAAA, 9)
	assertTree32Delete(r, ok, TestTree32EmptyTree,
		"32-tree with contained node", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 18, "test")
	r, ok = r.Delete(0xBBBBBBBB, 9)
	assertTree32Delete(r, ok, "", "32-tree with not contained node", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 9, "test")
	r = r.Insert(0xAAAAAAAA, 18, "test")
	r, ok = r.Delete(0xBBBBBBBB, 10)
	assertTree32Delete(r, ok, "", "32-tree with not containing node", t)

	r, ok = r.Delete(0xAAEAAAAA, 10)
	assertTree32Delete(r, ok, "", "32-tree with empty branch", t)

	r, ok = r.Delete(0xAAABBBBB, 16)
	assertTree32Delete(r, ok, "", "32-tree with not contained branch", t)

	r, ok = r.Delete(0xAAAAAAAA, 16)
	assertTree32Delete(r, ok, TestTree32WithDeletedChildNode,
		"32-tree with deleted child node", t)

	r, ok = r.Delete(0, -10)
	assertTree32Delete(r, ok, TestTree32EmptyTree,
		"32-tree with deleted all nodes by negative number of significant bits", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 9, "test")
	r = r.Insert(0xAAAAAAAA, 32, "test")
	r, ok = r.Delete(0xAAAAAAAA, 35)
	assertTree32Delete(r, ok, TestTree32WithDeletedChildNode,
		"32-tree with deleted child node by too big number of significant bits", t)

	r = nil
	r = r.Insert(0xAAAAAAAA, 7, "L1")
	r = r.Insert(0xA8AAAAAA, 9, "L2.1")
	r = r.Insert(0xABAAAAAA, 9, "L2.2")
	r = r.Insert(0xAAAAAAAA, 18, "L3")
	r = r.Insert(0xAAABAAAA, 24, "L5")
	r = r.Insert(0xAABAAAAA, 19, "L4")

	r, ok = r.Delete(0xAABAAAAA, 19)
	assertTree32Delete(r, ok, TestTree32WithDeletedChildAndNonLeafNodes,
		"32-tree with deleted child and non-leaf node", t)

	r, ok = r.Delete(0xAAAAAAAA, 18)
	assertTree32Delete(r, ok, TestTree32WithDeletedTwoChildrenAndNonLeafNodes,
		"32-tree with deleted two children and non-leaf nodes", t)
}

func TestClz32(t *testing.T) {
	assertClz32(0x00000000, 32, t)
	assertClz32(0x00000001, 31, t)
	assertClz32(0x00000002, 30, t)
	assertClz32(0x00000003, 30, t)
	assertClz32(0x00000004, 29, t)
	assertClz32(0x00000005, 29, t)
	assertClz32(0x00000006, 29, t)
	assertClz32(0x00000007, 29, t)
	assertClz32(0x00000008, 28, t)

	assertClz32(0x00000010, 27, t)
	assertClz32(0x00000020, 26, t)
	assertClz32(0x00000040, 25, t)
	assertClz32(0x00000080, 24, t)

	assertClz32(0x00000100, 23, t)
	assertClz32(0x00000200, 22, t)
	assertClz32(0x00000400, 21, t)
	assertClz32(0x00000800, 20, t)

	assertClz32(0x00001000, 19, t)
	assertClz32(0x00002000, 18, t)
	assertClz32(0x00004000, 17, t)
	assertClz32(0x00008000, 16, t)

	assertClz32(0x00010000, 15, t)
	assertClz32(0x00020000, 14, t)
	assertClz32(0x00040000, 13, t)
	assertClz32(0x00080000, 12, t)

	assertClz32(0x00100000, 11, t)
	assertClz32(0x00200000, 10, t)
	assertClz32(0x00400000, 9, t)
	assertClz32(0x00800000, 8, t)

	assertClz32(0x01000000, 7, t)
	assertClz32(0x02000000, 6, t)
	assertClz32(0x04000000, 5, t)
	assertClz32(0x08000000, 4, t)

	assertClz32(0x10000000, 3, t)
	assertClz32(0x20000000, 2, t)
	assertClz32(0x40000000, 1, t)
	assertClz32(0x80000000, 0, t)
}

func assertTree32(r *Node32, e, desc string, t *testing.T) {
	assertStringLists(difflib.SplitLines(r.Dot()), difflib.SplitLines(e), desc, t)
}

func assertSequence32(ch chan *Node32, t *testing.T, desc string, e ...string) {
	items := []string{}
	for n := range ch {
		if n == nil {
			items = append(items, fmt.Sprintf("%#v\n", n))
			continue
		}

		s, ok := n.Value.(string)
		if ok {
			s = fmt.Sprintf("%q", s)
		} else {
			s = fmt.Sprintf("%#v (non-string type %T)", n.Value, n.Value)
		}

		items = append(items, fmt.Sprintf("0x%08x/%d: %s\n", n.Key, n.Bits, s))
	}

	eItems := make([]string, len(e))
	for i, item := range e {
		eItems[i] = item + "\n"
	}

	assertStringLists(items, eItems, desc, t)
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

func wrapStr(s string) *string {
	return &s
}

func assertTreeMatch(v interface{}, ok bool, e *string, desc string, t *testing.T) {
	if e == nil {
		if ok {
			t.Errorf("Expected no result for %s but got ok: true, value: %#v", desc, v)
		}

		return
	}

	if !ok {
		t.Errorf("Expected some result for %s but got ok: false, value: %#v", desc, v)
		return
	}

	s, ok := v.(string)
	if !ok {
		t.Errorf("Expected result for %s to be string but got %T: %#v", desc, v, v)
		return
	}

	if s != *e {
		t.Errorf("Expected %q as result for %s but got %q", desc, *e, s)
	}
}

func assertTree32Delete(r *Node32, ok bool, e string, desc string, t *testing.T) {
	if len(e) > 0 {
		if !ok {
			t.Errorf("Expected something to be deleted from %s but it isn't and got old root:\n%s\n", desc, r.Dot())
			return
		}

		assertTree32(r, e, desc, t)
	} else if ok {
		t.Errorf("Expected nothing to be deleted from %s but it is and got new root:\n%s\n", desc, r.Dot())
	}
}

func assertClz32(x uint32, c uint8, t *testing.T) {
	r := clz32(x)
	if r != c {
		t.Errorf("Expected %d as result of clz32(0x%08x) but got %d", c, x, r)
	}
}

const (
	TestTree32WithSingleNodeInserted = `digraph d {
N0 [label="k: 00000000, b: 32, v: \"\"test\"\""]
}
`

	TestTree32WithTopAfterBottomToLeftNodesInserted = `digraph d {
N0 [label="k: aaaaaaaa, b: 9, v: \"\"top\"\""]
N0 -> { N1 N2 }
N1 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
N2 [label="nil"]
}
`

	TestTree32WithTopAfterBottomToRightNodesInserted = `digraph d {
N0 [label="k: aaaaaaaa, b: 10, v: \"\"top\"\""]
N0 -> { N1 N2 }
N1 [label="nil"]
N2 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
}
`

	TestTree32WithTopAfterBottomAndAdditionalNotLeafNodesInserted = `digraph d {
N0 [label="k: aa000000, b: 7"]
N0 -> { N1 N2 }
N1 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
N2 [label="k: abaaaaaa, b: 10, v: \"\"top\"\""]
}
`

	TestTree32WithOldTopReplacingTopAfterBottomNodesInserted = `digraph d {
N0 [label="k: aa000000, b: 7"]
N0 -> { N1 N2 }
N1 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
N2 [label="k: abaaaaaa, b: 10, v: \"\"top\"\""]
}
`

	TestTree32WithNewTopReplacingTopAfterBottomNodesInserted = `digraph d {
N0 [label="k: abaaaaaa, b: 7, v: \"\"root\"\""]
N0 -> { N1 N2 }
N1 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
N2 [label="k: abaaaaaa, b: 10, v: \"\"top\"\""]
}
`

	TestTree32WithTopBeforeBottomToLeftNodesInserted = `digraph d {
N0 [label="k: aaaaaaaa, b: 9, v: \"\"top\"\""]
N0 -> { N1 N2 }
N1 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
N2 [label="nil"]
}
`

	TestTree32WithTopBeforeBottomToRightNodesInserted = `digraph d {
N0 [label="k: aaaaaaaa, b: 10, v: \"\"top\"\""]
N0 -> { N1 N2 }
N1 [label="nil"]
N2 [label="k: aaaaaaaa, b: 18, v: \"\"bottom\"\""]
}
`

	TestTree32WithTopBeforeBottomSeveralLevelNodesInserted = `digraph d {
N0 [label="k: aaaaaaaa, b: 7, v: \"\"L1\"\""]
N0 -> { N1 N2 }
N1 [label="k: aaa00000, b: 11"]
N1 -> { N3 N4 }
N2 [label="k: abaaaaaa, b: 9, v: \"\"L2\"\""]
N3 [label="k: aaaaaaaa, b: 18, v: \"\"L3\"\""]
N4 [label="k: aabaaaaa, b: 19, v: \"\"L4\"\""]
}
`

	TestTree32WithNegativeNumberOfBits = `digraph d {
N0 [label="k: 00000000, b: 0, v: \"<nil>\""]
}
`

	TestTree32WithTooBigNumberOfBits = `digraph d {
N0 [label="k: 00000000, b: 32, v: \"<nil>\""]
}
`

	TestTree32BigTreeInsertions = `digraph d {
N0 [label="k: 00000000, b: 24"]
N0 -> { N1 N2 }
N1 [label="k: 00000000, b: 25"]
N1 -> { N3 N4 }
N2 [label="k: 00000080, b: 25"]
N2 -> { N5 N6 }
N3 [label="k: 00000000, b: 26"]
N3 -> { N7 N8 }
N4 [label="k: 00000040, b: 26"]
N4 -> { N9 N10 }
N5 [label="k: 00000080, b: 26"]
N5 -> { N11 N12 }
N6 [label="k: 000000c0, b: 26"]
N6 -> { N13 N14 }
N7 [label="k: 00000000, b: 27"]
N7 -> { N15 N16 }
N8 [label="k: 00000020, b: 27"]
N8 -> { N17 N18 }
N9 [label="k: 00000040, b: 27"]
N9 -> { N19 N20 }
N10 [label="k: 00000060, b: 27"]
N10 -> { N21 N22 }
N11 [label="k: 00000080, b: 27"]
N11 -> { N23 N24 }
N12 [label="k: 000000a0, b: 27"]
N12 -> { N25 N26 }
N13 [label="k: 000000c0, b: 27"]
N13 -> { N27 N28 }
N14 [label="k: 000000e0, b: 27"]
N14 -> { N29 N30 }
N15 [label="k: 00000000, b: 28"]
N15 -> { N31 N32 }
N16 [label="k: 00000010, b: 28"]
N16 -> { N33 N34 }
N17 [label="k: 00000020, b: 28"]
N17 -> { N35 N36 }
N18 [label="k: 00000030, b: 28"]
N18 -> { N37 N38 }
N19 [label="k: 00000040, b: 28"]
N19 -> { N39 N40 }
N20 [label="k: 00000050, b: 28"]
N20 -> { N41 N42 }
N21 [label="k: 00000060, b: 28"]
N21 -> { N43 N44 }
N22 [label="k: 00000070, b: 28"]
N22 -> { N45 N46 }
N23 [label="k: 00000080, b: 28"]
N23 -> { N47 N48 }
N24 [label="k: 00000090, b: 28"]
N24 -> { N49 N50 }
N25 [label="k: 000000a0, b: 28"]
N25 -> { N51 N52 }
N26 [label="k: 000000b0, b: 28"]
N26 -> { N53 N54 }
N27 [label="k: 000000c0, b: 28"]
N27 -> { N55 N56 }
N28 [label="k: 000000d0, b: 28"]
N28 -> { N57 N58 }
N29 [label="k: 000000e0, b: 28"]
N29 -> { N59 N60 }
N30 [label="k: 000000f0, b: 28"]
N30 -> { N61 N62 }
N31 [label="k: 00000000, b: 29"]
N31 -> { N63 N64 }
N32 [label="k: 00000008, b: 29"]
N32 -> { N65 N66 }
N33 [label="k: 00000010, b: 29"]
N33 -> { N67 N68 }
N34 [label="k: 00000018, b: 29"]
N34 -> { N69 N70 }
N35 [label="k: 00000020, b: 29"]
N35 -> { N71 N72 }
N36 [label="k: 00000028, b: 29"]
N36 -> { N73 N74 }
N37 [label="k: 00000030, b: 29"]
N37 -> { N75 N76 }
N38 [label="k: 00000038, b: 29"]
N38 -> { N77 N78 }
N39 [label="k: 00000040, b: 29"]
N39 -> { N79 N80 }
N40 [label="k: 00000048, b: 29"]
N40 -> { N81 N82 }
N41 [label="k: 00000050, b: 29"]
N41 -> { N83 N84 }
N42 [label="k: 00000058, b: 29"]
N42 -> { N85 N86 }
N43 [label="k: 00000060, b: 29"]
N43 -> { N87 N88 }
N44 [label="k: 00000068, b: 29"]
N44 -> { N89 N90 }
N45 [label="k: 00000070, b: 29"]
N45 -> { N91 N92 }
N46 [label="k: 00000078, b: 29"]
N46 -> { N93 N94 }
N47 [label="k: 00000080, b: 29"]
N47 -> { N95 N96 }
N48 [label="k: 00000088, b: 29"]
N48 -> { N97 N98 }
N49 [label="k: 00000090, b: 29"]
N49 -> { N99 N100 }
N50 [label="k: 00000098, b: 29"]
N50 -> { N101 N102 }
N51 [label="k: 000000a0, b: 29"]
N51 -> { N103 N104 }
N52 [label="k: 000000a8, b: 29"]
N52 -> { N105 N106 }
N53 [label="k: 000000b0, b: 29"]
N53 -> { N107 N108 }
N54 [label="k: 000000b8, b: 29"]
N54 -> { N109 N110 }
N55 [label="k: 000000c0, b: 29"]
N55 -> { N111 N112 }
N56 [label="k: 000000c8, b: 29"]
N56 -> { N113 N114 }
N57 [label="k: 000000d0, b: 29"]
N57 -> { N115 N116 }
N58 [label="k: 000000d8, b: 29"]
N58 -> { N117 N118 }
N59 [label="k: 000000e0, b: 29"]
N59 -> { N119 N120 }
N60 [label="k: 000000e8, b: 29"]
N60 -> { N121 N122 }
N61 [label="k: 000000f0, b: 29"]
N61 -> { N123 N124 }
N62 [label="k: 000000f8, b: 29"]
N62 -> { N125 N126 }
N63 [label="k: 00000000, b: 30"]
N63 -> { N127 N128 }
N64 [label="k: 00000004, b: 30"]
N64 -> { N129 N130 }
N65 [label="k: 00000008, b: 30"]
N65 -> { N131 N132 }
N66 [label="k: 0000000c, b: 30"]
N66 -> { N133 N134 }
N67 [label="k: 00000010, b: 30"]
N67 -> { N135 N136 }
N68 [label="k: 00000014, b: 30"]
N68 -> { N137 N138 }
N69 [label="k: 00000018, b: 30"]
N69 -> { N139 N140 }
N70 [label="k: 0000001c, b: 30"]
N70 -> { N141 N142 }
N71 [label="k: 00000020, b: 30"]
N71 -> { N143 N144 }
N72 [label="k: 00000024, b: 30"]
N72 -> { N145 N146 }
N73 [label="k: 00000028, b: 30"]
N73 -> { N147 N148 }
N74 [label="k: 0000002c, b: 30"]
N74 -> { N149 N150 }
N75 [label="k: 00000030, b: 30"]
N75 -> { N151 N152 }
N76 [label="k: 00000034, b: 30"]
N76 -> { N153 N154 }
N77 [label="k: 00000038, b: 30"]
N77 -> { N155 N156 }
N78 [label="k: 0000003c, b: 30"]
N78 -> { N157 N158 }
N79 [label="k: 00000040, b: 30"]
N79 -> { N159 N160 }
N80 [label="k: 00000044, b: 30"]
N80 -> { N161 N162 }
N81 [label="k: 00000048, b: 30"]
N81 -> { N163 N164 }
N82 [label="k: 0000004c, b: 30"]
N82 -> { N165 N166 }
N83 [label="k: 00000050, b: 30"]
N83 -> { N167 N168 }
N84 [label="k: 00000054, b: 30"]
N84 -> { N169 N170 }
N85 [label="k: 00000058, b: 30"]
N85 -> { N171 N172 }
N86 [label="k: 0000005c, b: 30"]
N86 -> { N173 N174 }
N87 [label="k: 00000060, b: 30"]
N87 -> { N175 N176 }
N88 [label="k: 00000064, b: 30"]
N88 -> { N177 N178 }
N89 [label="k: 00000068, b: 30"]
N89 -> { N179 N180 }
N90 [label="k: 0000006c, b: 30"]
N90 -> { N181 N182 }
N91 [label="k: 00000070, b: 30"]
N91 -> { N183 N184 }
N92 [label="k: 00000074, b: 30"]
N92 -> { N185 N186 }
N93 [label="k: 00000078, b: 30"]
N93 -> { N187 N188 }
N94 [label="k: 0000007c, b: 30"]
N94 -> { N189 N190 }
N95 [label="k: 00000080, b: 30"]
N95 -> { N191 N192 }
N96 [label="k: 00000084, b: 30"]
N96 -> { N193 N194 }
N97 [label="k: 00000088, b: 30"]
N97 -> { N195 N196 }
N98 [label="k: 0000008c, b: 30"]
N98 -> { N197 N198 }
N99 [label="k: 00000090, b: 30"]
N99 -> { N199 N200 }
N100 [label="k: 00000094, b: 30"]
N100 -> { N201 N202 }
N101 [label="k: 00000098, b: 30"]
N101 -> { N203 N204 }
N102 [label="k: 0000009c, b: 30"]
N102 -> { N205 N206 }
N103 [label="k: 000000a0, b: 30"]
N103 -> { N207 N208 }
N104 [label="k: 000000a4, b: 30"]
N104 -> { N209 N210 }
N105 [label="k: 000000a8, b: 30"]
N105 -> { N211 N212 }
N106 [label="k: 000000ac, b: 30"]
N106 -> { N213 N214 }
N107 [label="k: 000000b0, b: 30"]
N107 -> { N215 N216 }
N108 [label="k: 000000b4, b: 30"]
N108 -> { N217 N218 }
N109 [label="k: 000000b8, b: 30"]
N109 -> { N219 N220 }
N110 [label="k: 000000bc, b: 30"]
N110 -> { N221 N222 }
N111 [label="k: 000000c0, b: 30"]
N111 -> { N223 N224 }
N112 [label="k: 000000c4, b: 30"]
N112 -> { N225 N226 }
N113 [label="k: 000000c8, b: 30"]
N113 -> { N227 N228 }
N114 [label="k: 000000cc, b: 30"]
N114 -> { N229 N230 }
N115 [label="k: 000000d0, b: 30"]
N115 -> { N231 N232 }
N116 [label="k: 000000d4, b: 30"]
N116 -> { N233 N234 }
N117 [label="k: 000000d8, b: 30"]
N117 -> { N235 N236 }
N118 [label="k: 000000dc, b: 30"]
N118 -> { N237 N238 }
N119 [label="k: 000000e0, b: 30"]
N119 -> { N239 N240 }
N120 [label="k: 000000e4, b: 30"]
N120 -> { N241 N242 }
N121 [label="k: 000000e8, b: 30"]
N121 -> { N243 N244 }
N122 [label="k: 000000ec, b: 30"]
N122 -> { N245 N246 }
N123 [label="k: 000000f0, b: 30"]
N123 -> { N247 N248 }
N124 [label="k: 000000f4, b: 30"]
N124 -> { N249 N250 }
N125 [label="k: 000000f8, b: 30"]
N125 -> { N251 N252 }
N126 [label="k: 000000fc, b: 30"]
N126 -> { N253 N254 }
N127 [label="k: 00000000, b: 31"]
N127 -> { N255 N256 }
N128 [label="k: 00000002, b: 31"]
N128 -> { N257 N258 }
N129 [label="k: 00000004, b: 31"]
N129 -> { N259 N260 }
N130 [label="k: 00000006, b: 31"]
N130 -> { N261 N262 }
N131 [label="k: 00000008, b: 31"]
N131 -> { N263 N264 }
N132 [label="k: 0000000a, b: 31"]
N132 -> { N265 N266 }
N133 [label="k: 0000000c, b: 31"]
N133 -> { N267 N268 }
N134 [label="k: 0000000e, b: 31"]
N134 -> { N269 N270 }
N135 [label="k: 00000010, b: 31"]
N135 -> { N271 N272 }
N136 [label="k: 00000012, b: 31"]
N136 -> { N273 N274 }
N137 [label="k: 00000014, b: 31"]
N137 -> { N275 N276 }
N138 [label="k: 00000016, b: 31"]
N138 -> { N277 N278 }
N139 [label="k: 00000018, b: 31"]
N139 -> { N279 N280 }
N140 [label="k: 0000001a, b: 31"]
N140 -> { N281 N282 }
N141 [label="k: 0000001c, b: 31"]
N141 -> { N283 N284 }
N142 [label="k: 0000001e, b: 31"]
N142 -> { N285 N286 }
N143 [label="k: 00000020, b: 31"]
N143 -> { N287 N288 }
N144 [label="k: 00000022, b: 31"]
N144 -> { N289 N290 }
N145 [label="k: 00000024, b: 31"]
N145 -> { N291 N292 }
N146 [label="k: 00000026, b: 31"]
N146 -> { N293 N294 }
N147 [label="k: 00000028, b: 31"]
N147 -> { N295 N296 }
N148 [label="k: 0000002a, b: 31"]
N148 -> { N297 N298 }
N149 [label="k: 0000002c, b: 31"]
N149 -> { N299 N300 }
N150 [label="k: 0000002e, b: 31"]
N150 -> { N301 N302 }
N151 [label="k: 00000030, b: 31"]
N151 -> { N303 N304 }
N152 [label="k: 00000032, b: 31"]
N152 -> { N305 N306 }
N153 [label="k: 00000034, b: 31"]
N153 -> { N307 N308 }
N154 [label="k: 00000036, b: 31"]
N154 -> { N309 N310 }
N155 [label="k: 00000038, b: 31"]
N155 -> { N311 N312 }
N156 [label="k: 0000003a, b: 31"]
N156 -> { N313 N314 }
N157 [label="k: 0000003c, b: 31"]
N157 -> { N315 N316 }
N158 [label="k: 0000003e, b: 31"]
N158 -> { N317 N318 }
N159 [label="k: 00000040, b: 31"]
N159 -> { N319 N320 }
N160 [label="k: 00000042, b: 31"]
N160 -> { N321 N322 }
N161 [label="k: 00000044, b: 31"]
N161 -> { N323 N324 }
N162 [label="k: 00000046, b: 31"]
N162 -> { N325 N326 }
N163 [label="k: 00000048, b: 31"]
N163 -> { N327 N328 }
N164 [label="k: 0000004a, b: 31"]
N164 -> { N329 N330 }
N165 [label="k: 0000004c, b: 31"]
N165 -> { N331 N332 }
N166 [label="k: 0000004e, b: 31"]
N166 -> { N333 N334 }
N167 [label="k: 00000050, b: 31"]
N167 -> { N335 N336 }
N168 [label="k: 00000052, b: 31"]
N168 -> { N337 N338 }
N169 [label="k: 00000054, b: 31"]
N169 -> { N339 N340 }
N170 [label="k: 00000056, b: 31"]
N170 -> { N341 N342 }
N171 [label="k: 00000058, b: 31"]
N171 -> { N343 N344 }
N172 [label="k: 0000005a, b: 31"]
N172 -> { N345 N346 }
N173 [label="k: 0000005c, b: 31"]
N173 -> { N347 N348 }
N174 [label="k: 0000005e, b: 31"]
N174 -> { N349 N350 }
N175 [label="k: 00000060, b: 31"]
N175 -> { N351 N352 }
N176 [label="k: 00000062, b: 31"]
N176 -> { N353 N354 }
N177 [label="k: 00000064, b: 31"]
N177 -> { N355 N356 }
N178 [label="k: 00000066, b: 31"]
N178 -> { N357 N358 }
N179 [label="k: 00000068, b: 31"]
N179 -> { N359 N360 }
N180 [label="k: 0000006a, b: 31"]
N180 -> { N361 N362 }
N181 [label="k: 0000006c, b: 31"]
N181 -> { N363 N364 }
N182 [label="k: 0000006e, b: 31"]
N182 -> { N365 N366 }
N183 [label="k: 00000070, b: 31"]
N183 -> { N367 N368 }
N184 [label="k: 00000072, b: 31"]
N184 -> { N369 N370 }
N185 [label="k: 00000074, b: 31"]
N185 -> { N371 N372 }
N186 [label="k: 00000076, b: 31"]
N186 -> { N373 N374 }
N187 [label="k: 00000078, b: 31"]
N187 -> { N375 N376 }
N188 [label="k: 0000007a, b: 31"]
N188 -> { N377 N378 }
N189 [label="k: 0000007c, b: 31"]
N189 -> { N379 N380 }
N190 [label="k: 0000007e, b: 31"]
N190 -> { N381 N382 }
N191 [label="k: 00000080, b: 31"]
N191 -> { N383 N384 }
N192 [label="k: 00000082, b: 31"]
N192 -> { N385 N386 }
N193 [label="k: 00000084, b: 31"]
N193 -> { N387 N388 }
N194 [label="k: 00000086, b: 31"]
N194 -> { N389 N390 }
N195 [label="k: 00000088, b: 31"]
N195 -> { N391 N392 }
N196 [label="k: 0000008a, b: 31"]
N196 -> { N393 N394 }
N197 [label="k: 0000008c, b: 31"]
N197 -> { N395 N396 }
N198 [label="k: 0000008e, b: 31"]
N198 -> { N397 N398 }
N199 [label="k: 00000090, b: 31"]
N199 -> { N399 N400 }
N200 [label="k: 00000092, b: 31"]
N200 -> { N401 N402 }
N201 [label="k: 00000094, b: 31"]
N201 -> { N403 N404 }
N202 [label="k: 00000096, b: 31"]
N202 -> { N405 N406 }
N203 [label="k: 00000098, b: 31"]
N203 -> { N407 N408 }
N204 [label="k: 0000009a, b: 31"]
N204 -> { N409 N410 }
N205 [label="k: 0000009c, b: 31"]
N205 -> { N411 N412 }
N206 [label="k: 0000009e, b: 31"]
N206 -> { N413 N414 }
N207 [label="k: 000000a0, b: 31"]
N207 -> { N415 N416 }
N208 [label="k: 000000a2, b: 31"]
N208 -> { N417 N418 }
N209 [label="k: 000000a4, b: 31"]
N209 -> { N419 N420 }
N210 [label="k: 000000a6, b: 31"]
N210 -> { N421 N422 }
N211 [label="k: 000000a8, b: 31"]
N211 -> { N423 N424 }
N212 [label="k: 000000aa, b: 31"]
N212 -> { N425 N426 }
N213 [label="k: 000000ac, b: 31"]
N213 -> { N427 N428 }
N214 [label="k: 000000ae, b: 31"]
N214 -> { N429 N430 }
N215 [label="k: 000000b0, b: 31"]
N215 -> { N431 N432 }
N216 [label="k: 000000b2, b: 31"]
N216 -> { N433 N434 }
N217 [label="k: 000000b4, b: 31"]
N217 -> { N435 N436 }
N218 [label="k: 000000b6, b: 31"]
N218 -> { N437 N438 }
N219 [label="k: 000000b8, b: 31"]
N219 -> { N439 N440 }
N220 [label="k: 000000ba, b: 31"]
N220 -> { N441 N442 }
N221 [label="k: 000000bc, b: 31"]
N221 -> { N443 N444 }
N222 [label="k: 000000be, b: 31"]
N222 -> { N445 N446 }
N223 [label="k: 000000c0, b: 31"]
N223 -> { N447 N448 }
N224 [label="k: 000000c2, b: 31"]
N224 -> { N449 N450 }
N225 [label="k: 000000c4, b: 31"]
N225 -> { N451 N452 }
N226 [label="k: 000000c6, b: 31"]
N226 -> { N453 N454 }
N227 [label="k: 000000c8, b: 31"]
N227 -> { N455 N456 }
N228 [label="k: 000000ca, b: 31"]
N228 -> { N457 N458 }
N229 [label="k: 000000cc, b: 31"]
N229 -> { N459 N460 }
N230 [label="k: 000000ce, b: 31"]
N230 -> { N461 N462 }
N231 [label="k: 000000d0, b: 31"]
N231 -> { N463 N464 }
N232 [label="k: 000000d2, b: 31"]
N232 -> { N465 N466 }
N233 [label="k: 000000d4, b: 31"]
N233 -> { N467 N468 }
N234 [label="k: 000000d6, b: 31"]
N234 -> { N469 N470 }
N235 [label="k: 000000d8, b: 31"]
N235 -> { N471 N472 }
N236 [label="k: 000000da, b: 31"]
N236 -> { N473 N474 }
N237 [label="k: 000000dc, b: 31"]
N237 -> { N475 N476 }
N238 [label="k: 000000de, b: 31"]
N238 -> { N477 N478 }
N239 [label="k: 000000e0, b: 31"]
N239 -> { N479 N480 }
N240 [label="k: 000000e2, b: 31"]
N240 -> { N481 N482 }
N241 [label="k: 000000e4, b: 31"]
N241 -> { N483 N484 }
N242 [label="k: 000000e6, b: 31"]
N242 -> { N485 N486 }
N243 [label="k: 000000e8, b: 31"]
N243 -> { N487 N488 }
N244 [label="k: 000000ea, b: 31"]
N244 -> { N489 N490 }
N245 [label="k: 000000ec, b: 31"]
N245 -> { N491 N492 }
N246 [label="k: 000000ee, b: 31"]
N246 -> { N493 N494 }
N247 [label="k: 000000f0, b: 31"]
N247 -> { N495 N496 }
N248 [label="k: 000000f2, b: 31"]
N248 -> { N497 N498 }
N249 [label="k: 000000f4, b: 31"]
N249 -> { N499 N500 }
N250 [label="k: 000000f6, b: 31"]
N250 -> { N501 N502 }
N251 [label="k: 000000f8, b: 31"]
N251 -> { N503 N504 }
N252 [label="k: 000000fa, b: 31"]
N252 -> { N505 N506 }
N253 [label="k: 000000fc, b: 31"]
N253 -> { N507 N508 }
N254 [label="k: 000000fe, b: 31"]
N254 -> { N509 N510 }
N255 [label="k: 00000000, b: 32, v: \"\"00\"\""]
N256 [label="k: 00000001, b: 32, v: \"\"01\"\""]
N257 [label="k: 00000002, b: 32, v: \"\"02\"\""]
N258 [label="k: 00000003, b: 32, v: \"\"03\"\""]
N259 [label="k: 00000004, b: 32, v: \"\"04\"\""]
N260 [label="k: 00000005, b: 32, v: \"\"05\"\""]
N261 [label="k: 00000006, b: 32, v: \"\"06\"\""]
N262 [label="k: 00000007, b: 32, v: \"\"07\"\""]
N263 [label="k: 00000008, b: 32, v: \"\"08\"\""]
N264 [label="k: 00000009, b: 32, v: \"\"09\"\""]
N265 [label="k: 0000000a, b: 32, v: \"\"0a\"\""]
N266 [label="k: 0000000b, b: 32, v: \"\"0b\"\""]
N267 [label="k: 0000000c, b: 32, v: \"\"0c\"\""]
N268 [label="k: 0000000d, b: 32, v: \"\"0d\"\""]
N269 [label="k: 0000000e, b: 32, v: \"\"0e\"\""]
N270 [label="k: 0000000f, b: 32, v: \"\"0f\"\""]
N271 [label="k: 00000010, b: 32, v: \"\"10\"\""]
N272 [label="k: 00000011, b: 32, v: \"\"11\"\""]
N273 [label="k: 00000012, b: 32, v: \"\"12\"\""]
N274 [label="k: 00000013, b: 32, v: \"\"13\"\""]
N275 [label="k: 00000014, b: 32, v: \"\"14\"\""]
N276 [label="k: 00000015, b: 32, v: \"\"15\"\""]
N277 [label="k: 00000016, b: 32, v: \"\"16\"\""]
N278 [label="k: 00000017, b: 32, v: \"\"17\"\""]
N279 [label="k: 00000018, b: 32, v: \"\"18\"\""]
N280 [label="k: 00000019, b: 32, v: \"\"19\"\""]
N281 [label="k: 0000001a, b: 32, v: \"\"1a\"\""]
N282 [label="k: 0000001b, b: 32, v: \"\"1b\"\""]
N283 [label="k: 0000001c, b: 32, v: \"\"1c\"\""]
N284 [label="k: 0000001d, b: 32, v: \"\"1d\"\""]
N285 [label="k: 0000001e, b: 32, v: \"\"1e\"\""]
N286 [label="k: 0000001f, b: 32, v: \"\"1f\"\""]
N287 [label="k: 00000020, b: 32, v: \"\"20\"\""]
N288 [label="k: 00000021, b: 32, v: \"\"21\"\""]
N289 [label="k: 00000022, b: 32, v: \"\"22\"\""]
N290 [label="k: 00000023, b: 32, v: \"\"23\"\""]
N291 [label="k: 00000024, b: 32, v: \"\"24\"\""]
N292 [label="k: 00000025, b: 32, v: \"\"25\"\""]
N293 [label="k: 00000026, b: 32, v: \"\"26\"\""]
N294 [label="k: 00000027, b: 32, v: \"\"27\"\""]
N295 [label="k: 00000028, b: 32, v: \"\"28\"\""]
N296 [label="k: 00000029, b: 32, v: \"\"29\"\""]
N297 [label="k: 0000002a, b: 32, v: \"\"2a\"\""]
N298 [label="k: 0000002b, b: 32, v: \"\"2b\"\""]
N299 [label="k: 0000002c, b: 32, v: \"\"2c\"\""]
N300 [label="k: 0000002d, b: 32, v: \"\"2d\"\""]
N301 [label="k: 0000002e, b: 32, v: \"\"2e\"\""]
N302 [label="k: 0000002f, b: 32, v: \"\"2f\"\""]
N303 [label="k: 00000030, b: 32, v: \"\"30\"\""]
N304 [label="k: 00000031, b: 32, v: \"\"31\"\""]
N305 [label="k: 00000032, b: 32, v: \"\"32\"\""]
N306 [label="k: 00000033, b: 32, v: \"\"33\"\""]
N307 [label="k: 00000034, b: 32, v: \"\"34\"\""]
N308 [label="k: 00000035, b: 32, v: \"\"35\"\""]
N309 [label="k: 00000036, b: 32, v: \"\"36\"\""]
N310 [label="k: 00000037, b: 32, v: \"\"37\"\""]
N311 [label="k: 00000038, b: 32, v: \"\"38\"\""]
N312 [label="k: 00000039, b: 32, v: \"\"39\"\""]
N313 [label="k: 0000003a, b: 32, v: \"\"3a\"\""]
N314 [label="k: 0000003b, b: 32, v: \"\"3b\"\""]
N315 [label="k: 0000003c, b: 32, v: \"\"3c\"\""]
N316 [label="k: 0000003d, b: 32, v: \"\"3d\"\""]
N317 [label="k: 0000003e, b: 32, v: \"\"3e\"\""]
N318 [label="k: 0000003f, b: 32, v: \"\"3f\"\""]
N319 [label="k: 00000040, b: 32, v: \"\"40\"\""]
N320 [label="k: 00000041, b: 32, v: \"\"41\"\""]
N321 [label="k: 00000042, b: 32, v: \"\"42\"\""]
N322 [label="k: 00000043, b: 32, v: \"\"43\"\""]
N323 [label="k: 00000044, b: 32, v: \"\"44\"\""]
N324 [label="k: 00000045, b: 32, v: \"\"45\"\""]
N325 [label="k: 00000046, b: 32, v: \"\"46\"\""]
N326 [label="k: 00000047, b: 32, v: \"\"47\"\""]
N327 [label="k: 00000048, b: 32, v: \"\"48\"\""]
N328 [label="k: 00000049, b: 32, v: \"\"49\"\""]
N329 [label="k: 0000004a, b: 32, v: \"\"4a\"\""]
N330 [label="k: 0000004b, b: 32, v: \"\"4b\"\""]
N331 [label="k: 0000004c, b: 32, v: \"\"4c\"\""]
N332 [label="k: 0000004d, b: 32, v: \"\"4d\"\""]
N333 [label="k: 0000004e, b: 32, v: \"\"4e\"\""]
N334 [label="k: 0000004f, b: 32, v: \"\"4f\"\""]
N335 [label="k: 00000050, b: 32, v: \"\"50\"\""]
N336 [label="k: 00000051, b: 32, v: \"\"51\"\""]
N337 [label="k: 00000052, b: 32, v: \"\"52\"\""]
N338 [label="k: 00000053, b: 32, v: \"\"53\"\""]
N339 [label="k: 00000054, b: 32, v: \"\"54\"\""]
N340 [label="k: 00000055, b: 32, v: \"\"55\"\""]
N341 [label="k: 00000056, b: 32, v: \"\"56\"\""]
N342 [label="k: 00000057, b: 32, v: \"\"57\"\""]
N343 [label="k: 00000058, b: 32, v: \"\"58\"\""]
N344 [label="k: 00000059, b: 32, v: \"\"59\"\""]
N345 [label="k: 0000005a, b: 32, v: \"\"5a\"\""]
N346 [label="k: 0000005b, b: 32, v: \"\"5b\"\""]
N347 [label="k: 0000005c, b: 32, v: \"\"5c\"\""]
N348 [label="k: 0000005d, b: 32, v: \"\"5d\"\""]
N349 [label="k: 0000005e, b: 32, v: \"\"5e\"\""]
N350 [label="k: 0000005f, b: 32, v: \"\"5f\"\""]
N351 [label="k: 00000060, b: 32, v: \"\"60\"\""]
N352 [label="k: 00000061, b: 32, v: \"\"61\"\""]
N353 [label="k: 00000062, b: 32, v: \"\"62\"\""]
N354 [label="k: 00000063, b: 32, v: \"\"63\"\""]
N355 [label="k: 00000064, b: 32, v: \"\"64\"\""]
N356 [label="k: 00000065, b: 32, v: \"\"65\"\""]
N357 [label="k: 00000066, b: 32, v: \"\"66\"\""]
N358 [label="k: 00000067, b: 32, v: \"\"67\"\""]
N359 [label="k: 00000068, b: 32, v: \"\"68\"\""]
N360 [label="k: 00000069, b: 32, v: \"\"69\"\""]
N361 [label="k: 0000006a, b: 32, v: \"\"6a\"\""]
N362 [label="k: 0000006b, b: 32, v: \"\"6b\"\""]
N363 [label="k: 0000006c, b: 32, v: \"\"6c\"\""]
N364 [label="k: 0000006d, b: 32, v: \"\"6d\"\""]
N365 [label="k: 0000006e, b: 32, v: \"\"6e\"\""]
N366 [label="k: 0000006f, b: 32, v: \"\"6f\"\""]
N367 [label="k: 00000070, b: 32, v: \"\"70\"\""]
N368 [label="k: 00000071, b: 32, v: \"\"71\"\""]
N369 [label="k: 00000072, b: 32, v: \"\"72\"\""]
N370 [label="k: 00000073, b: 32, v: \"\"73\"\""]
N371 [label="k: 00000074, b: 32, v: \"\"74\"\""]
N372 [label="k: 00000075, b: 32, v: \"\"75\"\""]
N373 [label="k: 00000076, b: 32, v: \"\"76\"\""]
N374 [label="k: 00000077, b: 32, v: \"\"77\"\""]
N375 [label="k: 00000078, b: 32, v: \"\"78\"\""]
N376 [label="k: 00000079, b: 32, v: \"\"79\"\""]
N377 [label="k: 0000007a, b: 32, v: \"\"7a\"\""]
N378 [label="k: 0000007b, b: 32, v: \"\"7b\"\""]
N379 [label="k: 0000007c, b: 32, v: \"\"7c\"\""]
N380 [label="k: 0000007d, b: 32, v: \"\"7d\"\""]
N381 [label="k: 0000007e, b: 32, v: \"\"7e\"\""]
N382 [label="k: 0000007f, b: 32, v: \"\"7f\"\""]
N383 [label="k: 00000080, b: 32, v: \"\"80\"\""]
N384 [label="k: 00000081, b: 32, v: \"\"81\"\""]
N385 [label="k: 00000082, b: 32, v: \"\"82\"\""]
N386 [label="k: 00000083, b: 32, v: \"\"83\"\""]
N387 [label="k: 00000084, b: 32, v: \"\"84\"\""]
N388 [label="k: 00000085, b: 32, v: \"\"85\"\""]
N389 [label="k: 00000086, b: 32, v: \"\"86\"\""]
N390 [label="k: 00000087, b: 32, v: \"\"87\"\""]
N391 [label="k: 00000088, b: 32, v: \"\"88\"\""]
N392 [label="k: 00000089, b: 32, v: \"\"89\"\""]
N393 [label="k: 0000008a, b: 32, v: \"\"8a\"\""]
N394 [label="k: 0000008b, b: 32, v: \"\"8b\"\""]
N395 [label="k: 0000008c, b: 32, v: \"\"8c\"\""]
N396 [label="k: 0000008d, b: 32, v: \"\"8d\"\""]
N397 [label="k: 0000008e, b: 32, v: \"\"8e\"\""]
N398 [label="k: 0000008f, b: 32, v: \"\"8f\"\""]
N399 [label="k: 00000090, b: 32, v: \"\"90\"\""]
N400 [label="k: 00000091, b: 32, v: \"\"91\"\""]
N401 [label="k: 00000092, b: 32, v: \"\"92\"\""]
N402 [label="k: 00000093, b: 32, v: \"\"93\"\""]
N403 [label="k: 00000094, b: 32, v: \"\"94\"\""]
N404 [label="k: 00000095, b: 32, v: \"\"95\"\""]
N405 [label="k: 00000096, b: 32, v: \"\"96\"\""]
N406 [label="k: 00000097, b: 32, v: \"\"97\"\""]
N407 [label="k: 00000098, b: 32, v: \"\"98\"\""]
N408 [label="k: 00000099, b: 32, v: \"\"99\"\""]
N409 [label="k: 0000009a, b: 32, v: \"\"9a\"\""]
N410 [label="k: 0000009b, b: 32, v: \"\"9b\"\""]
N411 [label="k: 0000009c, b: 32, v: \"\"9c\"\""]
N412 [label="k: 0000009d, b: 32, v: \"\"9d\"\""]
N413 [label="k: 0000009e, b: 32, v: \"\"9e\"\""]
N414 [label="k: 0000009f, b: 32, v: \"\"9f\"\""]
N415 [label="k: 000000a0, b: 32, v: \"\"a0\"\""]
N416 [label="k: 000000a1, b: 32, v: \"\"a1\"\""]
N417 [label="k: 000000a2, b: 32, v: \"\"a2\"\""]
N418 [label="k: 000000a3, b: 32, v: \"\"a3\"\""]
N419 [label="k: 000000a4, b: 32, v: \"\"a4\"\""]
N420 [label="k: 000000a5, b: 32, v: \"\"a5\"\""]
N421 [label="k: 000000a6, b: 32, v: \"\"a6\"\""]
N422 [label="k: 000000a7, b: 32, v: \"\"a7\"\""]
N423 [label="k: 000000a8, b: 32, v: \"\"a8\"\""]
N424 [label="k: 000000a9, b: 32, v: \"\"a9\"\""]
N425 [label="k: 000000aa, b: 32, v: \"\"aa\"\""]
N426 [label="k: 000000ab, b: 32, v: \"\"ab\"\""]
N427 [label="k: 000000ac, b: 32, v: \"\"ac\"\""]
N428 [label="k: 000000ad, b: 32, v: \"\"ad\"\""]
N429 [label="k: 000000ae, b: 32, v: \"\"ae\"\""]
N430 [label="k: 000000af, b: 32, v: \"\"af\"\""]
N431 [label="k: 000000b0, b: 32, v: \"\"b0\"\""]
N432 [label="k: 000000b1, b: 32, v: \"\"b1\"\""]
N433 [label="k: 000000b2, b: 32, v: \"\"b2\"\""]
N434 [label="k: 000000b3, b: 32, v: \"\"b3\"\""]
N435 [label="k: 000000b4, b: 32, v: \"\"b4\"\""]
N436 [label="k: 000000b5, b: 32, v: \"\"b5\"\""]
N437 [label="k: 000000b6, b: 32, v: \"\"b6\"\""]
N438 [label="k: 000000b7, b: 32, v: \"\"b7\"\""]
N439 [label="k: 000000b8, b: 32, v: \"\"b8\"\""]
N440 [label="k: 000000b9, b: 32, v: \"\"b9\"\""]
N441 [label="k: 000000ba, b: 32, v: \"\"ba\"\""]
N442 [label="k: 000000bb, b: 32, v: \"\"bb\"\""]
N443 [label="k: 000000bc, b: 32, v: \"\"bc\"\""]
N444 [label="k: 000000bd, b: 32, v: \"\"bd\"\""]
N445 [label="k: 000000be, b: 32, v: \"\"be\"\""]
N446 [label="k: 000000bf, b: 32, v: \"\"bf\"\""]
N447 [label="k: 000000c0, b: 32, v: \"\"c0\"\""]
N448 [label="k: 000000c1, b: 32, v: \"\"c1\"\""]
N449 [label="k: 000000c2, b: 32, v: \"\"c2\"\""]
N450 [label="k: 000000c3, b: 32, v: \"\"c3\"\""]
N451 [label="k: 000000c4, b: 32, v: \"\"c4\"\""]
N452 [label="k: 000000c5, b: 32, v: \"\"c5\"\""]
N453 [label="k: 000000c6, b: 32, v: \"\"c6\"\""]
N454 [label="k: 000000c7, b: 32, v: \"\"c7\"\""]
N455 [label="k: 000000c8, b: 32, v: \"\"c8\"\""]
N456 [label="k: 000000c9, b: 32, v: \"\"c9\"\""]
N457 [label="k: 000000ca, b: 32, v: \"\"ca\"\""]
N458 [label="k: 000000cb, b: 32, v: \"\"cb\"\""]
N459 [label="k: 000000cc, b: 32, v: \"\"cc\"\""]
N460 [label="k: 000000cd, b: 32, v: \"\"cd\"\""]
N461 [label="k: 000000ce, b: 32, v: \"\"ce\"\""]
N462 [label="k: 000000cf, b: 32, v: \"\"cf\"\""]
N463 [label="k: 000000d0, b: 32, v: \"\"d0\"\""]
N464 [label="k: 000000d1, b: 32, v: \"\"d1\"\""]
N465 [label="k: 000000d2, b: 32, v: \"\"d2\"\""]
N466 [label="k: 000000d3, b: 32, v: \"\"d3\"\""]
N467 [label="k: 000000d4, b: 32, v: \"\"d4\"\""]
N468 [label="k: 000000d5, b: 32, v: \"\"d5\"\""]
N469 [label="k: 000000d6, b: 32, v: \"\"d6\"\""]
N470 [label="k: 000000d7, b: 32, v: \"\"d7\"\""]
N471 [label="k: 000000d8, b: 32, v: \"\"d8\"\""]
N472 [label="k: 000000d9, b: 32, v: \"\"d9\"\""]
N473 [label="k: 000000da, b: 32, v: \"\"da\"\""]
N474 [label="k: 000000db, b: 32, v: \"\"db\"\""]
N475 [label="k: 000000dc, b: 32, v: \"\"dc\"\""]
N476 [label="k: 000000dd, b: 32, v: \"\"dd\"\""]
N477 [label="k: 000000de, b: 32, v: \"\"de\"\""]
N478 [label="k: 000000df, b: 32, v: \"\"df\"\""]
N479 [label="k: 000000e0, b: 32, v: \"\"e0\"\""]
N480 [label="k: 000000e1, b: 32, v: \"\"e1\"\""]
N481 [label="k: 000000e2, b: 32, v: \"\"e2\"\""]
N482 [label="k: 000000e3, b: 32, v: \"\"e3\"\""]
N483 [label="k: 000000e4, b: 32, v: \"\"e4\"\""]
N484 [label="k: 000000e5, b: 32, v: \"\"e5\"\""]
N485 [label="k: 000000e6, b: 32, v: \"\"e6\"\""]
N486 [label="k: 000000e7, b: 32, v: \"\"e7\"\""]
N487 [label="k: 000000e8, b: 32, v: \"\"e8\"\""]
N488 [label="k: 000000e9, b: 32, v: \"\"e9\"\""]
N489 [label="k: 000000ea, b: 32, v: \"\"ea\"\""]
N490 [label="k: 000000eb, b: 32, v: \"\"eb\"\""]
N491 [label="k: 000000ec, b: 32, v: \"\"ec\"\""]
N492 [label="k: 000000ed, b: 32, v: \"\"ed\"\""]
N493 [label="k: 000000ee, b: 32, v: \"\"ee\"\""]
N494 [label="k: 000000ef, b: 32, v: \"\"ef\"\""]
N495 [label="k: 000000f0, b: 32, v: \"\"f0\"\""]
N496 [label="k: 000000f1, b: 32, v: \"\"f1\"\""]
N497 [label="k: 000000f2, b: 32, v: \"\"f2\"\""]
N498 [label="k: 000000f3, b: 32, v: \"\"f3\"\""]
N499 [label="k: 000000f4, b: 32, v: \"\"f4\"\""]
N500 [label="k: 000000f5, b: 32, v: \"\"f5\"\""]
N501 [label="k: 000000f6, b: 32, v: \"\"f6\"\""]
N502 [label="k: 000000f7, b: 32, v: \"\"f7\"\""]
N503 [label="k: 000000f8, b: 32, v: \"\"f8\"\""]
N504 [label="k: 000000f9, b: 32, v: \"\"f9\"\""]
N505 [label="k: 000000fa, b: 32, v: \"\"fa\"\""]
N506 [label="k: 000000fb, b: 32, v: \"\"fb\"\""]
N507 [label="k: 000000fc, b: 32, v: \"\"fc\"\""]
N508 [label="k: 000000fd, b: 32, v: \"\"fd\"\""]
N509 [label="k: 000000fe, b: 32, v: \"\"fe\"\""]
N510 [label="k: 000000ff, b: 32, v: \"\"ff\"\""]
}
`

	TestTree32BigTreeInvertedInsertions = `digraph d {
N0 [label="k: 00000000, b: 0"]
N0 -> { N1 N2 }
N1 [label="k: 00000000, b: 1"]
N1 -> { N3 N4 }
N2 [label="k: 80000000, b: 1"]
N2 -> { N5 N6 }
N3 [label="k: 00000000, b: 2"]
N3 -> { N7 N8 }
N4 [label="k: 40000000, b: 2"]
N4 -> { N9 N10 }
N5 [label="k: 80000000, b: 2"]
N5 -> { N11 N12 }
N6 [label="k: c0000000, b: 2"]
N6 -> { N13 N14 }
N7 [label="k: 00000000, b: 3"]
N7 -> { N15 N16 }
N8 [label="k: 20000000, b: 3"]
N8 -> { N17 N18 }
N9 [label="k: 40000000, b: 3"]
N9 -> { N19 N20 }
N10 [label="k: 60000000, b: 3"]
N10 -> { N21 N22 }
N11 [label="k: 80000000, b: 3"]
N11 -> { N23 N24 }
N12 [label="k: a0000000, b: 3"]
N12 -> { N25 N26 }
N13 [label="k: c0000000, b: 3"]
N13 -> { N27 N28 }
N14 [label="k: e0000000, b: 3"]
N14 -> { N29 N30 }
N15 [label="k: 00000000, b: 4"]
N15 -> { N31 N32 }
N16 [label="k: 10000000, b: 4"]
N16 -> { N33 N34 }
N17 [label="k: 20000000, b: 4"]
N17 -> { N35 N36 }
N18 [label="k: 30000000, b: 4"]
N18 -> { N37 N38 }
N19 [label="k: 40000000, b: 4"]
N19 -> { N39 N40 }
N20 [label="k: 50000000, b: 4"]
N20 -> { N41 N42 }
N21 [label="k: 60000000, b: 4"]
N21 -> { N43 N44 }
N22 [label="k: 70000000, b: 4"]
N22 -> { N45 N46 }
N23 [label="k: 80000000, b: 4"]
N23 -> { N47 N48 }
N24 [label="k: 90000000, b: 4"]
N24 -> { N49 N50 }
N25 [label="k: a0000000, b: 4"]
N25 -> { N51 N52 }
N26 [label="k: b0000000, b: 4"]
N26 -> { N53 N54 }
N27 [label="k: c0000000, b: 4"]
N27 -> { N55 N56 }
N28 [label="k: d0000000, b: 4"]
N28 -> { N57 N58 }
N29 [label="k: e0000000, b: 4"]
N29 -> { N59 N60 }
N30 [label="k: f0000000, b: 4"]
N30 -> { N61 N62 }
N31 [label="k: 00000000, b: 5"]
N31 -> { N63 N64 }
N32 [label="k: 08000000, b: 5"]
N32 -> { N65 N66 }
N33 [label="k: 10000000, b: 5"]
N33 -> { N67 N68 }
N34 [label="k: 18000000, b: 5"]
N34 -> { N69 N70 }
N35 [label="k: 20000000, b: 5"]
N35 -> { N71 N72 }
N36 [label="k: 28000000, b: 5"]
N36 -> { N73 N74 }
N37 [label="k: 30000000, b: 5"]
N37 -> { N75 N76 }
N38 [label="k: 38000000, b: 5"]
N38 -> { N77 N78 }
N39 [label="k: 40000000, b: 5"]
N39 -> { N79 N80 }
N40 [label="k: 48000000, b: 5"]
N40 -> { N81 N82 }
N41 [label="k: 50000000, b: 5"]
N41 -> { N83 N84 }
N42 [label="k: 58000000, b: 5"]
N42 -> { N85 N86 }
N43 [label="k: 60000000, b: 5"]
N43 -> { N87 N88 }
N44 [label="k: 68000000, b: 5"]
N44 -> { N89 N90 }
N45 [label="k: 70000000, b: 5"]
N45 -> { N91 N92 }
N46 [label="k: 78000000, b: 5"]
N46 -> { N93 N94 }
N47 [label="k: 80000000, b: 5"]
N47 -> { N95 N96 }
N48 [label="k: 88000000, b: 5"]
N48 -> { N97 N98 }
N49 [label="k: 90000000, b: 5"]
N49 -> { N99 N100 }
N50 [label="k: 98000000, b: 5"]
N50 -> { N101 N102 }
N51 [label="k: a0000000, b: 5"]
N51 -> { N103 N104 }
N52 [label="k: a8000000, b: 5"]
N52 -> { N105 N106 }
N53 [label="k: b0000000, b: 5"]
N53 -> { N107 N108 }
N54 [label="k: b8000000, b: 5"]
N54 -> { N109 N110 }
N55 [label="k: c0000000, b: 5"]
N55 -> { N111 N112 }
N56 [label="k: c8000000, b: 5"]
N56 -> { N113 N114 }
N57 [label="k: d0000000, b: 5"]
N57 -> { N115 N116 }
N58 [label="k: d8000000, b: 5"]
N58 -> { N117 N118 }
N59 [label="k: e0000000, b: 5"]
N59 -> { N119 N120 }
N60 [label="k: e8000000, b: 5"]
N60 -> { N121 N122 }
N61 [label="k: f0000000, b: 5"]
N61 -> { N123 N124 }
N62 [label="k: f8000000, b: 5"]
N62 -> { N125 N126 }
N63 [label="k: 00000000, b: 6"]
N63 -> { N127 N128 }
N64 [label="k: 04000000, b: 6"]
N64 -> { N129 N130 }
N65 [label="k: 08000000, b: 6"]
N65 -> { N131 N132 }
N66 [label="k: 0c000000, b: 6"]
N66 -> { N133 N134 }
N67 [label="k: 10000000, b: 6"]
N67 -> { N135 N136 }
N68 [label="k: 14000000, b: 6"]
N68 -> { N137 N138 }
N69 [label="k: 18000000, b: 6"]
N69 -> { N139 N140 }
N70 [label="k: 1c000000, b: 6"]
N70 -> { N141 N142 }
N71 [label="k: 20000000, b: 6"]
N71 -> { N143 N144 }
N72 [label="k: 24000000, b: 6"]
N72 -> { N145 N146 }
N73 [label="k: 28000000, b: 6"]
N73 -> { N147 N148 }
N74 [label="k: 2c000000, b: 6"]
N74 -> { N149 N150 }
N75 [label="k: 30000000, b: 6"]
N75 -> { N151 N152 }
N76 [label="k: 34000000, b: 6"]
N76 -> { N153 N154 }
N77 [label="k: 38000000, b: 6"]
N77 -> { N155 N156 }
N78 [label="k: 3c000000, b: 6"]
N78 -> { N157 N158 }
N79 [label="k: 40000000, b: 6"]
N79 -> { N159 N160 }
N80 [label="k: 44000000, b: 6"]
N80 -> { N161 N162 }
N81 [label="k: 48000000, b: 6"]
N81 -> { N163 N164 }
N82 [label="k: 4c000000, b: 6"]
N82 -> { N165 N166 }
N83 [label="k: 50000000, b: 6"]
N83 -> { N167 N168 }
N84 [label="k: 54000000, b: 6"]
N84 -> { N169 N170 }
N85 [label="k: 58000000, b: 6"]
N85 -> { N171 N172 }
N86 [label="k: 5c000000, b: 6"]
N86 -> { N173 N174 }
N87 [label="k: 60000000, b: 6"]
N87 -> { N175 N176 }
N88 [label="k: 64000000, b: 6"]
N88 -> { N177 N178 }
N89 [label="k: 68000000, b: 6"]
N89 -> { N179 N180 }
N90 [label="k: 6c000000, b: 6"]
N90 -> { N181 N182 }
N91 [label="k: 70000000, b: 6"]
N91 -> { N183 N184 }
N92 [label="k: 74000000, b: 6"]
N92 -> { N185 N186 }
N93 [label="k: 78000000, b: 6"]
N93 -> { N187 N188 }
N94 [label="k: 7c000000, b: 6"]
N94 -> { N189 N190 }
N95 [label="k: 80000000, b: 6"]
N95 -> { N191 N192 }
N96 [label="k: 84000000, b: 6"]
N96 -> { N193 N194 }
N97 [label="k: 88000000, b: 6"]
N97 -> { N195 N196 }
N98 [label="k: 8c000000, b: 6"]
N98 -> { N197 N198 }
N99 [label="k: 90000000, b: 6"]
N99 -> { N199 N200 }
N100 [label="k: 94000000, b: 6"]
N100 -> { N201 N202 }
N101 [label="k: 98000000, b: 6"]
N101 -> { N203 N204 }
N102 [label="k: 9c000000, b: 6"]
N102 -> { N205 N206 }
N103 [label="k: a0000000, b: 6"]
N103 -> { N207 N208 }
N104 [label="k: a4000000, b: 6"]
N104 -> { N209 N210 }
N105 [label="k: a8000000, b: 6"]
N105 -> { N211 N212 }
N106 [label="k: ac000000, b: 6"]
N106 -> { N213 N214 }
N107 [label="k: b0000000, b: 6"]
N107 -> { N215 N216 }
N108 [label="k: b4000000, b: 6"]
N108 -> { N217 N218 }
N109 [label="k: b8000000, b: 6"]
N109 -> { N219 N220 }
N110 [label="k: bc000000, b: 6"]
N110 -> { N221 N222 }
N111 [label="k: c0000000, b: 6"]
N111 -> { N223 N224 }
N112 [label="k: c4000000, b: 6"]
N112 -> { N225 N226 }
N113 [label="k: c8000000, b: 6"]
N113 -> { N227 N228 }
N114 [label="k: cc000000, b: 6"]
N114 -> { N229 N230 }
N115 [label="k: d0000000, b: 6"]
N115 -> { N231 N232 }
N116 [label="k: d4000000, b: 6"]
N116 -> { N233 N234 }
N117 [label="k: d8000000, b: 6"]
N117 -> { N235 N236 }
N118 [label="k: dc000000, b: 6"]
N118 -> { N237 N238 }
N119 [label="k: e0000000, b: 6"]
N119 -> { N239 N240 }
N120 [label="k: e4000000, b: 6"]
N120 -> { N241 N242 }
N121 [label="k: e8000000, b: 6"]
N121 -> { N243 N244 }
N122 [label="k: ec000000, b: 6"]
N122 -> { N245 N246 }
N123 [label="k: f0000000, b: 6"]
N123 -> { N247 N248 }
N124 [label="k: f4000000, b: 6"]
N124 -> { N249 N250 }
N125 [label="k: f8000000, b: 6"]
N125 -> { N251 N252 }
N126 [label="k: fc000000, b: 6"]
N126 -> { N253 N254 }
N127 [label="k: 00000000, b: 7"]
N127 -> { N255 N256 }
N128 [label="k: 02000000, b: 7"]
N128 -> { N257 N258 }
N129 [label="k: 04000000, b: 7"]
N129 -> { N259 N260 }
N130 [label="k: 06000000, b: 7"]
N130 -> { N261 N262 }
N131 [label="k: 08000000, b: 7"]
N131 -> { N263 N264 }
N132 [label="k: 0a000000, b: 7"]
N132 -> { N265 N266 }
N133 [label="k: 0c000000, b: 7"]
N133 -> { N267 N268 }
N134 [label="k: 0e000000, b: 7"]
N134 -> { N269 N270 }
N135 [label="k: 10000000, b: 7"]
N135 -> { N271 N272 }
N136 [label="k: 12000000, b: 7"]
N136 -> { N273 N274 }
N137 [label="k: 14000000, b: 7"]
N137 -> { N275 N276 }
N138 [label="k: 16000000, b: 7"]
N138 -> { N277 N278 }
N139 [label="k: 18000000, b: 7"]
N139 -> { N279 N280 }
N140 [label="k: 1a000000, b: 7"]
N140 -> { N281 N282 }
N141 [label="k: 1c000000, b: 7"]
N141 -> { N283 N284 }
N142 [label="k: 1e000000, b: 7"]
N142 -> { N285 N286 }
N143 [label="k: 20000000, b: 7"]
N143 -> { N287 N288 }
N144 [label="k: 22000000, b: 7"]
N144 -> { N289 N290 }
N145 [label="k: 24000000, b: 7"]
N145 -> { N291 N292 }
N146 [label="k: 26000000, b: 7"]
N146 -> { N293 N294 }
N147 [label="k: 28000000, b: 7"]
N147 -> { N295 N296 }
N148 [label="k: 2a000000, b: 7"]
N148 -> { N297 N298 }
N149 [label="k: 2c000000, b: 7"]
N149 -> { N299 N300 }
N150 [label="k: 2e000000, b: 7"]
N150 -> { N301 N302 }
N151 [label="k: 30000000, b: 7"]
N151 -> { N303 N304 }
N152 [label="k: 32000000, b: 7"]
N152 -> { N305 N306 }
N153 [label="k: 34000000, b: 7"]
N153 -> { N307 N308 }
N154 [label="k: 36000000, b: 7"]
N154 -> { N309 N310 }
N155 [label="k: 38000000, b: 7"]
N155 -> { N311 N312 }
N156 [label="k: 3a000000, b: 7"]
N156 -> { N313 N314 }
N157 [label="k: 3c000000, b: 7"]
N157 -> { N315 N316 }
N158 [label="k: 3e000000, b: 7"]
N158 -> { N317 N318 }
N159 [label="k: 40000000, b: 7"]
N159 -> { N319 N320 }
N160 [label="k: 42000000, b: 7"]
N160 -> { N321 N322 }
N161 [label="k: 44000000, b: 7"]
N161 -> { N323 N324 }
N162 [label="k: 46000000, b: 7"]
N162 -> { N325 N326 }
N163 [label="k: 48000000, b: 7"]
N163 -> { N327 N328 }
N164 [label="k: 4a000000, b: 7"]
N164 -> { N329 N330 }
N165 [label="k: 4c000000, b: 7"]
N165 -> { N331 N332 }
N166 [label="k: 4e000000, b: 7"]
N166 -> { N333 N334 }
N167 [label="k: 50000000, b: 7"]
N167 -> { N335 N336 }
N168 [label="k: 52000000, b: 7"]
N168 -> { N337 N338 }
N169 [label="k: 54000000, b: 7"]
N169 -> { N339 N340 }
N170 [label="k: 56000000, b: 7"]
N170 -> { N341 N342 }
N171 [label="k: 58000000, b: 7"]
N171 -> { N343 N344 }
N172 [label="k: 5a000000, b: 7"]
N172 -> { N345 N346 }
N173 [label="k: 5c000000, b: 7"]
N173 -> { N347 N348 }
N174 [label="k: 5e000000, b: 7"]
N174 -> { N349 N350 }
N175 [label="k: 60000000, b: 7"]
N175 -> { N351 N352 }
N176 [label="k: 62000000, b: 7"]
N176 -> { N353 N354 }
N177 [label="k: 64000000, b: 7"]
N177 -> { N355 N356 }
N178 [label="k: 66000000, b: 7"]
N178 -> { N357 N358 }
N179 [label="k: 68000000, b: 7"]
N179 -> { N359 N360 }
N180 [label="k: 6a000000, b: 7"]
N180 -> { N361 N362 }
N181 [label="k: 6c000000, b: 7"]
N181 -> { N363 N364 }
N182 [label="k: 6e000000, b: 7"]
N182 -> { N365 N366 }
N183 [label="k: 70000000, b: 7"]
N183 -> { N367 N368 }
N184 [label="k: 72000000, b: 7"]
N184 -> { N369 N370 }
N185 [label="k: 74000000, b: 7"]
N185 -> { N371 N372 }
N186 [label="k: 76000000, b: 7"]
N186 -> { N373 N374 }
N187 [label="k: 78000000, b: 7"]
N187 -> { N375 N376 }
N188 [label="k: 7a000000, b: 7"]
N188 -> { N377 N378 }
N189 [label="k: 7c000000, b: 7"]
N189 -> { N379 N380 }
N190 [label="k: 7e000000, b: 7"]
N190 -> { N381 N382 }
N191 [label="k: 80000000, b: 7"]
N191 -> { N383 N384 }
N192 [label="k: 82000000, b: 7"]
N192 -> { N385 N386 }
N193 [label="k: 84000000, b: 7"]
N193 -> { N387 N388 }
N194 [label="k: 86000000, b: 7"]
N194 -> { N389 N390 }
N195 [label="k: 88000000, b: 7"]
N195 -> { N391 N392 }
N196 [label="k: 8a000000, b: 7"]
N196 -> { N393 N394 }
N197 [label="k: 8c000000, b: 7"]
N197 -> { N395 N396 }
N198 [label="k: 8e000000, b: 7"]
N198 -> { N397 N398 }
N199 [label="k: 90000000, b: 7"]
N199 -> { N399 N400 }
N200 [label="k: 92000000, b: 7"]
N200 -> { N401 N402 }
N201 [label="k: 94000000, b: 7"]
N201 -> { N403 N404 }
N202 [label="k: 96000000, b: 7"]
N202 -> { N405 N406 }
N203 [label="k: 98000000, b: 7"]
N203 -> { N407 N408 }
N204 [label="k: 9a000000, b: 7"]
N204 -> { N409 N410 }
N205 [label="k: 9c000000, b: 7"]
N205 -> { N411 N412 }
N206 [label="k: 9e000000, b: 7"]
N206 -> { N413 N414 }
N207 [label="k: a0000000, b: 7"]
N207 -> { N415 N416 }
N208 [label="k: a2000000, b: 7"]
N208 -> { N417 N418 }
N209 [label="k: a4000000, b: 7"]
N209 -> { N419 N420 }
N210 [label="k: a6000000, b: 7"]
N210 -> { N421 N422 }
N211 [label="k: a8000000, b: 7"]
N211 -> { N423 N424 }
N212 [label="k: aa000000, b: 7"]
N212 -> { N425 N426 }
N213 [label="k: ac000000, b: 7"]
N213 -> { N427 N428 }
N214 [label="k: ae000000, b: 7"]
N214 -> { N429 N430 }
N215 [label="k: b0000000, b: 7"]
N215 -> { N431 N432 }
N216 [label="k: b2000000, b: 7"]
N216 -> { N433 N434 }
N217 [label="k: b4000000, b: 7"]
N217 -> { N435 N436 }
N218 [label="k: b6000000, b: 7"]
N218 -> { N437 N438 }
N219 [label="k: b8000000, b: 7"]
N219 -> { N439 N440 }
N220 [label="k: ba000000, b: 7"]
N220 -> { N441 N442 }
N221 [label="k: bc000000, b: 7"]
N221 -> { N443 N444 }
N222 [label="k: be000000, b: 7"]
N222 -> { N445 N446 }
N223 [label="k: c0000000, b: 7"]
N223 -> { N447 N448 }
N224 [label="k: c2000000, b: 7"]
N224 -> { N449 N450 }
N225 [label="k: c4000000, b: 7"]
N225 -> { N451 N452 }
N226 [label="k: c6000000, b: 7"]
N226 -> { N453 N454 }
N227 [label="k: c8000000, b: 7"]
N227 -> { N455 N456 }
N228 [label="k: ca000000, b: 7"]
N228 -> { N457 N458 }
N229 [label="k: cc000000, b: 7"]
N229 -> { N459 N460 }
N230 [label="k: ce000000, b: 7"]
N230 -> { N461 N462 }
N231 [label="k: d0000000, b: 7"]
N231 -> { N463 N464 }
N232 [label="k: d2000000, b: 7"]
N232 -> { N465 N466 }
N233 [label="k: d4000000, b: 7"]
N233 -> { N467 N468 }
N234 [label="k: d6000000, b: 7"]
N234 -> { N469 N470 }
N235 [label="k: d8000000, b: 7"]
N235 -> { N471 N472 }
N236 [label="k: da000000, b: 7"]
N236 -> { N473 N474 }
N237 [label="k: dc000000, b: 7"]
N237 -> { N475 N476 }
N238 [label="k: de000000, b: 7"]
N238 -> { N477 N478 }
N239 [label="k: e0000000, b: 7"]
N239 -> { N479 N480 }
N240 [label="k: e2000000, b: 7"]
N240 -> { N481 N482 }
N241 [label="k: e4000000, b: 7"]
N241 -> { N483 N484 }
N242 [label="k: e6000000, b: 7"]
N242 -> { N485 N486 }
N243 [label="k: e8000000, b: 7"]
N243 -> { N487 N488 }
N244 [label="k: ea000000, b: 7"]
N244 -> { N489 N490 }
N245 [label="k: ec000000, b: 7"]
N245 -> { N491 N492 }
N246 [label="k: ee000000, b: 7"]
N246 -> { N493 N494 }
N247 [label="k: f0000000, b: 7"]
N247 -> { N495 N496 }
N248 [label="k: f2000000, b: 7"]
N248 -> { N497 N498 }
N249 [label="k: f4000000, b: 7"]
N249 -> { N499 N500 }
N250 [label="k: f6000000, b: 7"]
N250 -> { N501 N502 }
N251 [label="k: f8000000, b: 7"]
N251 -> { N503 N504 }
N252 [label="k: fa000000, b: 7"]
N252 -> { N505 N506 }
N253 [label="k: fc000000, b: 7"]
N253 -> { N507 N508 }
N254 [label="k: fe000000, b: 7"]
N254 -> { N509 N510 }
N255 [label="k: 00000000, b: 32, v: \"\"00\"\""]
N256 [label="k: 01000000, b: 32, v: \"\"01\"\""]
N257 [label="k: 02000000, b: 32, v: \"\"02\"\""]
N258 [label="k: 03000000, b: 32, v: \"\"03\"\""]
N259 [label="k: 04000000, b: 32, v: \"\"04\"\""]
N260 [label="k: 05000000, b: 32, v: \"\"05\"\""]
N261 [label="k: 06000000, b: 32, v: \"\"06\"\""]
N262 [label="k: 07000000, b: 32, v: \"\"07\"\""]
N263 [label="k: 08000000, b: 32, v: \"\"08\"\""]
N264 [label="k: 09000000, b: 32, v: \"\"09\"\""]
N265 [label="k: 0a000000, b: 32, v: \"\"0a\"\""]
N266 [label="k: 0b000000, b: 32, v: \"\"0b\"\""]
N267 [label="k: 0c000000, b: 32, v: \"\"0c\"\""]
N268 [label="k: 0d000000, b: 32, v: \"\"0d\"\""]
N269 [label="k: 0e000000, b: 32, v: \"\"0e\"\""]
N270 [label="k: 0f000000, b: 32, v: \"\"0f\"\""]
N271 [label="k: 10000000, b: 32, v: \"\"10\"\""]
N272 [label="k: 11000000, b: 32, v: \"\"11\"\""]
N273 [label="k: 12000000, b: 32, v: \"\"12\"\""]
N274 [label="k: 13000000, b: 32, v: \"\"13\"\""]
N275 [label="k: 14000000, b: 32, v: \"\"14\"\""]
N276 [label="k: 15000000, b: 32, v: \"\"15\"\""]
N277 [label="k: 16000000, b: 32, v: \"\"16\"\""]
N278 [label="k: 17000000, b: 32, v: \"\"17\"\""]
N279 [label="k: 18000000, b: 32, v: \"\"18\"\""]
N280 [label="k: 19000000, b: 32, v: \"\"19\"\""]
N281 [label="k: 1a000000, b: 32, v: \"\"1a\"\""]
N282 [label="k: 1b000000, b: 32, v: \"\"1b\"\""]
N283 [label="k: 1c000000, b: 32, v: \"\"1c\"\""]
N284 [label="k: 1d000000, b: 32, v: \"\"1d\"\""]
N285 [label="k: 1e000000, b: 32, v: \"\"1e\"\""]
N286 [label="k: 1f000000, b: 32, v: \"\"1f\"\""]
N287 [label="k: 20000000, b: 32, v: \"\"20\"\""]
N288 [label="k: 21000000, b: 32, v: \"\"21\"\""]
N289 [label="k: 22000000, b: 32, v: \"\"22\"\""]
N290 [label="k: 23000000, b: 32, v: \"\"23\"\""]
N291 [label="k: 24000000, b: 32, v: \"\"24\"\""]
N292 [label="k: 25000000, b: 32, v: \"\"25\"\""]
N293 [label="k: 26000000, b: 32, v: \"\"26\"\""]
N294 [label="k: 27000000, b: 32, v: \"\"27\"\""]
N295 [label="k: 28000000, b: 32, v: \"\"28\"\""]
N296 [label="k: 29000000, b: 32, v: \"\"29\"\""]
N297 [label="k: 2a000000, b: 32, v: \"\"2a\"\""]
N298 [label="k: 2b000000, b: 32, v: \"\"2b\"\""]
N299 [label="k: 2c000000, b: 32, v: \"\"2c\"\""]
N300 [label="k: 2d000000, b: 32, v: \"\"2d\"\""]
N301 [label="k: 2e000000, b: 32, v: \"\"2e\"\""]
N302 [label="k: 2f000000, b: 32, v: \"\"2f\"\""]
N303 [label="k: 30000000, b: 32, v: \"\"30\"\""]
N304 [label="k: 31000000, b: 32, v: \"\"31\"\""]
N305 [label="k: 32000000, b: 32, v: \"\"32\"\""]
N306 [label="k: 33000000, b: 32, v: \"\"33\"\""]
N307 [label="k: 34000000, b: 32, v: \"\"34\"\""]
N308 [label="k: 35000000, b: 32, v: \"\"35\"\""]
N309 [label="k: 36000000, b: 32, v: \"\"36\"\""]
N310 [label="k: 37000000, b: 32, v: \"\"37\"\""]
N311 [label="k: 38000000, b: 32, v: \"\"38\"\""]
N312 [label="k: 39000000, b: 32, v: \"\"39\"\""]
N313 [label="k: 3a000000, b: 32, v: \"\"3a\"\""]
N314 [label="k: 3b000000, b: 32, v: \"\"3b\"\""]
N315 [label="k: 3c000000, b: 32, v: \"\"3c\"\""]
N316 [label="k: 3d000000, b: 32, v: \"\"3d\"\""]
N317 [label="k: 3e000000, b: 32, v: \"\"3e\"\""]
N318 [label="k: 3f000000, b: 32, v: \"\"3f\"\""]
N319 [label="k: 40000000, b: 32, v: \"\"40\"\""]
N320 [label="k: 41000000, b: 32, v: \"\"41\"\""]
N321 [label="k: 42000000, b: 32, v: \"\"42\"\""]
N322 [label="k: 43000000, b: 32, v: \"\"43\"\""]
N323 [label="k: 44000000, b: 32, v: \"\"44\"\""]
N324 [label="k: 45000000, b: 32, v: \"\"45\"\""]
N325 [label="k: 46000000, b: 32, v: \"\"46\"\""]
N326 [label="k: 47000000, b: 32, v: \"\"47\"\""]
N327 [label="k: 48000000, b: 32, v: \"\"48\"\""]
N328 [label="k: 49000000, b: 32, v: \"\"49\"\""]
N329 [label="k: 4a000000, b: 32, v: \"\"4a\"\""]
N330 [label="k: 4b000000, b: 32, v: \"\"4b\"\""]
N331 [label="k: 4c000000, b: 32, v: \"\"4c\"\""]
N332 [label="k: 4d000000, b: 32, v: \"\"4d\"\""]
N333 [label="k: 4e000000, b: 32, v: \"\"4e\"\""]
N334 [label="k: 4f000000, b: 32, v: \"\"4f\"\""]
N335 [label="k: 50000000, b: 32, v: \"\"50\"\""]
N336 [label="k: 51000000, b: 32, v: \"\"51\"\""]
N337 [label="k: 52000000, b: 32, v: \"\"52\"\""]
N338 [label="k: 53000000, b: 32, v: \"\"53\"\""]
N339 [label="k: 54000000, b: 32, v: \"\"54\"\""]
N340 [label="k: 55000000, b: 32, v: \"\"55\"\""]
N341 [label="k: 56000000, b: 32, v: \"\"56\"\""]
N342 [label="k: 57000000, b: 32, v: \"\"57\"\""]
N343 [label="k: 58000000, b: 32, v: \"\"58\"\""]
N344 [label="k: 59000000, b: 32, v: \"\"59\"\""]
N345 [label="k: 5a000000, b: 32, v: \"\"5a\"\""]
N346 [label="k: 5b000000, b: 32, v: \"\"5b\"\""]
N347 [label="k: 5c000000, b: 32, v: \"\"5c\"\""]
N348 [label="k: 5d000000, b: 32, v: \"\"5d\"\""]
N349 [label="k: 5e000000, b: 32, v: \"\"5e\"\""]
N350 [label="k: 5f000000, b: 32, v: \"\"5f\"\""]
N351 [label="k: 60000000, b: 32, v: \"\"60\"\""]
N352 [label="k: 61000000, b: 32, v: \"\"61\"\""]
N353 [label="k: 62000000, b: 32, v: \"\"62\"\""]
N354 [label="k: 63000000, b: 32, v: \"\"63\"\""]
N355 [label="k: 64000000, b: 32, v: \"\"64\"\""]
N356 [label="k: 65000000, b: 32, v: \"\"65\"\""]
N357 [label="k: 66000000, b: 32, v: \"\"66\"\""]
N358 [label="k: 67000000, b: 32, v: \"\"67\"\""]
N359 [label="k: 68000000, b: 32, v: \"\"68\"\""]
N360 [label="k: 69000000, b: 32, v: \"\"69\"\""]
N361 [label="k: 6a000000, b: 32, v: \"\"6a\"\""]
N362 [label="k: 6b000000, b: 32, v: \"\"6b\"\""]
N363 [label="k: 6c000000, b: 32, v: \"\"6c\"\""]
N364 [label="k: 6d000000, b: 32, v: \"\"6d\"\""]
N365 [label="k: 6e000000, b: 32, v: \"\"6e\"\""]
N366 [label="k: 6f000000, b: 32, v: \"\"6f\"\""]
N367 [label="k: 70000000, b: 32, v: \"\"70\"\""]
N368 [label="k: 71000000, b: 32, v: \"\"71\"\""]
N369 [label="k: 72000000, b: 32, v: \"\"72\"\""]
N370 [label="k: 73000000, b: 32, v: \"\"73\"\""]
N371 [label="k: 74000000, b: 32, v: \"\"74\"\""]
N372 [label="k: 75000000, b: 32, v: \"\"75\"\""]
N373 [label="k: 76000000, b: 32, v: \"\"76\"\""]
N374 [label="k: 77000000, b: 32, v: \"\"77\"\""]
N375 [label="k: 78000000, b: 32, v: \"\"78\"\""]
N376 [label="k: 79000000, b: 32, v: \"\"79\"\""]
N377 [label="k: 7a000000, b: 32, v: \"\"7a\"\""]
N378 [label="k: 7b000000, b: 32, v: \"\"7b\"\""]
N379 [label="k: 7c000000, b: 32, v: \"\"7c\"\""]
N380 [label="k: 7d000000, b: 32, v: \"\"7d\"\""]
N381 [label="k: 7e000000, b: 32, v: \"\"7e\"\""]
N382 [label="k: 7f000000, b: 32, v: \"\"7f\"\""]
N383 [label="k: 80000000, b: 32, v: \"\"80\"\""]
N384 [label="k: 81000000, b: 32, v: \"\"81\"\""]
N385 [label="k: 82000000, b: 32, v: \"\"82\"\""]
N386 [label="k: 83000000, b: 32, v: \"\"83\"\""]
N387 [label="k: 84000000, b: 32, v: \"\"84\"\""]
N388 [label="k: 85000000, b: 32, v: \"\"85\"\""]
N389 [label="k: 86000000, b: 32, v: \"\"86\"\""]
N390 [label="k: 87000000, b: 32, v: \"\"87\"\""]
N391 [label="k: 88000000, b: 32, v: \"\"88\"\""]
N392 [label="k: 89000000, b: 32, v: \"\"89\"\""]
N393 [label="k: 8a000000, b: 32, v: \"\"8a\"\""]
N394 [label="k: 8b000000, b: 32, v: \"\"8b\"\""]
N395 [label="k: 8c000000, b: 32, v: \"\"8c\"\""]
N396 [label="k: 8d000000, b: 32, v: \"\"8d\"\""]
N397 [label="k: 8e000000, b: 32, v: \"\"8e\"\""]
N398 [label="k: 8f000000, b: 32, v: \"\"8f\"\""]
N399 [label="k: 90000000, b: 32, v: \"\"90\"\""]
N400 [label="k: 91000000, b: 32, v: \"\"91\"\""]
N401 [label="k: 92000000, b: 32, v: \"\"92\"\""]
N402 [label="k: 93000000, b: 32, v: \"\"93\"\""]
N403 [label="k: 94000000, b: 32, v: \"\"94\"\""]
N404 [label="k: 95000000, b: 32, v: \"\"95\"\""]
N405 [label="k: 96000000, b: 32, v: \"\"96\"\""]
N406 [label="k: 97000000, b: 32, v: \"\"97\"\""]
N407 [label="k: 98000000, b: 32, v: \"\"98\"\""]
N408 [label="k: 99000000, b: 32, v: \"\"99\"\""]
N409 [label="k: 9a000000, b: 32, v: \"\"9a\"\""]
N410 [label="k: 9b000000, b: 32, v: \"\"9b\"\""]
N411 [label="k: 9c000000, b: 32, v: \"\"9c\"\""]
N412 [label="k: 9d000000, b: 32, v: \"\"9d\"\""]
N413 [label="k: 9e000000, b: 32, v: \"\"9e\"\""]
N414 [label="k: 9f000000, b: 32, v: \"\"9f\"\""]
N415 [label="k: a0000000, b: 32, v: \"\"a0\"\""]
N416 [label="k: a1000000, b: 32, v: \"\"a1\"\""]
N417 [label="k: a2000000, b: 32, v: \"\"a2\"\""]
N418 [label="k: a3000000, b: 32, v: \"\"a3\"\""]
N419 [label="k: a4000000, b: 32, v: \"\"a4\"\""]
N420 [label="k: a5000000, b: 32, v: \"\"a5\"\""]
N421 [label="k: a6000000, b: 32, v: \"\"a6\"\""]
N422 [label="k: a7000000, b: 32, v: \"\"a7\"\""]
N423 [label="k: a8000000, b: 32, v: \"\"a8\"\""]
N424 [label="k: a9000000, b: 32, v: \"\"a9\"\""]
N425 [label="k: aa000000, b: 32, v: \"\"aa\"\""]
N426 [label="k: ab000000, b: 32, v: \"\"ab\"\""]
N427 [label="k: ac000000, b: 32, v: \"\"ac\"\""]
N428 [label="k: ad000000, b: 32, v: \"\"ad\"\""]
N429 [label="k: ae000000, b: 32, v: \"\"ae\"\""]
N430 [label="k: af000000, b: 32, v: \"\"af\"\""]
N431 [label="k: b0000000, b: 32, v: \"\"b0\"\""]
N432 [label="k: b1000000, b: 32, v: \"\"b1\"\""]
N433 [label="k: b2000000, b: 32, v: \"\"b2\"\""]
N434 [label="k: b3000000, b: 32, v: \"\"b3\"\""]
N435 [label="k: b4000000, b: 32, v: \"\"b4\"\""]
N436 [label="k: b5000000, b: 32, v: \"\"b5\"\""]
N437 [label="k: b6000000, b: 32, v: \"\"b6\"\""]
N438 [label="k: b7000000, b: 32, v: \"\"b7\"\""]
N439 [label="k: b8000000, b: 32, v: \"\"b8\"\""]
N440 [label="k: b9000000, b: 32, v: \"\"b9\"\""]
N441 [label="k: ba000000, b: 32, v: \"\"ba\"\""]
N442 [label="k: bb000000, b: 32, v: \"\"bb\"\""]
N443 [label="k: bc000000, b: 32, v: \"\"bc\"\""]
N444 [label="k: bd000000, b: 32, v: \"\"bd\"\""]
N445 [label="k: be000000, b: 32, v: \"\"be\"\""]
N446 [label="k: bf000000, b: 32, v: \"\"bf\"\""]
N447 [label="k: c0000000, b: 32, v: \"\"c0\"\""]
N448 [label="k: c1000000, b: 32, v: \"\"c1\"\""]
N449 [label="k: c2000000, b: 32, v: \"\"c2\"\""]
N450 [label="k: c3000000, b: 32, v: \"\"c3\"\""]
N451 [label="k: c4000000, b: 32, v: \"\"c4\"\""]
N452 [label="k: c5000000, b: 32, v: \"\"c5\"\""]
N453 [label="k: c6000000, b: 32, v: \"\"c6\"\""]
N454 [label="k: c7000000, b: 32, v: \"\"c7\"\""]
N455 [label="k: c8000000, b: 32, v: \"\"c8\"\""]
N456 [label="k: c9000000, b: 32, v: \"\"c9\"\""]
N457 [label="k: ca000000, b: 32, v: \"\"ca\"\""]
N458 [label="k: cb000000, b: 32, v: \"\"cb\"\""]
N459 [label="k: cc000000, b: 32, v: \"\"cc\"\""]
N460 [label="k: cd000000, b: 32, v: \"\"cd\"\""]
N461 [label="k: ce000000, b: 32, v: \"\"ce\"\""]
N462 [label="k: cf000000, b: 32, v: \"\"cf\"\""]
N463 [label="k: d0000000, b: 32, v: \"\"d0\"\""]
N464 [label="k: d1000000, b: 32, v: \"\"d1\"\""]
N465 [label="k: d2000000, b: 32, v: \"\"d2\"\""]
N466 [label="k: d3000000, b: 32, v: \"\"d3\"\""]
N467 [label="k: d4000000, b: 32, v: \"\"d4\"\""]
N468 [label="k: d5000000, b: 32, v: \"\"d5\"\""]
N469 [label="k: d6000000, b: 32, v: \"\"d6\"\""]
N470 [label="k: d7000000, b: 32, v: \"\"d7\"\""]
N471 [label="k: d8000000, b: 32, v: \"\"d8\"\""]
N472 [label="k: d9000000, b: 32, v: \"\"d9\"\""]
N473 [label="k: da000000, b: 32, v: \"\"da\"\""]
N474 [label="k: db000000, b: 32, v: \"\"db\"\""]
N475 [label="k: dc000000, b: 32, v: \"\"dc\"\""]
N476 [label="k: dd000000, b: 32, v: \"\"dd\"\""]
N477 [label="k: de000000, b: 32, v: \"\"de\"\""]
N478 [label="k: df000000, b: 32, v: \"\"df\"\""]
N479 [label="k: e0000000, b: 32, v: \"\"e0\"\""]
N480 [label="k: e1000000, b: 32, v: \"\"e1\"\""]
N481 [label="k: e2000000, b: 32, v: \"\"e2\"\""]
N482 [label="k: e3000000, b: 32, v: \"\"e3\"\""]
N483 [label="k: e4000000, b: 32, v: \"\"e4\"\""]
N484 [label="k: e5000000, b: 32, v: \"\"e5\"\""]
N485 [label="k: e6000000, b: 32, v: \"\"e6\"\""]
N486 [label="k: e7000000, b: 32, v: \"\"e7\"\""]
N487 [label="k: e8000000, b: 32, v: \"\"e8\"\""]
N488 [label="k: e9000000, b: 32, v: \"\"e9\"\""]
N489 [label="k: ea000000, b: 32, v: \"\"ea\"\""]
N490 [label="k: eb000000, b: 32, v: \"\"eb\"\""]
N491 [label="k: ec000000, b: 32, v: \"\"ec\"\""]
N492 [label="k: ed000000, b: 32, v: \"\"ed\"\""]
N493 [label="k: ee000000, b: 32, v: \"\"ee\"\""]
N494 [label="k: ef000000, b: 32, v: \"\"ef\"\""]
N495 [label="k: f0000000, b: 32, v: \"\"f0\"\""]
N496 [label="k: f1000000, b: 32, v: \"\"f1\"\""]
N497 [label="k: f2000000, b: 32, v: \"\"f2\"\""]
N498 [label="k: f3000000, b: 32, v: \"\"f3\"\""]
N499 [label="k: f4000000, b: 32, v: \"\"f4\"\""]
N500 [label="k: f5000000, b: 32, v: \"\"f5\"\""]
N501 [label="k: f6000000, b: 32, v: \"\"f6\"\""]
N502 [label="k: f7000000, b: 32, v: \"\"f7\"\""]
N503 [label="k: f8000000, b: 32, v: \"\"f8\"\""]
N504 [label="k: f9000000, b: 32, v: \"\"f9\"\""]
N505 [label="k: fa000000, b: 32, v: \"\"fa\"\""]
N506 [label="k: fb000000, b: 32, v: \"\"fb\"\""]
N507 [label="k: fc000000, b: 32, v: \"\"fc\"\""]
N508 [label="k: fd000000, b: 32, v: \"\"fd\"\""]
N509 [label="k: fe000000, b: 32, v: \"\"fe\"\""]
N510 [label="k: ff000000, b: 32, v: \"\"ff\"\""]
}
`

	TestTree32EmptyTree = `digraph d {
N0 [label="nil"]
}
`

	TestTree32WithDeletedChildNode = `digraph d {
N0 [label="k: aaaaaaaa, b: 9, v: \"\"test\"\""]
}
`

	TestTree32WithDeletedChildAndNonLeafNodes = `digraph d {
N0 [label="k: a8000000, b: 6"]
N0 -> { N1 N2 }
N1 [label="k: a8aaaaaa, b: 9, v: \"\"L2.1\"\""]
N2 [label="k: aaaaaaaa, b: 7, v: \"\"L1\"\""]
N2 -> { N3 N4 }
N3 [label="k: aaaa0000, b: 15"]
N3 -> { N5 N6 }
N4 [label="k: abaaaaaa, b: 9, v: \"\"L2.2\"\""]
N5 [label="k: aaaaaaaa, b: 18, v: \"\"L3\"\""]
N6 [label="k: aaabaaaa, b: 24, v: \"\"L5\"\""]
}
`

	TestTree32WithDeletedTwoChildrenAndNonLeafNodes = `digraph d {
N0 [label="k: a8000000, b: 6"]
N0 -> { N1 N2 }
N1 [label="k: a8aaaaaa, b: 9, v: \"\"L2.1\"\""]
N2 [label="k: aaaaaaaa, b: 7, v: \"\"L1\"\""]
N2 -> { N3 N4 }
N3 [label="k: aaabaaaa, b: 24, v: \"\"L5\"\""]
N4 [label="k: abaaaaaa, b: 9, v: \"\"L2.2\"\""]
}
`
)

var inv32 = []uint32{
	0x0, 0x80, 0x40, 0xc0, 0x20, 0xa0, 0x60, 0xe0, 0x10, 0x90, 0x50, 0xd0, 0x30, 0xb0, 0x70, 0xf0,
	0x8, 0x88, 0x48, 0xc8, 0x28, 0xa8, 0x68, 0xe8, 0x18, 0x98, 0x58, 0xd8, 0x38, 0xb8, 0x78, 0xf8,
	0x4, 0x84, 0x44, 0xc4, 0x24, 0xa4, 0x64, 0xe4, 0x14, 0x94, 0x54, 0xd4, 0x34, 0xb4, 0x74, 0xf4,
	0xc, 0x8c, 0x4c, 0xcc, 0x2c, 0xac, 0x6c, 0xec, 0x1c, 0x9c, 0x5c, 0xdc, 0x3c, 0xbc, 0x7c, 0xfc,
	0x2, 0x82, 0x42, 0xc2, 0x22, 0xa2, 0x62, 0xe2, 0x12, 0x92, 0x52, 0xd2, 0x32, 0xb2, 0x72, 0xf2,
	0xa, 0x8a, 0x4a, 0xca, 0x2a, 0xaa, 0x6a, 0xea, 0x1a, 0x9a, 0x5a, 0xda, 0x3a, 0xba, 0x7a, 0xfa,
	0x6, 0x86, 0x46, 0xc6, 0x26, 0xa6, 0x66, 0xe6, 0x16, 0x96, 0x56, 0xd6, 0x36, 0xb6, 0x76, 0xf6,
	0xe, 0x8e, 0x4e, 0xce, 0x2e, 0xae, 0x6e, 0xee, 0x1e, 0x9e, 0x5e, 0xde, 0x3e, 0xbe, 0x7e, 0xfe,
	0x1, 0x81, 0x41, 0xc1, 0x21, 0xa1, 0x61, 0xe1, 0x11, 0x91, 0x51, 0xd1, 0x31, 0xb1, 0x71, 0xf1,
	0x9, 0x89, 0x49, 0xc9, 0x29, 0xa9, 0x69, 0xe9, 0x19, 0x99, 0x59, 0xd9, 0x39, 0xb9, 0x79, 0xf9,
	0x5, 0x85, 0x45, 0xc5, 0x25, 0xa5, 0x65, 0xe5, 0x15, 0x95, 0x55, 0xd5, 0x35, 0xb5, 0x75, 0xf5,
	0xd, 0x8d, 0x4d, 0xcd, 0x2d, 0xad, 0x6d, 0xed, 0x1d, 0x9d, 0x5d, 0xdd, 0x3d, 0xbd, 0x7d, 0xfd,
	0x3, 0x83, 0x43, 0xc3, 0x23, 0xa3, 0x63, 0xe3, 0x13, 0x93, 0x53, 0xd3, 0x33, 0xb3, 0x73, 0xf3,
	0xb, 0x8b, 0x4b, 0xcb, 0x2b, 0xab, 0x6b, 0xeb, 0x1b, 0x9b, 0x5b, 0xdb, 0x3b, 0xbb, 0x7b, 0xfb,
	0x7, 0x87, 0x47, 0xc7, 0x27, 0xa7, 0x67, 0xe7, 0x17, 0x97, 0x57, 0xd7, 0x37, 0xb7, 0x77, 0xf7,
	0xf, 0x8f, 0x4f, 0xcf, 0x2f, 0xaf, 0x6f, 0xef, 0x1f, 0x9f, 0x5f, 0xdf, 0x3f, 0xbf, 0x7f, 0xff}

package iptree

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/infobloxopen/go-trees/numtree"
)

func TestInsertNet(t *testing.T) {
	r := NewTree()

	newR := r.InsertNet(nil, "test")
	if newR != r {
		t.Errorf("Expected no changes inserting nil network but got:\n%s\n", newR.root32.Dot())
	}

	newR = r.InsertNet(&net.IPNet{IP: nil, Mask: nil}, "test")
	if newR != r {
		t.Errorf("Expected no changes inserting invalid network but got:\n%s\n", newR.root32.Dot())
	}

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	newR = r.InsertNet(n, "test")
	if newR == r {
		t.Errorf("Expected new root after insertion of new IPv4 address but got previous")
	} else {
		assertTree32Node(newR, 0xc0000200, 24, "test", "tree with single IPv4 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8::/32")
	r1 := r.InsertNet(n, "test")
	if r1 == r {
		t.Errorf("Expected new root after insertion of new IPv6 address but got previous")
	} else {
		assertTree64Node(r1, 0x20010db800000000, 32, 0x0, 0, "test",
			"tree with single IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r2 := r1.InsertNet(n, "test 1")
	if r2 == r1 {
		t.Errorf("Expected new root after insertion of second IPv6 address but got previous")
	} else {
		assertTree64Node(r2, 0x20010db800000000, 64, 0x000000ff00000000, 32, "test 1",
			"tree with second IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	r3 := r2.InsertNet(n, "test 2")
	if r3 == r1 {
		t.Errorf("Expected new root after insertion of third IPv6 address but got previous")
	} else {
		assertTree64Node(r3, 0x20010db800000000, 64, 0x000000fe00000000, 32, "test 2",
			"tree with third IPv6 address inserted", t)
	}

	invR := NewTree()
	invR.root64 = invR.root64.Insert(0x20010db800000000, 64, "test")
	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	assertPanic(func() { invR.InsertNet(n, "panic") }, "inserting to invalid IPv6 tree", t)
}

func TestInplaceInsertNet(t *testing.T) {
	r := NewTree()

	r.InplaceInsertNet(nil, "test")
	if r.root32 != nil || r.root64 != nil {
		t.Error("Expected empty tree after inserting nil network")
	}

	r.InplaceInsertNet(&net.IPNet{IP: nil, Mask: nil}, "test")
	if r.root32 != nil || r.root64 != nil {
		t.Error("Expected empty tree after inserting invalid network")
	}

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	r.InplaceInsertNet(n, "test")
	if r.root32 == nil {
		t.Error("Expected some data in 32-bit tree")
	} else {
		assertTree32Node(r, 0xc0000200, 24, "test", "tree with single IPv4 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8::/32")
	r.InplaceInsertNet(n, "test")
	if r.root64 == nil {
		t.Error("Expected some data in 64-bit tree")
	} else {
		assertTree64Node(r, 0x20010db800000000, 32, 0x0, 0, "test",
			"tree with single IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r.InplaceInsertNet(n, "test 1")
	if r.root64 == nil {
		t.Error("Expected some data in 64-bit tree")
	} else {
		assertTree64Node(r, 0x20010db800000000, 64, 0x000000ff00000000, 32, "test 1",
			"tree with second IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	r.InplaceInsertNet(n, "test 2")
	if r.root64 == nil {
		t.Error("Expected some data in 64-bit tree")
	} else {
		assertTree64Node(r, 0x20010db800000000, 64, 0x000000fe00000000, 32, "test 2",
			"tree with third IPv6 address inserted", t)
	}

	invR := NewTree()
	invR.root64 = invR.root64.Insert(0x20010db800000000, 64, "test")
	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	assertPanic(func() { invR.InplaceInsertNet(n, "panic") }, "inserting to invalid IPv6 tree", t)
}

func (p Pair) String() string {
	s, ok := p.Value.(string)
	if ok {
		return fmt.Sprintf("%s: %q", p.Key, s)
	}

	return fmt.Sprintf("%s: %T (%#v)", p.Key, p.Value, p.Value)
}

func TestEnumerate(t *testing.T) {
	var r *Tree

	for p := range r.Enumerate() {
		t.Errorf("Expected no nodes in empty tree but got at least one: %s", p)
		break
	}

	r = NewTree()

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	r = r.InsertNet(n, "test 1")

	_, n, _ = net.ParseCIDR("2001:db8::/32")
	r = r.InsertNet(n, "test 2.1")

	_, n, _ = net.ParseCIDR("2001:db8:1::/48")
	r = r.InsertNet(n, "test 2.2")

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r = r.InsertNet(n, "test 3")

	items := []string{}
	for p := range r.Enumerate() {
		items = append(items, p.String())
	}

	s := strings.Join(items, ", ")
	e := "192.0.2.0/24: \"test 1\", " +
		"2001:db8::/32: \"test 2.1\", " +
		"2001:db8::ff:0:0/96: \"test 3\", " +
		"2001:db8:1::/48: \"test 2.2\""
	if s != e {
		t.Errorf("Expected following nodes %q but got %q", e, s)
	}
}

func TestGetByNet(t *testing.T) {
	r := NewTree()

	_, n4, _ := net.ParseCIDR("192.0.2.0/24")
	r = r.InsertNet(n4, "test 1")

	_, n6Short1, _ := net.ParseCIDR("2001:db8::/32")
	r = r.InsertNet(n6Short1, "test 2.1")

	_, n6Short2, _ := net.ParseCIDR("2001:db8:1::/48")
	r = r.InsertNet(n6Short2, "test 2.2")

	_, n6Long, _ := net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r = r.InsertNet(n6Long, "test 3")

	v, ok := r.GetByNet(nil)
	if ok {
		t.Errorf("Expected no result for nil network but got %T (%#v)", v, v)
	}

	v, ok = r.GetByNet(&net.IPNet{IP: nil, Mask: nil})
	if ok {
		t.Errorf("Expected no result for invalid network but got %T (%#v)", v, v)
	}

	v, ok = r.GetByNet(n4)
	assertResult(v, ok, "test 1", fmt.Sprintf("%s", n4), t)

	v, ok = r.GetByNet(n6Short1)
	assertResult(v, ok, "test 2.1", fmt.Sprintf("%s", n6Short1), t)

	v, ok = r.GetByNet(n6Long)
	assertResult(v, ok, "test 3", fmt.Sprintf("%s", n6Long), t)

	_, n6, _ := net.ParseCIDR("2001:db8:1::/64")
	v, ok = r.GetByNet(n6)
	assertResult(v, ok, "test 2.2", fmt.Sprintf("%s", n6), t)

	_, n6, _ = net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	v, ok = r.GetByNet(n6)
	assertResult(v, ok, "test 2.1", fmt.Sprintf("%s", n6), t)
}

func TestDeleteByNet(t *testing.T) {
	var r *Tree

	_, n4, _ := net.ParseCIDR("192.0.2.0/24")
	r, ok := r.DeleteByNet(n4)
	if ok {
		t.Errorf("Expected no deletion in empty tree but got one")
	}

	r = r.InsertNet(n4, "test 1")

	_, n6Short1, _ := net.ParseCIDR("2001:db8::/32")
	r = r.InsertNet(n6Short1, "test 2.1")

	_, n6Short2, _ := net.ParseCIDR("2001:db8:1::/48")
	r = r.InsertNet(n6Short2, "test 2.2")

	_, n6Long1, _ := net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r = r.InsertNet(n6Long1, "test 3.1")

	_, n6Long2, _ := net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	r = r.InsertNet(n6Long2, "test 3.2")

	r, ok = r.DeleteByNet(nil)
	if ok {
		t.Errorf("Expected no deletion by nil network but got one")
	}

	r, ok = r.DeleteByNet(&net.IPNet{IP: nil, Mask: nil})
	if ok {
		t.Errorf("Expected no deletion by invalid network but got one")
	}

	r, ok = r.DeleteByNet(n6Long2)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Long2)
	}

	r, ok = r.DeleteByNet(n6Long1)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Long1)
	}

	v, ok := r.root64.ExactMatch(0x20010db800000000, 64)
	if ok {
		t.Errorf("Expected no subtree node at 0x%016x, %d after deleting all long mask addresses but got %#v",
			0x20010db800000000, 64, v)
	}

	r, ok = r.DeleteByNet(n6Short2)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Short2)
	}

	r, ok = r.DeleteByNet(n6Short1)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Short1)
	}

	r, ok = r.DeleteByNet(n4)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n4)
	}

	if r.root32 != nil || r.root64 != nil {
		t.Errorf("Expected expected empty tree at the end but have root32: %#v and root64: %#v", r.root32, r.root64)
	}

	r.root64 = r.root64.Insert(0x20010db800000000, 64, "panic")
	assertPanic(func() { r.DeleteByNet(n6Long1) }, "deletion from invalid tree", t)
}

func TestTreeByIP(t *testing.T) {
	ip := net.ParseIP("2001:db8::1")

	var r *Tree
	r = r.InsertIP(ip, "test")
	if r == nil {
		t.Errorf("Expected some tree after insert %s but got %#v", ip, r)
	}

	v, ok := r.GetByIP(ip)
	assertResult(v, ok, "test", fmt.Sprintf("address %s", ip), t)

	r, ok = r.DeleteByIP(ip)
	if !ok {
		t.Errorf("Expected deletion by address %s but got nothing", ip)
	}

	r.InplaceInsertIP(ip, "test")
	if r.root64 == nil {
		t.Errorf("Expected some tree after inplace insert %s", ip)
	}
}

func TestIPv4NetToUint32(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	key, bits := iPv4NetToUint32(n)
	if key != 0xc0000200 || bits != 24 {
		t.Errorf("Expected 0xc0000200, 24 pair but got 0x%08x, %d", key, bits)
	}

	n = &net.IPNet{
		IP:   net.IP{0xc, 0x00},
		Mask: net.IPMask{0xff, 0xff, 0xff, 0x00}}
	key, bits = iPv4NetToUint32(n)
	if bits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv4 address but got 0x%08x, %d", key, bits)
	}

	n = &net.IPNet{
		IP:   net.IP{0xc, 0x00, 0x02, 0x00},
		Mask: net.IPMask{0xff, 0x00, 0xff, 0x00}}
	key, bits = iPv4NetToUint32(n)
	if bits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv4 mask but got 0x%08x, %d", key, bits)
	}
}

func TestIPv6NetToUint64Pair(t *testing.T) {
	_, n, _ := net.ParseCIDR("2001:db8::/32")
	MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n)
	if MSKey != 0x20010db800000000 || MSBits != 32 || LSKey != 0x0 || LSBits != 0 {
		t.Errorf("Expected 0x20010db800000000, 32 and 0x0000000000000000, 0 pairs bit got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	MSKey, MSBits, LSKey, LSBits = iPv6NetToUint64Pair(n)
	if MSKey != 0x20010db800000000 || MSBits != 64 || LSKey != 0x000000ff00000000 || LSBits != 32 {
		t.Errorf("Expected 0x20010db800000000, 32 and 0x0000000000000000, 0 pairs bit got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}

	n = &net.IPNet{
		IP: net.IP{
			0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00},
		Mask: net.IPMask{
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	MSKey, MSBits, LSKey, LSBits = iPv6NetToUint64Pair(n)
	if MSBits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv6 address but got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}

	n = &net.IPNet{
		IP: net.IP{
			0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Mask: net.IPMask{
			0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	MSKey, MSBits, LSKey, LSBits = iPv6NetToUint64Pair(n)
	if MSBits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv6 mask but got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}
}

func TestNewIPNetFromIP(t *testing.T) {
	n := newIPNetFromIP(net.ParseIP("192.0.2.1"))
	if n.String() != "192.0.2.1/32" {
		t.Errorf("Expected %s for IPv4 conversion but got %s", "192.0.2.1/32", n)
	}

	n = newIPNetFromIP(net.ParseIP("2001:db8::1"))
	if n.String() != "2001:db8::1/128" {
		t.Errorf("Expected %s for IPv6 conversion but got %s", "2001:db8::1/128", n)
	}

	n = newIPNetFromIP(net.IP{0xc, 0x00})
	if n != nil {
		t.Errorf("Expected %#v for invalid IP address but got %s", nil, n)
	}
}

func assertTree32Node(r *Tree, key uint32, bits int, e, desc string, t *testing.T) {
	v, ok := r.root32.ExactMatch(key, bits)
	assertResult(v, ok, e, fmt.Sprintf("0x%08x, %d for %s", key, bits, desc), t)
}

func assertTree64Node(r *Tree, MSKey uint64, MSBits int, LSKey uint64, LSBits int, e, desc string, t *testing.T) {
	desc = fmt.Sprintf("0x%016x, %d and 0x%016x, %d for %s", MSKey, MSBits, LSKey, LSBits, desc)
	v, ok := r.root64.ExactMatch(MSKey, MSBits)
	if ok {
		if MSBits < 64 {
			assertResult(v, ok, e, desc, t)
		} else {
			r, ok := v.(subTree64)
			if ok {
				v, ok := (*numtree.Node64)(r).ExactMatch(LSKey, LSBits)
				if ok {
					assertResult(v, ok, e, desc, t)
				} else {
					t.Errorf("Expected string %q at %s but got nothing at second hop", e, desc)
				}
			} else {
				t.Errorf("Expected subTree64 at %s (first hop) but got %T (%#v)", desc, v, v)
			}
		}
	} else {
		if MSBits < 64 {
			t.Errorf("Expected string %q at %s but got nothing", e, desc)
		} else {
			t.Errorf("Expected string %q at %s but got nothing even at first hop", e, desc)
		}
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

func assertResult(v interface{}, ok bool, e, desc string, t *testing.T) {
	if ok {
		s, ok := v.(string)
		if ok {
			if s != e {
				t.Errorf("Expected string %q at %s but got %q", e, desc, s)
			}
			return
		}

		t.Errorf("Expected string %q at %s but got %T (%#v)", e, desc, v, v)
		return
	}

	t.Errorf("Expected string %q at %s but got nothing", e, desc)
}

package domaintree

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/infobloxopen/go-trees/domain"
)

func TestInsert(t *testing.T) {
	var r *Node

	r1 := r.Insert(makeTestDN(t, "com"), "1")
	if r1 == nil {
		t.Error("Expected new tree but got nothing")
	}

	r2 := r1.Insert(makeTestDN(t, "test.com"), "2")
	r3 := r2.Insert(makeTestDN(t, "test.net"), "3")
	r4 := r3.Insert(makeTestDN(t, "example.com"), "4")
	r5 := r4.Insert(makeTestDN(t, "www.test.com."), "5")
	r6 := r5.Insert(makeTestDN(t, "."), "6")
	r6a := r6.Insert(makeTestDN(t, ""), "6a")

	assertTree(r, "empty tree", t)

	assertTree(r1, "single element tree", t,
		"\"com\": \"1\"\n")

	assertTree(r2, "two elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n")

	assertTree(r3, "three elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"test.net\": \"3\"\n")

	assertTree(r4, "four elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n")

	assertTree(r5, "five elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"www.test.com\": \"5\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n")

	assertTree(r6, "siz elements tree", t,
		"\"\": \"6\"\n",
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"www.test.com\": \"5\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n")

	assertTree(r6a, "five elements tree", t,
		"\"\": \"6a\"\n",
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"www.test.com\": \"5\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n")

	r = r.Insert(makeTestDN(t, "AbCdEfGhIjKlMnOpQrStUvWxYz.aBcDeFgHiJkLmNoPqRsTuVwXyZ"), "test")
	assertTree(r, "case-check tree", t,
		"\"abcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyz\": \"test\"\n")
}

func TestInplaceInsert(t *testing.T) {
	r := &Node{}
	assertTree(r, "empty inplace tree", t)

	r.InplaceInsert(makeTestDN(t, "com"), "1")
	assertTree(r, "single element inplace tree", t,
		"\"com\": \"1\"\n")

	r.InplaceInsert(makeTestDN(t, "test.com"), "2")
	assertTree(r, "two elements inplace tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n")

	r.InplaceInsert(makeTestDN(t, "test.net"), "3")
	assertTree(r, "three elements inplace tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"test.net\": \"3\"\n")

	r.InplaceInsert(makeTestDN(t, "example.com"), "4")
	assertTree(r, "four elements inplace tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n")

	r.InplaceInsert(makeTestDN(t, "www.test.com"), "5")
	assertTree(r, "five elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"www.test.com\": \"5\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n")
}

func TestGet(t *testing.T) {
	var r *Node

	v, ok := r.Get(makeTestDN(t, "test.com"))
	assertValue(v, ok, "", false, "fetching from empty tree", t)

	r = r.Insert(makeTestDN(t, "com"), "1")
	r = r.Insert(makeTestDN(t, "test.com"), "2")
	r = r.Insert(makeTestDN(t, "test.net"), "3")
	r = r.Insert(makeTestDN(t, "example.com"), "4")
	r = r.Insert(makeTestDN(t, "www.test.com"), "5")

	v, ok = r.Get(makeTestDN(t, "test.com"))
	assertValue(v, ok, "2", true, "fetching \"test.com\" from tree", t)

	v, ok = r.Get(makeTestDN(t, "www.test.com"))
	assertValue(v, ok, "5", true, "fetching \"www.test.com\" from tree", t)

	v, ok = r.Get(makeTestDN(t, "ns.test.com"))
	assertValue(v, ok, "2", true, "fetching \"ns.test.com\" from tree", t)

	v, ok = r.Get(makeTestDN(t, "test.org"))
	assertValue(v, ok, "", false, "fetching \"test.org\" from tree", t)

	v, ok = r.Get(makeTestDN(t, "nS.tEsT.cOm"))
	assertValue(v, ok, "2", true, "fetching \"nS.tEsT.cOm\" from tree", t)
}

func TestDeleteSubdomains(t *testing.T) {
	var r *Node

	r, ok := r.DeleteSubdomains(makeTestDN(t, "test.com"))
	if ok {
		t.Error("Expected no deletion from empty tree but got deleted something")
	}

	r = r.Insert(makeTestDN(t, "com"), "1")
	r = r.Insert(makeTestDN(t, "test.com"), "2")
	r = r.Insert(makeTestDN(t, "test.net"), "3")
	r = r.Insert(makeTestDN(t, "example.com"), "4")
	r = r.Insert(makeTestDN(t, "www.test.com"), "5")
	r = r.Insert(makeTestDN(t, "www.test.org"), "6")

	r, ok = r.DeleteSubdomains(makeTestDN(t, "ns.test.com"))
	if ok {
		t.Error("Expected \"ns.test.com\" to be not deleted as it's absent in the tree")
	}

	r, ok = r.DeleteSubdomains(makeTestDN(t, "test.com"))
	if !ok {
		t.Error("Expected \"test.com\" to be deleted")
	}

	r, ok = r.DeleteSubdomains(makeTestDN(t, "www.test.com"))
	if ok {
		t.Error("Expected \"www.test.com\" to be not deleted as it should be deleted with \"test.com\"")
	}

	r, ok = r.DeleteSubdomains(makeTestDN(t, "com"))
	if !ok {
		t.Error("Expected \"com\" to be deleted")
	}

	assertTree(r, "tree with no \"com\"", t,
		"\"test.net\": \"3\"\n",
		"\"www.test.org\": \"6\"\n")

	r, ok = r.DeleteSubdomains(makeTestDN(t, "test.net"))
	if !ok {
		t.Error("Expected \"test.net\" to be deleted")
	}

	r, ok = r.DeleteSubdomains(makeTestDN(t, ""))
	if !ok {
		t.Error("Expected not empty tree to be cleaned up")
	}

	r, ok = r.DeleteSubdomains(makeTestDN(t, ""))
	if ok {
		t.Error("Expected nothing to clean up from empty tree")
	}

	r = r.Insert(makeTestDN(t, "com"), "1")
	r = r.Insert(makeTestDN(t, "test.com"), "2")
	r = r.Insert(makeTestDN(t, "test.net"), "3")
	r = r.Insert(makeTestDN(t, "example.com"), "4")
	r = r.Insert(makeTestDN(t, "www.test.com"), "5")
	r = r.Insert(makeTestDN(t, "www.test.org"), "6")

	r, ok = r.DeleteSubdomains(makeTestDN(t, "WwW.tEsT.cOm"))
	if !ok {
		t.Error("Expected \"WwW.tEsT.cOm\" to be deleted")
	}
}

func TestDelete(t *testing.T) {
	var r *Node

	r, ok := r.Delete(makeTestDN(t, "test.com"))
	if ok {
		t.Error("Expected no deletion from empty tree but got deleted something")
	}

	r = r.Insert(makeTestDN(t, "com"), "1")
	r = r.Insert(makeTestDN(t, "test.com"), "2")
	r = r.Insert(makeTestDN(t, "test.net"), "3")
	r = r.Insert(makeTestDN(t, "example.com"), "4")
	r = r.Insert(makeTestDN(t, "www.test.com"), "5")
	r = r.Insert(makeTestDN(t, "www.test.org"), "6")

	r, ok = r.Delete(makeTestDN(t, "ns.test.com"))
	if ok {
		t.Error("Expected \"ns.test.com\" to be not deleted as it's absent in the tree")
	}

	r, ok = r.Delete(makeTestDN(t, "test.com"))
	if !ok {
		t.Error("Expected \"test.com\" to be deleted")
	}

	r, ok = r.Delete(makeTestDN(t, "www.test.com"))
	if !ok {
		t.Error("Expected \"www.test.com\" to be deleted")
	}

	r, ok = r.Delete(makeTestDN(t, "com"))
	if !ok {
		t.Error("Expected \"com\" to be deleted")
	}

	assertTree(r, "tree", t,
		"\"example.com\": \"4\"\n",
		"\"test.net\": \"3\"\n",
		"\"www.test.org\": \"6\"\n")

	r, ok = r.Delete(makeTestDN(t, "test.net"))
	if !ok {
		t.Error("Expected \"test.net\" to be deleted")
	}

	r, ok = r.Delete(makeTestDN(t, ""))
	if ok {
		t.Error("Expected nothing to delete from root node as it has no any value")
	}

	r = r.Insert(makeTestDN(t, ""), "root")
	assertTree(r, "tree", t,
		"\"\": \"root\"\n",
		"\"example.com\": \"4\"\n",
		"\"www.test.org\": \"6\"\n")

	r, ok = r.Delete(makeTestDN(t, ""))
	if !ok {
		t.Error("Expected root node to be deleted")
	}
	assertTree(r, "tree", t,
		"\"example.com\": \"4\"\n",
		"\"www.test.org\": \"6\"\n")

	r = r.Insert(makeTestDN(t, ""), "root")
	r, ok = r.Delete(makeTestDN(t, "example.com"))
	if !ok {
		t.Error("Expected \"example.com\" to be deleted")
	}
	r, ok = r.Delete(makeTestDN(t, "www.test.org"))
	if !ok {
		t.Error("Expected \"www.test.org\" to be deleted")
	}
	assertTree(r, "tree", t,
		"\"\": \"root\"\n")

	r, ok = r.Delete(makeTestDN(t, ""))
	if !ok {
		t.Error("Expected root node to be deleted")
	}
	assertTree(r, "tree", t)

	r = r.Insert(makeTestDN(t, "com"), "1")
	r = r.Insert(makeTestDN(t, "test.com"), "2")
	r = r.Insert(makeTestDN(t, "test.net"), "3")
	r = r.Insert(makeTestDN(t, "example.com"), "4")
	r = r.Insert(makeTestDN(t, "www.test.com"), "5")
	r = r.Insert(makeTestDN(t, "www.test.org"), "6")

	r, ok = r.Delete(makeTestDN(t, "WwW.tEsT.cOm"))
	if !ok {
		t.Error("Expected \"WwW.tEsT.cOm\" to be deleted")
	}

	r = nil
	r = r.Insert(makeTestDN(t, "escaped.\\\\label.com"), "root")
	assertTree(r, "tree", t,
		"\"escaped.\\\\\\\\label.com\": \"root\"\n")

	r, ok = r.Delete(makeTestDN(t, "escaped.\\\\label.com"))
	if !ok {
		t.Error("Expected \"escaped.\\\\\\\\label.com\" to be deleted")
	}
}

func makeTestDN(t *testing.T, s string) domain.Name {
	d, err := domain.MakeNameFromString(s)
	if err != nil {
		t.Fatalf("can't create domain name from string %q: %s", s, err)
	}

	return d
}

func assertTree(r *Node, desc string, t *testing.T, e ...string) {
	pairs := []string{}
	for p := range r.Enumerate() {
		s, ok := p.Value.(string)
		if ok {
			pairs = append(pairs, fmt.Sprintf("%q: %q\n", p.Key, s))
		} else {
			pairs = append(pairs, fmt.Sprintf("%q: %T (%#v)\n", p.Key, p.Value, p.Value))
		}
	}

	ctx := difflib.ContextDiff{
		A:        e,
		B:        pairs,
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("can't compare \"%s\": %s", desc, err))
	}

	if len(diff) > 0 {
		t.Errorf("\"%s\" doesn't match:\n%s", desc, diff)
	}
}

func assertValue(v interface{}, vok bool, e string, eok bool, desc string, t *testing.T) {
	if eok {
		if vok {
			s, ok := v.(string)
			if !ok {
				t.Errorf("Expected string %q for %s but got %T (%#v)", e, desc, v, v)
			} else if s != e {
				t.Errorf("Expected %q for %s but got %q", e, desc, s)
			}
		} else {
			t.Errorf("Expected %q for %s but got nothing", e, desc)
		}
	} else {
		if vok {
			s, ok := v.(string)
			if ok {
				t.Errorf("Expected no value for %s but got %q", desc, s)
			} else {
				t.Errorf("Expected no value for %s but got %T (%#v)", desc, v, v)
			}
		}
	}
}

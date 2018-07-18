package pep

import (
	"sync"
	"testing"

	"google.golang.org/grpc/naming"
)

func TestStaticWatcher(t *testing.T) {
	addrs := []string{
		"192.0.2.1",
		"192.0.2.2",
		"192.0.2.3",
	}
	w := newStaticWatcher(addrs)
	if w == nil {
		t.Fatalf("expected pointer to staticWatcher but got %v", w)
	}

	upds := [][]*naming.Update{}
	var (
		wg  sync.WaitGroup
		u   []*naming.Update
		err error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			u, err = w.Next()
			if err != nil {
				return
			}

			if u == nil {
				return
			}

			upds = append(upds, u)
		}
	}()

	w.Close()
	wg.Wait()

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if len(upds) != 1 {
		t.Fatalf("expected one update got %d: %#v", len(upds), upds)
	}

	u = upds[0]
	if len(u) != len(addrs) {
		t.Fatalf("expected %d addresses added but got %d: %#v", len(addrs), len(u), u)
	}

	for i, addr := range addrs {
		if u[i].Op != naming.Add || u[i].Addr != addr {
			t.Errorf("expected address %d %q added (%d) but got operation %d for address %q",
				i, addr, naming.Add, u[i].Op, u[i].Addr)
		}
	}
}

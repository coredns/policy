package pep

import "testing"

func TestStaticResolver(t *testing.T) {
	r := newStaticResolver(virtualServerAddress,
		"192.0.2.1",
		"192.0.2.2",
		"192.0.2.3",
	)
	if r == nil {
		t.Fatalf("expected pointer to static resolver but got %v", r)
	}

	w, err := r.Resolve(virtualServerAddress)
	if err != nil {
		t.Errorf("expected no error but got %s", err)
	}

	if w == nil {
		t.Errorf("expected pointer to static watcher but got %v", w)
	}

	_, err = r.Resolve("wrong.name")
	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

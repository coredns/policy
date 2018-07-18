package perf

import "time"

type timing struct {
	s time.Time
	r time.Time
	e error
}

func (t *timing) setSend() {
	t.s = time.Now()
}

func (t *timing) setReceive() {
	t.r = time.Now()
}

func (t *timing) setError(err error) {
	t.e = err
}

type bySend []timing

func (s bySend) Len() int           { return len(s) }
func (s bySend) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySend) Less(i, j int) bool { return s[i].s.Before(s[j].s) }

type byRecive []timing

func (s byRecive) Len() int      { return len(s) }
func (s byRecive) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byRecive) Less(i, j int) bool {
	if s[i].e != nil && s[j].e != nil {
		return s[i].s.Before(s[j].s)
	}

	if s[i].e != nil {
		return false
	}

	if s[j].e != nil {
		return true
	}

	return s[i].r.Before(s[j].r)
}

type timings struct {
	Sends    []int64   `json:"sends"`
	Receives []int64   `json:"receives"`
	Pairs    [][]int64 `json:"pairs"`
}

package perf

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func dump(recs []timing, path string) error {
	tm := timings{
		Sends:    make([]int64, len(recs)),
		Receives: make([]int64, len(recs)),
		Pairs:    make([][]int64, len(recs)),
	}

	sort.Sort(bySend(recs))
	for i, t := range recs {
		tm.Sends[i] = t.s.UnixNano()
		if t.e != nil {
			tm.Pairs[i] = []int64{t.s.UnixNano()}
		} else {
			tm.Pairs[i] = []int64{
				t.s.UnixNano(),
				t.r.UnixNano(),
				t.r.UnixNano() - t.s.UnixNano(),
			}
		}
	}

	sort.Sort(byRecive(recs))
	for i, t := range recs {
		if t.e == nil {
			tm.Receives[i] = t.r.UnixNano()
		}
	}

	b, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshal timings to JSON: %s", err)
	}

	f := os.Stdout
	dstName := "stdout"
	if len(path) > 0 {
		f, err = os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		dstName = fmt.Sprintf("file %s", path)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("can't dump JSON timings to %s: %s", dstName, err)
	}

	return nil
}

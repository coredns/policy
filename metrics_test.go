package policy

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/infobloxopen/themis/pdp"
	dto "github.com/prometheus/client_model/go"
)

// ============= SlicedCounter tests ==============

func TestSlicedCounterInc(t *testing.T) {
	var utime uint32 = 100
	sc := NewSlicedCounter(100)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				ut := atomic.LoadUint32(&utime)
				if !sc.Inc(ut) {
					if ut < 100+BucketCnt {
						t.Errorf("Inc unexpectedly returned false")
					}
					break
				}
				if ut >= 100+BucketCnt {
					t.Errorf("Inc unexpectedly returned true")
					break
				}
				time.Sleep(time.Nanosecond)
			}
			wg.Done()
		}()
	}
	for i := 0; i <= BucketCnt; i++ {
		atomic.AddUint32(&utime, 1)
		time.Sleep(time.Microsecond)
	}
	wg.Wait()

	checkTotal(t, sc)
}

func TestNoStale(t *testing.T) {
	sc := newTestSlicedCounter(100)
	sc.EraseStale(105)

	for i := 0; i < BucketCnt; i++ {
		if sc.buckets[i] != uint32(i+10) {
			t.Errorf("bucket[%d] unexpectedly was erased", i)
		}
	}

	checkTotal(t, sc)
}

func Test6Stale(t *testing.T) {
	sc := newTestSlicedCounter(100)
	sc.EraseStale(115)

	for i := 0; i < 4; i++ {
		if sc.buckets[i] != uint32(i+10) {
			t.Errorf("bucket[%d] unexpectedly was erased", i)
		}
	}
	for i := 4; i < 10; i++ {
		if sc.buckets[i] != 0 {
			t.Errorf("bucket[%d] unexpectedly was not erased", i)
		}
	}
	for i := 10; i < 16; i++ {
		if sc.buckets[i] != uint32(i+10) {
			t.Errorf("bucket[%d] unexpectedly was erased", i)
		}
	}

	checkTotal(t, sc)
}

func TestAllStale(t *testing.T) {
	sc := newTestSlicedCounter(100)
	sc.EraseStale(130)

	for i := 0; i < BucketCnt; i++ {
		if sc.buckets[i] != 0 {
			t.Errorf("bucket[%d] unexpectedly was not erased", i)
		}
	}

	checkTotal(t, sc)
}

func TestIncVsAllStale(t *testing.T) {
	sc := newTestSlicedCounter(100)
	var testTime uint32 = 140
	var utime uint32 = testTime

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				ut := atomic.LoadUint32(&utime)
				sc.Inc(ut)
				if ut > testTime {
					break
				}
				time.Sleep(time.Nanosecond)
			}
			wg.Done()
		}()
	}
	sc.EraseStale(testTime)

	atomic.AddUint32(&utime, 1)
	wg.Wait()

	if sc.buckets[utime%BucketCnt] != 10 {
		t.Errorf("Unexpected counter after Erase, expected 10, got %d", sc.buckets[utime%BucketCnt])
	}
	checkTotal(t, sc)
}

func TestIncVsPartiallyStale(t *testing.T) {
	sc := newTestSlicedCounter(100)
	var testTime uint32 = 114
	var utime uint32 = testTime

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				ut := atomic.LoadUint32(&utime)
				if !sc.Inc(ut) {
					t.Errorf("Inc unexpectedly returned false")
				}
				if ut > testTime {
					break
				}
				time.Sleep(time.Nanosecond)
			}
			wg.Done()
		}()
	}
	sc.EraseStale(testTime)

	atomic.AddUint32(&utime, 1)
	wg.Wait()

	checkTotal(t, sc)
}

// ============= AttrGauge tests ==============

func TestStartStop(t *testing.T) {
	ag := newTestAttrGauge()

	ag.Start(10*time.Millisecond, 20)
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadUint32(&ag.state) != AttrGaugeStarted {
		t.Errorf("AttrGauge has not started")
	}

	ag.Stop()
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadUint32(&ag.state) != AttrGaugeStopped {
		t.Errorf("AttrGauge has not stopped")
	}
}

func TestNegativeStop(t *testing.T) {
	ag := newTestAttrGauge()

	ag.state = AttrGaugeStopped
	ag.Stop()
	if atomic.LoadUint32(&ag.state) == AttrGaugeStopping {
		t.Errorf("AttrGauge is unexpectedly stopping")
	}

	ag.state = AttrGaugeStopping
	ag.Stop()
	if atomic.LoadUint32(&ag.state) != AttrGaugeStopping {
		t.Errorf("AttrGauge is unexpectedly stopped")
	}
}

func TestAttrGaugeNil(t *testing.T) {
	var ag *AttrGauge

	ag.Inc(testAttr())
	ag.Stop()
}

func TestAttrGaugeInc(t *testing.T) {
	ag := newTestAttrGauge()
	attr := testAttr()

	setTestTime(100)
	ag.synchInc(attr)
	ag.tick()
	scVal := totalVal(t, ag, attr)
	gVal, err := gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, scVal, 1); err != nil {
		t.Error(err)
	}

	setTestTime(101)
	ag.synchInc(attr)
	ag.tick()
	scVal = totalVal(t, ag, attr)
	gVal, err = gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, scVal, 2); err != nil {
		t.Error(err)
	}
}

func TestAttrGaugeTick(t *testing.T) {
	ag := newTestAttrGauge()
	attr := testAttr()

	for i := 100; i < 120; i++ {
		setTestTime(uint32(i))
		ag.synchInc(attr)
	}
	scVal := totalVal(t, ag, attr)
	if scVal != 16 {
		t.Errorf("unexpected counter, expected %d, got %d", 16, scVal)
	}

	setTestTime(120)
	eCnt := ag.tick()
	if eCnt != 4 {
		t.Errorf("unexpected error count, expected %d, got %d", 4, eCnt)
	}
	scVal = totalVal(t, ag, attr)
	gVal, err := gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, scVal, 5); err != nil {
		t.Error(err)
	}
}

func TestAttrGaugeEraseValue(t *testing.T) {
	ag := newTestAttrGauge()
	attr := testAttr()

	setTestTime(100)
	ag.Inc(attr)
	setTestTime(120)
	ag.tick()
	scVal := totalVal(t, ag, attr)
	gVal, err := gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, scVal, 0); err != nil {
		t.Error(err)
	}
}

func TestAttrGaugeSubsequentTicks(t *testing.T) {
	ag := newTestAttrGauge()
	ag.Start(10*time.Millisecond, 20)
	attr := testAttr()

	setTestTime(93)
	ag.Inc(attr)
	time.Sleep(100 * time.Millisecond)
	setTestTime(94)
	ag.Inc(attr)
	time.Sleep(100 * time.Millisecond)
	gVal, err := gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, gVal, 2); err != nil {
		t.Error(err)
	}

	setTestTime(103)
	time.Sleep(100 * time.Millisecond)
	gVal, err = gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, gVal, 1); err != nil {
		t.Error(err)
	}

	setTestTime(104)
	time.Sleep(100 * time.Millisecond)
	gVal, err = gaugeVal(t, ag, attr)
	if err := checkVal(err, gVal, gVal, 0); err != nil {
		t.Error(err)
	}

	ag.Stop()
}

func TestAttrGaugeErrorInc(t *testing.T) {
	ag := newTestAttrGauge()
	attr := testAttr()

	setTestTime(93)
	ag.Inc(attr)
	ag.Inc(attr)
	ag.Inc(attr)

	if ag.errCnt != 3 {
		t.Errorf("unexpected error count, expected %d, got %d", 3, ag.errCnt)
	}
}

func TestAttrGaugeAddAttributes(t *testing.T) {
	ag := newTestAttrGauge()
	setTestTime(100)

	ag.Start(10*time.Millisecond, DefaultQueryChanSize)
	// Just make sure the test doesn't panic

	ag.AddAttributes("test_attr1")
	ag.Inc(pdp.MakeStringAssignment("test_attr1", "test_value1"))

	ag.AddAttributes("test_attr2")
	ag.Inc(pdp.MakeStringAssignment("test_attr2", "test_value2"))

	ag.Stop()
}

func TestAttrGaugeAddAttributeAgain(t *testing.T) {
	ag := newTestAttrGauge()
	setTestTime(100)

	ag.Start(10*time.Millisecond, DefaultQueryChanSize)
	attr := testAttr()

	ag.Inc(attr)
	ag.AddAttributes("test_attr")
	ag.Inc(attr)

	time.Sleep(100 * time.Millisecond)
	ag.Stop()
	gVal, err := gaugeVal(t, ag, attr)
	if err = checkVal(err, gVal, gVal, 2); err != nil {
		t.Error(err)
	}
}

// ============== utility functions ===============

func newTestSlicedCounter(ut uint32) *SlicedCounter {
	sc := NewSlicedCounter(ut)
	for i := 0; i < BucketCnt; i++ {
		sc.buckets[i] = uint32(i + 10)
		sc.total += sc.buckets[i]
	}
	return sc
}

func logSc(t *testing.T, sc *SlicedCounter) {
	for i := 0; i < BucketCnt; i++ {
		t.Logf("bucket[%d] == %d", i, sc.buckets[i])
	}
}

func checkTotal(t *testing.T, sc *SlicedCounter) {
	var total uint32
	for i := 0; i < BucketCnt; i++ {
		total += sc.buckets[i]
	}
	if total != sc.Total() {
		t.Errorf("Unexpected total, expected=%d, actual=%d", total, sc.Total())
	}
}

func testAttr() pdp.AttributeAssignment {
	return pdp.MakeStringAssignment("test_attr", "test_value")
}

var utime uint32

func testTime() uint32 {
	return atomic.LoadUint32(&utime)
}

func setTestTime(t uint32) {
	atomic.StoreUint32(&utime, t)
}

func newTestAttrGauge() *AttrGauge {
	ag := NewAttrGauge()
	ag.addAttribute("test_attr")
	ag.timeFunc = testTime
	return ag
}

func totalVal(t *testing.T, ag *AttrGauge, attr pdp.AttributeAssignment) uint32 {
	if vMap, ok := ag.perAttr[attr.GetID()]; ok {
		if sc, ok := vMap[serializeOrPanic(attr)]; ok {
			return sc.Total()
		}
	}
	return 0
}

func gaugeVal(t *testing.T, ag *AttrGauge, attr pdp.AttributeAssignment) (uint32, error) {
	g, e := ag.pgv.GetMetricWithLabelValues(attr.GetID(), serializeOrPanic(attr))
	if e != nil {
		return 0, e
	}
	metric := &dto.Metric{}
	g.Write(metric)
	out, e := json.Marshal(metric)
	if e != nil {
		return 0, e
	}

	result := make(map[string]interface{})
	e = json.Unmarshal(out, &result)
	if e != nil {
		return 0, e
	}
	if v, ok := result["gauge"]; ok {
		if vMap, ok := v.(map[string]interface{}); ok {
			if v, ok := vMap["value"]; ok {
				if v, ok := v.(float64); ok {
					return uint32(v), nil
				}
			}
		}
	}
	return 0, fmt.Errorf("Gauge value not found")
}

func checkVal(err error, gVal, scVal, expVal uint32) error {
	if err != nil {
		return fmt.Errorf("Failed to get gauge value - %s", err)
	}
	if gVal != expVal {
		return fmt.Errorf("unexpected gauge value, expected %d, got %d", expVal, gVal)
	}
	if gVal != scVal {
		return fmt.Errorf("gauge value mismatch, gauge=%d, slicedCounter=%d", gVal, scVal)
	}
	return nil
}

func resetGlobals() {

}

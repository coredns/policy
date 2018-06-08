package server

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"time"

	log "github.com/sirupsen/logrus"
)

// MemLimits structure contains memory limit levels to manage GC
type MemLimits struct {
	limit uint64
	reset float64
	soft  float64
	frag  float64
	back  float64
}

// MakeMemLimits fills MemLimits structure with given parameters
func MakeMemLimits(limit uint64, reset, soft, back, frag float64) (MemLimits, error) {
	m := MemLimits{limit: limit}
	if m.limit > 0 {
		if reset < 0 || reset > 100 {
			return MemLimits{}, fmt.Errorf("reset limit should be in range 0 - 100 but got %f", reset)
		}

		if soft < 0 || soft > 100 {
			return MemLimits{}, fmt.Errorf("soft limit should be in range 0 - 100 but got %f", soft)
		}

		if soft >= reset {
			return MemLimits{},
				fmt.Errorf("reset limit should be higher than soft limit "+
					"but got %f >= %f", soft, reset)
		}

		m.reset = reset / 100 * float64(m.limit)
		m.soft = soft / 100 * float64(m.limit)

		if back < 0 || back > 100 {
			return MemLimits{}, fmt.Errorf("back percentage should be in range 0 - 100 but got %f", back)
		}
		m.back = back / 100

		if frag < 0 || frag > 100 {
			return MemLimits{},
				fmt.Errorf("fragmentation warning percentage should be in range 0 - 100 "+
					"but got %f", frag)
		}
		m.frag = frag / 100
	}

	return m, nil
}

func (s *Server) checkMemory(c *MemLimits) {
	if c.limit <= 0 {
		return
	}

	now := time.Now()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	total := float64(m.Sys - m.HeapReleased)
	if total >= c.reset {
		s.opts.logger.WithFields(log.Fields{
			"allocated": fmtMemSize(m.Sys),
			"limit":     fmtMemSize(c.limit)}).Error("Memory usage is too high. Exiting...")

		s.Stop()
		c.limit = 0
		return
	}

	if total >= c.soft {
		if s.softMemWarn == nil {
			tmp := now
			s.softMemWarn = &tmp

			s.opts.logger.WithFields(log.Fields{
				"allocated": fmtMemSize(m.Sys),
				"limit":     fmtMemSize(c.limit)}).Warn("Memory usage essentially increased")
		} else if now.Sub(*s.softMemWarn) > time.Minute {
			*s.softMemWarn = now

			s.opts.logger.WithFields(log.Fields{
				"allocated": fmtMemSize(m.Sys),
				"limit":     fmtMemSize(c.limit)}).Warn("Memory usage remains high")
		}
	} else {
		s.softMemWarn = nil
	}

	limit := float64(c.limit)
	if total > 0.1*limit && float64(m.HeapInuse-m.HeapAlloc)/total >= c.frag {
		if s.fragMemWarn == nil {
			tmp := now
			s.fragMemWarn = &tmp

			s.opts.logger.WithFields(log.Fields{
				"allocated":    fmtMemSize(m.Sys),
				"in-use":       fmtMemSize(m.HeapAlloc),
				"in-use-spans": fmtMemSize(m.HeapInuse)}).Warn("Amount of fragmented memory essentially increased")
		} else if now.Sub(*s.fragMemWarn) > time.Minute {
			*s.fragMemWarn = now

			s.opts.logger.WithFields(log.Fields{
				"allocated":    fmtMemSize(m.Sys),
				"in-use":       fmtMemSize(m.HeapAlloc),
				"in-use-spans": fmtMemSize(m.HeapInuse)}).Warn("Amount of fragmented memory remains high")
		}
	} else {
		s.fragMemWarn = nil
	}

	if total > 0.1*limit && (total-float64(m.HeapAlloc))/total >= c.back {
		if s.backMemWarn == nil {
			tmp := now
			s.backMemWarn = &tmp

			s.opts.logger.WithFields(log.Fields{
				"allocated": fmtMemSize(m.Sys),
				"in-use":    fmtMemSize(m.HeapAlloc)}).Warn("Amount of unused memory essentially increased")
		} else if now.Sub(*s.backMemWarn) > time.Minute {
			*s.backMemWarn = now

			s.opts.logger.WithFields(log.Fields{
				"allocated": fmtMemSize(m.Sys),
				"in-use":    fmtMemSize(m.HeapAlloc)}).Warn("Amount of unused memory remains high")
		}
	} else {
		s.backMemWarn = nil
	}
}

func (s *Server) memoryChecker() {
	c := s.opts.memLimits
	if c == nil || c.limit <= 0 {
		return
	}

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for c.limit > 0 {
		<-t.C
		s.checkMemory(c)
	}
}

func fmtMemSize(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%d", size)
	}

	s := float32(size) / 1024
	if s < 1024 {
		return fmt.Sprintf("%.2f KB", s)
	}

	s /= 1024
	if s < 1024 {
		return fmt.Sprintf("%.2f MB", s)
	}

	s /= 1024
	if s < 1024 {
		return fmt.Sprintf("%.2f GB", s)
	}

	s /= 1024
	return fmt.Sprintf("%.2f TB", s)
}

type memStatsLogEntry struct {
	Timestamp int64
	MemStats  runtime.MemStats
}

func (s *Server) memStatsLogging(done <-chan struct{}) {
	if s.opts.memStatsLogInterval < 0 {
		return
	}

	f, err := os.Create(s.opts.memStatsLogPath)
	if err != nil {
		s.opts.logger.WithError(err).Fatal("Failed to create file for runtime.MemStats logs")
	}
	defer f.Close()

	out := json.NewEncoder(f)

	var e memStatsLogEntry
	m := &e.MemStats

	if s.opts.memStatsLogInterval > 0 {
		ticker := time.NewTicker(s.opts.memStatsLogInterval)
		for {
			select {
			case <-done:
				ticker.Stop()
				return

			case now := <-ticker.C:
				e.Timestamp = now.UnixNano()
				runtime.ReadMemStats(m)

				out.Encode(e)
				f.Sync()
			}
		}
	} else {
		e.Timestamp = time.Now().UnixNano()
		runtime.ReadMemStats(m)

		minE := e
		maxE := e
		sameE := true

		var prevNumGC uint32 = 0

		ticker := time.NewTicker(memStatsCheckInterval)
		for {
			select {
			case <-done:
				ticker.Stop()
				return

			case now := <-ticker.C:
				e.Timestamp = now.UnixNano()
				runtime.ReadMemStats(m)

				if prevNumGC != e.MemStats.NumGC {
					if sameE {
						out.Encode(minE)
					} else if minE.Timestamp < maxE.Timestamp {
						out.Encode(minE)
						out.Encode(maxE)
					} else {
						out.Encode(maxE)
						out.Encode(minE)
					}

					f.Sync()

					minE = e
					maxE = e
					sameE = true
				} else if e.MemStats.Alloc < minE.MemStats.Alloc && e.MemStats.Alloc >= maxE.MemStats.Alloc {
					minE = e
					maxE = e
					sameE = true
				} else if e.MemStats.Alloc < minE.MemStats.Alloc {
					minE = e
					sameE = false
				} else if e.MemStats.Alloc >= maxE.MemStats.Alloc {
					maxE = e
					sameE = false
				}

				prevNumGC = e.MemStats.NumGC
			}
		}
	}
}

func (s *Server) memProfCleanup() {
	if err := os.RemoveAll(s.opts.memProfDumpPath); err != nil {
		s.opts.logger.WithError(err).Warn("Failed to cleanup directory for memory profiles")
	}

	if err := os.MkdirAll(s.opts.memProfDumpPath, 0755); err != nil {
		s.opts.logger.WithError(err).Fatal("Failed to create directory for memory profiles")
	}
}

func (s *Server) memProfDump(numGC uint32) {
	name := fmt.Sprintf("mem-%09d.pprof", numGC)
	f, err := os.Create(path.Join(s.opts.memProfDumpPath, name))
	if err != nil {
		s.opts.logger.WithFields(log.Fields{
			"numGC": numGC,
			"err":   err,
		}).Debug("Failed to create file for memory profile dump. Skipping...")
		return
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		s.opts.logger.WithFields(log.Fields{
			"name":  name,
			"numGC": numGC,
			"err":   err,
		}).Debug("Falied to dump memory profile. Skipping...")
	}
}

func (s *Server) memProfBaseDump() {
	if s.memProfBaseDumpDone == nil {
		return
	}

	defer close(s.memProfBaseDumpDone)
	s.memProfCleanup()

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	s.memProfDump(m.NumGC)
	s.memProfBaseDumpDone <- m.NumGC
}

func (s *Server) memProfDumping(done <-chan struct{}) {
	if s.opts.memProfNumGC == 0 {
		return
	}

	m := new(runtime.MemStats)
	var (
		startNumGC uint32
		lastNumGC  uint32 = math.MaxUint32
	)

	if s.memProfBaseDumpDone == nil {
		s.memProfCleanup()
	} else {
		select {
		case <-done:
			return

		case n, ok := <-s.memProfBaseDumpDone:
			if !ok {
				return
			}

			startNumGC = n
		}
	}

	ticker := time.NewTicker(memStatsCheckInterval)
	for {
		select {
		case <-done:
			ticker.Stop()
			return

		case <-ticker.C:
			runtime.ReadMemStats(m)
			if (m.NumGC-startNumGC)%s.opts.memProfNumGC == 0 && lastNumGC != m.NumGC {
				s.memProfDump(m.NumGC)

				lastNumGC = m.NumGC
			}
		}
	}
}

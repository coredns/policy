package file

import (
	"os"
	"time"
)

// TickTime is the default time we use to reload zone. Exported to be tweaked in tests.
var TickTime = 1 * time.Minute

// Reload reloads a zone when it is changed on disk. If z.NoRoload is true, no reloading will be done.
func (z *Zone) Reload() error {
	if z.NoReload {
		return nil
	}

	tick := time.NewTicker(TickTime)

	go func() {

		for {
			select {

			case <-tick.C:
				zFile := z.File()
				reader, err := os.Open(zFile)
				if err != nil {
					log.Errorf("Failed to open zone %q in %q: %v", z.origin, zFile, err)
					continue
				}

				serial := z.SOASerialIfDefined()
				zone, err := Parse(reader, z.origin, zFile, serial)
				if err != nil {
					if _, ok := err.(*serialErr); !ok {
						log.Errorf("Parsing zone %q: %v", z.origin, err)
					}
					continue
				}

				// copy elements we need
				z.reloadMu.Lock()
				z.Apex = zone.Apex
				z.Tree = zone.Tree
				z.reloadMu.Unlock()

				log.Infof("Successfully reloaded zone %q in %q with serial %d", z.origin, zFile, z.Apex.SOA.Serial)
				z.Notify()

			case <-z.reloadShutdown:
				tick.Stop()
				return
			}
		}
	}()
	return nil
}

// SOASerialIfDefined returns the SOA's serial if the zone has a SOA record in the Apex, or
// -1 otherwise.
func (z *Zone) SOASerialIfDefined() int64 {
	z.reloadMu.Lock()
	defer z.reloadMu.Unlock()
	if z.Apex.SOA != nil {
		return int64(z.Apex.SOA.Serial)
	}
	return -1
}

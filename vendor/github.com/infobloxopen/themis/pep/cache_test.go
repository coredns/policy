package pep

import (
	"testing"
	"time"

	"github.com/allegro/bigcache"
)

func TestAdjustCacheConfig(t *testing.T) {
	cfg := bigcache.DefaultConfig(15 * time.Minute)
	cfg = adjustCacheConfig(cfg)
	if cfg.Shards != 1024 || cfg.MaxEntriesInWindow != 536870 {
		t.Errorf("Expected %d shards and %d entries in window but got %d and %d",
			1024, 536870, cfg.Shards, cfg.MaxEntriesInWindow)
	}

	cfg = bigcache.DefaultConfig(15 * time.Minute)
	cfg.HardMaxCacheSize = 128
	cfg.MaxEntrySize = 10240
	cfg = adjustCacheConfig(cfg)
	if cfg.Shards != 256 || cfg.MaxEntriesInWindow != 16384 {
		t.Errorf("Expected %d shards and %d entries in window but got %d and %d",
			1024, 536870, cfg.Shards, cfg.MaxEntriesInWindow)
	}
}

func TestRound(t *testing.T) {
	if n := round(1.4); n != 1 {
		t.Errorf("Expected 1.4 rounded to 1 but got %g", n)
	}

	if n := round(1.5); n != 2 {
		t.Errorf("Expected 1.5 rounded to 2 but got %g", n)
	}

	if n := round(-1.4); n != -1 {
		t.Errorf("Expected -1.4 rounded to -1 but got %g", n)
	}

	if n := round(-1.5); n != -2 {
		t.Errorf("Expected -1.5 rounded to -2 but got %g", n)
	}
}

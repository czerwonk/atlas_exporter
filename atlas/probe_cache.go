package atlas

import (
	"time"

	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/common/log"
)

var cache *probe.Cache

// InitCache initializes the cache
func InitCache(ttl, cleanup time.Duration) {
	cache = probe.NewCache(ttl)
	startCacheCleanupFunc(cleanup)
}

func startCacheCleanupFunc(d time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(d):
				log.Infoln("Cleaning up cache...")
				r := cache.CleanUp()
				log.Infof("Items removed: %d", r)
			}
		}
	}()
}

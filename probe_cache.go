package main

import (
	"time"

	"github.com/czerwonk/atlas_exporter/probe"
	"github.com/prometheus/common/log"
)

var cache *probe.ProbeCache

func initCache() {
	cache = probe.NewCache(time.Duration(*cacheTTL) * time.Second)
	startCacheCleanupFunc(time.Duration(*cacheCleanUp) * time.Second)
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

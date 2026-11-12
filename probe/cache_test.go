// SPDX-License-Identifier: LGPL-3.0-or-later

package probe

import (
	"sync"
	"testing"
	"time"
)

func TestCacheCleanUpRemovesExpired(t *testing.T) {
	c := NewCache(10 * time.Millisecond)
	c.Add(1, &Probe{})

	if _, ok := c.Get(1); !ok {
		t.Fatal("expected probe 1 to be cached")
	}

	time.Sleep(20 * time.Millisecond)
	if n := c.CleanUp(); n != 1 {
		t.Fatalf("CleanUp removed %d, want 1", n)
	}
	if _, ok := c.Get(1); ok {
		t.Fatal("expired probe 1 was still cached")
	}
}

func TestCacheCleanUpConcurrentWithGetAndAdd(t *testing.T) {
	c := NewCache(time.Millisecond)
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p := &Probe{}
			for j := 0; j < 1000; j++ {
				c.Add(id, p)
				c.Get(id)
				c.CleanUp()
			}
		}(i)
	}

	wg.Wait()
}

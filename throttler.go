package main

import (
	"context"
	"sync"
	"time"

	"github.com/jonnrb/speedtest"
)

// Helper to make a new Throttler. Source field is left zero because it is
// assumed that multiple throttlers will share a Linearizer but will have
// different Sources.
//
func NewThrottler(lin chan func()) *Throttler {
	return &Throttler{
		Interval:   *throttleTime,
		Linearizer: lin,
	}
}

// Throttles repeated probe queries by caching the last result for a specified
// interval. When that interval expires and the Throttler is read again, the
// source RadarSign gets read from through the global Linearizer.
//
type Throttler struct {
	Interval   time.Duration
	Source     RadarSign
	Linearizer chan func()

	mu    sync.Mutex
	last  time.Time
	gauge speedtest.BytesPerSecond
}

func (t *Throttler) Read(ctx context.Context) (speed speedtest.BytesPerSecond, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if time.Now().Before(t.last.Add(t.Interval)) {
		speed = t.gauge
		return
	}

	done := make(chan struct{})

	select {
	case t.Linearizer <- func() {
		speed, err = t.Source.Read(ctx)
		close(done)
	}:
	case <-ctx.Done():
		err = ctx.Err()
		close(done)
	}

	<-done
	if err == nil {
		t.last = time.Now()
		t.gauge = speed
	}
	return
}
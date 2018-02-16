package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
)

func TestLinearizer(t *testing.T) {
	const trials = 100

	defer leaktest.Check(t)()

	c := Linearizer()
	var (
		val int64 = 0
		wg  sync.WaitGroup
	)

	for i := 0; i < trials; i++ {
		wg.Add(1)
		go func(i int64) {
			c <- func() {
				if !atomic.CompareAndSwapInt64(&val, 0, i) {
					t.Fail()
				}
				time.Sleep(time.Millisecond)
				if !atomic.CompareAndSwapInt64(&val, i, 0) {
					t.Fail()
				}
				wg.Done()
			}
		}(int64(i))
	}

	wg.Wait()
	close(c)
}

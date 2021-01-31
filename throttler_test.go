package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.jonnrb.io/speedtest/units"
)

type mockRadarSign struct {
	ReadImpl func(ctx context.Context) (units.BytesPerSecond, error)
}

func (m *mockRadarSign) Read(ctx context.Context) (units.BytesPerSecond, error) {
	return m.ReadImpl(ctx)
}

func TestThrottler_Read_Throttled(t *testing.T) {
	const (
		expected = 42 * units.GBps
		runs     = 10
	)

	l := Linearizer()
	th := NewThrottler(l)
	th.Interval = 10 * time.Hour

	i := 0
	r := func(ctx context.Context) (units.BytesPerSecond, error) {
		if i > 0 {
			t.Log("Throttler should have only allowed 1 read from source.")
			t.Fail()
		}
		return expected, nil
	}
	th.Source = &mockRadarSign{r}

	for i := 0; i < runs; i++ {
		s, err := th.Read(context.Background())
		if err != nil {
			t.Logf("Read failed: %v", err)
			t.Fail()
		}
		if s != expected {
			t.Logf("Expected %v; got %v", expected, s)
			t.Fail()
		}
	}
}

func TestThrottler_Read_IntervalExceeded(t *testing.T) {
	const (
		expected = 42 * units.GBps
		runs     = 10
		interval = 2 * time.Millisecond
	)

	l := Linearizer()
	th := NewThrottler(l)
	th.Interval = interval

	i := 0
	r := func(ctx context.Context) (units.BytesPerSecond, error) {
		i++
		if i > runs {
			t.Logf("Throttler should have only done at most %v reads from source.", runs)
			t.Fail()
		}
		return expected, nil
	}
	th.Source = &mockRadarSign{r}

	for i := 0; i < runs; i++ {
		s, err := th.Read(context.Background())
		if err != nil {
			t.Logf("Read failed: %v", err)
			t.Fail()
		}
		if s != expected {
			t.Logf("Expected %v; got %v", expected, s)
			t.Fail()
		}
		time.Sleep(interval)
	}

	if i < runs {
		t.Logf("Expected %v reads from source.", runs)
		t.Fail()
	}
}

func TestThrottler_Read_ContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	l := Linearizer()

	// Occupy the Linearizer
	done := make(chan struct{})
	l <- func() {
		<-done
	}

	th := NewThrottler(l)
	badErr := fmt.Errorf("should not get here")
	r := func(ctx context.Context) (units.BytesPerSecond, error) {
		t.Fail()
		t.Log("Read should never execute.")
		return 42 * units.GBps, badErr
	}
	th.Source = &mockRadarSign{r}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
		time.Sleep(5 * time.Second)
		close(done)
	}()

	s, err := th.Read(ctx)
	if s != units.BytesPerSecond(0) {
		t.Fail()
		t.Log("Should have gotten zero speed.")
	}
	if err == badErr {
		t.Log(err)
		t.Fail()
	}
}

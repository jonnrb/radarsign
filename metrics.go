package main

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	downloadSpeedGaugeOpts = prometheus.GaugeOpts{
		Name: "download_speed",
		Help: "Gauge representing network download speed in bytes per second.",
	}
	uploadSpeedGaugeOpts = prometheus.GaugeOpts{
		Name: "upload_speed",
		Help: "Gauge representing network upload speed in bytes per second.",
	}
)

// Binds a RadarSign (speed prober) to a Prometheus Collector that has a timeout
// on gauge reads.
//
type RadarSignMetric struct {
	RadarSign
	prometheus.Collector
	Timeout time.Duration
}

func NewRadarSignMetric(s RadarSign, opts prometheus.GaugeOpts) *RadarSignMetric {
	r := &RadarSignMetric{}

	r.RadarSign = s
	r.Collector = prometheus.NewGaugeFunc(opts,
		func() float64 {
			ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
			defer cancel()

			speed, err := s.Read(ctx)

			if err != nil {
				glog.Errorf("Error probing RadarSign metric %q: %v", opts.Name, err)
				return 0
			} else {
				return float64(speed)
			}
		})

	return r
}

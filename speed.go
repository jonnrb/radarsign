package main

import (
	"context"
	"fmt"

	"go.jonnrb.io/speedtest"
)

type RadarSign interface {
	Read(context.Context) (speedtest.BytesPerSecond, error)
}

type RepeatableSpeedTest struct {
	// SpeedTest client to use on each probe.
	//
	Client speedtest.Client

	// Selected few servers (should be by geographic proximity) to select from
	// in fastestServer().
	//
	CandidateServers []speedtest.Server
}

type DownloadRadarSign struct {
	RepeatableSpeedTest
}

type UploadRadarSign struct {
	RepeatableSpeedTest
}

func NewRepeatableSpeedTest(ctx context.Context, candidates int) (st RepeatableSpeedTest, err error) {
	if candidates <= 0 {
		err = fmt.Errorf("must try for at least 1 candidate; got %v", candidates)
		return
	}

	cfg, err := st.Client.Config(ctx)
	if err != nil {
		return
	}

	servers, err := st.Client.LoadAllServers(ctx)
	if err != nil {
		return
	}
	if len(servers) == 0 {
		err = fmt.Errorf("no servers found")
		return
	}

	_ = speedtest.SortServersByDistance(servers, cfg.Coordinates)
	if len(servers) > candidates {
		servers = servers[:candidates]
	}

	st.CandidateServers = servers
	return
}

func (d *DownloadRadarSign) Read(ctx context.Context) (speed speedtest.BytesPerSecond, err error) {
	if s, err := d.RepeatableSpeedTest.fastestServer(ctx); err == nil {
		speed, err = s.ProbeDownloadSpeed(ctx, &d.RepeatableSpeedTest.Client, nil)
	}
	return
}

func (u *UploadRadarSign) Read(ctx context.Context) (speed speedtest.BytesPerSecond, err error) {
	if s, err := u.RepeatableSpeedTest.fastestServer(ctx); err == nil {
		speed, err = s.ProbeUploadSpeed(ctx, &u.RepeatableSpeedTest.Client, nil)
	}
	return
}

// Picks the server with the lowest latency from the CandidateServers.
//
func (r *RepeatableSpeedTest) fastestServer(ctx context.Context) (speedtest.Server, error) {
	s := make([]speedtest.Server, len(r.CandidateServers))
	copy(s, r.CandidateServers)
	_, err := speedtest.StableSortServersByAverageLatency(s, ctx, &r.Client, speedtest.DefaultLatencySamples)
	return s[0], err
}

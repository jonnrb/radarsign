package main

import (
	"context"

	"go.jonnrb.io/speedtest/fastdotcom"
	"go.jonnrb.io/speedtest/units"
)

type RadarSign interface {
	Read(context.Context) (units.BytesPerSecond, error)
}

type RepeatableSpeedTest interface {
	ProbeDownloadSpeed(ctx context.Context) (units.BytesPerSecond, error)
	ProbeUploadSpeed(ctx context.Context) (units.BytesPerSecond, error)
}

type DownloadRadarSign struct {
	RepeatableSpeedTest
}

type UploadRadarSign struct {
	RepeatableSpeedTest
}

func NewRepeatableSpeedTest() RepeatableSpeedTest {
	return &FastDotComSpeedTest{}
}

func (d *DownloadRadarSign) Read(ctx context.Context) (units.BytesPerSecond, error) {
	return d.ProbeDownloadSpeed(ctx)
}

func (u *UploadRadarSign) Read(ctx context.Context) (units.BytesPerSecond, error) {
	return u.ProbeUploadSpeed(ctx)
}

type FastDotComSpeedTest struct {
	Client fastdotcom.Client
}

func (f *FastDotComSpeedTest) ProbeDownloadSpeed(ctx context.Context) (units.BytesPerSecond, error) {
	m, err := f.getManifest(ctx)
	if err != nil {
		return units.BytesPerSecond(0), err
	}
	return m.ProbeDownloadSpeed(ctx, &f.Client, nil)
}

func (f *FastDotComSpeedTest) ProbeUploadSpeed(ctx context.Context) (units.BytesPerSecond, error) {
	m, err := f.getManifest(ctx)
	if err != nil {
		return units.BytesPerSecond(0), err
	}
	return m.ProbeUploadSpeed(ctx, &f.Client, nil)
}

func (f *FastDotComSpeedTest) getManifest(ctx context.Context) (*fastdotcom.Manifest, error) {
	ctx, cancel := context.WithTimeout(ctx, *configureTimeout)
	defer cancel()
	return fastdotcom.GetManifest(ctx, *candidateServers)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	*http.ServeMux
	*prometheus.Registry
	Linearizer chan func()

	Download *RadarSignMetric
	Upload   *RadarSignMetric
}

func NewServer() *Server {
	s := &Server{
		ServeMux:   http.NewServeMux(),
		Linearizer: Linearizer(),
		Registry:   prometheus.NewRegistry(),
	}

	handler := promhttp.HandlerFor(s.Registry, promhttp.HandlerOpts{
		ErrorLog:      &glogger{},
		ErrorHandling: 0,
	})
	s.ServeMux.Handle("/metrics", handler)

	return s
}

func (s *Server) Close() {
	close(s.Linearizer)
}

type glogger struct{}

func (_ glogger) Println(v ...interface{}) {
	glog.Infoln(v...)
}

func (s *Server) Register(st RepeatableSpeedTest) error {
	if err := s.registerDownload(st); err != nil {
		return fmt.Errorf("registering download metric: %v", err)
	}
	if err := s.registerUpload(st); err != nil {
		return fmt.Errorf("registering upload metric: %v", err)
	}
	return nil
}

func (s *Server) registerDownload(st RepeatableSpeedTest) error {
	t := NewThrottler(s.Linearizer)
	t.Source = &DownloadRadarSign{st}

	s.Download = NewRadarSignMetric(t, downloadSpeedGaugeOpts)
	s.Download.Timeout = *downloadTime + *configureTimeout

	return s.Registry.Register(s.Download)
}

func (s *Server) registerUpload(st RepeatableSpeedTest) error {
	t := NewThrottler(s.Linearizer)
	t.Source = &UploadRadarSign{st}

	s.Upload = NewRadarSignMetric(t, uploadSpeedGaugeOpts)
	s.Upload.Timeout = *uploadTime + *configureTimeout

	return s.Registry.Register(s.Upload)
}

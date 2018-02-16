package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	s := NewServer()
	if err := s.Register(newRepeatableSpeedTestOrFail()); err != nil {
		glog.Fatalf("Error registering metrics: %v", err)
	}

	glog.V(1).Infof("Starting HTTP server listening on %q.", *addr)
	http.ListenAndServe(*addr, s)
}

func newRepeatableSpeedTestOrFail() RepeatableSpeedTest {
	ctx, cancel := context.WithTimeout(context.Background(), *configureTimeout)
	defer cancel()

	glog.V(1).Infof("Determining %d servers to use.", *candidateServers)

	st, err := NewRepeatableSpeedTest(ctx, *candidateServers)
	if err != nil {
		glog.Fatalf("Could not get initial speedtest.net configuration information: %v", err)
	}
	return st
}

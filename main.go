package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	s := NewServer()
	if err := s.Register(NewRepeatableSpeedTest()); err != nil {
		glog.Fatalf("Error registering metrics: %v", err)
	}

	glog.V(1).Infof("Starting HTTP server listening on %q.", *addr)
	http.ListenAndServe(*addr, s)
}

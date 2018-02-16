package main

import (
	"flag"
	"time"
)

var (
	addr = flag.String("addr", "0.0.0.0:8080",
		"an address, port pair the HTTP server will listen on")

	downloadTime = flag.Duration("time.download", 10*time.Second,
		"total time to spend trying to probe download speed")

	uploadTime = flag.Duration("time.upload", 10*time.Second,
		"total time to spend trying to probe upload speed")

	candidateServers = flag.Int("candidate_servers", 5,
		"number of servers whose latency gets checked during a speed probe")

	configureTimeout = flag.Duration("time.configure", 30*time.Second,
		"maximum time to spend getting the initial speedtest.net configuration")

	throttleTime = flag.Duration("time.throttle", 5*time.Minute,
		"how long to wait before allowing speed probe metrics to be refreshed")
)

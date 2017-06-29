package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"github.com/streadway/handy/report"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	start = time.Now()
)

func main() {
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api/", handleAPI)

	// Log every received HTTP request to stdout.
	go http.ListenAndServe(*addr, alice.New(
		report.JSONMiddleware(os.Stdout),
	).Then(http.DefaultServeMux))

	startClient(*addr)

	select {}
}

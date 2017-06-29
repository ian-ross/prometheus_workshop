package main

import (
	"bytes"
	"flag"
	"math"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var oscillationPeriod = flag.Duration("oscillation-period", 5*time.Minute, "The duration of the rate oscillation period.")

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP request counts",
		},
		[]string{"method", "url", "code"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds",
			Help: "HTTP request latencies",
			Buckets: []float64{0.05, },
		},
		[]string{"method", "url", "code"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpRequestDuration)
}

func incRequests(start time.Time, method string, code string, url string) {
	duration := float64(time.Since(start) / time.Millisecond)
	httpRequests.With(
		prometheus.Labels{"method": method, "url": url, "code": code},
	).Inc()
	httpRequestDuration.With(
		prometheus.Labels{"method": method, "url": url, "code": code},
	).Observe(duration)
}

func startClient(servAddr string) {

	oscillationFactor := func() float64 {
		return 2 + math.Sin(math.Sin(2*math.Pi*float64(time.Since(start))/float64(*oscillationPeriod)))
	}

	ignoreRequest := func(start time.Time, resp *http.Response, err error) {
		if err != nil {
			incRequests(start, resp.Request.Method, resp.Status, resp.Request.URL.Path)
			return
		}
		resp.Body.Close()
		incRequests(start, resp.Request.Method, resp.Status, resp.Request.URL.Path)
	}

	// GET /api/foo.
	go func() {
		for {
			start := time.Now()
			resp, err := http.Get("http://" + servAddr + "/api/foo")
			ignoreRequest(start, resp, err)
			time.Sleep(time.Duration(10*oscillationFactor()) * time.Millisecond)
		}
	}()
	// POST /api/foo.
	go func() {
		for {
			start := time.Now()
			resp, err := http.Post("http://"+servAddr+"/api/foo", "text/plain", &bytes.Buffer{})
			ignoreRequest(start, resp, err)
			time.Sleep(time.Duration(150*oscillationFactor()) * time.Millisecond)
		}
	}()
	// GET /api/bar.
	go func() {
		for {
			start := time.Now()
			resp, err := http.Get("http://" + servAddr + "/api/bar")
			ignoreRequest(start, resp, err)
			time.Sleep(time.Duration(20*oscillationFactor()) * time.Millisecond)
		}
	}()
	// POST /api/bar.
	go func() {
		for {
			start := time.Now()
			resp, err := http.Post("http://"+servAddr+"/api/bar", "text/plain", &bytes.Buffer{})
			ignoreRequest(start, resp, err)
			time.Sleep(time.Duration(100*oscillationFactor()) * time.Millisecond)
		}
	}()
	// GET /api/nonexistent.
	go func() {
		for {
			start := time.Now()
			resp, err := http.Get("http://" + servAddr + "/api/nonexistent")
			ignoreRequest(start, resp, err)
			time.Sleep(time.Duration(500*oscillationFactor()) * time.Millisecond)
		}
	}()
}

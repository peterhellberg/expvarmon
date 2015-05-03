package main

import (
	"expvar"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"
)

var startTime = time.Now().UTC()

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

// uptime is an expvar.Func compliant wrapper for uptime info.
func uptime() interface{} {
	uptime := time.Since(startTime)
	return int64(uptime)
}

// startPort defines lower port for bind
const startPort = 32768

// StartSelfMonitor starts http server on random port and exports expvars.
//
// It tries 1024 ports, starting from startPort and registers some expvars if ok.
func StartSelfMonitor() (string, error) {
	for port := startPort; port < startPort+1024; port++ {
		bind := fmt.Sprintf("localhost:%d", port)
		l, err := net.Listen("tcp", bind)
		if err != nil {
			continue
		}
		l.Close()

		expvar.Publish("Goroutines", expvar.Func(goroutines))
		expvar.Publish("Uptime", expvar.Func(uptime))
		go http.ListenAndServe(bind, nil)
		return bind, nil
	}

	return "", fmt.Errorf("no free ports found")
}
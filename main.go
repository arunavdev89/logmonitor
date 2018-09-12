package main

import (
	"fmt"
	"os"
	"os/signal"
)

const (
	AccessLogMetricsName = "access_log_stat"
)

func main() {
	cfg := NewConfig()
	server := NewHttpServer(cfg)

	done := make(chan bool)
	reportChan := make(chan *Stat)

	monitor := NewAccessLogMonitor(cfg.File, AccessLogMetricsName, dimensionNames, done)
	collector := NewAccessLogStatCollector(reportChan, done, AccessLogMetricsName, cfg)
	reporter := NewAccessLogReporter(cfg.PrometheusAddress, reportChan)
	go monitor.Monitor()
	go WaitForShutdown(done)
	go collector.Collect()
	go reporter.Report()
	server.Serve()

	fmt.Print("\nServer stopped!!")
	fmt.Print("\nWaiting for file monitor to complete or shutdown!!")
	<-done

	fmt.Print("\nShutdown!.\n")
}

//Listen for shitdown
func WaitForShutdown(done chan bool) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	for range sig {
		done <- true
	}
}

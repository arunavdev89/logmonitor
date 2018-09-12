package main

import (
	"flag"
)

const (
	//default arguments
	defaultStatDisplayInterval = 10  //10s
	defaultAlertInterval       = 120 //2min
	defaultHitsAlertThreshold  = 10  //alert after 10 hits/2min
	defaultTopK                = 10  //Number of top websites to display
	defaultFile                = "/var/log/access.log"
	defaultPort                = "8080"
	defaultPrometheusUrl       = "http://localhost:9090"
)

type Config struct {
	StatDisplayInterval int
	AlertInterval       int
	HitsThreshold       int
	File                string
	TopK                int
	Port                string
	PrometheusAddress   string
}

func NewConfig() *Config {
	sdi := flag.Int("stat-interval", defaultStatDisplayInterval, "Stat Display Interval")
	ai := flag.Int("alert-interval", defaultAlertInterval, "Stat Alert Interval")
	hi := flag.Int("hit-threshold", defaultHitsAlertThreshold, "Number of hits alert Interval")
	file := flag.String("log-file", defaultFile, "Log file to monitor")
	topk := flag.Int("topk", defaultTopK, "Top K sections to display")
	port := flag.String("port", defaultPort, "Http Server Port")
	prometheusUrl := flag.String("prometheus-url", defaultPrometheusUrl, "Prometheus Server Address")
	flag.Parse()

	//TODO: Perform some input validation
	return &Config{
		StatDisplayInterval: *sdi,
		AlertInterval:       *ai,
		HitsThreshold:       *hi,
		File:                *file,
		TopK:                *topk,
		Port:                *port,
		PrometheusAddress:   *prometheusUrl,
	}
}

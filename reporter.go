package main

import (
	"fmt"
	"io"
	"os"
)

type Reporter interface {
	Report()
}

type reporter struct {
	Source     string
	Sink       io.Writer
	ReportChan chan *Stat
}

func NewAccessLogReporter(prometheusUrl string, channel chan *Stat) Reporter {
	return &reporter{
		Source:     prometheusUrl,
		Sink:       os.Stdout,
		ReportChan: channel,
	}
}

//TODO: Make some better reporting - e.g. color the terminals to yellow/green based on
//alert type
func (r *reporter) Report() {
	for stat := range r.ReportChan {
		details, err := stat.Details()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(details)
	}
}

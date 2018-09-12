package main

import (
	"fmt"

	"github.com/hpcloud/tail"
)

type Monitor interface {
	Monitor()
}

type AccessLogMonitor struct {
	Target  string
	Metrics Metrics
	Done    chan bool
}

func NewAccessLogMonitor(file string, metricName string, dimensions []string, done chan bool) Monitor {
	return &AccessLogMonitor{
		Target:  file,
		Metrics: NewAccessLogMetrics(metricName, dimensions),
		Done:    done,
	}
}

func (m *AccessLogMonitor) Monitor() {
	defer func() { m.Done <- true }()
	t, err := tail.TailFile(m.Target,
		tail.Config{
			Follow: true,
			ReOpen: true,
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	for line := range t.Lines {
		if line.Err != nil {
			fmt.Println(line.Err)
			continue
		}
		entry, err := NewEntry(line.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if m.Metrics.AddMetrics(entry.Dimension(), nil) != err {
			fmt.Println(err)
			continue
		}
	}
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

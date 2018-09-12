package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type StatType int

const (
	trafficHigh StatType = iota
	trafficBackToNormal
	summary
)

type Stat struct {
	Type  StatType
	Stats map[string]string
	Time  time.Time
}

func NewTrafficHighStat(avgTraffic string) *Stat {
	return &Stat{
		Type:  trafficHigh,
		Stats: map[string]string{"avgTraffic": avgTraffic},
		Time:  time.Now(),
	}
}

func NewTrafficNormalStat() *Stat {
	return &Stat{
		Type:  trafficBackToNormal,
		Stats: map[string]string{},
		Time:  time.Now(),
	}
}

func NewTrafficSummaryStat(totalTraffic string, avgTraffic string, errorTraffic string, topK map[string]string) *Stat {
	topKHits := []string{}
	stats := make(map[string]string)
	stats["totalTraffic"] = totalTraffic
	stats["avgTraffic"] = avgTraffic
	stats["errorTraffic"] = errorTraffic

	for k, v := range topK {
		topKHits = append(topKHits, fmt.Sprintf("URL [%s] (Total Hits %s)", k, v))
	}
	fmt.Println("Top K Hits", topKHits)
	stats["topKWebsites"] = fmt.Sprintf("[%s]", strings.Join(topKHits, ","))

	return &Stat{
		Type:  summary,
		Stats: stats,
		Time:  time.Now(),
	}
}

func (s *Stat) Details() (string, error) {
	if s.Type == trafficHigh {
		return fmt.Sprintf("High traffic generated an alert - hits = {%s}, triggered at %s ", s.Stats["avgTraffic"], s.Time.Format(time.RFC3339)), nil
	}
	if s.Type == trafficBackToNormal {
		return fmt.Sprintf("High traffic alarm recovered - alarm back to normal %s", s.Time.Format(time.RFC3339)), nil
	}
	if s.Type == summary {
		details, _ := json.MarshalIndent(s.Stats, "", "  ")
		return fmt.Sprintf("Stat summary [%s]\n%s", s.Time.Format(time.RFC3339), details), nil
	}
	return "", fmt.Errorf("Unknown stat type.")
}

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arunavdev89/logmonitor/prometheus"
)

var (
	avgTrafficQueryFmt   = "sum(rate(%s[%ss]))"
	totalTrafficQueryFmt = "sum(increase(%s[%ss]))"
	errorTrafficQueryFmt = "sum(increase(%s{status!~\"2*\"}[%ss]))"
	topKQueryFmt         = "topk(%s,sum(increase(%s[%ss])) by (website))"
)

type Collector interface {
	Collect() error
}

type collector struct {
	ReportChan    chan *Stat
	DoneChan      chan bool
	Cfg           *Config
	MetricName    string
	DisplayTicker *time.Ticker
	AlertTicker   *time.Ticker
	StatClient    *prometheus.PrometheusClient
	State         StatType
}

func NewAccessLogStatCollector(reportChan chan *Stat, doneChan chan bool, metricName string, cfg *Config) Collector {
	displayTicker := time.NewTicker(time.Second * time.Duration(cfg.StatDisplayInterval))
	alertTicker := time.NewTicker(time.Duration(cfg.AlertInterval) * time.Second)
	statClient, err := prometheus.NewPrometheusClient(cfg.PrometheusAddress)

	if err != nil {
		fmt.Println(err)
		panic("Unable to create prometheus client")
	}

	return &collector{
		ReportChan:    reportChan,
		DoneChan:      doneChan,
		Cfg:           cfg,
		MetricName:    metricName,
		DisplayTicker: displayTicker,
		AlertTicker:   alertTicker,
		StatClient:    statClient,
	}
}

func (c *collector) Collect() error {
	for {
		select {
		case <-c.DoneChan:
			c.DisplayTicker.Stop()
			c.AlertTicker.Stop()
			close(c.ReportChan)
			break

		case <-c.DisplayTicker.C:
			stat, err := c.CollectSummary()
			if err != nil {
				fmt.Println(err)
				continue
			}
			c.ReportChan <- stat

		case <-c.AlertTicker.C:
			stat, err := c.CollectAlert()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if stat != nil {
				c.ReportChan <- stat
			}
		}
	}
}

func (c *collector) CollectSummary() (*Stat, error) {
	avgTrafficQuery := fmt.Sprintf(avgTrafficQueryFmt, c.MetricName, strconv.Itoa(c.Cfg.StatDisplayInterval))
	totalTrafficQuery := fmt.Sprintf(totalTrafficQueryFmt, c.MetricName, strconv.Itoa(c.Cfg.StatDisplayInterval))
	topKQuery := fmt.Sprintf(topKQueryFmt, strconv.Itoa(c.Cfg.TopK), c.MetricName, strconv.Itoa(c.Cfg.StatDisplayInterval))
	totalErrorQuery := fmt.Sprintf(errorTrafficQueryFmt, c.MetricName, strconv.Itoa(c.Cfg.StatDisplayInterval))
	instant := time.Now()

	//find avg traffic/s in last 10s
	r, err := c.StatClient.Query(avgTrafficQuery, instant)
	if err != nil {
		return nil, err
	}
	avgTraffic := GetFirstResultElement(r.Data.Result)

	//find total traffic in last 10s
	r, err = c.StatClient.Query(totalTrafficQuery, instant)
	if err != nil {
		return nil, err
	}
	totalTraffic := GetFirstResultElement(r.Data.Result)

	//find total traffic that was not 2XX status
	r, err = c.StatClient.Query(totalErrorQuery, instant)
	if err != nil {
		return nil, err
	}
	totalErrorTraffic := GetFirstResultElement(r.Data.Result)

	//find top K traffic websites
	r, err = c.StatClient.Query(topKQuery, instant)
	if err != nil {
		return nil, err
	}
	topKHitsResult := r.Data.Result
	topKHits := make(map[string]string)
	for _, result := range topKHitsResult {
		website := result.Metric["website"]
		value := result.Value()
		topKHits[website] = value
	}

	return NewTrafficSummaryStat(avgTraffic, totalTraffic, totalErrorTraffic, topKHits), nil
}

func (c *collector) CollectAlert() (*Stat, error) {
	avgTrafficQuery := fmt.Sprintf(avgTrafficQueryFmt, c.MetricName, strconv.Itoa(c.Cfg.StatDisplayInterval))
	instant := time.Now()
	//find avg traffic/s in last 10s
	r, err := c.StatClient.Query(avgTrafficQuery, instant)
	if err != nil {
		return nil, err
	}
	avgTraffic := GetFirstResultElement(r.Data.Result)
	avgTrafficFloat, err := strconv.ParseFloat(avgTraffic, 64)
	if err != nil {
		return nil, err
	}
	//check if threshold was crossed
	if int(avgTrafficFloat) > c.Cfg.HitsThreshold {
		c.State = trafficHigh
		return NewTrafficHighStat(avgTraffic), nil
	} else if c.State == trafficHigh {
		c.State = trafficBackToNormal
		return NewTrafficNormalStat(), nil
	}
	return nil, nil
}

func GetFirstResultElement(result []*prometheus.QueryResponseResult) string {
	if result != nil && len(result) > 0 {
		return result[0].Value()
	}
	return "0"
}

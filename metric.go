package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

//MetricDimension is a list of dimension name and their values. Example: dimensionNames: {"method", "status"} dimensionValues: {"GET", "200"}
type MetricDimension struct {
	dimensionNames  []string
	dimensionValues []string
}

//Return a new MetricDimension with just the names set.
//Useful for initial declaration
func NewMetricDimension(dimensions []string) *MetricDimension {
	return &MetricDimension{
		dimensionNames: dimensions,
	}
}

//Adds the dimension values and returns a new MetricDimension
func (d *MetricDimension) WithDimensionValues(values []string) *MetricDimension {
	return &MetricDimension{
		dimensionNames:  d.dimensionNames,
		dimensionValues: values,
	}
}

//Metrics interface implements one method AddMetrics on a given MetricDimension
type Metrics interface {
	AddMetrics(dimension *MetricDimension, value interface{}) error
}

//AccessLogMetrics implements Metrics interface.
type AccessLogMetrics struct {
	MetricDimensions []string
	Counter          *prometheus.CounterVec
}

func NewAccessLogMetrics(metricName string, dimensions []string) Metrics {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricName,
			Help: "A counter to count access log entries.",
		},
		dimensions,
	)
	prometheus.MustRegister(counter)

	return &AccessLogMetrics{
		MetricDimensions: dimensions,
		Counter:          counter,
	}
}

func (m *AccessLogMetrics) AddMetrics(dimension *MetricDimension, value interface{}) error {
	if len(dimension.dimensionNames) != len(dimension.dimensionValues) {
		return fmt.Errorf("Metric dimensions must be of same length")
	}
	labels := make(map[string]string)
	dimensionNames := dimension.dimensionNames
	dimensionValues := dimension.dimensionValues

	for i, _ := range dimensionNames {
		labels[dimensionNames[i]] = dimensionValues[i]
	}
	m.Counter.With(labels).Inc()
	return nil
}

// Package sql - metrics describes prom collector for getting metrics from sql.DB
package sql

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusCollector struct {
	db *sql.DB

	MaxOpenConnections *prometheus.Desc

	OpenConnections *prometheus.Desc
	InUse           *prometheus.Desc
	Idle            *prometheus.Desc

	WaitCount         *prometheus.Desc
	WaitDuration      *prometheus.Desc
	MaxIdleClosed     *prometheus.Desc
	MaxIdleTimeClosed *prometheus.Desc
	MaxLifetimeClosed *prometheus.Desc
}

// NewPrometheusCollector returns a collector that exports metrics about the db.
func NewPrometheusCollector(name string, db *sql.DB) prometheus.Collector {
	var labels = prometheus.Labels{
		"name": name,
	}

	return &prometheusCollector{
		db: db,

		MaxOpenConnections: prometheus.NewDesc(
			"max_open_connections",
			"Maximum number of open connections to the database",
			nil, labels,
		),
		OpenConnections: prometheus.NewDesc(
			"open_connections",
			"The number of established connections both in use and idle",
			nil, labels,
		),
		InUse: prometheus.NewDesc(
			"in_use_connections",
			"The number of connections currently in use",
			nil, labels,
		),
		Idle: prometheus.NewDesc(
			"idle_connections",
			"The number of idle connections",
			nil, labels,
		),

		WaitCount: prometheus.NewDesc(
			"wait_connections",
			"The total number of connections waited for",
			nil, labels,
		),
		WaitDuration: prometheus.NewDesc(
			"wait_duration_connections",
			"The total time blocked waiting for a new connection",
			nil, labels,
		),
		MaxIdleClosed: prometheus.NewDesc(
			"max_idle_closed_connections",
			"The total number of connections closed due to SetMaxIdleConns",
			nil, labels,
		),
		MaxIdleTimeClosed: prometheus.NewDesc(
			"max_idle_time_closed_connections",
			"The total number of connections closed due to SetConnMaxIdleTime",
			nil, labels,
		),
		MaxLifetimeClosed: prometheus.NewDesc(
			"max_lifetime_closed_connections",
			"The total number of connections closed due to SetConnMaxLifetime",
			nil, labels,
		),
	}
}

// Describe returns all descriptions of the collector.
func (c *prometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.MaxOpenConnections
	ch <- c.OpenConnections
	ch <- c.InUse
	ch <- c.Idle
	ch <- c.WaitCount
	ch <- c.WaitDuration
	ch <- c.MaxIdleClosed
	ch <- c.MaxIdleTimeClosed
	ch <- c.MaxLifetimeClosed
}

// Collect returns the current state of all metrics of the collector.
func (c *prometheusCollector) Collect(ch chan<- prometheus.Metric) {
	var stats = c.db.Stats()

	ch <- prometheus.MustNewConstMetric(c.MaxOpenConnections, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.OpenConnections, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.InUse, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.Idle, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.WaitCount, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.WaitDuration, prometheus.CounterValue, float64(stats.WaitDuration))
	ch <- prometheus.MustNewConstMetric(c.MaxIdleClosed, prometheus.CounterValue, float64(stats.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.MaxIdleTimeClosed, prometheus.CounterValue, float64(stats.MaxIdleTimeClosed))
	ch <- prometheus.MustNewConstMetric(c.MaxLifetimeClosed, prometheus.CounterValue, float64(stats.MaxLifetimeClosed))
}

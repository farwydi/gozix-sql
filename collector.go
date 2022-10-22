// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package sql

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusCollector struct {
	db *sql.DB

	maxOpenConnections *prometheus.Desc
	openConnections    *prometheus.Desc
	inUse              *prometheus.Desc
	idle               *prometheus.Desc
	waitCount          *prometheus.Desc
	waitDuration       *prometheus.Desc
	maxIdleClosed      *prometheus.Desc
	maxIdleTimeClosed  *prometheus.Desc
	maxLifetimeClosed  *prometheus.Desc
}

// newPrometheusCollector returns a collector that exports metrics about the db.
func newPrometheusCollector(name string, db *sql.DB) prometheus.Collector {
	var labels = prometheus.Labels{
		"name": name,
	}

	return &prometheusCollector{
		db: db,

		maxOpenConnections: prometheus.NewDesc(
			"max_open_connections",
			"Maximum number of open connections to the database",
			nil, labels,
		),
		openConnections: prometheus.NewDesc(
			"open_connections",
			"The number of established connections both in use and idle",
			nil, labels,
		),
		inUse: prometheus.NewDesc(
			"in_use_connections",
			"The number of connections currently in use",
			nil, labels,
		),
		idle: prometheus.NewDesc(
			"idle_connections",
			"The number of idle connections",
			nil, labels,
		),
		waitCount: prometheus.NewDesc(
			"wait_connections",
			"The total number of connections waited for",
			nil, labels,
		),
		waitDuration: prometheus.NewDesc(
			"wait_duration_connections",
			"The total time blocked waiting for a new connection",
			nil, labels,
		),
		maxIdleClosed: prometheus.NewDesc(
			"max_idle_closed_connections",
			"The total number of connections closed due to SetMaxIdleConns",
			nil, labels,
		),
		maxIdleTimeClosed: prometheus.NewDesc(
			"max_idle_time_closed_connections",
			"The total number of connections closed due to SetConnMaxIdleTime",
			nil, labels,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			"max_lifetime_closed_connections",
			"The total number of connections closed due to SetConnMaxLifetime",
			nil, labels,
		),
	}
}

// Describe returns all descriptions of the collector.
func (c *prometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConnections
	ch <- c.openConnections
	ch <- c.inUse
	ch <- c.idle
	ch <- c.waitCount
	ch <- c.waitDuration
	ch <- c.maxIdleClosed
	ch <- c.maxIdleTimeClosed
	ch <- c.maxLifetimeClosed
}

// Collect returns the current state of all metrics of the collector.
func (c *prometheusCollector) Collect(ch chan<- prometheus.Metric) {
	var stats = c.db.Stats()

	ch <- prometheus.MustNewConstMetric(c.maxOpenConnections, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.openConnections, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUse, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitCount, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDuration, prometheus.CounterValue, float64(stats.WaitDuration))
	ch <- prometheus.MustNewConstMetric(c.maxIdleClosed, prometheus.CounterValue, float64(stats.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.maxIdleTimeClosed, prometheus.CounterValue, float64(stats.MaxIdleTimeClosed))
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosed, prometheus.CounterValue, float64(stats.MaxLifetimeClosed))
}

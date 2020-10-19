package sql

import (
	"database/sql"
	"fmt"
	"strings"

	promBundle "github.com/gozix/prometheus"
	"github.com/gozix/viper/v2"
	"github.com/iqoption/nap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sarulabs/di/v2"

	"github.com/gozix/sql/v2/internal/metric"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Row is type alias of sql.Row.
	Row = sql.Row

	// Rows is type alias of sql.Row.
	Rows = sql.Rows
)

// BundleName is default definition name.
const BundleName = "sql"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Key implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cfg *viper.Viper
			if err = ctn.Fill(viper.BundleName, &cfg); err != nil {
				return nil, err
			}

			var promRegistry *prometheus.Registry
			if err = ctn.Fill(promBundle.DefRegistryName, &promRegistry); err != nil {
				return nil, err
			}

			// use this is hack, not UnmarshalKey
			// see https://github.com/spf13/viper/issues/188
			var (
				keys = cfg.Sub(BundleName).AllKeys()
				conf = make(Configs, len(keys))
			)

			for _, key := range keys {
				var name = strings.Split(key, ".")[0]
				if _, ok := conf[name]; ok {
					continue
				}

				var (
					c      Config
					suffix = fmt.Sprintf("%s.%s.", BundleName, name)
				)

				if cfg.IsSet(suffix + "nodes") {
					c.Nodes = cfg.GetStringSlice(suffix + "nodes")
				}

				if cfg.IsSet(suffix + "driver") {
					c.Driver = cfg.GetString(suffix + "driver")
				}

				if cfg.IsSet(suffix + "max_open_conns") {
					c.MaxOpenConns = cfg.GetInt(suffix + "max_open_conns")
				}

				if cfg.IsSet(suffix + "max_idle_conns") {
					c.MaxIdleConns = cfg.GetInt(suffix + "max_idle_conns")
				}

				if cfg.IsSet(suffix + "conn_max_lifetime") {
					c.ConnMaxLifetime = cfg.GetDuration(suffix + "conn_max_lifetime")
				}

				conf[name] = c
			}
			var (
				registry = NewRegistry(conf)
				c        prometheus.Collector
				dbs      *nap.DB
			)
			for name, _ := range conf {
				if dbs, err = registry.ConnectionWithName(name); err != nil {
					return nil, err
				}

				for i, db := range dbs.Databases() {
					n := fmt.Sprintf("%s_%d", name, i)
					c = metric.NewPrometheusCollector(n, db)
					promRegistry.MustRegister(c)
				}
			}

			return registry, nil
		},
		Close: func(obj interface{}) error {
			return obj.(*Registry).Close()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{"viper", "prometheus"}
}

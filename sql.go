// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package sql

import (
	"fmt"
	"strings"

	"github.com/gozix/di"
	"github.com/gozix/glue/v3"
	gzOTEL "github.com/gozix/otel"
	gzViper "github.com/gozix/viper/v3"

	"github.com/XSAM/otelsql"
	"github.com/iqoption/nap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Bundle implements the glue.Bundle interface.
type Bundle struct{}

// BundleName is default definition name.
const BundleName = "sql"

// Bundle implements glue.Bundle interface.
var _ glue.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder di.Builder) error {
	return builder.Provide(b.provideRegistry)
}

func (b *Bundle) DependsOn() []string {
	return []string{
		gzViper.BundleName,
		gzOTEL.BundleName,
	}
}

func (b *Bundle) provideRegistry(
	cfg *viper.Viper,
	registry *prometheus.Registry,
	tracerProvider trace.TracerProvider,
) (_ *Registry, _ func() error, err error) {
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

		c.AfterOpen = func(name string, db *nap.DB) {
			for i, dbItem := range db.Databases() {
				n := fmt.Sprintf("%s_%d", name, i)
				registry.MustRegister(
					newPrometheusCollector(n, dbItem),
				)

				err = otelsql.RegisterDBStatsMetrics(dbItem,
					otelsql.WithTracerProvider(tracerProvider),
					otelsql.WithAttributes(attribute.String("name", n)),
				)
				if err != nil {
					panic(err)
				}
			}
		}

		conf[name] = c
	}

	var sqlRegistry *Registry
	if sqlRegistry, err = NewRegistry(conf); err != nil {
		return nil, nil, err
	}

	var closer = func() error {
		return sqlRegistry.Close()
	}

	return sqlRegistry, closer, nil
}

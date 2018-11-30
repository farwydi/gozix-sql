package sql

import (
	"database/sql"

	"github.com/gozix/viper"
	"github.com/sarulabs/di"
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

			var conf Config
			if err = cfg.UnmarshalKey(BundleName, &conf); err != nil {
				return nil, err
			}

			return NewRegistry(conf), nil
		},
		Close: func(obj interface{}) error {
			return obj.(*Registry).Close()
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{"viper"}
}

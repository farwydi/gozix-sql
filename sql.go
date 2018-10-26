package sql

import (
	"github.com/gozix/viper"
	"github.com/sarulabs/di"
)

// Bundle implements the glue.Bundle interface.
type Bundle struct{}

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
	builder.Add(di.Def{
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

	return nil
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{"viper"}
}

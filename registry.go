package sql

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/iqoption/nap"
)

// DEFAULT is default connection name.
const DEFAULT = "default"

type (
	// Config is registry configuration item.
	Config struct {
		Nodes           []string      `json:"nodes"`
		Driver          string        `json:"driver"`
		MaxOpenConns    int           `json:"max_open_conns"`
		MaxIdleConns    int           `json:"max_idle_conns"`
		ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	}

	// Configs is registry configurations.
	Configs map[string]Config

	// DB is type alias to nap.DB.
	DB = nap.DB

	// Registry is database connection registry.
	Registry struct {
		mux  sync.Mutex
		dbs  map[string]*nap.DB
		conf Configs
	}
)

var (
	// ErrUnknownConnection is error triggered when connection with provided name not founded.
	ErrUnknownConnection = errors.New("unknown connection")
)

// NewRegistry is registry constructor.
func NewRegistry(conf Configs) *Registry {
	return &Registry{
		dbs:  make(map[string]*nap.DB, 4),
		conf: conf,
	}
}

// Close is method for close connections.
func (r *Registry) Close() (err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for key, db := range r.dbs {
		if err = db.Close(); err != nil {
			return err
		}

		delete(r.dbs, key)
	}

	return nil
}

// Connection is default connection getter.
func (r *Registry) Connection() (*nap.DB, error) {
	return r.ConnectionWithName(DEFAULT)
}

// ConnectionWithName is connection getter by name.
func (r *Registry) ConnectionWithName(name string) (_ *nap.DB, err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var db, initialized = r.dbs[name]
	if initialized {
		return db, nil
	}

	var value, exists = r.conf[name]
	if !exists {
		return nil, ErrUnknownConnection
	}

	if db, err = nap.Open(value.Driver, strings.Join(value.Nodes, ";")); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(value.MaxOpenConns)
	db.SetMaxIdleConns(value.MaxIdleConns)
	db.SetConnMaxLifetime(value.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	r.dbs[name] = db

	return db, nil
}

// Driver is default connection driver name getter.
func (r *Registry) Driver() (string, error) {
	return r.DriverWithName(DEFAULT)
}

// DriverWithName is driver name getter by name.
func (r *Registry) DriverWithName(name string) (string, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var value, exists = r.conf[name]
	if !exists {
		return "", ErrUnknownConnection
	}

	return value.Driver, nil
}

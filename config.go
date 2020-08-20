package sql

import (
	"github.com/go-ginger/models"
	"reflect"
)

type Config struct {
	models.IConfig

	Debug                 bool
	SqlDialect            string
	SqlConnectionString   string
	ConnectionPoolEnabled bool
	MaxIdleConnections    *int
	MaxOpenConnections    *int
	MaxLifetimeSeconds    *int
}

var config Config

func InitializeConfig(input interface{}) {
	v := reflect.Indirect(reflect.ValueOf(input))
	sqlDialect := v.FieldByName("SqlDialect")
	sqlConnectionString := v.FieldByName("SqlConnectionString")

	config = Config{
		SqlDialect:          sqlDialect.String(),
		SqlConnectionString: sqlConnectionString.String(),
	}
	debug := v.FieldByName("Debug")
	if !debug.IsZero() {
		config.Debug = debug.Bool()
	}
	connectionPoolEnabled := v.FieldByName("ConnectionPoolEnabled")
	if !connectionPoolEnabled.IsZero() {
		config.ConnectionPoolEnabled = connectionPoolEnabled.Bool()
	}
	maxIdleConnections := v.FieldByName("MaxIdleConnections")
	if !maxIdleConnections.IsZero() {
		config.MaxIdleConnections = maxIdleConnections.Interface().(*int)
	}
	maxOpenConnections := v.FieldByName("MaxOpenConnections")
	if !maxOpenConnections.IsZero() {
		config.MaxOpenConnections = maxOpenConnections.Interface().(*int)
	}
	maxLifetimeSeconds := v.FieldByName("MaxLifetimeSeconds")
	if !maxLifetimeSeconds.IsZero() {
		config.MaxLifetimeSeconds = maxLifetimeSeconds.Interface().(*int)
	}
}

package sql

import (
	"github.com/go-ginger/models"
	"reflect"
)

type Config struct {
	models.IConfig

	SqlDialect          string
	SqlConnectionString string
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
}

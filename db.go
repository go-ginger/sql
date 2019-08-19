package sql

import "github.com/jinzhu/gorm"

func GetDb() (db *gorm.DB, err error) {
	db, err = gorm.Open(config.SqlDialect, config.SqlConnectionString)
	return
}

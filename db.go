package sql

import (
	"github.com/go-ginger/models"
	"github.com/jinzhu/gorm"
)

func GetDb(request models.IRequest) (db *gorm.DB, closeAtEnd bool, err error) {
	if request != nil {
		tx := request.GetTemp("tx")
		if tx != nil {
			db = tx.(*gorm.DB)
			return
		}
	}
	closeAtEnd = true
	db, err = gorm.Open(config.SqlDialect, config.SqlConnectionString)
	return
}

package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/mts"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func (handler *DbHandler) FixDB(database *gorm.DB) {
	if !config.ConnectionPoolEnabled {
		log.Println("connection pool is disabled and you're trying to set fixed database")
		return
	}
	handler.GormDB = database
	return
}

func (handler *DbHandler) GetDb(request models.IRequest) (db *gorm.DB, closeAtEnd bool, err error) {
	if request != nil {
		tx := request.GetTemp("tx")
		if tx != nil {
			db = tx.(*gorm.DB)
			return
		}
	}
	if handler.GormDB != nil {
		db = handler.GormDB
		return
	}
	closeAtEnd = true
	db, err = gorm.Open(config.SqlDialect, config.SqlConnectionString)
	if err != nil {
		return
	}
	if db != nil {
		if config.Debug {
			db = db.Debug()
		}
		if !config.ConnectionPoolEnabled {
			log.Println("sql connection pool is disabled")
			return
		}
		closeAtEnd = false
		if config.MaxIdleConnections != nil {
			db.DB().SetMaxIdleConns(*config.MaxIdleConnections)
		}
		if config.MaxOpenConnections != nil {
			db.DB().SetMaxOpenConns(*config.MaxOpenConnections)
		}
		if config.MaxLifetimeSeconds != nil {
			db.DB().SetConnMaxLifetime(time.Second * time.Duration(*config.MaxLifetimeSeconds))
		}
	}
	return
}

func (handler *DbHandler) HandleRequestFilters(request models.IRequest, existingQuery *gorm.DB) (query *gorm.DB, err error) {
	query = existingQuery
	req := request.GetBaseRequest()
	var q interface{}
	var params []interface{}
	if req.Filters != nil {
		q, params = mts.Parse(*req.Filters)
	}
	if q != nil && params != nil {
		query = query.Where(q, params...)
	}
	return
}

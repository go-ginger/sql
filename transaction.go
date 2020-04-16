package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"github.com/jinzhu/gorm"
)

func (handler *DbHandler) StartTransaction(request models.IRequest) (err error) {
	db, _, err := GetDb(request)
	if err != nil {
		return
	}
	tx := db.Begin()
	request.SetTemp("tx", tx)
	return
}

func (handler *DbHandler) GetTransaction(request models.IRequest) (tx *gorm.DB) {
	itx := request.GetTemp("tx")
	if itx == nil {
		return
	}
	tx = itx.(*gorm.DB)
	return
}

func (handler *DbHandler) CommitTransaction(request models.IRequest) (err error) {
	tx := handler.GetTransaction(request)
	if tx == nil {
		err = errors.GetInternalServiceError(request)
		return
	}
	tx.Commit()
	request.SetTemp("tx", nil)
	return
}

func (handler *DbHandler) RollbackTransaction(request models.IRequest) (err error) {
	tx := handler.GetTransaction(request)
	if tx == nil {
		err = errors.GetInternalServiceError(request)
		return
	}
	tx.Rollback()
	request.SetTemp("tx", nil)
	return
}

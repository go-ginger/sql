package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
)

func (handler *DbHandler) Insert(request models.IRequest) (result models.IBaseModel, err error) {
	req := request.GetBaseRequest()
	db, closeAtEnd, err := handler.GetDb(request)
	if err != nil {
		return
	}
	defer func() {
		if closeAtEnd {
			e := db.Close()
			if e != nil {
				err = e
			}
		}
	}()
	dbc := db.Create(req.Body)
	if dbc.Error != nil {
		return nil, errors.HandleError(dbc.Error)
	}
	result = req.Body
	_, err = handler.BaseDbHandler.Insert(req)
	return
}

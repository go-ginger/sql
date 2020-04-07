package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
)

func (handler *DbHandler) Insert(request models.IRequest) (result interface{}, err error) {
	req := request.GetBaseRequest()
	db, err := GetDb()
	if err != nil {
		return
	}
	defer func() {
		e := db.Close()
		if e != nil {
			err = e
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

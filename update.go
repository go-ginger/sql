package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
)

func (handler *DbHandler) Update(request models.IRequest) (err error) {
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
	req := request.GetBaseRequest()
	model := handler.GetModelInstance()
	query := db.
		Model(model).
		Where("id=?", req.ID)

	dbc := query.Update(req.Body)
	if dbc.Error != nil {
		return errors.HandleError(dbc.Error)
	}
	if req.ExtraQuery != nil {
		for key, value := range req.ExtraQuery {
			query.UpdateColumn(key, value)
		}
	}
	//if dbc.RowsAffected == 0 {
	//	return models.GetError(models.NotFoundError)
	//}
	return handler.BaseDbHandler.Update(request)
}

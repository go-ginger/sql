package sql

import (
	"errors"
	"github.com/go-ginger/models"
	ge "github.com/go-ginger/models/errors"
)

func (handler *DbHandler) Delete(request models.IRequest) (err error) {
	db, closeDb, err := handler.GetDb(request)
	if err != nil {
		return
	}
	defer func() {
		if closeDb {
			e := db.Close()
			if e != nil {
				err = e
			}
		}
	}()
	req := request.GetBaseRequest()
	model := handler.GetModelInstance()
	query := db.Model(model)
	filtered := false
	if req.ID != nil {
		query = query.Where(`id=?`, req.ID)
		filtered = true
	}
	if req.Body != nil {
		query = query.Where(req.Body)
		filtered = true
	}
	if req.ExtraQuery != nil {
		for key, value := range req.ExtraQuery {
			query = query.Where(key, value)
		}
		filtered = true
	}
	if !filtered {
		return ge.HandleError(errors.New("you can not delete all items"))
	}
	dbc := query.Delete(model)
	if dbc.Error != nil {
		return ge.HandleError(dbc.Error)
	}
	if dbc.RowsAffected == 0 {
		return ge.GetError(request, ge.NotFoundError)
	}
	return
}

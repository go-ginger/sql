package sql

import (
	"github.com/go-ginger/models"
	ge "github.com/go-ginger/models/errors"
)

func (handler *DbHandler) Update(request models.IRequest) (err error) {
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
	req := request.GetBaseRequest()
	model := handler.GetModelInstance()
	query := db.
		Model(model)
	if req.ID != nil {
		query = query.Where("id=?", req.ID)
	}
	query, err = handler.HandleRequestFilters(request, query)
	if req.Body != nil {
		dbc := query.Update(req.Body)
		if dbc.Error != nil {
			return ge.HandleError(dbc.Error)
		}
	}
	iExtraUpdates := req.GetTemp("extra_updates")
	if iExtraUpdates != nil {
		extraUpdates := iExtraUpdates.([]interface{})
		dbc := query.Update(extraUpdates...)
		if dbc.Error != nil {
			return ge.HandleError(dbc.Error)
		}
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

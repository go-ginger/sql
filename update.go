package sql

import (
	"github.com/kulichak/models"
)

func (handler *DbHandler) Update(request models.IRequest) error {
	db, err := GetDb()
	if err != nil {
		return err
	}
	defer db.Close()
	req := request.GetBaseRequest()

	query := db.
		Model(req.Model).
		Where("id=?", req.ID)

	dbc := query.Update(req.Body)
	if dbc.Error != nil {
		return models.HandleError(dbc.Error)
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

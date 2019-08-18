package sql

import (
	"github.com/kulichak/models"
)

func (handler *DbHandler) Update(request *models.Request) error {
	db, err := GetDb()
	if err != nil {
		return err
	}
	defer db.Close()

	query := db.
		Model(request.Model).
		Where("id=?", request.ID)

	dbc := query.Update(request.Body)
	if dbc.Error != nil {
		return models.HandleError(dbc.Error)
	}
	return nil
}

package sql

import "github.com/kulichak/models"

func (handler *DbHandler) Insert(request *models.Request) (obj interface{}, err error) {
	db, err := GetDb()
	if err != nil {
		return
	}
	defer db.Close()
	dbc := db.Create(request.Body)
	if dbc.Error != nil {
		return nil, models.HandleError(dbc.Error)
	}
	obj = request.Body
	handler.BaseDbHandler.Insert(request)
	return
}

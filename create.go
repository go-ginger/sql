package sql

import "github.com/kulichak/models"

func (handler *DbHandler) Insert(request *models.Request) (obj interface{}, err error) {
	db, err := GetDb()
	if err != nil {
		return
	}
	defer db.Close()
	db.Create(request.Body)
	obj = request.Body
	return
}

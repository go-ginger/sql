package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/kulichak/dl"
	"github.com/kulichak/models"
	"github.com/kulichak/mts"
	"math"
)

type DbHandler struct {
	dl.BaseDbHandler

	Model  interface{}
}

func (handler *DbHandler) countRecords(db *gorm.DB, done chan bool, count *uint64) {
	db.Count(count)
	done <- true
}

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

func (handler *DbHandler) Paginate(request *models.Request) (*models.PaginateResult, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q, params := mts.Parse(*request.Filters)
	offset := (*request.Page - 1) * *request.PerPage

	//done := make(chan bool, 1)

	query := db.
		Model(handler.Model).
		Where(q, params...)

	var totalCount uint64
	//go handler.countRecords(query, done, &totalCount)

	var result = request.Result //helpers.CreateArray(reflect.TypeOf(handler.Model), 0)
	query.Limit(*request.PerPage).Offset(offset).Find(&result)
	//<-done

	return &models.PaginateResult{
		Items: result,
		Pagination: models.PaginationInfo{
			Page:       *request.Page,
			PerPage:    *request.PerPage,
			//HasNext:    uint64(len(result)) >= *request.PerPage,
			PageCount:  uint64(math.Ceil(float64(totalCount) / float64(*request.PerPage))),
			TotalCount: totalCount,
		},
	}, nil
}

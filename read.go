package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/kulichak/models"
	"github.com/kulichak/mts"
	"math"
)

func (handler *DbHandler) countRecords(db *gorm.DB, done chan bool, count *uint64) {
	db.Count(count)
	done <- true
}

func (handler *DbHandler) Paginate(request *models.Request) (*models.PaginateResult, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q, params := mts.Parse(*request.Filters)
	offset := (request.Page - 1) * request.PerPage

	done := make(chan bool, 1)

	query := db.
		Model(request.Models).
		Where(q, params...)

	var totalCount uint64
	go handler.countRecords(query, done, &totalCount)

	//var result = &request.Result //helpers.CreateArray(reflect.TypeOf(handler.Model), 0)
	if request.Sort != nil {
		for _, s := range *request.Sort {
			sort := s.Name
			if !s.Ascending{
				sort += " DESC"
			}
			query = query.Order(sort)
		}
	}
	query.Limit(request.PerPage).Offset(offset).Find(request.Models)
	<-done

	pageCount := uint64(math.Ceil(float64(totalCount) / float64(request.PerPage)))
	return &models.PaginateResult{
		Items: request.Models,
		Pagination: models.PaginationInfo{
			Page:       request.Page,
			PerPage:    request.PerPage,
			HasNext:    pageCount > request.Page,
			PageCount:  pageCount,
			TotalCount: totalCount,
		},
	}, nil
}


func (handler *DbHandler) Get(request *models.Request) (*models.IBaseModel, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q, params := mts.Parse(*request.Filters)

	query := db.
		Model(request.Model).
		Where(q, params...)

	query.Find(request.Model)

	return &request.Model, nil
}

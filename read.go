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

func (handler *DbHandler) Paginate(request models.IRequest) (*models.PaginateResult, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	req := request.GetBaseRequest()

	q, params := mts.Parse(*req.Filters)
	offset := (req.Page - 1) * req.PerPage

	done := make(chan bool, 1)

	query := db.
		Model(req.Models).
		Where(q, params...)

	var totalCount uint64
	go handler.countRecords(query, done, &totalCount)

	//var result = &request.Result //helpers.CreateArray(reflect.TypeOf(handler.Model), 0)
	if req.Sort != nil {
		for _, s := range *req.Sort {
			sort := s.Name
			if !s.Ascending {
				sort += " DESC"
			}
			query = query.Order(sort)
		}
	}
	query.Limit(req.PerPage).Offset(offset).Find(req.Models)
	<-done

	pageCount := uint64(math.Ceil(float64(totalCount) / float64(req.PerPage)))
	return &models.PaginateResult{
		Items: req.Models,
		Pagination: models.PaginationInfo{
			Page:       req.Page,
			PerPage:    req.PerPage,
			PageCount:  pageCount,
			TotalCount: totalCount,
		},
	}, nil
}

func (handler *DbHandler) Get(request models.IRequest) (models.IBaseModel, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	req := request.GetBaseRequest()

	var q interface{}
	var params []interface{}
	if req.Filters != nil {
		q, params = mts.Parse(*req.Filters)
	}
	query := db.Model(&req.Model)
	if q != nil && params != nil {
		query = query.Where(q, params...)
	}

	dbc := query.Find(req.Model)
	if dbc.Error != nil {
		if dbc.RecordNotFound() {
			return nil, models.GetError(models.NotFoundError)
		}
		return nil, models.HandleError(dbc.Error)
	}

	return req.Model, nil
}

func (handler *DbHandler) Select(request models.IRequest, tableName string, selectQuery string,
	dest interface{}, args ...interface{}) (err error) {
	db, err := GetDb()
	if err != nil {
		return
	}
	defer db.Close()
	req := request.GetBaseRequest()
	var q interface{}
	var params []interface{}
	if req.Filters != nil {
		q, params = mts.Parse(*req.Filters)
	}
	query := db.Table(tableName).Select(selectQuery, args...)
	if q != nil && params != nil {
		query = query.Where(q, params...)
	}
	query.Scan(dest)
	return
}

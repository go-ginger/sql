package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"github.com/go-ginger/mts"
	"github.com/jinzhu/gorm"
	"math"
)

func (handler *DbHandler) countRecords(db *gorm.DB, done chan bool, count *uint64) {
	db.Count(count)
	done <- true
}

func (handler *DbHandler) Paginate(request models.IRequest) (*models.PaginateResult, error) {
	db, closeAtEnd, err := GetDb(request)
	if err != nil {
		return nil, err
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

	var q interface{}
	var params []interface{}
	if req.Filters != nil {
		q, params = mts.Parse(*req.Filters)
	}
	offset := (req.Page - 1) * req.PerPage

	done := make(chan bool, 1)

	model := handler.GetModelInstance()
	query := db.Model(model)

	if q != nil && params != nil {
		query = query.Where(q, params...)
	}

	var totalCount uint64
	go handler.countRecords(query, done, &totalCount)

	if req.Sort != nil {
		for _, s := range *req.Sort {
			sort := s.Name
			if !s.Ascending {
				sort += " DESC"
			}
			query = query.Order(sort)
		}
	}
	items := handler.GetModelsInstancePtr()
	query.Limit(req.PerPage).Offset(offset).Find(items)
	<-done

	pageCount := uint64(math.Ceil(float64(totalCount) / float64(req.PerPage)))
	return &models.PaginateResult{
		Items: items,
		Pagination: models.PaginationInfo{
			Page:       req.Page,
			PerPage:    req.PerPage,
			PageCount:  pageCount,
			TotalCount: totalCount,
			HasNext:    req.Page < pageCount,
		},
	}, nil
}

func (handler *DbHandler) Get(request models.IRequest) (result models.IBaseModel, err error) {
	db, closeAtEnd, err := GetDb(request)
	if err != nil {
		return nil, err
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

	var q interface{}
	var params []interface{}
	if req.Filters != nil {
		q, params = mts.Parse(*req.Filters)
	}
	model := handler.GetModelInstance()
	query := db.Model(model)
	if q != nil && params != nil {
		query = query.Where(q, params...)
	}

	dbc := query.Find(model)
	if dbc.Error != nil {
		if dbc.RecordNotFound() {
			return nil, errors.GetError(request, errors.NotFoundError)
		}
		return nil, errors.HandleError(dbc.Error)
	}
	result = model.(models.IBaseModel)
	return
}

func (handler *DbHandler) Select(request models.IRequest, tableName string, selectQuery string,
	dest interface{}, args ...interface{}) (err error) {
	db, closeAtEnd, err := GetDb(request)
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

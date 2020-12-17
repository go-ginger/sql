package sql

import (
	"github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"github.com/go-ginger/mts"
	"math"
	"reflect"
)

//func (handler *DbHandler) countRecords(db *gorm.DB, done chan bool, count *uint64) {
//	db.Count(count)
//	done <- true
//}

func (handler *DbHandler) Paginate(request models.IRequest) (*models.PaginateResult, error) {
	db, closeAtEnd, err := handler.GetDb(request)
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
	offset := (req.Page - 1) * req.PerPage

	//done := make(chan bool, 1)

	model := handler.GetModelInstance()
	query := db.Model(model)

	query, err = handler.HandleRequestFilters(request, query)

	var totalCount uint64
	if ignoreCount := request.GetTemp("ignore_count"); ignoreCount == nil || !ignoreCount.(bool) {
		//go handler.countRecords(query, done, &totalCount)
		query.Count(&totalCount)
	}

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
	query = query.Limit(req.PerPage).Offset(offset).Find(items)
	if query.Error != nil {
		return nil, errors.HandleError(query.Error)
	}
	//<-done

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
	db, closeAtEnd, err := handler.GetDb(request)
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
	model := handler.GetModelInstance()
	query := db.Model(model)
	query, err = handler.HandleRequestFilters(request, query)

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
	var q interface{}
	var params []interface{}
	if req.Filters != nil {
		q, params = mts.Parse(*req.Filters)
	}
	query := db.Table(tableName).Select(selectQuery, args...)
	if q != nil && params != nil {
		query = query.Where(q, params...)
	}
	query = query.Scan(dest)
	if query.Error != nil {
		return errors.HandleError(query.Error)
	}
	return
}

func (handler *DbHandler) First(request models.IRequest) (result models.IBaseModel, err error) {
	request.SetTemp("ignore_count", true)
	req := request.GetBaseRequest()
	req.Page = 1
	req.PerPage = 1
	pr, e := handler.DoPaginate(req)
	if e != nil {
		err = e
		return
	}
	arrValue := reflect.ValueOf(pr.Items)
	if arrValue.Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Len() > 0 {
		ind0 := arrValue.Index(0)
		if ind0.Kind() != reflect.Ptr {
			ind0 = ind0.Addr()
		}
		result = ind0.Interface().(models.IBaseModel)
	}
	return
}

func (handler *DbHandler) Exists(request models.IRequest) (exists bool, err error) {
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
	var totalCount uint64
	model := handler.GetModelInstance()
	query := db.Model(model)
	query, err = handler.HandleRequestFilters(request, query)

	dbc := query.Count(&totalCount)
	if dbc.Error != nil {
		if dbc.RecordNotFound() {
			return false, errors.GetError(request, errors.NotFoundError)
		}
		return false, errors.HandleError(dbc.Error)
	}
	if totalCount > 0 {
		exists = true
	}
	return
}

func (handler *DbHandler) Count(request models.IRequest) (count uint64, err error) {
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
	var totalCount uint64
	model := handler.GetModelInstance()
	query := db.Model(model)
	query, err = handler.HandleRequestFilters(request, query)

	dbc := query.Count(&totalCount)
	if dbc.Error != nil {
		if dbc.RecordNotFound() {
			return 0, errors.GetError(request, errors.NotFoundError)
		}
		return 0, errors.HandleError(dbc.Error)
	}
	if totalCount > 0 {
		count = totalCount
	}
	return
}

package sql

import (
	"github.com/go-ginger/dl"
	"github.com/jinzhu/gorm"
)

type DbHandler struct {
	dl.BaseDbHandler
	GormDB *gorm.DB
}

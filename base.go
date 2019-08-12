package sql

import (
	"github.com/jinzhu/gorm"
	"time"
)

type BaseModel struct {
	gorm.Model
	id         int64
	CreateDate time.Time `json:"create_date,omitempty"`
	ModifyDate time.Time `json:"modify_date,omitempty"`
}

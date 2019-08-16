package sql

import (
	"time"
)

type IBaseModel interface {

}

type BaseModel struct {
	IBaseModel

	ID        uint64     `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/kulichak/models"
	"time"
)

type BaseModel struct {
	models.BaseModel
	gorm.Model

	ID        uint64     `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

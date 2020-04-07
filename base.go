package sql

import (
	"github.com/go-ginger/models"
	"time"
)

type BaseModel struct {
	models.BaseModel `json:"-" gorm:"-"`

	ID        uint64     `json:"id,omitempty" gorm:"column:id;primary_key,AUTO_INCREMENT"`
	CreatedAt time.Time  `json:"created_at,omitempty" sql:"index" gorm:"column:created_at;not null"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (base *BaseModel) updateFromBase() {
	base.CreatedAt = base.BaseModel.CreatedAt
	base.UpdatedAt = base.BaseModel.UpdatedAt
	base.DeletedAt = base.BaseModel.DeletedAt
}

func (base *BaseModel) HandleCreateDefaultValues() {
	base.BaseModel.HandleCreateDefaultValues()
	base.updateFromBase()
}

func (base *BaseModel) HandleUpdateDefaultValues() {
	base.BaseModel.HandleUpdateDefaultValues()
	base.updateFromBase()
}

func (base *BaseModel) HandleDeleteDefaultValues() {
	base.BaseModel.HandleDeleteDefaultValues()
	base.updateFromBase()
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	Id int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

type ResourceInfo struct {
	Model
	Name      string         `gorm:"not null;size:128;unique;COMMENT:resource name" json:"name"`
	Status    bool           `gorm:"not null;COMMENT:resource status； true：open" json:"status"`
	Url       string         `gorm:"not null;COMMENT:resource base url" json:"url"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

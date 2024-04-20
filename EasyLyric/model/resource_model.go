package model

import "gorm.io/gorm"

type Model struct {
	Id int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

type ResourceInfo struct {
	Model
	Name      string `gorm:"not null;COMMENT:resource name"`
	Status    bool   `gorm:"not null;COMMENT:resource status； true：open"`
	Url       string `gorm:"not null;COMMENT:resource base url"`
	DeletedAt gorm.DeletedAt
}

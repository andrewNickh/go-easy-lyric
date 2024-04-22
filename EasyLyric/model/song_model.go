package model

import (
	"gorm.io/gorm"
	"time"
)

type Song struct {
	Model
	ResourceId int64          `gorm:"not null;index:idx_resource_id" json:"resourceId"`
	Title      string         `gorm:"not null;size:128;unique;COMMENT:song title" json:"title"`
	Url        string         `gorm:"not null;COMMENT:song url" json:"url"`
	Lyric      string         `gorm:"not null;COMMENT:song lyric" json:"lyric"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

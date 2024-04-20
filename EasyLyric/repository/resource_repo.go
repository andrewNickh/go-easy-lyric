package repository

import (
	"easy-lyric/EasyLyric/model"
	"gorm.io/gorm"
)

var ResourceRepo = new(resourceRepo)

type resourceRepo struct {
}

func (r *resourceRepo) Take(db *gorm.DB, where ...interface{}) *model.ResourceInfo {
	ret := &model.ResourceInfo{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

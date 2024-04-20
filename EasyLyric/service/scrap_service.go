package service

import (
	"easy-lyric/EasyLyric/model"
	"easy-lyric/EasyLyric/repository"
	"easy-lyric/util/db"
)

var ScrapService = new(scrapService)

type scrapService struct {
}

func (s *scrapService) GetResource(id int64) *model.ResourceInfo {
	return repository.ResourceRepo.Take(db.DBSlave(), "status = ? AND id = ?", true, id)
}

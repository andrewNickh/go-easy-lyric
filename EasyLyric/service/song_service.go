package service

import (
	"easy-lyric/EasyLyric/model"
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/util/db"
	"fmt"
)

var SongService = new(songService)

type songService struct {
}

func (s *songService) GetSongService(title string, page, limit int) (song []*model.Song, err error) {
	offset := (page - 1) * limit
	whereClause := fmt.Sprintf("title LIKE '%%%s%%'", title)
	err = db.DBSlave().Where(whereClause).Offset(offset).Limit(limit).Find(&song).Error
	return
}

func (s *songService) CreateSongService(req request.SaveSongReq) error {
	song := new(model.Song)
	song.ResourceId = req.ResourceId
	song.Title = req.Title
	song.Url = req.Url
	song.Lyric = req.Lyric

	err := db.Master().Create(&song).Error
	return err
}

func (s *songService) UpdateSongService(req request.UpdateSongReq) error {
	var existingSong *model.Song
	if err := db.DBSlave().First(&existingSong, req.Id).Error; err != nil { // todo: check from redis
		return err
	}

	existingSong.Title = req.SaveSongReq.Title
	existingSong.Url = req.SaveSongReq.Url
	existingSong.Lyric = req.SaveSongReq.Lyric

	if err := db.Master().Save(&existingSong).Error; err != nil {
		return err
	}

	return nil
}

func (s *songService) DeleteSongService(songId int64) error {
	var song *model.Song
	err := db.Master().Debug().Delete(&song, "id = ?", songId).Error
	if err != nil {
		return err
	}
	return nil
}

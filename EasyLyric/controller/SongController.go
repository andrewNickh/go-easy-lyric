package controller

import (
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/EasyLyric/model/response"
	"easy-lyric/EasyLyric/service"
	"easy-lyric/util/log"
	"github.com/kataras/iris/v12"
)

var SongController = new(songController)

type songController struct {
}

func (s *songController) SearchSong(ctx iris.Context) {
	req := request.JsonBodyToMap(ctx)
	title := request.JSONValueString(req, "title")
	page, limit := request.JSONValuePageInfo(req)

	songs, err := service.SongService.GetSongService(title, page, limit)
	if err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to get song", ctx)
		return
	}
	paging := &response.Paging{
		Page:  page,
		Limit: limit,
		Total: len(songs),
	}

	response.OkWithPagination("ok", songs, paging, ctx)
}

func (s *songController) CreateSong(ctx iris.Context) {
	var req request.SaveSongReq
	err := request.ReadBody(ctx, &req)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	if err = service.SongService.CreateSongService(req); err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to create song", ctx)
		return
	}
	response.OkWithMessageV2("", nil, ctx)
}

func (s *songController) UpdateSong(ctx iris.Context) {
	var req request.UpdateSongReq
	err := request.ReadBody(ctx, &req)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	if err = service.SongService.UpdateSongService(req); err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to update song", ctx)
		return
	}
	response.OkWithMessageV2("", nil, ctx)
}

func (s *songController) DeleteSong(ctx iris.Context) {
	var req request.GetById
	err := request.ReadBody(ctx, &req)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	if err = service.SongService.DeleteSongService(req.Id); err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to delete song", ctx)
		return
	}
	response.OkWithMessageV2("", nil, ctx)
}

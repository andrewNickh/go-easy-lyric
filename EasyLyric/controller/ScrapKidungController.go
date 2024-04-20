package controller

import (
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/EasyLyric/model/response"
	"easy-lyric/EasyLyric/service"
	"github.com/kataras/iris/v12"
)

var ScrapController = new(scrapController)

type scrapController struct {
}

func (s *scrapController) GetScrappedLyric(ctx iris.Context) {
	req := request.JsonBodyToMap(ctx)
	page, limit := request.JSONValuePageInfo(req)
	keyword := request.JSONValueString(req, "keyword")

	songs, total, err := service.ScrapService.GetScrapService(keyword, page, limit)
	if err != nil {
		response.FailWithMessageV2(err.Error(), ctx)
		return
	}

	paging := &response.Paging{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	response.OkWithPagination("ok", songs, paging, ctx)
}

package controller

import (
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/EasyLyric/model/response"
	"easy-lyric/EasyLyric/resources"
	"easy-lyric/EasyLyric/service"
	"github.com/kataras/iris/v12"
)

var ScrapController = new(scrapController)

type scrapController struct {
}

func (s *scrapController) GetLyrics(ctx iris.Context) {
	var req request.ScrapReq
	err := request.ReadBody(ctx, &req)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	src := service.ScrapService.GetResource(req.ResourceId)

	base := resources.Get(src.Name)
	if base == nil {
		response.FailWithMessageV2("invalid resource", ctx)
		return
	}

	songs, total, err := base.Scrape(req)
	if err != nil {
		response.FailWithMessageV2(err.Error(), ctx)
		return
	}

	paging := &response.Paging{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	response.OkWithPagination("ok", songs, paging, ctx)
}

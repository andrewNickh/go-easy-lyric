package controller

import (
	"easy-lyric/EasyLyric/model"
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/EasyLyric/model/response"
	"easy-lyric/EasyLyric/service"
	"easy-lyric/util/log"
	"github.com/kataras/iris/v12"
)

var ResourceController = new(resourceController)

type resourceController struct {
}

func (r *resourceController) GetResourceList(ctx iris.Context) {
	req := request.JsonBodyToMap(ctx)
	name := request.JSONValueString(req, "name")
	status := request.JSONValueBoolDefault(req, "status", true)
	page, limit := request.JSONValuePageInfo(req)
	log.Debug(status)

	list, err := service.ResourceService.GetResourceListService(name, status, page, limit)
	if err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to get resource list", ctx)
		return
	}
	paging := &response.Paging{
		Page:  page,
		Limit: limit,
		Total: len(list),
	}
	response.OkWithPagination("ok", list, paging, ctx)
}

func (r *resourceController) GetResource(ctx iris.Context) {
	req := request.JsonBodyToMap(ctx)
	rscId := request.JSONValueInt64Default(req, "id", 0)

	src, err := service.ResourceService.GetResourceService(rscId)
	if err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to get resource record", ctx)
		return
	}
	response.OkWithMessageV2("ok", src, ctx)
}

func (r *resourceController) CreateResource(ctx iris.Context) {
	var req request.CreateResourceReq
	err := request.ReadBody(ctx, &req)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	err = service.ResourceService.CreateResourceService(req)
	if err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to create resource record", ctx)
		return
	}
	response.OkWithMessageV2("", nil, ctx)
}

func (r *resourceController) UpdateResource(ctx iris.Context) {
	var rsc model.ResourceInfo
	err := request.ReadBody(ctx, &rsc)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	err = service.ResourceService.UpdateResourceService(rsc)
	if err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to update resource record", ctx)
		return
	}
	response.OkWithMessageV2("", nil, ctx)
}

func (r *resourceController) DeleteResource(ctx iris.Context) {
	var req request.GetById
	err := request.ReadBody(ctx, &req)
	if err != nil {
		response.FailWithMessageV2("failed to parse JSON body", ctx)
		return
	}

	err = service.ResourceService.DeleteResourceService(req.Id)
	if err != nil {
		log.Error(err)
		response.FailWithMessageV2("failed to delete resource record", ctx)
		return
	}
	response.OkWithMessageV2("", nil, ctx)
}

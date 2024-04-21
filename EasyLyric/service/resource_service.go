package service

import (
	"easy-lyric/EasyLyric/model"
	"easy-lyric/EasyLyric/model/request"
	"easy-lyric/util/db"
	"easy-lyric/util/log"
	"encoding/json"
	"fmt"
)

var ResourceService = new(resourceService)

type resourceService struct {
}

func (r *resourceService) GetResourceListService(name string, status bool, page, limit int) ([]*model.ResourceInfo, error) {
	offset := (page - 1) * limit
	var resources []*model.ResourceInfo
	var whereClause string
	whereClause = fmt.Sprintf("status = %t", status)
	if name != "" {
		whereClause += fmt.Sprintf(" AND name LIKE '%%%s%%'", name)
	}
	err := db.DBSlave().Where(whereClause).
		Offset(offset).Limit(limit).Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *resourceService) GetResourceService(id int64) (*model.ResourceInfo, error) {
	var rsc *model.ResourceInfo
	if err := db.DBSlave().First(&rsc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return rsc, nil
}

func (r *resourceService) CreateResourceService(req request.CreateResourceReq) error {
	src := &model.ResourceInfo{
		Name:   req.Name,
		Url:    req.Url,
		Status: req.Status,
	}
	return db.Master().Create(&src).Error
}

func (r *resourceService) UpdateResourceService(rsc model.ResourceInfo) error {
	var existingResource *model.ResourceInfo
	if err := db.DBSlave().First(&existingResource, rsc.Id).Error; err != nil {
		return err
	}

	existingResource.Name = rsc.Name
	existingResource.Status = rsc.Status
	existingResource.Url = rsc.Url

	err := db.Master().Save(&existingResource).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *resourceService) DeleteResourceService(id int64) error {
	var rsc *model.ResourceInfo
	err := db.Master().Debug().Delete(&rsc, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func PrettyPrintJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error(err)
	}
	log.Info(string(jsonData))
}

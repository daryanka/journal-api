package services

import (
	"api/clients"
	"api/domains/entry"
	"api/utils"
	"api/utils/xerror"
	"github.com/daryanka/myorm"
)

type EntryServiceI interface {
	CreateEntry(request entry.CreateEntryRequest) (int64, *xerror.XerrorT)
	UpdateEntry(userID int64, request entry.UpdateEntryRequest) *xerror.XerrorT
	DeleteEntry(id, userID int64) *xerror.XerrorT
}

type entryService struct{}

var EntryService EntryServiceI = &entryService{}

func (i *entryService) CreateEntry(req entry.CreateEntryRequest) (int64, *xerror.XerrorT) {
	id, err := clients.ClientOrm.Table("entries").Insert(myorm.H{
		"user_id":     req.UserID,
		"day":         req.Day,
		"start_time":  req.StartTime,
		"end_time":    req.EndTime,
		"title":       req.Title,
		"description": req.Description,
	})

	if err != nil {
		utils.ErrorLogger(err)
		return 0, xerror.NewInternalError("server error")
	}

	return id, nil
}

func (i *entryService) UpdateEntry(userID int64, req entry.UpdateEntryRequest) *xerror.XerrorT {
	err := clients.ClientOrm.Table("entries").
		Where("id", "=", req.ID).
		Where("user_id", "=", userID).
		Update(myorm.H{
			"start_time":  req.StartTime,
			"end_time":    req.EndTime,
			"title":       req.Title,
			"description": req.Description,
		})

	if err != nil {
		utils.ErrorLogger(err)
		return xerror.NewInternalError("server error")
	}

	return nil
}

func (i *entryService) DeleteEntry(id, userID int64) *xerror.XerrorT {
	err := clients.ClientOrm.Table("entries").
		Where("id", "=", id).
		Where("user_id", "=", userID).
		Delete()

	if err != nil {
		utils.ErrorLogger(err)
		return xerror.NewInternalError("server error")
	}

	return nil
}

package services

import (
	"api/clients"
	"api/domains/entry"
	"api/utils"
	"api/utils/xerror"
	"database/sql"
	"github.com/daryanka/myorm"
)

type EntryServiceI interface {
	CreateEntry(request entry.CreateEntryRequest) (int64, *xerror.XerrorT)
	UpdateEntry(userID int64, request entry.UpdateEntryRequest) *xerror.XerrorT
	DeleteEntry(id, userID int64) *xerror.XerrorT
	ViewRangeEntries(from, to string, userID int64) ([]entry.Entry, *xerror.XerrorT)
	ViewDayEntries(day string, userID int64) ([]entry.Entry, *xerror.XerrorT)
}

type entryService struct{}

var EntryService EntryServiceI = &entryService{}

func (i *entryService) CreateEntry(req entry.CreateEntryRequest) (int64, *xerror.XerrorT) {
	// Validate tag exists and belongs to user
	if req.TagID != nil {
		var tagID int64

		err := clients.ClientOrm.Table("tags").
			Select("id").
			WhereRaw("id = ? AND (user_id = ? OR user_id IS NULL)", req.UserID, req.TagID).
			First(&tagID)

		if err != nil {
			return 0, xerror.NewBadRequest("invalid tag")
		}
	}

	id, err := clients.ClientOrm.Table("entries").Insert(myorm.H{
		"user_id":     req.UserID,
		"day":         req.Day,
		"start_time":  req.StartTime,
		"end_time":    req.EndTime,
		"title":       req.Title,
		"description": req.Description,
		"tag_id":      req.TagID,
	})

	if err != nil {
		utils.ErrorLogger(err)
		return 0, xerror.NewInternalError("server error")
	}

	return id, nil
}

func (i *entryService) UpdateEntry(userID int64, req entry.UpdateEntryRequest) *xerror.XerrorT {
	if req.TagID != nil {
		var tagID int64

		err := clients.ClientOrm.Table("tags").
			Select("id").
			WhereRaw("id = ? AND (user_id = ? OR user_id IS NULL)", userID, req.TagID).
			First(&tagID)

		if err != nil {
			return xerror.NewBadRequest("invalid tag")
		}
	}

	err := clients.ClientOrm.Table("entries").
		Where("id", "=", req.ID).
		Where("user_id", "=", userID).
		Update(myorm.H{
			"start_time":  req.StartTime,
			"end_time":    req.EndTime,
			"title":       req.Title,
			"description": req.Description,
			"tag_id":      req.TagID,
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

func (i *entryService) ViewDayEntries(day string, userID int64) ([]entry.Entry, *xerror.XerrorT) {
	result := []entry.Entry{}

	err := clients.ClientOrm.Table("entries").
		Select(
			"entries.id AS id",
			"day",
			"start_time",
			"end_time",
			"title",
			"description",
			"tag_id",
			"tag_name",
		).
		LeftJoin("tags", "tags.id", "=", "entries.tag_id").
		Where("day", "=", day).
		Where("entries.user_id", "=", userID).
		Get(&result)

	if err != nil && err != sql.ErrNoRows {
		utils.ErrorLogger(err)
		return nil, xerror.NewInternalError("server error")
	}

	return result, nil
}

func (i *entryService) ViewRangeEntries(from, to string, userID int64) ([]entry.Entry, *xerror.XerrorT) {
	result := []entry.Entry{}

	err := clients.ClientOrm.Table("entries").
		Select(
			"entries.id AS id",
			"day",
			"start_time",
			"end_time",
			"title",
			"description",
			"tag_id",
			"tag_name",
		).
		LeftJoin("tags", "tags.id", "=", "entries.tag_id").
		Where("day", ">=", from).
		Where("day", "<=", to).
		Where("entries.user_id", "=", userID).
		Get(&result)

	if err != nil && err != sql.ErrNoRows {
		utils.ErrorLogger(err)
		return nil, xerror.NewInternalError("server error")
	}

	return result, nil
}

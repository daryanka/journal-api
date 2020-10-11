package services

import (
	"api/clients"
	"api/domains/tag"
	"api/utils"
	"api/utils/xerror"
	"database/sql"
	"github.com/daryanka/myorm"
)

type TagServiceI interface {
	CreateTag(request tag.CreateTagRequest, userID int64) (int64, *xerror.XerrorT)
	UpdateTag(request tag.UpdateTagRequest, userID int64) *xerror.XerrorT
	DeleteTag(tagID, userID int64) *xerror.XerrorT
	MyTags(userID int64) ([]tag.TagType, *xerror.XerrorT)
	TagInUse(tagID int64) (bool, *xerror.XerrorT)
}

type tagService struct{}

var TagService TagServiceI = &tagService{}

func (t *tagService) CreateTag(req tag.CreateTagRequest, userID int64) (int64, *xerror.XerrorT) {
	// Check if tag already exists
	var id int64
	err := clients.ClientOrm.Table("tags").
		Select("id").
		WhereRaw("tag_name = ? AND (user_id = ? OR user_id IS NULL)", req.TagName, userID).
		First(&id)

	if err != nil && err != sql.ErrNoRows {
		return 0, xerror.NewInternalError("error creating tag")
	}
	if id != 0 {
		return 0, xerror.NewBadRequest("tag already exists")
	}

	id, err = clients.ClientOrm.Table("tags").
		Insert(myorm.H{
			"tag_name":  req.TagName,
			"user_id":   userID,
			"hex_color": req.HexColor,
		})

	if err != nil {
		utils.ErrorLogger(err)
		return 0, xerror.NewInternalError("unable to create tag")
	}

	return id, nil
}

func (t *tagService) UpdateTag(req tag.UpdateTagRequest, userID int64) *xerror.XerrorT {
	err := clients.ClientOrm.Table("tags").
		Where("user_id", "=", userID).
		Where("id", "=", req.TagID).
		Update(myorm.H{
			"tag_name":  req.TagName,
			"hex_color": req.HexColor,
		})

	if err != nil {
		utils.ErrorLogger(err)
		return xerror.NewInternalError("unable to update tag")
	}

	return nil
}

func (t *tagService) DeleteTag(tagID, userID int64) *xerror.XerrorT {
	// Check that tag isn't is use already
	var id int64
	err := clients.ClientOrm.Table("entries").
		Select("id").
		Where("user_id", "=", userID).
		Where("tag_id", "=", tagID).
		First(&id)

	if err != nil && err != sql.ErrNoRows {
		utils.ErrorLogger(err)
		return xerror.NewInternalError("server error")
	}

	if id != 0 {
		return xerror.NewBadRequest("cannot delete tag that is already in use")
	}

	err = clients.ClientOrm.Table("tags").
		Where("user_id", "=", userID).
		Where("id", "=", tagID).
		Delete()

	if err != nil {
		utils.ErrorLogger(err)
		return xerror.NewInternalError("error deleting tag")
	}

	return nil
}

func (t *tagService) MyTags(userID int64) ([]tag.TagType, *xerror.XerrorT) {
	var result []tag.TagType

	err := clients.ClientOrm.Table("tags").
		Select("id", "tag_name", "user_id", "hex_color").
		WhereRaw("user_id = ? OR user_id IS NULL", userID).
		Get(&result)

	if err != nil {
		utils.ErrorLogger(err)
		return nil, xerror.NewInternalError("error getting tags")
	}

	return result, nil
}

func (t *tagService) TagInUse(id int64) (bool, *xerror.XerrorT) {
	var idInUse int

	err := clients.ClientOrm.Table("entries").
		Select("id").
		Where("tag_id", "=", id).
		First(&idInUse)

	if err != nil && err != sql.ErrNoRows {
		return true, xerror.NewInternalError("server error")
	}

	if err == sql.ErrNoRows || idInUse == 0 {
		return false, nil
	}

	return true, nil
}

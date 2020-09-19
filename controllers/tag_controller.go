package controllers

import (
	"api/domains/tag"
	"api/services"
	"api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MyTags(c *gin.Context) {
	user := services.GetAuthUser(c)

	res, err := services.TagService.MyTags(user.ID)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func CreateTag(c *gin.Context) {
	var req tag.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "true",
			"message": "invalid request body",
		})
		return
	}

	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, validationErr)
		return
	}

	user := services.GetAuthUser(c)

	id, err := services.TagService.CreateTag(req, user.ID)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, tag.TagType{
		TagID:   id,
		TagName: req.TagName,
	})
}

func UpdateTag(c *gin.Context) {
	var req tag.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "true",
			"message": "invalid request body",
		})
		return
	}

	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, validationErr)
		return
	}

	user := services.GetAuthUser(c)

	err := services.TagService.UpdateTag(req, user.ID)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, tag.TagType{
		TagID:   req.TagID,
		TagName: req.TagName,
	})
}

func DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "invalid tag id",
		})
		return
	}

	user := services.GetAuthUser(c)

	e := services.TagService.DeleteTag(int64(idInt), user.ID)

	if e != nil {
		c.JSON(e.StatusCode, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

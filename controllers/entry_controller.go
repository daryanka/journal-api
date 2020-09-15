package controllers

import (
	"api/domains/entry"
	"api/services"
	"api/utils"
	"api/utils/xerror"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateEntry(c *gin.Context) {
	var req entry.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("here", err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "true",
			"message": "invalid request body",
		})
		return
	}

	user := services.GetAuthUser(c)
	req.UserID = user.ID

	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, validationErr)
		return
	}

	// Parse Day
	t, e := utils.ParseDate(req.DayString)
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"day": "invalid time structure",
		})
		return
	}
	req.Day = t

	id, err := services.EntryService.CreateEntry(req)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"id":    id,
	})
}

func DeleteEntry(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "invalid entry",
		})
		return
	}

	user := services.GetAuthUser(c)

	e := services.EntryService.DeleteEntry(int64(id), user.ID)

	if e != nil {
		c.JSON(e.StatusCode, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func UpdateEntry(c *gin.Context) {
	var req entry.UpdateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "true",
			"message": "invalid request body",
		})
		return
	}

	user := services.GetAuthUser(c)

	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, validationErr)
		return
	}

	err := services.EntryService.UpdateEntry(user.ID, req)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
	})
}

func ViewDayEntries(c *gin.Context) {
	var req entry.ViewDayRequest
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

	// Parse Day
	_, e := utils.ParseDate(req.Day)
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"day": "invalid time structure",
		})
		return
	}

	user := services.GetAuthUser(c)

	result, err := services.EntryService.ViewDayEntries(req.Day, user.ID)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, utils.AdjustDay(result))
}

func ViewRangeEntries(c *gin.Context) {
	var req entry.EntriesRangeRequest
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

	// Parse Day
	from, e := utils.ParseDate(req.From)
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"from": "invalid time structure",
		})
		return
	}

	to, e := utils.ParseDate(req.To)
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"to": "invalid time structure",
		})
		return
	}

	// Check that from is before to
	if to.Sub(from) < 0 {
		e := xerror.NewBadRequest("Invalid date range", "RESET_RANGE")
		c.JSON(e.StatusCode, e)
		return
	}

	user := services.GetAuthUser(c)

	result, err := services.EntryService.ViewRangeEntries(req.From, req.To, user.ID)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, utils.AdjustDay(result))
}

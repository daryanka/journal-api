package controllers

import (
	"api/domains/entry"
	"api/services"
	"api/utils"
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
	fmt.Println(t)
	if e != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"day": "invalid time structure",
			"t":   req,
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

}

func ViewWeekEntries(c *gin.Context) {

}

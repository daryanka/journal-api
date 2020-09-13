package controllers

import (
	"api/domains/auth"
	"api/services"
	"api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(c *gin.Context)  {
	var req auth.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "true",
			"message": "invalid request body",
		})
		return
	}

	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		c.JSON(http.StatusUnprocessableEntity, validationErr)
		return
	}

	token, err := services.AuthService.CreateUser(&req)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, token)
}

func Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "true",
			"message": "invalid request body",
		})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	token, err := services.AuthService.Login(&req)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, token)
}
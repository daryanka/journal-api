package routes

import (
	"api/controllers"
	"github.com/gin-gonic/gin"
)

func StartRouting() {
	r := gin.Default()
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.CreateUser)
		authRoutes.POST("/login", controllers.Login)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err.Error())
	}
}
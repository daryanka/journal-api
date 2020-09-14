package routes

import (
	"api/controllers"
	"github.com/gin-gonic/gin"
)

func StartRouting() {
	r := gin.Default()
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}

	entriesRoutes := r.Group("entries")
	{
		entriesRoutes.POST("/", controllers.CreateEntry)
		entriesRoutes.PUT("/", controllers.UpdateEntry)
		entriesRoutes.DELETE("/:id", controllers.DeleteEntry)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

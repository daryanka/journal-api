package routes

import (
	"api/controllers"
	"api/middleware"
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
		entriesRoutes.POST("/day", middleware.ValidateAuthToken(), controllers.ViewDayEntries)
		entriesRoutes.POST("/range", middleware.ValidateAuthToken(), controllers.ViewRangeEntries)
		entriesRoutes.POST("/", middleware.ValidateAuthToken(), controllers.CreateEntry)
		entriesRoutes.PUT("/", middleware.ValidateAuthToken(), controllers.UpdateEntry)
		entriesRoutes.DELETE("/:id", middleware.ValidateAuthToken(), controllers.DeleteEntry)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

package routes

import (
	"api/controllers"
	"api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func StartRouting() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000", "https://localhost:5000"},
		AllowOriginFunc:  nil,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, Authorization, host, referer"},
		AllowCredentials: true,
		ExposeHeaders:    nil,
		MaxAge:           12 * time.Hour,
	}))

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.GET("/me", middleware.ValidateAuthToken(), controllers.Me)
	}

	entriesRoutes := r.Group("/entries", middleware.ValidateAuthToken())
	{
		entriesRoutes.POST("/day", controllers.ViewDayEntries)
		entriesRoutes.POST("/range", controllers.ViewRangeEntries)
		entriesRoutes.POST("/", controllers.CreateEntry)
		entriesRoutes.PUT("/", controllers.UpdateEntry)
		entriesRoutes.DELETE("/:id", controllers.DeleteEntry)
	}

	tagRoutes := r.Group("/tag", middleware.ValidateAuthToken())
	{
		tagRoutes.GET("/", controllers.MyTags)
		tagRoutes.POST("/", controllers.CreateTag)
		tagRoutes.PUT("/", controllers.UpdateTag)
		tagRoutes.DELETE("/:id", controllers.DeleteTag)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

package routes

import (
	"api/clients"
	"api/controllers"
	"api/middleware"
	"api/utils/xerror"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

func StartRouting() {
	r := gin.Default()

	corsConf := cors.Config{
		AllowOrigins:     []string{"http://localhost:5000"},
		AllowOriginFunc:  nil,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, Authorization, host, referer"},
		AllowCredentials: true,
		ExposeHeaders:    nil,
		MaxAge:           12 * time.Hour,
	}

	if strings.ToLower(os.Getenv("app")) == "prod" {
		corsConf.AllowOrigins = []string{"https://journal.daryanamin.co.uk"}
	}

	r.Use(cors.New(corsConf))

	r.NoRoute(func(c *gin.Context) {
		e := xerror.XerrorT{
			StatusCode: http.StatusNotFound,
			Message:    "Invalid route",
			Error:      true,
			Type:       "",
		}
		c.JSON(e.StatusCode, e)
	})

	r.GET("/health", func(c *gin.Context) {
		db, _ := clients.ClientOrm.GetDBConnection("defaultCon")
		dbOK := "OK"
		if err := db.Ping(); err != nil {
			dbOK = "ERROR"
		}

		c.JSON(200, gin.H{
			"status":   "OK",
			"database": dbOK,
		})
	})

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
		tagRoutes.GET("/in-use/:id", controllers.TagInUse)
	}

	// Default port 5000
	port := "5000"
	if os.Getenv("port") != "" {
		port = os.Getenv("port")
	}

	if err := r.Run(fmt.Sprintf(":%v", port)); err != nil {
		panic(err.Error())
	}
}

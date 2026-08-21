package router

import (
	"personal-erp-backend/internal/handler"
	"personal-erp-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(health *handler.HealthHandler, categories *handler.CategoryHandler, activities *handler.ActivityHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", health.Check)

		api.GET("/categories", categories.List)
		api.POST("/categories", categories.Create)
		api.PUT("/categories/:id", categories.Update)
		api.DELETE("/categories/:id", categories.Delete)

		api.GET("/activities", activities.List)
		api.GET("/activities/stats", activities.Stats)
		api.POST("/activities", activities.Create)
		api.GET("/activities/:id", activities.Get)
		api.PUT("/activities/:id", activities.Update)
		api.DELETE("/activities/:id", activities.Delete)
	}
	return r
}
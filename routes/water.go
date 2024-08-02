package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func WaterRoute(router *gin.Engine) {
	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/water", controller.CreateWater)
	authRoutes.GET("/water", controller.ListWater)
	authRoutes.DELETE("/water/:id", controller.DeleteWater)
	authRoutes.PATCH("/water/:id", controller.UpdateWater)
}

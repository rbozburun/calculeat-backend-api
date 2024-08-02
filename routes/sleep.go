package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func SleepRoute(router *gin.Engine) {
	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/sleep", controller.CreateSleep)
	authRoutes.GET("/sleep", controller.ListSleep)
	authRoutes.DELETE("/sleep/:id", controller.DeleteSleep)
	authRoutes.PATCH("/sleep/:id", controller.UpdateSleep)
}

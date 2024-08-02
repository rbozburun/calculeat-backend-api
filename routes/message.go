package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func MessageRoute(router *gin.Engine) {
	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/message", controller.CreateMessage)
	authRoutes.GET("/message", controller.ListMessages)
	authRoutes.DELETE("/message/:id", controller.DeleteMessage)
	authRoutes.PATCH("/message/:id", controller.UpdateMessage)

}

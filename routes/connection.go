package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func ConnectionRoute(router *gin.Engine) {
	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/connection", controller.CreateConnection)
	authRoutes.GET("/connection", controller.ListConnections)
	authRoutes.DELETE("/connection/:id", controller.DeleteConnection)
	authRoutes.PATCH("/connection/:id", controller.UpdateConnection)

}

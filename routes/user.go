package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	logger.Log.Debugln("UserRoute called.")

	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/user", controller.CreateUser)
	authRoutes.GET("/user", controller.ListUsers)
	authRoutes.DELETE("/user/:id", controller.DeleteUser)
	authRoutes.PATCH("/user/:id", controller.UpdateUser)

}

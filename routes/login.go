package routes

import (
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func LoginRoute(router *gin.Engine) {
	v1 := router.Group("/api/1.0")
	v1.POST("/login", controller.Login)
}

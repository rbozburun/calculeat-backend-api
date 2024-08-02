package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func MeetingRoute(router *gin.Engine) {
	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/meeting", controller.CreateMeeting)
	authRoutes.GET("/meeting", controller.ListMeetings)
	authRoutes.DELETE("/meeting/:id", controller.DeleteMeeting)
	authRoutes.PATCH("/meeting/:id", controller.UpdateMeeting)

}

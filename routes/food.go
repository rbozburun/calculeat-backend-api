package routes

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/controller"
	"github.com/gin-gonic/gin"
)

func FoodRoute(router *gin.Engine) {
	authRoutes := router.Group("/api/1.0")
	authRoutes.Use(config.AuthMiddleware())

	authRoutes.POST("/food", controller.CreateFood)
	authRoutes.GET("/food", controller.ListFoods)
	authRoutes.DELETE("/food/:id", controller.DeleteFood)
	authRoutes.PATCH("/food/:id", controller.UpdateFood)

}

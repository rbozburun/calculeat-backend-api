package main

import (
	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/middlewares"
	"github.com/calculeat/main_rest_api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Connect()
	logger.Log.Debugln("API Main function started...")

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())
	router.Use(middlewares.CORSMiddleware())

	routes.UserRoute(router)
	logger.Log.Debugln("UserRoute initiliazed.")

	routes.ConnectionRoute(router)
	logger.Log.Debugln("ConnectionRoute initiliazed.")

	routes.FoodRoute(router)
	logger.Log.Debugln("FoodRoute initiliazed.")

	routes.MeetingRoute(router)
	logger.Log.Debugln("MeetingRoute initiliazed.")

	routes.MessageRoute(router)
	logger.Log.Debugln("MessageRoute initiliazed.")

	routes.WaterRoute(router)
	logger.Log.Debugln("WaterRoute initiliazed.")

	routes.SleepRoute(router)
	logger.Log.Debugln("SleepRoute initiliazed.")

	routes.LoginRoute(router)
	logger.Log.Debugln("LoginRoute initiliazed.")

	routes.PredictRoute(router)
	logger.Log.Debugln("PredictRoute initiliazed.")

	router.Run("0.0.0.0:7854")
}

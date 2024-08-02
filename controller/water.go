package controller

import (
	"net/http"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// POST /api/1.0/water - Creates a new water object in DB
func CreateWater(ctx *gin.Context) {
	var water models.WaterObject
	ctx.ShouldBindJSON(&water)
	water.UserID = config.CURRENT_USER.ID

	if err := config.DB.Create(&water).Error; err != nil {
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		logger.Log.Error(logger.BindError)
		return
	} else {
		// Create water object succesfuly and respond with the object
		ctx.JSON(http.StatusCreated, &water)
		logger.Log.Infoln(logger.ObjectCreated)
	}
}

// GET /api/1.0/water - Lists all water objects in database
func ListWater(ctx *gin.Context) {
	query_water_model := models.WaterObject{}
	query_water_model.UserID = config.CURRENT_USER.ID
	water_objects := []models.WaterObject{}

	// Check for the date query is used or not
	if ctx.Query("date") != "" {
		// Query to db with the given date
		date := ctx.Query("date")
		logger.Log.Debugf("Date sent in request: %v, Current user id: %v", date, query_water_model.UserID)
		config.DB.Where("updated_at::date = ? AND user_id = ?", date, query_water_model.UserID).Find(&water_objects)
		ctx.JSON(http.StatusOK, &water_objects)
		return
	}

	// ** Enabled for development purposes **
	config.DB.Where(&query_water_model).Find(&water_objects)
	ctx.JSON(http.StatusOK, &water_objects)
}

// DELETE /api/1.0/water/:id - Deletes the specified water object with ID
func DeleteWater(ctx *gin.Context) {
	var water_object models.WaterObject
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&water_object).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	if config.CURRENT_USER.ID != water_object.UserID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.UnauthorizedAccessOperation)
		return
	}

	config.DB.Delete(&water_object, id)
	logger.Log.Infof("Water object with %v ID is deleted.", id)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})

}

// PATCH /api/1.0/water/:id - Updates the specified water object's value according to the body data
func UpdateWater(ctx *gin.Context) {
	var water_object models.WaterObject
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&water_object).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	// Validate input
	var input models.UpdateWaterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(err.Error())
		return
	}

	isUpdated := false
	water_object.UserID = config.CURRENT_USER.ID

	if input.Count != 0 && input.Count != water_object.Count {
		config.DB.Model(&water_object).Updates(map[string]interface{}{"count": input.Count})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": &water_object, "msg": msg})
}

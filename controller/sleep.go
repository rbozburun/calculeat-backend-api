package controller

import (
	"net/http"
	"strconv"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// POST /api/1.0/sleep - Creates a new sleep object to the database.
func CreateSleep(ctx *gin.Context) {
	var sleep models.SleepObject
	ctx.ShouldBindJSON(&sleep)
	sleep.UserID = config.CURRENT_USER.ID

	if err := config.DB.Create(&sleep).Error; err != nil {
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.InstanceCreationErrorSQL)
		return
	} else {
		// Create sleep object succesfuly and respond with the object
		ctx.JSON(http.StatusCreated, &sleep)
		return
	}
}

// GET /api/1.0/sleep - Lists all sleep objects in database
func ListSleep(ctx *gin.Context) {
	//Bind URL query parameters to a model
	query_sleep_model := models.SleepObject{}
	query_sleep_model.UserID = config.CURRENT_USER.ID
	sleep_objects := []models.SleepObject{}

	query_sleep_model.ID, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	query_sleep_model.StartTime = ctx.Query("start_time")
	query_sleep_model.EndTime = ctx.Query("end_time")

	//Query with the model
	config.DB.Where(&query_sleep_model).Find(&sleep_objects)
	ctx.JSON(200, &sleep_objects)
}

// DELETE /api/1.0/sleep/:id - Deletes the specified sleep object with ID
func DeleteSleep(ctx *gin.Context) {
	var sleep_object models.SleepObject
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&sleep_object).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		return
	}

	if config.CURRENT_USER.ID != sleep_object.UserID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.UnauthorizedAccessOperation)
		return
	}

	config.DB.Delete(&sleep_object, id)
	logger.Log.Infof("Sleep object with %v ID is deleted.", id)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})
}

// PATCH /api/1.0/sleep/:id - Updates the specified sleep object's value according to the body data
func UpdateSleep(ctx *gin.Context) {
	var sleep_object models.SleepObject
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&sleep_object).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		return
	}

	if config.CURRENT_USER.ID != sleep_object.UserID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.UnauthorizedAccessOperation)
		return
	}

	// Validate input
	var input models.UpdateSleepInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	isUpdated := false
	if input.StartTime != "" && input.StartTime != sleep_object.StartTime {
		config.DB.Model(&sleep_object).Updates(map[string]interface{}{"start_time": input.StartTime})
		isUpdated = true
	}

	if input.EndTime != "" && input.EndTime != sleep_object.EndTime {
		config.DB.Model(&sleep_object).Updates(map[string]interface{}{"end_time": input.EndTime})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": &sleep_object, "msg": msg})
}

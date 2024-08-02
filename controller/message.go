package controller

import (
	"net/http"
	"strconv"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// POST /api/1.0/message - Creates a new message in DB
func CreateMessage(ctx *gin.Context) {
	var message models.Message
	err_bind := ctx.ShouldBindJSON(&message)
	if err_bind != nil {
		logger.Log.Error(logger.BindError + " ->" + err_bind.Error())
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err_bind.Error()})
		return
	}

	message.SenderID = config.CURRENT_USER.ID
	if err := config.DB.Create(&message).Error; err != nil {
		// Throw error
		logger.Log.Error(logger.InstanceCreationErrorSQL)
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	} else {
		// Create message succesfuly and respond with message data
		ctx.JSON(http.StatusCreated, &message)
		logger.Log.Debugln(logger.ObjectCreated)
	}

}

// GET /api/1.0/message - Lists all of the messages in DB
func ListMessages(ctx *gin.Context) {
	//Bind URL query parameters to a model
	query_message_model := models.Message{}
	query_message_model.SenderID = config.CURRENT_USER.ID
	messages := []models.Message{}

	query_message_model.ID, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	query_message_model.RecieverID, _ = strconv.ParseUint(ctx.Query("reciever_id"), 10, 64)
	query_message_model.MessageText = ctx.Query("message_text")

	//Query with the model
	config.DB.Where(&query_message_model).Find(&messages)
	ctx.JSON(200, &messages)

}

// DELETE /message/:id - Deletes the given message from DB
func DeleteMessage(ctx *gin.Context) {
	var message models.Message
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&message).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	if config.CURRENT_USER.ID != message.SenderID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.UnauthorizedAccessOperation)
		return
	}

	config.DB.Delete(&message, id)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})
}

// PACTH /message/:id - Updates the message according to the data
func UpdateMessage(ctx *gin.Context) {
	var message models.Message
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&message).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Errorln(logger.RecordCouldNotFound)
		return
	}

	if config.CURRENT_USER.ID != message.SenderID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.UnauthorizedAccessOperation)
		return
	}

	// Validate input
	var input models.UpdateMessageInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Error(logger.BindError + " ->" + err.Error())
		return
	}

	isUpdated := false

	if input.MessageText != message.MessageText {
		config.DB.Model(&message).Updates(map[string]interface{}{"message_text": input.MessageText})
		isUpdated = true
	}

	if input.RecieverID != 0 && input.RecieverID != message.RecieverID {
		config.DB.Model(&message).Updates(map[string]interface{}{"reciever_id": input.RecieverID})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": &message, "msg": msg})

}

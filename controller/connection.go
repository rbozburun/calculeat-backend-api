package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/helpers"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// TO DO:
// [x] Query for ID (Bu şekilde oluşturulan modelde ID yok!)
// [x] is_accepted false oldugunda hata alıyoruz!
// [x] create endpointi authenticated olmalı!
// [] Update endpointine gelen istekte update edilen bir attr yoksa sadece `There is nothing to update, please change an attribute of the object.` mesajı basılacka.

// POST /api/1.0/connection - Creates a new connection in DB
func CreateConnection(ctx *gin.Context) {
	var connection models.Connection
	connectionInput := models.CreateConnectionInput{}
	connectionInput.TargetUserID, _ = strconv.ParseUint(ctx.Query("target_user_id"), 10, 64)

	// Parse the user ids from target_user_id
	dieticianUserIDParsed, clientUserIDParsed := helpers.GetDieticianAndClientUserIDs(connectionInput.TargetUserID)
	connection.DieticianUserID = dieticianUserIDParsed
	connection.ClientUserID = clientUserIDParsed

	// Check the current user has access to create the connection
	if config.CURRENT_USER.RoleID == 1 {
		if connection.DieticianUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	} else if config.CURRENT_USER.RoleID == 2 {
		if connection.ClientUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	}

	// Default FALSE
	connection.IsAccepted = "f"

	if err := config.DB.Create(&connection).Error; err != nil {
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		logger.Log.Error(logger.InstanceCreationErrorSQL)
		return
	} else {
		// Create connection succesfuly and respond with access token (non-hashed)
		ctx.JSON(http.StatusCreated, &connection)
		logger.Log.Debugln(logger.ObjectCreated)
	}
}

// GET /api/1.0/connection - Lists all of the connections in DB
func ListConnections(ctx *gin.Context) {
	//Bind URL query parameters to a model
	query_connection_model := models.Connection{}
	connections := []models.Connection{}

	query_connection_model.ID, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	// The client should see the other user's connections
	//query_connection_model.DieticianUserID, _ = strconv.ParseUint(ctx.Query("dietician_user_id"), 10, 64)
	//query_connection_model.ClientUserID, _ = strconv.ParseUint(ctx.Query("client_user_id"), 10, 64)

	// Parse the user ids from target_user_id
	dieticianUserIDParsed, clientUserIDParsed := helpers.GetDieticianAndClientUserIDs(config.CURRENT_USER.ID)
	query_connection_model.DieticianUserID = dieticianUserIDParsed
	query_connection_model.ClientUserID = clientUserIDParsed

	// Check the current user has access to list the connections
	if config.CURRENT_USER.RoleID == 1 {
		if query_connection_model.DieticianUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	} else if config.CURRENT_USER.RoleID == 2 {
		if query_connection_model.ClientUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	}
	query_connection_model.IsAccepted = strings.ToLower(ctx.Query("is_accepted"))

	//Query with the model
	config.DB.Where(&query_connection_model).Find(&connections)
	ctx.JSON(200, &connections)
}

// DELETE /connection/:id - Deletes the given connection from DB
func DeleteConnection(ctx *gin.Context) {
	var connection models.Connection
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&connection).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	// Check the current user has access to delete the connection
	if config.CURRENT_USER.RoleID == 1 {
		if connection.DieticianUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	} else if config.CURRENT_USER.RoleID == 2 {
		if connection.ClientUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	}

	config.DB.Delete(&connection, id)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})
}

// NOT:
///
/// PATCH requestleri target_user id alıp client/dieticien ayrımını kendisi yapsın.
// Bunun için bir helper fonksiyon yazılabilir.

// PATCH /connection/:id - Updates the connection according to the data
func UpdateConnection(ctx *gin.Context) {
	var connection models.Connection
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&connection).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Errorln(logger.RecordCouldNotFound)
		return
	}

	isUpdated := false
	// Validate input, client sends only {target_user_id, is_accepted}, need to parse client & dietician
	var updateConnectionInput models.CreateConnectionInput

	if err := ctx.ShouldBindJSON(&updateConnectionInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Error(logger.BindError + " ->" + err.Error())
		return
	}

	dieticianUserIDParsed, clientUserIDParsed := helpers.GetDieticianAndClientUserIDs(updateConnectionInput.TargetUserID)

	// Bind the user ids after parsing
	var input models.UpdateConnectionInput
	input.DieticianUserID = dieticianUserIDParsed
	input.ClientUserID = clientUserIDParsed
	input.IsAccepted = strings.ToLower(updateConnectionInput.IsAccepted)
	// If there is no non-zero value, update the model with given request body
	if input.IsAccepted == "t" {
		config.DB.Model(&connection).Updates(map[string]interface{}{"is_accepted": "t"})
		isUpdated = true
	} else if input.IsAccepted == "f" {
		config.DB.Model(&connection).Updates(map[string]interface{}{"is_accepted": "f"})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": &connection, "msg": msg})
}

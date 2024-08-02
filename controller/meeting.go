package controller

import (
	"net/http"
	"strconv"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/helpers"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// POST /api/1.0/meeting - Creates a new meeting in DB
func CreateMeeting(ctx *gin.Context) {
	var meeting models.Meeting
	var meetingInput models.CreateMeetingInput
	err_bind := ctx.ShouldBindJSON(&meetingInput)
	if err_bind != nil {
		logger.Log.Error(logger.BindError + " ->" + err_bind.Error())
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err_bind.Error()})
		return
	}

	// Current user is a dietician
	if config.CURRENT_USER.RoleID == 1 {
		meeting.DieticianUserID = config.CURRENT_USER.ID
		meeting.ClientUserID = meetingInput.TargetUserID
	} else if config.CURRENT_USER.RoleID == 2 {
		// Current user is a client
		meeting.ClientUserID = config.CURRENT_USER.ID
		meeting.DieticianUserID = meetingInput.TargetUserID
	}

	meeting.MeetLink = meetingInput.MeetLink

	if err := config.DB.Create(&meeting).Error; err != nil {
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		logger.Log.Error(logger.InstanceCreationErrorSQL)
		return
	} else {
		// Create meeting succesfuly
		ctx.JSON(http.StatusCreated, &meeting)
		logger.Log.Debugln(logger.ObjectCreated)
	}
}

// GET /api/1.0/meeting - Lists all of the meetings in DB
func ListMeetings(ctx *gin.Context) {
	//Bind URL query parameters to a model
	meetingInput := models.CreateMeetingInput{}
	var query_meeting_model models.Meeting
	meetings := []models.Meeting{}

	meetingInput.TargetUserID, _ = strconv.ParseUint(ctx.Query("target_user_id"), 10, 64)

	// Current user is a dietician
	if config.CURRENT_USER.RoleID == 1 {
		query_meeting_model.DieticianUserID = config.CURRENT_USER.ID
		query_meeting_model.ClientUserID = meetingInput.TargetUserID
	} else if config.CURRENT_USER.RoleID == 2 {
		// Current user is a client
		query_meeting_model.ClientUserID = config.CURRENT_USER.ID
		query_meeting_model.DieticianUserID = meetingInput.TargetUserID
	}

	query_meeting_model.MeetLink = ctx.Query("meet_link")

	//Query with the model
	config.DB.Where(&query_meeting_model).Find(&meetings)
	ctx.JSON(200, &meetings)
}

// DELETE /meeting/:id - Deletes the given meeting from DB
func DeleteMeeting(ctx *gin.Context) {
	var meeting models.Meeting
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&meeting).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	// Check the current user has access to delete the meeting
	if config.CURRENT_USER.RoleID == 1 {
		if meeting.DieticianUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	} else if config.CURRENT_USER.RoleID == 2 {
		if meeting.ClientUserID != config.CURRENT_USER.ID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
			logger.Log.Error(logger.UnauthorizedAccessOperation)
			return
		}
	}

	config.DB.Delete(&meeting, id)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})
}

// PACTH /meeting/:id - Updates the meeting according to the data
func UpdateMeeting(ctx *gin.Context) {
	var meeting models.Meeting
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&meeting).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Errorln(logger.RecordCouldNotFound)
		return
	}

	// Validate input, client sends only {target_user_id, is_accepted}, need to parse client & dietician
	var updateMeetingInput models.CreateMeetingInput

	if err := ctx.ShouldBindJSON(&updateMeetingInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Error(logger.BindError + " ->" + err.Error())
		return
	}

	dieticianUserIDParsed, clientUserIDParsed := helpers.GetDieticianAndClientUserIDs(updateMeetingInput.TargetUserID)

	// Validate input
	var input models.UpdateMeetingInput
	input.DieticianUserID = dieticianUserIDParsed
	input.ClientUserID = clientUserIDParsed
	isUpdated := false

	if input.MeetLink != "" && input.MeetLink != meeting.MeetLink {
		config.DB.Model(&meeting).Updates(map[string]interface{}{"meet_link": input.MeetLink})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": &meeting, "msg": msg})
}

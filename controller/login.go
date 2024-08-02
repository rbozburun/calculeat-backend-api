package controller

import (
	"context"
	"net/http"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/helpers"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// POST /api/1.0/login
func Login(ctx *gin.Context) {
	logger.Log.Debugln("Login controller started.")
	// Get credentials from context
	var input models.LoginUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user
	user, err := helpers.FindUserByEmail(input.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Errorf("User cannot find -> " + err.Error())
		return
	}

	// Verify credentials from Firebase & DB
	// Compare supplied password and the user's password in DB.
	is_match, err := helpers.ComparePasswords(input.Password, user.Password)
	if is_match {
		// Get Custom Token from firebase
		logger.Log.Debug("Getting Custom Token from Firebase...")
		f_user, err := config.FIRE_AUTH_CLIENT.GetUserByEmail(context.Background(), user.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			logger.Log.Errorf("Error fetching user: %v\n", err)
			return
		}

		// Get the user's UID
		uid := f_user.UID

		// Generate a custom JWT token if needed
		customToken, err := config.FIRE_AUTH_CLIENT.CustomToken(context.Background(), uid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			logger.Log.Errorf("Error generating custom token: %v\n", err)
			return
		}
		logger.Log.Debugf("Custom Token fetched! UID: %v - CustomToken: %v", uid, customToken)

		// Login with Custom Token
		logger.Log.Debug("Signing in with Custom Token...")
		idToken, err := helpers.SignInWithCustomToken(customToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			logger.Log.Errorf("Error while signing with custom token: %v\n", err)
			return
		}

		logger.Log.Debug("ID Token fetched and the user signed in.")
		ctx.JSON(http.StatusOK, gin.H{"access_token": idToken, "user": &user})

	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		logger.Log.Info("There is something wrong, possible the password does not match. -> " + err.Error())
		return
	}

}

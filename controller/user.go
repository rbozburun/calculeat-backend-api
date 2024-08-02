package controller

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/helpers"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

var FIREBASE_CREDENTIALS_FILE string

// To Do
// [] Update endpointine gelen istekte update edilen bir attr yoksa sadece `There is nothing to update, please change an attribute of the object.` mesajı basılacka.

// GET /api/1.0/user - Lists all of the users in DB
func ListUsers(ctx *gin.Context) {
	logger.Log.Debugln("ListUsers endpoint called.")
	//Get the current user's model
	var current_user models.User
	logger.Log.Debugf("Current user email: %v", config.FIRE_USER.Email)
	logger.Log.Debugf("Current user RoleID: %v", current_user.RoleID)

	// Fetch the current user from DB
	config.DB.First(&current_user, "email = ?", config.FIRE_USER.Email)
	logger.Log.Debugf("Current user (Fetched from DB) email: %v", current_user.Email)

	//Bind URL query parameters to a model
	query_user_model := models.User{}
	users := []models.User{}

	query_user_model.ID, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	//RoleID shouldn't get from the client, due to privacy restrictions.
	//query_user_model.RoleID, _ = strconv.Atoi(ctx.Query("role_id"))
	query_user_model.FirstName = ctx.Query("first_name")
	query_user_model.LastName = ctx.Query("last_name")
	query_user_model.Age, _ = strconv.Atoi(ctx.Query("age"))
	//query_user_model.Email = ctx.Query("email")
	query_user_model.Gender = ctx.Query("gender")
	//query_user_model.Born = ctx.Query("born")
	//query_user_model.Height, _ = strconv.Atoi(ctx.Query("height"))
	//query_user_model.Weight, _ = strconv.ParseFloat(ctx.Query("weight"), 32)
	//query_user_model.Goal = ctx.Query("goal")

	logger.Log.Debugf("Query parameters: %v", query_user_model)

	//Query with the model
	config.DB.Where(&query_user_model).Find(&users)
	users_redacted := []models.User{}

	for _, u := range users {
		if u.RoleID != current_user.RoleID {
			u.AccessToken = ""
			u.Password = ""
			u.Email = ""
			u.SignInProvider = ""
			u.Born = ""
			u.Height = 0
			u.Weight = 0
			u.Goal = ""
			u.Diseases = nil
			u.CreatedAt = time.Time{}
			u.UpdatedAt = time.Time{}
			u.DeletedOn = time.Time{}
			users_redacted = append(users_redacted, u)
		}
	}

	logger.Log.Debugln("Query operation successfull, if any user found it will be listed.")
	current_user = models.User{}
	ctx.JSON(200, &users_redacted)
}

// POST /api/1.0/user
func CreateUser(ctx *gin.Context) {
	logger.Log.Debugln("CreateUser controller called.")

	var new_user models.User
	firebase_user := config.FIRE_USER
	user_info := firebase_user.ProviderUserInfo[0]

	json_user_info, err := json.Marshal(user_info)
	if err != nil {
		logger.Log.Errorf("Cannot marshal the user info fetched from Firebase. Err: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot marshal the user info fetched from Firebase."})
	}

	var user_info_provider config.UserInfoProvider
	err = json.Unmarshal(json_user_info, &user_info_provider)
	if err != nil {
		logger.Log.Errorf("Cannot unmarshal the user info fetched from Firebase. Err: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot unmarshal the user info fetched from Firebase."})
	}

	provider := user_info_provider.ProviderId
	if provider == "password" || provider == "google.com" {
		// Create the user object in DB for signed user - [Password can be null for the google user.]
		err := ctx.ShouldBindJSON(&new_user)
		if err != nil {
			logger.Log.Errorf("Cannot bind POST body to user object. Err: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot bind the POST body"})
			return
		}
		if user_info_provider.Email != new_user.Email {
			logger.Log.Errorf("JWT email and POST body email does not match!")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "JWT email and POST body email does not match! Are you sure you are using your email in body?"})
			return
		}
		new_user.SignInProvider = provider

		// If the provider is 'password' & the password is not null, set password as SHA256 in DB.
		if new_user.SignInProvider == "password" {
			if new_user.Password != "" {
				pwdHash_byte := sha256.Sum256([]byte(new_user.Password))
				pwdHash_input := fmt.Sprintf("%x", pwdHash_byte[:])
				new_user.Password = pwdHash_input
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Empty password for email/password login."})
				logger.Log.Error("Empty password for email/password login.")
				return
			}
		}
		logger.Log.Debugf("User creation on db... User: %v", &new_user)

		if err := config.DB.Create(&new_user).Error; err != nil {
			// Throw error
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			logger.Log.Error(logger.InstanceCreationErrorSQL + err.Error())
			return
		} else {
			// The user created successfully. Fetch the user from DB then create water for it.

			// Create user object succesfuly and respond with the object
			new_user.Password = ""    // Remove the from response due to security reasons
			new_user.AccessToken = "" // Remove the from response due to security reasons

			logger.Log.Debugf("A new user created. User ID: %d, Email: %s, Provider: %s", new_user.ID, new_user.Email, new_user.SignInProvider)
			// When the user created, create sleep & water objects as empty for it
			created_water_obj, err := helpers.InitiliazeWaterObject(new_user.ID)
			if err != nil {
				logger.Log.Debugf("Could not initiliaze water object for -> User ID: %d, Email: %s, Error: %v", new_user.ID, new_user.Email, err)
			} else {
				logger.Log.Debugf("The water object initiliazed for -> User ID: %d, Email: %s", new_user.ID, new_user.Email)
			}

			// Set the initial water object id (init_water_id) for the created user for client-side response
			new_user.InitWaterId = int(created_water_obj.ID)

			ctx.JSON(http.StatusCreated, &new_user)

			return
		}

	} else {
		logger.Log.Errorf("Cannot determine provider. None of the'password or google.com' Err: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot determine provider. None of the'password or google.com'"})
		return
	}

}

// DELETE /api/1.0/user/:id
func DeleteUser(ctx *gin.Context) {
	var user models.User
	// A user can only delete itself
	id := config.CURRENT_USER.ID

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	delete_item := config.DB.Delete(&user, id)
	if delete_item.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": delete_item.Error})
		return
	} else if delete_item.RowsAffected < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Row cannot be deleted because it doesn't exist"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})

}

// PACTH /api/1.0/user/:id
func UpdateUser(ctx *gin.Context) {
	var user models.User
	id := config.CURRENT_USER.ID

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	// Validate input
	var input models.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Error(logger.BindError)
		return
	}

	// Check for default null values
	isUpdated := false

	if input.FirstName != "" && input.FirstName != user.FirstName {
		config.DB.Model(&user).Updates(map[string]interface{}{"first_name": input.FirstName})
		isUpdated = true
	}

	if input.LastName != "" && input.LastName != user.LastName {
		config.DB.Model(&user).Updates(map[string]interface{}{"last_name": input.LastName})
		isUpdated = true
	}

	if input.ProfileImgFilepath != "" && input.ProfileImgFilepath != user.ProfileImgFilepath {
		config.DB.Model(&user).Updates(map[string]interface{}{"profile_img_filepath": input.ProfileImgFilepath})
		isUpdated = true
	}

	if input.Gender != "" && input.Gender != user.Gender {
		config.DB.Model(&user).Updates(map[string]interface{}{"gender": input.Gender})
		isUpdated = true
	}

	if input.Born != "" && input.Born != user.Born {
		config.DB.Model(&user).Updates(map[string]interface{}{"born": input.Born})
		isUpdated = true
	}

	if input.Age != 0 && input.Age != user.Age {
		config.DB.Model(&user).Updates(map[string]interface{}{"age": input.Age})
		isUpdated = true
	}

	if input.Email != "" && input.Email != user.Email {
		config.DB.Model(&user).Updates(map[string]interface{}{"email": input.Email})
		isUpdated = true
	}

	if input.Height != 0 && input.Height != user.Height {
		config.DB.Model(&user).Updates(map[string]interface{}{"height": input.Height})
		isUpdated = true
	}

	if input.Weight != 0 && input.Weight != user.Weight {
		config.DB.Model(&user).Updates(map[string]interface{}{"weight": input.Weight})
		isUpdated = true
	}

	if input.Goal != "" && input.Goal != user.Goal {
		config.DB.Model(&user).Updates(map[string]interface{}{"goal": input.Goal})
		isUpdated = true
	}

	if input.Diseases != nil && helpers.IsArraysEqual(input.Diseases, user.Diseases) == false {
		config.DB.Model(&user).Updates(map[string]interface{}{"diseases": input.Diseases})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": &user, "msg": msg})
	return
}

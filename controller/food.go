package controller

import (
	"net/http"
	"strconv"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

// POST /api/1.0/food - Creates a new food in DB
func CreateFood(ctx *gin.Context) {
	var food models.FoodObject
	food.UserID = config.CURRENT_USER.ID
	err_bind := ctx.ShouldBindJSON(&food)
	if err_bind != nil {
		logger.Log.Error(logger.BindError + " ->" + err_bind.Error())
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err_bind.Error()})
		return
	}

	if err := config.DB.Create(&food).Error; err != nil {
		// Throw error
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		logger.Log.Error(logger.InstanceCreationErrorSQL)
		return
	} else {
		// Create food succesfuly and respond with access token (non-hashed)
		ctx.JSON(http.StatusCreated, &food)
		logger.Log.Debugln(logger.ObjectCreated)
	}
}

// GET /api/1.0/food - Creates all foods in database
func ListFoods(ctx *gin.Context) {
	//Bind URL query parameters to a model
	query_food_model := models.FoodObject{}
	food_objects := []models.FoodObject{}

	query_food_model.ID, _ = strconv.ParseUint(ctx.Query("id"), 10, 64)
	query_food_model.Name = ctx.Query("name")
	query_food_model.Meal = ctx.Query("meal")
	query_food_model.UserID = config.CURRENT_USER.ID
	query_food_model.Calorie, _ = strconv.ParseUint(ctx.Query("calorie"), 10, 64)
	query_food_model.Carbonhydrate, _ = strconv.ParseFloat(ctx.Query("carbonhydrate"), 64)
	query_food_model.Protein, _ = strconv.ParseFloat(ctx.Query("protein"), 64)
	query_food_model.Fat, _ = strconv.ParseFloat(ctx.Query("fat"), 64)
	query_food_model.Gram, _ = strconv.ParseFloat(ctx.Query("gram"), 64)
	query_food_model.Sugar, _ = strconv.ParseFloat(ctx.Query("sugar"), 64)
	query_food_model.VitC, _ = strconv.ParseUint(ctx.Query("vit_c"), 10, 64)
	query_food_model.VitD, _ = strconv.ParseUint(ctx.Query("vit_d"), 10, 64)
	query_food_model.VitA, _ = strconv.ParseUint(ctx.Query("vit_a"), 10, 64)
	query_food_model.Calcium, _ = strconv.ParseFloat(ctx.Query("calcium"), 64)
	query_food_model.Magnesium, _ = strconv.ParseFloat(ctx.Query("magnesium"), 64)

	//Query with the model
	config.DB.Where(&query_food_model).Find(&food_objects)
	ctx.JSON(200, &food_objects)
}

// DELETE /api/1.0/food/:id - Deletes the specified food with ID
func DeleteFood(ctx *gin.Context) {
	var food models.FoodObject
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&food).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Error(logger.RecordCouldNotFound)
		return
	}

	if config.CURRENT_USER.ID != food.UserID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.SmthWrong})
		logger.Log.Error(logger.UnauthorizedAccessOperation)
		return
	}

	config.DB.Delete(&food, id)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Deleted"})
}

// PACTH /api/1.0/food/:id - Updates the specified food's value according to the body data
func UpdateFood(ctx *gin.Context) {
	var food models.FoodObject
	id := ctx.Param("id")

	// Get model if exist
	if err := config.DB.Where("id = ?", id).First(&food).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": logger.RecordCouldNotFound})
		logger.Log.Errorln(logger.RecordCouldNotFound)
		return
	}

	// Validate input
	var input models.UpdateFoodInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Log.Error(logger.BindError + " ->" + err.Error())
		return
	}
	input.UserID = config.CURRENT_USER.ID

	// Check for default null values
	// TO DO:
	// [] Fat should be checked for "Fat = UnsaturatedFat + SaturadFat"

	isUpdated := false
	if input.Name != "" && input.Name != food.Name {
		config.DB.Model(&food).Updates(map[string]interface{}{"name": input.Name})
		isUpdated = true
	}
	if input.UserID != 0 && input.UserID != food.UserID {
		config.DB.Model(&food).Updates(map[string]interface{}{"user_id": input.UserID})
		isUpdated = true
	}
	if input.ImgleftFilepath != "" && input.ImgleftFilepath != food.ImgleftFilepath {
		config.DB.Model(&food).Updates(map[string]interface{}{"imgleft_filepath": input.ImgleftFilepath})
		isUpdated = true
	}
	if input.ImgrightFilepath != "" && input.ImgrightFilepath != food.ImgrightFilepath {
		config.DB.Model(&food).Updates(map[string]interface{}{"imgright_filepath": input.ImgrightFilepath})
		isUpdated = true
	}
	if input.ImgtopFilepath != "" && input.ImgtopFilepath != food.ImgtopFilepath {
		config.DB.Model(&food).Updates(map[string]interface{}{"imgtop_filepath": input.ImgtopFilepath})
		isUpdated = true
	}
	if input.Calorie != 0 && input.Calorie != food.Calorie {
		config.DB.Model(&food).Updates(map[string]interface{}{"calorie": input.Calorie})
		isUpdated = true
	}
	if input.Protein != 0 && input.Protein != food.Protein {
		config.DB.Model(&food).Updates(map[string]interface{}{"protein": input.Protein})
		isUpdated = true
	}
	if input.Carbonhydrate != 0 && input.Carbonhydrate != food.Carbonhydrate {
		config.DB.Model(&food).Updates(map[string]interface{}{"carbonhydrate": input.Carbonhydrate})
		isUpdated = true
	}
	if input.Fat != 0 && input.Fat != food.Fat {
		config.DB.Model(&food).Updates(map[string]interface{}{"fat": input.Fat})
		isUpdated = true
	}

	if input.Magnesium != 0 && input.Magnesium != food.Magnesium {
		config.DB.Model(&food).Updates(map[string]interface{}{"magnesium": input.Magnesium})
		isUpdated = true
	}
	if input.Calcium != 0 && input.Calcium != food.Calcium {
		config.DB.Model(&food).Updates(map[string]interface{}{"calcium": input.Calcium})
		isUpdated = true
	}
	if input.Sodium != 0 && input.Sodium != food.Sodium {
		config.DB.Model(&food).Updates(map[string]interface{}{"sodium": input.Sodium})
		isUpdated = true
	}
	if input.Iron != 0 && input.Iron != food.Iron {
		config.DB.Model(&food).Updates(map[string]interface{}{"iron": input.Iron})
		isUpdated = true
	}
	if input.VitA != 0 && input.VitA != food.VitA {
		config.DB.Model(&food).Updates(map[string]interface{}{"vit_a": input.VitA})
		isUpdated = true
	}
	if input.VitC != 0 && input.VitC != food.VitC {
		config.DB.Model(&food).Updates(map[string]interface{}{"vit_c": input.VitC})
		isUpdated = true
	}
	if input.VitD != 0 && input.VitD != food.VitD {
		config.DB.Model(&food).Updates(map[string]interface{}{"vit_d": input.VitD})
		isUpdated = true
	}
	if input.Sugar != 0 && input.Sugar != food.Sugar {
		config.DB.Model(&food).Updates(map[string]interface{}{"sugar": input.Sugar})
		isUpdated = true
	}
	if input.SaturatedFat != 0 && input.SaturatedFat != food.SaturatedFat {
		config.DB.Model(&food).Updates(map[string]interface{}{"saturated_fat": input.SaturatedFat})
		isUpdated = true
	}
	if input.UnsaturatedFat != 0 && input.UnsaturatedFat != food.UnsaturatedFat {
		config.DB.Model(&food).Updates(map[string]interface{}{"unsaturated_fat": input.UnsaturatedFat})
		isUpdated = true
	}
	if input.Gram != 0 && input.Gram != food.Gram {
		config.DB.Model(&food).Updates(map[string]interface{}{"gram": input.Gram})
		isUpdated = true
	}
	if input.Meal != "" && input.Meal != food.Meal {
		config.DB.Model(&food).Updates(map[string]interface{}{"meal": input.Meal})
		isUpdated = true
	}

	msg := "Successful operation, object has been updated."
	if isUpdated == false {
		msg = "There is nothing to update, please change an attribute of the object."
		ctx.JSON(http.StatusOK, gin.H{"msg": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": &food, "msg": msg})
}

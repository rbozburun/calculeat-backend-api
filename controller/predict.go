package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/gin-gonic/gin"
)

var PREDICT_API string

func init() {
	PREDICT_API = logger.GoDotEnvVariable("PREDICT_API")
}

// POST /api/1.0/predict
func Predict(ctx *gin.Context) {
	var predict models.Predict
	err := ctx.ShouldBind(&predict)
	if err != nil {
		logger.Log.Errorln("Invalid Request Format, couldnt bind json. Request Header:", ctx.Request.Header)
		logger.Log.Errorln("Invalid Request Format, couldnt bind json. Request Body:", ctx.Request.Body)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid Request Format, couldnt bind predict object's json"})
		return
	}

	// Read the recieved file
	file, err := predict.File.Open()
	if err != nil {
		logger.Log.Errorln("Error retrieving the file: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Error retrieving the file"})
		return
	}
	defer file.Close()

	// Create a new file on the server
	out, err := os.Create("image.png")
	if err != nil {
		logger.Log.Errorln("Cannot create the local file: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Cannot create the local file"})
		return
	}
	defer out.Close()

	// Copy the file content to the server file
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Log.Errorln("Cannot copy content of remote file to local file: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Cannot copy content of remote file to local file"})
		return
	}
	filePath := "image.png"

	// Call the AI Predict endpoint
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		logger.Log.Errorln("Cannot prepare form-data request for AI endpoint -> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	fd, err := os.Open(filePath)
	if err != nil {
		logger.Log.Errorln("Cannot open the local file: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		logger.Log.Errorln("io.Copy operation error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", PREDICT_API, form)
	if err != nil {
		logger.Log.Errorln("Request creation error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Request creation error."})
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Errorln("Request error while requesting to : " + PREDICT_API + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Request error while requesting to: " + PREDICT_API,
			"err": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorln("Error reading response:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	var prediction models.PredictionResult
	err = json.Unmarshal(body, &prediction)
	if err != nil {
		err_msg := "Error decoding response from AI API, cannot unmarshal the json:" + err.Error()
		logger.Log.Errorln("Error decoding response from AI API, cannot unmarshal the prediction json:", err)
		logger.Log.Debugln("AI Response Body: ", string(body))
		logger.Log.Debugln("Prediction object: ", &prediction)
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err_msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Successfully calculeated", "data": &prediction})
	return
}

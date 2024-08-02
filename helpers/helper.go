package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/calculeat/main_rest_api/config"
	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
)

const (
	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
)

var (
	apiKey = logger.GoDotEnvVariable("API_KEY")
)

func SetColForIntVal(tb string, target_col string, val int, target_pk_id int) {
	logger.Log.Debugf("Updating the %v column of %v table for the value %v where the id is %v", target_col, tb, val, target_pk_id)
	sql_query := "UPDATE " + tb + " SET " + target_col + " = " + strconv.Itoa(val) + " WHERE " + tb + ".id = " + strconv.Itoa(target_pk_id)
	logger.Log.Debugf("Using SQL query: %v", sql_query)
	config.DB.Exec(sql_query)

}

// Returns the (dieticianUserID,clientUserID) after parsing the targetUserID
func GetDieticianAndClientUserIDs(target_user_id uint64) (uint64, uint64) {
	// Current user is a dietician
	if config.CURRENT_USER.RoleID == 1 {
		return config.CURRENT_USER.ID, target_user_id
	} else if config.CURRENT_USER.RoleID == 2 {
		// Current user is a client
		return target_user_id, config.CURRENT_USER.ID
	} else {
		return 0, 0
	}
}

// Initiliaze water object in db
func InitiliazeWaterObject(user_id uint64) (models.WaterObject, error) {
	var water models.WaterObject
	water.Count = 0
	water.UserID = user_id

	if err := config.DB.Create(&water).Error; err != nil {
		logger.Log.Errorf(logger.BindError+" Err: %v", err)
		logger.Log.Debugf("PROBLEMATIC WATER OBJ: %v", water)
		logger.Log.Debugf("PROBLEMATIC USER_ID: %v, WATER.USER_ID: %v", user_id, water.UserID)
		return models.WaterObject{}, err
	} else {
		logger.Log.Info(logger.ObjectCreated)
		// Update the user's init_water_id column in db
		SetColForIntVal("users", "init_water_id", int(water.ID), int(user_id))

		return water, nil
	}

}

// Find the user model with the provided email and return it
func FindUserByEmail(email string) (models.User, error) {
	var user models.User

	err := config.DB.Where("email=?", email).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	if user.Email != email {
		logger.Log.Debugf("Email not found! Input email: %v, Database model email: %v", email, user.Email)
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// Function to convert an array to string
func ArrayToString(arr []string) string {
	var str_rep_of_arr string
	for _, value := range arr {
		str_rep_of_arr += value
	}
	return str_rep_of_arr
}

// Fuction to check given string is in array or not
func IsExistsOnArray(array []string, elem string) bool {
	for _, v := range array {
		if v == elem {
			return true
		}
	}
	return false
}

// Compare two arrays and their contents to check they are equal or not.
func IsArraysEqual(arr1 []string, arr2 []string) bool {
	array1_as_str := ArrayToString(arr1)
	array2_as_str := ArrayToString(arr2)

	logger.Log.Debugln("Comparing the arrays... Arr1: " + array1_as_str + " Arr2: " + array2_as_str)
	if len(arr1) != len(arr2) {
		return false
	}
	for _, v1 := range arr1 {
		if IsExistsOnArray(arr2, v1) == false {
			return false
		}
	}
	return true
}

// Converts base64 string to provided filepath as a file
func ConvertBase64StringToFile(b64_string string, filePath string) {
	// Decode base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(b64_string)
	if err != nil {
		logger.Log.Errorln("Error decoding base64 string:", err)
		return
	}

	// Write decoded bytes to file
	err = ioutil.WriteFile(filePath, decodedBytes, 0644)
	if err != nil {
		logger.Log.Errorln("Error writing to file:", err)
		return
	}
}

// Compares two passwords (plaintext, sha256) and returns true/false
func ComparePasswords(input_pwd string, db_user_pwd_hash string) (bool, error) {
	logger.Log.Debugln("Comparing the supplied password and password stored in db.")
	pwdHash_input_byte := sha256.Sum256([]byte(input_pwd))
	pwdHash_input := fmt.Sprintf("%x", pwdHash_input_byte[:])

	if db_user_pwd_hash == pwdHash_input {
		logger.Log.Debugln("Password match!")
		return true, nil
	}

	logger.Log.Debugln("Password does not match!")
	return false, errors.New("invalid password")
}

func postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("unexpected http status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

// Firebase sign-in with Custom Token
func SignInWithCustomToken(token string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, apiKey), req)
	if err != nil {
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}
	return respBody.IDToken, err
}

// Validate the given date. Valid date format: YYYY-MM-DD
func ValidateDate(date string) error {
	// Validate the date
	match, _ := regexp.MatchString("^([0-9]{4}-[0-9]{2}-[0-9]{2})$", date)
	if !match {
		return errors.New("invalid date format, use YYYY-MM-DD")
	}
	return nil
}

// Find the water objects by date
func FindWaterObjectsByDate(date string) ([]models.WaterObject, error) {
	water_objects := []models.WaterObject{}

	// Validate the date
	err := ValidateDate(date)
	if err != nil {
		return water_objects, errors.New("invalid date format, use YYYY-MM-DD")
	}

	config.DB.Where("created_at::date = ?", date).Find(&water_objects)
	if len(water_objects) > 0 {
		return water_objects, nil
	} else {
		return water_objects, errors.New("there is no water object for the specified date")
	}

}

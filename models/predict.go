package models

import "mime/multipart"

type Predict struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type PredictionResult struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
}

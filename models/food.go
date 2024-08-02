package models

import "time"

type FoodObject struct {
	ID               uint64  `json:"id" gorm:"primary_key;auto_increment;default:null"`
	Name             string  `json:"name" gorm:""`
	UserID           uint64  `json:"user_id"`
	ImgleftFilepath  string  `json:"imgleft_filepath"`
	ImgrightFilepath string  `json:"imgright_filepath"`
	ImgtopFilepath   string  `json:"imgtop_filepath"`
	Calorie          uint64  `json:"calorie"`
	Protein          float64 `json:"protein"`
	Carbonhydrate    float64 `json:"carbonhydrate"`
	Fat              float64 `json:"fat"`
	Magnesium        float64 `json:"magnesium"`
	Calcium          float64 `json:"calcium"`
	Sodium           float64 `json:"sodium"`
	Iron             float64 `json:"iron"`
	VitA             uint64  `json:"vit_a"`
	VitC             uint64  `json:"vit_c"`
	VitD             uint64  `json:"vit_d"`
	Sugar            float64 `json:"sugar"`
	SaturatedFat     float64 `json:"saturated_fat"`
	UnsaturatedFat   float64 `json:"unsaturated_fat"`
	Gram             float64 `json:"gram"`
	Meal             string  `json:"meal"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedOn        time.Time `json:"deleted_on"`
}

type UpdateFoodInput struct {
	Name             string  `json:"name" gorm:""`
	UserID           uint64  `json:"user_id"`
	ImgleftFilepath  string  `json:"imgleft_filepath"`
	ImgrightFilepath string  `json:"imgright_filepath"`
	ImgtopFilepath   string  `json:"imgtop_filepath"`
	Calorie          uint64  `json:"calorie"`
	Protein          float64 `json:"protein"`
	Carbonhydrate    float64 `json:"carbonhydrate"`
	Fat              float64 `json:"fat"`
	Magnesium        float64 `json:"magnesium"`
	Calcium          float64 `json:"calcium"`
	Sodium           float64 `json:"sodium"`
	Iron             float64 `json:"iron"`
	VitA             uint64  `json:"vit_a"`
	VitC             uint64  `json:"vit_c"`
	VitD             uint64  `json:"vit_d"`
	Sugar            float64 `json:"sugar"`
	SaturatedFat     float64 `json:"saturated_fat"`
	UnsaturatedFat   float64 `json:"unsaturated_fat"`
	Gram             float64 `json:"gram"`
	Meal             string  `json:"meal"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedOn        time.Time `json:"deleted_on"`
}

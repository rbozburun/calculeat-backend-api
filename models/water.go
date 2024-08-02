package models

import "time"

type WaterObject struct {
	ID        uint64 `json:"id" gorm:"primary_key;auto_increment;default:null"`
	UserID    uint64 `json:"user_id"`
	Count     int    `json:"count"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedOn time.Time `json:"deleted_on"`
}

type UpdateWaterInput struct {
	UserID    uint64 `json:"user_id"`
	Count     int    `json:"count"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedOn time.Time `json:"deleted_on"`
}

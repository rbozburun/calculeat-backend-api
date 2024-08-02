package models

import "time"

type SleepObject struct {
	ID        uint64 `json:"id" gorm:"primary_key;auto_increment;default:null"`
	UserID    uint64 `json:"user_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedOn time.Time `json:"deleted_on"`
}

type UpdateSleepInput struct {
	UserID    uint64 `json:"user_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	UpdatedAt time.Time
}

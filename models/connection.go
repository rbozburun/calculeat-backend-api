package models

import "time"

type Connection struct {
	ID              uint64 `json:"id" gorm:"primary_key;auto_increment;default:null"`
	ClientUserID    uint64 `json:"client_user_id" gorm:"default:null"`
	DieticianUserID uint64 `json:"dietician_user_id" gorm:"default:null"`
	IsAccepted      string `json:"is_accepted" gorm:"default:null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedOn       time.Time `json:"deleted_on"`
}

type CreateConnectionInput struct {
	TargetUserID uint64 `json:"target_user_id" gorm:"default:null"`
	IsAccepted   string `json:"meet_link"`
}

type UpdateConnectionInput struct {
	ClientUserID    uint64 `json:"client_user_id" gorm:"default:null"`
	DieticianUserID uint64 `json:"dietician_user_id" gorm:"default:null"`
	IsAccepted      string `json:"is_accepted" gorm:"default:null"`
	UpdatedAt       time.Time
}

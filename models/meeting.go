package models

import "time"

type Meeting struct {
	ID              uint64 `json:"id" gorm:"primary_key;auto_increment;default:null"`
	DieticianUserID uint64 `json:"dietician_user_id" gorm:"default:null"`
	ClientUserID    uint64 `json:"client_user_id" gorm:"default:null"`
	MeetLink        string `json:"meet_link"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedOn       time.Time `json:"deleted_on"`
}

type CreateMeetingInput struct {
	TargetUserID uint64 `json:"target_user_id" gorm:"default:null"`
	MeetLink     string `json:"meet_link"`
}

type UpdateMeetingInput struct {
	DieticianUserID uint64 `json:"dietician_user_id" gorm:"default:null"`
	ClientUserID    uint64 `json:"client_user_id" gorm:"default:null"`
	MeetLink        string `json:"meet_link"`
	UpdatedAt       time.Time
}

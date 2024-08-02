package models

import "time"

type Message struct {
	ID          uint64 `json:"id" gorm:"primary_key;auto_increment;default:null"`
	SenderID    uint64 `json:"sender_id"`
	RecieverID  uint64 `json:"reciever_id"`
	MessageText string `json:"message_text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedOn   time.Time `json:"deleted_on"`
}

type UpdateMessageInput struct {
	SenderID    uint64 `json:"sender_id"`
	RecieverID  uint64 `json:"reciever_id"`
	MessageText string `json:"message_text"`
	UpdatedAt   time.Time
}

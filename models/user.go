package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID                 uint64  `json:"id" gorm:"primary_key;auto_increment"`
	AccessToken        string  `json:"access_token" gorm:"default:null"`
	Password           string  `json:"password" gorm:"not null;default:null"`
	RoleID             int     `json:"role_id" gorm:"not null;default:null" binding:"required"`
	InitWaterId        int     `json:"init_water_id" gorm:"default:null"`
	Email              string  `json:"email" gorm:"not null;default:null;unique" binding:"required"`
	FirstName          string  `json:"first_name" gorm:"not null;default:null" binding:"required"`
	LastName           string  `json:"last_name" gorm:"not null;default:null" binding:"required"`
	ProfileImgFilepath string  `json:"profile_img_filepath" gorm:"default:null"`
	SignInProvider     string  `json:"sign_in_provider"`
	Gender             string  `json:"gender" gorm:"not null;default:null" binding:"required"`
	Born               string  `json:"born" gorm:"not null;default:null" binding:"required"`
	Age                int     `json:"age" gorm:"not null;default:null" `
	Height             int     `json:"height" gorm:"not null;default:null" binding:"required"`
	Weight             float64 `json:"weight" gorm:"not null;default:null" binding:"required"`
	Goal               string  `json:"goal"`
	//Diseases should be converted to list. It comes from as a string with comma seperated from Db.
	Diseases  pq.StringArray `json:"diseases" gorm:"type:text[]" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedOn time.Time `json:"deleted_on"`
}

type Users struct {
	Users []User
}

type UpdateUserInput struct {
	RoleID             int     `json:"role_id"`
	InitWaterId        int     `json:"init_water_id" gorm:"default:null"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	ProfileImgFilepath string  `json:"profile_img_filepath"`
	Gender             string  `json:"gender"`
	Born               string  `json:"born"`
	Age                int     `json:"age"`
	Email              string  `json:"email"`
	Height             int     `json:"height"`
	Weight             float64 `json:"weight"`
	Goal               string  `json:"goal"`
	//Diseases should be converted to list. It comes from as a string with comma seperated from Db.
	Diseases pq.StringArray `json:"diseases" gorm:"type:text[]"`
	Password string         `json:"password"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

package config

import (
	"fmt"
	"time"

	"github.com/calculeat/main_rest_api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connects to the database
func Connect() {
	db, err := gorm.Open(postgres.Open("postgres://calculeat_test:AwsSwh57134zxcYt@localhost:5433/calculeat_test?sslmode=disable"), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	start := time.Now()
	for sqlDB.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			fmt.Println("Failed to connect db after 10 secs.")
			break
		}
	}
	fmt.Println("connected: ", sqlDB.Ping() == nil)

	db.AutoMigrate(&models.User{})
	DB = db
}

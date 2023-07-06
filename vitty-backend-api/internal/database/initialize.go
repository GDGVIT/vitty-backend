package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(debug string, dbUrls string) {
	var err error

	if debug == "true" {
		DB, err = gorm.Open(postgres.Open(dbUrls), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		DB, err = gorm.Open(postgres.Open(dbUrls), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	log.Println("Connected to database!")
}

package database

import (
	"log"
	"os"

	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	if os.Getenv("DEBUG") == "1" {
		// Sqlite
		DB, err = gorm.Open(sqlite.Open("vitty.db"), &gorm.Config{})
	} else {
		// Postgresql
		DB, err = gorm.Open(postgres.Open(os.Getenv("POSTGRES_URL")), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	// Migrating schemas
	DB.AutoMigrate(&models.User{}, &models.Timetable{}, &models.FriendRequest{})
}

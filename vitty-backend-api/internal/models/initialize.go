package models

import (
	"log"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
)

func InitializeModels() {
	MODELS := map[string]interface{}{
		"User":            &User{},
		"Timetable":       &Timetable{},
		"Friend Requests": &FriendRequest{},
	}

	for name, model := range MODELS {
		err := database.DB.AutoMigrate(model)
		if err != nil {
			log.Fatal("Failed to initialize model: ", name)
		} else {
			log.Println("Successfully initialized model: ", name)
		}
	}
}

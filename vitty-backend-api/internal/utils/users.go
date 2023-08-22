package utils

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
)

func CheckUserExists(username string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count != 0
}

func CheckUserByEmail(email string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count != 0
}

func GetUserByUsername(username string) models.User {
	var user models.User
	database.DB.Model(&models.User{}).Where("username = ?", username).Preload("Friends").First(&user)
	return user
}

func CheckUserByUUID(uuid string) bool {
	var count int64
	database.DB.Model(&models.User{}).Where("firebase_uuid = ?", uuid).Count(&count)
	return count != 0
}

func ValidateUsername(username string) (bool, string) {
	// Username should be between 3 and 20 characters
	if len(username) < 3 || len(username) > 20 {
		return false, "Username should be between 3 and 20 characters"
	}
	// Username should be alphanumeric with no spaces, only . and _
	for _, char := range username {
		if !(char >= 'a' && char <= 'z') &&
			!(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') &&
			char != '.' &&
			char != '_' {
			return false, "Username should be alphanumeric with no spaces, only . and _"
		}
	}

	if CheckUserExists(username) {
		return false, "Username already exists"
	}

	return true, ""
}

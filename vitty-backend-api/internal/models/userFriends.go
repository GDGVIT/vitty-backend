package models

import (
	"time"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
)

type UserFriends struct {
	UserUsername   string     `gorm:"primaryKey"`
	FriendUsername string     `gorm:"primaryKey"`
	Hide           bool       `gorm:"default:false"`
	UpdatedAt      *time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (uf *UserFriends) GetActiveFriends(username string) ([]UserFriends, error) {
	var userFriends []UserFriends

	err := database.DB.Where("user_username = ? AND hide = false", username).Find(&userFriends).Error

	return userFriends, err
}

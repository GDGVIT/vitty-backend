package models

import "time"

type UserFriends struct {
	UserUsername   string     `gorm:"primaryKey"`
	FriendUsername string     `gorm:"primaryKey"`
	Hide           bool       `gorm:"default:false"`
	UpdatedAt      *time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

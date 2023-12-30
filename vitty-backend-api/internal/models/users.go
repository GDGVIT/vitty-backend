package models

import (
	"time"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm/clause"
)

type User struct {
	Username     string  `gorm:"primaryKey"`
	RegNo        string  `gorm:"unique"`
	Name         string  `gorm:"not null"`
	Email        string  `gorm:"unique,not null"`
	Role         string  `gorm:"default:normal"`
	Picture      string  `gorm:"not null"`
	Friends      []*User `gorm:"many2many:user_friends;foreignKey:Username;joinForeignKey:UserUsername;References:Username;joinReferences:FriendUsername"`
	FirebaseUuid string  `gorm:"unique"`
}

func (u *User) GetCurrentStatus() map[string]interface{} {
	// Check if user is currently in class
	// If yes, return map with class details
	// If no, status = free
	// Get current time
	time := time.Now()
	daySlots := u.GetTimeTable().GetDaySlots(time.Weekday())
	for _, slot := range daySlots[time.Weekday().String()] {
		if slot.StartTime.Before(time) && slot.EndTime.After(time) {
			return map[string]interface{}{
				"status": "class",
				"slot":   slot.Slot,
				"venue":  slot.Venue,
			}
		}
	}
	return map[string]interface{}{
		"status": "free",
	}
}

func (u *User) GetFriendRequests() []FriendRequest {
	friend_requests := []FriendRequest{}
	database.DB.Where("to_username= ?", u.Username).Preload(clause.Associations).Find(&friend_requests)
	return friend_requests
}

func (u *User) RemoveFriend(friend User) {
	database.DB.Model(u).Association("Friends").Delete(friend)
	database.DB.Model(friend).Association("Friends").Delete(u)
}

func (u *User) IsFriendsWith(user User) bool {
	var count int64
	database.DB.Table("user_friends").
		Where("user_username = ? AND friend_username = ?", u.Username, user.Username).
		Count(&count)
	return count != 0
}

func (u *User) HasSentFriendRequest(user User) bool {
	var count int64
	database.DB.Model(&FriendRequest{}).Where("from_username = ? AND to_username = ?", u.Username, user.Username).Count(&count)
	return count != 0
}

func (u *User) HasReceivedFriendRequest(user User) bool {
	var count int64
	database.DB.Model(&FriendRequest{}).Where("from_username = ? AND to_username = ?", user.Username, u.Username).Count(&count)
	return count != 0
}

func (u *User) CheckFriendStatus(user User) string {
	if u.Username == user.Username {
		return "self"
	} else if u.IsFriendsWith(user) {
		return "friends"
	} else if u.HasSentFriendRequest(user) {
		return "sent"
	} else if u.HasReceivedFriendRequest(user) {
		return "received"
	}
	return "none"
}

func (u *User) GetTimeTable() Timetable {
	timetable := Timetable{}
	database.DB.Where("user_username = ?", u.Username).First(&timetable)
	return timetable
}

func (u *User) FriendsCount() int64 {
	return database.DB.Model(&u).Association("Friends").Count()
}

func (u *User) FindMutualFriends(otherUser User) []*User {
	mutualFriends := []*User{}

	database.DB.Model(u).
		Joins("JOIN user_friends AS uf ON uf.friend_username = users.username").
		Joins("JOIN users AS u2 ON u2.username = uf.user_username").
		Where("u2.username = ?", otherUser.Username).
		Find(&mutualFriends)

	return mutualFriends
}

func (u *User) CountMutualFriends(otherUser User) int64 {
	var count int64
	database.DB.Model(u).
		Joins("JOIN user_friends AS uf ON uf.friend_username = users.username").
		Joins("JOIN users AS u2 ON u2.username = uf.user_username").
		Where("u2.username = ?", otherUser.Username).
		Count(&count)
	return count
}

func (u *User) FindSuggestedOnMutualFriends() []*User {
	suggested := []*User{}

	database.DB.Raw(`
		SELECT users.*, COUNT(users.*) as mf_count 
		FROM users 
		JOIN user_friends AS uf1 ON users.username = uf1.friend_username 
		JOIN user_friends AS uf2 ON uf1.user_username = uf2.user_username 
		WHERE uf2.friend_username = ? 
		AND users.username <> ?
		AND users.username NOT IN (SELECT uf3.friend_username
								   FROM user_friends AS uf3
								   WHERE uf3.user_username = ?)
		GROUP BY users.username 
		ORDER BY mf_count;
	`, u.Username, u.Username, u.Username).Scan(&suggested)

	return suggested
}

func (u *User) Save() error {
	return database.DB.Save(&u).Error
}

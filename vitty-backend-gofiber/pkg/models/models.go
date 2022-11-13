package models

import "time"

type Slot struct {
	Slot           string `json:"Slot"`
	CourseName     string `json:"Course_Name"`
	CourseFullName string `json:"Course_Full_Name,omitempty"`
	CourseType     string `json:"Course_type"`
	Venue          string `json:"Venue"`
}

type Timetable struct {
	ID        uint   `json:"id"`
	Slots     []Slot `json:"slots" gorm:"serializer:json"`
	UserRegNo string `json:"user_regno"`
	User      User   `json:"-" gorm:"foreignKey:UserRegNo"`
}

type User struct {
	ID        string    `json:"regno"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	Friends   []*User   `gorm:"many2many:user_friends"`
}

type FriendRequest struct {
	ID        uint   `json:"id"`
	FromRegNo string `json:"from"`
	ToRegNo   string `json:"to"`
	From      User   `json:"from_user" gorm:"foreignKey:FromRegNo"`
	To        User   `json:"to_user" gorm:"foreignKey:ToRegNo"`
}

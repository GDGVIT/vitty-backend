package schemas

import (
	"time"

	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
)

type UserResponse struct {
	ID        string              `json:"regno"`
	Email     string              `json:"email"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Friends   []UserBlockResponse `json:"friends"`
	Timetable TimetableResponse   `json:"timetable"`
}

type UserBlockResponse struct {
	ID        string `json:"regno"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserLoginResponse struct {
	ID        string `json:"regno"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

type TimetableResponse struct {
	ID        uint          `json:"id"`
	Slots     []models.Slot `json:"slots"`
	UserRegNo string        `json:"user_regno"`
}

type FriendRequestResponse struct {
	ID        uint              `json:"id"`
	FromRegNo string            `json:"from_regno"`
	FromUser  UserBlockResponse `json:"from"`
	ToUser    UserBlockResponse `json:"to"`
	ToRegNo   string            `json:"to_regno"`
}

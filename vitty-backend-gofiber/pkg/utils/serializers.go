package utils

import (
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/database"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/schemas"
)

func UserLoginSerializer(user *models.User, token string) schemas.UserLoginResponse {
	return schemas.UserLoginResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
	}
}

func UserSerializer(user *models.User) schemas.UserResponse {
	timetable := models.Timetable{}

	if database.DB.Where("user_reg_no = ?", user.ID).First(&timetable).RowsAffected == 0 {
		timetable = models.Timetable{}
	}

	return schemas.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Friends:   UserBlockSerializerSlice(user.Friends),
		Timetable: TimetableSerializer(&timetable),
	}
}

func UserBlockSerializerSlice(users []*models.User) []schemas.UserBlockResponse {
	var userBlockResponses []schemas.UserBlockResponse
	for _, user := range users {
		userBlockResponses = append(userBlockResponses, UserBlockSerializer(user))
	}
	return userBlockResponses
}

func UserBlockSerializer(user *models.User) schemas.UserBlockResponse {
	return schemas.UserBlockResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func TimetableSerializer(timetable *models.Timetable) schemas.TimetableResponse {
	return schemas.TimetableResponse{
		ID:        timetable.ID,
		Slots:     timetable.Slots,
		UserRegNo: timetable.UserRegNo,
	}
}

func FriendRequestSerializer(friendRequest *models.FriendRequest) schemas.FriendRequestResponse {
	return schemas.FriendRequestResponse{
		ID:        friendRequest.ID,
		FromUser:  UserBlockSerializer(&friendRequest.From),
		FromRegNo: friendRequest.FromRegNo,
		ToUser:    UserBlockSerializer(&friendRequest.To),
		ToRegNo:   friendRequest.ToRegNo,
	}
}

func FriendRequestSerializerSlice(friendRequests []models.FriendRequest) []schemas.FriendRequestResponse {
	var friendRequestResponses []schemas.FriendRequestResponse
	for _, friendRequest := range friendRequests {
		friendRequestResponses = append(friendRequestResponses, FriendRequestSerializer(&friendRequest))
	}
	return friendRequestResponses
}

package serializers

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	vittyerrors "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/vittyErrors"
)

func UserLoginSerializer(user models.User, token string) map[string]interface{} {
	return map[string]interface{}{
		"username": user.Username,
		"name":     user.Name,
		"picture":  user.Picture,
		"role":     user.Role,
		"token":    token,
	}
}

func UserCardSerializer(user models.User, request_user models.User) (map[string]interface{}, error) {
	friendStatus := request_user.CheckFriendStatus(user)
	var err error
	var currStatus map[string]interface{}
	if friendStatus == "friends" {
		currStatus, err = user.GetCurrentStatus()
		if err != nil {
			return nil, vittyerrors.ErrFetchCurrentStatus
		}
	} else {
		currStatus = map[string]interface{}{
			"status": "unknown",
		}
	}
	return map[string]interface{}{
		"username":             user.Username,
		"name":                 user.Name,
		"picture":              user.Picture,
		"friends_count":        user.FriendsCount(),
		"friend_status":        friendStatus,
		"mutual_friends_count": request_user.CountMutualFriends(user),
		"current_status":       currStatus,
	}, nil
}

func UserListSerializer(users []*models.User, request_user models.User) ([]map[string]interface{}, error) {
	var users_list []map[string]interface{}
	for _, user := range users {
		userCard, err := UserCardSerializer(*user, request_user)
		if err != nil {
			return users_list, err
		}
		users_list = append(users_list, userCard)
	}
	return users_list, nil
}

func UserSerializer(user models.User, request_user models.User) map[string]interface{} {
	return map[string]interface{}{
		"username":             user.Username,
		"name":                 user.Name,
		"picture":              user.Picture,
		"email":                user.Email,
		"timetable":            TimetableSerializer(user.GetTimeTable()),
		"friend_status":        request_user.CheckFriendStatus(user),
		"friends_count":        user.FriendsCount(),
		"mutual_friends_count": request_user.CountMutualFriends(user),
	}
}

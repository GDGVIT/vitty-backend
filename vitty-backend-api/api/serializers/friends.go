package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func FriendRequestsSerializer(friend_requests []models.FriendRequest, request_user models.User) []map[string]interface{} {
	var friend_requests_list []map[string]interface{}
	for _, friend_request := range friend_requests {
		friend_requests_list = append(friend_requests_list, map[string]interface{}{
			"from": UserCardSerializer(friend_request.From, request_user),
		})
	}
	return friend_requests_list
}

func ActiveFriendsSerializer(userFriends *[]models.UserFriends) []map[string]interface{} {
	var activeFriends []map[string]interface{}

	for _, userFriend := range *userFriends {
		activeFriends = append(activeFriends, map[string]interface{}{
			"friend_username": userFriend.FriendUsername,
			"hide":            userFriend.Hide,
		})
	}
	return activeFriends
}

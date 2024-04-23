package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func FriendRequestsSerializer(friend_requests []models.FriendRequest, request_user models.User) ([]map[string]interface{}, error) {
	var friend_requests_list []map[string]interface{}
	for _, friend_request := range friend_requests {
		userCard, err := UserCardSerializer(friend_request.From, request_user)
		if err != nil {
			return nil, err
		}
		friend_requests_list = append(friend_requests_list, map[string]interface{}{
			"from": userCard,
		})
	}
	return friend_requests_list, nil
}

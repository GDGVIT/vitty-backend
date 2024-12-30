package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func CirclesListSerializer(ucj []models.UsersCirclesJoin) []map[string]interface{} {
	var result []map[string]interface{}

	for _, userCircle := range ucj {

		out := map[string]interface{}{
			"circle_id":   userCircle.CID,
			"circle_role": userCircle.CircleRole,
			"circle_name": userCircle.Circles.CircleName,
		}
		result = append(result, out)
	}

	return result

}

func CircleRequestsSerializer(circleRequests []models.CircleRequest) []map[string]interface{} {
	var result []map[string]interface{}

	for _, circleRequest := range circleRequests {
		out := map[string]interface{}{
			"from_username": circleRequest.FromUsername,
			"to_username":   circleRequest.ToUsername,
			"circle_id":     circleRequest.CID,
			"circle_name":   circleRequest.Circles.CircleName,
		}
		result = append(result, out)
	}

	return result
}

func UsersListCircleSerializer(users []models.User) []map[string]interface{} {
	var result []map[string]interface{}

	for _, user := range users {
		out := map[string]interface{}{
			"username": user.Username,
			"name":     user.Name,
			"picture":  user.Picture,
			"email":    user.Email,
		}
		result = append(result, out)
	}

	return result
}

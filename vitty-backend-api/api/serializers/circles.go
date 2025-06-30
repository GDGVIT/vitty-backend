package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func CirclesListSerializer(ucj []models.UsersCirclesJoin) []map[string]interface{} {
	var result []map[string]interface{}

	for _, userCircle := range ucj {

		out := map[string]interface{}{
			"circle_id":        userCircle.CID,
			"circle_role":      userCircle.CircleRole,
			"circle_name":      userCircle.Circles.CircleName,
			"circle_join_code": userCircle.Circles.CircleJoinCode,
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

		currStatus := user.GetCurrentStatus()

		out := map[string]interface{}{
			"current_status": currStatus,
			"username":       user.Username,
			"name":           user.Name,
			"picture":        user.Picture,
			"email":          user.Email,
		}
		result = append(result, out)
	}

	return result
}

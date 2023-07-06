package v2

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func getFriendRequests(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(serializers.FriendRequestsSerializer(request_user.GetFriendRequests(), request_user))
}

func createFriendRequest(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "User does not exist",
		})
	}
	user := utils.GetUserByUsername(username)

	if request_user.Username == user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You cannot send a friend request to yourself",
		})
	}
	if request_user.IsFriendsWith(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You are already friends with this user",
		})
	}
	if request_user.HasSentFriendRequest(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You have already sent a friend request to this user",
		})
	}
	if request_user.HasReceivedFriendRequest(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You have already received a friend request from this user",
		})
	}

	var friend_request models.FriendRequest
	friend_request.From = request_user
	friend_request.To = user
	err = database.DB.Create(&friend_request).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Friend request sent successfully",
	})
}

func acceptFriendRequest(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "User does not exist",
		})
	}
	user := utils.GetUserByUsername(username)

	if request_user.Username == user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You cannot accept a friend request from yourself",
		})
	}
	if request_user.IsFriendsWith(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You are already friends with this user",
		})
	}
	if !request_user.HasReceivedFriendRequest(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You have not received a friend request from this user",
		})
	}

	var friend_request models.FriendRequest
	database.DB.Where("from_username = ? AND to_username = ?", user.Username, request_user.Username).Preload(clause.Associations).First(&friend_request)
	friend_request.Accept()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Friend request accepted successfully!",
	})
}

func declineFriendRequest(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "User does not exist",
		})
	}
	user := utils.GetUserByUsername(username)

	if request_user.Username == user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You cannot reject a friend request from yourself",
		})
	}
	if request_user.IsFriendsWith(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You are already friends with this user",
		})
	}
	if !request_user.HasReceivedFriendRequest(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You have not received a friend request from this user",
		})
	}

	var friend_request models.FriendRequest
	database.DB.Where("from_username = ? AND to_username = ?", user.Username, request_user.Username).First(&friend_request)
	friend_request.Decline()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Friend request rejected successfully",
	})
}

func getFriends(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}
	user := utils.GetUserByUsername(username)
	if !user.IsFriendsWith(request_user) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"friend_status": request_user.CheckFriendStatus(user),
			"data":          serializers.UserListSerializer(request_user.FindMutualFriends(user), request_user),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"friend_status": request_user.CheckFriendStatus(user),
		"data":          serializers.UserListSerializer(user.Friends, request_user),
	})
}

func removeFriend(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "User does not exist",
		})
	}
	user := utils.GetUserByUsername(username)

	if request_user.Username == user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You cannot remove yourself as a friend",
		})
	}
	if !request_user.IsFriendsWith(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "You are not friends with this user",
		})
	}

	request_user.RemoveFriend(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Friend removed successfully",
	})
}

func friendHandler(api fiber.Router) {
	requestGroup := api.Group("/requests")
	requestGroup.Get("/", getFriendRequests)
	requestGroup.Post("/:username/send", createFriendRequest)
	requestGroup.Post("/:username/accept", acceptFriendRequest)
	requestGroup.Post("/:username/decline", declineFriendRequest)

	friendGroup := api.Group("/friends")
	friendGroup.Get("/:username", getFriends)
	friendGroup.Delete("/:username", removeFriend)
}

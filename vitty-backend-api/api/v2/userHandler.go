package v2

import (
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func searchUsers(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	query := c.Query("query")
	var users []*models.User
	database.DB.Where("username LIKE ? OR name LIKE ?", query+"%", query+"%").Find(&users)
	return c.Status(fiber.StatusOK).JSON(serializers.UserListSerializer(users, request_user))
}

func getSuggestedUsers(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(serializers.UserListSerializer(request_user.FindSuggestedOnMutualFriends(), request_user))
}

func getUsers(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	if request_user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to perform this action",
		})
	}
	var users []*models.User
	database.DB.Find(&users)
	return c.Status(fiber.StatusOK).JSON(serializers.UserListSerializer(users, request_user))
}

func getUser(c *fiber.Ctx) error {
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
	fmt.Println("Friends", user.IsFriendsWith(request_user))

	if (user.Username == request_user.Username) ||
		(user.IsFriendsWith(request_user)) ||
		(request_user.Role == "admin") {
		return c.Status(fiber.StatusOK).JSON(serializers.UserSerializer(user, request_user))
	}
	return c.Status(fiber.StatusOK).JSON(serializers.UserCardSerializer(user, request_user))
}

func deleteUser(c *fiber.Ctx) error {
	request_user, err := auth.GetUserFromJWTToken(c.Get("Authorization"), auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	c.Params("username")
	if request_user.Username != c.Params("username") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to delete this user",
		})
	}

	database.DB.Delete(&request_user)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "User deleted successfully",
	})
}

func userHandler(api fiber.Router) {
	group := api.Group("/users")
	group.Get("/", getUsers)
	group.Get("/search", searchUsers)
	group.Get("/suggested", getSuggestedUsers)
	group.Get("/:username", getUser)
	group.Delete("/:username", deleteUser)
}

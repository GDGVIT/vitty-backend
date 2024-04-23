package v2

import (
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func userHandler(api fiber.Router) {
	group := api.Group("/users")
	group.Use(middleware.JWTAuthMiddleware)
	group.Get("/", getUsers)
	group.Get("/search", searchUsers)
	group.Get("/suggested", getSuggestedUsers)
	group.Get("/:username", getUser)
	group.Delete("/:username", deleteUser)
}

func searchUsers(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)

	query := c.Query("query")
	var users []*models.User
	database.DB.Where("username ILIKE ? OR name ILIKE ?", query+"%", query+"%").Find(&users)
	userList, err := serializers.UserListSerializer(users, request_user)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(userList)
}

func getSuggestedUsers(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)
	userList, err := serializers.UserListSerializer(request_user.FindSuggestedOnMutualFriends(), request_user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(userList)
}

func getUsers(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)
	if request_user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to perform this action",
		})
	}
	var users []*models.User
	database.DB.Find(&users)
	userList, err := serializers.UserListSerializer(users, request_user)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(userList)
}

func getUser(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)

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
	userList, err := serializers.UserCardSerializer(user, request_user)

	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(userList)
}

func deleteUser(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)

	c.Params("username")
	if request_user.Username != c.Params("username") && request_user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to delete this user",
		})
	}
	deleteUser := utils.GetUserByUsername(c.Params("username"))

	database.DB.Delete(&deleteUser)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "User deleted successfully",
	})
}

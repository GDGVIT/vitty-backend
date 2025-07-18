package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserHandler(api fiber.Router) {
	group := api.Group("/users")
	group.Use(middleware.JWTAuthMiddleware)
	group.Get("/", getUsers)
	group.Get("/search", searchUsers)
	group.Get("/suggested", getSuggestedUsers)
	group.Get("/emptyClassRooms", getEmptyClassRooms)
	group.Get("/:username", getUser)
	group.Patch("/campus", updateUserCampus)
	group.Delete("/:username", deleteUser)
}

func searchUsers(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)

	query := c.Query("query")
	var users []*models.User
	database.DB.Where("username ILIKE ? OR name ILIKE ?", query+"%", query+"%").Find(&users)
	return c.Status(fiber.StatusOK).JSON(serializers.UserListSerializer(users, request_user))
}

func getSuggestedUsers(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)
	return c.Status(fiber.StatusOK).JSON(serializers.UserListSerializer(request_user.FindSuggestedOnMutualFriends(), request_user))
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
	return c.Status(fiber.StatusOK).JSON(serializers.UserListSerializer(users, request_user))
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
		(request_user.Role == "admin") {
		return c.Status(fiber.StatusOK).JSON(serializers.UserSerializer(user, request_user))
	}

	if user.IsFriendsWith(request_user) {
		err, isGhosted := request_user.IsGhosted(username)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)

		}

		if isGhosted || errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":  "1811",
				"error": "",
			})
		}

		return c.Status(fiber.StatusOK).JSON(serializers.UserSerializer(user, request_user))
	}

	return c.Status(fiber.StatusOK).JSON(serializers.UserCardSerializer(user, request_user))
}

func deleteUser(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)

	username := c.Params("username")
	if request_user.Username != username && request_user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to delete this user",
		})
	}
	userToDelete := utils.GetUserByUsername(username)
	log.Println("User to delete:", userToDelete)

	if userToDelete.FirebaseUuid != "" {
		client, err := auth.FirebaseApp.Auth(context.Background())
		if err != nil {
			log.Println("Error connecting to Firebase:", err)
		} else {
			err = client.DeleteUser(context.Background(), userToDelete.FirebaseUuid)
			if err != nil {
				log.Println("Error deleting user from Firebase:", err)
			}
		}
	}

	if err := userToDelete.DeleteUser(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "Error deleting user: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "User deleted successfully",
	})
}

func getEmptyClassRooms(c *fiber.Ctx) error {
	filterSlot := strings.ToUpper(c.Query("slot"))

	if filterSlot == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "please mention the slot",
		})
	}

	file, err := os.Open("./data/freeClasses.json")
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "please contact vitty support",
		})
	}
	defer file.Close()

	var freeClasses map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&freeClasses)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
	}

	response := freeClasses[filterSlot]

	if response == nil {
		response = ""
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		filterSlot: response,
	})
}

func updateUserCampus(c *fiber.Ctx) error {
	type RequestBody struct {
		Campus *models.Campus `json:"campus"`
	}

	var body RequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	// Validate campus if provided
	if body.Campus != nil && !body.Campus.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "Invalid campus. Must be 'vellore' or 'chennai' or 'bhopal'",
		})
	}

	requestUser := c.Locals("user").(models.User)

	// Update the user's campus
	err = database.DB.Model(&models.User{}).Where("username = ?", requestUser.Username).Update("campus", body.Campus).Error
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "Failed to update campus",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Campus updated successfully",
		"campus": body.Campus,
	})
}

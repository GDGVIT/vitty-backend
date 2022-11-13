package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
)

func CheckFriends(user *models.User, request_user *models.User) bool {
	for _, friend := range user.Friends {
		if friend.ID == request_user.ID {
			return true
		}
	}
	return false
}

func CheckError(c *fiber.Ctx, err error) bool {
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"detail": err.Error(),
		})
		return true
	}
	return false
}

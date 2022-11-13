package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/database"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SendRequest(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	to_user_regno := c.Params("regno")

	if user.ID == to_user_regno {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You cannot send a friend request to yourself",
		})
	}

	to_user := models.User{}

	if database.DB.Where("id = ?", to_user_regno).First(&to_user).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}

	database.DB.Where("id = ?", to_user_regno).Preload(clause.Associations).First(&to_user)

	if utils.CheckFriends(&user, &to_user) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"detail": "Already friends",
		})
	}

	friend_request := models.FriendRequest{}

	if database.DB.Where("from_reg_no = ? AND to_reg_no = ?", user.ID, to_user.ID).First(&friend_request).RowsAffected != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"detail": "Request already sent",
		})
	}

	friend_request.From = user
	friend_request.To = to_user

	database.DB.Create(&friend_request)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"detail": "Request sent",
	})
}

func GetRequests(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	user_param := c.Params("regno")

	if user_param != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not allowed to view this",
		})
	}

	friend_requests := []models.FriendRequest{}

	database.DB.Where("to_reg_no = ?", user.ID).Preload(clause.Associations).Find(&friend_requests)

	return c.Status(fiber.StatusOK).JSON(utils.FriendRequestSerializerSlice(friend_requests))
}

func AcceptRequest(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	friend_request := models.FriendRequest{}
	database.DB.Where("id = ?", c.Params("id")).Preload(clause.Associations).First(&friend_request)

	if database.DB.Where("id = ?", c.Params("id")).First(&friend_request).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "Request not found",
		})
	}

	if friend_request.To.ID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not the receiver of this request",
		})
	}

	user.Friends = append(user.Friends, &friend_request.From)
	friend_request.From.Friends = append(friend_request.From.Friends, &user)

	database.DB.Save(&user)
	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
	database.DB.Save(&friend_request.From)
	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&friend_request.From)

	database.DB.Delete(&friend_request)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"detail": "Request accepted",
	})
}

func DeleteRequest(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	friend_request := models.FriendRequest{}

	if database.DB.Where("id = ?", c.Params("id")).First(&friend_request).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "Request not found",
		})
	}

	if friend_request.To.ID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not the receiver of this request",
		})
	}

	database.DB.Delete(&friend_request)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"detail": "Request deleted",
	})
}

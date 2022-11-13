package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/database"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUser(c *fiber.Ctx) error {
	request_user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	user := models.User{}
	database.DB.Where("id = ?", c.Params("regno")).Preload(clause.Associations).First(&user)

	if utils.CheckFriends(&user, &request_user) || request_user.ID == user.ID {
		return c.Status(fiber.StatusOK).JSON(utils.UserSerializer(&user))
	}

	return c.Status(fiber.StatusOK).JSON(utils.UserBlockSerializer(&user))
}

func UpdateUser(c *fiber.Ctx) error {
	request_user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	user := models.User{}
	database.DB.Where("id = ?", c.Params("regno")).First(&user)

	if request_user.ID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": "Unauthorized",
		})
	}

	user.FirstName = c.FormValue("first_name")
	user.LastName = c.FormValue("last_name")
	user.Email = c.FormValue("email")

	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"detail": "User updated",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	request_user, err := utils.GetUserFromHeader(c)
	user := models.User{}

	if utils.CheckError(c, err) {
		return nil
	}

	database.DB.Where("id = ?", c.Params("regno")).Preload(clause.Associations).First(&user)

	if request_user.ID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": "Unauthorized",
		})
	}

	database.DB.Delete(&user)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"detail": "User deleted",
	})
}

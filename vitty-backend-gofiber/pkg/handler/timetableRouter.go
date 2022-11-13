package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/database"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetTimetable(c *fiber.Ctx) error {
	request_user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	user := models.User{}
	database.DB.Where("id = ?", c.Params("regno")).Preload(clause.Associations).First(&user)

	timetable := models.Timetable{}
	database.DB.Where("user_reg_no = ?", user.ID).Preload(clause.Associations).First(&timetable)

	if request_user.ID == user.ID || utils.CheckFriends(&user, &request_user) {
		return c.Status(fiber.StatusOK).JSON(utils.TimetableSerializer(&timetable))
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"detail": "Unauthorized",
	})
}

func CreateTimetable(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	data := c.FormValue("request")
	timetable, err := utils.DetectTimetable(data)

	if utils.CheckError(c, err) {
		return nil
	}

	timetable.User = user

	// Check if user already has a timetable
	if database.DB.Where("user_reg_no = ?", user.ID).First(&timetable).RowsAffected != 0 {
		database.DB.Where("user_reg_no = ?", user.ID).Session(&gorm.Session{FullSaveAssociations: true}).Save(&timetable)
	} else {
		database.DB.Create(&timetable)
	}

	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&timetable)

	// Slots, err := utils.UnStringifySlots(timetable.Slots)

	if utils.CheckError(c, err) {
		return nil
	}

	return c.Status(fiber.StatusAccepted).JSON(utils.TimetableSerializer(&timetable))
}

func CreateTimetableV2(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if utils.CheckError(c, err) {
		return nil
	}

	data := c.FormValue("request")
	timetable, err := utils.DetectTimetableV2(data)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"detail": "Invalid Data",
		})
	}

	timetable.User = user

	// Check if user already has a timetable
	if database.DB.Where("user_reg_no = ?", user.ID).First(&timetable).RowsAffected != 0 {
		database.DB.Where("user_reg_no = ?", user.ID).Save(&timetable)
	} else {
		database.DB.Create(&timetable)
	}

	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&timetable)

	// Slots, err := utils.UnStringifySlots(timetable.Slots)

	if utils.CheckError(c, err) {
		return nil
	}

	return c.Status(fiber.StatusAccepted).JSON(utils.TimetableSerializer(&timetable))
}

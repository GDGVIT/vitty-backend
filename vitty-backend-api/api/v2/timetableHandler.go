package v2

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func timetableHandler(app fiber.Router) {
	group := app.Group("/timetable")
	group.Post("/parse", parseTimetable)
	group.Use(middleware.JWTAuthMiddleware)
	group.Post("/:username", createTimetable)
	group.Get("/:username", getTimetable)
	// group.Put("/:username", updateTimetable)
	group.Delete("/:username", deleteTimetable)
}

func parseTimetable(c *fiber.Ctx) error {
	// Get data from body
	var body struct {
		Timetable string `json:"timetable"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	var timetableV1 []utils.TimetableSlotV1
	timetableV1, err := utils.DetectTimetableV2(body.Timetable)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	slots := utils.SlotsV1ToSlotsV2(timetableV1)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"timetable": slots,
	})
}

func createTimetable(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)
	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}
	user := utils.GetUserByUsername(username)

	if user.Username != request_user.Username && request_user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to perform this action",
		})
	}

	// Get data from body
	var body struct {
		Timetable []models.Slot `json:"timetable"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	if !utils.CheckUserTimetableExists(user.Username) {
		var timetable models.Timetable
		timetable.User = user
		timetable.Slots = body.Timetable
		err := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&timetable).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
	} else {
		timetable := user.GetTimeTable()
		timetable.Slots = body.Timetable
		err := database.DB.Save(&timetable).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(serializers.TimetableSerializer(user.GetTimeTable()))
}

func getTimetable(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)

	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}
	user := utils.GetUserByUsername(username)

	if (user.Username != request_user.Username) &&
		(request_user.Role != "admin") &&
		(!user.IsFriendsWith(request_user)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to perform this action",
		})
	}

	timetable := user.GetTimeTable()
	return c.Status(fiber.StatusOK).JSON(serializers.TimetableSerializer(timetable))
}

func deleteTimetable(c *fiber.Ctx) error {
	request_user := c.Locals("user").(models.User)
	username := c.Params("username")
	if !utils.CheckUserExists(username) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}
	user := utils.GetUserByUsername(username)

	if user.Username != request_user.Username && request_user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "You are not authorized to perform this action",
		})
	}

	timetable := user.GetTimeTable()
	err := database.DB.Delete(&timetable).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Timetable deleted successfully",
	})
}

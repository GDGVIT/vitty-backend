package v2

import (
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/middleware"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
)

func reminderHandler(api fiber.Router) {
	group := api.Group("/reminders")
	group.Use(middleware.JWTAuthMiddleware)
	group.Get("/", getReminders)
	group.Post("/", createReminder)
	group.Patch("/", updateReminder)
	group.Delete("/:reminderId?", deleteReminder)
}

func getReminders(c *fiber.Ctx) error {
	var reminder models.Reminders

	username := c.Locals("user").(models.User).Username
	reminder.UserName = username

	err, reminders := reminder.GetReminders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "Reminders not fetched",
		})
	}

	if len(reminders) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"detail": "No reminders found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": serializers.RemindersSerializer(reminders),
	})
}

func createReminder(c *fiber.Ctx) error {
	var reminder models.Reminders

	if err := c.BodyParser(&reminder); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body ",
		})
	}

	if reminder.ReminderName == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reminder name is required",
		})
	}

	if reminder.ReminderContent == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reminder content is required",
		})
	}

	if !strings.Contains(reminder.ReminderId, "rem_") || len(reminder.ReminderId) < 32 {
		reminder.ReminderId = utils.UUIDWithPrefix("rem")
	}

	username := c.Locals("user").(models.User).Username
	reminder.UserName = username
	reminder.CreateReminder()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Reminder Saved Successfully",
	})
}

func updateReminder(c *fiber.Ctx) error {
	var reminder models.Reminders

	if err := c.BodyParser(&reminder); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body ",
		})
	}

	if reminder.ReminderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "Reminder id is required",
		})
	}

	err := reminder.UpdateReminder()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "Reminder not updated",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Reminder Saved Successfully",
	})
}

func deleteReminder(c *fiber.Ctx) error {
	var reminder models.Reminders

	reminderId := c.Params("reminderId")

	if reminderId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "Reminder id is missing",
		})
	}

	reminder.ReminderId = reminderId
	username := c.Locals("user").(models.User).Username
	reminder.UserName = username

	err := reminder.DeleteReminder()

	if err != nil {
		log.Info(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.ErrBadRequest)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Reminder deleted successfully",
	})
}

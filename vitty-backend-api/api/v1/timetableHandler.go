package v1

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func getTimetable(c *fiber.Ctx) error {
	response := struct {
		Data string `json:"data"`
	}{}

	if err := c.BodyParser(&response); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	slots, err := utils.DetectTimetable(response.Data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	if len(slots) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "No slots detected",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"slots": slots,
	})
}

func getTimetableV2(c *fiber.Ctx) error {
	response := struct {
		Data string `json:"data"`
	}{}

	if err := c.BodyParser(&response); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	slots, err := utils.DetectTimetableV2(response.Data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	if len(slots) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "No slots detected",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"slots": slots,
	})
}

func V1Handler(api fiber.Router) {
	group := api.Group("/v1")
	group.Post("/timetable", getTimetable)
	group.Post("/timetableV2", getTimetableV2)
}

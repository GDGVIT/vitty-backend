package v1

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func getTimetable(c *fiber.Ctx) error {
	data := c.FormValue("request")

	slots, err := utils.DetectTimetable(data)
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
	data := c.FormValue("request")

	slots, err := utils.DetectTimetableV2(data)
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
	group.Get("/timetable", getTimetable)
	group.Get("/timetable/v2", getTimetableV2)
}

package api

import (
	v1 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v1"
	v2 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v2"
	"github.com/gofiber/fiber/v2"
)

func InitializeApi(f *fiber.App) {
	// Root endpoint
	f.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to VITTY API!🎉")
	})

	// Ping endpoint
	f.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"detail": "pong",
			})
	})

	api := f.Group("/api")
	v1.V1Handler(api)
	v2.V2Handler(api)
}

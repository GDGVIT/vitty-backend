package api

import (
	v1 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v1"
	v2 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewWebApi() *fiber.App {
	// New fiber app
	fiberApp := fiber.New()

	fiberApp = fiber.New()
	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New(
		cors.Config{
			AllowOrigins:     "*",
			AllowHeaders:     "*",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,DELETE,PATCH,PUT,OPTIONS",
		},
	))

	// Root endpoint
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to VITTY API!🎉")
	})

	// Ping endpoint
	fiberApp.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"detail": "pong",
			})
	})

	api := fiberApp.Group("/api")
	v1.V1Handler(api)
	v2.V2Handler(api)

	return fiberApp
}

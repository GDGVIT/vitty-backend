package api

import (
	"time"

	v1 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v1"
	v2 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v2"
	v3 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewWebApi() *fiber.App {
	// New fiber app
	fiberApp := fiber.New()
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

	local, _ := time.LoadLocation("Asia/Kolkata")
	time.Local = local

	api := fiberApp.Group("/api")
	v1.V1Handler(api)
	v2.V2Handler(api)
	v3.V3Handler(api)

	return fiberApp
}

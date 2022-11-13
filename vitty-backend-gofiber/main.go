package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/database"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/handler"
)

func main() {
	database.Connect()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     "http://127.0.0.1:8000, http://0.0.0.0:8000, http://vittyapi.dscvit.com, https://vittyapi.dscvit.com, https://vitty.pages.dev, https://vitty.dscvit.com, http://vitty.dscvit.com,",
		AllowCredentials: true,
		AllowMethods:     "GET,POST",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Ok! Working!")
	})

	app.Get("/user/:regno/timetable", handler.GetTimetable)
	app.Post("/uploadtext/", handler.CreateTimetable)
	app.Post("/v2/uploadtext/", handler.CreateTimetableV2)

	app.Post("/signup/", handler.SignUp)
	app.Post("/login/", handler.Login)

	app.Get("/user/:regno", handler.GetUser)
	app.Patch("/user/:regno", handler.UpdateUser)
	app.Delete("/user/:regno", handler.DeleteUser)

	app.Get("/user/:regno/requests", handler.GetRequests)
	app.Post("/user/:regno/requests", handler.SendRequest)
	app.Post("/user/:regno/requests/:id", handler.AcceptRequest)
	app.Delete("/user/:regno/requests/:id", handler.DeleteRequest)

	app.Listen(os.Getenv("PORT"))
}

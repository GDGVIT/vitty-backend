package v2

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/admin/pkg"
	"github.com/gofiber/fiber/v2"
)

func V2Handler(api fiber.Router) {
	group := api.Group("/v2")
	authHandler(group)
	userHandler(group)
	timetableHandler(group)
	friendHandler(group)
	pkg.AdminHandler(group)
}

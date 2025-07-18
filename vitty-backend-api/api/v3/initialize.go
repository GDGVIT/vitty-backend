package v3

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/admin/pkg"
	v2 "github.com/GDGVIT/vitty-backend/vitty-backend-api/api/v2"
	"github.com/gofiber/fiber/v2"
)

func V3Handler(api fiber.Router) {
	group := api.Group("/v3")
	v2.AuthHandler(group)
	v2.UserHandler(group)
	v2.TimetableHandler(group)
	v2.FriendHandler(group)

	circleHandler(group)
	reminderHandler(group)
	noteHandler(group)
	pkg.AdminHandler(group)
}

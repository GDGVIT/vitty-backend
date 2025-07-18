package v2

import (
	"github.com/gofiber/fiber/v2"
)

func V2Handler(api fiber.Router) {
	group := api.Group("/v2")
	AuthHandler(group)
	UserHandler(group)
	TimetableHandler(group)
	FriendHandler(group)
}

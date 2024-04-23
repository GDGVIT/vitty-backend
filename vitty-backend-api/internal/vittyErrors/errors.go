package vittyerrors

import (
	"github.com/gofiber/fiber/v2"
)

var (
	ErrFetchCurrentStatus = fiber.NewError(fiber.StatusInternalServerError, "Getting User Current Status Failed")
)

package middleware

import "github.com/gofiber/fiber/v2"

// IsAdminMiddleware is a go-fiber middleware to check if the user is an admin
func IsAdminMiddleware(c *fiber.Ctx) error {
	// Get the team from the context
	role := c.Locals("role").(string)

	// If the team is not admin, return an error
	if role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": "You are not an admin",
		})
	}

	return c.Next()
}

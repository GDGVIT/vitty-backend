package v2

import (
	"context"
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

func googleLogin(c *fiber.Ctx) error {
	type RequestBody struct {
		Id_token string `json:"id_token"`
		RegNo    string `json:"reg_no,omitempty"`
		Username string `json:"username,omitempty"`
	}

	var body RequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	idtoken, err := idtoken.Validate(context.Background(), body.Id_token, auth.OauthConf.ClientID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	var user models.User
	user.Email = idtoken.Claims["email"].(string)
	if utils.CheckUserByEmail(user.Email) {
		err = database.DB.Where("email = ?", user.Email).First(&user).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
		if user.Picture != idtoken.Claims["picture"].(string) {
			database.DB.Model(&user).Update("picture", idtoken.Claims["picture"].(string))
		}
	} else {

		username := strings.ToLower(body.Username)
		if val, msg := utils.ValidateUsername(username); !val {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": msg,
			})
		}

		user.Name = idtoken.Claims["name"].(string)
		user.Picture = idtoken.Claims["picture"].(string)
		user.Username = username
		user.RegNo = body.RegNo

		err = database.DB.Create(&user).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
	}

	token, err := auth.CreateJWTToken(user.Username, user.Email, user.Role, auth.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(serializers.UserLoginSerializer(user, token))
}

func appleLogin(c *fiber.Ctx) error {
	type RequestBody struct {
		Id_token string `json:"id_token"`
		RegNo    string `json:"reg_no,omitempty"`
		Username string `json:"username,omitempty"`
	}

	var body RequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	return nil
}

func authHandler(api fiber.Router) {
	group := api.Group("/auth")
	group.Post("/google", googleLogin)
	group.Post("/apple", appleLogin)
}

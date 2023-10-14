package v2

import (
	"context"
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api/serializers"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

func authHandler(api fiber.Router) {
	group := api.Group("/auth")
	group.Post("/check-username", checkUsernameValidity)
	group.Post("/check-user-exists", checkUserExists)
	group.Post("/google", googleLogin)
	group.Post("/firebase", firebaseLogin)
}

func checkUsernameValidity(c *fiber.Ctx) error {
	type RequestBody struct {
		Username string `json:"username"`
	}
	var body RequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	username := strings.ToLower(body.Username)
	if val, msg := utils.ValidateUsername(username); !val {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": msg,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Username is valid",
	})
}

func checkUserExists(c *fiber.Ctx) error {
	type RequestBody struct {
		UUID string `json:"uuid"`
	}

	var body RequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	if !utils.CheckUserByUUID(body.UUID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "User does not exist",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "User exists",
	})
}

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

func firebaseLogin(c *fiber.Ctx) error {
	type RequestBody struct {
		UUID     string `json:"uuid"`
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

	client, err := auth.FirebaseApp.Auth(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}
	u_rec, err := client.GetUser(context.Background(), body.UUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	var user models.User
	user.FirebaseUuid = body.UUID
	if utils.CheckUserByUUID(user.FirebaseUuid) {
		err = database.DB.Where("firebase_uuid = ?", user.FirebaseUuid).First(&user).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": err.Error(),
			})
		}
		if user.Picture != u_rec.ProviderUserInfo[0].PhotoURL {
			database.DB.Model(&user).Update("picture", u_rec.ProviderUserInfo[0].PhotoURL)
		}
	} else {

		username := strings.ToLower(body.Username)
		if val, msg := utils.ValidateUsername(username); !val {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": msg,
			})
		}

		user.Name = u_rec.ProviderUserInfo[0].DisplayName
		user.Picture = u_rec.ProviderUserInfo[0].PhotoURL
		user.Email = u_rec.ProviderUserInfo[0].Email
		user.Name = u_rec.ProviderUserInfo[0].DisplayName
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

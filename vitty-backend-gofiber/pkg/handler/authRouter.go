package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/database"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
	"github.com/vitty-backend/vitty-backend-gofiber/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *fiber.Ctx) error {
	RegNo := c.FormValue("regno")
	FirstName := c.FormValue("first_name")
	LastName := c.FormValue("last_name")
	Email := c.FormValue("email")
	Password, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal("Error hashing password")
	}

	user := models.User{ID: RegNo, FirstName: FirstName, LastName: LastName, Email: Email, Password: string(Password)}

	if database.DB.Where("id = ?", user.ID).First(&user).RowsAffected == 1 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"detail": "User already exists",
		})
	}

	database.DB.Create(&user)
	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)

	token, err := utils.CreateToken(user.ID, user.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "Error creating token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(utils.UserLoginSerializer(&user, token))
}

func Login(c *fiber.Ctx) error {
	RegNo := c.FormValue("regno")
	Password := c.FormValue("password")

	user := models.User{}

	if database.DB.Where("id = ?", RegNo).First(&user).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"detail": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": "Incorrect password",
		})
	}

	token, err := utils.CreateToken(user.ID, user.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "Error creating token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.UserLoginSerializer(&user, token))
}

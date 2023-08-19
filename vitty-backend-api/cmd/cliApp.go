package cmd

import (
	"log"
	"os"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/cmd/management"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/urfave/cli/v2"
)

// VITTY CLI App
type VittyCliApp struct {
	env Env

	CliApp   *cli.App
	fiberApp *fiber.App
}

// Environment variables
type Env struct {
	// Fiber Variables
	fiberPort string
	debug     string

	// Database Variables
	postgresUrl string

	// Auth Variables
	jwtSecret string

	// Google Auth Variables
	google_client_id     string
	google_client_secret string
	google_redirect_uri  string
}

// Method to create a new VittyCliApp
func NewVittyCliApp() *VittyCliApp {
	var vittyCliApp VittyCliApp
	vittyCliApp.init()
	return &vittyCliApp
}

// Method to set environment variables
func (v *VittyCliApp) setEnv() {
	v.env.fiberPort = os.Getenv("FIBER_PORT")
	v.env.debug = os.Getenv("DEBUG")
	v.env.postgresUrl = os.Getenv("POSTGRES_URL")
	v.env.jwtSecret = os.Getenv("JWT_SECRET")
	v.env.google_client_id = os.Getenv("GOOGLE_CLIENT_ID")
	v.env.google_client_secret = os.Getenv("GOOGLE_CLIENT_SECRET")
	v.env.google_redirect_uri = os.Getenv("GOOGLE_REDIRECT_URI")
}

// Method to initialize the VittyCliApp
func (v *VittyCliApp) init() {
	v.setEnv()

	database.Connect(v.env.debug, v.env.postgresUrl)
	models.InitializeModels()
	auth.InitializeAuth(v.env.jwtSecret)
	auth.InitializeGoogleOauth(v.env.google_client_id, v.env.google_client_secret, v.env.google_redirect_uri)

	v.CliApp = cli.NewApp()

	// Set the name, usage and version of the app
	v.CliApp.Name = "Vitty"
	v.CliApp.Usage = "Vitty Backend API"
	v.CliApp.Version = "0.0.1"
	v.CliApp.Authors = []*cli.Author{
		{
			Name:  "Dhruv Shah",
			Email: "dhruvshahrds@gmail.com",
		}}
	v.CliApp.EnableBashCompletion = true

	v.fiberApp = fiber.New()
	v.fiberApp.Use(logger.New())
	v.fiberApp.Use(cors.New(
		cors.Config{
			AllowOrigins:     "*",
			AllowHeaders:     "Origin, Content-Type, Accept",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,DELETE,PATCH,PUT,OPTIONS",
		},
	))

	api.InitializeApi(v.fiberApp)

	runCommand := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run the server",
		Action: func(c *cli.Context) error {
			v.fiberApp.Listen(v.env.fiberPort)
			return nil
		},
	}

	v.CliApp.Commands = []*cli.Command{
		&runCommand,
		&management.CreateSuperuserCommand,
		&management.GetUsersCommand,
		&management.DeleteUserCommand,
	}
}

func (v *VittyCliApp) Run() {
	err := v.CliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package cmd

import (
	"log"
	"os"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/api"
	vittyCli "github.com/GDGVIT/vitty-backend/vitty-backend-api/cli"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/auth"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

// VITTY CLI App
type VittyApp struct {
	env Env

	cliApp *cli.App
	webApp *fiber.App
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

// Method to create a new VittyApp
func NewVittyCliApp() *VittyApp {
	var vittyApp VittyApp
	vittyApp.init()
	return &vittyApp
}

// Method to set environment variables
func (v *VittyApp) setEnv() {
	v.env.fiberPort = os.Getenv("FIBER_PORT")
	v.env.debug = os.Getenv("DEBUG")
	v.env.postgresUrl = os.Getenv("POSTGRES_URL")
	v.env.jwtSecret = os.Getenv("JWT_SECRET")
	v.env.google_client_id = os.Getenv("GOOGLE_CLIENT_ID")
	v.env.google_client_secret = os.Getenv("GOOGLE_CLIENT_SECRET")
	v.env.google_redirect_uri = os.Getenv("GOOGLE_REDIRECT_URI")
}

// Method to initialize CLI app
func (v *VittyApp) initCliApp() {
	v.cliApp = vittyCli.NewCliApp()

	// Adding Run command
	v.cliApp.Commands = append(v.cliApp.Commands, &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run the server",
		Action: func(c *cli.Context) error {
			v.webApp.Listen(v.env.fiberPort)
			return nil
		},
	})
}

// Method to initialize Web app
func (v *VittyApp) initWebApp() {
	v.webApp = api.NewWebApi()
}

// Method to initialize the VittyApp
func (v *VittyApp) init() {
	// Set environment variables
	v.setEnv()

	// Connect to database
	database.Connect(v.env.debug, v.env.postgresUrl)

	// Initialize models
	models.InitializeModels()

	// Initialize auth
	auth.InitializeAuth(v.env.jwtSecret)
	auth.InitializeGoogleOauth(v.env.google_client_id, v.env.google_client_secret, v.env.google_redirect_uri)
	auth.InitializeFirebaseApp()

	// Initialize Web app
	v.initWebApp()

	// Initialize CLI app
	v.initCliApp()
}

func (v *VittyApp) Run() {
	err := v.cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

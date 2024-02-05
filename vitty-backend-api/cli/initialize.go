package vittyCli

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/cli/commands"
	"github.com/urfave/cli/v2"
)

func NewCliApp() *cli.App {
	// New cli app
	cliApp := cli.NewApp()

	// Set the name, usage and version of the app
	cliApp.Name = "Vitty"
	cliApp.Usage = "Vitty Backend API"
	cliApp.Version = "0.0.1"
	cliApp.Authors = []*cli.Author{
		{
			Name:  "Dhruv Shah",
			Email: "dhruvshahrds@gmail.com",
		}}
	cliApp.EnableBashCompletion = true
	cliApp.Description = "CLI for Vitty Backend API"

	// Set the commands
	commands.AddCommands(cliApp)

	return cliApp
}

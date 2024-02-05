package commands

import "github.com/urfave/cli/v2"

// AddCommands adds the commands to the app
func AddCommands(cliApp *cli.App) {
	// Add the commands to the app
	cliApp.Commands = append(cliApp.Commands, userCommands...)
	cliApp.Commands = append(cliApp.Commands, TimetableCommands...)
}

package management

import (
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/urfave/cli/v2"
)

func getUsers(c *cli.Context) error {
	// Take arguements as name of filter and value of filter
	// If no arguements are given, return all users
	// If arguements are given, return users that match the filter

	filter := c.String("filter")
	value := c.String("value")

	var users []models.User
	filter_query := fmt.Sprintf("%s = ?", filter)

	if filter == "" || value == "" {
		database.DB.Find(&users)
	} else {
		database.DB.Where(filter_query, value).Find(&users)
	}

	// Display Users in a table
	fmt.Println("ID\tUsername\t\tRole\tRegNo\tEmail\tPicture")
	for _, user := range users {
		fmt.Printf("\t%s\t\t%s\t\t%s\t\t%s\t\t%s\n", user.Username, user.Role, user.RegNo, user.Email, user.Picture)
	}
	return nil
}

var GetUsersCommand = cli.Command{
	Name:    "getusers",
	Aliases: []string{"gu"},
	Usage:   "Get users",
	Action:  getUsers,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "filter",
			Aliases:  []string{"f"},
			Usage:    "Filter users by",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "value",
			Aliases:  []string{"v"},
			Usage:    "Value of filter",
			Required: false,
		},
	},
}

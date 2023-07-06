package management

import (
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/urfave/cli/v2"
)

func deleteUser(c *cli.Context) error {
	key := c.String("key")
	value := c.String("value")

	if key == "" || value == "" {
		fmt.Println("Enter key: ")
		fmt.Scanln(&key)
		fmt.Println("Enter value: ")
		fmt.Scanln(&value)
	}

	var users models.User
	filter_query := fmt.Sprintf("%s = ?", key)
	err := database.DB.Where(filter_query, value).Delete(&users).Error
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	return nil
}

var DeleteUserCommand = cli.Command{
	Name:    "deleteuser",
	Aliases: []string{"du"},
	Usage:   "Delete user",
	Action:  deleteUser,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "key",
			Aliases:  []string{"k"},
			Usage:    "Key of user to delete",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "value",
			Aliases:  []string{"v"},
			Usage:    "Value of key",
			Required: false,
		},
	},
}

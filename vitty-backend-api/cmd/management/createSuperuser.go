package management

import (
	"errors"
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/urfave/cli/v2"
)

func createSuperuser(c *cli.Context) error {
	fmt.Println("Enter username: ")
	var username string
	fmt.Scanln(&username)

	if !utils.CheckUserExists(username) {
		return errors.New("user does not exist")
	}
	user := utils.GetUserByUsername(username)
	user.Role = "superuser"
	return database.DB.Save(&user).Error
}

var CreateSuperuserCommand = cli.Command{
	Name:    "createsuperuser",
	Aliases: []string{"csu"},
	Usage:   "Create a superuser",
	Action:  createSuperuser,
}

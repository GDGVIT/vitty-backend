package commands

import (
	"errors"
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/urfave/cli/v2"
)

var userCommands = []*cli.Command{
	{
		Name:    "createsuperuser",
		Aliases: []string{"csu"},
		Usage:   "Create a superuser",
		Action:  createSuperuser,
	},
	{
		Name:    "createadminuser",
		Aliases: []string{"cau"},
		Usage:   "Create an admin user",
		Action:  createAdminuser,
	},
	{
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
	},
	{
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
	},
}

func createSuperuser(c *cli.Context) error {
	fmt.Println("Enter username: ")
	var username string
	fmt.Scanln(&username)

	if !utils.CheckUserExists(username) {
		return errors.New("user does not exist")
	}
	user := utils.GetUserByUsername(username)
	user.Role = "superuser"
	err := database.DB.Save(&user).Error
	if err != nil {
		return err
	}
	fmt.Println("Superuser created successfully")
	return nil
}

func createAdminuser(c *cli.Context) error {
	fmt.Println("Enter username: ")
	var username string
	fmt.Scanln(&username)

	if !utils.CheckUserExists(username) {
		return errors.New("user does not exist")
	}
	user := utils.GetUserByUsername(username)
	user.Role = "admin"
	err := database.DB.Save(&user).Error
	if err != nil {
		return err
	}
	fmt.Println("Admin user created successfully")
	return nil
}

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
	fmt.Println("User deleted successfully")
	return nil
}

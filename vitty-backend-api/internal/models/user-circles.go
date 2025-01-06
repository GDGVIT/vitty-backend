package models

import (
	"log"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm"
)

type UsersCirclesJoin struct {
	CID        string `gorm:"primaryKey"`
	Uname      string `gorm:"primaryKey"`
	CircleRole string
	Circles    Circles `gorm:"foreignKey:CID;references:CircleId;constraint:OnDelete:CASCADE"`
	User       User    `gorm:"foreignKey:Uname;references:Username"`
}

func (ucj *UsersCirclesJoin) AddUserToCircle() error {
	err := database.DB.Create(ucj).Error
	if err != nil {
		return err
	}

	var user User
	user.Username = ucj.Uname

	err = ucj.Circles.ComputeCircleSlots(user, false)

	return err
}

func (ucj *UsersCirclesJoin) UpdateUserCircleRole(role string) error {
	err := database.DB.Model(&UsersCirclesJoin{}).Where(&ucj).Update("circle_role", role).Error
	return err
}

func (ucj *UsersCirclesJoin) DeleteUserFromCircle() error {

	err := database.DB.Delete(ucj).Error

	if err != nil {
		return err
	}

	ucj.Uname = ""
	err, users := ucj.GetUsersofCircle()

	if err != nil {
		return err
	}

	if len(users) == 0 {
		var circle Circles
		circle.CircleId = ucj.CID

		err := circle.DeleteCircle()
		return err
	} else {
		for _, user := range users {
			ucj.Circles.CircleSlots = ""
			ucj.Circles.CircleId = ucj.CID
			err = ucj.Circles.ComputeCircleSlots(user, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ucj *UsersCirclesJoin) GetCircleofUserByCID() (error, UsersCirclesJoin) {
	var circleUser UsersCirclesJoin

	result := database.DB.Where(&ucj).Preload("Circles").Find(&circleUser)
	log.Println(result.Error)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound, circleUser
	}
	return result.Error, circleUser
}

func (ucj *UsersCirclesJoin) GetCirclesofUser() (error, []UsersCirclesJoin) {
	var circleUsers []UsersCirclesJoin

	err := database.DB.Where(&ucj).Preload("Circles").Find(&circleUsers).Error
	if err != nil {
		return err, nil
	}

	return err, circleUsers
}

func (ucj *UsersCirclesJoin) GetUsersofCircle() (error, []User) {
	var users []User
	var circleUsers []UsersCirclesJoin

	err := database.DB.Where(&ucj).Preload("User").Find(&circleUsers).Error
	if err != nil {
		return err, nil
	}

	for _, circleUser := range circleUsers {
		users = append(users, circleUser.User)
	}

	return err, users
}

func (ucj *UsersCirclesJoin) IsUserCircleAdmin() (error, bool) {

	err, userCircle := ucj.GetCircleofUserByCID()

	if err != nil {
		return err, false
	}

	if userCircle.CircleRole == "admin" {
		return nil, true
	}

	return nil, false
}

func (ucj *UsersCirclesJoin) IsUserOfCircle() (error, bool) {

	err, userCircle := ucj.GetCircleofUserByCID()
	if err != nil {
		return err, false
	}

	if userCircle.Uname == ucj.Uname {
		return nil, true
	}

	return nil, false
}

package models

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm"
)

type CircleRequest struct {
	FromUsername string  `gorm:"primaryKey"`
	ToUsername   string  `gorm:"primaryKey"`
	CID          string  `gorm:"primaryKey"`
	From         User    `gorm:"constraint:OnDelete:CASCADE;foreignKey:FromUsername;references:Username"`
	To           User    `gorm:"constraint:OnDelete:CASCADE;foreignKey:ToUsername;references:Username"`
	Circles      Circles `gorm:"constraint:OnDelete:CASCADE;foreignKey:CID;references:CircleId;constraint:OnDelete:CASCADE"`
}

func (cr *CircleRequest) CreateRequest() error {
	err := database.DB.Create(&cr).Error
	return err
}

func (cr *CircleRequest) AcceptRequest() error {
	var userCircle UsersCirclesJoin

	userCircle.CID = cr.CID
	userCircle.CircleRole = "member"
	userCircle.Uname = cr.ToUsername

	err, circleRequest := cr.GetRequestByCircleId()

	if circleRequest.ToUsername == "" {
		return gorm.ErrRecordNotFound
	} else if err != nil {
		return err
	}

	err = userCircle.AddUserToCircle()

	if err != nil {
		return err
	}

	err = cr.DeleteRequest()

	if err != nil {
		return err
	}

	return nil

}

func (cr *CircleRequest) DeclineRequest() error {
	result := database.DB.Where(cr).Delete(&CircleRequest{})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (cr *CircleRequest) GetRequestByCircleId() (error, CircleRequest) {
	var circleRequest CircleRequest
	err := database.DB.Where(cr).Find(&circleRequest).Error
	return err, circleRequest
}

func (cr *CircleRequest) GetReceivedRequests(username string) (error, []CircleRequest) {
	var circleRequests []CircleRequest
	err := database.DB.Where("to_username = ?", username).Preload("Circles").Find(&circleRequests).Error
	return err, circleRequests
}

func (cr *CircleRequest) GetSentRequests(username string) (error, []CircleRequest) {
	var circleRequests []CircleRequest
	err := database.DB.Where("from_username = ?", username).Preload("Circles").Find(&circleRequests).Error
	return err, circleRequests
}

func (cr *CircleRequest) DeleteRequest() error {
	result := database.DB.Where(&cr).Delete(&CircleRequest{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

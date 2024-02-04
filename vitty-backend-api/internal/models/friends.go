package models

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm"
)

type FriendRequest struct {
	FromUsername string `gorm:"primaryKey"`
	ToUsername   string `gorm:"primaryKey"`
	From         User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:FromUsername;references:Username"`
	To           User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ToUsername;references:Username"`
}

func (fr *FriendRequest) Accept() {
	fr.From.Friends = append(fr.From.Friends, &fr.To)
	fr.To.Friends = append(fr.To.Friends, &fr.From)
	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&fr.From)
	database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&fr.To)
	database.DB.Delete(&fr)
}

func (fr *FriendRequest) Decline() {
	database.DB.Delete(&fr)
}

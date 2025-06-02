package models

import (
	"time"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Reminders struct {
	ReminderId      string                `json:"reminder_id,omitempty" gorm:"unique"`
	ReminderName    *string               `json:"reminder_name" gorm:"primaryKey"`
	UserName        string                `json:"user_name" gorm:"primaryKey;foreignKey:Username;constraint:OnDelete:CASCADE"`
	ReminderContent *string               `json:"reminder_content"`
	ReminderTime    *time.Time            `json:"reminder_time,omitempty"`
	User            User                  `gorm:"foreignKey:UserName;references:Username;constraint:OnDelete:CASCADE"`
	DeletedAt       soft_delete.DeletedAt `json:"-,omitempty" gorm:"index"`
}

func (r *Reminders) CreateReminder() error {
	err := database.DB.Create(&r).Error
	return err
}

func (r *Reminders) GetReminders() (error, []Reminders) {
	var reminders []Reminders
	err := database.DB.Where(&r).Find(&reminders).Error
	return err, reminders
}

func (r *Reminders) UpdateReminder() error {
	updateInteface := make(map[string]interface{})

	if r.ReminderName != nil {
		updateInteface["reminder_name"] = *r.ReminderName
	}

	if r.ReminderContent != nil {
		updateInteface["reminder_content"] = *r.ReminderContent
	}

	if r.ReminderTime != nil {
		updateInteface["reminder_time"] = *r.ReminderTime
	}

	err := database.DB.Model(&Reminders{}).Where("reminder_id = ?", r.ReminderId).UpdateColumns(updateInteface).Error

	return err
}

func (r *Reminders) DeleteReminder() error {
	result := database.DB.Where("reminder_id = ? AND username like ?", r.ReminderId, r.UserName).Delete(&r)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (r *Reminders) SoftDeleteExpiredReminders() {
	now := time.Now()
	database.DB.Model(&Reminders{}).Where("reminder_time < ?", now).Delete(&Reminders{})
}

func (r *Reminders) CleanupOldReminders() {
	database.DB.Unscoped().Where("deleted_at < NOW() - INTERVAL '7 days'").Delete(&Reminders{})
}

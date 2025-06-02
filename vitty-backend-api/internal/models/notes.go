package models

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm"
)

type Notes struct {
	NoteID      string  `json:"note_id,omitempty" gorm:"unique"`
	NoteName    string  `json:"note_name" gorm:"primaryKey"`
	UserName    string  `json:"user_name" gorm:"primaryKey"`
	CourseID    string  `json:"course_id"`
	NoteContent string  `json:"note_content"`
	User        User    `gorm:"foreignKey:UserName;references:Username;constraint:OnDelete:CASCADE"`
	Courses     Courses ` gorm:"foreignKey:CourseID;references:CourseId;constraint:OnDelete:CASCADE"`
}

func (n *Notes) SaveNote() error {
	result := database.DB.Save(n)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (n *Notes) DeleteNote() error {
	result := database.DB.Where("user_name = ?", n.UserName).Delete(&n)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (n *Notes) GetNotesByCourseId() (error, []Notes) {
	var notes []Notes
	err := database.DB.Where(n).Preload("Courses").Find(&notes).Error
	return err, notes
}

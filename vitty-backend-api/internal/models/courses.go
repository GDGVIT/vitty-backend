package models

type Courses struct {
	CourseId   string `gorm:"primaryKey"`
	CourseName string
}

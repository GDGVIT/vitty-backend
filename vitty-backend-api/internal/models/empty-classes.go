package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ClassSlotsJoin struct {
	Class string         `gorm:"column:class;primarykey"`
	Slots datatypes.JSON `gorm:"column:slots"`
}

func (ec *ClassSlotsJoin) InsertEmptyClassRooms(class string, newSlot string) error {

	var classSlots ClassSlotsJoin
	result := database.DB.Where("class = ?", class).First(&classSlots)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Println("Class not found: ", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		slot := []string{newSlot}
		slotJson, _ := json.Marshal(slot)
		newClassSlot := ClassSlotsJoin{
			Class: class,
			Slots: datatypes.JSON(slotJson),
		}
		err := database.DB.Create(newClassSlot).Error

		if err != nil {
			log.Println("Failed to add new class: ", err)
			return fiber.ErrInternalServerError
		}

	} else {
		var existingSlots []string

		err := json.Unmarshal(classSlots.Slots, &existingSlots)
		if err != nil {
			log.Println("failed to unmarshal existing students: ", err)
			return fiber.ErrInternalServerError
		}

		slotExists := false

		for _, slot := range existingSlots {
			if slot == newSlot {
				slotExists = true
				break
			}
		}

		if !slotExists {
			// Append the new student to the array
			existingSlots = append(existingSlots, newSlot)
			updatedStudentsJSON, _ := json.Marshal(existingSlots)
			classSlots.Slots = datatypes.JSON(updatedStudentsJSON)
			err := database.DB.Save(&classSlots).Error

			if err != nil {
				log.Println("Failed to update array existing array: ", err)
				return fiber.ErrInternalServerError
			}
		}
	}

	return nil
}

func (ec *ClassSlotsJoin) FindEmptyClassRooms(slot string) []string {

	var freeClasses []string
	total := 0
	offset := 0
	limit := 1000

	query := database.DB.
		Model(&ClassSlotsJoin{}).
		Where("NOT (slots @> '[\"?\"]')", database.DB.Raw(slot)).
		Where("slots::text !~ '\\[\"L.*\"\\]'")

	err := query.
		Select("COUNT(class)").
		Scan(&total).Error

	if err != nil {
		fmt.Println("Error: ", err)
		return freeClasses
	}

	for total >= 0 {
		err := query.
			Select("class").
			Limit(limit).
			Offset(offset).
			Find(&freeClasses).Error

		if err != nil {
			fmt.Println("Error: ", err)
			return freeClasses
		}

		total -= limit
		offset += limit
	}

	return freeClasses
}

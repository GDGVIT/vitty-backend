package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
)

type TimetableSlotV1 struct {
	Slot           string `json:"Slot"`
	CourseName     string `json:"Course_Name"`
	CourseFullName string `json:"Course_Full_Name,omitempty"`
	CourseType     string `json:"Course_type"`
	Venue          string `json:"Venue"`
}

// Function to check if any array is empty
func isStrArrEmpty(arr []string) bool {
	return len(arr) == 0
}

func DetectTimetable(text string) ([]TimetableSlotV1, error) {
	re := regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}[\D]{1}[A-Z]{3,4}[0-9]{3,4}[A-Z]{0,1}[\D]{1}[A-Z]{2,3}[\D]{1}[A-Z]{2,6}[0-9]{2,4}[A-Za-z]{0,1}[\D]{1}[A-Z]{2,4}[0-9]{0,3}`)
	slots := re.FindAllString(text, -1)
	var Slots []TimetableSlotV1

	for _, slot := range slots {
		var obj TimetableSlotV1

		obj.Slot = regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}\b`).FindAllString(slot, -1)[0]
		obj.CourseName = regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}\b`).FindAllString(slot, -1)[0]
		course_type := regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}\b`).FindAllString(slot, -1)[0]

		var c_type string
		if course_type == "ELA" || course_type == "LO" {
			c_type = "Lab"
		} else {
			c_type = "Theory"
		}
		obj.CourseType = c_type
		obj.Venue = regexp.MustCompile(`[A-Z]{2,6}[0-9]{2,4}[A-Za-z]{0,1}\b`).FindAllString(slot, -1)[1]

		Slots = append(Slots, obj)
	}

	if len(Slots) == 0 {
		return Slots, nil
	}

	var err error

	if err != nil {
		return Slots, errors.New("error in detecting timetable")
	}

	return Slots, nil
}

func DetectTimetableV2(text string) ([]TimetableSlotV1, error) {
	fmt.Println("Detecting timetable...")
	text = strings.ReplaceAll(text, "\r", "")
	var Slots []TimetableSlotV1

	rows := regexp.MustCompile("(?s)[A-Z]{4}[0-9]{3}.+?Registered|[A-Z]{3}[0-9]{4}.+?Registered").FindAllString(text, -1)
	re_code_n_name := regexp.MustCompile("[A-Z]{4}[0-9]{3}.+\n|[A-Z]{3}[0-9]{4}.+\n")
	re_code := regexp.MustCompile("[A-Z]{4}[0-9]{3}[LPEM]|[A-Z]{3}[0-9]{4}[LPEM]")
	re_venue := regexp.MustCompile("\n{3}[A-Z]+[0-9]{1,3}.+\n|\n{3}NIL\n")
	re_slots := regexp.MustCompile(".+[1-9].+[-]\n|NIL.+[-]\n")

	for _, row := range rows {
		if isStrArrEmpty(re_code_n_name.FindAllString(row, -1)) {
			continue
		}
		code_n_name := re_code_n_name.FindAllString(row, -1)[0]
		code := re_code.FindAllString(code_n_name, -1)
		if isStrArrEmpty(code) {
			continue
		}
		code_n := code[0]
		name_n := strings.TrimRight(strings.TrimLeft(code_n_name, code_n)[3:], "\n")
		if code_n == "" && name_n == "" {
			continue
		}

		if isStrArrEmpty(re_venue.FindAllString(row, -1)) {
			continue
		}
		venue := re_venue.FindAllString(row, -1)[0]
		venue = venue[3 : len(venue)-1]

		if isStrArrEmpty(re_slots.FindAllString(row, -1)) {
			continue
		}
		slotStr := re_slots.FindAllString(row, -1)[0]
		slots := strings.Split(slotStr[:len(slotStr)-3], "+")

		for _, slot := range slots {
			var obj TimetableSlotV1
			obj.Slot = slot
			obj.CourseName = code_n
			obj.CourseFullName = name_n
			obj.Venue = venue
			if slot[0:1] == "L" {
				obj.CourseType = "Lab"
			} else {
				obj.CourseType = "Theory"
			}

			Slots = append(Slots, obj)
		}
	}

	if len(Slots) == 0 {
		goto throwerror
	}
	return Slots, nil

throwerror:
	return DetectTimetable(text)
}

func SlotsV1ToSlotsV2(slots []TimetableSlotV1) []models.Slot {
	var timetableSlots []models.Slot
	for _, slot := range slots {
		if slot.CourseFullName == "" {
			slot.CourseFullName = GetCourseFullNameIfExists(slot.CourseName)
		}
		slotV2 := models.Slot{
			Slot:  slot.Slot,
			Name:  slot.CourseFullName,
			Code:  slot.CourseName,
			Type:  slot.CourseType,
			Venue: slot.Venue,
		}
		slotV2.AddSlotTime()
		timetableSlots = append(timetableSlots, slotV2)
	}
	return timetableSlots
}

func CheckUserTimetableExists(username string) bool {
	var count int64
	database.DB.Model(&models.Timetable{}).Where("user_username = ?", username).Count(&count)
	return count != 0
}

func GetCourseFullNameIfExists(courseCode string) string {
	var slot models.Slot
	database.DB.Where("code = ?", courseCode).First(&slot)
	// If course name is not found, return the course code
	if slot.Name == "" {
		return courseCode
	}
	return slot.Name
}

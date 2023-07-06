package utils

import (
	"errors"
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
	text = strings.ReplaceAll(text, "\r", "")

	var code, name, venue []string
	var slots [][]string
	var Slots []TimetableSlotV1

	re := regexp.MustCompile("[A-Z]{4}[0-9]{3}.+\n|[A-Z]{3}[0-9]{4}.+\n").FindAllString(text, -1)

	for _, c := range re {
		code_n := regexp.MustCompile("[A-Z]{4}[0-9]{3}|[A-Z]{3}[0-9]{4}").FindAllString(c, -1)[0]
		name_n := strings.TrimRight(strings.TrimLeft(c, code_n), "\n")
		if code_n != "" && name_n != "" {
			code = append(code, code_n)
			name = append(name, name_n)
		}
	}

	re = regexp.MustCompile("\n{3}[A-Z]+[0-9]{1,3}.+\n|\n{3}NIL\n").FindAllString(text, -1)
	for _, c := range re {
		venue = append(venue, c[3:len(c)-1])
	}

	re = regexp.MustCompile(".+[1-9].+[-]\n|NIL.+[-]\n").FindAllString(text, -1)
	for _, c := range re {
		slots = append(slots, strings.Split(c[:len(c)-3], "+"))
	}

	if len(slots) == len(name) && len(slots) == len(code) && len(slots) == len(venue) {

		for i := 0; i < len(slots); i++ {
			if slots[i][0] != "NIL" {
				for _, slot := range slots[i] {
					var obj TimetableSlotV1
					obj.Slot = slot
					obj.CourseName = code[i]
					obj.CourseFullName = name[i]
					obj.Venue = venue[i]
					if slot[0:1] == "L" {
						obj.CourseType = "Lab"
					} else {
						obj.CourseType = "Theory"
					}

					Slots = append(Slots, obj)
				}
			}
		}

		if len(Slots) == 0 {
			goto throwerror
		}

		var err error

		if err != nil {
			return Slots, errors.New("error in detecting timetable")
		}
		return Slots, nil

	} else {
		goto throwerror
	}

throwerror:
	return DetectTimetable(text)
}

func CheckUserTimetableExists(username string) bool {
	var count int64
	database.DB.Model(&models.Timetable{}).Where("user_username = ?", username).Count(&count)
	return count != 0
}

package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
)

func DetectTimetable(text string) (models.Timetable, error) {
	re := regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}[\D]{1}[A-Z]{3,4}[0-9]{3,4}[A-Z]{0,1}[\D]{1}[A-Z]{2,3}[\D]{1}[A-Z]{2,6}[0-9]{2,4}[A-Za-z]{0,1}[\D]{1}[A-Z]{2,4}[0-9]{0,3}`)
	slots := re.FindAllString(text, -1)
	var Slots []models.Slot
	var timetable models.Timetable

	for _, slot := range slots {
		var obj models.Slot

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
		return timetable, nil
	}

	var err error
	// timetable.Slots, err = StringifySlots(Slots)
	timetable.Slots = Slots

	if err != nil {
		return timetable, errors.New("error in detecting timetable")
	}

	return timetable, nil
}

func DetectTimetableV2(text string) (models.Timetable, error) {
	text = strings.ReplaceAll(text, "\r", "")

	var code, name, venue []string
	var slots [][]string
	var Slots []models.Slot
	var timetable models.Timetable

	re := regexp.MustCompile("[A-Z]{4}[0-9]{3}.+\n|[A-Z]{3}[0-9]{4}.+\n").FindAllString(text, -1)

	for _, c := range re {
		split := strings.Split(c[:len(c)-1], " - ")

		if len(split) == 2 {
			code = append(code, split[0])
			name = append(name, split[1])
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
					var obj models.Slot
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
		// timetable.Slots, err = StringifySlots(Slots)
		timetable.Slots = Slots

		if err != nil {
			return timetable, errors.New("error in detecting timetable")
		}

		return timetable, nil

	} else {
		goto throwerror
	}

throwerror:
	return DetectTimetable(text)
}

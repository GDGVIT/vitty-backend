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

func DetectTimetable(text string) ([]TimetableSlotV1, error) {
	re := regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}[\D]{1}[A-Z]{3,4}[0-9]{3,4}[A-Z]{0,1}[\D]{1}[A-Z]{2,3}[\D]{1}[A-Z]{2,6}[0-9]{2,4}[A-Za-z]{0,1}[\D]{1}[A-Z]{2,4}[0-9]{0,3}`)
	slots := re.FindAllString(text, -1)
	var Slots []TimetableSlotV1

	for _, slot := range slots {
		var obj TimetableSlotV1

		obj.Slot = regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}\b`).FindAllString(slot, -1)[0]
		obj.CourseName = regexp.MustCompile(`[A-Z]{3,4}[0-9]{3,4}[A-Z]{0,1}\b`).FindAllString(slot, -1)[0]
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
		return Slots, errors.New("error in detecting timetable")
	}

	return Slots, nil
}

func DetectTimetableV2(text string, campus string) ([]TimetableSlotV1, error) {
	fmt.Println("Detecting timetable...")
	text = strings.ReplaceAll(text, "\r", "")
	var Slots []TimetableSlotV1

	switch strings.ToLower(campus) {
	case "chennai":
		if slots := parseChennaiFormat(text); len(slots) > 0 {
			return slots, nil
		}
	case "bhopal":
		if slots := parseBhopalFormat(text); len(slots) > 0 {
			return slots, nil
		}
	default:

		// Split the text into individual course entries based on numbered entries
		rows := regexp.MustCompile(`(?s)\n\d+\n.*?Registered(?:\s+and\s+Approved)?`).FindAllString(text, -1)

		// If no matches found, alternative
		if len(rows) == 0 {
			// splitting by course code
			rows = regexp.MustCompile(`(?s)[A-Z]{3,4}[0-9]{3,4}[LPEMJ]?\s*-\s*[^\n]+.*?Registered(?:\s+and\s+Approved)?`).FindAllString(text, -1)
		}

		re_code_n_name := regexp.MustCompile(`([A-Z]{3,4}[0-9]{3,4}[LPEMJ]?)\s*-\s*([^\n(]+)`)
		// captures venue names with optional spacing
		re_venue := regexp.MustCompile(`\n\s*([A-Z]+[0-9]{1,4}[A-Za-z]?|NIL)\s*(?:\n|$)`)
		re_slots := regexp.MustCompile(`\n([A-Z]*[0-9]*[A-Z]{1,3}[0-9]{1,2}(?:\+[TA]*[A-Z]{1,3}[0-9]{1,2})*|NIL)\s*-\s*\n`)

		for _, row := range rows {
			// course code and name
			codeNameMatches := re_code_n_name.FindStringSubmatch(row)
			if len(codeNameMatches) < 3 {
				continue
			}

			code_n := codeNameMatches[1]
			name_n := strings.TrimSpace(codeNameMatches[2])

			if code_n == "" && name_n == "" {
				continue
			}

			// venue
			venueMatches := re_venue.FindStringSubmatch(row)
			if len(venueMatches) < 2 {
				continue
			}
			venue := strings.TrimSpace(venueMatches[1])

			// slots
			slotMatches := re_slots.FindStringSubmatch(row)
			if len(slotMatches) < 2 {
				continue
			}
			slotStr := strings.TrimSpace(slotMatches[1])
			slots := strings.Split(slotStr, "+")

			for _, slot := range slots {
				slot = strings.TrimSpace(slot)
				if slot == "" || slot == "NIL" {
					continue
				}

				var obj TimetableSlotV1
				obj.Slot = slot
				obj.CourseName = code_n
				obj.CourseFullName = name_n
				obj.Venue = venue
				if len(slot) > 0 && slot[0:1] == "L" {
					obj.CourseType = "Lab"
				} else {
					obj.CourseType = "Theory"
				}

				Slots = append(Slots, obj)
			}
		}

		if len(Slots) == 0 {
			return DetectTimetable(text)
		}
	}

	return Slots, nil
}

func parseBhopalFormat(text string) []TimetableSlotV1 {
	var Slots []TimetableSlotV1

	rows := regexp.MustCompile(`(?s)\n\d+\n.*?Registered(?:\s+and\s+Approved)?`).FindAllString(text, -1)

	for _, row := range rows {
		codeNameMatches := regexp.MustCompile(`([A-Z]{3,4}[0-9]{3,4}[LPEMJ]?)\s*-\s*([^\n(]+)`).FindStringSubmatch(row)
		if len(codeNameMatches) < 3 {
			continue
		}

		courseCode := strings.TrimSpace(codeNameMatches[1])
		courseName := strings.TrimSpace(codeNameMatches[2])

		slotVenueMatches := regexp.MustCompile(`\n([A-Z][0-9]{2}(?:\+[A-Z][0-9]{2})*)\s*-\s*\n\s*([A-Z0-9-]+)`).FindStringSubmatch(row)
		if len(slotVenueMatches) < 3 {
			continue
		}

		slotStr := strings.TrimSpace(slotVenueMatches[1])
		venue := strings.TrimSpace(slotVenueMatches[2])

		if slotStr == "NIL" || venue == "NIL" {
			continue
		}

		slots := strings.Split(slotStr, "+")

		for _, slot := range slots {
			slot = strings.TrimSpace(slot)
			if slot == "" {
				continue
			}

			var obj TimetableSlotV1
			obj.Slot = slot
			obj.CourseName = courseCode
			obj.CourseFullName = courseName
			obj.Venue = venue
			obj.CourseType = "Theory"

			Slots = append(Slots, obj)
		}
	}

	return Slots
}

func parseChennaiFormat(text string) []TimetableSlotV1 {
	var Slots []TimetableSlotV1

	coursePattern := regexp.MustCompile(`(?s)\n(\d+)\n.*?Registered[^\n\d]*(?:\n|$)`)
	courseMatches := coursePattern.FindAllString(text, -1)

	for _, courseEntry := range courseMatches {
		courseCodeNameRegex := regexp.MustCompile(`([A-Z]{3,4}[0-9]{3,4}[LPNEJM]?)\s*-\s*([^\n(]+)`)
		codeNameMatch := courseCodeNameRegex.FindStringSubmatch(courseEntry)

		if len(codeNameMatch) < 3 {
			continue
		}

		courseCode := strings.TrimSpace(codeNameMatch[1])
		courseName := strings.TrimSpace(codeNameMatch[2])

		slotVenueRegex := regexp.MustCompile(`([A-Z]*[0-9]*[A-Z]{1,3}[0-9]{1,2}(?:\+[TA]*[A-Z]{1,3}[0-9]{1,2})*|NIL)\s*-\s*\n+\s*([A-Z0-9\s-]+?)(?:\n|$)`)
		slotVenueMatch := slotVenueRegex.FindStringSubmatch(courseEntry)

		if len(slotVenueMatch) < 3 {
			continue
		}

		slotStr := strings.TrimSpace(slotVenueMatch[1])
		venue := strings.TrimSpace(slotVenueMatch[2])

		if slotStr == "NIL" || venue == "NIL" {
			continue
		}

		slots := strings.Split(slotStr, "+")

		for _, slot := range slots {
			slot = strings.TrimSpace(slot)
			if slot == "" {
				continue
			}

			var obj TimetableSlotV1
			obj.Slot = slot
			obj.CourseName = courseCode
			obj.CourseFullName = courseName
			obj.Venue = venue

			if len(slot) > 0 && slot[0:1] == "L" {
				obj.CourseType = "Lab"
			} else {
				obj.CourseType = "Theory"
			}

			Slots = append(Slots, obj)
		}
	}

	return Slots
}

func SlotsV1ToSlotsV2(slots []TimetableSlotV1, campus string) []models.Slot {
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

		slotV2.AddSlotTime(campus)

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

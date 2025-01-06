package commands

import (
	"fmt"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/utils"
	"github.com/urfave/cli/v2"
)

var TimetableCommands = []*cli.Command{
	{
		Name:    "parse-timetable",
		Aliases: []string{"ptt"},
		Usage:   "Parse a timetable",
		Action:  parseTimetable,
	},
	{
		Name:    "fix-slot-times",
		Aliases: []string{"fst"},
		Usage:   "Fix slot times",
		Action:  fixSlotTimes,
	},
	{

		Name:    "seed-course-table",
		Aliases: []string{"sct"},
		Usage:   "populate course table",
		Action:  seedCourseTable,
	},
}

func parseTimetable(c *cli.Context) error {
	var timetableText string
	fmt.Println("Enter the timetable text:")
	fmt.Scanln(&timetableText)

	var timetableV1 []utils.TimetableSlotV1
	timetableV1, err := utils.DetectTimetableV2(timetableText)
	if err != nil {
		return err
	}

	fmt.Println("Parsed data: ")
	fmt.Println(timetableV1)
	fmt.Println("\n\n")

	var timetableSlots []models.Slot
	for _, slot := range timetableV1 {
		timetableSlots = append(timetableSlots, models.Slot{
			Slot:  slot.Slot,
			Name:  slot.CourseFullName,
			Code:  slot.CourseName,
			Type:  slot.CourseType,
			Venue: slot.Venue,
		})
	}

	fmt.Println("Slots: ")
	fmt.Println(timetableSlots)
	return nil
}

func fixSlotTimes(c *cli.Context) error {
	users := []models.User{}
	database.DB.Find(&users)
	for _, user := range users {
		timetable := user.GetTimeTable()
		var slots []models.Slot
		for _, slot := range timetable.Slots {
			err := slot.AddSlotTime()
			if err != nil {
				fmt.Println("Error adding slot time: ", err)
			}
			slots = append(slots, slot)
		}
		timetable.Slots = slots
		user.Save()
	}
	return nil
}

func seedCourseTable(c *cli.Context) error {
	reset := "\033[0m"
	red := "\033[31m"
	green := "\033[32m"
	cyan := "\033[36m "

	fmt.Print(cyan, "Seeding.. ", reset)
	err := database.DB.Exec(`
		INSERT INTO courses (course_id, course_name)
		SELECT
		DISTINCT  ON(elems.data->>'code')
			elems.data->>'code' AS CourseCode,
			elems.data->>'name' AS CourseName
		FROM
			timetables,
			jsonb_array_elements(timetables.slots::jsonb) AS elems(data)	
	`).Error

	if err != nil {
		fmt.Println(red, "Failed")
		fmt.Println("Error: ", err, reset)
	}

	fmt.Println(green, "Done", reset)

	return nil
}

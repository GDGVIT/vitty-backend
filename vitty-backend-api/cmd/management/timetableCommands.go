package management

import (
	"fmt"

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

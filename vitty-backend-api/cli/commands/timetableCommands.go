package commands

import (
	"encoding/json"
	"fmt"
	"os"

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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "campus",
				Value: "vellore",
				Usage: "Campus to parse timetable for (vellore, chennai, bhopal)",
			},
		},
		Action: parseTimetable,
	},
	{
		Name:    "fix-slot-times",
		Aliases: []string{"fst"},
		Usage:   "Fix slot times",
		Action:  fixSlotTimes,
	},
	{
		Name:    "empty-rooms",
		Aliases: []string{"er"},
		Usage:   "Generates empty classrooms file",
		Action:  GenerateEmptyRooms,
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
	campus := c.String("campus")
	if campus == "" {
		campus = "vellore"
	}

	var timetableV1 []utils.TimetableSlotV1
	timetableV1, err := utils.DetectTimetableV2(timetableText, campus)
	if err != nil {
		return err
	}

	fmt.Println("Parsed data: ")
	fmt.Println(timetableV1)
	fmt.Print("\n\n")

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

		campus := "vellore"
		if user.Campus != nil {
			campus = string(*user.Campus)
		}

		for _, slot := range timetable.Slots {
			err := slot.AddSlotTime(campus)
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

func GenerateEmptyRooms(c *cli.Context) error {
	reset := "\033[0m"
	red := "\033[31m"
	green := "\033[32m"
	cyan := "\033[36m "

	fmt.Print(cyan, "Initiating ", reset)
	fmt.Print("Extracting Class details... ")

	err := database.DB.Exec(`
		Drop table IF EXISTS joinData;
		CREATE TABLE joinData (
			class text,
			slots JSONB
		);

		INSERT INTO joinData (class, slots)
		SELECT
			elems.data->>'venue' AS venue,
			jsonb_agg( DISTINCT elems.data->>'slot') AS slots
			FROM
			timetables,
			jsonb_array_elements(timetables.slots::jsonb) AS elems(data)
			GROUP BY
			elems.data->>'venue';
	`).Error

	if err != nil {
		fmt.Println(red, "Failed")
		fmt.Println("Error: ", err, reset)
	}

	fmt.Println(green, "Complete", reset)

	fmt.Print(cyan, "Initiating ", reset)
	fmt.Print("Looking  for empty classes... ")

	emptyClassRoomsJson := make(map[string]interface{})

	allSlots := make(map[string]bool)

	for _, slot := range models.TimetableSlots {
		allSlots[slot] = true
	}

	for slot := range allSlots {
		freeClasses, err := findEmptyClassRooms(slot)

		if err != nil {
			fmt.Println(red, "Failed")
			fmt.Printf("Slot %s was not able to be processed\nError: %s %s", slot, err, reset)
		}

		emptyClassRoomsJson[slot] = freeClasses
	}

	fmt.Println(green, "Complete", reset)
	fmt.Print(cyan, "Initiating ", reset)
	fmt.Print("Saving result... ")

	jsonData, err := json.Marshal(emptyClassRoomsJson)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}

	err = database.DB.Exec(`
		Drop table joindata;
	`).Error

	if err != nil {
		fmt.Println(red, "Failed")
		fmt.Println("Error: ", err, reset)
	}

	err = os.WriteFile("./data/freeClasses.json", jsonData, 0644)

	if err != nil {
		fmt.Println(red, "Failed")
		fmt.Println("Error: ", err, reset)
		return err
	}

	fmt.Println(green, "Complete", reset)

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

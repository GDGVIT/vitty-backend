package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	app.Get("/", root)
	app.Post("/uploadtext/", getTimetable)
	app.Post("/v2/uploadtext/", getTimetableV2)

	log.Fatal(app.Listen(8080))
}

// Response Structs

type TimetableResponse struct {
	ParsedData string `json:"Parsed_Data"`
	Slot       string `json:"Slot"`
	CourseName string `json:"Course_Name"`
	CourseType string `json:"Course_type"`
	Venue      string `json:"Venue"`
}

type TimetableResponseV2 struct {
	ParsedData     string `json:"Parsed_Data"`
	Slot           string `json:"Slot"`
	CourseName     string `json:"Course_Name"`
	CourseFullName string `json:"Course_Full_Name"`
	CourseType     string `json:"Course_type"`
	Venue          string `json:"Venue"`
}

type ErrorMessage struct {
	Message string `json:"detail"`
}

// Error Handling
func errorMessage(message string) ErrorMessage {
	var Error ErrorMessage
	Error.Message = message
	return Error
}

func root(c *fiber.Ctx) {
	c.SendString("Ok! Working!")
}

func getTimetable(c *fiber.Ctx) {
	data := c.FormValue("request")

	re := regexp.MustCompile(`[A-Z]{1,3}[0-9]{1,2}[\D]{1}[A-Z]{3,4}[0-9]{3,4}[A-Z]{0,1}[\D]{1}[A-Z]{2,3}[\D]{1}[A-Z]{2,6}[0-9]{2,4}[A-Za-z]{0,1}[\D]{1}[A-Z]{2,4}[0-9]{0,3}`)
	slots := re.FindAllString(data, -1)
	var timetable []TimetableResponse
	response := make(map[string]interface{})

	for _, slot := range slots {
		var obj TimetableResponse

		obj.ParsedData = slot
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
		timetable = append(timetable, obj)
	}

	response["Slots"] = timetable
	c.Status(fiber.StatusAccepted).JSON(response)
}

func getTimetableV2(c *fiber.Ctx) {
	data := c.FormValue("request")
	data = strings.ReplaceAll(data, "\r", "")

	var code, name, venue []string
	var slots [][]string
	var timetable []TimetableResponseV2
	response := make(map[string]interface{})

	re := regexp.MustCompile("[A-Z]{4}[0-9]{3}.+\n|[A-Z]{3}[0-9]{4}.+\n").FindAllString(data, -1)

	for _, c := range re {
		split := strings.Split(c[:len(c)-1], " - ")

		if len(split) != 2 {
			goto throwerror
		}
		code = append(code, split[0])
		name = append(name, split[1])

	}

	re = regexp.MustCompile("\n{3}[A-Z]+[0-9]{1,3}.+\n|\n{3}NIL\n").FindAllString(data, -1)
	for _, c := range re {
		venue = append(venue, c[3:len(c)-1])
	}

	re = regexp.MustCompile(".+[1-9].+[-]\n|NIL.+[-]\n").FindAllString(data, -1)
	for _, c := range re {
		slots = append(slots, strings.Split(c[:len(c)-3], "+"))
	}

	if len(slots) == len(name) && len(slots) == len(code) && len(slots) == len(venue) {

		for i := 0; i < len(slots); i++ {
			if slots[i][0] != "NIL" {
				for _, slot := range slots[i] {
					var obj TimetableResponseV2
					obj.ParsedData = "NIL"
					obj.Slot = slot
					obj.CourseName = code[i]
					obj.CourseFullName = name[i]
					obj.Venue = venue[i]
					if slot[0:1] == "L" {
						obj.CourseType = "Lab"
					} else {
						obj.CourseType = "Theory"
					}
					timetable = append(timetable, obj)
				}
			}
		}
		response["Slots"] = timetable
		c.Status(fiber.StatusAccepted).JSON(response)
		return
	} else {
		goto throwerror
	}

throwerror:
	c.Status(fiber.StatusBadRequest).JSON(errorMessage("Invalid Data"))
}

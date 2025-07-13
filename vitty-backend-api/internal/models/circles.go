package models

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
	"gorm.io/gorm"
)

type Circles struct {
	CircleId       string `json:"circle_id" gorm:"unique"`
	CircleName     string `json:"circle_name" gorm:"primaryKey"`
	Uname          string `json:"omitempty" gorm:"primaryKey"`
	CircleSlots    string
	CircleJoinCode string `json:"circle_join_code" gorm:"default:null"`
	User           User   `gorm:"foreignKey:Uname;references:Username;"`
}

func (c *Circles) CreateCircle() error {
	err := database.DB.Create(c).Error
	return err
}

func (c *Circles) GetCircleByCircleId(circleId string) error {
	results := database.DB.Where("circle_join_code = ?", circleId).First(&c).Error
	return results
}

func (c *Circles) GetCircleByJoinCode(joinCode string) error {
	err := database.DB.Where("circle_join_code = ?", joinCode).First(&c).Error
	return err
}

func (c *Circles) GetCircleSlots() (error, map[string]string) {
	err := database.DB.Where(c).First(&c).Error
	return err, jsonToCircleSlots(c.CircleSlots)
}

func (c *Circles) UpdateCircleName(circleName string) error {
	result := database.DB.Model(&Circles{}).Where("circle_id like ?", c.CircleId).Update("circle_name", circleName)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (c *Circles) UpdateCircleUsername(username string) error {
	err := database.DB.Model(&Circles{}).Where(c).Update("uname", username).Error
	return err
}

func (c *Circles) updateCircleSlots(circleSlotMap map[string]string) error {
	jsonString := circleSlotsToJson(circleSlotMap)
	err := database.DB.Model(&Circles{}).Where(c).Update("circle_slots", jsonString).Error
	return err
}

func (c *Circles) CreateCircleJoinCode(joinCode string) error {
	err := database.DB.Model(&Circles{}).Where("circle_id like ?", c.CircleId).Update("circle_join_code", joinCode).Error
	return err
}

func (c *Circles) DeleteCircle() error {
	result := database.DB.Where(&c).Delete(&Circles{})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (c *Circles) ComputeCircleSlots(user User, recompute bool) error {

	resultMap := make(map[string]string)

	campus := "vellore"
	if user.Campus != nil {
		campus = string(*user.Campus)
	}

	dayWiseSlots := user.GetTimeTable().GetDaywiseTimetable(campus)

	var err error
	var circleSlotMap map[string]string

	if !recompute {
		err, circleSlotMap = c.GetCircleSlots()

		if err != nil {
			return err
		}
	}

	for day := time.Monday; day <= time.Friday; day++ {
		var result []string
		var dupEightSlot bool

		referenceTable := make(map[string]struct{})

		for _, slot := range dayWiseSlots[day.String()] {
			referenceTime := slot.StartTime.Format("15:04")
			referenceTable[referenceTime] = struct{}{}
		}

		timings := append(TheoryTimings, LabTimings...)

		if c.CircleSlots == "" {
			for _, timing := range timings {
				timeOfTiming := strings.Split(timing.StartTime, "T")[1]
				if _, found := referenceTable[timeOfTiming]; !found {
					if strings.Contains(timeOfTiming, "08:00") {
						if dupEightSlot {
							continue
						}
						dupEightSlot = true
					}

					result = append(result, timeOfTiming)
				}
			}

		} else {

			timings := strings.Split(circleSlotMap[day.String()], ",")

			for _, timing := range timings {
				if _, found := referenceTable[timing]; !found {
					if strings.Contains(timing, "08:00") {
						if dupEightSlot {
							continue
						}
						dupEightSlot = true
					}

					result = append(result, timing)
				}
			}

		}
		resultMap[day.String()] = strings.Join(result, ",")

	}

	c.updateCircleSlots(resultMap)

	return nil

}

// Helper Functions

func circleSlotsToJson(circleSlotMap map[string]string) string {
	jsonBytes, err := json.Marshal(circleSlotMap)

	if err != nil {
		log.Println(err)
	}
	return string(jsonBytes)
}

func jsonToCircleSlots(jsonString string) map[string]string {
	resultMap := make(map[string]string)
	if jsonString == "" {
		return resultMap
	}

	if err := json.Unmarshal([]byte(jsonString), &resultMap); err != nil {
		log.Fatal(err)
	}

	return resultMap
}

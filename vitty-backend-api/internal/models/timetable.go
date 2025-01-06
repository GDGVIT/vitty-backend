package models

import (
	"log"
	"time"

	"golang.org/x/exp/slices"
)

const STD_REF_TIME = "2006-01-02T15:04"

type Slot struct {
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Venue     string    `json:"venue"`
	Slot      string    `json:"slot"`
	Type      string    `json:"type"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

type Timetable struct {
	User         User   `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserUsername;references:Username"`
	UserUsername string `gorm:"primaryKey"`
	Slots        []Slot `gorm:"serializer:json"`
}

func (t Timetable) GetDaySlots(day time.Weekday) map[string][]Slot {
	resp := make(map[string][]Slot)
	var data []Slot
	daySlots := DailySlots[day.String()]
	labSlot := ""

	var err error
	// Theory slots
	for _, slot := range t.Slots {

		if slot.Type == "Theory" && slices.Contains(daySlots["Theory"], slot.Slot) {
			index := slices.Index(daySlots["Theory"], slot.Slot)
			slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, TheoryTimings[index].StartTime, time.Local)
			if err != nil {
				log.Println("Error parsing time: ", err)
				return nil
			}

			slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, TheoryTimings[index].EndTime, time.Local)

			if err != nil {
				log.Println("Error parsing time: ", err)
				return nil
			}

			data = append(data, slot)
		} else if slot.Type == "Lab" && slices.Contains(daySlots["Lab"], slot.Slot) {
			index := slices.Index(daySlots["Lab"], slot.Slot)

			if labSlot == "" {
				labSlot += slot.Slot + "+"
				continue
			} else {
				slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, LabTimings[index-1].StartTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}

				slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, LabTimings[index].EndTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}
				slot.Slot = labSlot + slot.Slot
				labSlot = ""
				data = append(data, slot)
			}

		}
	}
	resp[day.String()] = data
	return resp
}

func (t Timetable) GetDaywiseTimetable() map[string][]Slot {
	resp := make(map[string][]Slot)
	labSlot := ""

	for _, slot := range t.Slots {
		for day, value := range DailySlots {

			if slices.Contains(value["Theory"], slot.Slot) {
				index := slices.Index(value["Theory"], slot.Slot)
				var err error
				slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, TheoryTimings[index].StartTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}
				slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, TheoryTimings[index].EndTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}
				resp[day] = append(resp[day], slot)
			} else if slices.Contains(value["Lab"], slot.Slot) {
				index := slices.Index(value["Lab"], slot.Slot)

				if labSlot == "" {
					labSlot += slot.Slot + "+"
					continue
				} else {
					var err error
					slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, LabTimings[index-1].StartTime, time.Local)
					if err != nil {
						log.Println("Error parsing time: ", err)
						return nil
					}
					slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, LabTimings[index].EndTime, time.Local)
					if err != nil {
						log.Println("Error parsing time: ", err)
						return nil
					}

					slot.Slot = labSlot + slot.Slot
					labSlot = ""
					resp[day] = append(resp[day], slot)
				}
			}
		}
	}
	return resp
}

func (s *Slot) AddSlotTime() error {
	for _, value := range DailySlots {
		if slices.Contains(value["Theory"], s.Slot) {
			index := slices.Index(value["Theory"], s.Slot)
			var err error
			s.StartTime, err = time.ParseInLocation(STD_REF_TIME, TheoryTimings[index].StartTime, time.Local)
			if err != nil {
				return err
			}
			s.EndTime, err = time.ParseInLocation(STD_REF_TIME, TheoryTimings[index].EndTime, time.Local)
			if err != nil {
				return err
			}
		} else if slices.Contains(value["Lab"], s.Slot) {
			index := slices.Index(value["Lab"], s.Slot)
			var err error
			s.StartTime, err = time.ParseInLocation(STD_REF_TIME, LabTimings[index].StartTime, time.Local)
			if err != nil {
				return err
			}
			s.EndTime, err = time.ParseInLocation(STD_REF_TIME, LabTimings[index].EndTime, time.Local)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

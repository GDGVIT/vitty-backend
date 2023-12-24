package models

import (
	"log"
	"time"

	"golang.org/x/exp/slices"
)

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
	User         User   `gorm:"foreignKey:UserUsername;references:Username"`
	UserUsername string `gorm:"primaryKey"`
	Slots        []Slot `gorm:"serializer:json"`
}

func (t Timetable) GetDaySlots(day time.Weekday) map[string][]Slot {
	resp := make(map[string][]Slot)
	var data []Slot
	daySlots := DailySlots[day.String()]

	// Theory slots
	for _, slot := range t.Slots {
		if slot.Type == "Theory" && slices.Contains(daySlots["Theory"], slot.Slot) {
			data = append(data, slot)
		} else if slot.Type == "Lab" && slices.Contains(daySlots["Lab"], slot.Slot) {
			data = append(data, slot)
		}
	}
	resp[day.String()] = data
	return resp
}

func (t Timetable) GetDaywiseTimetable() map[string][]Slot {
	resp := make(map[string][]Slot)

	for _, slot := range t.Slots {
		for day, value := range DailySlots {
			if slices.Contains(value["Theory"], slot.Slot) {
				index := slices.Index(value["Theory"], slot.Slot)
				var err error
				slot.StartTime, err = time.Parse("15:04", TheoryTimings[index].StartTime)
				if err != nil {
					log.Println("Error parsing time: ", err)
				}
				slot.EndTime, err = time.Parse("15:04", TheoryTimings[index].EndTime)
				if err != nil {
					log.Println("Error parsing time: ", err)
				}
				resp[day] = append(resp[day], slot)
			} else if slices.Contains(value["Lab"], slot.Slot) {
				index := slices.Index(value["Lab"], slot.Slot)
				var err error
				slot.StartTime, err = time.Parse("15:04", LabTimings[index].StartTime)
				if err != nil {
					log.Println("Error parsing time: ", err)
				}
				slot.EndTime, err = time.Parse("15:04", LabTimings[index].EndTime)
				if err != nil {
					log.Println("Error parsing time: ", err)
				}
				resp[day] = append(resp[day], slot)
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
			s.StartTime, err = time.Parse("15:04", TheoryTimings[index].StartTime)
			if err != nil {
				return err
			}
			s.EndTime, err = time.Parse("15:04", TheoryTimings[index].EndTime)
			if err != nil {
				return err
			}
		} else if slices.Contains(value["Lab"], s.Slot) {
			index := slices.Index(value["Lab"], s.Slot)
			var err error
			s.StartTime, err = time.Parse("15:04", LabTimings[index].StartTime)
			if err != nil {
				return err
			}
			s.EndTime, err = time.Parse("15:04", LabTimings[index].EndTime)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

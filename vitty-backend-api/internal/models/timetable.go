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

func (t Timetable) GetDaySlots(day time.Weekday, campus string) map[string][]Slot {
	resp := make(map[string][]Slot)
	var data []Slot
	daySlots := GetDailySlotsForCampus(campus)[day.String()]

	if campus == "bhopal" {
		for _, slot := range t.Slots {
			if slot.Type == "Theory" && slices.Contains(daySlots["Theory"], slot.Slot) {
				if timingIndex, exists := GetBhopalSlotTimingIndex(slot.Slot); exists {
					theoryTimings := GetTheoryTimingsForCampus(campus)
					if timingIndex < len(theoryTimings) {
						var err error
						slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].StartTime, time.Local)
						if err != nil {
							log.Println("Error parsing time: ", err)
							return nil
						}
						slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].EndTime, time.Local)
						if err != nil {
							log.Println("Error parsing time: ", err)
							return nil
						}
						data = append(data, slot)
					}
				}
			}
		}
		resp[day.String()] = data
		return resp
	}

	if campus == "chennai" {
		for _, slot := range t.Slots {
			if slot.Type == "Theory" && slices.Contains(daySlots["Theory"], slot.Slot) {
				if timingIndex, exists := GetChennaiSlotTimingIndex(slot.Slot); exists {
					theoryTimings := GetTheoryTimingsForCampus(campus)
					if timingIndex < len(theoryTimings) {
						var err error
						slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].StartTime, time.Local)
						if err != nil {
							log.Println("Error parsing time: ", err)
							return nil
						}
						slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].EndTime, time.Local)
						if err != nil {
							log.Println("Error parsing time: ", err)
							return nil
						}
						data = append(data, slot)
					}
				}
			} else if slot.Type == "Lab" && slices.Contains(daySlots["Lab"], slot.Slot) {
				labTimings := GetLabTimingsForCampus(campus)
				index := slices.Index(daySlots["Lab"], slot.Slot)
				var err error
				slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].StartTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}
				slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].EndTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}
				data = append(data, slot)
			}
		}
		resp[day.String()] = data
		return resp
	}

	theoryTimings := GetTheoryTimingsForCampus(campus)
	labTimings := GetLabTimingsForCampus(campus)
	labSlot := ""

	var err error
	for _, slot := range t.Slots {
		if slot.Type == "Theory" && slices.Contains(daySlots["Theory"], slot.Slot) {
			index := slices.Index(daySlots["Theory"], slot.Slot)
			slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[index].StartTime, time.Local)
			if err != nil {
				log.Println("Error parsing time: ", err)
				return nil
			}

			slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[index].EndTime, time.Local)

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
				slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index-1].StartTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}

				slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].EndTime, time.Local)
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

func (t Timetable) GetDaywiseTimetable(campus string) map[string][]Slot {
	resp := make(map[string][]Slot)
	dailySlots := GetDailySlotsForCampus(campus)

	if campus == "bhopal" {
		for _, slot := range t.Slots {
			for day, value := range dailySlots {
				if slices.Contains(value["Theory"], slot.Slot) {
					if timingIndex, exists := GetBhopalSlotTimingIndex(slot.Slot); exists {
						theoryTimings := GetTheoryTimingsForCampus(campus)
						if timingIndex < len(theoryTimings) {
							var err error
							slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].StartTime, time.Local)
							if err != nil {
								log.Println("Error parsing time: ", err)
								return nil
							}
							slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].EndTime, time.Local)
							if err != nil {
								log.Println("Error parsing time: ", err)
								return nil
							}
							resp[day] = append(resp[day], slot)
						}
					}
				}
			}
		}
		return resp
	}

	if campus == "chennai" {
		for _, slot := range t.Slots {
			for day, value := range dailySlots {
				if slices.Contains(value["Theory"], slot.Slot) {
					if timingIndex, exists := GetChennaiSlotTimingIndex(slot.Slot); exists {
						theoryTimings := GetTheoryTimingsForCampus(campus)
						if timingIndex < len(theoryTimings) {
							var err error
							slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].StartTime, time.Local)
							if err != nil {
								log.Println("Error parsing time: ", err)
								return nil
							}
							slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].EndTime, time.Local)
							if err != nil {
								log.Println("Error parsing time: ", err)
								return nil
							}
							resp[day] = append(resp[day], slot)
						}
					}
				} else if slices.Contains(value["Lab"], slot.Slot) {
					index := slices.Index(value["Lab"], slot.Slot)
					labTimings := GetLabTimingsForCampus(campus)
					var err error
					slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].StartTime, time.Local)
					if err != nil {
						log.Println("Error parsing time: ", err)
						return nil
					}
					slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].EndTime, time.Local)
					if err != nil {
						log.Println("Error parsing time: ", err)
						return nil
					}
					resp[day] = append(resp[day], slot)
				}
			}
		}
		return resp
	}

	theoryTimings := GetTheoryTimingsForCampus(campus)
	labTimings := GetLabTimingsForCampus(campus)
	labSlot := ""

	for _, slot := range t.Slots {
		for day, value := range dailySlots {

			if slices.Contains(value["Theory"], slot.Slot) {
				index := slices.Index(value["Theory"], slot.Slot)
				var err error
				slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[index].StartTime, time.Local)
				if err != nil {
					log.Println("Error parsing time: ", err)
					return nil
				}
				slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[index].EndTime, time.Local)
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
					slot.StartTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index-1].StartTime, time.Local)
					if err != nil {
						log.Println("Error parsing time: ", err)
						return nil
					}
					slot.EndTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].EndTime, time.Local)
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

func (s *Slot) AddSlotTime(campus string) error {
	log.Printf("AddSlotTime called for slot: %s, campus: %s, type: %s", s.Slot, campus, s.Type)

	if campus == "bhopal" {
		if s.Type == "Theory" {
			if timingIndex, exists := GetBhopalSlotTimingIndex(s.Slot); exists {
				theoryTimings := GetTheoryTimingsForCampus(campus)
				if timingIndex < len(theoryTimings) {
					var err error
					s.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].StartTime, time.Local)
					if err != nil {
						return err
					}
					s.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].EndTime, time.Local)
					if err != nil {
						return err
					}
					log.Printf("Bhopal slot %s assigned timing: %s - %s", s.Slot, s.StartTime.Format("15:04"), s.EndTime.Format("15:04"))
				}
			}
		}
		return nil
	}

	if campus == "chennai" {
		if s.Type == "Theory" {
			if timingIndex, exists := GetChennaiSlotTimingIndex(s.Slot); exists {
				theoryTimings := GetTheoryTimingsForCampus(campus)
				if timingIndex < len(theoryTimings) {
					var err error
					s.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].StartTime, time.Local)
					if err != nil {
						return err
					}
					s.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[timingIndex].EndTime, time.Local)
					if err != nil {
						return err
					}
					log.Printf("Chennai slot %s assigned timing: %s - %s", s.Slot, s.StartTime.Format("15:04"), s.EndTime.Format("15:04"))
				}
			}
		} else if s.Type == "Lab" {
			// Chennai lab timing uses same logic as Vellore
			labTimings := GetLabTimingsForCampus(campus)
			dailySlots := GetDailySlotsForCampus(campus)

			for _, value := range dailySlots {
				if slices.Contains(value["Lab"], s.Slot) {
					index := slices.Index(value["Lab"], s.Slot)
					var err error
					s.StartTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].StartTime, time.Local)
					if err != nil {
						return err
					}
					s.EndTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].EndTime, time.Local)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		return nil
	}

	// Default Vellore campus logic
	dailySlots := GetDailySlotsForCampus(campus)
	theoryTimings := GetTheoryTimingsForCampus(campus)
	labTimings := GetLabTimingsForCampus(campus)

	for _, value := range dailySlots {
		if slices.Contains(value["Theory"], s.Slot) {
			index := slices.Index(value["Theory"], s.Slot)
			var err error
			s.StartTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[index].StartTime, time.Local)
			if err != nil {
				return err
			}
			s.EndTime, err = time.ParseInLocation(STD_REF_TIME, theoryTimings[index].EndTime, time.Local)
			if err != nil {
				return err
			}
		} else if slices.Contains(value["Lab"], s.Slot) {
			index := slices.Index(value["Lab"], s.Slot)
			var err error
			s.StartTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].StartTime, time.Local)
			if err != nil {
				return err
			}
			s.EndTime, err = time.ParseInLocation(STD_REF_TIME, labTimings[index].EndTime, time.Local)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

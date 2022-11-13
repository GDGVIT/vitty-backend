package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vitty-backend/vitty-backend-gofiber/pkg/models"
)

func StringifySlots(slots []models.Slot) (string, error) {
	var stringSlots []string

	for _, slot := range slots {
		slotString, err := json.Marshal(slot)

		if err != nil {
			return "", err
		}
		stringSlots = append(stringSlots, string(slotString))
	}

	return strings.Join(stringSlots, "|"), nil
}

func UnStringifySlots(slots string) ([]models.Slot, error) {
	var slotsArray []models.Slot

	slotsArrayString := strings.Split(slots, "|")

	for _, slot := range slotsArrayString {

		var slotObj models.Slot
		err := json.Unmarshal([]byte(slot), &slotObj)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		slotsArray = append(slotsArray, slotObj)
	}

	return slotsArray, nil
}

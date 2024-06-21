package commands

import (
	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/database"
)

func findEmptyClassRooms(slot string) ([]string, error) {
	var freeClasses []string
	total := 0
	offset := 0
	limit := 1000

	query := database.DB.
		Table("joindata").
		Where("NOT (slots @> '[\"?\"]')", database.DB.Raw(slot)).
		Where("slots::text !~ '\\[\"L.*\"\\]'")

	err := query.
		Select("COUNT(class)").
		Scan(&total).Error

	if err != nil {
		return freeClasses, err
	}

	for total >= 0 {
		err := query.
			Select("class").
			Limit(limit).
			Offset(offset).
			Find(&freeClasses).Error

		if err != nil {
			return freeClasses, err
		}

		total -= limit
		offset += limit
	}

	return freeClasses, nil
}

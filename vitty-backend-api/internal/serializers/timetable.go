package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func TimetableSerializer(timetable models.Timetable) map[string]interface{} {
	return map[string]interface{}{
		"data": timetable.GetDaywiseTimetable(),
	}
}

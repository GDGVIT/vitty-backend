package serializers

import "github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"

func RemindersSerializer(reminders []models.Reminders) []map[string]interface{} {
	var result []map[string]interface{}

	for _, reminder := range reminders {
		out := map[string]interface{}{
			"reminder_id":      reminder.ReminderId,
			"reminder_name":    reminder.ReminderName,
			"reminder_content": reminder.ReminderContent,
			"reminder_time":    reminder.ReminderTime,
		}
		result = append(result, out)
	}
	return result
}

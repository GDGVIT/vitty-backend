package jobs

import (
	"log"
	"time"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
)

func InitializeJobs(isRunJobs string) {
	if isRunJobs == "true" {
		startDailyJob()
		startWeeklyJob()
		log.Println("Jobs are running")
	} else {
		log.Println("Jobs disabled")
	}

}

var reminder models.Reminders

func startDailyJob() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reminder.SoftDeleteExpiredReminders()
		}
	}
}

func startWeeklyJob() {
	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reminder.CleanupOldReminders()
		}
	}
}

package jobs

import (
	"log"
	"time"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/models"
)

type JobTime struct {
	DAILY_JOB_MIN  int
	DAILY_JOB_SEC  int
	WEEKLY_JOB_MIN int
	WEEKLY_JOB_SEC int
}

func InitializeJobs(isRunJobs string, debug string, jobTimes JobTime) {
	if isRunJobs == "true" {
		if debug == "true" {

			timeDailyJob := time.Duration(jobTimes.DAILY_JOB_MIN)*time.Minute + time.Duration(jobTimes.DAILY_JOB_SEC)*time.Second
			timeWeeklyJob := time.Duration(jobTimes.WEEKLY_JOB_MIN)*time.Minute + time.Duration(jobTimes.WEEKLY_JOB_SEC)*time.Second
			go startDailyJob(timeDailyJob)
			go startWeeklyJob(timeWeeklyJob)
		} else {
			go startDailyJob(24 * time.Hour)
			go startWeeklyJob(7 * 24 * time.Hour)
		}
		log.Println("Jobs are running")
	} else {
		log.Println("Jobs disabled")
	}

}

var reminder models.Reminders

func startDailyJob(jobTime time.Duration) {
	ticker := time.NewTicker(jobTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reminder.SoftDeleteExpiredReminders()
		}
	}
}

func startWeeklyJob(jobTime time.Duration) {
	ticker := time.NewTicker(jobTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reminder.CleanupOldReminders()
		}
	}
}

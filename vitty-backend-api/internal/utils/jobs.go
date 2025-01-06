package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/GDGVIT/vitty-backend/vitty-backend-api/internal/jobs"
)

func ParseJobTimes() jobs.JobTime {
	var err error
	var jobTime jobs.JobTime

	jobTime.DAILY_JOB_MIN, err = strconv.Atoi(os.Getenv("DAILY_JOB_MIN"))

	if err != nil {
		log.Fatal(err, "JOB_DAILY_MIN_PARSE_ERR")
	}

	jobTime.DAILY_JOB_SEC, err = strconv.Atoi(os.Getenv("DAILY_JOB_SEC"))

	if err != nil {
		log.Fatal(err, "JOB_DAILY_SEC_PARSE_ERR")
	}

	jobTime.WEEKLY_JOB_MIN, err = strconv.Atoi(os.Getenv("WEEKLY_JOB_MIN"))

	if err != nil {
		log.Fatal(err, "JOB_WEEKLY_MIN_PARSE_ERR")
	}

	jobTime.WEEKLY_JOB_SEC, err = strconv.Atoi(os.Getenv("WEEKLY_JOB_SEC"))

	if err != nil {
		log.Fatal(err, "JOB_WEEKLY_SEC_PARSE_ERR")
	}

	log.Println("Parse Done")
	return jobTime
}

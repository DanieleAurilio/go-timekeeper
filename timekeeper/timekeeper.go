package timekeeper

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"timekeeper/models"
	"timekeeper/models/queue"
	"timekeeper/utils"
)

func Start() {
	fmt.Println("Starting timekeeper ...")

	// Initialize Queue
	var queue = &queue.Queue{}
	queue.Init()

	fmt.Println("Job Queue initialized. Start Pooling ...")

	// Set to 0 seconds
	t := time.Now()
	nextMinute := t.Truncate(time.Minute).Add(time.Minute)
	sleepDuration := nextMinute.Sub(t)
	time.Sleep(sleepDuration)

	poolms := utils.ReadConfigFile().PoolMs
	timePool := time.Duration(poolms) * time.Millisecond

	for {
		configJobs := utils.ReadConfigFile().Jobs
		poolJobs(queue, configJobs)
		queue.SortDescByTime()
		queue.ExecuteJobs()
		time.Sleep(timePool)
	}

}

func poolJobs(queue *queue.Queue, configJobs []models.JobConfig) {
	for i := 0; i < len(configJobs); i++ {
		currentJob := initJob(configJobs[i], i)

		if !currentJob.Enable {
			fmt.Println("Skip Job " + currentJob.ID + " is not enabled.")
			continue
		}

		queue.Set(currentJob)
	}
}

func initJob(config models.JobConfig, idx int) *queue.QueueJob {
	formattedDate, err := formatDate(&config.Schedule)

	if err != nil {
		fmt.Println("Skip Job " + config.ID + " on init job. Error: " + err.Error())
		return nil
	}

	return &queue.QueueJob{
		ID:       config.ID,
		Position: idx,
		Filename: config.Filename,
		Enable:   config.Enable,
		Params:   config.Params,
		Time:     formattedDate,
		Running:  false,
	}
}

func formatDate(schedule *models.Schedule) (string, error) {

	isEligibleMonth, month := isEligibleByMonth(schedule.Month)
	if !isEligibleMonth {
		return "", errors.New("no month found")
	}

	scheduledWeekdays := strings.Split(schedule.WeekDay, ",")
	day, isEligibleDay := isEligibleByDay(scheduledWeekdays)
	if !isEligibleDay {
		return "", errors.New("no day found")
	}
	scheduleTime := time.Date(time.Now().Year(), month, day, int(schedule.Hours), int(schedule.Minutes), int(schedule.Seconds), 0, time.Local)

	return scheduleTime.Format(time.DateTime), nil
}

func isEligibleByDay(scheduledWeekdays []string) (int, bool) {
	day := time.Now().Weekday().String()

	if strings.Contains(strings.Join(scheduledWeekdays, ","), "*") {
		return time.Now().Day(), true
	}

	for _, scheduledWeekday := range scheduledWeekdays {
		if strings.Contains(scheduledWeekday, day) {
			return time.Now().Day(), true
		}
	}

	return -1, false
}

func isEligibleByMonth(scheduledMonth string) (bool, time.Month) {
	month := time.Now().Month()

	if strings.Contains(scheduledMonth, "*") || strings.Contains(scheduledMonth, month.String()) {
		return true, month
	}

	if scheduledMonth != month.String() {
		fmt.Println("Month provided not correct: " + scheduledMonth + " instead of " + month.String())
		return false, month
	}

	return true, month
}

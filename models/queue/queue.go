package queue

import (
	"fmt"
	"reflect"
	"sort"
	"time"
	"timekeeper/utils"
)

type Queue map[QueueJobID]*QueueJob

type QueueJobID string

type QueueJob struct {
	ID       string
	Position int
	Filename string
	Enable   bool
	Params   map[string]string
	Time     string
	Running  bool
}

const mainFunc string = "main"

func (*Queue) Init() *Queue {
	var queue = make(Queue)
	return &queue
}

func (queue *Queue) Set(job *QueueJob) {
	fmt.Println("Append Job " + job.ID + " to queue.")
	(*queue)[QueueJobID(job.ID)] = job
}

func (queue *Queue) Get(jobID string) *QueueJob {
	if _, ok := (*queue)[QueueJobID(jobID)]; ok {
		return (*queue)[QueueJobID(jobID)]
	}

	return nil
}

func (queue *Queue) Delete(jobID string) {
	job := queue.Get(jobID)
	if jobRef := reflect.ValueOf(job); jobRef.IsNil() {
		fmt.Println("Cannot Delete " + jobID)
	} else {
		delete(*queue, QueueJobID(jobID))
	}
}

func (queue *Queue) SortDescByTime() {

	if len((*queue)) == 0 {
		return
	}

	var queueJob = make([]*QueueJob, 0)
	for _, job := range *queue {
		queueJob = append(queueJob, job)
	}

	fmt.Println("Sorting Jobs By Time... ")
	sort.Slice(queueJob, func(i int, j int) bool {
		timeA, _ := time.Parse(time.DateTime, queueJob[i].Time)
		timeB, _ := time.Parse(time.DateTime, queueJob[j].Time)
		return timeA.Before(timeB)
	})

	fmt.Println("Reinitialize Jobs Queue... ")
	for _, job := range queueJob {
		delete((*queue), QueueJobID(job.ID))
	}

	for _, job := range queueJob {
		(*queue).Set(job)
	}
}

func (queue *Queue) ExecuteJobs() {
	currentTime := time.Now().Format(time.DateTime)
	jobList := (*queue).GetJobsByTime(currentTime)

	if len(jobList) == 0 {
		fmt.Println("No job to execute.")
		return
	}

	for _, job := range jobList {
		fullFilePath := utils.ResolveJobPath(job.Filename)

		if len(fullFilePath) == 0 {
			continue
		}

		outputDir := utils.ReadConfigFile().OutputDir
		err, outputDirectoryPath := utils.BuildOutputDirectoryPath(outputDir)

		if err != nil {
			continue
		}

		utils.Mkdir(outputDirectoryPath)

		fmt.Println("Executing Datetime: " + time.Now().Format(time.DateTime))

		go utils.ExecuteJobCmd(fullFilePath, mainFunc, job.Params)
	}

}

func (queue *Queue) GetJobsByTime(timeVal string) []*QueueJob {
	var timeNow time.Time
	var errTime error
	if val := reflect.ValueOf(timeVal); val.IsValid() {
		timeNow, errTime = time.Parse(time.DateTime, timeVal)
	} else {
		timeNow, errTime = time.Parse(time.DateTime, time.Now().String())
	}

	if errTime != nil {
		fmt.Println(errTime.Error())
	}

	var jobList = make([]*QueueJob, 0)
	for _, currentJob := range *queue {
		jobTime, err := time.Parse(time.DateTime, currentJob.Time)

		if err != nil {
			fmt.Println("Please check time on job: " + currentJob.ID + ". Provide format i.e: 2006-01-02 15:04:05")
			return nil
		} else if utils.CheckTime(timeNow, jobTime) {
			jobList = append(jobList, currentJob)
			fmt.Println("Add Job " + currentJob.ID + " to list.")
		}
	}

	return jobList
}

func (queue *Queue) SetJobRun(jobID string) {
	(*queue)[QueueJobID(jobID)].Running = true
}

package cron

import (
	"sync"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/shashank-sharma/metadata/internal/logger"
)

type CronJob struct {
	Id          string
	Description string
	Interval    time.Duration
	Run         func()
	NextRun     binding.String
	mutex       sync.Mutex
	quit        chan struct{}
}

func (j *CronJob) Schedule() {
	j.quit = make(chan struct{})
	ticker := time.NewTicker(j.Interval)

	go func() {
		j.mutex.Lock()
		j.Run()
		nextRun := time.Now().Add(j.Interval).Format(time.RFC1123)
		j.NextRun.Set(nextRun)
		j.mutex.Unlock()
		for {
			select {
			case <-ticker.C:
				j.mutex.Lock()
				j.Run()
				nextRun := time.Now().Add(j.Interval).Format(time.RFC1123)
				j.NextRun.Set(nextRun)
				j.mutex.Unlock()
			case <-j.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (j *CronJob) Stop() {
	close(j.quit)
}

type CronService struct {
	cronJobs []*CronJob
}

func (cs *CronService) StopAllJobs() {
	logger.LogInfo("Stopping all CRON Jobs")
	for _, job := range cs.cronJobs {
		logger.LogInfo("Stopping CRON: ", job.Id)
		job.Stop()
	}
}

func (cs *CronService) AddJob(id, description string, interval time.Duration, run func()) *CronJob {
	logger.LogDebug("Starting job: ", id)
	nextRunBinding := binding.NewString()
	job := &CronJob{
		Id:          id,
		Interval:    interval,
		Description: description,
		Run:         run,
		NextRun:     nextRunBinding,
	}
	cs.cronJobs = append(cs.cronJobs, job)
	job.Schedule()
	return job
}

func (cs *CronService) GetJobs() []*CronJob {
	return cs.cronJobs
}

func NewCronService() *CronService {
	return &CronService{}
}

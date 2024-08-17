package cron

import (
	"sync"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/shashank-sharma/metadata/internal/logger"
)

type CronInfo struct {
	Id           string
	Description  string
	IsRunning    binding.Bool
	SuccessCount binding.Int
	FailedCount  binding.Int
	NextRun      binding.String
}

type CronJob struct {
	Interval time.Duration
	Run      func(cronInfo *CronInfo)
	CronInfo *CronInfo
	mutex    sync.Mutex
	quit     chan struct{}
}

func (j *CronJob) Execute() {
	j.CronInfo.IsRunning.Set(true)
	j.mutex.Lock()
	j.Run(j.CronInfo)
	nextRun := time.Now().Add(j.Interval).Format(time.RFC1123)
	j.CronInfo.NextRun.Set(nextRun)
	j.mutex.Unlock()
	j.CronInfo.IsRunning.Set(false)
}

func (j *CronJob) Schedule() {
	j.quit = make(chan struct{})
	ticker := time.NewTicker(j.Interval)

	go func() {
		j.Execute()
		for {
			select {
			case <-ticker.C:
				j.Execute()
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
		logger.LogInfo("Stopping CRON: ", job.CronInfo.Id)
		job.Stop()
	}
}

func (cs *CronService) AddJob(id, description string, interval time.Duration, run func(cronInfo *CronInfo)) *CronJob {
	logger.LogDebug("Starting job: ", id)
	nextRunBinding := binding.NewString()
	cronInfo := &CronInfo{
		Id:           id,
		Description:  description,
		IsRunning:    binding.NewBool(),
		SuccessCount: binding.NewInt(),
		FailedCount:  binding.NewInt(),
		NextRun:      nextRunBinding,
	}

	job := &CronJob{
		Interval: interval,
		Run:      run,
		CronInfo: cronInfo,
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

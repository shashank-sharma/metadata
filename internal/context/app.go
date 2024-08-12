package context

import (
	"fmt"
	"time"

	"github.com/shashank-sharma/metadata/internal/config"
	"github.com/shashank-sharma/metadata/internal/cron"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/models"
	"github.com/shashank-sharma/metadata/internal/services/activitywatch"
	"github.com/shashank-sharma/metadata/internal/services/backend"
	"github.com/shashank-sharma/metadata/internal/settings"
	"github.com/shashank-sharma/metadata/internal/theme"
)

type AppContext struct {
	BackendService *backend.BackendService
	AWService      *activitywatch.AWService
	Config         config.AppConfig
	State          models.BaseState
	Notification   *theme.Notification
	CronService    *cron.CronService
}

func NewAppContext(config config.AppConfig) *AppContext {
	bs := backend.NewBackendService(config.BackendEndpoint)
	aws := activitywatch.NewAWService(config.AWEndpoint)

	return &AppContext{
		BackendService: bs,
		AWService:      aws,
		State:          *models.NewGlobalState(),
		Notification:   theme.NewNotification(),
		Config:         config,
		CronService:    cron.NewCronService(),
	}
}

func (ac *AppContext) OnLoginSuccess() {
	logger.Debug.Println("Login Success")

	ac.InitializeAWState()
	if ac.Config.Settings.UserSettings.UserId != "" && ac.Config.Settings.UserSettings.Token != "" {
		ac.BackendService.Client.SetDevToken(fmt.Sprintf("%s.%s", ac.Config.Settings.UserSettings.UserId, ac.Config.Settings.UserSettings.Token))
	}

	ac.StartCronJob()
}

func (ac *AppContext) StartCronJob() {
	ac.CronService.StopAllJobs()
	userSettings := ac.Config.Settings.UserSettings
	for bucketName, bucketConfig := range userSettings.Bucket {
		if bucketConfig.IsEnabled {
			ac.CronService.AddJob(fmt.Sprintf("cron-%s", bucketName), "", 1*time.Minute,
				cron.SyncAWEventJob(*ac.AWService, *ac.BackendService, ac.Config, bucketName))
		}
	}
}

func (ac *AppContext) InitializeAWState() {
	if ac.AWService.AWInfo.Hostname != "" {
		userSettings := ac.Config.Settings.UserSettings
		buckets, err := ac.AWService.FetchBuckets()
		if err != nil {
			return
		}
		if userSettings.Bucket == nil {
			userSettings.Bucket = map[string]settings.BucketConfig{}
		}
		for _, bucket := range buckets {
			_, ok := userSettings.Bucket[bucket.ID]
			if !ok {
				userSettings.Bucket[bucket.ID] = settings.BucketConfig{
					IsEnabled: false,
				}
			}
		}

		ac.Config.SettingsManager.SaveSettings(userSettings)
	}
}

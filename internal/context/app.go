package context

import (
	"github.com/shashank-sharma/metadata/internal/config"
	"github.com/shashank-sharma/metadata/internal/cron"
	"github.com/shashank-sharma/metadata/internal/models"
	"github.com/shashank-sharma/metadata/internal/services/activitywatch"
	"github.com/shashank-sharma/metadata/internal/services/backend"
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

package controllers

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/router"
)

type CronController struct {
	router *router.Router
}

func NewCronController(router *router.Router) *CronController {
	return &CronController{
		router: router,
	}
}

func (c *CronController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	title := widget.NewLabel(meta.Title)
	vBox := container.NewVBox()
	vBox.Add(title)

	cronJobs := c.router.AppCtx.CronService.GetJobs()
	if len(cronJobs) == 0 {
		defaultLabel := widget.NewLabel("No Cron job found")
		vBox.Add(defaultLabel)
	} else {
		for _, cronJob := range cronJobs {
			label := widget.NewLabel(fmt.Sprintf("Cron Job: %s", cronJob.CronInfo.Id))
			runningStatus := binding.BoolToStringWithFormat(cronJob.CronInfo.IsRunning, "Running: %s")
			successCount := binding.IntToStringWithFormat(cronJob.CronInfo.SuccessCount, "Success: %s")
			failedCount := binding.IntToStringWithFormat(cronJob.CronInfo.FailedCount, "Failures: %s")

			labelWithData := widget.NewLabelWithData(runningStatus)
			labelNextRun := widget.NewLabelWithData(cronJob.CronInfo.NextRun)
			labelSuccessCount := widget.NewLabelWithData(successCount)
			labelFailedCount := widget.NewLabelWithData(failedCount)

			label.TextStyle = fyne.TextStyle{Bold: true}

			vBox.Add(label)
			vBox.Add(labelWithData)
			vBox.Add(labelNextRun)
			vBox.Add(labelSuccessCount)
			vBox.Add(labelFailedCount)
		}
	}

	return vBox
}

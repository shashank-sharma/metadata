package controllers

import (
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/router"
)

type HomeController struct {
	router *router.Router
}

func (hc *HomeController) GenerateBucketList() fyne.CanvasObject {
	var items []*widget.FormItem

	userSettings := hc.router.AppCtx.Config.Settings.UserSettings
	settingsManager := hc.router.AppCtx.Config.SettingsManager
	bucketMap := userSettings.Bucket
	for bucketName, bucketConfig := range bucketMap {
		bucketName := bucketName
		checkbox := widget.NewCheck(bucketName, func(checked bool) {
			tempBucket := bucketMap[bucketName]
			tempBucket.IsEnabled = checked
			bucketMap[bucketName] = tempBucket
		})
		checkbox.SetChecked(bucketConfig.IsEnabled)

		items = append(items, widget.NewFormItem("", checkbox))
	}

	saveButton := widget.NewButton("Save", func() {
		userSettings.Bucket = bucketMap
		if err := settingsManager.SaveSettings(userSettings); err != nil {
			logger.Error.Println("Failed saving bucketList: ", err)
			hc.router.AppCtx.Notification.Show("Failed saving", "error")
		}
		hc.router.AppCtx.Notification.Show("Settings saved", "info")
		hc.router.AppCtx.StartCronJob()

	})

	form := widget.NewForm(items...)
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), form, saveButton)
}

func (hc *HomeController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	hostname, _ := os.Hostname()
	osName := runtime.GOOS
	arch := runtime.GOARCH
	systemInfo := fmt.Sprintf("Hostname: %s\nOS: %s\nArch: %s", hostname, osName, arch)
	hc.router.AppCtx.Notification.Show("Welcome", "info")

	title := widget.NewLabel(meta.Title)
	info := widget.NewLabel(systemInfo)

	bucketList := hc.GenerateBucketList()

	cronButton := widget.NewButton("Go to CRON", func() {
		hc.router.Navigate(router.CronRoute)
	})
	aboutButton := widget.NewButton("Go to About", func() {
		hc.router.Navigate(router.AboutRoute)
	})
	return container.NewVBox(title, info, widget.NewLabel(meta.Content), bucketList, aboutButton, cronButton)
}

func NewHomeController(r *router.Router) *HomeController {
	return &HomeController{router: r}
}

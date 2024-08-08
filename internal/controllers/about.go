package controllers

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/distatus/battery"
	"github.com/shashank-sharma/metadata/internal/router"
)

type AboutController struct {
	router *router.Router
}

func (ac *AboutController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	title := widget.NewLabel(meta.Title)
	batteryCapacity := widget.NewLabel("")
	batteryFull := widget.NewLabel("")
	backendEndpoint := widget.NewLabel(fmt.Sprintf("Backend: %s", ac.router.AppCtx.Config.BackendEndpoint))
	awEndpoint := widget.NewLabel(fmt.Sprintf("AWEndpoint: %s", ac.router.AppCtx.Config.AWEndpoint))
	batteries, err := battery.GetAll()
	if err != nil {
		ac.router.AppCtx.Notification.Show("Failed to get battery info", "error")
	} else {
		b := batteries[0]
		batteryCapacity.SetText(fmt.Sprintf("Battery current: %f", b.Current))
		batteryFull.SetText(fmt.Sprintf("Battery full: %f", b.Full))
	}

	return container.NewVBox(title, widget.NewLabel(meta.Content),
		batteryCapacity,
		batteryFull,
		backendEndpoint,
		awEndpoint)
}

func NewAboutController(r *router.Router) *AboutController {
	return &AboutController{router: r}
}

package controllers

import (
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/router"
)

type HomeController struct {
	router *router.Router
}

func (hc *HomeController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	hostname, _ := os.Hostname()
	osName := runtime.GOOS
	arch := runtime.GOARCH
	systemInfo := fmt.Sprintf("Hostname: %s\nOS: %s\nArch: %s", hostname, osName, arch)

	title := widget.NewLabel(meta.Title)
	info := widget.NewLabel(systemInfo)

	button := widget.NewButton("Go to About", func() {
		hc.router.Navigate(router.AboutRoute)
	})
	return container.NewVBox(title, info, widget.NewLabel(meta.Content), button)
}

func NewHomeController(r *router.Router) *HomeController {
	return &HomeController{router: r}
}

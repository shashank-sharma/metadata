package controllers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/router"
)

type AboutController struct {
	router *router.Router
}

func (ac *AboutController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	title := widget.NewLabel(meta.Title)
	button := widget.NewButton("Go to Home", func() {
		ac.router.Navigate(router.HomeRoute)
	})
	return container.NewVBox(title, widget.NewLabel(meta.Content), button)
}

func NewAboutController(r *router.Router) *AboutController {
	return &AboutController{router: r}
}

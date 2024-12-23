package controllers

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/router"
)

func NewToolbar(r *router.Router) fyne.CanvasObject {
	// theme.NewColoredResource
	accountTheme := theme.ColorNameError
	return widget.NewToolbar(widget.NewToolbarAction(theme.HomeIcon(), func() {
		fmt.Println("Home")
		r.Navigate(router.HomeRoute)
	}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() { fmt.Println("Cut") }),
		widget.NewToolbarAction(theme.NewColoredResource(theme.AccountIcon(), accountTheme), func() { fmt.Println("Copy") }),
		widget.NewToolbarAction(theme.HistoryIcon(), func() {
			r.Navigate(router.FocusRoute)
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			r.Navigate(router.LoggerRoute)
		}),
	)
}

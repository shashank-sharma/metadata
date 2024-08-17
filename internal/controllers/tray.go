package controllers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/router"
)

func MakeTray(a fyne.App, w fyne.Window, r *router.Router) {
	icon, err := fyne.LoadResourceFromPath("tray.png")
	if err != nil {
		logger.LogError("Could not load tray icon: ", err)
	}

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(icon)

		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))
		desk.SetSystemTrayMenu(m)
		w.SetCloseIntercept(func() {
			logger.LogDebug("Closed intercept")
			w.Hide()
		})
	}
}

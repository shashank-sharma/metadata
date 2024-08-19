package controllers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/router"
)

func MakeTray(a fyne.App, w fyne.Window, r *router.Router) {
	if desk, ok := a.(desktop.App); ok {
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

package router

import (
	"fyne.io/fyne/v2"
	"github.com/shashank-sharma/metadata/internal/logger"
)

func (r *Router) Navigate(route Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.LogDebug("Route: ", route)
	if r.AppCtx.Config.Settings.UserSettings.Token == "" && route != "login" {
		return
	}

	meta, ok := r.metadata[route]
	if !ok {
		return
	}
	controller, exists := r.controllers[route]
	if !exists {
		return
	}

	screen := controller.Screen(meta)
	r.navStack.Objects = []fyne.CanvasObject{}
	r.navStack.Add(screen)
	r.navStack.Refresh()
}

func (r *Router) GetNavStack() *fyne.Container {
	return r.navStack
}

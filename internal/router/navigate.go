package router

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/shashank-sharma/metadata/internal/logger"
)

func (r *Router) Navigate(route Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Debug.Println("Route: ", route)

	meta, ok := r.metadata[route]
	if !ok {
		fmt.Printf("Route not found: %s\n", route)
		return
	}
	controller, exists := r.controllers[route]
	if !exists {
		fmt.Printf("No controller registered for route: %s\n", route)
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

package router

import "fyne.io/fyne/v2"

// Routable represents the interface that different screens/controllers need to implement.
type Routable interface {
	Screen(route RouteMetadata) fyne.CanvasObject
}

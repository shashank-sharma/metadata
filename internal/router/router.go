package router

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/shashank-sharma/metadata/internal/context"
	"github.com/shashank-sharma/metadata/internal/models"
)

type Route string

const (
	HomeRoute   Route = "home"
	AboutRoute  Route = "about"
	LoginRoute  Route = "login"
	CronRoute   Route = "cron"
	LoggerRoute Route = "logger"
	FocusRoute  Route = "focus"
)

type RouteMetadata struct {
	Title   string
	Content string
	State   models.BaseState
}

type Router struct {
	AppCtx      *context.AppContext
	controllers map[Route]Routable
	navStack    *fyne.Container
	metadata    map[Route]RouteMetadata
	mu          sync.Mutex
}

func NewRouter(appCtx *context.AppContext) *Router {
	router := &Router{
		AppCtx:      appCtx,
		controllers: make(map[Route]Routable),
		navStack:    container.NewStack(),
		metadata:    make(map[Route]RouteMetadata),
	}
	return router
}

func (r *Router) RegisterRoute(route Route, metadata RouteMetadata, controller Routable) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metadata[route] = metadata
	r.controllers[route] = controller
}

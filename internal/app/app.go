package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/shashank-sharma/metadata/internal/config"
	"github.com/shashank-sharma/metadata/internal/context"
	"github.com/shashank-sharma/metadata/internal/controllers"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/router"
)

type Application struct {
	RootApp    fyne.App
	RootWindow fyne.Window
	Context    *context.AppContext
	Config     config.AppConfig
	Router     *router.Router
}

func New() (*Application, error) {
	application := &Application{}
	application.RootApp = app.NewWithID("dashboard-metadata")
	application.RootWindow = application.RootApp.NewWindow("Metadata Manager")
	appConfig, err := config.LoadConfig(application.RootApp.Storage().RootURI())

	if err != nil {
		return nil, fmt.Errorf("Failed to load settings: %s", err)
	}
	application.Config = appConfig

	if application.Config.BackendEndpoint == "" || application.Config.AWEndpoint == "" {
		return nil, fmt.Errorf("Endpoint missing")
	}
	application.Context = context.NewAppContext(application.Config)
	application.Router = router.NewRouter(application.Context)

	return application, nil
}

func (app *Application) logLifecycle() {
	app.RootApp.Lifecycle().SetOnStarted(func() {
		logger.LogDebug("Lifecycle: Started")
	})
	app.RootApp.Lifecycle().SetOnStopped(func() {
		logger.LogDebug("Lifecycle: Stopped")
		app.Context.CronService.StopAllJobs()
	})
	/*
		app.RootApp.Lifecycle().SetOnEnteredForeground(func() {
			logger.Debug.Println("Lifecycle: Entered Foreground")
		})
		app.RootApp.Lifecycle().SetOnExitedForeground(func() {
			logger.Debug.Println("Lifecycle: Exited Foreground")
		})
	*/
}

func (app *Application) RegisterRoute() {
	app.Router.RegisterRoute(router.HomeRoute, router.RouteMetadata{
		Title:   "Home",
		Content: "Welcome to the home page!",
	}, controllers.NewHomeController(app.Router))

	app.Router.RegisterRoute(router.AboutRoute, router.RouteMetadata{
		Title:   "About",
		Content: "Created by me",
	}, controllers.NewAboutController(app.Router))

	app.Router.RegisterRoute(router.LoginRoute, router.RouteMetadata{
		Title:   "Login me",
		Content: "Created by login",
	}, controllers.NewLoginController(app.Router))

	app.Router.RegisterRoute(router.CronRoute, router.RouteMetadata{
		Title:   "CRON Jobs",
		Content: "Created by login",
	}, controllers.NewCronController(app.Router))

	app.Router.RegisterRoute(router.LoggerRoute, router.RouteMetadata{
		Title:   "Logger",
		Content: "Created by login",
	}, controllers.NewLoggerController(app.Router))
}

func (app *Application) Render() {
	toolbar := controllers.NewToolbar(app.Router)

	content := container.NewBorder(toolbar, nil, nil, nil, app.Router.GetNavStack())
	fullContent := container.NewBorder(app.Context.Notification.Container(), nil, nil, nil, content)
	app.RootWindow.SetContent(fullContent)

	if app.Config.Settings.UserSettings.Token == "" {
		app.Router.Navigate(router.LoginRoute)
	} else {
		app.Context.OnLoginSuccess()
		app.Router.Navigate(router.HomeRoute)
	}
	app.RootWindow.Resize(fyne.NewSize(520, 520))
	app.RootWindow.ShowAndRun()
}

func (app *Application) Start() error {
	app.logLifecycle()
	app.RegisterRoute()
	controllers.MakeTray(app.RootApp)
	app.Render()

	return nil
}

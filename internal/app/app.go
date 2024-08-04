package app

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/distatus/battery"
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
		return nil, fmt.Errorf("Failed to load settings")
	}
	application.Config = appConfig

	if application.Config.BackendEndpoint == "" || application.Config.AWEndpoint == "" {
		return nil, fmt.Errorf("Endpoint missing")
	}
	application.Context = context.NewAppContext(application.Config)
	application.Router = router.NewRouter(application.Context)

	return application, nil
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		logger.Debug.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		logger.Debug.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		logger.Debug.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		logger.Debug.Println("Lifecycle: Exited Foreground")
	})
}

func (app *Application) InitCronJobs() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if false {
					logger.Debug.Println("Running cron ping")
					// component.PingOnlineStatus(s)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
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
		Title:   "Login",
		Content: "Created by login",
	}, controllers.NewLoginController(app.Router))
}

func (app *Application) Start() error {
	// ui.MakeTray(app.RootApp)
	logLifecycle(app.RootApp)
	// app.RootWindow.SetMaster()
	app.RegisterRoute()

	content := container.NewBorder(app.Context.Notification.Container(), nil, nil, nil, app.Router.GetNavStack())
	app.RootWindow.SetContent(content)

	if app.Config.Settings.UserSettings.Token == "" {
		app.Router.Navigate(router.LoginRoute)
	} else {
		app.Router.Navigate(router.HomeRoute)
	}
	app.RootWindow.Resize(fyne.NewSize(520, 520))

	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return nil
	}
	for i, battery := range batteries {
		fmt.Printf("Bat%d: ", i)
		fmt.Printf("state: %s, ", battery.State.String())
		fmt.Printf("current capacity: %f mWh, ", battery.Current)
		fmt.Printf("last full capacity: %f mWh, ", battery.Full)
		fmt.Printf("design capacity: %f mWh, ", battery.Design)
		fmt.Printf("charge rate: %f mW, ", battery.ChargeRate)
		fmt.Printf("voltage: %f V, ", battery.Voltage)
		fmt.Printf("design voltage: %f V\n", battery.DesignVoltage)
	}

	app.RootWindow.ShowAndRun()

	return nil
}

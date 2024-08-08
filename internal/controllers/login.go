package controllers

import (
	"fmt"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/router"
	"github.com/shashank-sharma/metadata/internal/theme"
)

type LoginController struct {
	router *router.Router
}

func NewLoginController(router *router.Router) *LoginController {
	return &LoginController{router: router}
}

func (lc *LoginController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	title := theme.NewCenteredTitle(meta.Title)
	productEntry := widget.NewEntry()
	productEntry.SetPlaceHolder("Product Name")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	backendEndpoint := widget.NewLabel(fmt.Sprintf("Backend: %s", lc.router.AppCtx.Config.BackendEndpoint))
	awEndpoint := widget.NewLabel(fmt.Sprintf("AWEndpoint: %s", lc.router.AppCtx.Config.AWEndpoint))

	lc.router.AppCtx.Notification.Show("Login required for dev token", "info")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: productEntry},
			{Text: "Email", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		},
		OnSubmit: func() {
			productName := productEntry.Text
			username := usernameEntry.Text
			password := passwordEntry.Text
			loginResponse, err := lc.router.AppCtx.BackendService.Login(username, password)
			if err != nil {
				lc.router.AppCtx.Notification.Show("Failed to login", "error")
				return
			}

			lc.router.AppCtx.BackendService.Client.SetToken(loginResponse.Token)

			hostname, _ := os.Hostname()
			osName := runtime.GOOS
			arch := runtime.GOARCH

			devToken, err := lc.router.AppCtx.BackendService.GetDevToken()

			if err != nil {
				lc.router.AppCtx.Notification.Show("Failed fetching dev token", "error")
				return
			}

			productId, err := lc.router.AppCtx.BackendService.SetTrackingDevice(productName, hostname, osName, arch)
			if err != nil {
				lc.router.AppCtx.Notification.Show("Failed to register tracking device", "error")
				return
			}

			userSettings := lc.router.AppCtx.Config.Settings.UserSettings
			userSettings.UserId = loginResponse.Record.ID
			userSettings.ProductId = productId
			userSettings.Token = devToken

			lc.router.AppCtx.Config.SettingsManager.SaveSettings(userSettings)

			lc.router.AppCtx.Notification.Show("Welcome", "info")
			lc.router.Navigate(router.HomeRoute)
		},
	}

	formContainer := container.NewPadded(form)
	content := container.NewVBox(
		title,
		formContainer,
		backendEndpoint,
		awEndpoint,
	)

	return container.NewVBox(
		container.New(layout.NewVBoxLayout(), layout.NewSpacer(), content, layout.NewSpacer()),
	)
}

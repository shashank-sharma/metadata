package component

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/settings"
)

func MakeToolbarTab(setting *settings.Settings, window fyne.Window) fyne.CanvasObject {
	// theme.NewColoredResource
	accountTheme := theme.ColorNameError
	if setting.GetToken() != "" {
		accountTheme = theme.ColorNameSuccess
	}
	t := widget.NewToolbar(widget.NewToolbarAction(theme.HomeIcon(), func() {
		fmt.Println("Home")
		GenerateHomePage()
	}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() { fmt.Println("Cut") }),
		widget.NewToolbarAction(theme.NewColoredResource(theme.AccountIcon(), accountTheme), func() { fmt.Println("Copy") }),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fmt.Println("Settings")
			GenerateSettingsPage(setting, window)
		}),
	)

	return container.NewBorder(t, nil, nil, nil)
}

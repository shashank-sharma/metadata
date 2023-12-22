package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/shashank-sharma/metadata/component"
	"github.com/shashank-sharma/metadata/settings"
	"github.com/shashank-sharma/metadata/util"
)

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func runCronJobs(s *settings.Settings) {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if s.ShouldSync() {
					fmt.Println("RUnning cron ping")
					component.PingOnlineStatus(s)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	util.Init()
	a := util.GetApp()
	s := util.GetSetting()
	w := util.GetWindow()
	makeTray(a)
	logLifecycle(a)
	w.SetMaster()

	s.LoadAppearanceScreen(w)
	component.GenerateHomePage()

	w.Resize(fyne.NewSize(520, 520))

	runCronJobs(s)
	w.ShowAndRun()
}

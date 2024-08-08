package controllers

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/router"
)

type CronTask struct {
	Label    string
	Progress *widget.ProgressBar
}

type CronController struct {
	router *router.Router
	tasks  []CronTask
	mu     sync.Mutex
}

func NewCronController(router *router.Router) *CronController {
	return &CronController{
		router: router,
		tasks:  make([]CronTask, 10),
	}
}

func (c *CronController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	title := widget.NewLabel(meta.Title)
	c.mu.Lock()
	defer c.mu.Unlock()
	vBox := container.NewVBox()
	vBox.Add(title)

	for i := 0; i < 10; i++ {
		label := widget.NewLabel(fmt.Sprintf("Cron Job #%d", i+1))
		progressBar := widget.NewProgressBar()
		progressBar.Min = 0
		progressBar.Max = 60
		c.tasks[i] = CronTask{Label: fmt.Sprintf("Cron Job #%d", i+1), Progress: progressBar}

		vBox.Add(label)
		vBox.Add(progressBar)
	}

	go c.startTicking()

	return vBox
}

func (c *CronController) startTicking() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for _, task := range c.tasks {
				value := task.Progress.Value
				if value--; value < 0 {
					value = task.Progress.Max
				}
				task.Progress.SetValue(value)
			}
			c.mu.Unlock()
		}
	}
}

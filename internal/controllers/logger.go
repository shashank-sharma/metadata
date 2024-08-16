package controllers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/router"
)

type LoggerController struct {
	router *router.Router
}

func (lc *LoggerController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	logColors := map[string]color.Color{
		"DEBUG":   color.NRGBA{R: 0, G: 128, B: 128, A: 255},
		"INFO":    color.NRGBA{R: 0, G: 0, B: 255, A: 255},
		"WARNING": color.NRGBA{R: 255, G: 255, B: 0, A: 255},
		"ERROR":   color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}

	logEntries := logger.RetrieveLogs()
	logList := container.NewVBox()

	for _, entry := range logEntries {
		textColor, ok := logColors[entry.Level]
		if !ok {
			textColor = color.White
		}

		text := canvas.NewText(entry.Level+": "+entry.Message, textColor)
		text.Alignment = fyne.TextAlignLeading
		text.TextStyle = fyne.TextStyle{Bold: false}
		text.TextSize = 14

		logList.Add(text)
	}

	scrollContainer := container.NewScroll(logList)
	scrollContainer.SetMinSize(fyne.NewSize(480, 320))

	background := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	background.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	background.StrokeColor = color.Transparent
	background.StrokeWidth = 0

	return container.NewStack(background, scrollContainer)
}

func NewLoggerController(r *router.Router) *LoggerController {
	return &LoggerController{router: r}
}

package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func NewTitle(text string) *fyne.Container {
	title := canvas.NewText(text, nil)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 24

	titleContainer := container.NewStack()
	titleContainer.Add(title)
	return titleContainer
}

func NewCenteredTitle(text string) *fyne.Container {
	titleContainer := NewTitle(text)
	return container.NewHBox(layout.NewSpacer(), titleContainer, layout.NewSpacer())
}

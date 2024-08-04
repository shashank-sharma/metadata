package theme

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Notification struct {
	message    *canvas.Text
	background *canvas.Rectangle
	container  *fyne.Container
	animation  *fyne.Animation
}

func NewNotification() *Notification {
	message := canvas.NewText("", color.White)
	message.Alignment = fyne.TextAlignCenter
	background := canvas.NewRectangle(color.RGBA{200, 200, 200, 255})

	container := container.NewStack()
	container.Hide()

	return &Notification{
		message:    message,
		background: background,
		container:  container,
	}
}

func (n *Notification) animateOut(duration time.Duration) {
	n.animation = canvas.NewPositionAnimation(
		fyne.NewPos(0, 0), fyne.NewPos(0, -n.container.MinSize().Height-10),
		duration, func(p fyne.Position) {
			n.container.Move(p)
			if p == fyne.NewPos(0, -n.container.MinSize().Height) {
				// n.Hide()
				n.container.Refresh()
			}
		})
	n.animation.Start()
}

func (n *Notification) Show(text string, category string) {
	duration := 4 * time.Second
	n.message.Text = text
	if category == "error" {
		n.background.FillColor = color.RGBA{255, 0, 0, 130}
	} else if category == "info" {
		n.background.FillColor = color.RGBA{0, 240, 0, 130}
	} else {
		n.background.FillColor = color.RGBA{0, 0, 0, 130}

	}

	n.container.Objects = []fyne.CanvasObject{n.background, n.message}
	n.container.Move(fyne.NewPos(0, 0))
	n.message.Resize(n.container.Size())

	n.message.Refresh()
	n.container.Show()

	go func() {
		time.Sleep(duration)
		n.animateOut(time.Second)
	}()
}

func (n *Notification) Hide() {
	n.container.Hide()
}

func (n *Notification) Container() *fyne.Container {
	return n.container
}

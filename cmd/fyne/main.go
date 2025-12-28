package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const (
	greenR       = 46
	greenG       = 125
	greenB       = 50
	greenA       = 255
	windowWidth  = 1200
	windowHeight = 600
)

func main() {
	myApp := app.New()

	window := myApp.NewWindow("Irate Sol")

	// Green Rectangle.
	rect := canvas.NewRectangle(color.RGBA{R: greenR, G: greenG, B: greenB, A: greenA})
	window.SetContent(rect)

	width, height := float32(windowWidth), float32(windowHeight)
	window.Resize(fyne.NewSize(width, height))
	window.ShowAndRun()
}

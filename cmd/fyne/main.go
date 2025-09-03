package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	myApp := app.New()

	w := myApp.NewWindow("Irate Sol")

	// Green Rectangle.
	rect := canvas.NewRectangle(color.RGBA{R: 46, G: 125, B: 50, A: 255})
	w.SetContent(rect)

	width, height := float32(1200), float32(600)
	w.Resize(fyne.NewSize(width, height))
	w.ShowAndRun()
}

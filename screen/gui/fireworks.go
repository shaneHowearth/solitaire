package gui

import (
	"image/color"
	"math"
	"math/rand/v2"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type particle struct {
	shape *canvas.Circle
	velX  float32
	velY  float32
	clr   color.RGBA
}

func (d *Display) LaunchFireworks() {
	const maxSparks = 100
	const sparkRadius = 4.0

	sparkLayer := container.NewWithoutLayout()
	sparks := make([]*canvas.Circle, maxSparks)
	velocities := make([]fyne.Position, maxSparks)

	for i := 0; i < maxSparks; i++ {
		c := canvas.NewCircle(color.Transparent)
		c.Resize(fyne.NewSize(sparkRadius*2, sparkRadius*2))
		c.Hide()
		sparks[i] = c
		sparkLayer.Add(c)
	}

	d.CardLayer.Add(sparkLayer)

	// Fyne's Animation constructor
	anim := fyne.NewAnimation(time.Second*2, func(v float32) {
		// 1. Cleanup check: If animation is finished, remove layer
		if v >= 1.0 {
			d.CardLayer.Remove(sparkLayer)
			return
		}

		// 2. Burst trigger
		if int(v*10)%2 == 0 {
			size := d.Window.Canvas().Size()
			ox := 150 + rand.Float32()*(size.Width-300)
			oy := 150 + rand.Float32()*(size.Height-300)

			for i := 0; i < 20; i++ {
				idx := (int(v*10)/2)*20 + i
				if idx >= maxSparks {
					break
				}

				angle := rand.Float64() * 2 * math.Pi
				speed := 2.0 + rand.Float64()*5.0

				sparks[idx].FillColor = color.RGBA{uint8(rand.IntN(155) + 100), uint8(rand.IntN(155) + 100), uint8(rand.IntN(155) + 100), 255}
				sparks[idx].Move(fyne.NewPos(ox, oy))
				sparks[idx].Show()

				velocities[idx] = fyne.NewPos(float32(math.Cos(angle)*speed), float32(math.Sin(angle)*speed))
			}
		}

		// 3. Physics update
		for i := 0; i < maxSparks; i++ {
			if !sparks[i].Visible() {
				continue
			}

			pos := sparks[i].Position()
			velocities[i].Y += 0.15
			sparks[i].Move(fyne.NewPos(pos.X+velocities[i].X, pos.Y+velocities[i].Y))

			// Fading: modify alpha
			fade := uint8(255 * (1.0 - v))
			c := sparks[i].FillColor.(color.RGBA)
			c.A = fade
			sparks[i].FillColor = c
			sparks[i].Refresh()
		}
	})

	anim.Start()
}

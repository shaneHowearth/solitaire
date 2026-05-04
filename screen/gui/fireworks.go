package gui

import (
	"image/color"
	"math"
	"math/rand/v2"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type particle struct {
	shape *canvas.Circle
	velX  float32
	velY  float32
	clr   color.RGBA
}

func (d *Display) LaunchFireworks() {
	type spark struct {
		x, y   float32
		vx, vy float32
		clr    color.RGBA
		active bool
	}

	// We'll use a fixed pool of sparks to avoid memory allocation during animation
	const maxSparks = 100
	sparks := make([]spark, maxSparks)
	var activeCount int

	// Create a Raster: it calls a function to determine the color of every pixel
	// This is the fastest way to do custom drawing in Fyne
	raster := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		// Default background for the raster is transparent
		for i := 0; i < activeCount; i++ {
			s := sparks[i]
			if !s.active {
				continue
			}
			// If pixel (x,y) is within the spark radius, color it
			// Note: We use simple square-bounds check for speed
			dx := float32(x) - s.x
			dy := float32(y) - s.y
			if dx*dx+dy*dy < 16 { // 4px radius
				return s.clr
			}
		}
		return color.Transparent
	})

	d.CardLayer.Add(raster)
	raster.Hide() // Hide until first burst

	// Simple physics loop
	anim := fyne.NewAnimation(time.Second*3, func(v float32) {
		if v == 0 {
			raster.Show()
		}

		// 1. Every few frames, "spawn" a new burst if needed
		if int(v*100)%20 == 0 && activeCount < maxSparks-20 {
			size := d.Window.Canvas().Size()
			ox := 150 + rand.Float32()*(size.Width-300)
			oy := 150 + rand.Float32()*(size.Height-300)

			for i := 0; i < 20; i++ {
				angle := rand.Float64() * 2 * math.Pi
				speed := 2.0 + rand.Float64()*5.0
				sparks[activeCount] = spark{
					x: ox, y: oy,
					vx:     float32(math.Cos(angle) * speed),
					vy:     float32(math.Sin(angle) * speed),
					clr:    color.RGBA{uint8(rand.IntN(155) + 100), uint8(rand.IntN(155) + 100), uint8(rand.IntN(155) + 100), 255},
					active: true,
				}
				activeCount++
			}
		}

		// 2. Update positions
		for i := 0; i < activeCount; i++ {
			if !sparks[i].active {
				continue
			}
			sparks[i].x += sparks[i].vx
			sparks[i].y += sparks[i].vy
			sparks[i].vy += 0.15 // Gravity
			sparks[i].clr.A = uint8(255 * (1.0 - v))
		}

		raster.Refresh()

		if v == 1.0 {
			d.CardLayer.Remove(raster)
		}
	})

	anim.Start()
}

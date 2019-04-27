package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	backingColour = color.RGBA{0, 0, 0, 0}
	winBounds = pixel.R(0, 0, 1280, 720)
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title: "Pay You Way",
		Bounds: winBounds,
		VSync: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	lvlMan := NewLevelManager(winBounds)

	last := time.Now()
	second := time.Tick(time.Second)
	frames := 0

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(backingColour)

		lvlMan.Update(dt)
		lvlMan.Draw(win)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

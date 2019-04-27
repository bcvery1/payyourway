package main

import (
	"fmt"
	"image/color"
	"time"

	_ "image/png"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	backingColour = color.RGBA{0, 0, 0, 0}
	winBounds = pixel.R(0, 0, 1280, 720)
	tmxMap *tilepix.Map
)

func run() {
	var err error
	tmxMap, err = tilepix.ReadFile("assets/map.tmx")
	if err != nil {
		panic(err)
	}


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

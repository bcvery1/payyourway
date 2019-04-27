package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"time"

	_ "image/png"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	backingColour = color.RGBA{0, 0, 0, 0}
	winBounds = pixel.R(0, 0, 1920, 1080)
	tmxMap *tilepix.Map

	speed = 1280. /4

	camPos = pixel.ZV
	camZoom = 1.0

	player *Player
	lvlMan *LevelManager
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

	lvlMan = NewLevelManager(winBounds)
	player = NewPlayer()

	last := time.Now()
	second := time.Tick(time.Second)
	frames := 0

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(winBounds.Center().Sub(camPos))
		win.SetMatrix(cam)
		camZoom *= math.Pow(1.2, win.MouseScroll().Y)

		win.Clear(backingColour)

		player.Update(dt, cam.Unproject(winBounds.Center()))
		lvlMan.Update(dt, win)

		lvlMan.Draw(win)
		player.Draw(win)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | %v | %.2f", cfg.Title, frames, camPos, camZoom))
			frames = 0
		default:
		}
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main() {
	pixelgl.Run(run)
}

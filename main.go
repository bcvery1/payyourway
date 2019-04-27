package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func init() {
	var err error
	tilemapPic, err = loadPicture("assets/tilemap.png")
	if err != nil {
		panic(err)
	}

	rocketSprite = pixel.NewSprite(tilemapPic, pixel.R(0, 31*16, 16, 32*16))

	fireSprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, pixel.R(0, 32*16, 16, 33*16)),
		pixel.NewSprite(tilemapPic, pixel.R(16, 32*16, 32, 33*16)),
		//pixel.NewSprite(tilemapPic, pixel.R()),
		//pixel.NewSprite(tilemapPic, pixel.R()),
	}

	smokeSprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, pixel.R(64, 32*16, 80, 33*16)),
		pixel.NewSprite(tilemapPic, pixel.R(80, 32*16, 96, 33*16)),
		//pixel.NewSprite(tilemapPic, pixel.R()),
		//pixel.NewSprite(tilemapPic, pixel.R()),
	}

	rand.Seed(time.Now().UnixNano())
}

var (
	backingColour = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	winBounds     = pixel.R(0, 0, 1024, 720)
	tmxMap        *tilepix.Map

	speed = 1280. / 6

	camPos  = pixel.ZV
	camZoom = 1.0

	player *Player
	lvlMan *LevelManager

	atlas *text.Atlas

	tilemapPic pixel.Picture
)

func run() {
	var err error
	tmxMap, err = tilepix.ReadFile("assets/map.tmx")
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Pay Your Way",
		Bounds: winBounds,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	lvlMan = NewLevelManager(winBounds)
	lvlMan.StartLevel(Level1Ind)
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

		UpdateRockets(dt)
		DrawRockets(win)

		UpdateFires(dt)
		DrawFires(win)

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

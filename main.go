package main

import (
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font/basicfont"
)

func init() {
	logrus.SetLevel(logrus.FatalLevel)

	var err error
	tilemapPic, err = loadPicture("assets/tilemap.png")
	if err != nil {
		panic(err)
	}

	rocketSprite = pixel.NewSprite(tilemapPic, pixel.R(0, 31*16, 16, 32*16))

	fireSprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, pixel.R(0, 32*16, 16, 33*16)),
		pixel.NewSprite(tilemapPic, pixel.R(16, 32*16, 32, 33*16)),
		pixel.NewSprite(tilemapPic, pixel.R(32, 32*16, 48, 33*16)),
		pixel.NewSprite(tilemapPic, pixel.R(48, 32*16, 64, 33*16)),
	}

	smokeSprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, pixel.R(0, 33*16, 16, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(16, 33*16, 32, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(32, 33*16, 48, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(48, 33*16, 64, 34*16)),
	}

	enemySprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, pixel.R(0, 16, 16, 32)),
		pixel.NewSprite(tilemapPic, pixel.R(16, 16, 32, 32)),
		pixel.NewSprite(tilemapPic, pixel.R(32, 16, 48, 32)),
		pixel.NewSprite(tilemapPic, pixel.R(48, 16, 64, 32)),
	}

	flareSprites = []*pixel.Sprite{
		pixel.NewSprite(tilemapPic, pixel.R(64, 33*16, 80, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(80, 33*16, 96, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(96, 33*16, 112, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(112, 33*16, 128, 34*16)),
		pixel.NewSprite(tilemapPic, pixel.R(128, 33*16, 144, 34*16)),
	}

	invSprites = map[string]*pixel.Sprite{
		"Boat":          pixel.NewSprite(tilemapPic, pixel.R(0, 30*16, 16, 31*16)),
		"Flares":        pixel.NewSprite(tilemapPic, pixel.R(16, 30*16, 32, 31*16)),
		"Invincibility": pixel.NewSprite(tilemapPic, pixel.R(64, 33*16, 80, 34*16)),
	}

	powerupSprites = map[action]*pixel.Sprite{
		"health":    pixel.NewSprite(tilemapPic, pixel.R(32, 30*16, 48, 31*16)),
		"maxhealth": pixel.NewSprite(tilemapPic, pixel.R(32, 30*16, 48, 31*16)),
		"shield":    pixel.NewSprite(tilemapPic, pixel.R(48, 30*16, 64, 31*16)),
	}

	rand.Seed(time.Now().UnixNano())
}

var (
	backingColour = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	winBounds     = pixel.R(0, 0, 1024, 720)
	tmxMap        *tilepix.Map

	speed = 1280. / 6

	camPos = pixel.ZV

	player *Player
	lvlMan *LevelManager

	atlas *text.Atlas

	tilemapPic pixel.Picture

	invSprites map[string]*pixel.Sprite
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

	_ = tmxMap.DrawAll(win, color.Transparent, pixel.IM)

	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	lvlMan = NewLevelManager(winBounds)
	lvlMan.StartLevel(MenuInd)
	player = NewPlayer()

	SetupEnemies()
	SetupGuns()
	SetupAudio()
	SetupAnnouncements()
	SetupPowerups()

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(winBounds.Center().Sub(camPos))
		win.SetMatrix(cam)

		win.Clear(backingColour)

		player.Update(dt, cam.Unproject(winBounds.Center()))
		player.UpdateInventory(win)
		lvlMan.Update(dt, win)

		lvlMan.Draw(win)

		UpdateGuns(dt)
		DrawGuns(win)

		UpdateFires(dt)
		DrawFires(win)

		UpdateEnemies(dt)
		DrawEnemies(win)

		UpdateFlares(dt)
		DrawFlares(win)

		UpdateAnnouncements(dt)
		DrawAnnouncements(win)

		UpdatePowerups()
		DrawPowerups(win)

		win.Update()
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

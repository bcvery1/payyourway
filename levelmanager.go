package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	MenuInd = iota
	Level1Ind
	Level2Ind
	Level3Ind
	Level4Ind
	ShopInd
	EndInd
)

type Level interface {
	Init(pixel.Rect)
	Start()
	Update(float64, *pixelgl.Window)
	Draw(*pixelgl.Window)
	Collides(pixel.Rect) bool
	Hurt(pixel.Rect)
	ReachedShop() string
}

type LevelManager struct {
	currentLevel  int
	levels        []Level
	previousLevel int
}

func NewLevelManager(bounds pixel.Rect) *LevelManager {
	lm := LevelManager{
		currentLevel: 0,
		levels: []Level{
			&Menu{},
			&Level1{},
			&Level2{},
			&Level3{},
			&Level4{},
			&Shop{},
			&End{},
		},
	}

	for _, lvl := range lm.levels {
		lvl.Init(bounds)
	}

	return &lm
}

func (lm *LevelManager) Update(dt float64, win *pixelgl.Window) {
	lm.CurrentLevel().Update(dt, win)
}

func (lm *LevelManager) Draw(win *pixelgl.Window) {
	lm.CurrentLevel().Draw(win)
}

func (lm *LevelManager) CurrentLevel() Level {
	return lm.levels[lm.currentLevel]
}

func (lm *LevelManager) StartLevel(index int) {
	lm.previousLevel = lm.currentLevel
	lm.currentLevel = index
	lm.CurrentLevel().Start()

	rockets = make(map[int]*rocket)
	fires = make(map[int]*fire)
}

func (lm *LevelManager) RestartLevel() {
	lm.StartLevel(lm.currentLevel)
}

func (lm *LevelManager) Shop() *Shop {
	return lm.levels[ShopInd].(*Shop)
}

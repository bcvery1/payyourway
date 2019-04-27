package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level interface {
	Init(pixel.Rect)
	Start()
	Update(float64, *pixelgl.Window)
	Draw(*pixelgl.Window)
	Collides(pixel.Rect) bool
	Hurt(pixel.Rect)
}

type LevelManager struct {
	currentLevel int
	levels []Level
}

func NewLevelManager(bounds pixel.Rect) *LevelManager {
	lm := LevelManager{
		currentLevel: 0,
		levels: []Level{
			//&Menu{},
			&Level1{},
			//&Level2{},
			//&Level3{},
			//&Level4{},
			&Shop{},
		},
	}

	for _, lvl := range lm.levels {
		lvl.Init(bounds)
	}

	lm.levels[0].Start()

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
	lm.currentLevel = index
	lm.CurrentLevel().Start()
}

func (lm *LevelManager) RestartLevel() {
	lm.StartLevel(lm.currentLevel)
}

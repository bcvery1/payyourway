package main

import (
	"github.com/faiface/pixel"
)

type Level interface {
	Init(pixel.Rect)
	Start()
	Update(float64)
	Draw(pixel.Target)
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
			//&Shop{},
		},
	}

	for _, lvl := range lm.levels {
		lvl.Init(bounds)
	}

	return &lm
}

func (lm *LevelManager) Update(dt float64) {
	lm.levels[lm.currentLevel].Update(dt)
}

func (lm *LevelManager) Draw(target pixel.Target) {
	lm.levels[lm.currentLevel].Draw(target)
}

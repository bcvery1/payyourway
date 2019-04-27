package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level1 struct {

}

func (l *Level1) Init(pixel.Rect) {
	// do nothing
}

func (l *Level1) Draw(target pixel.Target) {
	_ = tmxMap.DrawAll(target, backingColour, pixel.IM)
}

func (l *Level1) Start() {
	camPos = pixel.V(3210, 405)
}

func (l *Level1) Update(dt float64, win *pixelgl.Window) {
	newPos := camPos
	if win.Pressed(pixelgl.KeyW) {
		newPos.Y += speed*dt
	}
	if win.Pressed(pixelgl.KeyS) {
		newPos.Y -= speed*dt
	}
	if win.Pressed(pixelgl.KeyA) {
		newPos.X -= speed*dt
	}
	if win.Pressed(pixelgl.KeyD) {
		newPos.X += speed*dt
	}

	if player.CanMove(newPos) {
		camPos = newPos
	}
}


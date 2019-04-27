package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level1 struct {
	collisionRects []pixel.Rect
}

func (l *Level1) Init(pixel.Rect) {
	objLayer := tmxMap.GetObjectLayerByName("Level1Collisions")

	for _, obj := range objLayer.Objects {
		if r, err :=  obj.GetRect(); err == nil {
			l.collisionRects = append(l.collisionRects, r)
		}
	}
}

func (l *Level1) Draw(target pixel.Target) {
	_ = tmxMap.DrawAll(target, backingColour, pixel.IM)
}

func (l *Level1) Start() {
	camPos = pixel.V(3210, 405)
}

func (l *Level1) Update(dt float64, win *pixelgl.Window) {
	deltaPos := pixel.ZV
	if win.Pressed(pixelgl.KeyW) {
		deltaPos.Y += speed*dt
	}
	if win.Pressed(pixelgl.KeyS) {
		deltaPos.Y -= speed*dt
	}
	if player.CanMove(deltaPos) {
		camPos = camPos.Add(deltaPos)
	}

	deltaPos = pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		deltaPos.X -= speed*dt
	}
	if win.Pressed(pixelgl.KeyD) {
		deltaPos.X += speed*dt
	}

	if player.CanMove(deltaPos) {
		camPos = camPos.Add(deltaPos)
	}
}

func (l *Level1) Collides(playerR pixel.Rect) bool {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.collisionRects {
		if r.Intersect(playerR) != zr {
			return true
		}
	}

	return false
}


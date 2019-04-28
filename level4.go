package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	gravity = pixel.V(0, speed/12)
	vel     = pixel.ZV
	jumping = false
)

type Level4 struct {
	collisionRects []pixel.Rect
	mineRects      []pixel.Rect
	gravity        bool
	gravityOn      pixel.Rect
	gravityOff     pixel.Rect
}

func (l *Level4) Init(pixel.Rect) {
	collisionOjLayer := tmxMap.GetObjectLayerByName("Level4Collisions")
	for _, obj := range collisionOjLayer.Objects {
		if r, err := obj.GetRect(); err == nil {
			l.collisionRects = append(l.collisionRects, r)
		}
	}

	for _, obj := range tmxMap.GetObjectLayerByName("CommonCollisions").Objects {
		if r, err := obj.GetRect(); err == nil {
			l.collisionRects = append(l.collisionRects, r)
		}
	}

	mineObjLayer := tmxMap.GetObjectLayerByName("mines")
	for _, obj := range mineObjLayer.Objects {
		if r, err := obj.GetRect(); err == nil {
			l.mineRects = append(l.mineRects, r)
		}
	}

	l.gravityOn, _ = tmxMap.GetObjectByName("platformer")[0].GetRect()
	l.gravityOff, _ = tmxMap.GetObjectByName("2d")[0].GetRect()

	l.gravity = true
}

func (l *Level4) Start() {
	camPos = pixel.V(45, 2262)
	Announce("Level 4")
	Announce("Gravity on")
}

func (l *Level4) Update(dt float64, win *pixelgl.Window) {
	l.checkGravity(player.CollisionBox())

	deltaPos := pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		deltaPos.X -= speed * dt
	}
	if win.Pressed(pixelgl.KeyD) {
		deltaPos.X += speed * dt
	}

	if deltaPos != pixel.ZV && player.CanMove(deltaPos) {
		camPos = camPos.Add(deltaPos)
		l.Hurt(player.CollisionBox().Moved(deltaPos))
	}

	if l.gravity {
		// Vertical
		if !jumping {
			if win.JustPressed(pixelgl.KeyW) {
				vel = pixel.V(0, speed*2*dt*player.jumpBoost)
				jumping = true
			}
		}
		vel = vel.Sub(gravity.Scaled(dt))

		deltaPos = vel
		if deltaPos != pixel.ZV && player.CanMove(deltaPos) {
			camPos = camPos.Add(deltaPos)
			l.Hurt(player.CollisionBox().Moved(deltaPos))
		} else {
			vel = pixel.ZV
			jumping = false
		}
	} else {
		deltaPos := pixel.ZV
		if win.Pressed(pixelgl.KeyS) {
			deltaPos.Y -= speed * dt
		}
		if win.Pressed(pixelgl.KeyW) {
			deltaPos.Y += speed * dt
		}

		if deltaPos != pixel.ZV && player.CanMove(deltaPos) {
			camPos = camPos.Add(deltaPos)
			l.Hurt(player.CollisionBox().Moved(deltaPos))
		}
	}

	// Check if we've reached a shop
	if shopName := l.ReachedShop(); shopName != "" {
		lvlMan.StartLevel(EndInd)
	}
}

func (l *Level4) Draw(win *pixelgl.Window) {
	_ = tmxMap.DrawAll(win, backingColour, pixel.IM)
	player.Draw(win)
}

func (l *Level4) Collides(playerR pixel.Rect) bool {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.collisionRects {
		if r.Intersect(playerR) != zr {
			return true
		}
	}

	return false
}

func (l *Level4) Hurt(playerR pixel.Rect) {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.mineRects {
		if r.Intersect(playerR) != zr {
			player.Hurt(10)
			return
		}
	}

	if playerR.Min.Y < 2240 {
		player.Hurt(player.health)
	}
}

func (l *Level4) ReachedShop() string {
	for _, obj := range tmxMap.GetObjectLayerByName("shops").Objects {

		if r, err := obj.GetRect(); err == nil && r.Intersect(player.CollisionBox()) != pixel.R(0, 0, 0, 0) {
			return obj.Name
		}
	}

	return ""
}

func (l *Level4) checkGravity(playerR pixel.Rect) {
	if l.gravityOn.Intersect(playerR) != pixel.R(0, 0, 0, 0) {
		l.gravity = true
	}
	if l.gravityOff.Intersect(playerR) != pixel.R(0, 0, 0, 0) {
		l.gravity = false
	}
}

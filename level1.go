package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level1 struct {
	collisionRects []pixel.Rect
	mineRects      []pixel.Rect
}

func (l *Level1) Init(pixel.Rect) {
	collisionOjLayer := tmxMap.GetObjectLayerByName("Level1Collisions")
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
}

func (l *Level1) Draw(win *pixelgl.Window) {
	_ = tmxMap.DrawAll(win, backingColour, pixel.IM)
	player.Draw(win)
}

func (l *Level1) Start() {
	camPos = pixel.V(3210, 405)
	Announce("Level 1")
}

func (l *Level1) Update(dt float64, win *pixelgl.Window) {
	deltaPos := pixel.ZV
	if win.Pressed(pixelgl.KeyW) {
		deltaPos.Y += speed * dt
	}
	if win.Pressed(pixelgl.KeyS) {
		deltaPos.Y -= speed * dt
	}
	if deltaPos != pixel.ZV && player.CanMove(deltaPos) {
		camPos = camPos.Add(deltaPos)
		l.Hurt(player.bounds.Moved(player.offSet.Add(deltaPos)))
	}

	deltaPos = pixel.ZV
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

	// Check if we've reached a shop
	if shopName := l.ReachedShop(); shopName != "" {
		lvlMan.StartLevel(ShopInd)
		lvlMan.Shop().Setup(shopName)
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

func (l *Level1) Hurt(playerR pixel.Rect) {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.mineRects {
		if r.Intersect(playerR) != zr {
			player.Hurt(10)
			return
		}
	}
}

func (l *Level1) ReachedShop() string {
	for _, obj := range tmxMap.GetObjectLayerByName("shops").Objects {

		if r, err := obj.GetRect(); err == nil && r.Intersect(player.CollisionBox()) != pixel.R(0, 0, 0, 0) {
			return obj.Name
		}
	}

	return ""
}

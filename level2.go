package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level2 struct {
	collisionRects []pixel.Rect
	mineRects      []pixel.Rect
	water          []pixel.Rect
}

func (l *Level2) Init(pixel.Rect) {
	collisionOjLayer := tmxMap.GetObjectLayerByName("Level2Collisions")
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

	for _, obj := range tmxMap.GetObjectLayerByName("Water").Objects {
		if r, err := obj.GetRect(); err != nil {
			l.water = append(l.water, r)
		}
	}
}

func (l *Level2) Start() {
	camPos = pixel.V(7330, 4660)
	Announce("Level 2")
}

func (l *Level2) Update(dt float64, win *pixelgl.Window) {
	deltaPos := pixel.ZV
	if win.Pressed(pixelgl.KeyW) {
		deltaPos.Y += speed * dt
	}
	if win.Pressed(pixelgl.KeyS) {
		deltaPos.Y -= speed * dt
	}
	if deltaPos != pixel.ZV && player.CanMove(deltaPos) {
		newPlayerPos := player.bounds.Moved(player.offSet.Add(deltaPos))

		camPos = camPos.Add(deltaPos)
		l.Hurt(newPlayerPos)
		l.Drown(newPlayerPos)
	}

	deltaPos = pixel.ZV
	if win.Pressed(pixelgl.KeyA) {
		deltaPos.X -= speed * dt
	}
	if win.Pressed(pixelgl.KeyD) {
		deltaPos.X += speed * dt
	}

	if deltaPos != pixel.ZV && player.CanMove(deltaPos) {
		newPlayerPos := player.bounds.Moved(player.offSet.Add(deltaPos))

		camPos = camPos.Add(deltaPos)
		l.Hurt(newPlayerPos)
		l.Drown(newPlayerPos)
	}

	// Check if we've reached a shop
	if shopName := l.ReachedShop(); shopName != "" {
		lvlMan.StartLevel(ShopInd)
		lvlMan.Shop().Setup(shopName)
	}
}

func (l *Level2) Draw(win *pixelgl.Window) {
	_ = tmxMap.DrawAll(win, backingColour, pixel.IM)
	player.Draw(win)
}

func (l *Level2) Collides(playerR pixel.Rect) bool {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.collisionRects {
		if r.Intersect(playerR) != zr {
			return true
		}
	}

	return false
}

func (l *Level2) Hurt(playerR pixel.Rect) {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.mineRects {
		if r.Intersect(playerR) != zr {
			player.Hurt(10)
			return
		}
	}
}

func (l *Level2) Drown(playerR pixel.Rect) {
	zr := pixel.R(0, 0, 0, 0)
	for _, r := range l.water {
		if r.Intersect(playerR) != zr {
			player.Drown(15)
			return
		}
	}
}

func (l *Level2) ReachedShop() string {
	for _, obj := range tmxMap.GetObjectLayerByName("shops").Objects {

		if r, err := obj.GetRect(); err == nil && r.Intersect(player.CollisionBox()) != pixel.R(0, 0, 0, 0) {
			return obj.Name
		}
	}

	return ""
}

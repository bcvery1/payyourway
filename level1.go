package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	lvl1gunPos = []pixel.Vec{
		pixel.V(3785, 3300),
		pixel.V(4260, 3604),
		pixel.V(3610, 3561),
		pixel.V(3610, 3689),
		pixel.V(3604, 3993),
		pixel.V(4046, 4138),
		pixel.V(3701, 4372),
		pixel.V(4261, 4596),
	}

	lvl1enemyPos = []pixel.Vec{
		pixel.V(3426, 1173),
		pixel.V(3313, 1283),
		pixel.V(3615, 1422),
		pixel.V(3353, 1479),
		pixel.V(3455, 1608),
		pixel.V(3099, 1545),
		pixel.V(3144, 1070),
		pixel.V(3295, 1003),
		pixel.V(3610, 963),
		pixel.V(3560, 1683),
		pixel.V(3946, 1177),
		pixel.V(3940, 1108),
		pixel.V(3933, 1000),
		pixel.V(5104, 1601),
	}
)

type Level1 struct {
	collisionRects []pixel.Rect
	mineRects      []pixel.Rect
	guns           []*Gun
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

	for _, gp := range lvl1gunPos {
		l.guns = append(l.guns, &Gun{
			pos:   gp,
			speed: 5,
		})
	}
}

func (l *Level1) Draw(win *pixelgl.Window) {
	_ = tmxMap.DrawAll(win, backingColour, pixel.IM)
	player.Draw(win)
}

func (l *Level1) Start() {
	camPos = pixel.V(3210, 405)

	for _, ep := range lvl1enemyPos {
		NewEnemy(true, ep)
	}
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

	for _, g := range l.guns {
		g.Update(dt)
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

package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level3 struct {
	collisionRects []pixel.Rect
	mineRects []pixel.Rect

}

func (l *Level3) Init(pixel.Rect) {
	collisionOjLayer := tmxMap.GetObjectLayerByName("Level3Collisions")
	for _, obj := range collisionOjLayer.Objects {
		if r, err :=  obj.GetRect(); err == nil {
			l.collisionRects = append(l.collisionRects, r)
		}
	}

	for _, obj := range tmxMap.GetObjectLayerByName("CommonCollisions").Objects {
		if r, err :=  obj.GetRect(); err == nil {
			l.collisionRects = append(l.collisionRects, r)
		}
	}

	mineObjLayer := tmxMap.GetObjectLayerByName("mines")
	for _, obj := range mineObjLayer.Objects {
		if r, err :=  obj.GetRect(); err == nil {
			l.mineRects = append(l.mineRects, r)
		}
	}
}

func (l *Level3) Start() {
	panic("implement me")
}

func (l *Level3) Update(float64, *pixelgl.Window) {
	panic("implement me")
}

func (l *Level3) Draw(*pixelgl.Window) {
	panic("implement me")
}

func (l *Level3) Collides(pixel.Rect) bool {
	panic("implement me")
}

func (l *Level3) Hurt(pixel.Rect) {
	panic("implement me")
}

func (l *Level3) ReachedShop() string {
	panic("implement me")
}


package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Level4 struct {
	collisionRects []pixel.Rect
	mineRects []pixel.Rect

}

func (l *Level4) Init(pixel.Rect) {
	collisionOjLayer := tmxMap.GetObjectLayerByName("Level4Collisions")
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

func (l *Level4) Start() {
	panic("implement me")
}

func (l *Level4) Update(float64, *pixelgl.Window) {
	panic("implement me")
}

func (l *Level4) Draw(*pixelgl.Window) {
	panic("implement me")
}

func (l *Level4) Collides(pixel.Rect) bool {
	panic("implement me")
}

func (l *Level4) Hurt(pixel.Rect) {
	panic("implement me")
}

func (l *Level4) ReachedShop() string {
	panic("implement me")
}


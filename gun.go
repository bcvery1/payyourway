package main

import (
	"github.com/faiface/pixel"
)

const (
	maxSafeDist = 48
	midSafeDist = 32
	minSafeDist = 16

	explodeDist = 8

	gunFireDist = 200
)

var (
	rockets     = make(map[int]*rocket)
	rocketCount = 0

	rocketAcc = .5

	rocketSprite *pixel.Sprite
)

type Gun struct {
	pos       pixel.Vec
	sinceLast float64
	speed     float64
}

func (g *Gun) Update(dt float64) {
	if g.sinceLast < g.speed {
		g.sinceLast += dt
	}

	if g.pos.To(player.CollisionBox().Center()).Len() > gunFireDist {
		return
	}

	if g.sinceLast > g.speed {
		g.Fire()
	}
}

func (g *Gun) Fire() {
	g.sinceLast = 0

	r := rocket{
		id:   rocketCount,
		pos:  g.pos,
		life: 5,
	}

	r.dir = r.toPlayer().Unit().Scaled(rocketAcc)

	rockets[rocketCount] = &r
	rocketCount++
}

type rocket struct {
	id   int
	pos  pixel.Vec
	dir  pixel.Vec
	life float64
}

func (r *rocket) Draw(target pixel.Target) {
	rocketSprite.Draw(target, pixel.IM.Moved(r.pos).Rotated(r.pos, r.dir.Angle()))
}

func (r *rocket) Update(dt float64) {
	r.life -= dt
	if r.life < 0 {
		r.explode()
		return
	}

	r.pos = r.pos.Add(r.dir)

	if r.toPlayer().Len() < explodeDist {
		r.explode()
	}

	r.dir = r.dir.Add(r.toPlayer().Unit().Scaled(rocketAcc))
}

func (r *rocket) toPlayer() pixel.Vec {
	return r.pos.To(player.CollisionBox().Center())
}

func (r *rocket) explode() {
	delete(rockets, r.id)

	if r.toPlayer().Len() < minSafeDist {
		player.Hurt(50)
	} else if r.toPlayer().Len() < midSafeDist {
		player.Hurt(30)
	} else if r.toPlayer().Len() < maxSafeDist {
		player.Hurt(15)
	}

	NewFire(r.pos, pixel.ZV)
}

func UpdateRockets(dt float64) {
	for _, r := range rockets {
		r.Update(dt)
	}
}

func DrawRockets(target pixel.Target) {
	for _, r := range rockets {
		r.Draw(target)
	}
}

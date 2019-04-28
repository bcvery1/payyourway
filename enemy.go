package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const (
	minFireSpeed = 0.3
	maxFireSpeed = 2.

	minBulletFireDist = 400

	bulletSpeed = 500
)

var (
	enemySprites []*pixel.Sprite

	enemies = make([]*Enemy, 0)

	bullets     = make(map[int]*bullet)
	bulletCount int
	bulletImd   = imdraw.New(nil)
)

type Enemy struct {
	pos       pixel.Vec
	static    bool
	angle     float64
	fireSpeed float64
	lastFire  float64
	spriteInd int
}

func NewEnemy(static bool, pos pixel.Vec) {
	e := Enemy{
		pos:       pos,
		static:    static,
		fireSpeed: (rand.Float64() * (maxFireSpeed - minFireSpeed)) + minFireSpeed,
		spriteInd: rand.Intn(len(enemySprites)),
	}

	enemies = append(enemies, &e)
}

func SetupEnemies() {
	for _, p := range tmxMap.GetObjectLayerByName("EnemyLocs").Objects {
		if p, err := p.GetPoint(); err == nil {
			NewEnemy(true, p)
		}
	}
}

func (e *Enemy) Update(dt float64) {
	toPlayer := e.pos.To(player.CollisionBox().Center())
	e.angle = toPlayer.Angle()

	e.lastFire += dt

	if e.lastFire > e.fireSpeed && toPlayer.Len() < minBulletFireDist {
		e.fire()
	}
}

func (e *Enemy) Draw(target pixel.Target) {
	enemySprites[e.spriteInd].Draw(target, pixel.IM.Moved(e.pos).Rotated(e.pos, e.angle))
}

func (e *Enemy) fire() {
	PlaySound(bulletSound)

	e.lastFire = 0

	NewBullet(e.pos, e.pos.To(player.CollisionBox().Center()))
}

func UpdateEnemies(dt float64) {
	for _, e := range enemies {
		e.Update(dt)
	}

	for _, b := range bullets {
		b.Update(dt)
	}
}

func DrawEnemies(target pixel.Target) {
	for _, e := range enemies {
		e.Draw(target)
	}

	bulletImd.Clear()
	for _, b := range bullets {
		b.Draw()
	}

	bulletImd.Draw(target)
}

type bullet struct {
	id  int
	pos pixel.Vec
	dir pixel.Vec
}

func NewBullet(pos, dir pixel.Vec) {
	b := bullet{
		id:  bulletCount,
		pos: pos,
		dir: dir,
	}

	bulletCount++

	bullets[b.id] = &b
}

func (b *bullet) Update(dt float64) {
	b.pos = b.pos.Add(b.dir.Unit().Scaled(bulletSpeed * dt))

	if b.pos.To(player.CollisionBox().Center()).Len() < 5 {
		player.Hurt(5)
		delete(bullets, b.id)
	}
}

func (b *bullet) Draw() {
	bulletImd.Push(b.pos)
	bulletImd.Circle(1, 0)
}

package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
)

type action string

const (
	health    action = "health"
	maxhealth action = "maxhealth"
	shield    action = "shield"

	powerupPickUp = 16.
)

var (
	powerups     = make(map[int]*powerup)
	powerupCount int

	powerupSprites map[action]*pixel.Sprite
)

func getAction(point *tilepix.Object) action {
	switch point.Type {
	case string(health):
		return health
	case string(maxhealth):
		return maxhealth
	case string(shield):
		return shield
	default:
		panic(point)
	}
}

func getVal(obj *tilepix.Object, act action) int {
	actStr := "health"
	if act == shield {
		actStr = "shield"
	}

	for _, p := range obj.Properties {
		if p.Name == actStr {
			i, err := strconv.Atoi(p.Value)
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	fmt.Println(obj)
	fmt.Println(string(act))
	panic("not found")
}

func SetupPowerups() {
	for _, obj := range tmxMap.GetObjectLayerByName("PowerUps").Objects {
		NewPowerUp(obj)
	}
}

type powerup struct {
	id     int
	pos    pixel.Vec
	action action
	value  int
}

func NewPowerUp(obj *tilepix.Object) {
	point, err := obj.GetPoint()
	if err != nil {
		panic(err)
	}

	p := powerup{
		id:     powerupCount,
		pos:    point,
		action: getAction(obj),
		value:  getVal(obj, getAction(obj)),
	}

	powerupCount++

	powerups[p.id] = &p
}

func (p *powerup) Draw(target pixel.Target) {
	sprite := powerupSprites[p.action]
	sprite.Draw(target, pixel.IM.Moved(p.pos))
}

func UpdatePowerups() {
	for _, p := range powerups {
		if player.CollisionBox().Center().To(p.pos).Len() < powerupPickUp {
			delete(powerups, p.id)

			switch p.action {
			case health:
				player.health = math.Min(float64(p.value)+player.health, player.maxHealth)
				Announce(fmt.Sprintf("%d hp", p.value))
			case maxhealth:
				player.maxHealth += float64(p.value)
				Announce(fmt.Sprintf("+%d max hp", p.value))
			case shield:
				player.shield += float64(p.value)
				player.maxShield += float64(p.value)
				Announce(fmt.Sprintf("+%d shield", p.value))
			}
		}
	}
}

func DrawPowerups(target pixel.Target) {
	for _, p := range powerups {
		p.Draw(target)
	}
}

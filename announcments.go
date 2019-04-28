package main

import (
	"fmt"
	"image/color"
	"sync/atomic"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

var (
	announcementsChan = make(chan string, 10)

	annoucing atomic.Value
	scale     float64
	alpha     uint8

	scaleRate = 1.5
	alphaRate = 2.

	t *text.Text
)

func SetupAnnouncements() {
	resetAnnouncemnts()

	t = text.New(pixel.ZV, atlas)
	t.Color = color.Black

	go listen()
}

func UpdateAnnouncements(dt float64) {
	if !annoucing.Load().(bool) {
		return
	}

	scale += scaleRate * dt
	alphaInt := int(alpha)
	alphaInt += int(alphaRate)
	if alphaInt > 255 {
		resetAnnouncemnts()
	} else {
		alpha = uint8(alphaInt)
	}
}

func resetAnnouncemnts() {
	scale = 1
	alpha = 0
	annoucing.Store(false)
}

func DrawAnnouncements(target pixel.Target) {
	if !annoucing.Load().(bool) {
		return
	}

	offset := player.offSet.Sub(t.Bounds().Center())
	mask := color.RGBA{R: alpha, G: alpha, B: alpha, A: 255 - alpha}
	t.DrawColorMask(target, pixel.IM.Moved(offset).Scaled(offset, scale), mask)
	//t.Draw(target, pixel.IM.Moved(offset).Scaled(offset, scale))
}

func Announce(text string) {
	announcementsChan <- text
}

func listen() {
	for {
		a := <-announcementsChan

		for annoucing.Load().(bool) {
			time.Sleep(time.Millisecond * 50)
		}

		annoucing.Store(true)
		t.Clear()
		_, _ = fmt.Fprint(t, a)
	}
}

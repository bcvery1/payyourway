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

	scaleRate = 2.5
	alphaRate = 2.

	t *text.Text
)

func SetupAnnouncements() {
	resetAnnouncemnts()

	t = text.New(pixel.ZV, atlas)
	t.Color = color.RGBA{R: 0x00, G: 0x2f, B: 0x5f, A: 0xff}

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

	offset := player.offSet.Sub(t.Bounds().Max)
	mask := color.RGBA{R: 255 - alpha, G: 255 - alpha, B: 255 - alpha, A: 0x00}
	t.DrawColorMask(target, pixel.IM.Moved(offset).Scaled(offset, scale), mask)
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

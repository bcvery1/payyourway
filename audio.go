package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var (
	mainTrackStreamer beep.Streamer

	sounds = make(map[Track]*beep.Buffer)

	hasErrored = false
)

type Track string

const (
	bulletSound    Track = "bullet"
	rocketSound    Track = "rocket"
	explosionSound Track = "explosion"
	hurtSound      Track = "hurt"
	drownSound     Track = "drown"
)

func SetupAudio() {
	mainTrack, err := os.Open("assets/mainTrack.wav")
	if err != nil {
		hasErrored = true
		return
	}

	streamer, format, err := wav.Decode(mainTrack)
	if err != nil {
		hasErrored = true
		return
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		hasErrored = true
		return
	}

	mainTrackStreamer = &effects.Volume{
		Streamer: beep.Loop(-1, streamer),
		Base:     2,
		Volume:   -3,
	}
	speaker.Play(mainTrackStreamer)

	loadSound(bulletSound)
	loadSound(rocketSound)
	loadSound(explosionSound)
	loadSound(hurtSound)
	loadSound(drownSound)

}

func loadSound(sound Track) {
	filename := filepath.Join("assets", fmt.Sprintf("%s.wav", sound))
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	_ = streamer.Close()

	sounds[sound] = buffer
}

func PlaySound(sound Track) {
	if hasErrored {
		return
	}

	b, ok := sounds[sound]
	if !ok {
		return
	}

	speaker.Play(b.Streamer(0, b.Len()))
}

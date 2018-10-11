package player_test

import (
	"testing"

	"github.com/code560/audigo/player"
)

func TestCreate(t *testing.T) {
	p := player.NewProxy()
	if p == nil {
		t.Error("dont created proxy: player.NewProxy()")
	}
}

func TestChan(t *testing.T) {
	p := player.NewProxy()
	p.Play <- player.NewPlayArgs(player.Src("bgm_wave.wav"))
	p.Volume <- &player.VolumeArgs{Vol: 1}
	p.Pause <- struct{}{}
	p.Resume <- struct{}{}
	p.Stop <- struct{}{}
}

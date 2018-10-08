package player

import (
	"testing"
	"time"
)

func TestSound(t *testing.T) {
	args := []string{
		"bgm_wave.wav",
		"se_jump.wav",
	}

	plist := make([]*Player, len(args))
	for i, arg := range args {
		p := newPlayerImpl()
		plist[i] = p

		go func(p *Player, name string) {
			p.Play(NewPlayArgs(
				Wave(name),
				Loop(true)))
		}(p, arg)
	}

	log.Debug("plaing sound")

	sec := time.Duration(2)
	time.Sleep(time.Second * sec)
	for _, p := range plist {
		log.Debug("call stop sound")
		p.Stop(nil)
	}

	log.Debug("done routines")
}

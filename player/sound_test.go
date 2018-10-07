package player

import (
	"os"
	"runtime/trace"
	"testing"
	"time"
)

func TestSound(t *testing.T) {
	args := []string{
		"bgm_wave.wav",
		"se_jump.wav",
	}

	// trace
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	// trace end

	plist := make([]*Player, len(args))
	for i, arg := range args {
		p := NewPlayer()
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
		p.Stop()
	}

	log.Debug("done routines")
}

package player

import (
	"sync"
	"testing"
	"time"

	"github.com/code560/audigo/util"
)

var l = util.GetLogger()

func TestSound(t *testing.T) {
	args := []string{
		"bgm_wave.wav",
		"se_jump.wav",
	}

	wg := &sync.WaitGroup{}
	plist := make([]*player, len(args))
	for i, arg := range args {
		p := newPlayerImpl()
		plist[i] = p
		wg.Add(1)
		go func(p *player, name string) {
			p.Play(NewPlayArgs(
				Src(name),
				Loop(true)))
			wg.Done()
		}(p, arg)
	}

	wg.Wait()
	l.Debug("plaing sound")

	sec := time.Duration(2)
	time.Sleep(time.Second * sec)
	for _, p := range plist {
		log.Debug("call stop sound")
		p.Stop(nil)
	}

	l.Debug("done routines")
}

func TestUnexpectedSound(t *testing.T) {
	loop := 10
	p := newPlayerImpl()
	for i := 0; i < loop; i++ {
		p.Stop(nil)
	}
	for i := 0; i < loop; i++ {
		p.Pause()
	}
	for i := 0; i < loop; i++ {
		p.Resume()
	}
	arg := 1.8
	for i := 0; i < loop; i++ {
		p.Volume(arg)
	}
	// for i := 0; i < loop; i++ {
	// 	p.Play(&PlayArgs{Src: "../asset/audio/bgm_wave.wav"})
	// 	fmt.Printf("%d ", i)
	// }
}

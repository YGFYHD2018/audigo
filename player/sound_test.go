package player_test

import (
	"math"
	"os"
	"runtime/trace"
	"testing"
	"time"

	"github.com/code560/audigo"
	"github.com/code560/audigo/util"
)

var l = util.GetLogger()

func TestSound(t *testing.T) {
	args := []string{
		"bgm_wave.wav",
		"se_jump.wav",
	}

	// trace
	f, err := os.Create("trace.out")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	// trace end

	plist := make([]audigo.Player, len(args))
	for _, arg := range args {
		p := audigo.NewPlayer()
		plist = append(plist, *p)
		go func(p *audigo.Player, name string) {
			p.Play(name, math.MaxInt64)
		}(p, arg)
	}

	l.Debug("started wav")
	time.Sleep(time.Second * 2)
	for _, p := range plist[0:] {
		l.Debug("call stop")
		p.Stop()
	}

	l.Debug("done routines")
}

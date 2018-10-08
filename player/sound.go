package player

import (
	"math"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"

	"github.com/code560/audigo/util"
)

var (
	mtx sync.Mutex
	log = util.GetLogger()
)

const (
	CH    = 2
	BPS   = 2
	CHUNK = 100

	VOLUME_BASE = 1.2
	VOLUME_INIT = 1
	// TODO マジックナンバー撲滅活動
)

type Player struct {
	ctrl  *beep.Ctrl
	vol   *effects.Volume
	mixer beep.Mixer

	samples [][CH]float64
	buf     []byte
	player  *oto.Player

	done *util.Closing
}

func newPlayerImpl() *Player {
	p := &Player{
		ctrl:  &beep.Ctrl{},
		vol:   &effects.Volume{Base: VOLUME_BASE, Volume: VOLUME_INIT},
		mixer: beep.Mixer{},

		done: util.NewClosing(),
	}
	return p
}

func (p *Player) Play(args *PlayArgs) {
	// open file
	f, err := os.Open(args.wav)
	if err != nil {
		log.Error(err)
		return
	}
	// decode file
	closer, format, err := wav.Decode(f)
	if err != nil {
		log.Error(err)
		return
	}
	defer closer.Close()

	// set middleware
	s := p.setMiddleware(closer, args)
	playing := make(chan struct{})
	s = beep.Seq(s, beep.Callback(func() {
		close(playing)
	}))
	// play sound
	p.setPlayer(format.SampleRate, format.SampleRate.N(time.Millisecond*CHUNK))
	p.mixer.Play(s)
	<-playing
	p.done.Reset()
}

func (p *Player) Stop(callback chan bool) {
	p.done.Close()
	if callback != nil {
		close(callback)
	}
}

func (p *Player) Volume(vol float64) {
	if vol == 0 {
		p.vol.Silent = true
		return
	} else {
		p.vol.Silent = false
	}
	p.vol.Volume = vol
}

func (p *Player) Pause() {
	p.ctrl.Paused = true
}

func (p *Player) Resume() {
	p.ctrl.Paused = false
}

func (p *Player) setMiddleware(closer beep.StreamSeekCloser, args *PlayArgs) beep.Streamer {
	s := beep.Loop(loopCount(args.loop), closer)
	p.ctrl.Streamer = s
	p.vol.Streamer = p.ctrl
	return p.vol
}

func (p *Player) setPlayer(sampleRate beep.SampleRate, bufferSize int) error {
	var err error
	bufferNum := bufferSize * CH * BPS
	mtx.Lock()
	p.player, err = oto.NewPlayer(int(sampleRate), CH, BPS, bufferNum)
	mtx.Unlock()
	if err != nil {
		return errors.Wrap(err, log.Error("failed to initialize oto.Player"))
	}
	p.samples = make([][CH]float64, bufferSize)
	p.buf = make([]byte, bufferNum)

	go func(done <-chan bool) {
		for {
			select {
			case <-done:
				log.Info("closing player")
				return
			default:
				p.sampling()
			}
		}
	}(p.done.GetDone())
	return nil
}

func (p *Player) sampling() {
	// mtx.Lock()
	p.mixer.Stream(p.samples)
	// mtx.Unlock()

	for s := range p.samples {
		for rl := range p.samples[s] {
			val := p.samples[s][rl]
			if val < -1 {
				val = -1
			}
			if val > +1 {
				val = +1
			}
			i16 := int16(val * (1<<15 - 1))
			l := byte(i16)
			h := byte(i16 >> 8)
			p.buf[s*4+rl*2+0] = l
			p.buf[s*4+rl*2+1] = h
		}
	}
	p.player.Write(p.buf)
}

func loopCount(enable bool) int {
	if enable {
		return math.MaxInt64
	} else {
		return 1
	}
}

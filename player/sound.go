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
	CH  = 2
	BPS = 2
	// TODO マジックナンバー撲滅活動
)

type Player struct {
	ctrl *beep.Ctrl
	vol  *effects.Volume

	mixer   beep.Mixer
	samples [][CH]float64
	buf     []byte
	player  *oto.Player

	done *util.Closing
}

func NewPlayer() *Player {
	p := &Player{}
	p.done = util.NewClosing()
	return p
}

func (p *Player) Play(args *PlayArgs) {
	// open file
	playing := make(chan struct{})
	f, err := os.Open(args.wav)
	if err != nil {
		log.Error(err)
		return
	}
	// decode file
	sc, format, err := wav.Decode(f)
	if err != nil {
		log.Error(err)
		return
	}
	defer sc.Close()

	// set middleware
	s := beep.Loop(loopCount(args.loop), sc)
	s = beep.Seq(s, beep.Callback(func() {
		close(playing)
	}))
	p.ctrl = &beep.Ctrl{Streamer: s}
	p.vol = &effects.Volume{Streamer: p.ctrl, Base: 1.2, Volume: 1}
	s = p.vol
	// play sound
	p.set(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.mixer = beep.Mixer{}
	p.mixer.Play(s)
	<-playing
}

func (p *Player) Stop() {
	p.done.Close()
}

func (p *Player) set(sampleRate beep.SampleRate, bufferSize int) error {
	var err error
	bufferNum := bufferSize * CH * BPS
	mtx.Lock()
	p.player, err = oto.NewPlayer(int(sampleRate), CH, BPS, bufferNum)
	mtx.Unlock()
	if err != nil {
		return errors.Wrap(err, log.Error("failed to initialize oto.Player"))
	}
	p.samples = make([][2]float64, bufferSize)
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

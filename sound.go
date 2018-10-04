package audigo

import (
	"context"
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
	dir string
	mtx sync.Mutex
	log = util.GetLogger()
)

func init() {
	dir = "assets/audio/"
}

type Player struct {
	ctrl beep.Ctrl
	vol  effects.Volume

	mixer   beep.Mixer
	samples [][2]float64
	buf     []byte
	player  *oto.Player

	ctx context.Context
}

func NewPlayer() *Player {
	p := new(Player)
	p.ctx = context.Background()
	return p
}

func (p *Player) Play(file string, loop int) {
	playing := make(chan struct{})
	f, err := os.Open(dir + file)
	if err != nil {
		log.Error(err)
		return
	}

	ssc, format, err := wav.Decode(f)
	if err != nil {
		log.Error(err)
		return
	}
	defer ssc.Close()

	l := beep.Loop(loop, ssc)
	s := beep.Seq(l, beep.Callback(func() {
		close(playing)
	}))
	p.ctrl = *(&beep.Ctrl{Streamer: s})
	p.vol = *(&effects.Volume{Streamer: &p.ctrl, Base: 1.2, Volume: 1})
	stm := &p.vol

	p.set(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.mixer = beep.Mixer{}
	p.mixer.Play(stm)
	<-playing
}

func (p *Player) Stop() {
	if isDone(p.ctx) {
		// 既にクローズしている
	} else {
		_, cancel := context.WithCancel(p.ctx)
		log.Debug("1")
		cancel()
		log.Debug("2")
		p.player.Close()
		log.Debug("3")
	}
}

func (p *Player) set(sampleRate beep.SampleRate, bufferSize int) error {
	bufferNum := bufferSize * 4

	var err error
	mtx.Lock()
	p.player, err = oto.NewPlayer(int(sampleRate), 2, 2, bufferNum)
	mtx.Unlock()
	if err != nil {
		return errors.Wrap(err, log.Error("failed to initialize oto.Player"))
	}
	p.samples = make([][2]float64, bufferSize)
	p.buf = make([]byte, bufferNum)

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				log.Info("closing player")
				return
			default:
				p.sampling()
			}
		}
	}()
	return nil
}

func (p *Player) sampling() {
	mtx.Lock()
	p.mixer.Stream(p.samples)
	mtx.Unlock()

	for i := range p.samples {
		for c := range p.samples[i] {
			val := p.samples[i][c]
			if val < -1 {
				val = -1
			}
			if val > +1 {
				val = +1
			}
			valInt16 := int16(val * (1<<15 - 1))
			low := byte(valInt16)
			high := byte(valInt16 >> 8)
			p.buf[i*4+c*2+0] = low
			p.buf[i*4+c*2+1] = high
		}
	}

	p.player.Write(p.buf)
}

// func ext(file string) string {
// 	pos := strings.LastIndex(file, ".")
// 	return file[pos:]
// }

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

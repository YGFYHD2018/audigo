package player

import (
	"math"
	"os"
	"sync"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"

	"github.com/code560/audigo/util"
)

// statics

var (
	mtx sync.Mutex
	log = util.GetLogger()
)

// consts

const (
	CH    = 2
	BPS   = 2
	CHUNK = 100

	VOLUME_BASE = 1.2
	VOLUME_INIT = 1
	// TODO マジックナンバー撲滅活動
)

// functions

func loopCount(enable bool) int {
	if enable {
		return math.MaxInt64
	} else {
		return 1
	}
}

// create

type playerMaker struct {
	close *util.Closing

	ctrl    *beep.Ctrl
	vol     *effects.Volume
	mixer   *beep.Mixer
	oto     *oto.Player
	samples [][CH]float64
	buf     []byte
}

// interface methods

func (p *playerMaker) Stop(callback chan<- struct{}) {
	if p.close == nil {
		return
	}
	p.close.Reset()
	if callback != nil {
		close(callback)
	}
}

func (p *playerMaker) Volume(args *VolumeArgs) {
	p.volume(args.Vol)
}

func (p *playerMaker) volume(vol float64) {
	if vol == 0 {
		p.vol.Silent = true
	} else {
		p.vol.Silent = false
	}
	p.vol.Volume = vol
}

func (p *playerMaker) Pause() {
	p.ctrl.Paused = true
}

func (p *playerMaker) Resume() {
	p.ctrl.Paused = false
}

// maker methods

func (p *playerMaker) makeMixer() *beep.Mixer {
	return new(beep.Mixer)
}

func (p *playerMaker) makeCtrl() *beep.Ctrl {
	return new(beep.Ctrl)
}

func (p *playerMaker) makeVolume() *effects.Volume {
	return &effects.Volume{Base: VOLUME_BASE, Volume: VOLUME_INIT}
}

func (p *playerMaker) makeClosing() *util.Closing {
	return util.NewClosing()
}

func (p *playerMaker) makeOtoPlayer(sampleRate beep.SampleRate, bufferSize int) error {
	var err error
	bufferNum := bufferSize * CH * BPS
	mtx.Lock()
	p.oto, err = oto.NewPlayer(int(sampleRate), CH, BPS, bufferNum)
	mtx.Unlock()
	if err != nil {
		return errors.Wrap(err, log.Error("failed to initialize oto.Player"))
	}
	p.samples = make([][CH]float64, bufferSize)
	p.buf = make([]byte, bufferNum)

	go func(close <-chan struct{}) {
		for {
			select {
			case <-close:
				log.Info("closing player")
				return
			default:
				p.sampling()
			}
		}
	}(p.close.GetDone())
	return nil
}

func (p *playerMaker) sampling() {
	mtx.Lock()
	p.mixer.Stream(p.samples)
	mtx.Unlock()

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
	p.oto.Write(p.buf)
}

func (p *playerMaker) setCtrlStream(s beep.Streamer) beep.Streamer {
	if p.ctrl == nil {
		p.ctrl = p.makeCtrl()
	}
	p.ctrl.Streamer = s
	return p.ctrl
}

func (p *playerMaker) setVolumeStream(s beep.Streamer) beep.Streamer {
	if p.vol == nil {
		p.vol = p.makeVolume()
	}
	p.vol.Streamer = s
	return p.vol
}

func (p *playerMaker) openFile(src string) (beep.StreamSeekCloser, *beep.Format) {

	// open file
	f, err := os.Open(src)
	if err != nil {
		log.Errorf("dont open file: %s\n%s", src, err.Error())
		return nil, nil
	}
	// decode file
	closer, format, err := wav.Decode(f)
	if err != nil {
		log.Error("dont decode file: %s\n%s", src, err.Error())
		return nil, nil
	}

	return closer, &format
}

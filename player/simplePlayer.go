package player

import (
	"os"
	"time"

	"github.com/faiface/beep"
)

type simplePlayer struct {
	playerMaker
}

func newSimplePlayer() Player {
	p := new(simplePlayer)
	p.close = p.makeClosing()
	return p
}

func (p *simplePlayer) Play(args *PlayArgs) {
	// init
	p.reset()
	// and stop
	if args.Stop {
		p.Pause()
	}
	// open file
	if _, err := os.Stat(args.Src); err != nil {
		log.Warn("not found music file: %s", args.Src)
		return
	}
	closer, format := p.openFile(args.Src)
	defer closer.Close()
	// set middleware
	s := beep.Loop(loopCount(args.Loop), closer)
	s = p.setCtrlStream(s)
	s = p.setVolumeStream(s)
	playing := make(chan struct{})
	s = beep.Seq(s, beep.Callback(func() {
		p.Stop(nil)
		close(playing)
	}))
	// play sound
	p.makeOtoPlayer(format.SampleRate, format.SampleRate.N(time.Millisecond*CHUNK))
	p.mixer = p.makeMixer()
	p.mixer.Play(s)
	<-playing
}

func (p *simplePlayer) Volume(args *VolumeArgs) {
	if p.vol == nil {
		p.vol = p.makeVolume()
	}
	p.volume(args.Vol)
}

func (p *simplePlayer) Pause() {
	if p.ctrl == nil {
		p.ctrl = p.makeCtrl()
	}
	p.ctrl.Paused = true
}

func (p *simplePlayer) Resume() {
	if p.ctrl == nil {
		p.ctrl = p.makeCtrl()
	}
	p.ctrl.Paused = false
}

func (p *simplePlayer) reset() {
	p.close = p.makeClosing()

	p.streaMutex.Lock()
	p.ctrl = nil
	p.vol = nil
	p.mixer = nil
	p.oto = nil
	p.samples = nil
	p.buf = nil
	p.streaMutex.Unlock()
}

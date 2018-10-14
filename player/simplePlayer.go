package player

import (
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
	if args.Stop {
		p.Pause()
	}

	closer, format := p.openFile(args.Src)
	defer closer.Close()

	// set middleware
	s := beep.Loop(loopCount(args.Loop), closer)
	s = p.setCtrlStream(s)
	s = p.setVolumeStream(s)
	playing := make(chan struct{})
	s = beep.Seq(s, beep.Callback(func() {
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

package player

import (
	"time"

	"github.com/faiface/beep"
)

type rePlayer struct {
	playerMaker
}

func newRePlayer() Player {
	p := new(rePlayer)
	p.ctrl = p.makeCtrl()
	p.vol = p.makeVolume()
	p.close = p.makeClosing()
	return p
}

func (p *rePlayer) Play(args *PlayArgs) {
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

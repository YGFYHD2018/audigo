package player

import "github.com/code560/audigo/util"

const (
	ChanSize = 10
)

// Proxy は、sound player proxyです。
type Proxy struct {
	Play   chan *PlayArgs
	Stop   chan struct{}
	Volume chan *VolumeArgs
	Pause  chan struct{}
	Resume chan struct{}

	closing chan struct{}

	sp *player
}

// NewProxy は、Playerを生成して返します。
func newProxyImpl() *Proxy {
	return &Proxy{
		Play:   make(chan *PlayArgs, ChanSize),
		Stop:   make(chan struct{}, ChanSize),
		Volume: make(chan *VolumeArgs, ChanSize),
		Pause:  make(chan struct{}, ChanSize),
		Resume: make(chan struct{}, ChanSize),

		closing: make(chan struct{}),
	}
}

func (p *Proxy) setPlayer(player *player) {
	if p.sp != nil && util.IsDone(p.sp.done.GetDone()) {
		p.sp.Stop(nil)
	}
	p.sp = player
	p.rechan()
}

func (p *Proxy) rechan() {
	// stop and restart worker
	close(p.closing)
	p.closing = make(chan struct{})
	if p.sp == nil {
		return
	}
	go p.work()
}

func (p *Proxy) work() {
	for {
		select {
		case a := <-p.Play:
			log.Debug("call chan Proxy.Play")
			if isDone(p.closing) {
				goto END
			}
			p.sp.Play(a)
		case <-p.Stop:
			log.Debug("call chan Proxy.Stop")
			if isDone(p.closing) {
				goto END
			}
			p.sp.Stop(nil)
		case a := <-p.Volume:
			log.Debug("call chan Proxy.Volume")
			if isDone(p.closing) {
				goto END
			}
			p.sp.Volume(a.Vol)
		case <-p.Pause:
			log.Debug("call chan Proxy.Pause")
			if isDone(p.closing) {
				goto END
			}
			p.sp.Pause()
		case <-p.Resume:
			log.Debug("call chan Proxy.Resume")
			if isDone(p.closing) {
				goto END
			}
			p.sp.Resume()
		}
	}
END:
}

func isDone(c chan struct{}) bool {
	return util.IsDone(c)
}

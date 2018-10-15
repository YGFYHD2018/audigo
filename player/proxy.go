package player

import (
	"sync"

	"github.com/code560/audigo/util"
)

const (
	ChanSize = 20
)

// Proxy は、sound player proxyです。
type Proxy interface {
	GetChannel() chan<- *Action
}

type Action struct {
	Act  Actions
	Args interface{}
}

type Actions int

const (
	_ Actions = iota
	Play
	Stop
	Volume
	Pause
	Resume
)

type simpleProxy struct {
	playerPool *sync.Pool
	act        chan *Action
	closing    chan struct{}
}

// NewProxy は、Playerを生成して返します。
func newSimpleProxy() Proxy {
	return &simpleProxy{
		act:     make(chan *Action, ChanSize),
		closing: make(chan struct{}),
	}
}

func (p *simpleProxy) GetChannel() chan<- *Action {
	return p.act
}

// func (p *simpleProxy) getPlayer() Player {
// 	return p.playerPool.Get().(Player)
// }

// func (p *simpleProxy) setPlayer(player Player) {
// 	if p.sp != nil {
// 		p.sp.Stop(nil)
// 	}
// 	p.sp = player
// 	p.rechan()
// }

// func (p *simpleProxy) rechan() {
// 	// stop and restart worker
// 	close(p.closing)
// 	p.closing = make(chan struct{})
// 	if p.sp == nil {
// 		return
// 	}
// 	go p.work()
// }

func (p *simpleProxy) work() {
	for {
		select {
		case v := <-p.act:
			if isDone(p.closing) {
				return
			}
			p.call(v)
		}
	}
}

func (p *simpleProxy) call(arg *Action) {
	switch arg.Act {
	case Play:
		log.Debug("call chan Proxy.Play")
		a := arg.Args.(*PlayArgs)
		a.Src = dir + a.Src
		// go p.sp.Play(a)
		go func(pool *sync.Pool, a *PlayArgs) {
			p := pool.Get().(Player)
			p.Play(a)
			pool.Put(p)
		}(p.playerPool, a)
	case Stop:
		log.Debug("call chan Proxy.Stop")
		// go p.sp.Stop(nil)
	case Pause:
		log.Debug("call chan Proxy.Pause")
		// go p.sp.Pause()
	case Resume:
		log.Debug("call chan Proxy.Resume")
		// go p.sp.Resume()
	case Volume:
		log.Debug("call chan Proxy.Volume")
		a := arg.Args.(*VolumeArgs)
		// go p.sp.Volume(a)
	default:
		log.Warn("nothing call player function")
	}
}

func isDone(c chan struct{}) bool {
	return util.IsDone(c)
}

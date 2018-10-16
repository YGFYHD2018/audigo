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

	plays   map[int]Player
	playMtx sync.Mutex
}

// NewProxy は、Playerを生成して返します。
func newSimpleProxy() Proxy {
	p := &simpleProxy{
		act:     make(chan *Action, ChanSize),
		closing: make(chan struct{}),

		plays: make(map[int]Player, 32),
	}
	go p.work()
	return p
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
		go func(p *simpleProxy, a *PlayArgs) {
			player := p.playerPool.Get().(Player)
			var i int
			p.playLock(func() {
				i = p.pushPlayer(player)
			})
			player.Play(a)
			p.playLock(func() {
				p.popPlayer(i)
			})
			p.playerPool.Put(player)
		}(p, a)
	case Stop:
		log.Debug("call chan Proxy.Stop")
		p.playLock(func() {
			for _, player := range p.plays {
				go player.Stop(nil)
			}
		})
	case Pause:
		log.Debug("call chan Proxy.Pause")
		p.playLock(func() {
			for _, player := range p.plays {
				go player.Pause()
			}
		})
	case Resume:
		log.Debug("call chan Proxy.Resume")
		p.playLock(func() {
			for _, player := range p.plays {
				go player.Resume()
			}
		})
	case Volume:
		log.Debug("call chan Proxy.Volume")
		a := arg.Args.(*VolumeArgs)
		p.playLock(func() {
			for _, player := range p.plays {
				go player.Volume(a)
			}
		})
	default:
		log.Warn("nothing call player function")
	}
}

func (p *simpleProxy) playLock(f func()) {
	p.playMtx.Lock()
	f()
	p.playMtx.Unlock()
}

func (p *simpleProxy) atPlayer(i int) Player {
	v, ok := p.plays[i]
	if ok {
		return v
	} else {
		return nil
	}
}

func (p *simpleProxy) pushPlayer(player Player) int {
	i := len(p.plays)
	p.plays[i] = player
	return i
}

func (p *simpleProxy) popPlayer(i int) {
	_, ok := p.plays[i]
	if ok {
		delete(p.plays, i)
	}
}

func isDone(c chan struct{}) bool {
	return util.IsDone(c)
}

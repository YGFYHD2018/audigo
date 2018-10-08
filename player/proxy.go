package player

import (
	"sync"
)

// Proxy は、sound player proxyです。
type Proxy struct {
	Play   chan *PlayArgs
	Stop   chan struct{}
	Volume chan *VolumeArgs
	Pause  chan struct{}
	Resume chan struct{}

	closing chan struct{}

	players []*Player
}

// NewProxy は、Playerを生成して返します。
func newProxyImpl() *Proxy {
	return &Proxy{
		Play:   make(chan *PlayArgs),
		Stop:   make(chan struct{}),
		Volume: make(chan *VolumeArgs),
		Pause:  make(chan struct{}),
		Resume: make(chan struct{}),

		closing: make(chan struct{}),
	}
}

func (proxy *Proxy) setPlayer(players []*Player) {
	proxy.stopall()
	proxy.players = players
	proxy.rechan()
}

func (proxy *Proxy) rechan() {
	// stop preset
	close(proxy.closing)
	proxy.closing = make(chan struct{})

	// nothing
	if proxy.players != nil {
		return
	}
	// reworker
	go proxy.work()
}

func (proxy *Proxy) stopall() {
	wg := sync.WaitGroup{}
	for _, p := range proxy.players {
		wg.Add(1)
		go func(p *Player) {
			p.Stop(nil)
			// TODO 終了待ち
			wg.Done()
		}(p)
	}

	wg.Wait()
}

func (proxy *Proxy) work() {
	for range proxy.closing {
		select {
		case a := <-proxy.Play:
			for _, p := range proxy.players {
				p.Play(a)
			}
		case <-proxy.Stop:
			for _, p := range proxy.players {
				p.Stop(nil)
			}
		case a := <-proxy.Volume:
			for _, p := range proxy.players {
				p.Volume(a.val)
			}
		case <-proxy.Pause:
			for _, p := range proxy.players {
				p.Pause()
			}
		case <-proxy.Resume:
			for _, p := range proxy.players {
				p.Resume()
			}
		}
	}
}

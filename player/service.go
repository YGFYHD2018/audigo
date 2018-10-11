package player

import (
	"sync"
)

var (
	playerPool sync.Pool
	proxyPool  sync.Pool
)

func NewProxy() *Proxy {
	proxy := proxyPool.Get().(*Proxy)
	player := playerPool.Get().(*player)
	proxy.setPlayer(player)
	return proxy
}

func CloseProxy(p *Proxy) {
	if p.sp != nil {
		playerPool.Put(p.sp)
	}
	proxyPool.Put(p)
}

func init() {
	initPool()
}

func initPool() {
	playerPool = sync.Pool{
		New: func() interface{} {
			return newPlayerImpl() // todo trans internal
		},
	}

	proxyPool = sync.Pool{
		New: func() interface{} {
			return newProxyImpl() // todo trans internal
		},
	}
}

package player

import (
	"sync"
)

var (
	playerPool sync.Pool
	proxyPool  sync.Pool
)

func init() {
	initPool()
}

func NewProxy() Proxy {
	proxy := proxyPool.Get().(*simpleProxy)
	player := playerPool.Get().(Player)
	proxy.setPlayer(player)
	return proxy
}

func CloseProxy(p Proxy) {
	p_ := p.(*simpleProxy)
	if p_.sp != nil {
		playerPool.Put(p_.sp)
	}
	proxyPool.Put(p)
}

func initPool() {
	playerPool = sync.Pool{
		New: func() interface{} {
			return newSimplePlayer() // todo trans internal
		},
	}

	proxyPool = sync.Pool{
		New: func() interface{} {
			return newSimpleProxy() // todo trans internal
		},
	}
}
